package main

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func Xor(m1, m2 []byte) []byte {
	result := make([]byte, len(m1))
	for i := 0; i < len(m1); i++ {
		result[i] = m1[i] ^ m2[i]
	}
	return result
}

func hmac(mac_key, message []byte) [32]byte {
	blocksize := 64
	key := make([]byte, blocksize)
	if len(mac_key) > blocksize {
		temp := sha256.Sum256(mac_key)
		copy(key, temp[:])
	}
	if len(mac_key) < blocksize {
		key = mac_key
		for i := 0; i < (blocksize - len(mac_key)); i++ {
			key = append(key, 0X00)
		}
	}

	o_pad := make([]byte, blocksize)
	for i := 0; i < blocksize; i++ {
		o_pad[i] = 0x5c
	}

	i_pad := make([]byte, blocksize)
	for i := 0; i < blocksize; i++ {
		i_pad[i] = 0x36
	}

	o_key_pad := Xor(o_pad, key)
	i_key_pad := Xor(i_pad, key)

	temp2 := sha256.Sum256(append(i_key_pad, message...))
	ihash := make([]byte, len(temp2))
	copy(ihash, temp2[:])
	ohash := sha256.Sum256(append(o_key_pad, ihash...))

	return ohash

}

func Enc(plaintext []byte, iv []byte, Enc_key []byte, Mac_key []byte) []byte {

	cipher_block, err := aes.NewCipher(Enc_key)
	if err != nil {
		fmt.Println("Key Error\n")
	}

	temp := hmac(Mac_key, plaintext)
	mac := make([]byte, len(temp))
	copy(mac, temp[:])
	plaintext = append(plaintext, mac...)

	if len(plaintext)%16 == 0 {
		for i := 0; i < 16; i++ {
			plaintext = append(plaintext, 0x10)
		}
	} else {
		remain := 16 - len(plaintext)%16
		pad := byte(remain)
		for i := 0; i < remain; i++ {
			plaintext = append(plaintext, pad)
		}
	}
	//fmt.Println("Plaintext:", plaintext, "\nlenth:", len(plaintext), "\n")
	blocksize := 16
	blockcount := len(plaintext) / 16
	ciphertext := make([]byte, blocksize*blockcount)
	xor_product := Xor(iv, plaintext[0:blocksize])
	cipher_block.Encrypt(ciphertext[0:blocksize], xor_product)
	for i := 1; i < blockcount; i++ {
		n1 := (i - 1) * blocksize
		n2 := i * blocksize
		n3 := (i + 1) * blocksize
		xor_product = Xor(ciphertext[n1:n2], plaintext[n2:n3])
		cipher_block.Encrypt(ciphertext[n2:n3], xor_product)
	}
	fin_ciphertext := make([]byte, len(iv)+len(ciphertext))
	fin_ciphertext = append(iv, ciphertext...)
	//fmt.Println("fin_ciphertext:", len(fin_ciphertext), "\n")
	return fin_ciphertext

}

func Dec(fin_ciphertext []byte, Dec_key []byte, Mac_key []byte) []byte {

	cipher_block, err := aes.NewCipher(Dec_key)
	if err != nil {
		fmt.Println("Key Error\n")
	}
	blocksize := 16
	iv := fin_ciphertext[0:16]
	ciphertext := fin_ciphertext[16:len(fin_ciphertext)]
	blockcount := len(ciphertext) / blocksize
	plaintext := make([]byte, len(ciphertext))
	cipher_block.Decrypt(plaintext[0:blocksize], ciphertext[0:blocksize])
	copy(plaintext[0:blocksize], Xor(iv, plaintext[0:blocksize]))
	for i := 1; i < blockcount; i++ {
		n1 := blocksize * (i - 1)
		n2 := blocksize * i
		n3 := blocksize * (i + 1)
		cipher_block.Decrypt(plaintext[n2:n3], ciphertext[n2:n3])
		copy(plaintext[n2:n3], Xor(ciphertext[n1:n2], plaintext[n2:n3]))
	}
	//fmt.Println("recovered plaintext:", plaintext, "\nlenth:", len(plaintext), "\n")

	paddingBlock := plaintext[len(plaintext)-1]
	padding_point := len(plaintext) - int(paddingBlock)
	for i := len(plaintext) - 1; i >= padding_point; i-- {
		if plaintext[i] != paddingBlock {
			fmt.Println("INVALID PADDING\n")
			os.Exit(1)
		}
	}

	mactag_point := padding_point - 32
	mac := make([]byte, 32)
	mac = plaintext[mactag_point:padding_point]
	message := make([]byte, mactag_point)
	message = plaintext[0:mactag_point]
	re_mac := hmac(Mac_key, message)
	for i := 0; i < 32; i++ {
		if mac[i] != re_mac[i] {
			fmt.Println("INVALID MAC\n")
			os.Exit(1)
		}
	}

	return message

}

func main() {
	mode := os.Args[1]
	allkey := os.Args[3]
	inputfile := os.Args[5]
	outputfile := os.Args[7]
	hexkey := allkey[0:32]
	hexmackey := allkey[32:64]
	key, _ := hex.DecodeString(hexkey)
	Mac_key, _ := hex.DecodeString(hexmackey)
	input, err := ioutil.ReadFile(inputfile)
	if err != nil {
		fmt.Println("Error while reading the input file:\n")
		os.Exit(1)
	}
	if mode == "encrypt" {
		iv := make([]byte, 16)
		_, err := rand.Read(iv)
		if err != nil {
			fmt.Println("Error while generating IV\n")
		}
		ciphertext := Enc(input, iv, key, Mac_key)
		err = ioutil.WriteFile(outputfile, ciphertext, 0644)
		if err != nil {
			fmt.Println("Error while writing the output file\n")
		} else {
			fmt.Println("File encrypted!")
		}
	}

	if mode == "decrypt" {
		plaintext := Dec(input, key, Mac_key)
		err := ioutil.WriteFile(outputfile, plaintext, 0644)
		if err != nil {
			fmt.Println("Error while writing the output file\n")
		} else {
			fmt.Println("File decrypted!")
		}
	}

}
