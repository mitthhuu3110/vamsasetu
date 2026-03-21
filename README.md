# VamsaSetu - Family Tree & Event Management System

VamsaSetu is an intelligent, tree-based family relationship visualization and event system that helps users understand how they're related to others in their family, track events (birthdays, anniversaries, ceremonies), and receive smart reminders through WhatsApp, email, and SMS.
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
- **Frontend**: React + TypeScript + Vite + Tailwind CSS v4 + React Flow + Framer Motion
- **Backend**: Go + Fiber + GORM
- **Databases**: Neo4j (Graph) + PostgreSQL (Relational)
- **Cache**: Redis
- **Authentication**: JWT
- **Messaging**: Twilio (SMS/WhatsApp) + SendGrid (Email)
- **Deployment**: Docker + Docker Compose

### System Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React App     │    │   Go + Fiber    │    │   Databases     │
│   (Frontend)    │◄──►│   (Backend)     │◄──►│ Neo4j + Postgres│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   WebSocket     │    │  Notification   │    │     Redis       │
│   (Real-time)   │    │   Scheduler     │    │    (Cache)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd vamsasetu
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:5173
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
   go mod download
   go run cmd/server/main.go
   ```

3. **Frontend setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## 📁 Project Structure

```
vamsasetu/
├── frontend/                 # React TypeScript frontend
│   ├── src/
│   │   ├── components/      # Reusable UI components
│   │   │   ├── auth/       # Authentication components
│   │   │   ├── common/     # Common components (Layout, Navbar, etc.)
│   │   │   ├── events/     # Event-related components
│   │   │   ├── family/     # Family tree components
│   │   │   └── ui/         # Base UI components
│   │   ├── pages/          # Page components
│   │   ├── services/       # API services
│   │   ├── hooks/          # React Query hooks
│   │   ├── stores/         # Zustand state management
│   │   ├── types/          # TypeScript type definitions
│   │   └── utils/          # Utility functions
│   ├── public/             # Static assets
│   └── package.json
├── backend/                 # Go + Fiber backend
│   ├── cmd/server/         # Main application entry point
│   ├── internal/
│   │   ├── config/         # Configuration management
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # Middleware (auth, CORS, logging)
│   │   ├── models/         # Data models
│   │   ├── repository/     # Data access layer
│   │   ├── scheduler/      # Background jobs
│   │   ├── service/        # Business logic
│   │   └── utils/          # Utility functions
│   ├── pkg/                # Shared packages
│   │   ├── neo4j/         # Neo4j client
│   │   ├── postgres/      # PostgreSQL client
│   │   └── redis/         # Redis client
│   ├── go.mod
│   └── go.sum
├── .kiro/                  # Kiro spec files
│   └── specs/
│       └── vamsasetu-full-system/
│           ├── requirements.md
│           ├── design.md
│           └── tasks.md
├── docker-compose.yml      # Multi-container setup
├── .env.example           # Environment variables template
└── README.md
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory (use `.env.example` as template):

```env
# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=vamsasetu
POSTGRES_USER=vamsasetu
POSTGRES_PASSWORD=vamsasetu123

NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=vamsasetu123

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Server Configuration
PORT=8080
FRONTEND_URL=http://localhost:5173

# Notification Services
TWILIO_ACCOUNT_SID=your-account-sid
TWILIO_AUTH_TOKEN=your-auth-token
TWILIO_PHONE_NUMBER=+1234567890
TWILIO_WHATSAPP_NUMBER=whatsapp:+1234567890

SENDGRID_API_KEY=your-sendgrid-api-key
SENDGRID_FROM_EMAIL=noreply@vamsasetu.com
SENDGRID_FROM_NAME=VamsaSetu

# Application Settings
LOG_LEVEL=info
NOTIFICATION_DAYS_AHEAD=7
```

## 📚 API Documentation

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/auth/profile` - Get user profile

### Member Endpoints
- `GET /api/members` - Get all family members (with pagination)
- `GET /api/members/:id` - Get member by ID
- `POST /api/members` - Add family member (owner/admin only)
- `PUT /api/members/:id` - Update member (owner/admin only)
- `DELETE /api/members/:id` - Soft delete member (owner/admin only)
- `GET /api/members/search?q=name` - Search members by name

### Relationship Endpoints
- `GET /api/relationships` - Get all relationships
- `POST /api/relationships` - Add relationship (owner/admin only)
- `DELETE /api/relationships/:id` - Delete relationship (owner/admin only)
- `GET /api/relationships/path?from=uuid1&to=uuid2` - Find relationship path

### Family Tree Endpoints
- `GET /api/family/tree` - Get family tree data (React Flow format)

### Event Endpoints
- `GET /api/events` - Get events (with filters)
- `GET /api/events/:id` - Get event by ID
- `POST /api/events` - Create event (owner/admin only)
- `PUT /api/events/:id` - Update event (owner/admin only)
- `DELETE /api/events/:id` - Delete event (owner/admin only)
- `GET /api/events/upcoming` - Get upcoming events

### WebSocket Endpoint
- `WS /api/ws` - WebSocket connection for real-time updates

### Health Check
- `GET /health` - Application health status

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Build Verification
```bash
# Backend
cd backend
go build -o bin/server cmd/server/main.go

# Frontend
cd frontend
npm run build
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
- Go Fiber for high-performance backend
- Neo4j for graph database capabilities
- Tailwind CSS v4 for beautiful UI components
- Framer Motion for smooth animations
- React Query for efficient data fetching
- Zustand for lightweight state management