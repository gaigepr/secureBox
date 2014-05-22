package main

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	//	"path/filepath"
)

// TODO:  ADD MAKING PADDING AND STRIPPING IT AS WELL AS PUTTING THE BASE IV ONTO THE FILE
func ReadAndEncrypt(filename string) {
	// read a file in chunks, encrypt, send!
	// Errors here should maybe result in pushing those event back onto the queue?
	// or push them onto an error channel so that they can be handled else where?

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
func CreateAES() []byte {
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

func EncryptFile(filepath string, key []byte) {
	plainText, err := os.Open(filepath)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	defer plainText.Close()

	cipherText, err := os.Create(filepath + ".aes")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	defer cipherText.Close()

	iv := make([]byte, 16) // AES const
	if _, err = rand.Read(iv); err != nil {
		fmt.Println("ERROR:", err)
	}
	cipherText.Write(iv)
	fmt.Println("iv=", iv[:])

	encBlock := make([]byte, 16) // change to AES const for non-hardcoding
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	// we now have the encrypter, we can loop over the file contents now
	stats, err := plainText.Stat()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	//var fileRead int64 = 0
	pad := 16 - (stats.Size() % 16) //AES const
	if pad == 0 {
		pad = 16
	}
	message := make([]byte, 16) // AES const
	endFile := false

	for !endFile { //stats.Size() >= fileRead  { // AES const
		//read in file, check for end, if so pad
		_, err := plainText.Read(message)
		//fileRead += 16 // AES const
		if err != nil && err != io.EOF {
			fmt.Println("ERROR:", err) // if delete close down encryption
		}
		if err == io.EOF {
			endFile = true
			// we have to pad
			for i := pad; i > 0; i-- {
				message[16-i] = byte(pad)
			}
		}
		fmt.Println("message =", message[:])
		// encrypt the message slice
		mode.CryptBlocks(encBlock, message)
		fmt.Println("encrypt=", encBlock[:])
		//write and reset the message slice
		cipherText.Write(encBlock)
	}
}

func DecryptFile(filepath string, key []byte) {

	fmt.Println()

	cipherText, err := os.Open(filepath)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	defer cipherText.Close()

	plainText, err := os.Create(filepath + ".dec")
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	defer plainText.Close()

	iv := make([]byte, 16) //AES const

	_, err = cipherText.Read(iv)
	fmt.Println("iv=", iv[:])
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	encBlock := make([]byte, 16) // AES const
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	// we have decrypter, can loop over file contents now

	// know size so can not get EOF'd
	stats, err := cipherText.Stat()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	var ReadAmt int64 = 16
	message := make([]byte, 16)

	for ReadAmt != stats.Size() {
		//read in file, decrypt, check for end, if so strip padding
		if _, err := cipherText.Read(encBlock); err != nil {
			fmt.Println("ERROR:", err)
		}

		ReadAmt += 16 // AES const
		mode.CryptBlocks(message, encBlock)

		//check for end, if so strip pad
		fmt.Println("encrypt=", encBlock[:])
		fmt.Println("message=", message[:])
		if ReadAmt == stats.Size() {
			pad := message[15]
			message = message[:16-pad]
		}
		plainText.Write(message)
	}
}

func encryptAES(key []byte, iv []byte, message []byte) []byte {
	encBlock := make([]byte, 16)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	//loop now
	mode.CryptBlocks(encBlock, message)
	return encBlock
}

func decryptAES(key []byte, iv []byte, encBlock []byte) []byte {
	decBlock := make([]byte, 16)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decBlock, encBlock)
	return decBlock
}
