package cartman

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

var Debug bool

type client struct {
	Name string
	Fingerprint string
}

type FileStore struct {
	root string
	certs map[string]client
}

func debug(msg ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stderr, "[cartman] DEBUG:", msg)
	}
}

func NewFileStore(root string) (FileStore, error) {
	s := FileStore{root: root, certs: map[string]client{}}
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
		content, err := os.ReadFile(path.Join(root, name))
		if err != nil {
			return s, err
		}
		fp := string(content)

		s.certs[fp] = client{Name: name, Fingerprint: fp}
		debug("loaded", name, "with fingerprint", fp)
	}
	debug("loaded", len(s.certs), "clients")
	
	return s, nil
}

func (s FileStore) GetClientFromCert(cert *x509.Certificate) (string, error) {
	client, ok := s.certs[fingerprint(cert)]
	if !ok {
		return "", fmt.Errorf("client not found for certificate")
	}
	return client.Name, nil
}

func (s FileStore) AddClient(name string, cert *x509.Certificate) error {
	debug("adding certificate for", name)

	f, err := os.Create(path.Join(s.root, name))
	if err != nil {
		return err
	}
	defer f.Close()

	fp := fingerprint(cert)
	if _, err := f.WriteString(fp); err != nil {
		return err
	}

	s.certs[fingerprint(cert)] = client{Name: name, Fingerprint: fp}
	return nil
}

func fingerprint(cert *x509.Certificate) string {
	fp := sha1.Sum(cert.Raw)
	return hex.EncodeToString(fp[:])
}

