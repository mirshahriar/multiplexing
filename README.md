# multiplexing

## Run TCP Server

```bash
$ go run main.go

Server is running on port 8080
gRPC Server is running on port 8080
```

### HTTP Call

```bash
$ curl --location '127.0.0.1:8080/echo'

echo from HTTP!
```


### gRPC Call

```bash
$ grpcurl -plaintext -d '{"message": "test"}' localhost:8080 echo.EchoService/EchoMessage

{
  "message": "echo test from grpc!"
}
```
