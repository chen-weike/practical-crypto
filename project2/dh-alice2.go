package main

import (
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
	frombob := os.Args[1]
	stored := os.Args[2]
	data1, err := ioutil.ReadFile(frombob)
	if err != nil {
		fmt.Println(err)
	}
	data1s := string(data1)
	g_b, _ := big.NewInt(0).SetString(data1s[1:len(data1s)-1], 10)

	data2, err := ioutil.ReadFile(stored)
	if err != nil {
		fmt.Println(err)
	}
	data2s := strings.Split(string(data2), ",")
	p, _ := big.NewInt(0).SetString(data2s[0][1:len(data2s[0])], 10)
	//g, _ := big.NewInt(0).SetString(data2s[1][0:len(data2s[1])], 10)
	a, _ := big.NewInt(0).SetString(data2s[2][0:len(data2s[2])-1], 10)

	public_key := pow_mod(g_b, a, p)
	fmt.Println(public_key)
}
