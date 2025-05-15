import { environment } from '../environments/environment';
import { authService } from './authService';

export interface Board {
  id: number;
  name: string;
}

export interface BoardMember {
  id: number;
  name: string;
  email: string;
  picture?: string;
  role: string;
}

class BoardService {
  // Get all boards
  async getAllBoards(): Promise<Board[]> {
    const response = await fetch(`${environment.apiUrl}/boards`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch boards');
    }
    
    return response.json();
  }
  
  // Create a new board
  async createBoard(name: string): Promise<void> {
    const response = await fetch(`${environment.apiUrl}/board`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({ name })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      throw new Error(errorData.message || 'Failed to create board');
    }
  }
  
  // Get board by ID
  async getBoardById(boardId: number): Promise<Board> {
    const response = await fetch(`${environment.apiUrl}/boards/${boardId}`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch board details');
    }
    
    return response.json();
  }
  
  // Get board members (for stakeholders and admins)
  async getBoardMembers(boardId: number): Promise<BoardMember[]> {
    // Check if user has permission
    if (!authService.isAppAdmin() && !authService.isBoardStakeholder(boardId)) {
      throw new Error('Unauthorized: Only stakeholders and admins can view board members');
    }
    
    const response = await fetch(`${environment.apiUrl}/boards/${boardId}/members`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch board members');
    }
    
    return response.json();
  }
  
  // Invite a user to a board (for stakeholders and admins)
  async inviteUserToBoard(boardId: number, email: string, role: 'stakeholder' | 'user'): Promise<void> {
    // Check if user has permission
    if (!authService.isAppAdmin() && !authService.isBoardStakeholder(boardId)) {
      throw new Error('Unauthorized: Only stakeholders and admins can invite users');
    }
    
    const response = await fetch(`${environment.apiUrl}/boards/${boardId}/members`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({ email, role })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      throw new Error(errorData.message || 'Failed to invite user');
    }
  }
  
  // Remove a user from a board (for stakeholders and admins)
  async removeUserFromBoard(boardId: number, userId: number): Promise<void> {
    // Check if user has permission
    if (!authService.isAppAdmin() && !authService.isBoardStakeholder(boardId)) {
      throw new Error('Unauthorized: Only stakeholders and admins can remove users');
    }
    
    const response = await fetch(`${environment.apiUrl}/boards/${boardId}/members/${userId}`, {
      method: 'DELETE',
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to remove user from board');
    }
  }
  
  // Client-side validation for board name
  validateBoardName(name: string): string | null {
    if (!name.trim()) {
      return 'Board name cannot be empty';
    }
    
    if (name.length < 3) {
      return 'Board name must be at least 3 characters';
    }
    
    if (name.length > 50) {
      return 'Board name must be less than 50 characters';
    }
    
    if (/[^a-zA-Z0-9 ]/.test(name)) {
      return 'Board name can only contain letters, numbers, and spaces';
    }
    
    return null;
  }
}

export const boardService = new BoardService();
