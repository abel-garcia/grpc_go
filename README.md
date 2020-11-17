# grpc-golang

### What is gRPC?

gRPC is a framework to connect different applications, where a client application can directly call a method on a server application on a different machine as if it were a local object, Doesn't matter the environments of the applications, this uses Protocol Buffers, making it easier for you to create distributed applications and services.

[documentation...](https://github.com/grpc/grpc-go)


****

#### Application

This example have 2 applications, server and client where the server expose some methods to management an address book, the methods used are PUT, GET, DELETE, which are called from the application client using gRPC and Protocol Buffers.


![Ilustration](./img/Simple%20application%20gRPC.png)

#### Content

```bash
.
├── README.md
├── client
│   └── gRPClient.go
├── config.json.example
├── go.mod
├── go.sum
├── main.go
├── server
│   └── gRPCServer.go
└── test
    └── srv_test.go
```

#### Application config.json

Change the port and host on config.json

```json
{
    "server":{
        "port": 9000,
        "host": "localhost"
    }
}
```

#### Running test

Running only Benchmark

```bash
> go test -v -bench=. ./...
```

Running all test 

```bash
> go test -v ./...
```