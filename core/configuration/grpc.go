package configuration

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"

	"github.com/quic-go/quic-go"
	log "github.com/sirupsen/logrus"
	"github.com/xans-me/grpcquic"
	"google.golang.org/grpc"

	// used for side effects
	_ "net/http/pprof"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer()
}

// ListenGRPC is a func to start http server
func ListenGRPC(conf *AppConfig) net.Listener {
	addr := conf.App.Host + ":" + conf.App.Port
	log.Info("run GRPC & HTTP Transport on " + addr)

	lis, err := net.Listen(conf.App.Protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	return lis
}

func ListenQGRPC(conf *AppConfig, certFile, keyFile string) net.Listener {
	tlsConf, err := generateTLSConfig(certFile, keyFile)
	if err != nil {
		log.Fatalf("QServer + gRPC => failed to generateTLSConfig. %s", err.Error())
	}

	addr := conf.App.Host + ":" + conf.App.Port

	qListener, err := quic.ListenAddr(addr, tlsConf, nil)
	if err != nil {
		log.Fatalf("QServer + gRPC => failed to ListenAddr. %s", err.Error())
	}
	listener := grpcquic.Listen(qListener)
	return listener
}

func generateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	if len(certFile) > 0 && len(keyFile) > 0 {
		log.Printf("generateTLSConfig => certFile=%s, keyFile=%s", certFile, keyFile)

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Printf("failed to tls.LoadX509KeyPair. %s", err.Error())
			return nil, err
		}
		return &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"q-echo-example"},
		}, nil
	} else {
		// log.Printf("generateTLSConfig => GenerateKey")
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			log.Printf("failed to rsa.GenerateKey. %s", err.Error())
			return nil, err
		}
		template := x509.Certificate{SerialNumber: big.NewInt(1)}
		certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
		if err != nil {
			log.Printf("failed to x509.CreateCertificate. %s", err.Error())
			return nil, err
		}
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

		tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			log.Printf("failed to tls.X509KeyPair. %s", err.Error())
			return nil, err
		}

		return &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			NextProtos:   []string{"q-echo-example"},
		}, nil
	}
}
