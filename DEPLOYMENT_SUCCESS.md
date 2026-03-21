# VamsaSetu - Deployment Successful! 🎉

## Application is Running

All services have been successfully built and deployed. Your VamsaSetu application is now accessible!

## Access Points

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Backend Health**: http://localhost:8080/health
- **Neo4j Browser**: http://localhost:7474 (username: `neo4j`, password: `vamsasetu123`)

## Services Status ✅

All services are healthy and running:
- ✅ PostgreSQL (port 5432)
- ✅ Neo4j Graph Database (ports 7474, 7687)
- ✅ Redis Cache (port 6379)
- ✅ Backend API (port 8080)
- ✅ Frontend (port 5173)

## Demo User Created ✅

A demo user has been created for you:
- **Email**: demo@vamsasetu.com
- **Password**: Demo@1234

## Getting Started

### 1. Login to the Application

Open http://localhost:5173 in your browser and login with:
- Email: demo@vamsasetu.com
- Password: Demo@1234

### 2. Start Building Your Family Tree

After login, you can:
- Add family members
- Create relationships between members
- Add events (birthdays, anniversaries, ceremonies)
- Visualize your family tree

### 3. Create Additional Users (Optional)

You can register more users via the frontend or API:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "YourPassword123",
    "name": "Your Name",
    "role": "owner"
  }'
```

## Managing the Application

### Stop Services
```bash
docker-compose down
```

### Restart Services
```bash
docker-compose up -d
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs backend
docker-compose logs frontend
```

### Check Status
```bash
docker-compose ps
```

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh token
- `GET /api/auth/profile` - Get profile (requires auth)

### Members
- `GET /api/members` - List all members
- `POST /api/members` - Create member
- `GET /api/members/:id` - Get member
- `PUT /api/members/:id` - Update member
- `DELETE /api/members/:id` - Delete member
- `GET /api/members/search?q=name` - Search members

### Relationships
- `GET /api/relationships` - List relationships
- `POST /api/relationships` - Create relationship
- `DELETE /api/relationships/:id` - Delete relationship
- `GET /api/relationships/path?from=id1&to=id2` - Find path

### Events
- `GET /api/events` - List events
- `POST /api/events` - Create event
- `GET /api/events/:id` - Get event
- `PUT /api/events/:id` - Update event
- `DELETE /api/events/:id` - Delete event
- `GET /api/events/upcoming` - Upcoming events

### Family Tree
- `GET /api/family/tree` - Get complete family tree (React Flow format)

## Technology Stack

### Frontend
- React + TypeScript
- Vite
- TailwindCSS
- React Query
- React Flow (tree visualization)
- Zustand (state management)

### Backend
- Go + Fiber framework
- PostgreSQL (user data, events)
- Neo4j (family relationships graph)
- Redis (caching)
- JWT authentication
- WebSocket support

## Troubleshooting

### Check Service Health
```bash
curl http://localhost:8080/health
```

### View Service Logs
```bash
docker-compose logs [service-name]
```

### Restart a Service
```bash
docker-compose restart [service-name]
```

### Remove All Data and Start Fresh
```bash
docker-compose down -v
docker-compose up -d
```

## Note About Seed Script

The seed script (`backend/cmd/seed/main.go`) is designed to run from your local machine but requires access to the Docker network. Since the databases are running in containers, you have two options:

1. **Use the API to create data** (recommended) - Use curl or the frontend to create users, members, and relationships
2. **Modify the .env file** - Change database hosts from `localhost` to the container names when running the seed script locally

For now, you can use the demo user and create family members through the frontend interface.

## Next Steps

1. ✅ Login with demo credentials (demo@vamsasetu.com / Demo@1234)
2. Add family members through the UI
3. Create relationships between members
4. Add important events
5. Explore the interactive family tree visualization
6. Test mobile responsiveness
7. Invite family members to collaborate

Enjoy building your family tree with VamsaSetu! 🌳👨‍👩‍👧‍👦
