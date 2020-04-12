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
func BigSqrt(n *big.Int) (x *big.Int) {
	var px, nx big.Int
	x = big.NewInt(0)
	x.SetBit(x, n.BitLen()/2+1, 1)
	for {
		nx.Rsh(nx.Add(x, nx.Div(n, x)), 1)
		if nx.Cmp(x) == 0 || nx.Cmp(&px) == 0 {
			break
		}
		px.Set(x)
		x.Set(&nx)
	}
	return

}
func Bigstep_Babystep(p, g, h *big.Int) *big.Int {

	m := BigSqrt(p)
	m = m.Add(m, big.NewInt(1))
	lookup := make(map[string]*big.Int)

	i := big.NewInt(1)
	res := big.NewInt(1)

	lookup["0"] = big.NewInt(1)
	for i.Cmp(m) != 1 {
		res = res.Mul(res, g)
		res = res.Mod(res, p)
		if res.Cmp(big.NewInt(0)) == 0 || res.Cmp(big.NewInt(0)) == 0 {
			break
		}
		lookup[res.String()] = new(big.Int).Set(i)
		i = i.Add(i, big.NewInt(1))
	}
	ginv := new(big.Int).ModInverse(g, p)
	ginv = pow_mod(ginv, m, p)
	q := new(big.Int).Set(h)
	i = big.NewInt(0)

	for i.Cmp(m) < 1 {

		j, ok := lookup[q.String()]
		if ok {
			x := new(big.Int).Set(i)
			x = x.Mul(x, m)
			x = x.Add(x, j)
			return x
		}
		q = q.Mul(q, ginv)
		q = q.Mod(q, p)
		i = i.Add(i, big.NewInt(1))

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
	x := Bigstep_Babystep(p, g, h)
	if x == nil {
		fmt.Println("No X Found.")
	} else {
		fmt.Println(x)
	}
}
