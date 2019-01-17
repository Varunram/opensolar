package aes

// the aes package implements AES-256 GCM encryption and decrpytion functions
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	utils "github.com/YaleOpenLab/opensolar/utils"
)

// Encrypt encrypts a given data stream with a given passphrase
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(utils.SHA3hash(passphrase)[96:128]) // last 32 characters in hash
	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Failed to initialize a new AES GCM while encrypting", err)
		return data, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println("Error while reading gcm bytes", err)
		return data, err
	}
	log.Println("RANDOM ENCRYPTION NONCE IS: ", nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts a given data stream with a given passphrase
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 || len(passphrase) == 0 {
		return data, fmt.Errorf("Length of data is zero, can't decrpyt!")
	}
	tempParam := utils.SHA3hash(passphrase)
	key := []byte(tempParam[96:128]) // last 32 characters in hash
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error while initializing cipher decryption")
		return data, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Failed to initialize a new AES GCM while decrypting")
		return data, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Failed to open gcm while decrypting", err)
		return plaintext, err
	}
	return plaintext, nil
}

// EncryptFile encrypts a given file with the given passphrase
func EncryptFile(filename string, data []byte, passphrase string) error {
	f, _ := os.Create(filename)
	defer f.Close()
	data, err := Encrypt(data, passphrase)
	if err != nil {
		return err
	}
	f.Write(data)
	return nil
}

// DecryptFile encrypts a given file with the given passphrase
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}
