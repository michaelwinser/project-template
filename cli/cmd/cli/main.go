package main

import (
	"fmt"
	"os"
	"strings"

	"project-template/cli/internal/client"
	"project-template/cli/internal/config"
)

const version = "0.1.0"

func main() {
	// Parse global flags and extract command
	configPath, args := parseGlobalFlags(os.Args[1:])

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	// Load config from specified path or default
	var cfg *config.Config
	var err error
	if configPath != "" {
		cfg, err = config.LoadFromPath(configPath)
	} else {
		cfg, err = config.Load()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case "login":
		cmdLogin(cfg)
	case "logout":
		cmdLogout(cfg)
	case "me", "whoami":
		cmdMe(cfg)
	case "health":
		cmdHealth(cfg)
	case "config":
		cmdConfig(cfg)
	case "version":
		fmt.Printf("project-template-cli v%s\n", version)
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// parseGlobalFlags extracts --config flag and returns remaining args
func parseGlobalFlags(args []string) (configPath string, remaining []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Handle --config=path or --config path
		if arg == "--config" || arg == "-c" {
			if i+1 < len(args) {
				configPath = args[i+1]
				i++ // Skip the next arg
				continue
			}
			fmt.Fprintf(os.Stderr, "Error: %s requires a path argument\n", arg)
			os.Exit(1)
		}
		if strings.HasPrefix(arg, "--config=") {
			configPath = strings.TrimPrefix(arg, "--config=")
			continue
		}
		if strings.HasPrefix(arg, "-c=") {
			configPath = strings.TrimPrefix(arg, "-c=")
			continue
		}

		remaining = append(remaining, arg)
	}
	return
}

func printUsage() {
	fmt.Println(`project-template-cli - CLI for project-template API

Usage:
  project-template-cli [options] <command>

Options:
  --config, -c <path>  Path to config file (default: ~/.project-template-cli.json)

Commands:
  login     Start the login process
  logout    Log out and clear session
  me        Show current user info (alias: whoami)
  health    Check server health
  config    Show current configuration
  version   Show CLI version
  help      Show this help message

Configuration:
  By default, config is stored in ~/.project-template-cli.json
  Session data is stored alongside the config file.

  Use --config to specify a different config file for isolated contexts:
    project-template-cli --config ./myproject.json login
    project-template-cli -c /path/to/config.json health`)
}

func cmdLogin(cfg *config.Config) {
	apiClient := client.NewClient(cfg.ServerURL)

	// In a real CLI, we'd open a browser and handle the OAuth callback
	// For now, we provide instructions and accept a manual cookie input
	fmt.Println("To log in:")
	fmt.Printf("1. Open this URL in your browser: %s\n", apiClient.GetLoginURL())
	fmt.Println("2. Complete the OAuth flow")
	fmt.Println("3. After successful login, copy the 'session' cookie value from your browser")
	fmt.Println()

	// For development mode, we can use a simple flow
	fmt.Print("Enter session cookie value (or 'dev' for development mode): ")
	var cookie string
	fmt.Scanln(&cookie)

	if cookie == "" {
		fmt.Println("No cookie provided, aborting")
		os.Exit(1)
	}

	// Test the cookie by getting user info
	apiClient.SetCookie(cookie)
	user, err := apiClient.GetCurrentUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Save the session
	session := &config.Session{
		Cookie: cookie,
		UserID: user.ID,
		Email:  user.Email,
	}
	if err := cfg.SaveSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Logged in as %s (%s)\n", user.Name, user.Email)
}

func cmdLogout(cfg *config.Config) {
	session, err := cfg.LoadSession()
	if err != nil {
		fmt.Println("Not logged in")
		return
	}

	apiClient := client.NewClient(cfg.ServerURL)
	apiClient.SetCookie(session.Cookie)

	// Call server logout
	_, err = apiClient.Logout()
	if err != nil {
		// Even if server logout fails, clear local session
		fmt.Fprintf(os.Stderr, "Warning: server logout failed: %v\n", err)
	}

	// Clear local session
	if err := cfg.ClearSession(); err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing session: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Logged out successfully")
}

func cmdMe(cfg *config.Config) {
	session, err := cfg.LoadSession()
	if err != nil {
		fmt.Println("Not logged in. Run 'project-template-cli login' first.")
		os.Exit(1)
	}

	apiClient := client.NewClient(cfg.ServerURL)
	apiClient.SetCookie(session.Cookie)

	user, err := apiClient.GetCurrentUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ID:    %s\n", user.ID)
	fmt.Printf("Email: %s\n", user.Email)
	if user.Name != "" {
		fmt.Printf("Name:  %s\n", user.Name)
	}
}

func cmdHealth(cfg *config.Config) {
	apiClient := client.NewClient(cfg.ServerURL)

	health, err := apiClient.GetHealth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Status:    %s\n", health.Status)
	fmt.Printf("Timestamp: %s\n", health.Timestamp)
	if health.Version != "" {
		fmt.Printf("Version:   %s\n", health.Version)
	}
}

func cmdConfig(cfg *config.Config) {
	fmt.Printf("Config file:  %s\n", cfg.ConfigPath())
	fmt.Printf("Server URL:   %s\n", cfg.ServerURL)
	fmt.Printf("Session file: %s\n", cfg.SessionFile)

	session, err := cfg.LoadSession()
	if err == nil {
		fmt.Printf("Logged in as: %s\n", session.Email)
	} else {
		fmt.Println("Not logged in")
	}
}
