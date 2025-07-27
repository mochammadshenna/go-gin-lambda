# PostgreSQL Setup Guide

## 🐘 Quick Setup for AI Service

### Step 1: Install PostgreSQL

#### Option A: Using Homebrew (macOS)
```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Add Homebrew to PATH (if needed)
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc
source ~/.zshrc

# Install PostgreSQL
brew install postgresql@14

# Start PostgreSQL service
brew services start postgresql@14
```

#### Option B: Download from Official Website
1. Visit: https://www.postgresql.org/download/macosx/
2. Download and install PostgreSQL

### Step 2: Create Database and Run Migrations

```bash
# 1. Create the database
createdb ai_service

# 2. Run the initial schema migration
psql -d ai_service -f scripts/migrations/001_initial_schema.sql

# 3. Run performance indexes (if available)
psql -d ai_service -f scripts/migrations/002_add_performance_indexes.sql
```

### Step 3: Verify Setup

```bash
# Test database connection
psql -d ai_service -c "SELECT 'PostgreSQL is working!' as status;"

# Check if tables were created
psql -d ai_service -c "\dt"

# Check if providers were inserted
psql -d ai_service -c "SELECT name, is_active FROM providers;"
```

### Step 4: Run the Application

```bash
# Set environment variables (optional, defaults are already set)
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=""  # Empty password for local development
export DB_NAME=ai_service

# Run the application
go run cmd/main/main.go
```

## 🔧 Troubleshooting

### PostgreSQL Service Not Running
```bash
# Check if PostgreSQL is running
pg_isready

# Start PostgreSQL service
brew services start postgresql@14  # macOS
sudo systemctl start postgresql    # Linux
```

### Database Connection Issues
```bash
# Test connection
psql -h localhost -p 5432 -U postgres -d ai_service

# If you get permission errors, you might need to:
# 1. Create a user
createuser -s postgres

# 2. Or connect as the default user
psql -d ai_service
```

### Tables Don't Exist
```bash
# Check if migration was run
psql -d ai_service -c "\dt"

# If no tables, run migration again
psql -d ai_service -f scripts/migrations/001_initial_schema.sql
```

## 📋 Default Configuration

The application uses these default database settings:
- **Host**: localhost
- **Port**: 5432
- **User**: postgres
- **Password**: (empty for local development)
- **Database**: ai_service
- **SSL Mode**: disable

You can override these by setting environment variables:
```bash
export DB_HOST=your_host
export DB_PORT=your_port
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=your_database
```

## 🎯 Success Indicators

When everything is working correctly, you should see:
1. ✅ PostgreSQL service running
2. ✅ Database `ai_service` exists
3. ✅ Tables: `generations`, `providers`, `stats`, `api_keys`
4. ✅ Default providers inserted
5. ✅ Application starts without database errors
6. ✅ Web interface accessible at http://localhost:8080 