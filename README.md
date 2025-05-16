# Stori Technical Challenge

## Overview

A transaction processing system that reads CSV files containing credits and debits,
calculates some stats, stores the data in PostgreSQL,
and emails you a nice summary with Stori branding. Built with Go and containerized with Docker.

## Features

- Processes transaction CSVs (with +/- notation for credits/debits)
- Calculates monthly stats and overall balance
- Sends pretty HTML emails with transaction summaries
- Stores everything in PostgreSQL for later reference
- Supports both local files and AWS S3 bucket files
- Containerized with Docker for easy deployment
- AWS ready with Lambda and EC2 support

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.24.3 (for local development)

### Quick Setup

1. Clone the repo

   ```bash
   git clone https://github.com/yourusername/stori-challenge.git
   cd stori-challenge
   ```

2. Configure your environment

   ```bash
   cp .env.example .env
   # Edit .env with your email credentials and other settings
   ```

3. Start the app

   ```bash
   docker-compose up -d
   ```

4. Make a request

    - Local

      ```bash
      curl -X POST http://localhost:8080/transactions/stats \
        -H "Content-Type: application/json" \
        -d '{"email": "your@email.com", "file_content": "the_CSV_in_base64"}'
      ```

    - Or with S3

      ```bash
      curl -X POST http://localhost:8080/transactions/stats \
        -H "Content-Type: application/json" \
        -d '{"bucket": "your-bucket", "key": "file_key", "email": "your@email.com"}'
      ```

Example file in base64: "SWQsRGF0ZSxUcmFuc2FjdGlvbgowLDcvMTUsKzYwLjUKMSw3LzI4LC0xMC4zCjIsNy8yLC0yMC40NgozLDcvMTMsKzEwCjQsOC8yLCsyNS4xNQo1LDgvNCwtMTIuNDUKNiw5LzMwLCsyLjUwCg=="

### Local Development

1. Start just the database

   ```bash
   docker-compose up -d db
   ```

2. Run the app locally

   ```bash
   go run .
   ```

3. Put your CSV file in the data directory (or modify the path in .env)

## API/Interface - Go Service

### POST /transactions/stats

Process transactions and send summary email

**Request Body:**

```json
{
  "bucket": "your-s3-bucket",    // S3 bucket name
  "key": "path/to/file.csv",     // S3 object key
  "email": "recipient@email.com" // Email to send summary
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Transactions processed and email sent successfully"
}
```

## Transaction CSV Format

```csv
Id,Date,Transaction
0,7/15,+60.5    // Credit transaction (+)
1,7/28,-10.3    // Debit transaction (-)
```

### Design Approach

I've implemented two architectural approaches:

1. **Naive Approach** (Current Implementation) - Simpler and straightforward that accomplished all the requirements and offer a fast release/feedback cycle
2. **Event-Driven CQRS Approach** (Planned - In progress) - A scalable solution addressing security and compliance concerns

The Naive approach (far from being trivial), is the solution to go live _as fast as possible_,
there are a few security concerns (like request authentication and authorization) that I left
out of scope in sake of simplicity when looking for a rapid implementation. Making use of the cloud infra leverage some advantage over these concerns, since I could just setup an specific IAM Roles structure and permission over the AWS resources.

The Event-Drive approach is absolutely the _should_ have solution. First of all, it involves compliance and auditability concerns due we are dealing with monetary transactions, because every event is immutable - even more with an Event Store db instead a traditional SQL - and we can replicate the system state at a specific time.

There are a few decisions that I haven't made and answer business questions:

- What should we prioritize? consistency, availability or partition tolerance?
- Will the system being accessed from multiple locations?
- How many teams needs access to the mounted directory? Should we move to a better AWS product like EFS?
- What kind of fault tolerance we are willing to risk?
- And many others...

## Architecture

Follows clean architecture principles with:

- **Domain Models**: Transaction entities and stats
- **Services**: File processing, transaction analysis, email sending
- **Infrastructure**: Database and API interfaces
- **Interface Layer**: HTTP handler and CLI entry points

Design decisions and potential improvements are documented throughout the codebase via TODO comments.

## AWS Deployment

### EC2 Deployment

1. Create required AWS resources (EC2, S3)
2. Run the setup script:

   ```bash
   ./deploy/setup-ec2.sh
   ```

3. Store your config in AWS Secrets Manager as "stori-challenge"

### Lambda Functions

The system includes two AWS Lambda functions in a serverless pipeline:

1. **File Upload Lambda**
   - Receives transaction file in base64-encoded format and recipient email
   - Decodes the file and uploads it to the specified S3 bucket
   - Triggers the ServiceCaller Lambda automatically

2. **ServiceCaller Lambda**
   - Triggered by the upload completion event
   - Calls the Go service API with bucket name, file key, and email
   - The Go service then processes the file, stores data, and sends the email

To deploy the Lambda functions:

```bash
# Configure your AWS credentials first
aws configure
```

Deploy based on your current IAM roles, make sure the FileProcessor has S3 push permissions.

## Todo

- Add a test suite
- Implement CI/CD pipeline
