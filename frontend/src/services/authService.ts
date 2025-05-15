import { environment } from '../environments/environment';

interface User {
  id: number;
  name: string;
  email: string;
  picture?: string;
  token: string;
  role: string; // "app_admin", "stakeholder", or "user"
  boardRoles?: Record<number, string>; // Board ID -> Role mappings
}

class AuthService {
  private currentUser: User | null = null;

  constructor() {
    // Load user from localStorage if available
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      try {
        this.currentUser = JSON.parse(savedUser);
      } catch (e) {
        // Invalid user in localStorage
        localStorage.removeItem('user');
      }
    }
  }

  // Get the Google sign-in URL
  async getGoogleSignInUrl(): Promise<string> {
    const response = await fetch(`${environment.apiUrl}/auth/google/login`);
    const data = await response.json();
    return data.url;
  }
  // Process the token received from Google OAuth
  async handleGoogleCallback(code: string): Promise<User> {
    const response = await fetch(`${environment.apiUrl}/auth/google/callback?code=${code}`);
    if (!response.ok) {
      throw new Error('Authentication failed');
    }
    
    const userData = await response.json();
    const user: User = {
      id: userData.userId,
      name: userData.name,
      email: userData.email,
      picture: userData.picture,
      token: userData.token,
      role: userData.role || 'user', // Default to 'user' if not provided
      boardRoles: userData.boardRoles || {}
    };
    
    // Store user in localStorage and memory
    localStorage.setItem('user', JSON.stringify(user));
    this.currentUser = user;
    
    return user;
  }

  // Get current authenticated user
  getCurrentUser(): User | null {
    return this.currentUser;
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return this.currentUser !== null;
  }

  // Get authentication token
  getToken(): string | null {
    return this.currentUser?.token || null;
  }

  // Logout user
  logout(): void {
    localStorage.removeItem('user');
    this.currentUser = null;
  }
  
  // Get auth header for API requests
  getAuthHeader(): HeadersInit {
    const token = this.getToken();
    return token ? { 'Authorization': `Bearer ${token}` } : {};
  }
  
  // Get user ID
  getUserId(): number | null {
    return this.currentUser?.id || null;
  }
  
  // Get user's global role
  getUserRole(): string | null {
    return this.currentUser?.role || null;
  }
  
  // Check if user is app admin
  isAppAdmin(): boolean {
    return this.currentUser?.role === 'app_admin';
  }
  
  // Check if user is stakeholder (globally)
  isStakeholder(): boolean {
    return this.currentUser?.role === 'stakeholder';
  }
  
  // Get user's role for a specific board
  getBoardRole(boardId: number): string | null {
    return this.currentUser?.boardRoles?.[boardId] || null;
  }
  
  // Check if user is a stakeholder for a specific board
  isBoardStakeholder(boardId: number): boolean {
    // App admins have stakeholder privileges on all boards
    if (this.isAppAdmin()) return true;
    
    return this.currentUser?.boardRoles?.[boardId] === 'stakeholder';
  }
  
  // Check if user has access to a specific board
  hasBoardAccess(boardId: number): boolean {
    // App admins have access to all boards
    if (this.isAppAdmin()) return true;
    
    return !!this.currentUser?.boardRoles?.[boardId];
  }
  
  // Refresh user profile and roles
  async refreshUserProfile(): Promise<User> {
    const response = await fetch(`${environment.apiUrl}/auth/profile`, {
      headers: this.getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Failed to refresh user profile');
    }
    
    const userData = await response.json();
    const user: User = {
      id: userData.id,
      name: userData.name,
      email: userData.email,
      picture: userData.picture,
      token: this.getToken() || '',
      role: userData.role,
      boardRoles: userData.boardRoles
    };
    
    // Update stored user
    localStorage.setItem('user', JSON.stringify(user));
    this.currentUser = user;
    
    return user;
  }
}

export const authService = new AuthService();
