import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { authService } from '../services/authService';
import { feedbackService, Feedback } from '../services/feedbackService';
import { categoryService, Category } from '../services/categoryService';
import { boardService, BoardMember } from '../services/boardService';
import FeedbackForm from '../components/FeedbackForm';
import FeedbackList from '../components/FeedbackList';
import CommentSection from '../components/CommentSection';
import BoardMemberManagement from '../components/BoardMemberManagement';

function BoardDetails() {
  const { boardId } = useParams<{ boardId: string }>();
  const navigate = useNavigate();
  
  const [board, setBoard] = useState<{ id: number; name: string } | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [selectedFeedback, setSelectedFeedback] = useState<Feedback | null>(null);
  const [boardMembers, setBoardMembers] = useState<BoardMember[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showFeedbackForm, setShowFeedbackForm] = useState(false);
  
  // Get current user's role on this board
  const boardIdNum = boardId ? parseInt(boardId) : 0;
  const isAppAdmin = authService.isAppAdmin();
  const isStakeholder = authService.isBoardStakeholder(boardIdNum);
  const currentUserId = authService.getUserId();

  useEffect(() => {
    // Redirect to sign-in if not authenticated
    if (!authService.isAuthenticated()) {
      navigate('/signin');
      return;
    }

    // Check if user has access to this board
    if (boardId && !authService.hasBoardAccess(parseInt(boardId)) && !isAppAdmin) {
      navigate('/forbidden');
      return;
    }

    if (boardId) {
      fetchData();
    }
  }, [boardId, navigate, isAppAdmin]);

  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Fetch board details
      if (boardId) {
        const boardData = await boardService.getBoardById(parseInt(boardId));
        setBoard(boardData);
      }

      // Fetch categories
      const categoriesData = await categoryService.getAllCategories();
      setCategories(categoriesData);
      
      // Fetch feedbacks for this board
      if (boardId) {
        const feedbacksData = await feedbackService.getFeedbacksByBoardId(parseInt(boardId));
        setFeedbacks(feedbacksData);
      }
      
      // If stakeholder or admin, fetch board members
      if ((isStakeholder || isAppAdmin) && boardId) {
        try {
          const response = await fetch(`${process.env.REACT_APP_API_URL}/boards/${boardId}/members`, {
            headers: authService.getAuthHeader()
          });
          
          if (response.ok) {
            const membersData = await response.json();
            setBoardMembers(membersData);
          }
        } catch (err) {
          console.error('Error fetching board members:', err);
          // Non-critical error, don't show to user
        }
      }
    } catch (err) {
      console.error('Error fetching data:', err);
      setError('Failed to load data. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleFeedbackSubmit = async (title: string, description: string, categoryId: number) => {
    setError(null);
    
    try {
      if (!boardId) return;
      
      await feedbackService.submitFeedback({
        boardId: parseInt(boardId),
        title,
        description,
        categoryId
      });
      
      // Refresh feedbacks
      fetchData();
      setShowFeedbackForm(false);
    } catch (err) {
      console.error('Error submitting feedback:', err);
      setError('Failed to submit feedback. Please try again.');
    }
  };

  const handleVote = async (feedbackId: number, voteType: 'upvote' | 'downvote') => {
    try {
      await feedbackService.voteFeedback(feedbackId, voteType);
      
      // Update local state to reflect vote
      setFeedbacks(prevFeedbacks => 
        prevFeedbacks.map(feedback => {
          if (feedback.id === feedbackId) {
            return {
              ...feedback,
              upvotes: voteType === 'upvote' ? feedback.upvotes + 1 : feedback.upvotes,
              downvotes: voteType === 'downvote' ? feedback.downvotes + 1 : feedback.downvotes,
              userVote: voteType
            };
          }
          return feedback;
        })
      );
    } catch (err) {
      console.error('Error voting on feedback:', err);
      setError('Failed to register vote. Please try again.');
    }
  };

  const handleUpdateStatus = async (feedbackId: number, newStatus: string) => {
    if (!isStakeholder && !isAppAdmin) {
      setError("You don't have permission to update feedback status");
      return;
    }
    
    try {
      await feedbackService.updateFeedbackStatus(feedbackId, newStatus as any);
      
      // Update local state to reflect the status change
      setFeedbacks(prevFeedbacks => 
        prevFeedbacks.map(feedback => {
          if (feedback.id === feedbackId) {
            return {
              ...feedback,
              status: newStatus as any
            };
          }
          return feedback;
        })
      );
      
      // If this is the currently selected feedback, update it too
      if (selectedFeedback && selectedFeedback.id === feedbackId) {
        setSelectedFeedback({
          ...selectedFeedback,
          status: newStatus as any
        });
      }
    } catch (err) {
      console.error('Error updating feedback status:', err);
      setError('Failed to update feedback status. Please try again.');
    }
  };

  const handleFeedbackClick = (feedback: Feedback) => {
    setSelectedFeedback(feedback);
  };

  const handleBackToList = () => {
    setSelectedFeedback(null);
  };

  if (isLoading && !categories.length) {
    return (
      <div className="container mx-auto p-6 text-center">
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-6">
        <div className="bg-red-100 p-4 rounded text-red-700 mb-6">{error}</div>
        <button 
          onClick={fetchData} 
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6 max-w-6xl">
      <div className="mb-6 flex justify-between items-center">
        <div>
          <button 
            onClick={() => navigate('/boards')}
            className="text-blue-500 hover:text-blue-700 flex items-center"
          >
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Back to Boards
          </button>
          <h1 className="text-3xl font-bold mt-2">
            {board?.name || `Board #${boardId}`}
          </h1>
        </div>
        {!selectedFeedback && (
          <button
            onClick={() => setShowFeedbackForm(!showFeedbackForm)}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            {showFeedbackForm ? 'Cancel' : 'Add Feedback'}
          </button>
        )}
      </div>

      {selectedFeedback ? (
        <div className="bg-white rounded-lg shadow-md p-6">
          <button 
            onClick={handleBackToList}
            className="text-blue-500 hover:text-blue-700 mb-4 flex items-center"
          >
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Back to List
          </button>
          
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-2">{selectedFeedback.title}</h2>
            <div className="text-gray-700 mb-4">
              {selectedFeedback.description}
            </div>
            <div className="flex items-center justify-between">
              <div>
                <span className="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                  {categories.find(c => c.id === selectedFeedback.categoryId)?.name || 'Uncategorized'}
                </span>
                
                {/* Status Badge */}
                <span className={`ml-2 inline-block px-2 py-1 text-xs font-semibold rounded-full ${
                  selectedFeedback.status === 'approved' ? 'bg-green-100 text-green-800' :
                  selectedFeedback.status === 'reviewing' ? 'bg-blue-100 text-blue-800' :
                  selectedFeedback.status === 'declined' ? 'bg-red-100 text-red-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {selectedFeedback.status || 'Pending'}
                </span>
              </div>
              
              <div className="flex space-x-2">                <button 
                  onClick={() => handleVote(selectedFeedback.id, 'upvote')}
                  className={`flex items-center space-x-1 px-2 py-1 rounded ${
                    selectedFeedback.userVote === 'upvote' 
                      ? 'bg-green-100 text-green-700' 
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <span>{selectedFeedback.upvotes}</span>
                </button>                <button 
                  onClick={() => handleVote(selectedFeedback.id, 'downvote')}
                  className={`flex items-center space-x-1 px-2 py-1 rounded ${
                    selectedFeedback.userVote === 'downvote' 
                      ? 'bg-red-100 text-red-700' 
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7" />
                  </svg>
                  <span>{selectedFeedback.downvotes}</span>
                </button>
              </div>
            </div>
          </div>
          
          <CommentSection feedbackId={selectedFeedback.id} />
          
          {/* Status Control for Stakeholders and Admins */}
          {(isStakeholder || isAppAdmin) && (
            <div className="mt-4 border-t pt-4">
              <h3 className="text-sm font-semibold text-gray-500 mb-2">Update Status:</h3>
              <div className="flex flex-wrap gap-2">
                {['pending', 'reviewing', 'approved', 'declined'].map(status => (
                  <button
                    key={status}
                    onClick={() => handleUpdateStatus(selectedFeedback.id, status)}
                    className={`px-3 py-1 text-xs rounded-full ${
                      selectedFeedback.status === status 
                        ? 'bg-blue-500 text-white' 
                        : 'bg-gray-200 hover:bg-gray-300 text-gray-800'
                    }`}
                  >
                    {status.charAt(0).toUpperCase() + status.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      ) : (
        <>
          {showFeedbackForm && (
            <div className="mb-8">
              <FeedbackForm 
                categories={categories} 
                onSubmit={handleFeedbackSubmit} 
                onCancel={() => setShowFeedbackForm(false)}
              />
            </div>
          )}
          
          <FeedbackList 
            feedbacks={feedbacks} 
            categories={categories} 
            onFeedbackClick={handleFeedbackClick} 
            onVote={handleVote}
          />
        </>
      )}
    </div>
  );
}

export default BoardDetails;
