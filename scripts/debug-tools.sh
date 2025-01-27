#!/bin/bash

# Set strict error handling
set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to check pod logs
check_pod_logs() {
    local namespace=${1:-social-network}
    local service=${2:-}
    local tail_lines=${3:-100}

    if [ -z "$service" ]; then
        log_error "Service name is required"
        echo "Usage: $0 logs <service-name> [tail-lines]"
        return 1
    fi

    log_info "Fetching logs for $service in namespace $namespace..."
    kubectl logs -n "$namespace" -l app="$service" --tail="$tail_lines" --all-containers=true
}

# Function to check pod status
check_pod_status() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_info "Checking all pods in namespace $namespace"
        kubectl get pods -n "$namespace" -o wide
    else
        log_info "Checking pods for service $service in namespace $namespace"
        kubectl get pods -n "$namespace" -l app="$service" -o wide
    fi
}

# Function to check resource usage
check_resource_usage() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_info "Checking resource usage for all pods in namespace $namespace"
        kubectl top pods -n "$namespace"
    else
        log_info "Checking resource usage for service $service in namespace $namespace"
        kubectl top pods -n "$namespace" -l app="$service"
    fi
}

# Function to check database connections
check_db_connections() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_error "Service name is required"
        echo "Usage: $0 db-check <service-name>"
        return 1
    fi

    log_info "Checking database connections for service $service..."
    kubectl exec -n "$namespace" -l app="$service" -- netstat -ant | grep 5432
}

# Function to check service dependencies
check_dependencies() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_error "Service name is required"
        echo "Usage: $0 deps <service-name>"
        return 1
    fi

    log_info "Checking service dependencies for $service..."
    kubectl describe deployment -n "$namespace" "$service" | grep -A 5 "Environment:"
}

# Function to check network connectivity
check_network() {
    local namespace=${1:-social-network}
    local service=${2:-}
    local target_service=${3:-}
    local port=${4:-80}

    if [ -z "$service" ] || [ -z "$target_service" ]; then
        log_error "Both source and target service names are required"
        echo "Usage: $0 network <source-service> <target-service> [port]"
        return 1
    fi

    log_info "Performing network diagnostics from $service to $target_service..."

    # DNS Resolution Test
    log_info "Testing DNS resolution..."
    kubectl exec -n "$namespace" -l app="$service" -- nslookup "$target_service" || {
        log_error "DNS resolution failed for $target_service"
        return 1
    }

    # TCP Connection Test
    log_info "Testing TCP connection on port $port..."
    kubectl exec -n "$namespace" -l app="$service" -- timeout 5 bash -c "echo >/dev/tcp/$target_service/$port" || {
        log_error "TCP connection failed to $target_service:$port"
        return 1
    }

    # Latency Test
    log_info "Measuring latency..."
    kubectl exec -n "$namespace" -l app="$service" -- ping -c 3 "$target_service" || {
        log_warn "Unable to measure latency (ICMP might be blocked)"
    }

    # HTTP Connectivity Test
    log_info "Testing HTTP connectivity..."
    kubectl exec -n "$namespace" -l app="$service" -- wget -q -O- --timeout=5 "http://$target_service:$port/health" || {
        log_error "HTTP connectivity test failed"
        return 1
    }

    log_info "Network diagnostics completed successfully"
}

# Function to check system metrics
check_metrics() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_error "Service name is required"
        echo "Usage: $0 metrics <service-name>"
        return 1
    fi

    log_info "Fetching metrics for service $service..."
    kubectl exec -n "$namespace" -l app="$service" -- curl -s localhost:9090/metrics | grep "^${service}"
}

# Function to describe pod
describe_pod() {
    local namespace=${1:-social-network}
    local service=${2:-}

    if [ -z "$service" ]; then
        log_error "Service name is required"
        echo "Usage: $0 describe <service-name>"
        return 1
    fi

    log_info "Describing pods for service $service in namespace $namespace"
    kubectl describe pods -n "$namespace" -l app="$service"
}

# Main script
command=${1:-}
shift || true

case "$command" in
    "logs")
        check_pod_logs "social-network" "$@"
        ;;
    "status")
        check_pod_status "social-network" "$@"
        ;;
    "resources")
        check_resource_usage "social-network" "$@"
        ;;
    "describe")
        describe_pod "social-network" "$@"
        ;;
    "db-check")
        check_db_connections "social-network" "$@"
        ;;
    "deps")
        check_dependencies "social-network" "$@"
        ;;
    "network")
        check_network "social-network" "$@"
        ;;
    "metrics")
        check_metrics "social-network" "$@"
        ;;
    *)
        echo "Usage: $0 <command> [options]"
        echo "Commands:"
        echo "  logs <service-name> [tail-lines]  - Check logs for a service"
        echo "  status [service-name]            - Check pod status"
        echo "  resources [service-name]         - Check resource usage"
        echo "  describe <service-name>          - Describe pods for a service"
        echo "  db-check <service-name>          - Check database connections"
        echo "  deps <service-name>              - Check service dependencies"
        echo "  network <src-svc> <dst-svc>     - Check network connectivity"
        echo "  metrics <service-name>          - Check service metrics"
        exit 1
        ;;
esac