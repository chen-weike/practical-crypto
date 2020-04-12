package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

func pow(x, y *big.Int) *big.Int {

	output := big.NewInt(1)
	for i := 0; i <= y.BitLen(); i++ {
		if y.Bit(i) == 1 {
			output = big.NewInt(0).Mul(output, x)
		}
		x = big.NewInt(0).Mul(x, x)
	}
	return output
}

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

func primetest(n *big.Int, k int) bool {
	one := big.NewInt(1)
	two := big.NewInt(2)
	subn := big.NewInt(0).Sub(n, one)
	d := subn
	s := 0
	for big.NewInt(0).Mod(d, two).Cmp(big.NewInt(0)) == 0 {
		d = big.NewInt(0).Div(d, two)
		s++
	}
	for i := 0; i < k; i++ {
		flag := false
		a, err := rand.Int(rand.Reader, n)
		if err != nil {
			fmt.Println(err)
		}
		for a.Cmp(two) == -1 {
			a, err = rand.Int(rand.Reader, n)
			if err != nil {
				fmt.Println(err)
			}
		}
		x := pow_mod(a, d, n)
		if x.Cmp(one) == 0 || x.Cmp(subn) == 0 {
			continue
		}

		for r := 1; r <= s-1; r++ {
			x = pow_mod(x, two, n)
			if x.Cmp(subn) == 0 {
				flag = true
				r = s - 1
			}
		}
		if flag == false {
			return false
		}

	}
	return true

}

func generator(p *big.Int, q *big.Int) *big.Int {
	factors := make(map[string]*big.Int)
	gsize := big.NewInt(0).Sub(p, big.NewInt(1))
	a := big.NewInt(0).Div(gsize, q)
	if big.NewInt(0).Mod(a, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		factors[big.NewInt(2).String()] = big.NewInt(2)
		for big.NewInt(0).Mod(a, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			a = big.NewInt(0).Div(a, big.NewInt(2))
		}
	}
	i := big.NewInt(3)
	for a.Cmp(big.NewInt(1)) == 1 {
		if big.NewInt(0).Mod(a, i).Cmp(big.NewInt(0)) == 0 {
			factors[i.String()] = i
			a = big.NewInt(0).Div(a, i)
		} else {
			i = big.NewInt(0).Add(i, big.NewInt(1))
		}
	}
	factors[q.String()] = q
	j := big.NewInt(2)
	for gsize.Cmp(j) > 0 {
		flag := true
		for _, value := range factors {
			if pow_mod(j, big.NewInt(0).Div(gsize, value), p).Cmp(big.NewInt(1)) == 0 {
				flag = false
			}
		}
		if flag == true {
			return j
		}
		j = big.NewInt(0).Add(j, big.NewInt(1))

	}
	fmt.Println("cannot get generator!")
	return gsize

}

func writefile(file string, x, y, z *big.Int) {
	var data bytes.Buffer
	data.WriteString("(")
	data.WriteString(x.String())
	data.WriteString(",")
	data.WriteString(y.String())
	data.WriteString(",")
	data.WriteString(z.String())
	data.WriteString(")")
	err := ioutil.WriteFile(file, []byte(data.String()), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	tobobfile := os.Args[1]
	toalicefile := os.Args[2]
	one := big.NewInt(1)
	//pmax := big.NewInt(0).Sub(pow(big.NewInt(2), big.NewInt(1024)), one)
	qmax := big.NewInt(0).Sub(pow(big.NewInt(2), big.NewInt(1020)), one)
	q, err := rand.Int(rand.Reader, qmax)
	if err != nil {
		fmt.Println(err)
	}
	for {

		if q.Cmp(big.NewInt(2048)) >= 0 && primetest(q, 6) {
			break
		}
		q, err = rand.Int(rand.Reader, qmax)
		if err != nil {
			fmt.Println(err)
		}

	}

	i := int64(2)
	p := new(big.Int).Add(big.NewInt(0).Mul(q, big.NewInt(i)), one)
	for {
		if (p.Cmp(big.NewInt(0).Sub(pow(big.NewInt(2), big.NewInt(1023)), one)) == 1) && (primetest(p, 6)) {
			break
		}
		i++
		p = new(big.Int).Add(big.NewInt(0).Mul(q, big.NewInt(i)), one)
	}
	g := generator(p, q)
	a, err := rand.Int(rand.Reader, big.NewInt(0).Sub(p, big.NewInt(1)))
	if err != nil {
		fmt.Println(err)
	}
	g_a := pow_mod(g, a, p)
	writefile(tobobfile, p, g, g_a)
	writefile(toalicefile, p, g, a)

}
