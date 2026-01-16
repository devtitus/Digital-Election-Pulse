import React from 'react';
import { BarChart, Bar, XAxis, YAxis, ResponsiveContainer, Tooltip } from 'recharts';
import './ComparisonChart.css';

const ComparisonChart = ({ data, selectedParty }) => {
    const CustomTooltip = ({ active, payload, label }) => {
        if (active && payload && payload.length) {
            const party = data.find(p => p.party === label);
            return (
                <div className="comparison-tooltip">
                    <p className="tooltip-party">{label}</p>
                    <p className="tooltip-value" style={{ color: party?.color }}>
                        Sentiment: {payload[0].value}/100
                    </p>
                </div>
            );
        }
        return null;
    };

    const getBarColor = (party) => {
        if (selectedParty && party === selectedParty.name) {
            // Highlight selected party with brighter color
            return data.find(p => p.party === party)?.color || 'var(--accent-gold)';
        }
        return data.find(p => p.party === party)?.color || 'var(--text-muted)';
    };

    if (!data || data.length === 0) {
        return (
            <div className="comparison-empty">
                <span>⚖️</span>
                <p>No comparison data available</p>
            </div>
        );
    }

    return (
        <div className="comparison-chart">
            <ResponsiveContainer width="100%" height={120}>
                <BarChart data={data} margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                    <XAxis
                        dataKey="party"
                        axisLine={false}
                        tickLine={false}
                        tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                        interval={0}
                    />
                    <YAxis
                        domain={[0, 100]}
                        axisLine={false}
                        tickLine={false}
                        tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                        ticks={[0, 25, 50, 75, 100]}
                    />
                    <Tooltip content={<CustomTooltip />} />
                    <Bar
                        dataKey="sentiment"
                        radius={[2, 2, 0, 0]}
                        fill={(entry) => getBarColor(entry.party)}
                    >
                        {data.map((entry, index) => (
                            <Bar key={`bar-${index}`} fill={entry.color} />
                        ))}
                    </Bar>
                </BarChart>
            </ResponsiveContainer>
            <div className="comparison-footer">
                <span className="comparison-label">Party Sentiment Comparison</span>
            </div>
        </div>
    );
};

export default ComparisonChart;