# multiplexing

## Run TCP Server

```bash
$ go run main.go

Server is running on port 8080
gRPC Server is running on port 8080
```

### HTTP Call

```bash
$ curl --location --request GET '127.0.0.1:8080/echo' \
--header 'Content-Type: application/json' \
--data '{
    "message": "test"
}'
```
#### Response

```json
{
  "from": "http",
  "response": {
    "message": "echo test from grpc!"
  }
}
```


### gRPC Call

```bash
$ grpcurl -plaintext -d '{"message": "test"}' localhost:8080 echo.EchoService/EchoMessage
```

#### Response
```json
{
  "message": "echo test from grpc!"
}
```
