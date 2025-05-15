import React, { useState } from 'react';

interface User {
  id: number;
  name: string;
}

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

interface CommentProps {
  comment: {
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
  };
  onLike: (commentId: number, isLike: boolean) => void;
  onReplyLike: (replyId: number, isLike: boolean) => void;
  onAddReply: (commentId: number, content: string) => void;
}

const Comment: React.FC<CommentProps> = ({ comment, onLike, onReplyLike, onAddReply }) => {
  const [showReplyForm, setShowReplyForm] = useState(false);
  const [replyContent, setReplyContent] = useState('');
  const [error, setError] = useState('');

  const handleSubmitReply = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!replyContent.trim()) {
      setError('Reply cannot be empty');
      return;
    }
    
    onAddReply(comment.id, replyContent);
    setReplyContent('');
    setShowReplyForm(false);
    setError('');
  };

  return (
    <div className="bg-white rounded-lg shadow p-4 my-4">
      <div className="flex justify-between items-start">
        <p className="font-semibold">User {comment.userId}</p>
        <span className="text-sm text-gray-500">{new Date(comment.createdAt).toLocaleString()}</span>
      </div>
      
      <p className="my-2">{comment.content}</p>
      
      <div className="flex space-x-4">
        <button 
          className={`flex items-center ${comment.isLiked ? 'text-blue-500' : 'text-gray-500'}`}
          onClick={() => onLike(comment.id, true)}
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
            <path d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z" />
          </svg>
          {comment.likes}
        </button>
        
        <button 
          className={`flex items-center ${comment.isDisliked ? 'text-red-500' : 'text-gray-500'}`}
          onClick={() => onLike(comment.id, false)}
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
            <path d="M18 9.5a1.5 1.5 0 11-3 0v-6a1.5 1.5 0 013 0v6zM14 9.667v-5.43a2 2 0 00-1.105-1.79l-.05-.025A4 4 0 0011.055 2H5.64a2 2 0 00-1.962 1.608l-1.2 6A2 2 0 004.44 12H8v4a2 2 0 002 2 1 1 0 001-1v-.667a4 4 0 01.8-2.4l1.4-1.866a4 4 0 00.8-2.4z" />
          </svg>
          {comment.dislikes}
        </button>
        
        <button 
          className="text-blue-500"
          onClick={() => setShowReplyForm(!showReplyForm)}
        >
          Reply
        </button>
      </div>
      
      {showReplyForm && (
        <form onSubmit={handleSubmitReply} className="mt-4">
          <textarea
            className="w-full border rounded p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Write a reply..."
            value={replyContent}
            onChange={(e) => setReplyContent(e.target.value)}
          />
          {error && <p className="text-red-500 text-sm">{error}</p>}
          <div className="flex justify-end space-x-2 mt-2">
            <button 
              type="button"
              className="px-4 py-2 text-gray-500"
              onClick={() => {
                setShowReplyForm(false);
                setError('');
              }}
            >
              Cancel
            </button>
            <button 
              type="submit"
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Submit
            </button>
          </div>
        </form>
      )}
      
      {/* Replies */}
      <div className="mt-4 pl-6 border-l-2 border-gray-200">
        {comment.replies && comment.replies.map(reply => (
          <div key={reply.id} className="bg-gray-50 rounded-lg p-3 my-2">
            <div className="flex justify-between items-start">
              <p className="font-semibold">User {reply.userId}</p>
              <span className="text-sm text-gray-500">{new Date(reply.createdAt).toLocaleString()}</span>
            </div>
            <p className="my-2">{reply.content}</p>
            <div className="flex space-x-4">
              <button 
                className={`flex items-center ${reply.isLiked ? 'text-blue-500' : 'text-gray-500'}`}
                onClick={() => onReplyLike(reply.id, true)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z" />
                </svg>
                {reply.likes}
              </button>
              <button 
                className={`flex items-center ${reply.isDisliked ? 'text-red-500' : 'text-gray-500'}`}
                onClick={() => onReplyLike(reply.id, false)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M18 9.5a1.5 1.5 0 11-3 0v-6a1.5 1.5 0 013 0v6zM14 9.667v-5.43a2 2 0 00-1.105-1.79l-.05-.025A4 4 0 0011.055 2H5.64a2 2 0 00-1.962 1.608l-1.2 6A2 2 0 004.44 12H8v4a2 2 0 002 2 1 1 0 001-1v-.667a4 4 0 01.8-2.4l1.4-1.866a4 4 0 00.8-2.4z" />
                </svg>
                {reply.dislikes}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Comment;
