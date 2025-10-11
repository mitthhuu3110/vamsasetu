# VamsaSetu - Family Tree & Event Management System

VamsaSetu is an intelligent, tree-based family relationship visualization and event system that helps users understand how they're related to others in their family, track events (birthdays, anniversaries, ceremonies), and receive smart reminders through WhatsApp, email, and SMS.

## 🌟 Features

### Core Features
- **Interactive Family Tree Visualizer**: Dynamic tree graph showing members as nodes and relations as edges
- **Smart Relationship Engine**: Auto-derive relationships based on input with complex Indian relations support
- **Event Reminder System**: Automated notifications via WhatsApp, Email, and SMS
- **Event Highlighting on Tree**: Highlight paths connecting members involved in events
- **AI-Generated Relationship Descriptions**: Natural language family relationship summaries
- **User Authentication & Profiles**: OAuth support with role-based access
- **Search & Filter**: Search by name, relation, or event with advanced filtering
- **Multi-Family Support**: Users can belong to multiple family trees
- **Backup & Versioning**: Auto backup and export capabilities

### Technical Features
- **Real-time Updates**: WebSocket support for live family tree updates
- **Responsive Design**: Mobile-first approach with modern UI/UX
- **Graph Database**: Neo4j for complex relationship modeling
- **Caching**: Redis for improved performance
- **Security**: JWT-based authentication with role-based authorization
- **API-First**: RESTful APIs with comprehensive documentation

## 🏗️ Architecture

### Tech Stack
- **Frontend**: React + TypeScript + Tailwind CSS + React Flow
- **Backend**: Spring Boot + Java 17
- **Databases**: Neo4j (Graph) + PostgreSQL (Relational)
- **Cache**: Redis
- **Authentication**: JWT + Spring Security
- **Messaging**: Twilio (SMS/WhatsApp) + SendGrid (Email)
- **Deployment**: Docker + Docker Compose

### System Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │  Spring Boot    │    │   Databases     │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│ Neo4j + Postgres│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │  Notification   │    │     Redis       │
│   (Real-time)   │    │   Services      │    │    (Cache)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Java 17+ (for local development)
- Maven 3.6+ (for local development)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd vamsasetu
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080/api
   - Neo4j Browser: http://localhost:7474
   - PostgreSQL: localhost:5432

### Local Development

1. **Start databases**
   ```bash
   docker-compose up -d neo4j postgres redis
   ```

2. **Backend setup**
   ```bash
   cd backend
   ./mvnw clean install
   ./mvnw spring-boot:run
   ```

3. **Frontend setup**
   ```bash
   cd frontend
   npm install
   npm start
   ```

## 📁 Project Structure

```
vamsasetu/
├── frontend/                 # React TypeScript frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── types/          # TypeScript type definitions
│   │   └── contexts/       # React contexts
│   ├── public/             # Static assets
│   └── package.json
├── backend/                 # Spring Boot backend
│   ├── src/main/java/
│   │   ├── controller/     # REST controllers
│   │   ├── service/        # Business logic
│   │   ├── repository/     # Data access layer
│   │   ├── model/          # Entity models
│   │   ├── security/       # Security configuration
│   │   └── config/         # Application configuration
│   └── pom.xml
├── docker-compose.yml      # Multi-container setup
└── README.md
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# Database
POSTGRES_URL=jdbc:postgresql://localhost:5432/vamsasetu
POSTGRES_USERNAME=vamsasetu
POSTGRES_PASSWORD=vamsasetu123

NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=vamsasetu123

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key-here

# Email
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password

# Twilio
TWILIO_ACCOUNT_SID=your-account-sid
TWILIO_AUTH_TOKEN=your-auth-token
TWILIO_PHONE_NUMBER=your-phone-number

# WhatsApp
WHATSAPP_API_URL=your-whatsapp-api-url
WHATSAPP_ACCESS_TOKEN=your-access-token
```

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET /api/auth/profile` - Get user profile

### Family Tree Endpoints
- `GET /api/family/members` - Get family members
- `POST /api/family/members` - Add family member
- `GET /api/family/relationships` - Get relationships
- `POST /api/family/relationships` - Add relationship
- `GET /api/family/tree` - Get family tree data

### Event Endpoints
- `GET /api/events` - Get events
- `POST /api/events` - Create event
- `PUT /api/events/{id}` - Update event
- `DELETE /api/events/{id}` - Delete event

## 🧪 Testing

### Backend Tests
```bash
cd backend
./mvnw test
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Integration Tests
```bash
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 🚀 Deployment

### Production Deployment

1. **Build production images**
   ```bash
   docker-compose -f docker-compose.prod.yml build
   ```

2. **Deploy to production**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Cloud Deployment

The application is designed to be cloud-native and can be deployed on:
- AWS (ECS, EKS, RDS, ElastiCache)
- Google Cloud (GKE, Cloud SQL, Memorystore)
- Azure (AKS, Database, Cache)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support, email support@vamsasetu.com or join our Slack channel.

## 🗺️ Roadmap

- [ ] Voice query support ("Hey, how is Sathvika related to Charan?")
- [ ] Legacy tree mode (ancestors only)
- [ ] Smart relationship suggestions
- [ ] Data privacy mode with encryption
- [ ] Mobile app (React Native)
- [ ] Advanced analytics and insights
- [ ] Multi-language support
- [ ] Integration with genealogy services

## 🙏 Acknowledgments

- React Flow for tree visualization
- Spring Boot for robust backend
- Neo4j for graph database capabilities
- Tailwind CSS for beautiful UI components