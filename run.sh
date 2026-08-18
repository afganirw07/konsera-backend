#!/bin/bash

# ============================================================
# SAKU / KONSERA BACKEND RUNNER
# ============================================================

set -o pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}      SAKU Backend Runner${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# ------------------------------------------------------------
# Load .env
# ------------------------------------------------------------
load_env() {
    if [ -f ".env" ]; then
        echo -e "${BLUE}Loading environment from .env...${NC}"

        set -a
        source .env
        set +a

        echo -e "${GREEN}✓ .env loaded${NC}"
    else
        echo -e "${YELLOW}! .env file not found${NC}"
        echo -e "Using environment variables from shell."
    fi

    echo ""
}

# ------------------------------------------------------------
# Check required environment variables
# ------------------------------------------------------------
check_env_vars() {
    echo -e "${BLUE}Checking environment variables...${NC}"

    local missing=false

    # PostgreSQL
    if [ -z "$DB_HOST" ]; then
        echo -e "${RED}✗ Missing DB_HOST${NC}"
        missing=true
    fi

    if [ -z "$DB_PORT" ]; then
        echo -e "${RED}✗ Missing DB_PORT${NC}"
        missing=true
    fi

    if [ -z "$DB_USER" ]; then
        echo -e "${RED}✗ Missing DB_USER${NC}"
        missing=true
    fi

    if [ -z "$DB_PASSWORD" ]; then
        echo -e "${RED}✗ Missing DB_PASSWORD${NC}"
        missing=true
    fi

    if [ -z "$DB_NAME" ]; then
        echo -e "${RED}✗ Missing DB_NAME${NC}"
        missing=true
    fi

    if [ -z "$DB_SSLMODE" ]; then
        echo -e "${RED}✗ Missing DB_SSLMODE${NC}"
        missing=true
    fi

    # SMTP
    if [ -z "$SMTP_HOST" ]; then
        echo -e "${RED}✗ Missing SMTP_HOST${NC}"
        missing=true
    fi

    if [ -z "$SMTP_PORT" ]; then
        echo -e "${RED}✗ Missing SMTP_PORT${NC}"
        missing=true
    fi

    if [ -z "$SMTP_EMAIL" ]; then
        echo -e "${RED}✗ Missing SMTP_EMAIL${NC}"
        missing=true
    fi

    if [ -z "$SMTP_PASSWORD" ]; then
        echo -e "${RED}✗ Missing SMTP_PASSWORD${NC}"
        missing=true
    fi

    if [ "$missing" = true ]; then
        echo ""
        echo -e "${RED}Some required environment variables are missing.${NC}"
        return 1
    fi

    echo -e "${GREEN}✓ All required environment variables are set${NC}"
    echo ""

    return 0
}

# ------------------------------------------------------------
# Check PostgreSQL
# ------------------------------------------------------------
check_postgres() {
    echo -e "${BLUE}Checking PostgreSQL connection...${NC}"

    if ! command -v psql &> /dev/null; then
        echo -e "${YELLOW}! psql not installed${NC}"
        echo -e "${YELLOW}Skipping PostgreSQL connection check.${NC}"
        echo ""

        return 0
    fi

    if PGPASSWORD="$DB_PASSWORD" psql \
        "host=$DB_HOST port=$DB_PORT user=$DB_USER dbname=$DB_NAME sslmode=$DB_SSLMODE" \
        -c "SELECT 1;" &> /dev/null
    then
        echo -e "${GREEN}✓ PostgreSQL connection successful${NC}"
        echo -e "  Host: ${DB_HOST}"
        echo -e "  Port: ${DB_PORT}"
        echo -e "  Database: ${DB_NAME}"
        echo -e "  SSL Mode: ${DB_SSLMODE}"
        echo ""

        return 0
    else
        echo -e "${RED}✗ Failed to connect to PostgreSQL${NC}"
        echo -e "${YELLOW}Check DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME and DB_SSLMODE.${NC}"
        echo ""

        return 1
    fi
}


# ------------------------------------------------------------
# Run application locally
# ------------------------------------------------------------
run_local() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Starting backend locally...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if ! command -v go &> /dev/null; then
        echo -e "${RED}✗ Go is not installed or not in PATH.${NC}"
        exit 1
    fi

    echo -e "${BLUE}Running:${NC} go run cmd/api/main.go"
    echo ""

    go run cmd/api/main.go
}

# ------------------------------------------------------------
# Run application with Docker
# ------------------------------------------------------------
run_docker() {
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}Starting backend with Docker...${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""

    if ! command -v docker &> /dev/null; then
        echo -e "${RED}✗ Docker is not installed.${NC}"
        exit 1
    fi

    if docker compose version &> /dev/null; then
        docker compose up -d
    elif command -v docker-compose &> /dev/null; then
        docker-compose up -d
    else
        echo -e "${RED}✗ Docker Compose is not installed.${NC}"
        exit 1
    fi
}

# ------------------------------------------------------------
# Main
# ------------------------------------------------------------
main() {
    local mode="${1:-local}"

    load_env

    # Environment validation
    if ! check_env_vars; then
        echo -e "${YELLOW}Environment validation failed.${NC}"
        echo ""
        read -r -p "Start using Docker instead? (y/n): " use_docker

        if [[ "$use_docker" =~ ^[Yy]$ ]]; then
            run_docker
            exit $?
        else
            echo -e "${RED}Exiting.${NC}"
            exit 1
        fi
    fi

    # Docker mode skips local connection checks
    if [ "$mode" = "docker" ]; then
        run_docker
        exit $?
    fi

    # PostgreSQL check
    check_postgres
    db_status=$?



    echo -e "${YELLOW}========================================${NC}"

    if [ $db_status -eq 0 ]; then
        echo -e "${GREEN}✓ All systems are ready!${NC}"
        echo -e "${YELLOW}========================================${NC}"
        echo ""

        run_local
    else
        echo -e "${RED}✗ One or more connection checks failed.${NC}"
        echo -e "${YELLOW}========================================${NC}"
        echo ""

        echo -e "${YELLOW}What do you want to do?${NC}"
        echo "1) Start with Docker"
        echo "2) Continue locally anyway"
        echo "3) Exit"
        echo ""

        read -r -p "Choose [1-3]: " choice

        case "$choice" in
            1)
                run_docker
                ;;
            2)
                run_local
                ;;
            3)
                echo -e "${RED}Exiting.${NC}"
                exit 1
                ;;
            *)
                echo -e "${RED}Invalid choice.${NC}"
                exit 1
                ;;
        esac
    fi
}

# ------------------------------------------------------------
# Entry point
# ------------------------------------------------------------
case "$1" in
    docker)
        main docker
        ;;
    local|"")
        main local
        ;;
    *)
        echo "Usage:"
        echo "  ./run.sh          Run locally"
        echo "  ./run.sh local    Run locally"
        echo "  ./run.sh docker   Run with Docker"
        exit 1
        ;;
esac