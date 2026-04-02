package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// обработка флагов командной строки
	// -enc для кодирования, -dec - декодирования
	// -codeword кодовое слово
	// последний параметр - шифруемая фраза
	isEncode := flag.Bool("enc", false, "Encode message with a codeword specified with --codeword")
	isDecode := flag.Bool("dec", false, "Decode message with a codeword specified with --codeword")
	codeword := flag.String("codeword", "", "A codeword to use with a phrase")
	isFile := flag.String("file", "", "A path to the file, if not specified encodes/decodes a string")

	flag.Parse()

	if len(os.Args) > 1 {
		var phrase string
		if *isFile == "" {
			phrase = os.Args[len(os.Args)-1]
		} else {
			phraseByte, err := os.ReadFile(*isFile)
			if err != nil {
				fmt.Println("Open file error: ", err)
				return
			}
			phrase = string(phraseByte)
		}

		if *codeword != "" {
			if !(*isEncode && *isDecode) {
				if *isEncode {
					// Encrypt
					ciphertext, err := encryptAES([]byte(phrase), []byte(*codeword))
					if err != nil {
						fmt.Println("Encryption error:", err)
						return
					}
					fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))
				}
				if *isDecode {
					// Decrypt
					decrypted, err := decryptAES([]byte(phrase), []byte(*codeword))
					if err != nil {
						fmt.Println("Decryption error:", err)
						return
					}
					fmt.Println(decrypted)
				}
			}
		}
	}
}

func encryptAES(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad plaintext to block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	ciphertext := make([]byte, aes.BlockSize+len(padtext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padtext)

	return ciphertext, nil
}

func decryptAES(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad
	padding := int(ciphertext[len(ciphertext)-1])
	return ciphertext[:len(ciphertext)-padding], nil
}
