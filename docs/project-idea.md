# 🗳️ Project Title: The TN Election Pulse (AI-Powered Forecasting Engine)

## Tagline
Decodig the gap between Social Media "Hype" and Electoral "Reality" in Tamil Nadu Politics using Generative AI.

## 1. 📖 Executive Summary
In Tamil Nadu, politics is highly digital. With the entry of new players (like TVK) and the dominance of legacy parties (DMK, AIADMK), social media is flooded with noise. However, Twitter trends rarely translate to vote shares.

The TN Election Pulse is a real-time analytics dashboard that predicts the "Digital Momentum" of political parties. Unlike traditional opinion polls which are expensive and slow, this system uses Golang to concurrently harvest data from diverse sources (News, YouTube, Reddit) and uses Google Gemini 1.5 Flash (LLM) to linguistically analyze Tamil and "Tanglish" sentiment. It calculates a weighted "Win Probability" score, helping analysts differentiate between a "Bot Army" and genuine voter enthusiasm.

## 2. 🚩 The Problem Statement
- **The Echo Chamber Effect**: Social media platforms like Reddit often represent urban, English-speaking demographics, ignoring the rural voter base.
- **Linguistic Complexity**: TN politics involves deep local slang (e.g., "Sanghi", "Upee", "200 Rupees", "Vadinokkan") which standard NLP libraries (like NLTK or TextBlob) fail to understand.
- **Data Overload**: There is too much noise on YouTube comments and News feeds for a human to track manually.

## 3. 💡 The Solution: How It Works
The system operates on a "Tri-Source" Analysis Engine:

1. **The "Official" Narrative (Google News)**: Scrapes headlines to understand how the media portrays a party.
2. **The "Mass" Pulse (YouTube)**: Fetches the top 50 comments from trending political videos. In TN, YouTube comments are the closest proxy to the "Tea Shop" conversation of the common voter.
3. **The "Elite" Opinion (Reddit)**: Analyzes subreddit discussions for urban sentiment.

**The AI Core**: The aggregated text is fed into Google Gemini 1.5 Flash with a specific system prompt to decode Tamil slang and sarcasm.

### The Math
The system applies a weighted formula to calculate the Digital Momentum Score (DMS):

$$
DMS = (News \times 0.3) + (Reddit \times 0.2) + (YouTube \times 0.5)
$$

(YouTube is weighted higher to reflect the larger voter base).

## 4. 🛠️ Tech Stack & Architecture

### Frontend (The User Interface)
- **Library**: React.js (Vite)
- **Styling**: Plain CSS 3 (CSS Modules & Grid Layout) - No frameworks, pure performance.
- **Visualization**: Recharts (for Gauge Charts and Trend Lines).
- **State Management**: React Context API.

### Backend (The Engine)
- **Language**: Go (Golang)
- **Framework**: Fiber (Express-style, high-performance web framework).
- **Concurrency**: Uses Goroutines and sync.WaitGroup to fetch data from YouTube, News, and Reddit simultaneously, reducing latency by 70%.

### Artificial Intelligence
- **Model**: Google Gemini 1.5 Flash (via Google GenAI SDK).
- **Role**: Sentiment Classification, Slang Detection, and Topic Extraction.

### Data & Infrastructure
- **Database**: PostgreSQL (Stores historical scores for trend analysis).
- **APIs**: YouTube Data API v3, Google News RSS, PRAW (Reddit).

## 5. ✨ Key Features

### 🔹 1. "Tanglish" Semantic Analysis
Standard tools read "Thalaiva mass" as neutral text. Our Gemini integration identifies this as High Positive Sentiment (Hero Worship). It detects sarcasm and local political slurs used in Tamil Nadu context.

### 🔹 2. The "Real-Time" Gauge
A speedometer-style visualization that moves instantly based on live data.
- **0-30**: Losing Relevance (Negative Sentiment).
- **31-60**: Neutral / Standard Opposition.
- **61-100**: Wave Election (High Positive Momentum).

### 🔹 3. "Key Issue" Word Cloud
Instead of just a score, the AI extracts why people are happy or angry.
- **Example Output**: "Anger points: Electricity Bill, Flood Relief." | "Support points: Women's Rights Fund."

### 🔹 4. Competitor Comparison
Users can select "DMK vs. TVK" to see a side-by-side comparison of their digital footprint and sentiment share.

## 6. 🚀 Future Scope
- **WhatsApp Forensics**: Integrating a module to analyze forwarded messages (the "WhatsApp University" factor).
- **Geo-Tagging**: Mapping sentiment to specific districts (e.g., "Madurai vs. Chennai sentiment").
- **DeepFake Detector**: Adding an image processing layer to flag AI-generated political images.

## 7. 📸 Project Modules (Mental Model)
- **Scraper Service (Go)**: The worker bots that go out to the internet to fetch text.
- **Analyst Service (Gemini)**: The brain that reads the text and assigns scores.
- **Dashboard (React)**: The control center where the user views the "Election Pulse".
