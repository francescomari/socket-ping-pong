# Socket Ping Pong

This project shows how to implement an HTTP [server](cmd/server/main.go) and [client](cmd/client/main.go) on Unix sockets.

## Run

Start the server:

```
go run ./cmd/server -socket socket.sock
```

Start the client:

```
go run ./cmd/client -socket socket.sock
```

If everything works, the client should emit messages similar to these:

```
2021/04/13 23:09:35 received 'pong'
2021/04/13 23:09:36 received 'pong'
2021/04/13 23:09:37 received 'pong'
2021/04/13 23:09:38 received 'pong'
2021/04/13 23:09:39 received 'pong'
```
