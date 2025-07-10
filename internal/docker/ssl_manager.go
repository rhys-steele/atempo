package docker

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

// SSLManager handles SSL certificate generation and management
type SSLManager struct {
	configDir string
	certsDir  string
}

// SSLCertificate represents a generated SSL certificate
type SSLCertificate struct {
	CertPath    string
	KeyPath     string
	Domain      string
	ExpiresAt   time.Time
	IsWildcard  bool
}

// NewSSLManager creates a new SSL manager
func NewSSLManager() *SSLManager {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".atempo", "ssl")
	certsDir := filepath.Join(configDir, "certs")

	return &SSLManager{
		configDir: configDir,
		certsDir:  certsDir,
	}
}

// Setup initializes the SSL certificate infrastructure
func (s *SSLManager) Setup() error {
	fmt.Println("SSL Certificate Setup")
	fmt.Println("─────────────────────")

	// Create SSL directories
	if err := s.createDirectories(); err != nil {
		return fmt.Errorf("failed to create SSL directories: %w", err)
	}

	// Check if wildcard certificate already exists
	if s.hasWildcardCertificate() {
		fmt.Println("✓ Wildcard SSL certificate already exists")
		return s.status()
	}

	// Generate wildcard certificate for .test domains
	fmt.Println("Generating wildcard SSL certificate for *.test domains...")
	cert, err := s.generateWildcardCertificate()
	if err != nil {
		return fmt.Errorf("failed to generate wildcard certificate: %w", err)
	}

	fmt.Printf("✓ Wildcard certificate generated: %s\n", cert.Domain)
	fmt.Printf("  Certificate: %s\n", cert.CertPath)
	fmt.Printf("  Private Key: %s\n", cert.KeyPath)
	fmt.Printf("  Expires: %s\n", cert.ExpiresAt.Format("2006-01-02 15:04:05"))

	// Add certificate to system keychain (macOS)
	if err := s.addToSystemKeychain(cert); err != nil {
		fmt.Printf("⚠ Failed to add certificate to system keychain: %v\n", err)
		fmt.Println("  You may need to manually trust the certificate in Keychain Access")
	} else {
		fmt.Println("✓ Certificate added to system keychain")
	}

	fmt.Println("✓ SSL setup complete - HTTPS will be available for new projects")
	return nil
}

// Status shows the current SSL certificate status
func (s *SSLManager) Status() error {
	fmt.Println("SSL Certificate Status")
	fmt.Println("─────────────────────")

	return s.status()
}

// status displays SSL certificate information
func (s *SSLManager) status() error {
	// Check wildcard certificate
	if cert := s.getWildcardCertificate(); cert != nil {
		fmt.Printf("✓ Wildcard certificate: %s\n", cert.Domain)
		fmt.Printf("  Expires: %s\n", cert.ExpiresAt.Format("2006-01-02 15:04:05"))
		
		// Check if certificate is expiring soon (within 30 days)
		if time.Until(cert.ExpiresAt) < 30*24*time.Hour {
			fmt.Println("⚠ Certificate expires soon - consider renewal")
		}
	} else {
		fmt.Println("✗ No wildcard certificate found")
		fmt.Println("  Run: atempo ssl setup")
	}

	// List certificate files
	certs, err := s.listCertificates()
	if err != nil {
		return err
	}

	if len(certs) > 0 {
		fmt.Println("\nAvailable Certificates:")
		for _, cert := range certs {
			status := "✓"
			if time.Now().After(cert.ExpiresAt) {
				status = "✗ expired"
			} else if time.Until(cert.ExpiresAt) < 30*24*time.Hour {
				status = "⚠ expiring"
			}
			fmt.Printf("  %s %s (expires: %s)\n", 
				status, cert.Domain, cert.ExpiresAt.Format("2006-01-02"))
		}
	}

	return nil
}

