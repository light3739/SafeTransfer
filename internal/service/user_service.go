package service

import (
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	UserRepo     repository.UserRepository
	JWTSecretKey string
}

// NewUserService function now accepts a JWTSecretKey as an argument
func NewUserService(userRepo repository.UserRepository, JWTSecretKey string) *UserService {
	return &UserService{
		UserRepo:     userRepo,
		JWTSecretKey: JWTSecretKey,
	}
}

func (us *UserService) GenerateNonceForUser(ethereumAddress string) (string, error) {
	nonceBytes := make([]byte, 16) // 128-bit nonce
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	nonce := hex.EncodeToString(nonceBytes)

	user := &model.User{
		EthereumAddress: ethereumAddress,
		Nonce:           nonce,
	}
	err := us.UserRepo.SaveOrUpdateUser(user)
	return nonce, err
}

func (us *UserService) GetNonceForUser(ethereumAddress string) (string, error) {
	user, err := us.UserRepo.FindByEthereumAddress(ethereumAddress)
	if err != nil {
		return "", err
	}
	return user.Nonce, nil
}

// VerifySignature verifies the signature of the nonce signed by the user.
func (us *UserService) VerifySignature(message, signature string) (string, error) {
	log.Printf("Verifying signature for message: %s", message)
	log.Printf("Signature: %s", signature)

	// Prepare the message for Ethereum signature verification
	message = "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message
	messageHash := crypto.Keccak256Hash([]byte(message))

	// Remove the '0x' prefix from the signature if present
	signature = strings.TrimPrefix(signature, "0x")

	// Decode the signature from hex
	sigBytes, err := hexutil.Decode("0x" + signature)
	if err != nil {
		return "", err
	}

	// Adjust the S/V values
	if sigBytes[64] != 27 && sigBytes[64] != 28 {
		return "", errors.New("invalid Ethereum signature (V value is incorrect)")
	}
	sigBytes[64] -= 27

	// Recover the public key from the signature
	publicKeyECDSA, err := crypto.SigToPub(messageHash.Bytes(), sigBytes)
	if err != nil {
		return "", err
	}

	// Convert the public key to an Ethereum address
	recoveredAddr := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	log.Printf("Recovered Ethereum address: %s", recoveredAddr)

	return recoveredAddr, nil
}

// GenerateJWT generates a JWT for a given user.
func (us *UserService) GenerateJWT(ethereumAddress string) (string, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return "", errors.New("JWT secret key is not set")
	}

	claims := jwt.MapClaims{
		"ethereumAddress": ethereumAddress,
		"exp":             time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
