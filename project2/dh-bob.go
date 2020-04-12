package main

import (
	"bytes"
	"crypto/rand"
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

func main() {
	input := os.Args[1]
	output := os.Args[2]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println(err)
	}
	datastring := strings.Split(string(data), ",")
	p, _ := big.NewInt(0).SetString(datastring[0][1:len(datastring[0])], 10)
	g, _ := big.NewInt(0).SetString(datastring[1][0:len(datastring[1])], 10)
	g_a, _ := big.NewInt(0).SetString(datastring[2][0:len(datastring[2])-1], 10)
	b, err := rand.Int(rand.Reader, big.NewInt(0).Sub(p, big.NewInt(1)))
	if err != nil {
		fmt.Println(err)
	}
	g_b := pow_mod(g, b, p)
	var toalice bytes.Buffer
	toalice.WriteString("(")
	toalice.WriteString(g_b.String())
	toalice.WriteString(")")
	err = ioutil.WriteFile(output, []byte(toalice.String()), 0644)
	if err != nil {
		fmt.Println(err)
	}
	public_key := pow_mod(g_a, b, p)
	fmt.Println(public_key)
}
