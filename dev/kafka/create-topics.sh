#!/bin/bash
#==============================================================================
# Medical Equipment Platform - Kafka Topic Creation Script
# 
# This script creates all required Kafka topics for the platform with proper
# configuration for healthcare compliance, partitioning strategy, and retention.
#
# Usage: ./create-topics.sh [--reset] [--verify-only] [--list]
#   --reset        : Delete all existing topics before creation
#   --verify-only  : Only verify topics exist without creating
#   --list         : List all existing topics
#==============================================================================

set -e

# Default Kafka configuration
KAFKA_BOOTSTRAP_SERVERS=${KAFKA_BOOTSTRAP_SERVERS:-"kafka:9092"}
ZOOKEEPER_CONNECT=${ZOOKEEPER_CONNECT:-"zookeeper:2181"}
TOPIC_PREFIX=${TOPIC_PREFIX:-""}

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default configuration for topics
DEFAULT_REPLICATION_FACTOR=1  # Set to 3 in production
DEFAULT_PARTITIONS=3
HIGH_VOLUME_PARTITIONS=6
LOW_VOLUME_PARTITIONS=1
# Healthcare data retention - 7 years (in milliseconds)
HEALTHCARE_RETENTION="220752000000"  # 7 years in ms
# Standard retention - 90 days (in milliseconds)
STANDARD_RETENTION="7776000000"  # 90 days in ms
# Short retention - 7 days (in milliseconds)
SHORT_RETENTION="604800000"  # 7 days in ms

#==============================================================================
# Helper Functions
#==============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Kafka is available
check_kafka_available() {
    log_info "Checking Kafka availability at $KAFKA_BOOTSTRAP_SERVERS..."
    
    kafka-broker-api-versions --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS > /dev/null 2>&1
    
    if [ $? -ne 0 ]; then
        log_error "Kafka is not available at $KAFKA_BOOTSTRAP_SERVERS"
        exit 1
    else
        log_success "Kafka is available at $KAFKA_BOOTSTRAP_SERVERS"
    fi
}

# Create a topic with specified configuration
create_topic() {
    local topic_name="$TOPIC_PREFIX$1"
    local partitions=$2
    local retention=$3
    local description="$4"
    
    # Check if topic already exists
    kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS --describe --topic $topic_name > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        log_warning "Topic '$topic_name' already exists. Skipping creation."
    else
        log_info "Creating topic '$topic_name' with $partitions partition(s)..."
        
        kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS \
            --create --topic $topic_name \
            --partitions $partitions \
            --replication-factor $DEFAULT_REPLICATION_FACTOR \
            --config retention.ms=$retention \
            --config cleanup.policy=delete \
            --config min.insync.replicas=1 \
            --config unclean.leader.election.enable=false \
            --config min.cleanable.dirty.ratio=0.01
        
        if [ $? -eq 0 ]; then
            log_success "Topic '$topic_name' created successfully."
            # Add description as a comment in the documentation file
            echo "$topic_name: $description" >> kafka-topics-documentation.txt
        else
            log_error "Failed to create topic '$topic_name'."
            return 1
        fi
    fi
    
    return 0
}

# Create a dead letter queue topic
create_dlq_topic() {
    local source_topic="$1"
    local dlq_topic="${source_topic}.dlq"
    
    create_topic "$dlq_topic" $LOW_VOLUME_PARTITIONS $STANDARD_RETENTION "Dead Letter Queue for $source_topic"
}

# Verify a topic exists
verify_topic() {
    local topic_name="$TOPIC_PREFIX$1"
    
    log_info "Verifying topic '$topic_name'..."
    
    kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS --describe --topic $topic_name > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        log_success "Topic '$topic_name' exists."
        return 0
    else
        log_error "Topic '$topic_name' does not exist."
        return 1
    fi
}

# Delete all topics (for reset)
delete_all_topics() {
    log_warning "Deleting all existing topics..."
    
    # List all topics and delete them
    for topic in $(kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS --list); do
        log_info "Deleting topic '$topic'..."
        kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS --delete --topic $topic
    done
    
    # Wait for deletion to complete
    sleep 5
    log_success "All topics deleted."
}

# List all topics
list_topics() {
    log_info "Listing all topics:"
    kafka-topics --bootstrap-server $KAFKA_BOOTSTRAP_SERVERS --list
}

#==============================================================================
# Topic Definitions
# Format: create_topic "topic_name" partitions retention "description"
#==============================================================================

