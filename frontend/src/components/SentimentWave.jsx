import React from 'react';
import './SentimentWave.css';

const SentimentWave = ({ score, party }) => {
    // Calculate wave properties based on score
    const normalizedScore = score / 100;
    const waveHeight = 20 + (normalizedScore * 40); // 20-60px height
    const waveFrequency = 0.02 + (normalizedScore * 0.04); // Slower waves for higher scores

    // Generate SVG path for the wave
    const generateWavePath = () => {
        const width = 300;
        const height = 80;
        const points = [];
        for (let x = 0; x <= width; x += 5) {
            const y = height/2 + Math.sin(x * waveFrequency) * waveHeight * Math.sin(normalizedScore * Math.PI);
            points.push(`${x},${y}`);
        }
        return `M0,${height} L${points.join(' L')} L${width},${height} Z`;
    };

    // Determine color based on score and party
    const getWaveColor = () => {
        // Use party color if available, otherwise fallback to sentiment-based colors
        if (party?.color_hex) {
            return party.color_hex;
        }
        if (score < 30) return 'var(--sentiment-negative)';
        if (score < 60) return 'var(--sentiment-neutral)';
        return 'var(--sentiment-positive)';
    };

    const getGlowColor = () => {
        // Use party glow if available, otherwise fallback to sentiment-based glows
        if (party?.glow) {
            return party.glow;
        }
        if (score < 30) return 'rgba(220, 38, 38, 0.3)';
        if (score < 60) return 'rgba(217, 119, 6, 0.3)';
        return 'rgba(22, 163, 74, 0.3)';
    };

    return (
        <div className="sentiment-wave">
            <h3>Digital Sentiment Pulse</h3>
            <div className="wave-container">
                <svg width="300" height="80" viewBox="0 0 300 80" className="wave-svg">
                    <defs>
                        <linearGradient id={`waveGradient-${party?.id}`} x1="0%" y1="0%" x2="100%" y2="0%">
                            <stop offset="0%" stopColor={getWaveColor()} stopOpacity="0.8" />
                            <stop offset="50%" stopColor={getWaveColor()} stopOpacity="1" />
                            <stop offset="100%" stopColor={getWaveColor()} stopOpacity="0.8" />
                        </linearGradient>
                        <filter id={`waveGlow-${party?.id}`}>
                            <feGaussianBlur stdDeviation="3" result="coloredBlur"/>
                            <feMerge>
                                <feMergeNode in="coloredBlur"/>
                                <feMergeNode in="SourceGraphic"/>
                            </feMerge>
                        </filter>
                    </defs>
                    <path
                        d={generateWavePath()}
                        fill={`url(#waveGradient-${party?.id})`}
                        filter={`url(#waveGlow-${party?.id})`}
                        className="wave-path animate-wave-gentle"
                    />
                </svg>
                <div className="score-display">
                    <span
                        className="score-number"
                        data-score-range={score < 30 ? 'negative' : score < 60 ? 'neutral' : 'positive'}
                        style={{ color: getWaveColor() }}
                    >
                        {score}
                    </span>
                    <span className="score-label">/ 100</span>
                </div>
            </div>
            <div className="pulse-indicator">
                <div
                    className="pulse-dot animate-pulse-subtle"
                    style={{ backgroundColor: getWaveColor(), boxShadow: `0 0 20px ${getGlowColor()}` }}
                ></div>
                <span className="pulse-label">Live Data</span>
            </div>
        </div>
    );
};

export default SentimentWave;