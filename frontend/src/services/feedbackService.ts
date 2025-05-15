import { environment } from '../environments/environment';
import { authService } from './authService';

export type FeedbackStatus = 'pending' | 'reviewing' | 'approved' | 'declined';

export interface Feedback {
  id: number;
  boardId: number;
  title: string;
  description: string;
  categoryId: number;
  upvotes: number;
  downvotes: number;
  status?: FeedbackStatus;
  userVote?: 'upvote' | 'downvote' | null;
}

export interface FeedbackSubmission {
  boardId: number;
  title: string;
  description: string;
  categoryId: number;
}

class FeedbackService {
  // Get all feedback for a board
  async getFeedbacksByBoardId(boardId: number): Promise<Feedback[]> {
    const response = await fetch(`${environment.apiUrl}/feedbacks?boardId=${boardId}`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch feedbacks');
    }
    
    return response.json();
  }
  
  // Submit new feedback
  async submitFeedback(feedback: FeedbackSubmission): Promise<void> {
    const response = await fetch(`${environment.apiUrl}/feedback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify(feedback)
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Unknown error' }));
      throw new Error(errorData.message || 'Failed to submit feedback');
    }
  }
  
  // Get feedback details by ID
  async getFeedbackById(id: number): Promise<Feedback> {
    const response = await fetch(`${environment.apiUrl}/feedback/${id}`, {
      headers: {
        ...authService.getAuthHeader()
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to fetch feedback details');
    }
    
    return response.json();
  }
    // Vote on a feedback
  async voteFeedback(feedbackId: number, voteType: 'upvote' | 'downvote'): Promise<void> {
    const response = await fetch(`${environment.apiUrl}/vote`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authService.getAuthHeader()
      },
      body: JSON.stringify({
        feedbackId,
        voteType
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to vote on feedback');
    }
  }
    // Get the board ID for a feedback item
  async getFeedbackBoardId(feedbackId: number): Promise<number> {
    try {
      const feedback = await this.getFeedbackById(feedbackId);
      return feedback.boardId;
    } catch (err) {
      console.error('Error getting feedback board ID:', err);
      throw err;
    }
  }
    // Update feedback status (stakeholders and admins only)
  async updateFeedbackStatus(feedbackId: number, status: FeedbackStatus): Promise<void> {
    try {
      // First, get the board ID for this feedback
      const feedback = await this.getFeedbackById(feedbackId);
      const boardId = feedback.boardId;
      
      // Check if user has permission
      if (!authService.isAppAdmin() && !authService.isBoardStakeholder(boardId)) {
        throw new Error('Unauthorized: Only stakeholders and admins can update feedback status');
      }
      
      const response = await fetch(`${environment.apiUrl}/feedbacks/${feedbackId}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...authService.getAuthHeader()
        },
        body: JSON.stringify({ status })
      });
      
      if (!response.ok) {
        throw new Error('Failed to update feedback status');
      }
    } catch (error) {
      console.error('Error updating feedback status:', error);
      throw error;
    }
  }
}

export const feedbackService = new FeedbackService();
