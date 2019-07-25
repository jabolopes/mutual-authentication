package main

import (
  "context"
  "crypto/tls"
  "crypto/x509"
  "errors"
  "flag"
  "fmt"
  "io/ioutil"
  "net"

  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/peer"
  "google.golang.org/grpc/status"

  "github.com/jabolopes/grpccert/grpctls"
  echopb "github.com/jabolopes/grpccert/proto"
)

var (
  caCert   = flag.String("ca_cert", "", "Root certificate authority")
  selfCert = flag.String("self_cert", "", "Certificate for this server")
  selfKey  = flag.String("self_key", "", "Private key for this server")
)

func NewServerTLSFromFileCustomCA(caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
  // Load certificates from disk.
  certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
  if err != nil {
    return nil, fmt.Errorf("failed to load X509 key pair: %v", err)
  }

  // Create a certificate pool from the certificate authority.
  certPool := x509.NewCertPool()
  ca, err := ioutil.ReadFile(caFile)
  if err != nil {
    return nil, fmt.Errorf("failed to read CA certificate: %v", err)
  }

  // Append the certificates from the CA.
  if ok := certPool.AppendCertsFromPEM(ca); !ok {
    return nil, errors.New("failed to append certs")
  }

  return credentials.NewTLS(&tls.Config{
    ClientAuth:   tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{certificate},
    ClientCAs:    certPool,
  }), nil
}

func getPeerName(ctx context.Context) (string, error) {
  p, ok := peer.FromContext(ctx)
  if !ok {
    return "", status.Error(codes.Unauthenticated, "no peer found")
  }

  tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
  if !ok {
    return "", status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
  }

  if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
    return "", status.Error(codes.Unauthenticated, "failed to verify peer certificate")
  }

  return tlsAuth.State.VerifiedChains[0][0].Subject.CommonName, nil
}

type server struct {
}

func (s *server) Echo(ctx context.Context, request *echopb.EchoRequest) (*echopb.EchoResponse, error) {
  peerName, err := getPeerName(ctx)
  if err != nil {
    return nil, err
  }

  return &echopb.EchoResponse{Data: fmt.Sprintf("Hello %s: you said %s", peerName, request.Data)}, nil
}

func run(caCert, selfCert, selfKey string) error {
  address, err := grpctls.GetCommonName(selfCert)
  if err != nil {
    return err
  }

  lis, err := net.Listen("tcp", address)
  if err != nil {
    return err
  }

  creds, err := NewServerTLSFromFileCustomCA(caCert, selfCert, selfKey)
  if err != nil {
    return err
  }

  fmt.Println("Ctrl-C to terminate the process")

  grpcServer := grpc.NewServer(grpc.Creds(creds))
  echopb.RegisterEchoServiceServer(grpcServer, &server{})
  if err := grpcServer.Serve(lis); err != nil {
    return err
  }

  return nil
}

func main() {
  flag.Parse()

  if err := run(*caCert, *selfCert, *selfKey); err != nil {
    panic(err)
  }
}
