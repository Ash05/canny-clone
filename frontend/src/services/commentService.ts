import { environment } from '../environments/environment';
import { authService } from './authService';

export interface Reply {
  id: number;
  commentId: number;
  userId: number;
  content: string;
  likes: number;
  dislikes: number;
  createdAt: string;
  isLiked?: boolean;
  isDisliked?: boolean;
}

export interface Comment {
  id: number;
  feedbackId: number;
  userId: number;
  content: string;
  likes: number;
  dislikes: number;
  createdAt: string;
  isLiked?: boolean;
  isDisliked?: boolean;
  replies: Reply[];
}

class CommentService {
  // Get comments for a feedback
  async getCommentsByFeedbackId(feedbackId: number): Promise<Comment[]> {
    const response = await fetch(`${environment.apiUrl}/comments?feedbackId=${feedbackId}`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch comments');
    }
    
    return response.json();
  }
  
  // Add a comment to a feedback
  async addComment(feedbackId: number, content: string): Promise<number> {
    const response = await fetch(`${environment.apiUrl}/comment`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({
        feedbackId,
        content
      })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      throw new Error(errorData.message || 'Failed to add comment');
    }
    
    const result = await response.json();
    return result.id;
  }
  
  // Add a reply to a comment
  async addReply(commentId: number, content: string): Promise<number> {
    const response = await fetch(`${environment.apiUrl}/reply`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({
        commentId,
        content
      })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      throw new Error(errorData.message || 'Failed to add reply');
    }
    
    const result = await response.json();
    return result.id;
  }
  
  // Like/dislike a comment
  async likeComment(commentId: number, isLike: boolean): Promise<void> {
    const response = await fetch(`${environment.apiUrl}/comment-like`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({
        commentId,
        isLike
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to like/dislike comment');
    }
  }
  
  // Like/dislike a reply
  async likeReply(replyId: number, isLike: boolean): Promise<void> {
    const response = await fetch(`${environment.apiUrl}/comment-like`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({
        replyId,
        isLike
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to like/dislike reply');
    }
  }
}

export const commentService = new CommentService();
