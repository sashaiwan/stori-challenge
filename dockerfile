FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stori-challenge .

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=builder /app/stori-challenge .
COPY --from=builder /app/templates ./templates

RUN mkdir -p /data

CMD [ "./stori-challenge" ]
