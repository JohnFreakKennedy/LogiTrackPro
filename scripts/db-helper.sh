#!/bin/bash
# Database Helper Script
# Provides common database operations for LogiTrackPro

set -e

DB_NAME="logitrackpro"
DB_USER="postgres"
DB_HOST="localhost"
DB_PORT="5432"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

show_help() {
    echo "LogiTrackPro Database Helper"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  reset      - Drop and recreate database (⚠️  DESTRUCTIVE)"
    echo "  backup     - Create a backup of the database"
    echo "  restore    - Restore database from backup"
    echo "  migrate    - Run database migrations"
    echo "  seed       - Seed database with sample data"
    echo "  status     - Show database connection status"
    echo "  tables     - List all tables"
    echo "  stats      - Show database statistics"
    echo ""
}

reset_db() {
    echo -e "${BLUE}🔄 Resetting database...${NC}"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
    echo -e "${GREEN}✅ Database reset complete${NC}"
}

backup_db() {
    BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
    echo -e "${BLUE}💾 Creating backup: $BACKUP_FILE${NC}"
    PGPASSWORD=postgres pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME > $BACKUP_FILE
    echo -e "${GREEN}✅ Backup created: $BACKUP_FILE${NC}"
}

restore_db() {
    if [ -z "$1" ]; then
        echo "❌ Please provide backup file: $0 restore <backup_file.sql>"
        exit 1
    fi
    echo -e "${BLUE}📥 Restoring from: $1${NC}"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME < $1
    echo -e "${GREEN}✅ Database restored${NC}"
}

status_db() {
    echo -e "${BLUE}📊 Database Status${NC}"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            version() as postgres_version,
            current_database() as database_name,
            pg_size_pretty(pg_database_size(current_database())) as database_size;
    "
}

list_tables() {
    echo -e "${BLUE}📋 Database Tables${NC}"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        ORDER BY table_name;
    "
}

show_stats() {
    echo -e "${BLUE}📈 Database Statistics${NC}"
    PGPASSWORD=postgres psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        SELECT 
            schemaname,
            tablename,
            pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
            n_live_tup as row_count
        FROM pg_tables
        LEFT JOIN pg_stat_user_tables USING (tablename)
        WHERE schemaname = 'public'
        ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
    "
}

# Main command handler
case "$1" in
    reset)
        reset_db
        ;;
    backup)
        backup_db
        ;;
    restore)
        restore_db "$2"
        ;;
    status)
        status_db
        ;;
    tables)
        list_tables
        ;;
    stats)
        show_stats
        ;;
    *)
        show_help
        exit 1
        ;;
esac

