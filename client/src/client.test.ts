import { ApiClient, ApiError } from './index';

// Mock fetch
const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('ApiClient', () => {
  let client: ApiClient;

  beforeEach(() => {
    client = new ApiClient({ baseUrl: 'http://localhost:8080' });
    mockFetch.mockReset();
  });

  describe('getHealth', () => {
    it('should return health status', async () => {
      const mockResponse = {
        status: 'healthy',
        timestamp: '2024-01-15T12:00:00Z',
        version: '0.1.0',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await client.getHealth();

      expect(result).toEqual(mockResponse);
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/health',
        expect.objectContaining({
          method: 'GET',
          credentials: 'include',
        })
      );
    });
  });

  describe('getCurrentUser', () => {
    it('should return the current user', async () => {
      const mockUser = {
        id: 'user-123',
        email: 'test@example.com',
        name: 'Test User',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      });

      const result = await client.getCurrentUser();

      expect(result).toEqual(mockUser);
    });

    it('should throw ApiError when not authenticated', async () => {
      const mockErrorResponse = {
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({
          error: 'not_authenticated',
          message: 'Not authenticated',
        }),
      };

      mockFetch.mockResolvedValueOnce(mockErrorResponse);
      await expect(client.getCurrentUser()).rejects.toThrow(ApiError);

      mockFetch.mockResolvedValueOnce(mockErrorResponse);
      await expect(client.getCurrentUser()).rejects.toMatchObject({
        statusCode: 401,
        errorCode: 'not_authenticated',
      });
    });
  });

  describe('isAuthenticated', () => {
    it('should return true when authenticated', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({ id: 'user-123', email: 'test@example.com' }),
      });

      const result = await client.isAuthenticated();

      expect(result).toBe(true);
    });

    it('should return false when not authenticated', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({
          error: 'not_authenticated',
          message: 'Not authenticated',
        }),
      });

      const result = await client.isAuthenticated();

      expect(result).toBe(false);
    });
  });

  describe('logout', () => {
    it('should log out successfully', async () => {
      const mockResponse = {
        success: true,
        message: 'Successfully logged out',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      });

      const result = await client.logout();

      expect(result).toEqual(mockResponse);
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/auth/logout',
        expect.objectContaining({
          method: 'POST',
        })
      );
    });
  });

  describe('getLoginUrl', () => {
    it('should return the login URL', () => {
      const url = client.getLoginUrl();

      expect(url).toBe('http://localhost:8080/auth/login');
    });
  });
});
