import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

// Components
import ProtectedRoute from './components/ProtectedRoute';

// Pages
import SignIn from './pages/SignIn';
import Register from './pages/Register';
import Boards from './pages/Boards';
import BoardDetails from './pages/BoardDetails';
import Forbidden from './pages/Forbidden';

const AppRouter: React.FC = () => {
  return (
    <Router>
      <Routes>
        {/* Public routes */}
        <Route path="/signin" element={<SignIn />} />
        <Route path="/register" element={<Register />} />
        <Route path="/forbidden" element={<Forbidden />} />
        
        {/* Protected routes - require authentication */}
        <Route 
          path="/boards" 
          element={
            <ProtectedRoute>
              <Boards />
            </ProtectedRoute>
          } 
        />
        
        {/* Board details - require board access */}
        <Route 
          path="/boards/:id" 
          element={
            <ProtectedRoute>
              <BoardDetails />
            </ProtectedRoute>
          } 
        />
        
        {/* Admin routes - require app_admin role */}
        <Route 
          path="/admin/*" 
          element={
            <ProtectedRoute requiredRole="app_admin">
              <div>Admin Panel (To be implemented)</div>
            </ProtectedRoute>
          } 
        />
        
        {/* Default route */}
        <Route path="/" element={<Navigate to="/boards" />} />
        <Route path="*" element={<Navigate to="/boards" />} />
      </Routes>
    </Router>
  );
};

export default AppRouter;
