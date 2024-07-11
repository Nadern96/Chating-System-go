
FROM golang:1.21

WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping


ENV PORT 80
ENV GRPC_PORT 50051

ENV ENVIRONMENT ""
ENV CASSANDRA_URL ""
ENV CASSANDRA_KEYSPACE ""
ENV CASSANDRA_USERNAME ""
ENV CASSANDRA_PASSWORD ""
ENV JWT_SECRET_KEY ""
ENV REDIS_URL ""

# Run
CMD ["/app/main"]