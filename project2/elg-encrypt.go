package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

func encrypt(msg string, k []byte) ([]byte, []byte) {
	aesCipher, err := aes.NewCipher(k)
	if err != nil {
		panic(err.Error())
	}
	IV := make([]byte, 16)
	_, err1 := io.ReadFull(rand.Reader, IV)
	if err1 != nil {
		panic(err.Error())
	}
	aesgcm, err2 := cipher.NewGCMWithNonceSize(aesCipher, 16)
	if err2 != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, IV, []byte(msg), nil)

	return ciphertext, IV

}

func main() {
	msg := os.Args[1]
	pkeyfile := os.Args[2]
	cipherfile := os.Args[3]
	pkey, err := ioutil.ReadFile(pkeyfile)
	if err != nil {
		fmt.Println(err)
	}
	pkeystring := strings.Split(string(pkey), ",")
	p, _ := big.NewInt(0).SetString(pkeystring[0][1:len(pkeystring[0])], 10)
	g, _ := big.NewInt(0).SetString(pkeystring[1][0:len(pkeystring[1])], 10)
	g_a, _ := big.NewInt(0).SetString(pkeystring[2][0:len(pkeystring[2])-1], 10)
	b, err := rand.Int(rand.Reader, big.NewInt(0).Sub(p, big.NewInt(1)))
	if err != nil {
		fmt.Println(err)
	}
	g_b := pow_mod(g, b, p)
	g_ab := pow_mod(g_a, b, p)
	var key bytes.Buffer
	key.WriteString(g_a.String())
	key.WriteString(" ")
	key.WriteString(g_b.String())
	key.WriteString(" ")
	key.WriteString(g_ab.String())
	h := sha256.New()
	h.Write([]byte(key.String()))
	k := h.Sum(nil)
	cipher, IV := encrypt(msg, k)

	var cmsg bytes.Buffer
	cmsg.WriteString(hex.EncodeToString(IV))
	cmsg.WriteString(hex.EncodeToString(cipher))

	var output bytes.Buffer
	output.WriteString("(")
	output.WriteString(g_b.String())
	output.WriteString(",")
	output.WriteString(cmsg.String())
	output.WriteString(")")
	err3 := ioutil.WriteFile(cipherfile, []byte(output.String()), 0644)
	if err3 != nil {
		fmt.Println("err")
		os.Exit(0)
	}

}
