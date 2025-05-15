# Feedback Management System

This project is a feedback management system inspired by [Canny.io](https://canny.io/). It is built using React for the frontend and Go for the backend.

## Features
- Google SSO for sign-up and log-in.
- Sign-in and Register pages.

## Getting Started

### Prerequisites
- Node.js and npm installed.
- Go installed.
- Docker and Docker Compose (optional, for containerized setup).

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
#### Option 1: Local Development
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

#### Option 2: Using Docker
1. Make sure you have Docker and Docker Compose installed.
2. Build and start all the containers using Docker Compose:
   ```bash
   docker-compose up --build
   ```
3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
4. To stop the containers, press `Ctrl+C` or run:
   ```bash
   docker-compose down
   ```

### Docker Container Details
- **Frontend Container**: Built with Node.js and served via Nginx, accessible on port 3000
- **Backend Container**: Built with Go, accessible on port 8080
- **Database Container**: PostgreSQL 13, data is persisted in a Docker volume

### Folder Structure
- `frontend/`: React application.
- `backend/`: Go application.

## License
This project is licensed under the MIT License.
