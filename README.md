# gRPC cert

This project shows how to perform mutual TLS authentication with
certificates (using gRPC). Normally, the client authenticates the
server by verifying the server's certificate (this is the typical
scenario in HTTPS). But this example shows also how the server can
also authenticate the client. This is an authentication mechanism that
does not require username / password, only certificates.

Finally, both the client and server extract the common name from the
certificate. The client uses the common name as the server address to
connect to. The server uses the common name as the port to listen
to. This also means that the client and the server never need to be
given an additional address to connect or bind to (respectively), they
only need the certificates.

This project combines the following sources:

* https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
* https://jbrandhorst.com/post/grpc-auth/

## Dependencies

Install the compiler for the Protocol Buffer version 3 (aka
proto3). It should be available in your distribution's packages.

Install the compiler for the Go programming language. It should also
be available in your distribution's packages.

Slackware packages:

* proto3: http://slackbuilds.org/repository/14.2/misc/protobuf3
* Go: http://slackbuilds.org/repository/14.2/development/google-go-lang

## Run the server and client

Run the server with automatically generated certificates:

```shell
make run-server
```

In a different terminal, run the client with automatically generated
certificates:

```shell
make run-client
```