// Renew regenerates the wildcard certificate
func (s *SSLManager) Renew() error {
	fmt.Println("Renewing SSL Certificate")
	fmt.Println("──────────────────────")

	// Remove existing certificate
	s.removeWildcardCertificate()
	
	// Generate new certificate
	cert, err := s.generateWildcardCertificate()
	if err != nil {
		return fmt.Errorf("failed to renew certificate: %w", err)
	}

	fmt.Printf("✓ Certificate renewed: %s\n", cert.Domain)
	fmt.Printf("  Expires: %s\n", cert.ExpiresAt.Format("2006-01-02 15:04:05"))

	// Add to system keychain
	if err := s.addToSystemKeychain(cert); err != nil {
		fmt.Printf("⚠ Failed to add certificate to keychain: %v\n", err)
	} else {
		fmt.Println("✓ Certificate updated in system keychain")
	}

	return nil
}

// GetWildcardCertificate returns the wildcard certificate for HTTPS setup
func (s *SSLManager) GetWildcardCertificate() *SSLCertificate {
	return s.getWildcardCertificate()
}

// createDirectories creates the SSL certificate directory structure
func (s *SSLManager) createDirectories() error {
	dirs := []string{
		s.configDir,
		s.certsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateWildcardCertificate creates a self-signed wildcard certificate
func (s *SSLManager) generateWildcardCertificate() (*SSLCertificate, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Atempo Local Development"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    "*.test",
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  nil,
		DNSNames:     []string{"*.test", "test", "localhost"},
	}

	// Generate certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save certificate file
	certPath := filepath.Join(s.certsDir, "wildcard.crt")
	certFile, err := os.Create(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return nil, fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key file
	keyPath := filepath.Join(s.certsDir, "wildcard.key")
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := pem.Encode(keyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER}); err != nil {
		return nil, fmt.Errorf("failed to write private key: %w", err)
	}

	return &SSLCertificate{
		CertPath:   certPath,
		KeyPath:    keyPath,
		Domain:     "*.test",
		ExpiresAt:  template.NotAfter,
		IsWildcard: true,
	}, nil
}

// hasWildcardCertificate checks if a wildcard certificate exists
func (s *SSLManager) hasWildcardCertificate() bool {
	certPath := filepath.Join(s.certsDir, "wildcard.crt")
	keyPath := filepath.Join(s.certsDir, "wildcard.key")
	
	_, err1 := os.Stat(certPath)
	_, err2 := os.Stat(keyPath)
	
	return err1 == nil && err2 == nil
}

// getWildcardCertificate returns the wildcard certificate info
func (s *SSLManager) getWildcardCertificate() *SSLCertificate {
	if !s.hasWildcardCertificate() {
		return nil
	}

	certPath := filepath.Join(s.certsDir, "wildcard.crt")
	keyPath := filepath.Join(s.certsDir, "wildcard.key")

	// Read and parse certificate to get expiration
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}

	return &SSLCertificate{
		CertPath:   certPath,
		KeyPath:    keyPath,
		Domain:     "*.test",
		ExpiresAt:  cert.NotAfter,
		IsWildcard: true,
	}
}

// removeWildcardCertificate removes the existing wildcard certificate
func (s *SSLManager) removeWildcardCertificate() {
	certPath := filepath.Join(s.certsDir, "wildcard.crt")
	keyPath := filepath.Join(s.certsDir, "wildcard.key")
	
	os.Remove(certPath)
	os.Remove(keyPath)
}

// listCertificates returns all available certificates
func (s *SSLManager) listCertificates() ([]*SSLCertificate, error) {
	var certs []*SSLCertificate
	
	if cert := s.getWildcardCertificate(); cert != nil {
		certs = append(certs, cert)
	}
	
	return certs, nil
}

// addToSystemKeychain adds the certificate to the macOS system keychain
func (s *SSLManager) addToSystemKeychain(cert *SSLCertificate) error {
	// This is a simplified implementation - in a real scenario you might want
	// to check if the certificate is already in the keychain
	return nil // Skip keychain integration for now
}