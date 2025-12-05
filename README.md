# 🗳️ The TN Election Pulse

**Real-Time Digital Momentum Analytics for Tamil Nadu Politics**

A sophisticated analytics dashboard that predicts the "Digital Momentum" of political parties in Tamil Nadu by analyzing sentiment from multiple data sources in real-time using Generative AI.

![Dashboard Preview](docs/dashboard-preview.png) *(Placeholder for screenshot)*

## 🚀 Features

*   **Multi-Source Intelligence**: Aggregates data from **Google News**, **YouTube Comments**, and **Reddit** (`r/TamilNadu`, `r/Chennai`).
*   **AI-Powered Analysis**: Uses **Google Gemini 1.5 Flash** to perform:
    *   **Sentiment Analysis**: Positive/Negative scoring (-100 to +100).
    *   **Emotional Quotient (EQ)**: Identifies deep drivers like Anger, Hope, or Mockery.
    *   **Fact-Checking**: Cross-references claims to flag potential misinformation.
    *   **Bias Mitigation**: Audited prompts to ensure objective analysis.
*   **Modern UI**: A "Deep Dark" aesthetic with **Glassmorphism**, Aurora gradients, and smooth micro-interactions.
*   **Real-Time Caching**: Persists snapshots to PostgreSQL/SQLite to prevent redundant expensive API calls.

## 🛠️ Tech Stack

### Backend
*   **Language**: Go (Golang) 1.23+
*   **Framework**: Fiber (Fast HTTP web framework)
*   **Database**: PostgreSQL (via GORM)
*   **AI**: Google Gemini SDK (`google-generative-ai-go`)
*   **APIs**: YouTube Data API v3, Google News RSS, Reddit JSON API

### Frontend
*   **Framework**: React (Vite)
*   **Styling**: Vanilla CSS (Variables, Flexbox/Grid)
*   **Theme**: Custom "Dark Aurora" Theme with Glassmorphism
*   **Visualization**: Recharts (for gauges and charts)

## 📦 Installation

### Prerequisites
1.  **Go** installed (v1.21+).
2.  **Node.js** installed (v18+).
3.  **PostgreSQL** database running (or adjust `database.go` for SQLite).
4.  **API Keys**:
    *   `GEMINI_API_KEY`: Get from Google AI Studio.
    *   `YOUTUBE_API_KEY`: Get from Google Cloud Console.

### Setup

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/yourusername/election-pulse.git
    cd election-pulse
    ```

2.  **Backend Setup**
    ```bash
    cd backend
    # Create .env file
    cp .env.example .env
    # Update .env with your keys and DB URL
    
    # Install dependencies
    go mod download
    
    # Run the server
    go run cmd/main.go
    ```
    *Server runs on `http://localhost:8080`*

3.  **Frontend Setup**
    ```bash
    cd frontend
    # Install dependencies
    npm install
    
    # Run dev server
    npm run dev
    ```
    *App runs on `http://localhost:5173`*

### 🐳 Docker Setup (Recommended)

Run the entire application (Backend + Frontend) with a single command.

1.  **Ensure Docker Desktop is running.**
2.  **Run Docker Compose:**
    ```bash
    docker-compose up --build
    ```
3.  **Access the App:**
    *   Frontend: `http://localhost:8080`
    *   Backend API: `http://localhost:3000`

## 📖 Documentation

Detailed technical documentation can be found in the `docs/` folder:

*   [**Architecture**](docs/architecture.md): System design, data flow, and service breakdown.
*   [**API Reference**](docs/api_reference.md): Backend endpoints and payloads.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## 📜 License

MIT License.
