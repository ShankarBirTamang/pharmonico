# Kafka Verification Guide (Task 5.4)

This guide explains how to verify Kafka producer and consumer functionality using the test programs and Kafka UI.

## Overview

Task 5.4 involves three verification steps:
1. **5.4.1**: Produce a test event to Kafka
2. **5.4.2**: Consume a test event from Kafka
3. **5.4.3**: Validate event trace in Kafka UI

---

## Prerequisites

1. **Docker Compose running**: Ensure all services are up
   ```bash
   docker-compose up -d
   ```

2. **Kafka UI accessible**: Open http://localhost:8085 in your browser

3. **Go environment**: Make sure Go is installed and the project dependencies are available

4. **Kafka Broker Address**: 
   - When running from host machine: Uses `localhost:29092` by default
   - When running inside Docker: Uses `kafka:9092` (set via environment variable)
   - You can override by setting `KAFKA_BROKERS` environment variable:
     ```bash
     export KAFKA_BROKERS=localhost:29092
     ```

---

## Step 1: Produce Test Event (5.4.1)

### Run the Producer Test

```bash
cd backend-go
go run cmd/test-kafka-producer/main.go
```

### Expected Output

```
🚀 Testing Kafka Producer - Producing Test Event...
============================================================

📡 Connecting to Kafka...
   ✓ Using topic: prescription.intake.received

📤 Producing test event...
   Event ID: test-event-20240115-143022
   Event Type: test.verification
   Topic: prescription.intake.received
   Payload: {"event_id":"test-event-...","event_type":"test.verification",...}

✅ Test event produced successfully!

📊 Event Details:
   ✓ Topic: prescription.intake.received
   ✓ Message Key: test-rx-12345
   ✓ Event ID: test-event-20240115-143022
   ✓ Timestamp: 2024-01-15T14:30:22Z
   ✓ Payload Size: 234 bytes
```

### What Happens

- Creates a Kafka producer connection
- Publishes a test event to the `prescription.intake.received` topic
- Uses the prescription ID as the message key for partitioning
- Event includes metadata: event ID, timestamp, test data, prescription ID, patient ID

---

## Step 2: Consume Test Event (5.4.2)

### Run the Consumer Test

```bash
cd backend-go
go run cmd/test-kafka-consumer/main.go
```

### Expected Output

```
👂 Testing Kafka Consumer - Consuming Test Event...
============================================================

📡 Connecting to Kafka...
   ✓ Subscribing to topic: prescription.intake.received
   ✓ Consumer Group: test-consumer-group

🔍 Polling for messages (timeout: 30 seconds)...
   Waiting for test event...

✅ Message received!

📨 Message Details:
   ✓ Topic: prescription.intake.received
   ✓ Partition: 0
   ✓ Offset: 42
   ✓ Key: test-rx-12345
   ✓ Value Size: 234 bytes

📋 Event Payload:
   ✓ Event ID: test-event-20240115-143022
   ✓ Event Type: test.verification
   ✓ Timestamp: 2024-01-15T14:30:22Z
   ✓ Test Data: This is a test event for Kafka verification (Task 5.4.1)
   ✓ Prescription ID: test-rx-12345
   ✓ Patient ID: test-patient-67890

📊 Kafka Message Info:
[Kafka] Received: topic=prescription.intake.received, partition=0, offset=42, key=test-rx-12345
   ✓ Message committed successfully

============================================================
✅ Test event consumed successfully!
```

### What Happens

- Creates a Kafka consumer connection
- Subscribes to the `prescription.intake.received` topic
- Polls for messages (with 30-second timeout)
- Receives and parses the test event
- Commits the message offset

---

## Step 3: Validate Event Trace in Kafka UI (5.4.3)

### Access Kafka UI

1. **Open Kafka UI**: Navigate to http://localhost:8085

2. **Select Cluster**: The cluster "phil-my-meds-cluster" should be visible

### View Topics

1. **Navigate to Topics**: Click on "Topics" in the left sidebar

2. **Find Test Topic**: Look for `prescription.intake.received` in the topics list

