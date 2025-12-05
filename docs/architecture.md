# 🏗️ System Architecture

This document outlines the high-level architecture and data flow of **The TN Election Pulse**.

## High-Level Overview

The system follows a standard **Client-Server** architecture. The standard Go backend serves as a REST API that orchestrates data fetching and AI analysis, while the React frontend consumes these APIs to display the dashboard.

```mermaid
graph TD
    User[User / Browser] <-->|HTTP| FE[React Frontend]
    FE <-->|REST API| BE[Go Fiber Backend]
    
    subgraph "Backend Services"
        Orch[Orchestrator]
        News[News Service]
        YT[YouTube Service]
        Reddit[Reddit Service]
        AI[AI Service (Gemini)]
    end
    
    BE <--> Orch
    Orch -->|Concurrent Fetch| News
    Orch -->|Concurrent Fetch| YT
    Orch -->|Concurrent Fetch| Reddit
    
    News -->|RSS| GNews[Google News]
    YT -->|API| YApi[YouTube Data API]
    Reddit -->|JSON| RApi[Reddit Public API]
    
    Orch -->|Aggregated Text| AI
    AI -->|Prompt & Context| LLM[Google Gemini 1.5 Flash]
    
    BE <-->|Read/Write| DB[(PostgreSQL Database)]
```

## Components

### 1. Orchestrator (`orchestrator.go`)
*   **Role**: The central coordinator designed to handle data gathering efficiently.
*   **Mechanism**: Uses Go `sync.WaitGroup` to launch concurrent goroutines for each data source (News, YouTube, Reddit).
*   **Aggregation**: Collects results (or errors) from all sources and compiles a single "Corpus" string for the AI.

### 2. Data Services
*   **`news_service.go`**: Parses Google News RSS feeds for specific queries (e.g., "DMK Tamil Nadu").
*   **`youtube_service.go`**: Searches for recent videos and fetches top-level comments. Implements retry logic if comments are disabled.
*   **`reddit_service.go`**: Scrapes recent posts from target subreddits (`r/TamilNadu`, `r/India`) using the JSON API.

### 3. AI Service (`ai_service.go`)
*   **Role**: The intelligence layer.
*   **Input**: A raw text corpus of headlines and comments.
*   **Process**:
    1.  Constructs a prompt with strict guidelines (Bias Check, EQ, Fact-Check).
    2.  Sends to Google Gemini 1.5 Flash.
    3.  Validates and parses the JSON response.
*   **Output**: Structure containing Sentiment Score, Emotion, Key Topics, and Fact Check Notes.

### 4. Database Layer
*   Uses **GORM** for ORM capabilities.
*   **`Party` Model**: Static data about political parties (Name, Color).
*   **`SentimentSnapshot` Model**: Time-series record of each analysis run. Stores `KeyTopics` as JSONB for flexibility.

## Data Flow (Analysis Request)

1.  User clicks "Refresh" on the Dashboard.
2.  Frontend calls `POST /api/v1/analyze` with `{ party_id: 1 }`.
3.  Backend checks if a fresh snapshot exists (< 12 hours old).
    *   *If yes*: Returns cached data immediately.
    *   *If no*: Triggers the Orchestrator.
4.  Orchestrator concurrently fetches data from News, YouTube, Reddit.
5.  Aggregated text is sent to Gemini.
6.  Result is saved to DB as a new `SentimentSnapshot`.
7.  JSON response is sent back to Frontend.
