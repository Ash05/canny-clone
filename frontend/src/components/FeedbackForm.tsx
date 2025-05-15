import React from 'react';
import { useForm } from 'react-hook-form';

interface FeedbackFormProps {
  categories: { id: number; name: string }[];
  onSubmit: (title: string, description: string, categoryId: number) => void;
  onCancel?: () => void;
}

function FeedbackForm({ categories, onSubmit, onCancel }: FeedbackFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<{ title: string; description: string; categoryId: number }>();

  const handleFormSubmit = (data: { title: string; description: string; categoryId: number }) => {
    onSubmit(data.title, data.description, Number(data.categoryId));
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4 p-4 border rounded shadow">
      <input
        type="text"
        placeholder="Title"
        {...register('title', { required: 'Title is required.' })}
        className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      {errors.title && <p className="text-red-500 text-sm">{errors.title.message}</p>}

      <textarea
        placeholder="Description"
        {...register('description', { required: 'Description is required.' })}
        className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      ></textarea>
      {errors.description && <p className="text-red-500 text-sm">{errors.description.message}</p>}

      <select
        {...register('categoryId', { required: 'Category is required.' })}
        className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">Select Category</option>
        {categories.map((category) => (
          <option key={category.id} value={category.id}>
            {category.name}
          </option>
        ))}
      </select>
      {errors.categoryId && <p className="text-red-500 text-sm">{errors.categoryId.message}</p>}      <div className="flex space-x-2">
        <button type="submit" className="flex-1 bg-blue-500 text-white py-2 rounded hover:bg-blue-600">
          Submit
        </button>
        {onCancel && (
          <button 
            type="button" 
            onClick={onCancel}
            className="flex-1 bg-gray-300 text-gray-700 py-2 rounded hover:bg-gray-400"
          >
            Cancel
          </button>
        )}
      </div>
    </form>
  );
}

export default FeedbackForm;
