import React from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, ResponsiveContainer, Tooltip } from 'recharts';
import './TrendGraph.css';

const TrendGraph = ({ data, party }) => {
    // Transform data for recharts
    const chartData = data.map(item => ({
        date: new Date(item.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
        sentiment: item.sentiment,
        fullDate: item.date
    }));

    const getLineColor = () => {
        // Use party color if available, otherwise fallback to default
        return party?.color_hex || 'var(--accent-gold)';
    };

    const CustomTooltip = ({ active, payload, label }) => {
        if (active && payload && payload.length) {
            return (
                <div className="trend-tooltip">
                    <p className="tooltip-date">{label}</p>
                    <p className="tooltip-value" style={{ color: getLineColor() }}>
                        Sentiment: {payload[0].value}/100
                    </p>
                </div>
            );
        }
        return null;
    };

    if (!data || data.length === 0) {
        return (
            <div className="trend-empty">
                <span>📈</span>
                <p>No trend data available</p>
            </div>
        );
    }

    return (
        <div className="trend-graph">
            <ResponsiveContainer width="100%" height={120}>
                <LineChart data={chartData} margin={{ top: 5, right: 5, left: 5, bottom: 5 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="rgba(251, 191, 36, 0.1)" />
                    <XAxis
                        dataKey="date"
                        axisLine={false}
                        tickLine={false}
                        tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                    />
                    <YAxis
                        domain={[0, 100]}
                        axisLine={false}
                        tickLine={false}
                        tick={{ fontSize: 10, fill: 'var(--text-muted)' }}
                        ticks={[0, 25, 50, 75, 100]}
                    />
                    <Tooltip content={<CustomTooltip />} />
                    <Line
                        type="monotone"
                        dataKey="sentiment"
                        stroke={getLineColor()}
                        strokeWidth={2}
                        dot={{ fill: getLineColor(), strokeWidth: 2, r: 3 }}
                        activeDot={{ r: 5, stroke: getLineColor(), strokeWidth: 2, fill: 'var(--bg-charcoal)' }}
                    />
                </LineChart>
            </ResponsiveContainer>
            <div className="trend-footer">
                <span className="trend-indicator animate-pulse-subtle" style={{ backgroundColor: getLineColor() }}></span>
                <span className="trend-label">7-Day Trend</span>
            </div>
        </div>
    );
};

export default TrendGraph;