#!/bin/bash

# =====================================================
# Pemdes Payroll - Quick Deploy Script
# =====================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        print_error "Please run as root (use sudo)"
        exit 1
    fi
}

# Check Docker installation
check_docker() {
    print_info "Checking Docker installation..."
    if ! command -v docker &> /dev/null; then
        print_warning "Docker not found. Installing Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh get-docker.sh
        systemctl enable docker
        systemctl start docker
        rm get-docker.sh
    else
        print_info "Docker is already installed"
    fi
}

# Check Docker Compose installation
check_docker_compose() {
    print_info "Checking Docker Compose installation..."
    if ! command -v docker-compose &> /dev/null; then
        print_warning "Docker Compose not found. Installing..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    else
        print_info "Docker Compose is already installed"
    fi
}

# Create .env file if not exists
setup_env() {
    if [ ! -f .env ]; then
        print_info "Creating .env file..."
        cp .env.example .env

        # Generate random password
        MYSQL_PASSWORD=$(openssl rand -base64 16)
        JWT_SECRET=$(openssl rand -base64 32)

        sed -i "s/MYSQL_PASSWORD=.*/MYSQL_PASSWORD=$MYSQL_PASSWORD/" .env
        sed -i "s/JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" .env

        print_warning "Please update .env file with your preferences"
        print_info "Generated passwords saved to .env"
    else
        print_info ".env file already exists"
    fi
}

# Verify main.go exists
verify_main_go() {
    if [ ! -f main.go ]; then
        print_error "main.go not found in project root!"
        print_info "Please ensure main.go exists before deploying"
        exit 1
    fi
    print_info "main.go found"
}

# Build and start containers
deploy() {
    print_info "Building Docker images..."
    docker-compose build

    print_info "Starting containers..."
    docker-compose up -d

    print_info "Waiting for services to be ready..."
    sleep 10

    print_info "Checking container status..."
    docker-compose ps
}

# Show status
show_status() {
    echo ""
    print_info "=== DEPLOYMENT STATUS ==="
    echo ""

    # Show running containers
    print_info "Running containers:"
    docker-compose ps

    echo ""
    print_info "=== ACCESS INFORMATION ==="
    echo ""

    # Get server IP
    SERVER_IP=$(hostname -I | awk '{print $1}')

    echo "Frontend: http://$SERVER_IP (or http://localhost)"
    echo "Backend API: http://$SERVER_IP/api"
    echo ""

    print_info "To view logs, run: docker-compose logs -f"
    print_info "To stop, run: docker-compose stop"
    echo ""
}

# Main execution
main() {
    echo "=========================================="
    echo "  Pemdes Payroll - Quick Deploy Script"
    echo "=========================================="
    echo ""

    check_root
    check_docker
    check_docker_compose
    setup_env
    verify_main_go
    deploy
    show_status

    print_info "Deployment completed successfully!"
    print_warning "Please save your database credentials from .env file"
}

# Run main function
main
