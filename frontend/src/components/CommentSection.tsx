import React, { useState, useEffect } from 'react';
import Comment from './Comment';
import AddComment from './AddComment';
import { environment } from '../environments/environment';

interface CommentReply {
  id: number;
  commentId: number;
  userId: number;
  content: string;
  likes: number;
  dislikes: number;
  createdAt: string;
  isLiked: boolean;
  isDisliked: boolean;
}

interface CommentType {
  id: number;
  feedbackId: number;
  userId: number;
  content: string;
  likes: number;
  dislikes: number;
  createdAt: string;
  isLiked: boolean;
  isDisliked: boolean;
  replies: CommentReply[];
}

interface CommentSectionProps {
  feedbackId: number;
}

const CommentSection: React.FC<CommentSectionProps> = ({ feedbackId }) => {
  const [comments, setComments] = useState<CommentType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchComments = async () => {
    try {
      setLoading(true);
      const res = await fetch(`${environment.apiUrl}/comments?feedbackId=${feedbackId}`);
      
      if (!res.ok) {
        throw new Error('Failed to fetch comments');
      }
      
      const data = await res.json();
      setComments(data);
      setError('');
    } catch (err) {
      setError('Error loading comments');
      console.error('Error fetching comments:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComments();
  }, [feedbackId]);

  const handleAddComment = async (content: string) => {
    try {
      const res = await fetch(`${environment.apiUrl}/comment`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          feedbackId,
          content,
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to add comment');
      }
      
      // Refresh comments after adding
      fetchComments();
    } catch (err) {
      setError('Error adding comment');
      console.error('Error adding comment:', err);
    }
  };

  const handleAddReply = async (commentId: number, content: string) => {
    try {
      const res = await fetch(`${environment.apiUrl}/reply`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          commentId,
          content,
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to add reply');
      }
      
      // Refresh comments after adding
      fetchComments();
    } catch (err) {
      setError('Error adding reply');
      console.error('Error adding reply:', err);
    }
  };

  const handleLikeComment = async (commentId: number, isLike: boolean) => {
    try {
      const res = await fetch(`${environment.apiUrl}/comment-like`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          commentId,
          isLike,
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to like comment');
      }
      
      // Update local state to show immediate feedback
      setComments(prevComments => 
        prevComments.map(comment => {
          if (comment.id === commentId) {
            let likes = comment.likes;
            let dislikes = comment.dislikes;
            
            // Logic to update like/dislike counts
            if (comment.isLiked && isLike) {
              // User is un-liking (toggling off)
              likes--;
            } else if (comment.isDisliked && !isLike) {
              // User is un-disliking (toggling off)
              dislikes--;
            } else if (comment.isLiked && !isLike) {
              // User is switching from like to dislike
              likes--;
              dislikes++;
            } else if (comment.isDisliked && isLike) {
              // User is switching from dislike to like
              dislikes--;
              likes++;
            } else if (isLike) {
              // New like
              likes++;
            } else {
              // New dislike
              dislikes++;
            }
            
            return {
              ...comment,
              likes,
              dislikes,
              isLiked: comment.isLiked === isLike ? !comment.isLiked : isLike,
              isDisliked: comment.isDisliked === !isLike ? !comment.isDisliked : !isLike
            };
          }
          return comment;
        })
      );
    } catch (err) {
      setError('Error liking comment');
      console.error('Error liking comment:', err);
    }
  };

  const handleLikeReply = async (replyId: number, isLike: boolean) => {
    try {
      const res = await fetch(`${environment.apiUrl}/comment-like`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          replyId,
          isLike,
        }),
      });
      
      if (!res.ok) {
        throw new Error('Failed to like reply');
      }
      
      // Update local state to show immediate feedback
      setComments(prevComments => 
        prevComments.map(comment => ({
          ...comment,
          replies: comment.replies.map(reply => {
            if (reply.id === replyId) {
              let likes = reply.likes;
              let dislikes = reply.dislikes;
              
              // Logic to update like/dislike counts
              if (reply.isLiked && isLike) {
                // User is un-liking (toggling off)
                likes--;
              } else if (reply.isDisliked && !isLike) {
                // User is un-disliking (toggling off)
                dislikes--;
              } else if (reply.isLiked && !isLike) {
                // User is switching from like to dislike
                likes--;
                dislikes++;
              } else if (reply.isDisliked && isLike) {
                // User is switching from dislike to like
                dislikes--;
                likes++;
              } else if (isLike) {
                // New like
                likes++;
              } else {
                // New dislike
                dislikes++;
              }
              
              return {
                ...reply,
                likes,
                dislikes,
                isLiked: reply.isLiked === isLike ? !reply.isLiked : isLike,
                isDisliked: reply.isDisliked === !isLike ? !reply.isDisliked : !isLike
              };
            }
            return reply;
          })
        }))
      );
    } catch (err) {
      setError('Error liking reply');
      console.error('Error liking reply:', err);
    }
  };

  if (loading && comments.length === 0) {
    return <div className="text-center py-8">Loading comments...</div>;
  }

  return (
    <div className="mt-8">
      <h2 className="text-xl font-bold mb-4">Comments ({comments.length})</h2>
      
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}
      
      <AddComment feedbackId={feedbackId} onAddComment={handleAddComment} />
      
      <div className="space-y-4">
        {comments.map(comment => (
          <Comment 
            key={comment.id} 
            comment={comment} 
            onLike={handleLikeComment}
            onReplyLike={handleLikeReply}
            onAddReply={handleAddReply}
          />
        ))}
        
        {comments.length === 0 && !loading && (
          <div className="text-gray-500 text-center py-8">
            No comments yet. Be the first to comment!
          </div>
        )}
      </div>
    </div>
  );
};

export default CommentSection;
