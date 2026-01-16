import axios from 'axios';

const API_BASE_URL = 'http://localhost:3000/api/v1';

export const api = axios.create({
    baseURL: API_BASE_URL,
});

export const getParties = async () => {
    try {
        const response = await api.get('/parties');
        if (Array.isArray(response.data)) {
            return response.data;
        } else {
            console.warn("API returned non-array for parties:", response.data);
            throw new Error("Invalid API response format");
        }
    } catch (error) {
        console.error("API Error fetching parties:", error);
        // Return mock data if API fails for dev/demo purposes until backend is fully up
        return [
            { id: 1, name: "DMK", color_hex: "#dd2e44", leader: "M.K. Stalin" },
            { id: 2, name: "AIADMK", color_hex: "#27ae60", leader: "Edappadi Palaniswami" },
            { id: 3, name: "TVK", color_hex: "#f1c40f", leader: "Vijay" },
        ];
    }
};

export const analyzeParty = async (partyName) => {
    try {
        const response = await api.post('/analyze', { party_name: partyName });
        return response.data;
    } catch (error) {
        console.error("API Error analyzing party:", error);
        throw error;
    }
};

export const getLatestSnapshot = async (partyName) => {
    try {
        const response = await api.get('/latest', { params: { party_name: partyName } });
        return response.data;
    } catch (error) {
        console.error("API Error fetching latest snapshot:", error);
        throw error;
    }
};

export const getTrends = async (partyName) => {
    try {
        const response = await api.get('/trends', { params: { party_name: partyName } });
        return response.data;
    } catch (error) {
        console.error("API Error fetching trends:", error);
        // Return mock trend data
        return [
            { date: '2024-01-01', sentiment: 45 },
            { date: '2024-01-02', sentiment: 52 },
            { date: '2024-01-03', sentiment: 48 },
            { date: '2024-01-04', sentiment: 61 },
            { date: '2024-01-05', sentiment: 55 },
            { date: '2024-01-06', sentiment: 58 },
            { date: '2024-01-07', sentiment: 62 }
        ];
    }
};

export const getComparison = async () => {
    try {
        const response = await api.get('/comparison');
        return response.data;
    } catch (error) {
        console.error("API Error fetching comparison:", error);
        // Return mock comparison data
        return [
            { party: 'DMK', sentiment: 62, color: '#dc2626' },
            { party: 'AIADMK', sentiment: 58, color: '#16a34a' },
            { party: 'BJP', sentiment: 45, color: '#ea580c' },
            { party: 'NTK', sentiment: 38, color: '#b91c1c' },
            { party: 'TVK', sentiment: 52, color: '#d97706' }
        ];
    }
};
