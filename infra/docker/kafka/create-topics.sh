#!/bin/bash
# ============================================
# Kafka Topic Initialization Script
# Run after Kafka broker is ready
# ============================================

KAFKA_BROKER=${KAFKA_BROKER:-localhost:9092}
PARTITIONS=${PARTITIONS:-3}
REPLICATION=${REPLICATION:-1}

echo "🚀 Creating Kafka topics..."

# Wait for Kafka to be ready
echo "⏳ Waiting for Kafka broker at $KAFKA_BROKER..."
until kafka-topics.sh --bootstrap-server $KAFKA_BROKER --list > /dev/null 2>&1; do
    echo "Kafka not ready yet..."
    sleep 2
done
echo "✅ Kafka is ready!"

# Define topics
TOPICS=(
    "intake_received"
    "validate_prescription"
    "validation_completed"
    "enrollment_requested"
    "enrollment_completed"
    "pharmacy_recommendation_requested"
    "pharmacy_routed"
    "prescription_fulfilled"
    "dead_letter_queue"
)

# Create each topic
for TOPIC in "${TOPICS[@]}"; do
    echo "📝 Creating topic: $TOPIC"
    kafka-topics.sh --bootstrap-server $KAFKA_BROKER \
        --create \
        --if-not-exists \
        --topic "$TOPIC" \
        --partitions $PARTITIONS \
        --replication-factor $REPLICATION
done

echo ""
echo "✅ All topics created! Current topics:"
kafka-topics.sh --bootstrap-server $KAFKA_BROKER --list

