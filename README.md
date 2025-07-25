# checkpoint
API health monitoring for production environments. Written in Golang for speed and concurrency.

# Features
- **Fast**: Written in Golang, designed for high performance and low latency.
- **Concurrent**: Utilises Go's goroutines for efficient concurrent processing.
- **Configurable**: Easily configure endpoints and monitoring parameters via a YAML file.
- **HTTP Support**: Supports GET and POST requests with custom headers.

# HTTP Connection Design
In order to efficiently handle HTTP connections, the design is layed out to utilise a connection pool.
HTTP client warming is also used, so that timed requests are not delayed by connection setup time. The healthchecker effectively replicates the behaviour of a long-lived HTTP client; the average is not thrown off by the initial connection setup time.