create_all_topics() {
    # Documentation file header
    echo "# Kafka Topics Documentation" > kafka-topics-documentation.txt
    echo "# Generated on $(date)" >> kafka-topics-documentation.txt
    echo "# Format: topic_name: description" >> kafka-topics-documentation.txt
    echo "" >> kafka-topics-documentation.txt
    
    log_info "Creating health check topic..."
    create_topic "health.check.v1" $LOW_VOLUME_PARTITIONS $SHORT_RETENTION "Health check topic for monitoring Kafka"
    
    log_info "Creating marketplace domain topics..."
    # Marketplace Domain - High volume topics
    create_topic "marketplace.rfq.created.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "RFQ creation events"
    create_topic "marketplace.rfq.updated.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "RFQ update events"
    create_topic "marketplace.rfq.expired.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "RFQ expiration events"
    create_topic "marketplace.quote.created.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Quote creation events"
    create_topic "marketplace.quote.updated.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Quote update events"
    create_topic "marketplace.contract.created.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Contract creation events"
    create_topic "marketplace.catalog.item.created.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Catalog item creation events"
    create_topic "marketplace.catalog.item.updated.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Catalog item update events"
    
    log_info "Creating service domain topics..."
    # Service Domain - High volume topics
    create_topic "service.ticket.created.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Service ticket creation events"
    create_topic "service.ticket.updated.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Service ticket update events"
    create_topic "service.ticket.closed.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Service ticket closure events"
    create_topic "service.asset.registered.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Asset registration events"
    create_topic "service.asset.updated.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Asset update events"
    create_topic "service.qr.deployed.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "QR code deployment events"
    create_topic "service.qr.scanned.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "QR code scan events"
    create_topic "service.workflow.started.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Workflow start events"
    create_topic "service.workflow.completed.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Workflow completion events"
    create_topic "service.workflow.stage.completed.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Workflow stage completion events"
    create_topic "service.diagnostic.started.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Diagnostic procedure start events"
    create_topic "service.diagnostic.completed.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Diagnostic procedure completion events"
    
    log_info "Creating AI/ML domain topics..."
    # AI/ML Domain
    create_topic "ai.triage.completed.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "AI triage completion events"
    create_topic "ai.prediction.created.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "AI prediction creation events"
    create_topic "ai.negotiation.recommendation.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "AI negotiation recommendation events"
    create_topic "ai.dispatch.recommendation.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "AI dispatch recommendation events"
    create_topic "ai.predictive.maintenance.alert.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Predictive maintenance alert events"
    create_topic "ai.demand.forecast.v1" $LOW_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Demand forecast events"
    
    log_info "Creating WhatsApp integration topics..."
    # WhatsApp Integration
    create_topic "whatsapp.message.received.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "WhatsApp message received events"
    create_topic "whatsapp.message.sent.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "WhatsApp message sent events"
    
    log_info "Creating geography domain topics..."
    # Geography Domain
    create_topic "geography.facility.created.v1" $LOW_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Healthcare facility creation events"
    create_topic "geography.facility.updated.v1" $LOW_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Healthcare facility update events"
    create_topic "geography.service.area.updated.v1" $LOW_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Service area update events"
    
    log_info "Creating audit topics..."
    # Audit and Compliance - Longer retention
    create_topic "audit.user.login.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "User login events"
    create_topic "audit.user.logout.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "User logout events"
    create_topic "audit.access.denied.v1" $DEFAULT_PARTITIONS $HEALTHCARE_RETENTION "Access denied events"
    create_topic "audit.data.access.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Data access events"
    create_topic "audit.data.change.v1" $HIGH_VOLUME_PARTITIONS $HEALTHCARE_RETENTION "Data change events"
    
    log_info "Creating notification topics..."
    # Notifications
    create_topic "notification.email.v1" $DEFAULT_PARTITIONS $STANDARD_RETENTION "Email notification events"
    create_topic "notification.sms.v1" $DEFAULT_PARTITIONS $STANDARD_RETENTION "SMS notification events"
    create_topic "notification.push.v1" $DEFAULT_PARTITIONS $STANDARD_RETENTION "Push notification events"
    
    log_info "Creating dead letter queue topics..."
    # Create DLQ topics for high-volume topics
    create_dlq_topic "marketplace.rfq.created.v1"
    create_dlq_topic "marketplace.quote.created.v1"
    create_dlq_topic "service.ticket.created.v1"
    create_dlq_topic "service.ticket.updated.v1"
    create_dlq_topic "service.qr.scanned.v1"
    create_dlq_topic "whatsapp.message.received.v1"
    create_dlq_topic "whatsapp.message.sent.v1"
    create_dlq_topic "audit.user.login.v1"
    create_dlq_topic "audit.data.access.v1"
    create_dlq_topic "audit.data.change.v1"
    
    log_success "All topics created successfully."
    log_info "Topic documentation saved to kafka-topics-documentation.txt"
}

