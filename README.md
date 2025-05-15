# Feedback Management System

This project is a feedback management system inspired by [Canny.io](https://canny.io/). It is built using React for the frontend and Go for the backend.

## Features
- Google SSO for sign-up and log-in.
- Sign-in and Register pages.

## Getting Started

### Prerequisites
- Node.js and npm installed.
- Go installed.

### Installation
1. Clone the repository.
2. Navigate to the project directory.
3. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   ```
4. Install backend dependencies:
   ```bash
   cd backend
   go mod tidy
   ```

### Running the Application
1. Start the backend server:
   ```bash
   cd backend
   go run main.go
   ```
2. Start the frontend development server:
   ```bash
   cd frontend
   npm start
   ```

### Folder Structure
- `frontend/`: React application.
- `backend/`: Go application.

## License
This project is licensed under the MIT License.
