import React, { useState, useEffect } from 'react';
import PartySelector from './PartySelector';
import SentimentGauge from './SentimentGauge';
import { getParties, analyzeParty, getLatestSnapshot } from '../api/api';
import './Dashboard.css';

const Dashboard = () => {
    const [parties, setParties] = useState([]);
    const [selectedParty, setSelectedParty] = useState(null);
    const [loading, setLoading] = useState(false);
    const [analysisResult, setAnalysisResult] = useState(null);
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
            // Try to get cached/latest data first
            const latest = await getLatestSnapshot(party.name);

            if (latest && latest.exists) {
                setAnalysisResult(latest);
            } else {
                // If no data exists, trigger fresh analysis automatically (optional, or ask user?)
                // For now, let's trigger it automatically if it's the first time
                await handleAnalyze(party);
            }
        } catch (err) {
            console.error("Error fetching latest data:", err);
            // If error (e.g. 404), maybe try fresh analysis?
            // await handleAnalyze(party);
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
            <header className="dashboard-header">
                <h1>🗳️ TN Election Pulse</h1>
                <p>AI-Powered Political Sentiment Analysis</p>
            </header>

            <main>
                <PartySelector
                    parties={parties}
                    selectedPartyId={selectedParty?.id}
                    onSelect={setSelectedParty}
                />

                <div className="controls">
                    {/* Add Refresh Button */}
                    {selectedParty && (
                        <button
                            className="refresh-btn"
                            onClick={() => handleAnalyze(selectedParty)}
                            disabled={loading}
                        >
                            {loading ? 'Analyzing...' : '🔄 Refresh Analysis'}
                        </button>
                    )}
                </div>

                {/* Show loading overlay if data exists, or full loader if not */}
                {/* Show loading overlay if data exists, or full loader if not */}
                {loading && !analysisResult && (
                    <div className="loading-state">
                        <div className="spinner">🗳️</div>
                        <p>Crunching numbers for {selectedParty?.name || 'election'}...</p>
                    </div>
                )}

                {error && (
                    <div className="error-state">
                        <div className="error-icon">⚠️</div>
                        <h3>Oops! Something went wrong.</h3>
                        <p>{error}</p>
                        <button className="retry-btn" onClick={() => handleAnalyze(selectedParty)}>
                            Try Again
                        </button>
                    </div>
                )}

                {!loading && !error && !analysisResult && (
                    <div className="empty-state">
                        <div className="empty-icon">📊</div>
                        <h3>No Data Yet</h3>
                        <p>Select a party to view their digital momentum.</p>
                        <button className="start-btn" onClick={() => handleAnalyze(selectedParty)}>
                            Start Analysis
                        </button>
                    </div>
                )}

                {analysisResult && (
                    <div className={`dashboard-grid ${loading ? 'updating' : ''}`}>
                        {loading && (
                            <div className="updating-overlay">
                                <span>Refreshing data...</span>
                            </div>
                        )}
                        <SentimentGauge score={analysisResult.sentiment_score} />

                        <div className="card topics-card">
                            <h3>Key Topics</h3>
                            <div className="topics-list">
                                {analysisResult.key_topics && analysisResult.key_topics.length > 0 ? (
                                    analysisResult.key_topics.map((topic, index) => (
                                        <span key={index} className="topic-tag">{topic}</span>
                                    ))
                                ) : (
                                    <div className="mini-empty-state">
                                        <span className="mini-icon">💭</span>
                                        <p>No specific topics found</p>
                                    </div>
                                )}
                            </div>
                        </div>

                        <div className="card emotion-card">
                            <h3>Dominant Emotion</h3>
                            {analysisResult.emotion ? (
                                <div className="emotion-value">{analysisResult.emotion}</div>
                            ) : (
                                <div className="mini-empty-state">
                                    <span className="mini-icon">😐</span>
                                    <p>Neutral / Unclear</p>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
};

export default Dashboard;
