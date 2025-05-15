# Docker Setup for Canny Clone

This document provides instructions for running the Canny Clone application using Docker.

## Prerequisites

- Docker installed on your machine
- Docker Compose installed on your machine

## Getting Started

1. Make sure you're in the root directory of the project.

2. Build and start all the containers using Docker Compose:

```bash
docker-compose up --build
```

This command will:
- Build the frontend container
- Build the backend container
- Start the PostgreSQL database container
- Setup the network between the containers

3. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

4. To stop the containers, press `Ctrl+C` or run:

```bash
docker-compose down
```

## Container Details

### Frontend Container
- Built with Node.js and served via Nginx
- Accessible on port 3000
- React application with TypeScript

### Backend Container
- Built with Go
- Accessible on port 8080
- Connected to the PostgreSQL database

### Database Container
- PostgreSQL 13
- Data is persisted in a Docker volume
- Database migrations are automatically run at startup

## Configuration

- The backend uses the `config.docker.json` configuration file
- Database credentials:
  - Username: user
  - Password: password
  - Database name: canny_clone

## Development Notes

- To make changes to the frontend or backend, rebuild the containers using `docker-compose up --build`
- Database data is persisted even after containers are removed, thanks to the Docker volume
- The backend automatically runs migration scripts located in the `backend/migrations` directory
