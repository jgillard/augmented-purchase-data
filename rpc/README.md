# gRPC Implementation

## Compile

`protoc --go_out=plugins=grpc:. *.proto`

## Testing

`docker run -v ~/go/src/github.com/jgillard/practising-go-tdd/rpc:/rpc lequoctuan/grpcc --insecure --proto /rpc/service.proto --address host.docker.internal:7777 --exec /rpc/grpcc_tests.js`

prints for the first call

```
{
  "status": "ok"
}
```