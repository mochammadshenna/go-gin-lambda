#!/bin/bash

# PostgreSQL Setup Script for AI Service
# This script helps set up PostgreSQL and run migrations

set -e

echo "🚀 Setting up PostgreSQL for AI Service..."

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed."
    echo ""
    echo "📋 Installation Options:"
    echo "1. Install Homebrew first, then PostgreSQL:"
    echo "   /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    echo "   brew install postgresql@14"
    echo "   brew services start postgresql@14"
    echo ""
    echo "2. Download from PostgreSQL official website:"
    echo "   https://www.postgresql.org/download/macosx/"
    echo ""
    echo "3. Use Docker (if you have Docker installed):"
    echo "   docker run --name postgres-ai -e POSTGRES_PASSWORD=ai_service_password -e POSTGRES_DB=ai_service -p 5432:5432 -d postgres:14"
    echo ""
    echo "After installing PostgreSQL, run this script again."
    exit 1
fi

# Database configuration
DB_NAME="ai_service"
DB_USER="ai_service_user"
DB_PASSWORD="ai_service_password"
DB_HOST="localhost"
DB_PORT="5432"

echo "✅ PostgreSQL is installed."

# Check if PostgreSQL service is running
if ! pg_isready -h $DB_HOST -p $DB_PORT &> /dev/null; then
    echo "⚠️  PostgreSQL service is not running."
    echo "Starting PostgreSQL service..."
    
    # Try to start PostgreSQL service
    if command -v brew &> /dev/null; then
        brew services start postgresql@14
    else
        echo "❌ Cannot start PostgreSQL automatically."
        echo "Please start PostgreSQL service manually and run this script again."
        exit 1
    fi
    
    # Wait for service to start
    echo "⏳ Waiting for PostgreSQL to start..."
    sleep 5
fi

echo "✅ PostgreSQL service is running."

# Create database and user
echo "📊 Creating database and user..."

# Connect as postgres superuser to create database and user
psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE DATABASE $DB_NAME;" 2>/dev/null || echo "Database already exists"
psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "User already exists"
psql -h $DB_HOST -p $DB_PORT -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" 2>/dev/null || echo "Privileges already granted"

echo "✅ Database and user created successfully."

# Run migrations
echo "🔄 Running database migrations..."

# Set environment variables for the application
export DB_TYPE=postgres
export DB_HOST=$DB_HOST
export DB_PORT=$DB_PORT
export DB_USER=$DB_USER
export DB_PASSWORD=$DB_PASSWORD
export DB_NAME=$DB_NAME
export DB_SSLMODE=disable

# Run the migration
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f scripts/migrations/001_initial_schema.sql

echo "✅ Migrations completed successfully."

# Test the connection
echo "🧪 Testing database connection..."

# Create a simple test script
cat > test_db.sql << EOF
SELECT 'Database connection successful!' as status;
SELECT COUNT(*) as providers_count FROM providers;
SELECT COUNT(*) as generations_count FROM generations;
EOF

psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f test_db.sql

# Clean up test file
rm test_db.sql

echo ""
echo "🎉 PostgreSQL setup completed successfully!"
echo ""
echo "📋 Database Configuration:"
echo "   Host: $DB_HOST"
echo "   Port: $DB_PORT"
echo "   Database: $DB_NAME"
echo "   User: $DB_USER"
echo "   Password: $DB_PASSWORD"
echo ""
echo "🔧 To run the application with PostgreSQL:"
echo "   export DB_TYPE=postgres"
echo "   export DB_HOST=$DB_HOST"
echo "   export DB_PORT=$DB_PORT"
echo "   export DB_USER=$DB_USER"
echo "   export DB_PASSWORD=$DB_PASSWORD"
echo "   export DB_NAME=$DB_NAME"
echo "   export DB_SSLMODE=disable"
echo "   go run cmd/main/main.go"
echo ""
echo "🌐 Or use the Makefile:"
echo "   make run-postgres" 