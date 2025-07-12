# GoFleet

**GoFleet** is a simple load balancer written in Go that distributes incoming work requests across multiple backend nodes using a **round-robin** strategy. Each node processes jobs using a **worker pool** and a **queue system** to simulate real-world load handling.

---

## 🚀 Features

- Round-robin load balancing
- Node health checks via `/status`
- Job queuing with channels
- Worker pool using goroutines
- Simulated CPU-heavy workload: finding all prime numbers up to a random `n`

---

## 🧠 How It Works

1. The **load balancer** listens on port `3000` and forwards `/work` requests to healthy nodes (`3001`, `3002`, ...).
2. Each **node** maintains:
   - A **work queue** (buffered channel)
   - A **pool of worker goroutines**
3. Workers pick jobs from the queue and run the **FindPrimes** function.

---

## 🛠️ Run It Locally

```bash
# Build everything
make build

# Start load balancer and all nodes
make start-all

# Run load tests
make run-tests

# Stop everything
make stop-all
```

---

## 📂 Project Structure

```
gofleet/
├── loadbalancer/     # Handles request routing
├── nodes/            # Backend workers (node1, node2)
├── tests/            # Load testing tool
├── Makefile
└── README.md
```

---

## 🔧 Future Ideas

- Round-robin with weight/fairness
- Dockerized deployment
- Web dashboard with metrics
