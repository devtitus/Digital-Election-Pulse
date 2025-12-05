import React from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer } from 'recharts';
import './SentimentGauge.css';

const SentimentGauge = ({ score }) => {
    const chartData = [
        { name: 'Score', value: score },
        { name: 'Remaining', value: 100 - score },
    ];

    let color = '#718096';
    if (score < 30) color = '#e74c3c'; // Negative
    else if (score < 60) color = '#f1c40f'; // Neutral
    else color = '#27ae60'; // Positive

    return (
        <div className="sentiment-gauge">
            <h3>Digital Momentum Score</h3>
            <div className="gauge-container">
                <ResponsiveContainer width="100%" height={200}>
                    <PieChart>
                        <Pie
                            data={chartData}
                            cx="50%"
                            cy="70%"
                            startAngle={180}
                            endAngle={0}
                            innerRadius={60}
                            outerRadius={80}
                            paddingAngle={5}
                            dataKey="value"
                            stroke="none"
                        >
                            <Cell key="cell-0" fill={color} />
                            <Cell key="cell-1" fill="#e2e8f0" />
                        </Pie>
                    </PieChart>
                </ResponsiveContainer>
                <div className="score-display">
                    <span className="score-number" style={{ color }}>{score}</span>
                    <span className="score-label">/ 100</span>
                </div>
            </div>
        </div>
    );
};

export default SentimentGauge;
