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
func brute(p, g, h *big.Int) *big.Int {
	x := big.NewInt(0)
	for x.Cmp(p) == -1 {
		if pow_mod(g, x, p).Cmp(h) == 0 {
			return x
		} else {
			x = big.NewInt(0).Add(x, big.NewInt(1))
		}
	}
	return nil
}

func main() {
	input := os.Args[1]
	data, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Println(err)
	}
	datastring := strings.Split(string(data), ",")
	p, _ := big.NewInt(0).SetString(datastring[0][1:len(datastring[0])], 10)
	g, _ := big.NewInt(0).SetString(datastring[1][0:len(datastring[1])], 10)
	h, _ := big.NewInt(0).SetString(datastring[2][0:len(datastring[2])-1], 10)
	x := brute(p, g, h)
	if x == nil {
		fmt.Println("No X Found.")
	} else {
		fmt.Println(x)
	}
}
