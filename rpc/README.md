# gRPC Implementation

## Compile

`protoc --go_out=plugins=grpc:. *.proto`

## Testing

`docker run -v ~/go/src/github.com/jgillard/practising-go-tdd/rpc:/rpc lequoctuan/grpcc --insecure --proto /rpc/status.proto --address host.docker.internal:7777 --eval 'client.getStatus({},printReply)'`

returns 

```
{
  "status": "ok"
}
```