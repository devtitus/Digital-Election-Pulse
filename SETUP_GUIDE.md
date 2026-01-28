# 🚀 Setup Guide for The TN Election Pulse

This guide provides step-by-step instructions to set up and run The TN Election Pulse project on your local machine or using Docker.

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

### 🐹 Backend Requirements
- **Go 1.23+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL** - [Download PostgreSQL](https://www.postgresql.org/download/)

### 🎨 Frontend Requirements
- **Node.js 18+** - [Download Node.js](https://nodejs.org/)
- **npm** (comes with Node.js)

### 🐳 Docker Requirements (Optional)
- **Docker Desktop** - [Download Docker](https://www.docker.com/products/docker-desktop/)
- **Docker Compose** (usually comes with Docker Desktop)

### 🔑 API Keys Required
You'll need the following API keys:
- **Google Gemini API Key** - Get from [Google AI Studio](https://aistudio.google.com/)
- **YouTube Data API Key** - Get from [Google Cloud Console](https://console.cloud.google.com/)
- **NewsData API Key** (Optional) - Get from [NewsData.io](https://newsdata.io/)

## 📥 Installation Methods

Choose one of the following installation methods:

## Method 1: Local Development Setup

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/yourusername/election-pulse.git
cd election-pulse
```

### 2️⃣ Backend Setup

#### Install Go Dependencies

```bash
cd backend
# Install Go modules
go mod download
```

#### Configure Environment Variables

Create a `.env` file in the `backend` directory:

```bash
cp .env.example .env
```

Edit the `.env` file and add your API keys:

```env
DATABASE_URL="postgresql://username:password@localhost:5432/election_pulse?sslmode=disable"
GEMINI_API_KEY="your-gemini-api-key"
YOUTUBE_API_KEY="your-youtube-api-key"
NEWSDATA_API_KEY="your-newsdata-api-key"
```

#### Set Up PostgreSQL Database

1. Create a new database:
   ```sql
   CREATE DATABASE election_pulse;
   ```

2. Update the `DATABASE_URL` in your `.env` file with your PostgreSQL credentials.

#### Run the Backend Server

```bash
# Start the backend server
go run cmd/main.go
```

The backend will start on `http://localhost:3000`

### 3️⃣ Frontend Setup

#### Install Node.js Dependencies

```bash
cd ../frontend
npm install
```

#### Configure Frontend Environment

Create a `.env` file in the `frontend` directory:

```env
VITE_API_URL=http://localhost:3000/api/v1
```

#### Run the Frontend Development Server

```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

## Method 2: Docker Setup (Recommended)

### 1️⃣ Clone the Repository

```bash
git clone https://github.com/yourusername/election-pulse.git
cd election-pulse
```

### 2️⃣ Configure Environment Variables

#### Backend Configuration

Edit the `backend/.env` file:

```env
DATABASE_URL="postgresql://username:password@postgres:5432/election_pulse?sslmode=disable"
GEMINI_API_KEY="your-gemini-api-key"
YOUTUBE_API_KEY="your-youtube-api-key"
NEWSDATA_API_KEY="your-newsdata-api-key"
```

#### Frontend Configuration

Edit the `frontend/.env` file:

```env
VITE_API_URL=http://backend:3000/api/v1
```

### 3️⃣ Build and Run with Docker Compose

```bash
docker-compose up --build
```

This will:
1. Build the backend Docker image
2. Build the frontend Docker image
3. Start both containers
4. Set up the network between them

### 4️⃣ Access the Application

- **Frontend**: `http://localhost:8080`
- **Backend API**: `http://localhost:3000`

## 🧪 Testing the Setup

### Verify Backend is Running

```bash
curl http://localhost:3000/api/v1/parties
```

You should see a JSON response with the list of political parties.

### Verify Frontend is Running

Open your browser and navigate to:
- Local setup: `http://localhost:5173`
- Docker setup: `http://localhost:8080`

You should see the Election Pulse dashboard with party selection options.

## 🔧 Troubleshooting

### Common Issues and Solutions

#### Backend Issues

**Issue**: Backend fails to start with database connection error

**Solution**:
1. Verify PostgreSQL is running
2. Check your `DATABASE_URL` in `.env` file
3. Ensure the database exists and credentials are correct

**Issue**: API keys not working

**Solution**:
1. Verify API keys are correct
2. Check for typos in `.env` file
3. Ensure keys have proper permissions

#### Frontend Issues

**Issue**: Frontend shows "Failed to fetch" errors

**Solution**:
1. Verify backend is running
2. Check CORS settings in backend
3. Ensure `VITE_API_URL` points to correct backend URL

**Issue**: White screen or blank page

**Solution**:
1. Check browser console for errors
2. Verify all dependencies installed (`npm install`)
3. Clear browser cache and try again

#### Docker Issues

**Issue**: Docker containers fail to start

**Solution**:
1. Check Docker Desktop is running
2. Verify port conflicts (3000, 8080)
3. Check container logs: `docker-compose logs`

**Issue**: Database connection refused in Docker

**Solution**:
1. Use `postgres` as hostname in `DATABASE_URL`
2. Ensure PostgreSQL container is running
3. Check network connectivity between containers

## 📦 Project Structure

```
election-pulse/
├── backend/              # Go backend
│   ├── cmd/              # Main application
│   ├── db/               # Database connection
│   ├── handlers/         # API routes
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   └── Dockerfile        # Backend container
├── frontend/             # React frontend
│   ├── src/              # Source code
│   ├── public/           # Static assets
│   └── Dockerfile        # Frontend container
├── docs/                 # Documentation
├── docker-compose.yml    # Docker setup
└── README.md             # Project info
```

## 🚀 Running in Production

### For Local Setup

#### Backend
```bash
# Build backend
go build -o main cmd/main.go
# Run with PM2 or similar process manager
pm2 start main --name election-pulse-backend
```

#### Frontend
```bash
npm run build
# Serve built files with nginx, Apache, or similar
```

### For Docker Setup

```bash
# Run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop containers
docker-compose down
```

## 🔄 Updating the Project

### Pull Latest Changes
```bash
git pull origin main
```

### Update Dependencies

**Backend**:
```bash
cd backend
go get -u ./...
```

**Frontend**:
```bash
cd frontend
npm update
```

### Rebuild Docker Images
```bash
docker-compose build --no-cache
docker-compose up -d
```

## 📖 Additional Resources

- **API Documentation**: See `docs/api_reference.md`
- **Architecture**: See `docs/architecture.md`
- **Project Details**: See `PROJECTS.md`

## 🤝 Support

If you encounter issues not covered in this guide:
- Check the project documentation
- Open an issue on GitHub
- Review the API reference for technical details