# Verify all topics exist
verify_all_topics() {
    local failed=0
    
    log_info "Verifying all topics..."
    
    # Health check topic
    verify_topic "health.check.v1" || failed=1
    
    # Marketplace domain
    verify_topic "marketplace.rfq.created.v1" || failed=1
    verify_topic "marketplace.rfq.updated.v1" || failed=1
    verify_topic "marketplace.rfq.expired.v1" || failed=1
    verify_topic "marketplace.quote.created.v1" || failed=1
    verify_topic "marketplace.quote.updated.v1" || failed=1
    verify_topic "marketplace.contract.created.v1" || failed=1
    verify_topic "marketplace.catalog.item.created.v1" || failed=1
    verify_topic "marketplace.catalog.item.updated.v1" || failed=1
    
    # Service domain
    verify_topic "service.ticket.created.v1" || failed=1
    verify_topic "service.ticket.updated.v1" || failed=1
    verify_topic "service.ticket.closed.v1" || failed=1
    verify_topic "service.asset.registered.v1" || failed=1
    verify_topic "service.asset.updated.v1" || failed=1
    verify_topic "service.qr.deployed.v1" || failed=1
    verify_topic "service.qr.scanned.v1" || failed=1
    verify_topic "service.workflow.started.v1" || failed=1
    verify_topic "service.workflow.completed.v1" || failed=1
    verify_topic "service.workflow.stage.completed.v1" || failed=1
    verify_topic "service.diagnostic.started.v1" || failed=1
    verify_topic "service.diagnostic.completed.v1" || failed=1
    
    # AI/ML domain
    verify_topic "ai.triage.completed.v1" || failed=1
    verify_topic "ai.prediction.created.v1" || failed=1
    verify_topic "ai.negotiation.recommendation.v1" || failed=1
    verify_topic "ai.dispatch.recommendation.v1" || failed=1
    verify_topic "ai.predictive.maintenance.alert.v1" || failed=1
    verify_topic "ai.demand.forecast.v1" || failed=1
    
    # WhatsApp integration
    verify_topic "whatsapp.message.received.v1" || failed=1
    verify_topic "whatsapp.message.sent.v1" || failed=1
    
    # Geography domain
    verify_topic "geography.facility.created.v1" || failed=1
    verify_topic "geography.facility.updated.v1" || failed=1
    verify_topic "geography.service.area.updated.v1" || failed=1
    
    # Audit topics
    verify_topic "audit.user.login.v1" || failed=1
    verify_topic "audit.user.logout.v1" || failed=1
    verify_topic "audit.access.denied.v1" || failed=1
    verify_topic "audit.data.access.v1" || failed=1
    verify_topic "audit.data.change.v1" || failed=1
    
    # Notification topics
    verify_topic "notification.email.v1" || failed=1
    verify_topic "notification.sms.v1" || failed=1
    verify_topic "notification.push.v1" || failed=1
    
    # DLQ topics
    verify_topic "marketplace.rfq.created.v1.dlq" || failed=1
    verify_topic "marketplace.quote.created.v1.dlq" || failed=1
    verify_topic "service.ticket.created.v1.dlq" || failed=1
    verify_topic "service.ticket.updated.v1.dlq" || failed=1
    verify_topic "service.qr.scanned.v1.dlq" || failed=1
    verify_topic "whatsapp.message.received.v1.dlq" || failed=1
    verify_topic "whatsapp.message.sent.v1.dlq" || failed=1
    verify_topic "audit.user.login.v1.dlq" || failed=1
    verify_topic "audit.data.access.v1.dlq" || failed=1
    verify_topic "audit.data.change.v1.dlq" || failed=1
    
    if [ $failed -eq 0 ]; then
        log_success "All topics verified successfully."
        return 0
    else
        log_error "Some topics failed verification."
        return 1
    fi
}

#==============================================================================
# Main Execution
#==============================================================================

# Parse command line arguments
RESET=false
VERIFY_ONLY=false
LIST_ONLY=false

for arg in "$@"; do
    case $arg in
        --reset)
            RESET=true
            shift
            ;;
        --verify-only)
            VERIFY_ONLY=true
            shift
            ;;
        --list)
            LIST_ONLY=true
            shift
            ;;
        *)
            # Unknown option
            ;;
    esac
done

# Check if Kafka is available
check_kafka_available

# Handle list only option
if [ "$LIST_ONLY" = true ]; then
    list_topics
    exit 0
fi

# Handle reset option
if [ "$RESET" = true ]; then
    delete_all_topics
fi

# Handle verify only option
if [ "$VERIFY_ONLY" = true ]; then
    verify_all_topics
    exit $?
fi

# Create all topics
create_all_topics

# Verify all topics
verify_all_topics

log_info "Script completed successfully."
exit 0
