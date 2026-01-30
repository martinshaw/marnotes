package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"path/filepath"
)

type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func GenerateKeyPair(bits int) (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

func LoadOrGenerateKeyPair(keyDir string, bits int) (*KeyPair, error) {
	privKeyPath := filepath.Join(keyDir, "private.pem")
	pubKeyPath := filepath.Join(keyDir, "public.pem")

	// Try to load existing keys
	if _, err := os.Stat(privKeyPath); err == nil {
		privateKey, err := LoadPrivateKey(privKeyPath)
		if err == nil && privateKey != nil {
			return &KeyPair{
				PrivateKey: privateKey,
				PublicKey:  &privateKey.PublicKey,
			}, nil
		}
	}

	// Generate new keys
	kp, err := GenerateKeyPair(bits)
	if err != nil {
		return nil, err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(keyDir, 0755); err != nil {
		return nil, err
	}

	// Save keys
	if err := SavePrivateKey(privKeyPath, kp.PrivateKey); err != nil {
		return nil, err
	}
	if err := SavePublicKey(pubKeyPath, kp.PublicKey); err != nil {
		return nil, err
	}

	return kp, nil
}

func SavePrivateKey(filePath string, privateKey *rsa.PrivateKey) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return os.WriteFile(filePath, privateKeyPEM, 0600)
}

func SavePublicKey(filePath string, publicKey *rsa.PublicKey) error {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return os.WriteFile(filePath, publicKeyPEM, 0644)
}

func LoadPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey, err
}

func LoadPublicKey(filePath string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func Encrypt(publicKey *rsa.PublicKey, data []byte) (string, error) {
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(privateKey *rsa.PrivateKey, encryptedData string) ([]byte, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, nil)
}

func GetPublicKeyPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}), nil
}
