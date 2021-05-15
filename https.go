package main

import (
	"crypto/tls"
	"log"
)

func setupTLS() *tls.Config {
	cert, err := tls.X509KeyPair([]byte(certBytes), []byte(keyBytes))
	if err != nil {
		log.Println(err)
		return nil
	}

	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		Certificates:             []tls.Certificate{cert},
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	return config
}
