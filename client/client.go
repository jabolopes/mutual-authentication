package main

import (
  "context"
  "flag"
  "fmt"

  "google.golang.org/grpc"

  "github.com/jabolopes/grpccert/grpctls"
  echopb "github.com/jabolopes/grpccert/proto"
)

var (
  caCert     = flag.String("ca_cert", "", "Root certificate authority")
  selfCert   = flag.String("self_cert", "", "Certificate for this client")
  selfKey    = flag.String("self_key", "", "Private key for this client")
  serverCert = flag.String("server_cert", "", "Remote server certificate")
)

func run(ctx context.Context, caCert, selfCert, selfKey, serverCert string) error {
  serverName, err := grpctls.GetCommonName(serverCert)
  if err != nil {
    return err
  }

  creds, err := grpctls.NewClientTLSFromFileCustomCA(serverName, caCert, selfCert, selfKey)
  if err != nil {
    return err
  }

  conn, err := grpc.Dial(serverName, grpc.WithTransportCredentials(creds))
  if err != nil {
    return err
  }

  client := echopb.NewEchoServiceClient(conn)
  response, err := client.Echo(ctx, &echopb.EchoRequest{Data: "Hello world!"})
  if err != nil {
    return err
  }

  fmt.Printf("Response: %s\n", response.Data)

  return nil
}

func main() {
  flag.Parse()

  if err := run(context.Background(), *caCert, *selfCert, *selfKey, *serverCert); err != nil {
    panic(err)
  }
}
