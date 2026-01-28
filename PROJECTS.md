# 🗳️ The TN Election Pulse - Project Overview

## 📋 Project Summary

**The TN Election Pulse** is a real-time digital momentum analytics dashboard designed to analyze and predict the sentiment and popularity of political parties in Tamil Nadu, India. The system aggregates data from multiple sources including news articles, social media comments, and discussion forums, then uses advanced AI to provide comprehensive sentiment analysis.

## 🎯 Key Objectives

* **Real-time Sentiment Analysis**: Provide up-to-date sentiment scores for political parties
* **Multi-source Data Aggregation**: Combine data from Google News, YouTube, and Reddit
* **AI-Powered Insights**: Use Google Gemini AI for advanced sentiment and emotion analysis
* **Visual Data Representation**: Present complex data in intuitive, interactive visualizations
* **Historical Trend Tracking**: Allow users to track sentiment changes over time

## 🏗️ System Architecture

The project follows a modern client-server architecture with:

### Backend Components
- **Go (Golang) API Server**: Built with Fiber framework for high performance
- **PostgreSQL Database**: For storing party information and sentiment snapshots
- **AI Integration**: Google Gemini SDK for sentiment analysis
- **Data Fetchers**: Concurrent services for Google News, YouTube, and Reddit data

### Frontend Components
- **React Application**: Built with Vite for fast development
- **Responsive UI**: Modern dark theme with glassmorphism effects
- **Data Visualization**: Interactive charts and gauges using Recharts
- **Real-time Updates**: Live data refresh capabilities

## 🔧 Technical Stack

### Backend
- **Language**: Go 1.23+
- **Framework**: Fiber (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **AI**: Google Gemini 1.5 Flash
- **APIs**: YouTube Data API v3, Google News RSS, Reddit JSON API

### Frontend
- **Framework**: React 18+
- **Build Tool**: Vite
- **Styling**: Vanilla CSS with CSS Variables
- **Visualization**: Recharts
- **State Management**: React Hooks

### DevOps
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Environment Management**: .env files

## 📊 Data Flow

1. **User Interaction**: User selects a political party from the dashboard
2. **API Request**: Frontend sends analysis request to backend
3. **Data Collection**: Backend concurrently fetches data from multiple sources
4. **AI Analysis**: Aggregated data is sent to Google Gemini for sentiment analysis
5. **Result Processing**: AI results are processed and stored in database
6. **Response**: Processed data is sent back to frontend
7. **Visualization**: Frontend displays results in interactive charts and gauges

## 🎨 UI/UX Features

* **Party Selector**: Dropdown to choose between different political parties
* **Sentiment Wave**: Visual representation of sentiment score
* **Emotion Display**: Shows dominant emotion detected
* **Key Topics**: Lists main discussion topics
* **Trend Graph**: Historical sentiment trends
* **Comparison Chart**: Compare sentiment across parties
* **Responsive Design**: Works on desktop and mobile devices

## 🔍 AI Analysis Capabilities

The system uses Google Gemini AI to perform:

* **Sentiment Scoring**: -100 (extremely negative) to +100 (extremely positive)
* **Emotion Detection**: Identifies emotions like Hope, Anger, Fear, Mockery
* **Topic Extraction**: Extracts key discussion topics
* **Fact Checking**: Flags potential misinformation
* **Bias Mitigation**: Ensures objective analysis

## 📈 Sentiment Calculation

The sentiment score is calculated using a weighted formula:
- Raw AI score (-1.0 to +1.0) is converted to 0-100 scale
- Formula: `FinalScore = 50 + (RawScore * 50)`
- This provides a balanced scale centered around 50 (neutral)

## 🗂️ Database Schema

### Party Table
- `id`: Primary key
- `name`: Party name
- `leader`: Party leader name
- `color_hex`: Party color for UI representation
- `created_at`, `updated_at`: Timestamps

### SentimentSnapshot Table
- `id`: Primary key
- `party_id`: Foreign key to Party
- `score`: Sentiment score (0-100)
- `key_topics`: JSON array of discussion topics
- `emotion`: Dominant emotion
- `source_breakdown`: Data source distribution
- `created_at`: Timestamp

## 🌐 API Endpoints

* `GET /api/v1/parties`: Get list of political parties
* `POST /api/v1/analyze`: Trigger sentiment analysis for a party
* `GET /api/v1/latest`: Get latest cached analysis
* `GET /api/v1/trends`: Get historical trends
* `GET /api/v1/comparison`: Compare sentiment across parties

## 🚀 Deployment Options

### Local Development
- Run backend and frontend separately
- Use local PostgreSQL instance
- Ideal for development and testing

### Docker Deployment
- Containerized backend and frontend
- Single command deployment with Docker Compose
- Production-ready setup
- Easy scaling and management

### Cloud Deployment
- Can be deployed to any cloud provider
- Supports environment variables for configuration
- Stateless design for easy scaling

## 📁 Project Structure

```
/
├── backend/              # Go backend code
│   ├── cmd/              # Main application
│   ├── db/               # Database connection
│   ├── handlers/         # API routes and handlers
│   ├── models/           # Data models
│   ├── services/         # Business logic and services
│   └── Dockerfile        # Backend container
├── frontend/             # React frontend code
│   ├── src/              # Source code
│   │   ├── api/          # API clients
│   │   ├── components/   # React components
│   │   ├── styles/       # CSS styles
│   │   └── App.jsx       # Main application
│   └── Dockerfile        # Frontend container
├── docs/                 # Documentation
├── docker-compose.yml    # Docker orchestration
└── README.md             # Project documentation
```

## 🎯 Future Enhancements

* **Additional Data Sources**: Twitter/X, Facebook, regional news
* **Advanced Analytics**: Predictive modeling and trend forecasting
* **User Authentication**: Role-based access control
* **Export Features**: PDF/CSV report generation
* **Mobile App**: Native mobile applications
* **Multi-language Support**: Tamil and other regional languages

## 🤝 Contribution Guidelines

Contributions are welcome! Please follow these guidelines:

1. **Fork the repository** and create your branch
2. **Follow coding standards** consistent with the project
3. **Write tests** for new features
4. **Update documentation** for changes
5. **Submit pull requests** with clear descriptions

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For questions, issues, or support, please:
- Open an issue on GitHub
- Check the documentation in the docs/ folder
- Review the API reference for technical details