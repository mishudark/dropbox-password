package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//Hash return the encrypted hash of a password
func Hash(password, masterKey string) (string, error) {
	plaintext := []byte(password)

	//1) First, the plaintext password is transformed into a hash value using SHA512
	hash512 := sha512.Sum512(plaintext)

	//2) SHA512 hash is hashed again using bcrypt with a cost of 10, and a unique, per-user salt
	bcryptHash, err := bcrypt.GenerateFromPassword(hash512[:], 10)
	if err != nil {
		return "", err
	}

	//3) Finally, the resulting bcrypt hash is encrypted with AES256 using a secret key
	key := []byte(masterKey)
	block, err := aes.NewCipher(key)

	if err != nil {
		return "", err
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, bcryptHash, nil)
	strHash := hex.EncodeToString(ciphertext)
	b64Nonce := base64.StdEncoding.EncodeToString(nonce)

	return "aes256$" + b64Nonce + "$" + strHash, nil
}

//IsValid checks if a password match this hash
func IsValid(password, hash, masterKey string) bool {
	parts := strings.Split(hash, "$")

	if len(parts) < 3 {
		return false
	}

	nonce, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	ciphertext, err := hex.DecodeString(parts[2])
	if err != nil {
		return false
	}

	block, err := aes.NewCipher([]byte(masterKey))
	if err != nil {
		return false
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false
	}

	hash512 := sha512.Sum512([]byte(password))
	if err = bcrypt.CompareHashAndPassword(plaintext, hash512[:]); err != nil {
		return false
	}

	return true
}
