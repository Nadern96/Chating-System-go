# Chat-System-go
 

## prerequisites
before you run services you need the following to be installed:
- **go 1.21**
- **docker**


## Run Locally
### for each service you want to run:
**running directly on local machine:**
- first you need to run the docker-compose for Cassandra and Redis
    ### Cassandra cluster (3 nodes) -> handling distributed database concerns
    ```
    docker-compose -f docker-compose.yaml up
    ``` 
    ### Cassandra Single Node
    ```
    docker-compose -f docker-compose-cassandra-single.yaml up
    ```
- load the required env vars with their corresponding values
    ### for bff
    ```
    export SERVICE_AUTH_URL="localhost:50051"
    export SERVICE_CHAT_URL="localhost:50052"
    export PORT="8080"
    ```
    ### for service-auth:
    ```
    export GRPC_PORT=50051
    export ENVIRONMENT=local

    export CASSANDRA_URL=localhost
    export CASSANDRA_KEYSPACE=auth
    export CASSANDRA_USERNAME=""
    export CASSANDRA_PASSWORD=""

    export JWT_SECRET_KEY="JWT_SECRET_AUTH_KEY_2024_CME"

    export REDIS_URL="localhost:6379"
    ```

    ### for service-auth:
    ```
    export GRPC_PORT=50052
    export ENVIRONMENT=local

    export CASSANDRA_URL=localhost
    export CASSANDRA_KEYSPACE=chat
    export CASSANDRA_USERNAME=""
    export CASSANDRA_PASSWORD=""

    export REDIS_URL="localhost:6379"
    ```
- inside the service directory, use:
    ```
    go run main.go
    ```

## Running Unit Tests
### running every service test cases:
after exporting the service env vars
inside the service directory, use:

```
go test ./service-chat/... 
```

```
go test ./service-auth/... 
```

![architecture](architecture.png)
