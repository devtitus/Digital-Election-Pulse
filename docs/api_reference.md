# 📡 API Reference

Base URL: `http://localhost:8080/api/v1`

## Endpoints

### 1. Get All Parties
Retrieves the list of configured political parties.

*   **URL**: `/parties`
*   **Method**: `GET`
*   **Response**: `200 OK`
    ```json
    [
      {
        "id": 1,
        "name": "DMK",
        "code": "DMK",
        "color_hex": "#ff3333"
      },
      ...
    ]
    ```

### 2. Analyze Party Sentiment
Triggers a real-time analysis for a specific party.

*   **URL**: `/analyze`
*   **Method**: `POST`
*   **Body**:
    ```json
    {
      "party_id": 1
    }
    ```
*   **Response**: `200 OK`
    ```json
    {
      "sentiment_score": 75.5,
      "emotion": "Hope",
      "key_topics": ["Flood Relief", "Metro Project"],
      "fact_check_notes": "None",
      "created_at": "2023-10-27T10:00:00Z"
    }
    ```

### 3. Get Latest Snapshot
Fetches the most recent cached analysis for a party without triggering a new AI run.

*   **URL**: `/latest`
*   **Method**: `GET`
*   **Query Params**: `party_id` (int)
*   **Response**: `200 OK`
    ```json
    {
      "exists": true,
      "sentiment_score": 75.5,
      "emotion": "Hope",
      "key_topics": ["Flood Relief", "Metro Project"],
      "created_at": "2023-10-27T10:00:00Z"
    }
    ```
    *(Returns `exists: false` if no prior data found)*
