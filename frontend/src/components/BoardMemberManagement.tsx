import React, { useState, FormEvent } from 'react';
import { BoardMember as IBoardMember, boardService } from '../services/boardService';

interface BoardMemberManagementProps {
  boardId: number;
  members: IBoardMember[];
  onMemberAdded: () => void;
}

const BoardMemberManagement: React.FC<BoardMemberManagementProps> = ({ 
  boardId, 
  members, 
  onMemberAdded 
}) => {
  const [email, setEmail] = useState('');
  const [role, setRole] = useState<'stakeholder' | 'user'>('user');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [showAddForm, setShowAddForm] = useState(false);
  
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setIsLoading(true);
    
    try {
      await boardService.inviteUserToBoard(boardId, email, role);
      setSuccess(`${role === 'stakeholder' ? 'Stakeholder' : 'User'} invited successfully!`);
      setEmail('');
      onMemberAdded();
      
      // Hide form after successful invite
      setTimeout(() => {
        setShowAddForm(false);
        setSuccess(null);
      }, 3000);
    } catch (err) {
      setError((err as Error).message || 'Failed to invite user');
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleRemoveMember = async (userId: number) => {
    setError(null);
    setSuccess(null);
    
    try {
      await boardService.removeUserFromBoard(boardId, userId);
      setSuccess('User removed successfully!');
      onMemberAdded();
      
      // Hide success message after a delay
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      setError((err as Error).message || 'Failed to remove user');
    }
  };
  
  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">Board Members</h2>
        <button
          onClick={() => setShowAddForm(!showAddForm)}
          className="px-3 py-1 bg-blue-500 text-white text-sm rounded hover:bg-blue-600"
        >
          {showAddForm ? 'Cancel' : 'Add Member'}
        </button>
      </div>
      
      {error && (
        <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4">
          <p>{error}</p>
        </div>
      )}
      
      {success && (
        <div className="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4">
          <p>{success}</p>
        </div>
      )}
      
      {showAddForm && (
        <form onSubmit={handleSubmit} className="mb-6 p-4 border border-gray-200 rounded">
          <div className="mb-4">
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
              Email Address
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              required
              placeholder="email@example.com"
            />
          </div>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Role
            </label>
            <div className="flex space-x-4">
              <label className="inline-flex items-center">
                <input
                  type="radio"
                  name="role"
                  value="user"
                  checked={role === 'user'}
                  onChange={() => setRole('user')}
                  className="form-radio h-4 w-4 text-blue-600"
                />
                <span className="ml-2 text-sm text-gray-700">User</span>
              </label>
              <label className="inline-flex items-center">
                <input
                  type="radio"
                  name="role"
                  value="stakeholder"
                  checked={role === 'stakeholder'}
                  onChange={() => setRole('stakeholder')}
                  className="form-radio h-4 w-4 text-blue-600"
                />
                <span className="ml-2 text-sm text-gray-700">Stakeholder</span>
              </label>
            </div>
          </div>
          
          <button
            type="submit"
            disabled={isLoading}
            className={`w-full px-4 py-2 ${
              isLoading ? 'bg-gray-400' : 'bg-blue-600 hover:bg-blue-700'
            } text-white rounded-md`}
          >
            {isLoading ? 'Inviting...' : 'Invite User'}
          </button>
        </form>
      )}
      
      <div className="overflow-y-auto max-h-80">
        {members.length === 0 ? (
          <p className="text-gray-500 py-2 text-center">No members yet</p>
        ) : (
          <ul className="space-y-3">
            {members.map((member) => (
              <li key={member.id} className="flex items-center justify-between p-2 hover:bg-gray-50 rounded">
                <div className="flex items-center">
                  {member.picture && (
                    <img 
                      src={member.picture} 
                      alt={member.name} 
                      className="w-8 h-8 rounded-full mr-2"
                    />
                  )}
                  <div>
                    <p className="font-medium">{member.name}</p>
                    <p className="text-xs text-gray-500">{member.email}</p>
                  </div>
                </div>
                <div className="flex items-center">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                    member.role === 'stakeholder' 
                      ? 'bg-purple-100 text-purple-800' 
                      : 'bg-blue-100 text-blue-800'
                  }`}>
                    {member.role}
                  </span>
                  <button
                    onClick={() => handleRemoveMember(member.id)}
                    className="ml-2 text-red-500 hover:text-red-700"
                    aria-label="Remove member"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default BoardMemberManagement;
