# Hadoop Distributed System

A custom implementation of a distributed big data solution written in Go, inspired by Apache Hadoop. This project provides a master-worker architecture for distributed file storage with heartbeat-based health monitoring using gRPC.

## Overview

This is a lightweight distributed file system implementation featuring:
- **Master Node**: Manages cluster metadata, file-to-chunk mappings, and worker node health
- **Worker Nodes**: Store data blocks and respond to heartbeat health checks
- **gRPC Communication**: Efficient inter-node communication protocol
- **Distributed Storage**: Configurable block size and replication factor
- **Console Application**: CLI interface for interacting with the distributed system

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Console Application                      │
│                    (CLI Interface)                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              │                         │
         ┌────▼────┐            ┌──────▼─────┐
         │  Master │            │   Worker   │
         │   Node  │◄──gRPC────►│  Nodes (1-5)│
         └────┬────┘            └────────────┘
              │
    ┌─────────┼─────────┐
    │         │         │
 Metadata  Heartbeat  Node Status
Tracking   Monitoring  Management
```

## Key Components

### gRPC Protocol (`gRPC/`)
- **heartbeat.proto**: Defines the HealthService and message formats for heartbeat communication
- Includes `HeartbeatRequest` with worker ID, timestamp, and resource information
- Node resource tracking: total storage, used storage, free storage

### Master Node (`Nodes/Master/`)
- **nodeMaster.go**: Core master node implementation
  - Maintains cluster metadata (file-to-chunk mappings, chunk-to-node locations)
  - Tracks worker node health and status
  - Manages distributed storage accounting
  - Thread-safe operations using SafeMap
- **heartbeat.go**: Health check mechanism
- **safemap.go**: Thread-safe map implementation for concurrent access

### Worker Nodes (`Nodes/Worker/`)
- **nodeManager.go**: Manages local storage operations
  - Initializes local chunk storage
  - Manages local chunk metadata
  - Tracks storage capacity and usage
- **heartbeat.go**: Sends periodic health updates to master
- **main.go**: Worker node entry point supporting multiple concurrent instances

### Console Application (`ConsolApp/`)
- **main.go**: Interactive CLI for cluster operations
- **file_interface.go**: File system operations interface
- Input processing with configurable block size and replication factor

## Configuration

Default constants (configurable in source):
- **BLOCK_SIZE**: 8096 bytes per block
- **REPLICATION_FACTOR**: 3 copies of each block
- **TOTAL_STORAGE**: 8096 × 10,000 bytes per worker node

## Prerequisites

- **Go 1.26.4** or later
- **Make** build tool
- Linux/Unix-like environment (for shell scripting)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/tymoneq/Hadoop.git
cd Hadoop
```

2. Download dependencies:
```bash
go mod download
```

3. Generate gRPC code (if not already generated):
```bash
cd gRPC
make
cd ..
```

## Building and Running

### Build Everything
```bash
# Build Master Node
cd Nodes/Master
make run
cd ../..

# Build Worker Nodes (5 instances)
cd Nodes/Worker
make run
cd ../..

# Build Console Application
cd ConsolApp
make run
cd ..
```

### Individual Components

**Start Master Node:**
```bash
cd Nodes/Master
go build -o main main.go heartbeat.go safemap.go nodeMaster.go
./main
```

**Start Worker Nodes:**
```bash
cd Nodes/Worker
go build -o node_app main.go heartbeat.go nodeManager.go
./node_app --node-id=node-01 &
./node_app --node-id=node-02 &
./node_app --node-id=node-03 &
./node_app --node-id=node-04 &
./node_app --node-id=node-05 &
```

**Run Console Application:**
```bash
cd ConsolApp
go build main.go file_interface.go
./main
```

## Project Structure

```
Hadoop/
├── go.mod                          # Go module definition
├── README.md                       # This file
├── gRPC/                          # Protocol Buffer definitions
│   ├── heartbeat.proto            # Service and message definitions
│   ├── Makefile                   # Build gRPC code
│   └── pb/                        # Generated Go code
│       ├── heartbeat.pb.go
│       └── heartbeat_grpc.pb.go
├── Nodes/
│   ├── Master/                    # Master node implementation
│   │   ├── main.go
│   │   ├── nodeMaster.go          # Core metadata management
│   │   ├── heartbeat.go           # Health check implementation
│   │   ├── safemap.go             # Thread-safe map
│   │   └── Makefile
│   └── Worker/                    # Worker node implementation
│       ├── main.go
│       ├── nodeManager.go         # Local storage management
│       ├── heartbeat.go           # Health update mechanism
│       ├── data/                  # Distributed data storage
│       │   ├── data-for-node-01/
│       │   ├── data-for-node-02/
│       │   ├── data-for-node-03/
│       │   ├── data-for-node-04/
│       │   └── data-for-node-05/
│       └── Makefile
└── ConsolApp/                     # Command-line interface
    ├── main.go
    ├── file_interface.go          # File operations
    └── Makefile
```

## How It Works

1. **Initialization**: Master node starts and initializes metadata structures. Worker nodes start and register with the master.

2. **Heartbeat Protocol**: Worker nodes periodically send heartbeat messages to the master containing:
   - Worker ID
   - Current timestamp
   - Resource information (storage stats)

3. **Storage Management**: 
   - Files are split into blocks of 8KB each
   - Each block is replicated across 3 worker nodes
   - Master maintains mapping of files to chunks and chunks to nodes

4. **Health Monitoring**: Master tracks worker node status and can detect failures through heartbeat timeouts


## License

Personal project - Open for reference and modification

## Author

tymoneq
