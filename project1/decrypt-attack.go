package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func attack(ciphertext []byte) []byte {
	IV := ciphertext[:16]
	//fmt.Println("IV:", IV, "\n")
	ciphertext = ciphertext[16:]
	//fmt.Println("Cipher:", ciphertext, "\n")
	blockcount := len(ciphertext) / 16
	preblock := make([]byte, 16)
	block := make([]byte, 16)
	plaintblock := make([]byte, 16)
	plaintext := make([]byte, 0)
	for i := blockcount - 1; i >= 0; i-- {
		if i == 0 {
			preblock = IV
		} else {
			preblock = ciphertext[16*(i-1) : 16*i]
		}
		block = ciphertext[16*i : 16*(i+1)]
		//fmt.Println("preblock:", preblock, "\n")
		//fmt.Println("block:", block, "\n")
		plaintblock = attackperblock(preblock, block)
		//fmt.Println("plaintblock:", plaintblock, "\n")
		plaintext = append(plaintblock, plaintext...)
	}
	pad := plaintext[len(plaintext)-1]
	n := int(pad)
	plaintext = plaintext[:len(plaintext)-n-32]
	return plaintext

}

func attackperblock(preblock, block []byte) []byte {
	tmpcipher := make([]byte, 64)
	for i := 0; i < 16; i++ {
		tmpcipher[48+i] = block[i]
	}
	tmppad := make([]byte, 16)
	try := make([]byte, 16)
	for j := 1; j <= 16; j++ {
		for k := 1; k <= j; k++ {
			tmppad[16-k] = byte(j)
		}

		for l := 2; l <= 256; l++ {
			try[16-j] = byte(l)
			for n := 0; n < 16; n++ {
				tmpcipher[32+n] = preblock[n] ^ try[n] ^ tmppad[n]
			}
			temp := "temp.txt"
			err_write := ioutil.WriteFile(temp, tmpcipher, 0644)
			if err_write != nil {
				fmt.Println("ERROR: ", err_write)
			}

			result, err_cmd := exec.Command("./decrypt-test", "-i", temp).Output()
			if err_cmd != nil {
				fmt.Println("ERROR: ", err_cmd)
			}
			if string(result) != "INVALID PADDING" {
				break
			}
		}
	}
	return try
}

func main() {
	file := os.Args[2]
	ciphertext, err_fileopen := ioutil.ReadFile(file)
	if err_fileopen != nil {
		fmt.Println("Error - Reading the file\n\n")
		os.Exit(1)
	}
	plaintext := attack(ciphertext)
	fmt.Println(string(plaintext))
}