3. **View Topic Details**: Click on the topic name

### Verify Message

1. **Check Partitions**: You should see 3 partitions (as configured in docker-compose.yml)

2. **View Messages**: 
   - Click on the "Messages" tab
   - You should see the test event in one of the partitions
   - The message should show:
     - **Offset**: The message offset (e.g., 42)
     - **Key**: `test-rx-12345`
     - **Value**: JSON payload with event details
     - **Timestamp**: When the message was produced

3. **Message Details**:
   - Click on a message to see full details
   - Verify the JSON structure matches the test event:
     ```json
     {
       "event_id": "test-event-20240115-143022",
       "event_type": "test.verification",
       "timestamp": "2024-01-15T14:30:22Z",
       "test_data": "This is a test event for Kafka verification (Task 5.4.1)",
       "prescription_id": "test-rx-12345",
       "patient_id": "test-patient-67890"
     }
     ```

### Check Consumer Groups

1. **Navigate to Consumer Groups**: Click on "Consumer Groups" in the left sidebar

2. **Find Consumer Group**: Look for `test-consumer-group`

3. **View Group Details**: 
   - Click on the consumer group
   - Verify it shows:
     - **Topic**: `prescription.intake.received`
     - **Partition**: The partition number (0, 1, or 2)
     - **Current Offset**: Should match the consumed message offset
     - **Lag**: Should be 0 (no unprocessed messages)

### Verify Event Flow

1. **Check Topic Metrics**:
   - View the topic's metrics tab
   - Verify message count increased
   - Check partition distribution

2. **View Consumer Lag**:
   - In Consumer Groups, verify lag is 0
   - This confirms the message was consumed successfully

---

## Troubleshooting

### Producer Issues

**Problem**: Producer fails to connect
- **Solution**: Ensure Kafka container is running: `docker ps | grep kafka`
- **Solution**: Check Kafka brokers configuration matches docker-compose.yml

**Problem**: Message not appearing in Kafka UI
- **Solution**: Wait a few seconds for Kafka UI to refresh
- **Solution**: Check Kafka logs: `docker logs phil-my-meds-kafka`

### Consumer Issues

**Problem**: Consumer times out (no messages)
- **Solution**: Make sure you ran the producer test first
- **Solution**: Check that both are using the same topic name
- **Solution**: Verify Kafka is running and accessible

**Problem**: Consumer receives old messages
- **Solution**: This is expected if `StartOffset` is set to `kafka.FirstOffset`
- **Solution**: The test consumer uses `kafka.LastOffset` to only get new messages

### Kafka UI Issues

**Problem**: Kafka UI not accessible
- **Solution**: Check container is running: `docker ps | grep kafka-ui`
- **Solution**: Verify port 8085 is not in use by another service
- **Solution**: Check logs: `docker logs phil-my-meds-kafka-ui`

**Problem**: Topic not visible in Kafka UI
- **Solution**: Refresh the page
- **Solution**: Check topic was created: `docker exec phil-my-meds-kafka kafka-topics --bootstrap-server localhost:9092 --list`

---

## Verification Checklist

- [ ] Producer test runs successfully
- [ ] Consumer test receives the event
- [ ] Event appears in Kafka UI under the correct topic
- [ ] Message details match the produced event
- [ ] Consumer group shows correct offset
- [ ] Consumer lag is 0 (message processed)
- [ ] Event JSON structure is valid
- [ ] Message key and value are correct

---

## Next Steps

After completing Task 5.4:

1. **Integration**: Use the producer/consumer in actual services
2. **Error Handling**: Add retry logic and dead letter queue handling
3. **Monitoring**: Set up metrics and alerts for Kafka operations
4. **Testing**: Create integration tests for Kafka workflows

---

## Additional Resources

- [Kafka UI Documentation](https://github.com/provectus/kafka-ui)
- [Kafka Go Library (segmentio/kafka-go)](https://github.com/segmentio/kafka-go)
- [Kafka Documentation](https://kafka.apache.org/documentation/)

