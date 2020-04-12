package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

func pow_mod(x, y, n *big.Int) *big.Int {

	output := big.NewInt(1)
	for i := 0; i <= y.BitLen(); i++ {
		if y.Bit(i) == 1 {
			output = big.NewInt(0).Mod(big.NewInt(0).Mul(output, x), n)
		}
		x = big.NewInt(0).Mod(big.NewInt(0).Mul(x, x), n)
	}
	return output
}

func decrypt(c, k []byte) string {
	IV := make([]byte, 16)
	for i := 0; i < 16; i++ {
		IV[i] = c[i]
	}

	cmsg := make([]byte, (len(c) - 16))
	for i := 16; i < len(c); i++ {
		cmsg[i-16] = c[i]
	}

	aescipher, err1 := aes.NewCipher(k)
	if err1 != nil {
		fmt.Println("Error!")
		os.Exit(0)
	}

	aesgcm, err2 := cipher.NewGCMWithNonceSize(aescipher, 16)
	if err2 != nil {
		fmt.Println("Error!")
		os.Exit(0)
	}
	plaintext, err3 := aesgcm.Open(nil, IV, cmsg, nil)
	if err3 != nil {
		fmt.Println("error")
		os.Exit(0)
	}

	return string(plaintext)
}

func main() {
	cipherfile := os.Args[1]
	secretfile := os.Args[2]
	secret, err := ioutil.ReadFile(secretfile)
	if err != nil {
		fmt.Println(err)
	}
	secretstring := strings.Split(string(secret), ",")
	p, _ := big.NewInt(0).SetString(secretstring[0][1:len(secretstring[0])], 10)
	g, _ := big.NewInt(0).SetString(secretstring[1][0:len(secretstring[1])], 10)
	a, _ := big.NewInt(0).SetString(secretstring[2][0:len(secretstring[2])-1], 10)

	cdata, err1 := ioutil.ReadFile(cipherfile)
	if err != nil {
		fmt.Println(err1)
	}
	cdatastring := strings.Split(string(cdata), ",")
	g_b, _ := big.NewInt(0).SetString(cdatastring[0][1:len(cdatastring[0])], 10)
	cipher, _ := hex.DecodeString(cdatastring[1][0 : len(cdatastring[1])-1])

	g_a := pow_mod(g, a, p)
	g_ab := pow_mod(g_b, a, p)

	var key bytes.Buffer
	key.WriteString(g_a.String())
	key.WriteString(" ")
	key.WriteString(g_b.String())
	key.WriteString(" ")
	key.WriteString(g_ab.String())
	h := sha256.New()
	h.Write([]byte(key.String()))
	k := h.Sum(nil)
	plaintext := decrypt(cipher, k)
	fmt.Println(plaintext)

}
