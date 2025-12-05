# 📘 Project Plan: The TN Election Pulse (AI-Powered Forecasting Engine)
**Version**: 1.0  
**Domain**: Political Intelligence / Predictive Analytics  
**Target Region**: Tamil Nadu, India  
**Goal**: To analyze real-time sentiment from Social Media & News to predict the "Digital Momentum" of political parties using AI.

## 1. 🏗️ High-Level Architecture
The system follows a Event-Driven Architecture where a user request triggers a concurrent data fetching process, analyzes it via LLM, and stores the snapshot for historical trending.

```mermaid
graph TD
    User[User Client (React)] -->|HTTP Request| API[Backend API (Go / Fiber)]
    API -->|Spawns Goroutines| Worker1[Worker 1: Google News RSS Scraper]
    API -->|Spawns Goroutines| Worker2[Worker 2: YouTube Data API]
    API -->|Spawns Goroutines| Worker3[Worker 3: Reddit Scraper]
    
    Worker1 --> Aggregator[Data Aggregator (Go Structs)]
    Worker2 --> Aggregator
    Worker3 --> Aggregator
    
    Aggregator -->|Raw Text Payload| AI[AI Engine (Google Gemini 1.5 Flash)]
    AI -->|JSON Analysis: Sentiment + Key Issues| Core[Logic Core (Go)]
    Core -->|Calculates Momentum Score| DB[(Database (PostgreSQL))]
    DB -->|Response| User
```

## 2. 🛠️ Tech Stack & Tools

### Frontend (Client Side)
- **Framework**: React.js (Vite Build Tool)
- **Language**: JavaScript (ES6+) or TypeScript
- **Styling**: Plain CSS (CSS Modules) using Flexbox/Grid. No Tailwind.
- **State Management**: React Context API or simply useState (keep it simple).
- **HTTP Client**: Axios.
- **Visualization**: Recharts (for Gauge charts and Trend lines).

### Backend (Server Side)
- **Language**: Go (Golang) - Chosen for concurrency.
- **Web Framework**: Go Fiber (Extremely fast, Express-like syntax).
- **Environment Config**: godotenv.
- **Data Parsing**: gofeed (RSS Parser).

### AI & Data
- **LLM**: Google Gemini 1.5 Flash (via google-generativeai Go SDK).
- **Primary Database**: PostgreSQL (Hosted on Neon.tech or Local).
- **External APIs**:
    - Google News (RSS)
    - YouTube Data API v3
    - Reddit JSON (via PRAW or direct endpoint).

## 3. 💾 Data Sources & Extraction Strategy

### A. Google News (Official Narrative)
- **Method**: RSS Parsing (No API key needed).
- **Target URL**: `https://news.google.com/rss/search?q={QUERY}&hl=ta&gl=IN&ceid=IN:ta`
- **Query Logic**: `"{Party Name}" AND "Tamil Nadu"`
- **Extraction**: Get Top 10 Headlines.

### B. YouTube (The Mass Pulse)
- **Method**: YouTube Data API v3.
- **Constraint**: Daily Quota is 10,000 units.
    - Search = 100 units.
    - Comments = 1 unit.
- **Optimization**: Strict Caching. Do not fetch live data more than once per hour per party.
- **Target**: Search for `"{Party Name} speech"` or `"{Party Name} news"` sorted by date (Last 24h).
- **Data Points**: Extract Top 50 comments from the top video.

### C. Reddit (Urban Sentiment)
- **Method**: Direct JSON scraping or Go wrapper.
- **Target**: r/TamilNadu, r/Chennai.
- **Filter**: Search within subreddit for Party Name in last 24h.

## 4. 🧮 Algorithms & Logic

### A. The Gemini Prompt (System Instruction)
We must instruct the AI to handle Tamil slang and political context.

**Prompt**: 
> "You are a Tamil Nadu political analyst. I will give you a list of headlines and comments regarding the party '{PARTY_NAME}'.
> Filter out spam.
> Analyze sentiment considering local slang (e.g., 'Sanghi', 'Upee', '200 rs', 'Vadinokkan').
> Output a JSON object with:
> sentiment_score: float (-1.0 to 1.0)
> emotion: string (Anger, Hope, Mockery, Support)
> key_topics: array of strings (e.g., ['NEET', 'Flood Relief'])"

### B. The Momentum Formula
We combine the raw sentiment with a "Volume Multiplier".

```go
// Go Logic Draft
const (
    WeightNews    = 0.3
    WeightYouTube = 0.5 // Higher because it represents mass voters
    WeightReddit  = 0.2
)

RawScore = (NewsSent * WeightNews) + (YTSent * WeightYouTube) + (RedditSent * WeightReddit)

// Normalize to 0-100 scale
WinningProbability = 50 + (RawScore * 50)
```

## 5. 🗄️ Database Schema (PostgreSQL)
We need two tables.

