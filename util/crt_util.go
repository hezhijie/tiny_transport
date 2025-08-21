package util

import (
	"crypto/tls"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
)

var (
	errBlockIsNotCertificate = errors.New("block is not a certificate, unable to load certificates")
	errNoCertificateFound    = errors.New("no certificate found, unable to load certificates")
)

func LoadKeyAndCertificate(keyPath string, certificatePath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certificatePath, keyPath)
}

func Check(err error) {
	var netError net.Error
	if errors.As(err, &netError) && netError.Temporary() { //nolint:staticcheck
		fmt.Printf("Warning: %v\n", err)
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}
}

func LoadCertificate(path string) (*tls.Certificate, error) {
	rawData, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var certificate tls.Certificate

	for {
		block, rest := pem.Decode(rawData)
		if block == nil {
			break
		}

		if block.Type != "CERTIFICATE" {
			return nil, errBlockIsNotCertificate
		}

		certificate.Certificate = append(certificate.Certificate, block.Bytes)
		rawData = rest
	}

	if len(certificate.Certificate) == 0 {
		return nil, errNoCertificateFound
	}

	return &certificate, nil
}
