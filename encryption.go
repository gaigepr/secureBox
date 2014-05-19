package main

import (
	"code.google.com/p/go.crypto/pbkdf2"
	//"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/sha256"
	"fmt"
	"os"
)

func ReadAndEncrypt(filename string) {
	// read a file in chunks, encrypt, send!
	// Errors here should maybe result in pushing those event back onto the queue?
	// or push them onto an error channel so that they can be handled else where?

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err, "\n\n")
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	fmt.Println(stat.Size())

	var amount int64 = 0
	var EOF bool = false
	for !EOF {
		if amount >= stat.Size() {
			EOF = true
		}

		data := make([]byte, 16)
		count, err := file.Read(data)

		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		amount += int64(count)
		fmt.Printf("read %d bytes: %q\n", count, data[:count])
	}

}

// test function to see if the symmetric RSA encryption worked.
func TestCrypto() {
	key := createRSA(2048)
	message := []byte("this is a test message")
	fmt.Println(len(message))
	encMessage := encryptRSA(&key.PublicKey, message)
	decMessage := decryptRSA(key, encMessage)
	fmt.Println("message is: {", string(message[:]), "}\ndecrypted message is: {", string(decMessage[:]), "}")
	if len(message) == len(decMessage) {
		fmt.Printf("It worked!\n")
	}
	return
}

// creates a unique(relying on crypto/rand) AES key for the per file keys
func createAES() []byte {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return key
}

// creates the master AES key for the user based on their password.
func createUserAES(password string) ([]byte, []byte) {
	salt := make([]byte, 10) // TODO: make 10 a constant variable
	_, err := rand.Read(salt)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	key := pbkdf2.Key([]byte(password), salt, 5000, 32, sha256.New) // TODO: make 5000 a const var

	return key, salt
}

// create an RSA key when given a size.
// TODO: add in customization based on password.
func createRSA(size int) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Printf("error occured generating the key:\n")
	}
	return key
}

// encrypts a given message with a provided RSA public key.
// the label must be the same for the encryption and decryption for it to work.
func encryptRSA(pub *rsa.PublicKey, message []byte) []byte {
	label := make([]byte, 10)
	encMessage, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pub, message, label)
	fmt.Println(encMessage)
	if err != nil {
		fmt.Println("encryption failed!: ", err)
	}
	return encMessage
}

// decrypts a given message with the provided RSA private key
// the label must be the same for the encryption and decryption for it to work.
func decryptRSA(private *rsa.PrivateKey, ciphertext []byte) []byte {
	label := make([]byte, 10)
	message, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, private, ciphertext, label)
	if err != nil {
		fmt.Println("decryption failed!: ", err)
	}
	return message
}
