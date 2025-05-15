import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { authService } from '../services/authService';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string;
  requiredBoardRole?: {
    boardId: number;
    role: string;
  };
}

/**
 * A wrapper for routes that require authentication and specific roles
 */
const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiredRole,
  requiredBoardRole,
}) => {
  const location = useLocation();
  const isAuthenticated = authService.isAuthenticated();
  
  // Check if user is authenticated
  if (!isAuthenticated) {
    // Redirect to sign in page, but save the current location they were trying to access
    return <Navigate to="/signin" state={{ from: location }} replace />;
  }
  
  // Check for required global role if specified
  if (requiredRole) {
    const userRole = authService.getUserRole();
    
    // For app_admin role requirement
    if (requiredRole === 'app_admin' && userRole !== 'app_admin') {
      return <Navigate to="/forbidden" replace />;
    }
    
    // For stakeholder role requirement (either app_admin or stakeholder)
    if (requiredRole === 'stakeholder' && 
        userRole !== 'app_admin' && 
        userRole !== 'stakeholder') {
      return <Navigate to="/forbidden" replace />;
    }
  }
  
  // Check for required board-specific role if specified
  if (requiredBoardRole) {
    const { boardId, role } = requiredBoardRole;
    
    // For board stakeholder requirement
    if (role === 'stakeholder') {
      if (!authService.isBoardStakeholder(boardId)) {
        return <Navigate to="/forbidden" replace />;
      }
    } 
    // For any board access requirement
    else if (!authService.hasBoardAccess(boardId)) {
      return <Navigate to="/forbidden" replace />;
    }
  }
  
  // If all checks pass, render the children
  return <>{children}</>;
};

export default ProtectedRoute;
