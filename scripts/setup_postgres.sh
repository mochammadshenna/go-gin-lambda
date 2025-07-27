#!/bin/bash

# PostgreSQL Setup Script for AI Service
# This script helps set up PostgreSQL for the AI service

set -e

echo "🐘 Setting up PostgreSQL for AI Service..."

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install PostgreSQL first."
    echo "   Visit: https://www.postgresql.org/download/"
    exit 1
fi

# Check if PostgreSQL service is running
if ! pg_isready -q; then
    echo "❌ PostgreSQL service is not running. Please start PostgreSQL first."
    echo "   On macOS: brew services start postgresql"
    echo "   On Ubuntu: sudo systemctl start postgresql"
    echo "   On Windows: Start PostgreSQL service from Services"
    exit 1
fi

# Database configuration
DB_NAME="ai_service"
DB_USER="postgres"
DB_PASSWORD=""  # Empty password for local development

echo "📋 Database Configuration:"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo "   Password: (empty for local development)"

# Create database if it doesn't exist
echo "🗄️  Creating database..."
createdb -U $DB_USER $DB_NAME 2>/dev/null || echo "   Database already exists"

# Run migrations
echo "📝 Running migrations..."
if [ -f "scripts/migrations/001_initial_schema.sql" ]; then
    psql -U $DB_USER -d $DB_NAME -f scripts/migrations/001_initial_schema.sql
    echo "   ✅ Initial schema applied"
else
    echo "   ⚠️  Migration file not found: scripts/migrations/001_initial_schema.sql"
fi

if [ -f "scripts/migrations/002_add_performance_indexes.sql" ]; then
    psql -U $DB_USER -d $DB_NAME -f scripts/migrations/002_add_performance_indexes.sql
    echo "   ✅ Performance indexes applied"
else
    echo "   ⚠️  Migration file not found: scripts/migrations/002_add_performance_indexes.sql"
fi

# Create test database
echo "🧪 Creating test database..."
TEST_DB_NAME="ai_service_test"
createdb -U $DB_USER $TEST_DB_NAME 2>/dev/null || echo "   Test database already exists"

# Test the setup
echo "🧪 Testing database setup..."
psql -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) as providers_count FROM providers;" 2>/dev/null || echo "   ⚠️  Could not query providers table"

echo ""
echo "✅ PostgreSQL setup completed!"
echo ""
echo "📋 Next steps:"
echo "   1. Copy .env.example to .env"
echo "   2. Update database credentials in .env if needed"
echo "   3. Set your AI provider API keys in .env"
echo "   4. Run: go run cmd/main/main.go"
echo ""
echo "🔗 Useful commands:"
echo "   - Connect to database: psql -U $DB_USER -d $DB_NAME"
echo "   - View tables: \dt"
echo "   - View data: SELECT * FROM generations LIMIT 5;" 