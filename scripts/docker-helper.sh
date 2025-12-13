#!/bin/bash
# Docker Helper Script for LogiTrackPro
# Provides common Docker operations

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

show_help() {
    echo "LogiTrackPro Docker Helper"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  up          - Start all services"
    echo "  down        - Stop all services"
    echo "  restart     - Restart all services"
    echo "  logs        - Show logs (all services)"
    echo "  logs [svc]  - Show logs for specific service"
    echo "  build       - Rebuild all containers"
    echo "  clean       - Remove containers and volumes (⚠️  DESTRUCTIVE)"
    echo "  status      - Show service status"
    echo "  shell [svc] - Open shell in service container"
    echo "  db-shell    - Open PostgreSQL shell"
    echo ""
}

start_services() {
    echo -e "${BLUE}🚀 Starting all services...${NC}"
    docker-compose up -d
    echo -e "${GREEN}✅ Services started${NC}"
    echo ""
    echo "📍 Services available at:"
    echo "   Frontend:  http://localhost:3000"
    echo "   Backend:   http://localhost:8080"
    echo "   Optimizer: http://localhost:8000"
    echo "   Database:  localhost:5432"
}

stop_services() {
    echo -e "${BLUE}🛑 Stopping all services...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ Services stopped${NC}"
}

restart_services() {
    echo -e "${BLUE}🔄 Restarting all services...${NC}"
    docker-compose restart
    echo -e "${GREEN}✅ Services restarted${NC}"
}

show_logs() {
    if [ -z "$1" ]; then
        docker-compose logs -f
    else
        docker-compose logs -f "$1"
    fi
}

rebuild_services() {
    echo -e "${BLUE}🔨 Rebuilding all containers...${NC}"
    docker-compose build --no-cache
    echo -e "${GREEN}✅ Containers rebuilt${NC}"
}

clean_all() {
    echo -e "${YELLOW}⚠️  This will remove all containers and volumes!${NC}"
    read -p "Are you sure? (yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        docker-compose down -v --remove-orphans
        echo -e "${GREEN}✅ Cleaned up${NC}"
    else
        echo "Cancelled"
    fi
}

show_status() {
    echo -e "${BLUE}📊 Service Status${NC}"
    docker-compose ps
    echo ""
    echo -e "${BLUE}💾 Volumes${NC}"
    docker volume ls | grep logitrackpro || echo "No volumes found"
}

open_shell() {
    if [ -z "$1" ]; then
        echo "❌ Please specify service: $0 shell [backend|optimizer|frontend|postgres]"
        exit 1
    fi
    
    case "$1" in
        backend)
            docker-compose exec backend sh
            ;;
        optimizer)
            docker-compose exec optimizer bash
            ;;
        frontend)
            docker-compose exec frontend sh
            ;;
        postgres|db)
            docker-compose exec postgres psql -U postgres -d logitrackpro
            ;;
        *)
            echo "❌ Unknown service: $1"
            exit 1
            ;;
    esac
}

# Main command handler
case "$1" in
    up)
        start_services
        ;;
    down)
        stop_services
        ;;
    restart)
        restart_services
        ;;
    logs)
        show_logs "$2"
        ;;
    build)
        rebuild_services
        ;;
    clean)
        clean_all
        ;;
    status)
        show_status
        ;;
    shell)
        open_shell "$2"
        ;;
    db-shell)
        docker-compose exec postgres psql -U postgres -d logitrackpro
        ;;
    *)
        show_help
        exit 1
        ;;
esac

