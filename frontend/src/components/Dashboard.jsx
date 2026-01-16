import React, { useState, useEffect } from "react";
import PartySelector from "./PartySelector";
import SentimentWave from "./SentimentWave";
import TrendGraph from "./TrendGraph";
import ComparisonChart from "./ComparisonChart";
import {
  getParties,
  analyzeParty,
  getLatestSnapshot,
  getTrends,
  getComparison,
} from "../api/api";
import "./Dashboard.css";

const Dashboard = () => {
  const [parties, setParties] = useState([]);
  const [selectedParty, setSelectedParty] = useState(null);
  const [loading, setLoading] = useState(false);
  const [analysisResult, setAnalysisResult] = useState(null);
  const [trends, setTrends] = useState([]);
  const [comparison, setComparison] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchParties = async () => {
      const data = await getParties();
      setParties(data);
      if (data && data.length > 0) setSelectedParty(data[0]);
    };
    fetchParties();
  }, []);

  const fetchLatestData = async (party) => {
    setLoading(true);
    setError(null);
    try {
      const latest = await getLatestSnapshot(party.name);
      if (latest && latest.exists) {
        setAnalysisResult(latest);
      } else {
        await handleAnalyze(party);
      }
    } catch (err) {
      console.error("Error fetching latest data:", err);
      await handleAnalyze(party);
    } finally {
      setLoading(false);
    }
  };

  const handleAnalyze = async (party) => {
    setLoading(true);
    setError(null);
    try {
      const result = await analyzeParty(party.name);
      setAnalysisResult(result);
      // Fetch additional data
      const trendData = await getTrends(party.name);
      setTrends(trendData);
      const comparisonData = await getComparison();
      setComparison(comparisonData);
    } catch (err) {
      console.error(err);
      setError("Analysis failed. Please check backend connection.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (selectedParty) {
      fetchLatestData(selectedParty);
    }
  }, [selectedParty]);

  return (
    <div className="dashboard-container">
      <header className="dashboard-header animate-fade-in-up">
        <h1 className="pulse-title">
          <span className="title-wave">Election</span>
          <span className="title-wave delay-1">Pulse</span>
        </h1>
        <p>AI-Powered Political Sentiment Analysis</p>
      </header>

      {/* Creative Magazine-Style Layout */}
      <div className="magazine-layout">
        {/* Party Banner Section */}
        <div className="magazine-section party-banner">
          <PartySelector
            parties={parties}
            selectedPartyId={selectedParty?.id}
            onSelect={setSelectedParty}
          />

          <div className="controls">
            {selectedParty && (
              <button
                className="analyze-btn"
                onClick={() => handleAnalyze(selectedParty)}
                disabled={loading}
              >
                {loading ? "Analyzing..." : "🔄 Refresh Analysis"}
              </button>
            )}
          </div>
        </div>

        {/* Hero Section - Large Feature */}
        {analysisResult && (
          <div className="magazine-section hero-feature">
            <div className="feature-card pulse-feature">
              <div className="feature-header">
                <h3>Digital Sentiment Pulse</h3>
                <div className="pulse-indicator">
                  <div
                    className="pulse-dot animate-pulse-subtle"
                    style={{
                      backgroundColor: selectedParty?.color_hex || 'var(--accent-gold)',
                      boxShadow: `0 0 20px ${selectedParty?.glow || 'rgba(251, 191, 36, 0.3)'}`
                    }}
                  ></div>
                  <span className="pulse-label">Live Data</span>
                </div>
              </div>
              <SentimentWave
                score={analysisResult.sentiment_score}
                party={selectedParty}
              />
            </div>
          </div>
        )}

        {/* Content Grid - Magazine Style */}
        <div className="magazine-grid">
          {/* Left Column */}
          <div className="magazine-column left-column">
            {/* Emotion Card */}
            <div className="magazine-card emotion-spotlight">
              <div className="card-content">
                <h3>Dominant Emotion</h3>
                {analysisResult?.emotion ? (
                  <div className="emotion-display">
                    <span className="emotion-text">{analysisResult.emotion}</span>
                    <div className="emotion-pulse animate-wave-gentle"></div>
                  </div>
                ) : (
                  <div className="empty-particle">
                    <span>😐</span>
                    <p>Neutral / Unclear</p>
                  </div>
                )}
              </div>
            </div>

            {/* Key Topics Card - Fixed Height with Scroll */}
            <div className="magazine-card topics-showcase">
              <div className="card-content">
                <h3>Key Topics</h3>
                <div className="topics-container">
                  {analysisResult?.key_topics &&
                  analysisResult.key_topics.length > 0 ? (
                    <div className="topics-scroll">
                      {analysisResult.key_topics.map((topic, index) => (
                        <span
                          key={index}
                          className="topic-particle animate-pulse-subtle"
                          style={{ animationDelay: `${index * 0.1}s` }}
                        >
                          {topic}
                        </span>
                      ))}
                    </div>
                  ) : (
                    <div className="empty-particle">
                      <span>💭</span>
                      <p>No specific topics found</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Right Column */}
          <div className="magazine-column right-column">
            {/* Party Comparison */}
            <div className="magazine-card comparison-highlight">
              <div className="card-content">
                <h3>Party Comparison</h3>
                <ComparisonChart data={comparison} selectedParty={selectedParty} />
              </div>
            </div>

            {/* Trend Analysis */}
            <div className="magazine-card trend-insight">
              <div className="card-content">
                <h3>Sentiment Trends</h3>
                <TrendGraph data={trends} party={selectedParty} />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Loading States */}
      {loading && !analysisResult && (
        <div className="pulse-loading">
          <div className="loading-pulse animate-pulse-subtle">
            <span>📊</span>
            <p>
              Crunching numbers for {selectedParty?.name || "election"}...
            </p>
          </div>
        </div>
      )}

      {error && (
        <div className="pulse-error">
          <div className="error-particle">
            <span>⚠️</span>
            <h3>Oops! Something went wrong.</h3>
            <p>{error}</p>
            <button
              className="retry-pulse"
              onClick={() => handleAnalyze(selectedParty)}
            >
              Try Again
            </button>
          </div>
        </div>
      )}

      {!loading && !error && !analysisResult && (
        <div className="pulse-empty">
          <div className="empty-pulse">
            <span>📈</span>
            <h3>No Data Yet</h3>
            <p>Select a party to view their digital momentum.</p>
            <button
              className="start-pulse"
              onClick={() => handleAnalyze(selectedParty)}
            >
              Start Analysis
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
