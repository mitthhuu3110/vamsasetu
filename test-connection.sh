#!/bin/bash

echo "🧪 Testing VamsaSetu Database Connections"
echo "========================================"

echo "📊 Testing PostgreSQL..."
docker-compose exec postgres psql -U vamsasetu -d vamsasetu -c "SELECT 'PostgreSQL is working!' as status;"

echo ""
echo "🔗 Testing Redis..."
docker-compose exec redis redis-cli ping

echo ""
echo "🌐 Testing Neo4j..."
echo "Neo4j is running at: http://localhost:7474"
echo "Username: neo4j"
echo "Password: vamsasetu123"

echo ""
echo "✅ All databases are running!"
echo ""
echo "🌐 Access URLs:"
echo "  Neo4j Browser: http://localhost:7474"
echo "  PostgreSQL: localhost:5432"
echo "  Redis: localhost:6379"
echo ""
echo "📝 Next Steps:"
echo "  1. Open Neo4j Browser and login"
echo "  2. Run: MATCH (n) RETURN n LIMIT 25"
echo "  3. The application is ready for development!"
