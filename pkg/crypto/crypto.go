// pkg/crypto/crypto.go

package crypto

import (
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

// EncryptFile encrypts the given file using AES.
func EncryptFile(file io.Reader, key []byte) (io.Reader, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Generate a random IV for CTR mode
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	// Use the IV with cipher.NewCTR
	stream := cipher.NewCTR(block, iv)

	// Create a reader that encrypts the input file
	encryptedFile := &cipher.StreamReader{S: stream, R: file}

	// Log the IV to ensure it's unique for each encryption
	fmt.Printf("IV: %x\n", iv)

	return encryptedFile, iv, nil
}

// DecryptFile decrypts the given file using AES.
func DecryptFile(encryptedFile io.Reader, key []byte, nonce []byte) (io.Reader, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Use the nonce as the IV for CTR mode
	stream := cipher.NewCTR(block, nonce)

	// Create a reader that decrypts the input file
	decryptedFile := &cipher.StreamReader{S: stream, R: encryptedFile}

	return decryptedFile, nil
}

// SignFile signs the given file using RSA.
func SignFile(file io.Reader, privateKey *rsa.PrivateKey) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyFile verifies the given file's signature using RSA.
func VerifyFile(file io.Reader, signature string, publicKey *rsa.PublicKey) error {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), signatureBytes)
}
