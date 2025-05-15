import React from 'react';

interface FeedbackCardProps {
  feedback: {
    id: number;
    title: string;
    description: string;
    upvotes: number;
    downvotes: number;
    categoryId?: number;
    status?: string;
    userVote?: 'upvote' | 'downvote' | null | undefined;
  };
  category?: {
    id: number;
    name: string;
  };
  onVote: (feedbackId: number, voteType: 'upvote' | 'downvote') => void;
  onViewMore: (feedback: any) => void;
}

function FeedbackCard({ feedback, category, onVote, onViewMore }: FeedbackCardProps) {
  return (
    <div className="border p-4 rounded shadow hover:shadow-lg transition">
      <h3 className="text-lg font-bold">{feedback.title}</h3>
      <p className="text-sm text-gray-600">{feedback.description.substring(0, 100)}...</p>
      {category && (
        <span className="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
          {category.name}
        </span>
      )}
      <div className="flex space-x-2 mt-2">
        <button
          onClick={() => onVote(feedback.id, 'upvote')}
          className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
        >
          Upvote ({feedback.upvotes})
        </button>
        <button
          onClick={() => onVote(feedback.id, 'downvote')}
          className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
        >
          Downvote ({feedback.downvotes})
        </button>
        <button
          onClick={() => onViewMore(feedback)}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          View More
        </button>
      </div>
    </div>
  );
}

export default FeedbackCard;
