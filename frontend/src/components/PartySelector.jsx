import React from 'react';
import './PartySelector.css';

const PartySelector = ({ parties, selectedPartyId, onSelect }) => {
    return (
        <div className="party-selector-container">
            <div className="party-selector-scroll">
                {parties.map((party) => (
                    <button
                        key={party.id}
                        className={`party-btn ${selectedPartyId === party.id ? 'active' : ''}`}
                        style={{ '--party-color': party.color_hex }}
                        onClick={() => onSelect(party)}
                    >
                        {party.name}
                    </button>
                ))}
            </div>
        </div>
    );
};

export default PartySelector;
