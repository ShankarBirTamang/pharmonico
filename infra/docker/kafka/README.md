# Kafka Configuration (KRaft Mode)

Apache Kafka provides event streaming for Pharmonico's async workflows.
Running in **KRaft mode** - no Zookeeper dependency required.

## What is KRaft?

KRaft (Kafka Raft) is Kafka's built-in consensus mechanism that replaces the dependency on Apache Zookeeper. Benefits include:
- Simpler architecture (single system to manage)
- Faster startup and recovery times
- Better scalability for metadata
- Reduced operational complexity

## Purpose
- Event-driven prescription processing
- Decouples API from workers
- Enables retry and dead-letter queue patterns

## Topics Used
| Topic | Purpose | Partitions |
|-------|---------|------------|
| `intake_received` | New prescriptions from intake | 3 |
| `validate_prescription` | Validation requests | 3 |
| `validation_completed` | Validation results | 3 |
| `enrollment_requested` | Patient enrollment requests | 3 |
| `enrollment_completed` | Patient enrollment events | 3 |
| `pharmacy_recommendation_requested` | Routing requests | 3 |
| `pharmacy_routed` | Pharmacy assignment events | 3 |
| `prescription_fulfilled` | Fulfillment completion events | 3 |
| `dead_letter_queue` | Failed message handling | 1 |

## Default Settings
- Broker Port: 9092 (internal), 29092 (external)
- Controller Port: 9093 (internal, for KRaft consensus)
- Default Replication Factor: 1 (dev only)
- Default Partitions: 3

## KRaft Environment Variables
| Variable | Description |
|----------|-------------|
| `KAFKA_NODE_ID` | Unique node identifier |
| `KAFKA_PROCESS_ROLES` | Node roles: `broker`, `controller`, or `broker,controller` |
| `CLUSTER_ID` | Unique cluster identifier (generated once) |
| `KAFKA_CONTROLLER_QUORUM_VOTERS` | Controller quorum configuration |
| `KAFKA_CONTROLLER_LISTENER_NAMES` | Listener name for controller communication |

## Health Check
Kafka is healthy when topics can be listed via `kafka-topics --bootstrap-server localhost:9092 --list`.

## Production Notes
- Use 3+ brokers for high availability
- Separate controller and broker roles for large clusters
- Set replication factor to 3
- Configure proper retention policies per topic
- Enable TLS and SASL authentication
- Generate a unique `CLUSTER_ID` for each environment

## Generating a Cluster ID
```bash
# Generate a new cluster ID (only needed once per cluster)
docker run --rm confluentinc/cp-kafka:7.5.0 kafka-storage random-uuid
```

## Migration from Zookeeper
If migrating from an existing Zookeeper-based setup:
1. Ensure Kafka version 3.3+ (Confluent 7.3+)
2. Follow Apache Kafka's migration guide
3. Generate new cluster ID for fresh installs
