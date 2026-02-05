#!/bin/bash

################################################################################
# ServQR Platform - System Prerequisites Installation
#
# Installs required system packages (excluding Docker which is handled separately)
#
# Usage: sudo bash install-prerequisites.sh
################################################################################

set -e
set -u

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }
warn() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# Detect OS
detect_os() {
    if [[ -f /etc/os-release ]]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    else
        error "Cannot detect operating system"
    fi
    
    log "Detected OS: $OS $VER"
}

# Install Go 1.23+
install_go() {
    log "Installing Go 1.23..."
    
    # Check if Go is already installed
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        log "Go already installed: $GO_VERSION"
        if [[ "$GO_VERSION" > "1.23" ]] || [[ "$GO_VERSION" == "1.23"* ]]; then
            log "Go version is sufficient"
            return 0
        else
            warn "Go version is old, upgrading..."
        fi
    fi
    
    # Download and install Go
    GO_VERSION="1.23.6"
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) GO_ARCH="amd64" ;;
        aarch64|arm64) GO_ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
    
    cd /tmp
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" || error "Failed to download Go"
    
    # Remove old Go installation
    rm -rf /usr/local/go
    
    # Extract new Go
    tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # Add to PATH
    if ! grep -q "/usr/local/go/bin" /etc/profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    fi
    
    export PATH=$PATH:/usr/local/go/bin
    
    # Verify installation
    /usr/local/go/bin/go version || error "Go installation failed"
    log "Go installed successfully: $(/usr/local/go/bin/go version)"
    
    # Cleanup
    rm -f "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
}

# Install Node.js 20+
install_nodejs() {
    log "Installing Node.js 20..."
    
    # Check if Node.js is already installed
    if command -v node &> /dev/null; then
        NODE_VERSION=$(node -v | sed 's/v//')
        log "Node.js already installed: v$NODE_VERSION"
        if [[ "${NODE_VERSION%%.*}" -ge 20 ]]; then
            log "Node.js version is sufficient"
            return 0
        else
            warn "Node.js version is old, upgrading..."
        fi
    fi
    
    # Install Node.js using NodeSource repository
    if [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
        apt-get install -y nodejs
    elif [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        curl -fsSL https://rpm.nodesource.com/setup_20.x | bash -
        yum install -y nodejs
    else
        error "Unsupported OS for Node.js installation: $OS"
    fi
    
    # Verify installation
    node -v || error "Node.js installation failed"
    npm -v || error "npm installation failed"
    
    log "Node.js installed successfully: $(node -v)"
    log "npm installed successfully: $(npm -v)"
}

# Install system packages
install_packages() {
    log "Installing system packages..."
    
    if [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        # Update package list
        apt-get update
        
        # Install packages
        apt-get install -y \
            git \
            curl \
            wget \
            build-essential \
            ca-certificates \
            gnupg \
            lsb-release \
            nginx \
            certbot \
            python3-certbot-nginx \
            htop \
            vim \
            net-tools \
            ufw \
            logrotate \
            cron
            
    elif [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        # Update package list
        yum update -y
        
        # Install EPEL repository
        yum install -y epel-release
        
        # Install packages
        yum install -y \
            git \
            curl \
            wget \
            gcc \
            gcc-c++ \
            make \
            ca-certificates \
            nginx \
            certbot \
            python3-certbot-nginx \
            htop \
            vim \
            net-tools \
            firewalld \
            logrotate \
            cronie
            
    else
        error "Unsupported OS: $OS"
    fi
    
    log "System packages installed successfully"
}

# Configure firewall
configure_firewall() {
    log "Configuring firewall..."
    
    if [[ "$OS" == "ubuntu" ]] || [[ "$OS" == "debian" ]]; then
        # UFW firewall
        if command -v ufw &> /dev/null; then
            ufw --force enable
            ufw allow ssh
            ufw allow 80/tcp
            ufw allow 443/tcp
            ufw allow 8081/tcp  # Backend (temporary, will be proxied via Nginx)
            ufw allow 3000/tcp  # Frontend (temporary, will be proxied via Nginx)
            ufw --force reload
            log "UFW firewall configured"
        fi
    elif [[ "$OS" == "centos" ]] || [[ "$OS" == "rhel" ]] || [[ "$OS" == "fedora" ]]; then
        # firewalld
        if command -v firewall-cmd &> /dev/null; then
            systemctl start firewalld
            systemctl enable firewalld
            firewall-cmd --permanent --add-service=ssh
            firewall-cmd --permanent --add-service=http
            firewall-cmd --permanent --add-service=https
            firewall-cmd --permanent --add-port=8081/tcp
            firewall-cmd --permanent --add-port=3000/tcp
            firewall-cmd --reload
            log "firewalld configured"
        fi
    fi
}

# Set up system limits
configure_system_limits() {
    log "Configuring system limits..."
    
    # Increase file descriptor limits
    cat >> /etc/security/limits.conf << 'EOF'
# ServQR Platform - Increased limits
* soft nofile 65536
* hard nofile 65536
* soft nproc 32768
* hard nproc 32768
EOF
    
    # Kernel parameters
    cat >> /etc/sysctl.conf << 'EOF'
# ServQR Platform - Network tuning
net.core.somaxconn = 4096
net.ipv4.tcp_max_syn_backlog = 4096
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_tw_reuse = 1
EOF
    
    sysctl -p
    
    log "System limits configured"
}

# Main installation
main() {
    echo "=========================================================================="
    echo "  ServQR Platform - Prerequisites Installation"
    echo "=========================================================================="
    echo ""
    
    # Check root
    if [[ $EUID -ne 0 ]]; then
        error "This script must be run as root (use sudo)"
    fi
    
    # Detect OS
    detect_os
    
    # Install packages in order
    log "Starting prerequisites installation..."
    
    install_packages
    install_go
    install_nodejs
    configure_firewall
    configure_system_limits
    
    log ""
    log "✓ Prerequisites installation completed successfully!"
    log ""
    log "Installed versions:"
    log "  - Go: $(/usr/local/go/bin/go version)"
    log "  - Node.js: $(node -v)"
    log "  - npm: $(npm -v)"
    log "  - Nginx: $(nginx -v 2>&1 | grep -o '[0-9]*\.[0-9]*\.[0-9]*')"
    log ""
}

main "$@"
