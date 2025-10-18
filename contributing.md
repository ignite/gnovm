# Contributing

This module is developed using [Ignite CLI](https://docs.ignite.com/), a developer-friendly blockchain development tools for building blockchain applications. To maintain code coherence with Ignite core principles, all development should be done using Ignite CLI commands and conventions.

## Prerequisites

1. Install [Ignite CLI](https://docs.ignite.com/guide/install):

   ```bash
   curl https://get.ignite.com/cli@latest! | bash
   ```

2. Ensure you have Go 1.25+ installed

## Development Workflow

### Getting Started

This project is an Ignite CLI-based blockchain module. All code generation and scaffolding should be done using Ignite commands to maintain consistency with the framework's architecture.

```bash
# Clone and navigate to the project
git clone <repository-url>
cd gnovm

# Serve the blockchain locally
ignite chain serve
```

### Creating New Components

**Always use Ignite CLI commands** to generate new blockchain components. This ensures proper integration with the Cosmos SDK and maintains code consistency.

#### Messages (Transactions)

```bash
# Create a new message type
ignite scaffold message [message-name] [field1] [field2:type] --module gnovm

# Example: Create a message to execute Gno code
ignite scaffold message execute-gno code:string --module gnovm
```

#### Queries

```bash
# Create a new query
ignite scaffold query [query-name] [field1] [field2:type] --module gnovm

# Example: Create a query to get Gno VM state
ignite scaffold query vm-state address:string --module gnovm
```

#### Types

```bash
# Create a new type
ignite scaffold type [type-name] [field1] [field2:type] --module gnovm

# Example: Create a type for Gno packages
ignite scaffold type gno-package path:string code:string --module gnovm
```

#### Maps and Lists

```bash
# Create a map structure
ignite scaffold map [map-name] [field1] [field2:type] --module gnovm

# Create a list structure
ignite scaffold list [list-name] [field1] [field2:type] --module gnovm
```

### Code Organization

The project follows Ignite CLI's standard module structure:

```
x/gnovm/
├── ante/          # AnteHandler logic
├── client/        # CLI commands and REST endpoints
├── keeper/        # Core business logic
├── module/        # Module definition and lifecycle
├── simulation/    # Simulation and testing
└── types/         # Protobuf types and interfaces
```

### Testing

Use Ignite's testing framework and follow test-driven development:

```bash
# Run all tests
ignite chain test

# Run specific module tests
go test ./x/gnovm/...

# Run with coverage
go test -cover ./x/gnovm/...
```

### Protocol Buffers

When modifying `.proto` files, regenerate the Go code using Ignite:

```bash
# Regenerate protobuf files
ignite generate proto-go
```

**Never manually edit generated `.pb.go` files** - always modify the source `.proto` files and regenerate.

### Configuration

The testing blockchain configuration is managed in `config.yml`.

### Development Commands

```bash
# Start development server with hot reload
ignite chain serve

# Build the blockchain binary
ignite chain build

# Initialize a new chain
ignite chain init
```

### Contributing Guidelines

1. **Fork and Clone**: Fork the repository and clone your fork
2. **Feature Branch**: Create a feature branch for your changes
3. **Use Ignite**: Generate code using appropriate Ignite commands
4. **Test**: Ensure all tests pass before submitting
5. **Document**: Update documentation for new features
6. **Pull Request**: Submit a PR with clear description of changes

### Code Review Process

- All changes must be made using Ignite CLI commands
- PRs should include appropriate tests
- Follow conventional commit messages
- Ensure documentation is updated
- Code must pass all CI checks

For more information on Ignite CLI features and commands, visit the [official documentation](https://docs.ignite.com/).
