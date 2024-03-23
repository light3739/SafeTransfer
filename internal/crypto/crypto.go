package crypto

import (
	"bufio"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// newAesCtrStream initializes an AES cipher in CTR mode with the given key and IV.
func newAesCtrStream(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}
	return cipher.NewCTR(block, iv), nil
}

func EncryptFile(file io.Reader, key []byte) (io.Reader, []byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	stream, err := newAesCtrStream(key, iv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize AES CTR stream: %w", err)
	}

	bufferedFile := bufio.NewReader(file)
	encryptedFile := &cipher.StreamReader{S: stream, R: bufferedFile}
	return encryptedFile, iv, nil
}

func DecryptFile(encryptedFile io.Reader, key []byte, iv []byte) (io.Reader, error) {
	stream, err := newAesCtrStream(key, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AES CTR stream: %w", err)
	}

	bufferedFile := bufio.NewReader(encryptedFile)
	decryptedFile := &cipher.StreamReader{S: stream, R: bufferedFile}
	return decryptedFile, nil
}

// SignFile signs the given file using RSA.
func SignFile(file io.Reader, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to hash file: %w", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to sign hash: %w", err)
	}

	return signature, nil
}

// VerifyFile verifies the given file's signature using RSA.
func VerifyFile(file io.Reader, signature string, publicKey *rsa.PublicKey) error {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to hash file for verification: %w", err)
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), signatureBytes); err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	return nil
}
