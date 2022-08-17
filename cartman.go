package cartman

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"
)

var Debug bool

type client struct {
	Name string
	Cert *x509.Certificate
}

type Store struct {
	root string
	certs map[string]client
}

func debug(msg ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stderr, "[cartman] DEBUG:", msg)
	}
}

func NewStore(root string) (Store, error) {
	s := Store{root: root, certs: map[string]client{}}
	dir, err := os.Open(root)
	if err != nil {
		return s, err
	}
	clients, err := dir.Readdirnames(0)
	if err != nil {
		return s, err
	}

	for _, name := range clients {
		debug("found client", name)

		// load certs
		der, err := os.ReadFile(path.Join(root, name, "cert.pem"))
		if err != nil {
			return s, err
		}
		block, _ := pem.Decode(der)
		if err != nil {
			return s, err
		}

		u := client{Name: name}
		u.Cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return s, err
		}
		s.certs[fingerprint(u.Cert)] = u
	}
	debug("loaded", len(s.certs), "certificates")
	
	return s, nil
}

func (s Store) GetClientFromCert(cert *x509.Certificate) (string, error) {
	client, ok := s.certs[fingerprint(cert)]
	if !ok {
		return "", fmt.Errorf("client not found for certificate")
	}
	return client.Name, nil
}

func (s Store) AddClient(name string, cert *x509.Certificate) error {
	debug("adding certificate for", name)

	f, err := os.Create(path.Join(s.root, name, "cert.pem"))
	if err != nil {
		return err
	}
	defer f.Close()

	block := pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	if err := pem.Encode(f, &block); err != nil {
		return err
	}

	s.certs[fingerprint(cert)] = client{Name: name, Cert: cert}
	return nil
}

func fingerprint(cert *x509.Certificate) string {
	fp := sha1.Sum(cert.Raw)
	return string(fp[:])
}

