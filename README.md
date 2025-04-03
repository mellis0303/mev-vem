## Project Structure

```
.
├── cmd/
│   ├── mev-omega/        # MEV Omega main application
│   ├── mev-oraclex/      # MEV OracleX main application
│   └── ...               # Other components
├── pkg/
│   ├── mev-omega/        # MEV Omega core package
│   ├── mev-oraclex/      # MEV OracleX core package
│   └── ...               # Other component packages
└── go.mod                # Go module definition
```

## Building

To build any component:

```bash
go build ./cmd/mev-omega
go build ./cmd/mev-oraclex
```

## Running

To run any component:

```bash
./mev-omega
./mev-oraclex
```

## Features

- Transaction dependency resolution
- Dynamic transaction ordering
- Profit optimization
- Front-running protection
- Gas price optimization
- Block space auctioning
- Flashbots integration

## Architecture Overview
This repository is composed of multiple components designed for MEV strategies:
  - **MEV Omega:** Core module for transaction ordering and profit maximization.
  - **MEV OracleX:** Decentralized oracle interactions.
  - **MEV Guardian Engine:** Secure transaction management and transaction decryption.
  - **MEV Arbitrage:** Module for arbitrage detection and execution.

## Requirements

- Go 1.21 or later
- Ethereum node access (for production use)
- Ethereum Node URL (set via the environment variable ETH_NODE_URL)

## Quick Start
1. Set your Ethereum node endpoint:
     export ETH_NODE_URL="https://your-eth-node.com"
2. Build a component (e.g., MEV Omega):
     go build ./cmd/mev-omega
3. Run the component:
     ./mev-omega

## CI/CD
Automated testing and build/deployment are configured in .github/workflows/ci.yml.

