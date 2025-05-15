import React from 'react';
import FeedbackCard from './FeedbackCard';

interface FeedbackListProps {
  feedbacks: {
    id: number;
    title: string;
    description: string;
    upvotes: number;
    downvotes: number;
    categoryId: number;
    status?: string;
    userVote?: 'upvote' | 'downvote' | null | undefined;
  }[];
  categories: { id: number; name: string }[];
  onVote: (feedbackId: number, voteType: 'upvote' | 'downvote') => void;
  onFeedbackClick: (feedback: any) => void;
}

function FeedbackList({ feedbacks, categories, onVote, onFeedbackClick }: FeedbackListProps) {
  return (
    <div className="space-y-4">
      {feedbacks.map((feedback) => (
        <FeedbackCard
          key={feedback.id}
          feedback={feedback}
          category={categories.find(c => c.id === feedback.categoryId)}
          onVote={onVote}
          onViewMore={onFeedbackClick}
        />
      ))}
    </div>
  );
}

export default FeedbackList;
