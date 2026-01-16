# Project Analysis: The TN Election Pulse

## Overview
This is a sophisticated real-time analytics dashboard for Tamil Nadu politics, predicting "Digital Momentum" of political parties using AI-powered sentiment analysis from multiple sources.

## Key Components
- **Backend**: Go with Fiber framework, concurrent data fetching from Google News, YouTube, Reddit
- **AI**: Google Gemini 1.5 Flash for sentiment analysis, slang detection, and topic extraction
- **Frontend**: React with Vite, plain CSS, Recharts for visualizations
- **Database**: PostgreSQL for storing historical snapshots

## Architecture Highlights
- Tri-source analysis: News (30%), YouTube (50%), Reddit (20%)
- Weighted scoring formula producing 0-100 momentum score
- Real-time caching to avoid redundant API calls
- Concurrent goroutines for efficient data collection

## Current Status
The project includes comprehensive documentation, backend services implementation, and basic frontend components. Ready for integration and testing with API keys.

## Agent Rules
The project includes AI agent rules for:
- React performance optimization
- LLM integration and prompt engineering
- PostgreSQL best practices
- React hooks best practices
- Go fundamentals

## Recommendations
- Obtain required API keys (Gemini, YouTube Data API)
- Test data collection services
- Implement error handling for API quotas
- Add authentication if needed for production