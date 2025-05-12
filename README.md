# Stori Technical Challenge

## Overview

Transaction processing system that reads CSV files and sends email summaries.

## Features

- CSV transaction processing
- PostgreSQL database storage
- HTML email notifications
- Docker containerization

## Requirements

- Docker
- Docker Compose
- Go 1.21+

## Setup

1. Clone the repository
2. `cd stori-challenge`
3. Create `.env` file with email credentials: `cp .env.example .env`
4. Run `docker-compose up`

## Development

1. Place your CSV file in the `data` directory
2. In a terminal, run `go run .`

## API/Interface

### Design Approach

I've implemented two architectural approaches:

1. **Naive Approach** (Current Implementation) - Simpler and straightforward that accomplished all the requirements and offer a fast release/feedback cycle
2. **Event-Driven CQRS Approach** (Planned - In progress) - A scalable solution addressing security and compliance concerns

### Architecture

The implementation follows Clean Architecture principles without strict layering:

- **Domain Layer**: Core business entities (`Transaction`, `TransactionStats`)
- **Application Layer**: Use cases (`processCSV`, `getTransactionStats`)
- **Infrastructure Layer**: Database and email implementations
- **Interface Layer**: Main entry point and configuration

Design decisions and potential improvements are documented throughout the codebase via TODO comments.

### Key Functions

| Function | Purpose |
|----------|---------|
| `processCSV()` | Parse and validate transaction data |
| `getTransactionStats()` | Calculate financial summaries |
| `sendEmail()` | Generate and send HTML reports |
| `saveToDatabase()` | Persist data with account management |

### Environment Variables

Configuration is managed through environment variables, with defaults where appropriate. See `.env.example` for required variables.

## Database Schema

There are two tables

## Testing

TBD...