### Table 1: parties
| Column | Type | Description |
| :--- | :--- | :--- |
| id | SERIAL (PK) | 1, 2, 3... |
| name | VARCHAR | "DMK", "AIADMK", "TVK" |
| leader | VARCHAR | "M.K. Stalin", "Edappadi Palaniswami" |
| color_hex | VARCHAR | "#FF0000" (for UI) |

### Table 2: sentiment_snapshots
| Column | Type | Description |
| :--- | :--- | :--- |
| id | SERIAL (PK) | |
| party_id | INT (FK) | Links to parties |
| score | FLOAT | The calculated 0-100 score |
| key_issue | TEXT | "Electricity Bill High" |
| source_breakdown| JSONB | `{"yt": 0.8, "news": -0.2}` |
| created_at | TIMESTAMP | To plot the line graph |

## 6. 🔌 API Endpoints Specification
**Group**: `/api/v1`

### 1. GET /parties
Returns list of supported parties (DMK, AIADMK, BJP, NTK, TVK).

### 2. POST /analyze
**Body**: `{"party_name": "TVK"}`
**Logic**:
- **Check DB**: Is there a snapshot for this party created in the last 1 hour?
- **If Yes**: Return DB data (Cache Hit).
- **If No**: Trigger Scrapers -> Call Gemini -> Save to DB -> Return Data.

### 3. GET /history/:party_id
Returns last 7 days of scores for the Line Chart.

## 7. 📅 Implementation Roadmap

### Phase 1: Backend Foundation (Go)
- [ ] Initialize Go module (`go mod init`).
- [ ] Set up Fiber server and Hello World route.
- [ ] Connect PostgreSQL using pgx or gorm.
- [ ] Create the Database Tables (SQL Scripts).

### Phase 2: The Data Connectors (Go)
- [ ] Implement `services/news_service.go` (RSS Parsing).
- [ ] Implement `services/youtube_service.go` (API Call).
- [ ] Create the Concurrency Manager (WaitGroups) to call these in parallel.

### Phase 3: The Brain (Gemini Integration)
- [ ] Get API Key from Google AI Studio.
- [ ] Implement `services/ai_service.go`.
- [ ] Write the "Prompt Engineering" logic to format the text sent to Gemini.

### Phase 4: Frontend Development (React + Plain CSS)
- [ ] Create React App (`npm create vite@latest`).
- [ ] CSS Architecture:
    - `src/styles/global.css` (Variables, Reset).
    - `src/components/Dashboard.css` (Grid layout).
    - `src/components/PartyCard.css` (Specific styling).
- [ ] Build Components:
    - `Header.jsx`
    - `PartySelector.jsx`
    - `SentimentGauge.jsx` (Using Recharts).
    - `TopicWordCloud.jsx`.

### Phase 5: Integration & Polish
- [ ] Connect React to Go API using Axios.
- [ ] Handle Loading States (Spinners while AI analyzes).
- [ ] **Crucial**: Add Error Handling (e.g., "YouTube Quota Exceeded").

## 8. 📂 Folder Structure
```text
/tn-election-pulse
├── /backend (Go)
│   ├── /cmd
│   │   └── main.go           # Entry point
│   ├── /handlers             # API Route Logic
│   ├── /models               # DB Structs
│   ├── /services             # External Logic (AI, YT, News)
│   │   ├── ai_service.go
│   │   ├── youtube_service.go
│   │   └── news_service.go
│   ├── /db                   # Database connection
│   └── .env                  # API Keys
│
├── /frontend (React)
│   ├── /src
│   │   ├── /api              # Axios calls
│   │   ├── /components       # React Components
│   │   ├── /styles           # Plain CSS files
│   │   ├── App.jsx
│   │   └── main.jsx
│   └── package.json
│
└── PROJECT_PLAN.md
```

## 9. 🎨 CSS Visual Plan (No Tailwind)
Since we are using plain CSS, we will define a color palette variable system in `global.css` to keep it looking modern.

```css
/* src/styles/global.css */
:root {
  --bg-color: #f4f6f8;
  --card-bg: #ffffff;
  --text-primary: #1a202c;
  --dmk-color: #dd2e44;
  --aiadmk-color: #27ae60;
  --tvk-color: #f1c40f;
  --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}

body {
  background-color: var(--bg-color);
  font-family: 'Inter', sans-serif;
}
```

We will use CSS Grid for the Dashboard layout:

```css
/* Dashboard.css */
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
  padding: 20px;
}
```

## 10. ✅ Checklist Before Starting
- [ ] **Get API Keys**:
    - Google Gemini API Key.
    - Google Cloud Console (YouTube Data API v3).
- [ ] **Install Go**: Verify with `go version`.
- [ ] **Install Node**: Verify with `node -v`.
- [ ] **Setup DB**: Install Postgres or use a Docker container.

This document covers every tiny piece needed to build the project. You can now start coding Phase 1!
