package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
	"golang.org/x/crypto/pbkdf2"
)

// EncryptAES_CBC encrypts the given plaintext using AES encryption in CBC mode.
// It returns the encrypted ciphertext as a base64-encoded string.
// The function takes a plaintext byte slice and a configuration object as input.
func EncryptAES_CBC(plaintext []byte, cfg *config.SecurityConfig) (string, error) {

	key := []byte(cfg.EncryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", eris.New("failed to create AES cipher")
	}

	paddedPlaintext := pkcs7Padding(plaintext, block.BlockSize())

	ciphertext := make([]byte, block.BlockSize()+len(paddedPlaintext))
	iv := ciphertext[:block.BlockSize()]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", eris.Wrap(err, "failed to read random bytes for IV")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[block.BlockSize():], paddedPlaintext)

	// Encode the ciphertext using base64 and URL-safe encoding
	encodedCiphertext := base64.URLEncoding.EncodeToString(ciphertext)

	return encodedCiphertext, nil
}

// DecryptAES_CBC decrypts the given ciphertext using AES decryption in CBC mode.
// It returns the decrypted plaintext as a string.
// The function takes a base64-encoded ciphertext string and a configuration object as input.
// It returns an error if the decryption fails.
// The ciphertext is expected to be in base64 URL-safe format.
func DecryptAES_CBC(ciphertext string, cfg *config.SecurityConfig) (string, error) {

	key := []byte(cfg.EncryptionKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", eris.Wrap(err, "failed to create AES cipher")
	}

	encrypted, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", eris.Wrap(err, "failed to decode base64 ciphertext")
	}

	if len(encrypted) < aes.BlockSize {
		return "", eris.New("ciphertext too short")
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	if len(encrypted)%aes.BlockSize != 0 {
		return "", eris.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)

	plaintext, err := pkcs7Unpadding(encrypted)
	if err != nil {
		return "", eris.Wrap(err, "failed to unpad decrypted data")
	}

	return string(plaintext), nil
}

func pkcs7Padding(input []byte, blockSize int) []byte {
	padding := blockSize - len(input)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padText...)
}

func pkcs7Unpadding(input []byte) ([]byte, error) {
	length := len(input)
	unpadding := int(input[length-1])
	if length < unpadding {
		return nil, eris.New("unpadding size is larger than input length")
	}
	return input[:(length - unpadding)], nil
}

// Password Utils
func GenerateRandomPassword(length int, cfg *config.SecurityConfig) string {
	charLength := big.NewInt(int64(len(cfg.AllowedCharacters)))
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charLength)
		password[i] = cfg.AllowedCharacters[randomIndex.Int64()]
	}

	return string(password)
}

// Password Storing Utils

// hexDecode decodes a hexadecimal string to bytes.
func hexDecode(hexStr string) ([]byte, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, eris.Wrap(err, "failed to decode hex string")
	}
	return data, nil
}

// HashPassword generates a PBKDF2 hash of the password and returns the hash and salt.
func HashPassword(password string, cfg *config.SecurityConfig) (string, string, error) {
	// Generate a random salt
	salt := make([]byte, cfg.SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", eris.Wrap(err, "failed to generate random salt")
	}

	// Compute the PBKDF2 hash of the password using the salt
	hashedPassword := pbkdf2.Key([]byte(password), salt, cfg.IterationCount, cfg.KeySize, sha256.New)

	// Encode the salt and hash as hexadecimal strings
	saltHex := fmt.Sprintf("%x", salt)
	hashHex := fmt.Sprintf("%x", hashedPassword)

	return hashHex, saltHex, nil
}

// VerifyPassword verifies if a given password matches a stored hash and salt.
func VerifyPassword(password, salt, storedHash string, cfg *config.SecurityConfig) bool {
	// Decode the salt and stored hash from hexadecimal strings
	saltBytes, _ := hexDecode(salt)
	storedHashBytes, _ := hexDecode(storedHash)

	// Compute the PBKDF2 hash of the input password using the stored salt
	computedHash := pbkdf2.Key([]byte(password), saltBytes, cfg.IterationCount, cfg.KeySize, sha256.New)

	// Use subtle.ConstantTimeCompare to compare the computed hash with the stored hash
	return subtle.ConstantTimeCompare(computedHash, storedHashBytes) == 1
}
