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

## Requirements

- Go 1.21 or later
- Ethereum node access (for production use)

