package grpctls

import (
  "crypto/tls"
  "crypto/x509"
  "encoding/pem"
  "errors"
  "fmt"
  "io/ioutil"

  "google.golang.org/grpc/credentials"
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

func NewClientTLSFromFileCustomCA(serverName, caFile, certFile, keyFile string) (credentials.TransportCredentials, error) {
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
    ServerName:   serverName,
    Certificates: []tls.Certificate{certificate},
    RootCAs:      certPool,
  }), nil
}

func GetCommonName(certFile string) (string, error) {
  certData, err := ioutil.ReadFile(certFile)
  if err != nil {
    return "", err
  }

  block, _ := pem.Decode(certData)
  if block == nil || block.Type != "CERTIFICATE" {
    return "", errors.New("failed to decode PEM certificate")
  }

  cert, err := x509.ParseCertificate(block.Bytes)
  if err != nil {
    return "", err
  }

  return cert.Subject.CommonName, nil
}
