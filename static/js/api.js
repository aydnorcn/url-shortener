/**
 * GoShort API Client
 * Wraps all REST API endpoints provided by the Gin backend.
 */
class ApiClient {
  constructor() {
    // Determine Base URL
    if (window.location.protocol.startsWith('http')) {
      this.baseUrl = window.location.origin;
    } else {
      this.baseUrl = 'http://localhost:8080';
    }
    this.tokenKey = 'goshort_jwt_token';
    this.userKey = 'goshort_user_data';
  }

  // Token & User Storage Helpers
  getToken() {
    return localStorage.getItem(this.tokenKey);
  }

  setToken(token) {
    if (token) {
      localStorage.setItem(this.tokenKey, token);
    } else {
      localStorage.removeItem(this.tokenKey);
    }
  }

  getUser() {
    const raw = localStorage.getItem(this.userKey);
    try {
      return raw ? JSON.parse(raw) : null;
    } catch {
      return null;
    }
  }

  setUser(user) {
    if (user) {
      localStorage.setItem(this.userKey, JSON.stringify(user));
    } else {
      localStorage.removeItem(this.userKey);
    }
  }

  logout() {
    this.setToken(null);
    this.setUser(null);
  }

  isAuthenticated() {
    return !!this.getToken();
  }

  // Core HTTP Request Handler
  async request(endpoint, options = {}) {
    const url = endpoint.startsWith('http') ? endpoint : `${this.baseUrl}${endpoint}`;
    
    const headers = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...(options.headers || {})
    };

    if (options.useAuth !== false) {
      const token = this.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    const config = {
      method: options.method || 'GET',
      headers,
      ...(options.body ? { body: typeof options.body === 'string' ? options.body : JSON.stringify(options.body) } : {})
    };

    const startTime = performance.now();
    try {
      const response = await fetch(url, config);
      const endTime = performance.now();
      const latency = Math.round(endTime - startTime);

      let data = null;
      const text = await response.text();
      try {
        data = text ? JSON.parse(text) : {};
      } catch {
        data = { raw: text };
      }

      if (!response.ok) {
        const errorMsg = data?.message || data?.error || `HTTP ${response.status}: ${response.statusText}`;
        const error = new Error(errorMsg);
        error.status = response.status;
        error.data = data;
        error.latency = latency;
        throw error;
      }

      return {
        status: response.status,
        statusText: response.statusText,
        data,
        latency
      };
    } catch (err) {
      if (!err.latency) {
        err.latency = Math.round(performance.now() - startTime);
      }
      throw err;
    }
  }

  // Health Endpoint
  async getHealth() {
    return this.request('/api/health', { useAuth: false });
  }

  // Auth Endpoints
  async register(email, password, confirmPassword) {
    return this.request('/api/v1/auth/register', {
      method: 'POST',
      body: {
        email,
        password,
        confirm_password: confirmPassword
      },
      useAuth: false
    });
  }

  async login(email, password) {
    const res = await this.request('/api/v1/auth/login', {
      method: 'POST',
      body: { email, password },
      useAuth: false
    });

    if (res.data?.token) {
      this.setToken(res.data.token);
      this.setUser(res.data.user || { email });
    }
    return res;
  }

  // URL Management Endpoints
  async createUrl({ originalUrl, customAlias, expiresAt }) {
    const payload = {
      original_url: originalUrl
    };

    if (customAlias && customAlias.trim() !== '') {
      payload.custom_alias = customAlias.trim();
    }

    if (expiresAt) {
      // Format to ISO 8601 string
      const dateObj = new Date(expiresAt);
      if (!isNaN(dateObj.getTime())) {
        payload.expires_at = dateObj.toISOString();
      }
    }

    return this.request('/api/v1/urls', {
      method: 'POST',
      body: payload
    });
  }

  async getUrls(page = 1, pageSize = 10) {
    return this.request(`/api/v1/urls?page=${page}&pageSize=${pageSize}`);
  }

  async getUrl(id) {
    return this.request(`/api/v1/urls/${id}`);
  }

  async updateUrl(id, { originalUrl, expiresAt, isActive }) {
    const payload = {};
    if (originalUrl !== undefined) payload.original_url = originalUrl;
    if (isActive !== undefined) payload.is_active = isActive;
    if (expiresAt !== undefined) {
      if (expiresAt === null || expiresAt === '') {
        payload.expires_at = null;
      } else {
        const dateObj = new Date(expiresAt);
        payload.expires_at = !isNaN(dateObj.getTime()) ? dateObj.toISOString() : null;
      }
    }

    return this.request(`/api/v1/urls/${id}`, {
      method: 'PATCH',
      body: payload
    });
  }

  async deleteUrl(id) {
    return this.request(`/api/v1/urls/${id}`, {
      method: 'DELETE'
    });
  }

  async activateUrl(id) {
    return this.request(`/api/v1/urls/${id}/activate`, {
      method: 'PATCH'
    });
  }

  async deactivateUrl(id) {
    return this.request(`/api/v1/urls/${id}/deactivate`, {
      method: 'PATCH'
    });
  }

  // Analytics Endpoint
  async getAnalytics(id) {
    return this.request(`/api/v1/urls/${id}/analytics`);
  }
}

// Global API instance
window.api = new ApiClient();
