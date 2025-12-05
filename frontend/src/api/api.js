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

export const getHistory = async (partyId) => {
    try {
        const response = await api.get(`/history/${partyId}`);
        return response.data;
    } catch (error) {
        console.error("API Error fetching history:", error);
        return [];
    }
};
