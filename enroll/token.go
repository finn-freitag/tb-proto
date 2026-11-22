// Package enroll contains helpers shared by the dashboard and the agent for
// node enrollment tokens and dashboard CA pinning.
//
// Token format:  tb1.<sha256 of dashboard CA certificate DER, hex>.<secret>
//
// The agent uses the fingerprint to verify the dashboard's TLS certificate
// chain on first contact (before it has the CA), then calls Enrollment.Enroll
// with the secret.
package enroll

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const tokenPrefix = "tb1"

// Token is a parsed enrollment token.
type Token struct {
	CAFingerprint string // lowercase hex sha256 of the CA certificate DER
	Secret        string
}

// String formats the token.
func (t Token) String() string {
	return fmt.Sprintf("%s.%s.%s", tokenPrefix, t.CAFingerprint, t.Secret)
}

// ParseToken parses an enrollment token.
func ParseToken(s string) (Token, error) {
	parts := strings.Split(strings.TrimSpace(s), ".")
	if len(parts) != 3 || parts[0] != tokenPrefix {
		return Token{}, errors.New("enroll: malformed token")
	}
	if len(parts[1]) != 64 {
		return Token{}, errors.New("enroll: malformed CA fingerprint")
	}
	if _, err := hex.DecodeString(parts[1]); err != nil {
		return Token{}, errors.New("enroll: malformed CA fingerprint")
	}
	if parts[2] == "" {
		return Token{}, errors.New("enroll: empty secret")
	}
	return Token{CAFingerprint: strings.ToLower(parts[1]), Secret: parts[2]}, nil
}

// Fingerprint returns the lowercase hex sha256 of a DER certificate.
func Fingerprint(der []byte) string {
	sum := sha256.Sum256(der)
	return hex.EncodeToString(sum[:])
}

// VerifyPinnedChain returns a tls.Config.VerifyPeerCertificate-compatible
// function that accepts a server chain only if one of its certificates is the
// pinned CA and the leaf verifies against it.
func VerifyPinnedChain(fingerprint string) func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return errors.New("enroll: no server certificate")
		}
		var ca *x509.Certificate
		for _, raw := range rawCerts {
			if Fingerprint(raw) == fingerprint {
				c, err := x509.ParseCertificate(raw)
				if err != nil {
					return err
				}
				ca = c
			}
		}
		if ca == nil {
			return errors.New("enroll: dashboard CA fingerprint mismatch")
		}
		leaf, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		if leaf.Equal(ca) {
			return nil
		}
		pool := x509.NewCertPool()
		pool.AddCert(ca)
		_, err = leaf.Verify(x509.VerifyOptions{Roots: pool, KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}})
		return err
	}
}
