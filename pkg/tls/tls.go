package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

var (
	LoadX509KeyPairFailedErr = errors.New("LoadX509KeyPair failed")
)

type ConfigTLS struct {
	certPath string
	keyPath  string
}

func NewConfigTLS(certPath string, keyPath string) *ConfigTLS {
	return &ConfigTLS{
		certPath: certPath,
		keyPath:  keyPath,
	}
}

func (cfg *ConfigTLS) BuildTLSConfig(cipherSuite []uint16) (*tls.Config, error) {
	var tlsConfig tls.Config

	if cfg.keyPath != "" {
		cert, err := tls.LoadX509KeyPair(cfg.certPath, cfg.keyPath)
		if err != nil {
			log.WithField("error", err).Error("can't LoadX509KeyPair")
			return nil, LoadX509KeyPairFailedErr
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	} else {
		log.Error(errors.New("Not loading cert key pair because key path not provided. Incoming connections will not work!"))
	}

	// The internal-wildcard.crt is not trusted by default so we have to add it to the CA store
	// Annoyingly, the CA stuff requires it to be loaded in a different way to cert list above
	x509pem, err := ioutil.ReadFile(cfg.certPath)
	if err != nil {
		return nil, fmt.Errorf("GetTLSConfig: read: %v", err)
	}

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("GetTLSConfig: certpool: %v", err)
	}

	ok := certPool.AppendCertsFromPEM(x509pem)
	if !ok {
		return nil, errors.New("GetTLSConfig: Failed to add cert to CA store")
	}

	tlsConfig.RootCAs = certPool
	tlsConfig.ClientCAs = certPool

	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.CipherSuites = cipherSuite
	tlsConfig.InsecureSkipVerify = true

	return &tlsConfig, nil
}

// ScalaServiceCipherSuites defines the Ecom cipher suite used in Scala services
func ScalaServiceCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}
}

// ActiveMqCipherSuites returns the cipher suites to be used for outbound connections to ActiveMQ
func ActiveMqCipherSuites() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}
}
