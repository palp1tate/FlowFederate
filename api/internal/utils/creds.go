package utils

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/palp1tate/FlowFederate/api/global"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

func GetClientCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("authorization/client.crt", "authorization/client.key")
	if err != nil {
		zap.S().Error("LoadX509KeyPair error ", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("internal/authorization/ca.crt")
	if err != nil {
		zap.S().Error("ReadFile ca.crt error ", err)
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		zap.S().Error("certPool.AppendCertsFromPEM error")
		return nil, err
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   global.ServerConfig.CommonName,
		RootCAs:      certPool,
	})
	return c, nil
}

func GetServerCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("authorization/server.crt", "authorization/server.key")
	if err != nil {
		zap.S().Error("LoadX509KeyPair error ", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("internal/authorization/ca.crt")
	if err != nil {
		zap.S().Error("ReadFile ca.crt error ", err)
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		zap.S().Error("certPool.AppendCertsFromPEM error")
		return nil, err
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})
	return c, nil
}
