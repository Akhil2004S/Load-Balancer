# HTTP Load Balancer
A Load Balancer created in go leveraging the power of concurrency. It has CLI integrated, which will make it easier to access and metrics to know the performance.
## What it does
The Load balancer gets the required config from the config.yaml file using viper. The config file requires the port at which the load balacer should start, the algorithm to choose and the server's addresses. The load balancer has three main algorithms right now and they are Weighted Least Response time, Least connections and IP Hash. The algorithm can either be chosen in the CLI or in the config file.  
## Architecture
## Design Decision
Generally the load balancers I know implement a single algorithm and force the user to choose the same. If the user is not technical, then they might not even know there are other ways to distribute the traffic. So my load balancer gives users the option to pick one of the three algorithms through CLI flag or via the config file. If no flag and config is set, then the default algorithm (Least Connections) will be used. The algorithms used are IP Hashing (Static), Least connections(Dynamic) and Weighted Least Response Time(Dynamic).
## Why these algorithms
Least connection is the simple and easy to understand approach where we don't want to flood the same system with many requests. It is a dynamic approach that distributes based on current traffic.\
Weighted Least response time is a practical approach that relies on the performance of the server for the distribution. This is dynamic approach as well. This algorithm uses EMA for the calcuation for average response time. EMA stands for Exponential Moving Average. It is a mathematical approach that ensures that the average is not too much biased by a single bad value that occured some time during the process. It uses an alpha value that acts as the confidence factor. Using this formula, the average response time is calculated.\
Static IP is a static algorithm used to maintain sticky sessions and to ensure request from the same client IP is handled by same server.
## Health Check
Each server consists of three health states.
1. Healthy - Depicting that the server is healthy and can handle requests sent to it.
2. Unhealthy - Depicting that the server is not working and is out of the rotation.
3. Evaluating - This means that the server is in a transition state. For a server to be marked as healthy or unhealthy, it has to pass the respective threshold. If the healthy threshold is 5, then it has to respond with 200 stauts code in order to be marked as healthy. During this checking period, the server is out of rotation and is marked as evaluating.
The health checker runs every second and pings each server's /health endpoint. If a server goes down between health checks, the reactive fallback in the request handler catches the failure and marks it as evaluating immediately.

## Graceful Shutdown
Graceful shutdown ensures that the load balancer exits properly when interupped. So when Ctrl+C is received, the health checker stops and the load balancer exits. The requests that are queued are dropped. This is a known limitation and can be fixed with http.Server().Shutdown().

## Metrics
Metrics of the load balancer is obtained using the cli command. Internally, a gRPC server is started when the load balancer starts which responds with the current metrics. The metrics can easily be obtained using API request like /metrics to the load balancer easily. But instead I have implemented a gRPC server to experience the protocol used in production systems.  

## How to run
### Prerequisites
Go 1.21+

### Installation
```
git clone https://github.com/Akhil2004S/Load-Balancer
cd Load-Balancer
go build -o loadbalancer .
```
### Configuration
Create a config.yaml in the project root
```
port: 3000
algorithm: "least_connections"
alpha: 0.1
servers:
  - address: "http://127.0.0.1:8080"
  - address: "http://127.0.0.1:8081"
```
Available algorithms: least_connections, least_response_time, ip_hash
### Usage
```
# Start with config file
./loadbalancer start

# Override port or algorithm via CLI flags
./loadbalancer start --port 4000 --algorithm least_response_time

# Query real-time metrics (load balancer must be running)
./loadbalancer metrics
```
