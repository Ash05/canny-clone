import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { boardService, Board } from '../services/boardService';
import { authService } from '../services/authService';

interface FormData {
  boardName: string;
}

function Boards() {
  const [boards, setBoards] = useState<Board[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>();
  // Get user's role
  const isAppAdmin = authService.isAppAdmin();
  
  useEffect(() => {
    // Redirect to sign-in if not authenticated
    if (!authService.isAuthenticated()) {
      navigate('/signin');
      return;
    }

    fetchBoards();
  }, [navigate]);

  const fetchBoards = async () => {
    setIsLoading(true);
    setErrorMessage(null);
    
    try {
      const fetchedBoards = await boardService.getAllBoards();
      setBoards(fetchedBoards);
    } catch (error) {
      console.error('Error fetching boards:', error);
      setErrorMessage('Failed to fetch boards. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  const validateBoardName = (name: string): string | undefined => {
    const error = boardService.validateBoardName(name);
    return error || undefined;
  };

  const onSubmit = async (data: FormData) => {
    setIsLoading(true);
    setErrorMessage(null);
    setSuccessMessage(null);
    
    try {
      await boardService.createBoard(data.boardName);
      setSuccessMessage('Board created successfully!');
      reset(); // Reset form
      fetchBoards(); // Refresh boards list
    } catch (error) {
      console.error('Error creating board:', error);
      setErrorMessage('Failed to create board. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };
  const handleBoardClick = (boardId: number) => {
    navigate(`/board/${boardId}`);
  };
  return (
    <div className="container mx-auto p-6 max-w-4xl">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">Feedback Boards</h1>
          {isAppAdmin && (
            <div className="mt-1 text-sm text-blue-600 font-semibold">
              Admin Dashboard
            </div>
          )}
        </div>
        <div className="flex items-center">
          <div className="mr-4 text-right">
            <span className="text-sm text-gray-600">
              {authService.getCurrentUser()?.email || 'Not signed in'}
            </span>
            <div className="text-xs text-gray-500">
              Role: {
                isAppAdmin 
                  ? 'App Admin' 
                  : authService.isStakeholder() 
                    ? 'Stakeholder' 
                    : 'User'
              }
            </div>
          </div>
          <button 
            onClick={() => {
              authService.logout();
              navigate('/signin');
            }}
            className="px-4 py-2 text-sm text-gray-600 hover:text-gray-800 border border-gray-300 rounded"
          >
            Sign Out
          </button>
        </div>
      </div>
      
      {errorMessage && (
        <div className="mb-4 p-4 bg-red-100 text-red-700 rounded-md">
          {errorMessage}
        </div>
      )}
      
      {successMessage && (
        <div className="mb-4 p-4 bg-green-100 text-green-700 rounded-md">
          {successMessage}
        </div>
      )}      {isAppAdmin && (
        <div className="bg-white rounded-lg shadow-md p-6 mb-8">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Create a New Board</h2>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label htmlFor="boardName" className="block text-sm font-medium text-gray-700 mb-1">
                Board Name
              </label>
              <input
                id="boardName"
                type="text"
                {...register('boardName', {
                  required: 'Board name is required',
                  validate: validateBoardName
                })}
                className="w-full px-4 py-2 border rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter a name for your feedback board"
              />
              {errors.boardName && (
                <p className="mt-1 text-sm text-red-600">{errors.boardName.message}</p>
              )}
            </div>
            
            <button
              type="submit"
              disabled={isSubmitting || isLoading}
              className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
            >
              {isSubmitting || isLoading ? 'Creating...' : 'Create Board'}
            </button>
          </form>
        </div>
      )}
      
      <div>
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Your Boards</h2>
        
        {isLoading && !boards.length ? (
          <div className="text-center p-8">
            <p className="text-gray-500">Loading boards...</p>
          </div>
        ) : !boards.length ? (
          <div className="bg-gray-50 rounded-lg p-8 text-center border-2 border-dashed border-gray-300">
            <p className="text-gray-500">You haven't created any boards yet.</p>
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {boards.map((board) => (
              <div
                key={board.id}
                onClick={() => handleBoardClick(board.id)}
                className="bg-white rounded-lg shadow-md p-6 cursor-pointer hover:shadow-lg transition-shadow border border-gray-200 hover:border-blue-300"
              >
                <h3 className="text-lg font-medium text-gray-800">{board.name}</h3>
                <div className="mt-2 flex justify-between items-center">
                  <span className="inline-block bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                    Click to view
                  </span>
                  
                  <span className={`inline-block text-xs px-2 py-1 rounded-full ${
                    authService.isBoardStakeholder(board.id)
                      ? 'bg-purple-100 text-purple-800'
                      : 'bg-green-100 text-green-800'
                  }`}>
                    {authService.isBoardStakeholder(board.id) 
                      ? 'Stakeholder' 
                      : isAppAdmin 
                        ? 'Admin' 
                        : 'Member'}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}      </div>
    </div>
  );
}

export default Boards;
