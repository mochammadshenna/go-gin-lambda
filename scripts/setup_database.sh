#!/bin/bash

# Simple Database Setup Script for AI Service
# This script helps set up the PostgreSQL database

set -e

echo "🐘 Setting up PostgreSQL database for AI Service..."

# Database configuration
DB_NAME="ai_service"
DB_USER="postgres"

echo "📋 Database Configuration:"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo "   Password: (empty for local development)"

# Check if PostgreSQL is running
if ! pg_isready -q; then
    echo "❌ PostgreSQL is not running. Please start PostgreSQL first."
    echo "   On macOS: brew services start postgresql@14"
    echo "   On Ubuntu: sudo systemctl start postgresql"
    echo "   On Windows: Start PostgreSQL service from Services"
    exit 1
fi

echo "✅ PostgreSQL is running"

# Create database if it doesn't exist
echo "🗄️  Creating database..."
createdb -U $DB_USER $DB_NAME 2>/dev/null || echo "   Database already exists"

# Run migrations
echo "📝 Running migrations..."
if [ -f "scripts/migrations/001_initial_schema.sql" ]; then
    psql -U $DB_USER -d $DB_NAME -f scripts/migrations/001_initial_schema.sql
    echo "   ✅ Initial schema applied"
else
    echo "   ❌ Migration file not found: scripts/migrations/001_initial_schema.sql"
    exit 1
fi

if [ -f "scripts/migrations/002_add_performance_indexes.sql" ]; then
    psql -U $DB_USER -d $DB_NAME -f scripts/migrations/002_add_performance_indexes.sql
    echo "   ✅ Performance indexes applied"
else
    echo "   ⚠️  Migration file not found: scripts/migrations/002_add_performance_indexes.sql"
fi

# Test the setup
echo "🧪 Testing database setup..."
psql -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) as providers_count FROM providers;" 2>/dev/null || echo "   ⚠️  Could not query providers table"
psql -U $DB_USER -d $DB_NAME -c "SELECT COUNT(*) as generations_count FROM generations;" 2>/dev/null || echo "   ⚠️  Could not query generations table"

echo ""
echo "✅ Database setup completed!"
echo ""
echo "📋 Next steps:"
echo "   1. Set your AI provider API keys in environment variables"
echo "   2. Run: make run"
echo ""
echo "🔗 Useful commands:"
echo "   - Connect to database: psql -U $DB_USER -d $DB_NAME"
echo "   - View tables: \dt"
echo "   - View data: SELECT * FROM generations LIMIT 5;" 