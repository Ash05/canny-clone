import React, { useState } from 'react';

interface AddCommentProps {
  feedbackId: number;
  onAddComment: (content: string) => void;
}

const AddComment: React.FC<AddCommentProps> = ({ feedbackId, onAddComment }) => {
  const [content, setContent] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!content.trim()) {
      setError('Comment cannot be empty');
      return;
    }
    
    onAddComment(content);
    setContent('');
    setError('');
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-lg shadow p-4 mb-6">
      <h3 className="text-lg font-semibold mb-3">Add a Comment</h3>
      <textarea
        className="w-full border rounded p-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="Write a comment..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={3}
      />
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
      <div className="flex justify-end mt-2">
        <button 
          type="submit"
          className="px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition"
        >
          Comment
        </button>
      </div>
    </form>
  );
};

export default AddComment;
