package api

import "crypto/tls"

// CreateHTTPS creates a new Server instance that listens on HTTPS
func CreateHTTPS(host string, c Config, gateway Gatewayer, certFile, keyFile string) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	logger.Infof("Using %s for the certificate", certFile)
	logger.Infof("Using %s for the key", keyFile)

	listener, err := tls.Listen("tcp", host, &tls.Config{ //nolint:gosec
		Certificates: []tls.Certificate{cert},
	})
	if err != nil {
		return nil, err
	}

	// If the host did not specify a port, allowing the kernel to assign one,
	// we need to get the assigned address to know the full hostname
	host = listener.Addr().String()

	s, err := create(host, c, gateway)
	if err != nil {
		if s != nil {
			if closeErr := s.listener.Close(); closeErr != nil {
				logger.WithError(err).Warning("s.listener.Close() error")
			}
		}
		return nil, err
	}

	s.listener = listener

	return s, nil
}
