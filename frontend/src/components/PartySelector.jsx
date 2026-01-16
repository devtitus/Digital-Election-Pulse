import React from 'react';
import './PartySelector.css';

const PartySelector = ({ parties, selectedPartyId, onSelect }) => {
    // Default party icons and glows - can be overridden by API data
    const defaultPartyData = {
        'DMK': { icon: '🏴', glow: 'var(--dmk-red-glow)' },
        'AIADMK': { icon: '🐘', glow: 'var(--aiadmk-green-glow)' },
        'BJP': { icon: '🌺', glow: 'var(--bjp-orange-glow)' },
        'NTK': { icon: '⚡', glow: 'var(--ntk-red-dark-glow)' },
        'TVK': { icon: '🎭', glow: 'var(--tvk-amber-glow)' },
    };

    const getPartyIcon = (party) => {
        // Use icon from database if available, otherwise fallback to defaults
        return party.icon || defaultPartyData[party.name]?.icon || '🏛️';
    };

    const getPartyGlow = (party) => {
        // Use glow from database if available, otherwise fallback to defaults
        return party.glow || defaultPartyData[party.name]?.glow || 'var(--accent-gold-glow)';
    };

    return (
        <div className="party-hub">
            <h3 className="hub-title">Select Political Party</h3>
            <div className="party-nodes">
                {parties.map((party, index) => {
                    const isSelected = selectedPartyId === party.id;
                    return (
                        <button
                            key={party.id}
                            className={`party-node ${isSelected ? 'active' : ''}`}
                            style={{
                                '--party-color': party.color_hex,
                                '--party-glow': getPartyGlow(party),
                                '--animation-delay': `${index * 0.1}s`
                            }}
                            onClick={() => onSelect(party)}
                        >
                            <div className="node-icon">
                                {getPartyIcon(party)}
                            </div>
                            <div className="node-content">
                                <span className="node-name">{party.name}</span>
                                {party.leader && (
                                    <span className="node-leader">{party.leader}</span>
                                )}
                            </div>
                            <div className="node-pulse animate-pulse-subtle"></div>
                        </button>
                    );
                })}
            </div>
        </div>
    );
};

export default PartySelector;
