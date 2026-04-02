package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func main() {
	// обработка флагов командной строки
	// -sign для подписи, -verify - проверки
	// -codeword кодовое слово
	// последний параметр - шифруемая фраза
	isSign := flag.Bool("sign", false, "Sign message with a EGSA signature, outputs the open key and the signature (two integer numbers)")
	isVerify := flag.Bool("verify", false, "Verify message signature by providing the open key with -openkey, signature with -s1 -s2 and the message (always last)")
	openkey := flag.String("openkey", "", "A codeword to use with a phrase")
	sign1 := flag.String("s1", "", "First part of the signature")
	sign2 := flag.String("s2", "", "Second part of the signature")
	isFile := flag.String("file", "", "A path to the file, if not specified works with a message string")

	flag.Parse()

	if len(os.Args) > 1 {
		var M string
		if *isFile == "" {
			M = os.Args[len(os.Args)-1]
		} else {
			MByte, err := os.ReadFile(*isFile)
			if err != nil {
				fmt.Println("Open file error: ", err)
				return
			}
			M = string(MByte)
		}

		p := big.NewInt(23)
		g := big.NewInt(5)

		var m big.Int = h(M)

		if !(*isSign && *isVerify) {
			if *isSign {
				x, _ := rand.Int(rand.Reader, p)
				x.Add(x, big.NewInt(1))
				y := new(big.Int).Exp(g, x, p)
				fmt.Println("Open key: ", y)
				r, s := sign(p, g, x, &m)
				fmt.Println("Signature", r, s)
			}
			if *isVerify && *openkey != "" && *sign1 != "" && *sign2 != "" {
				y, err := strconv.ParseInt(*openkey, 10, 64)
				if err != nil {
					fmt.Println(err)
					return
				}
				r, err := strconv.ParseInt(*sign1, 10, 64)
				if err != nil {
					fmt.Println(err)
					return
				}
				s, err := strconv.ParseInt(*sign2, 10, 64)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println("Valid: ", verify(p, g, big.NewInt(y), &m, big.NewInt(r), big.NewInt(s)))
			}
		}
	}
}

func sign(p, g, x, m *big.Int) (*big.Int, *big.Int) {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	k := randomK(pMinus1)
	kInv := new(big.Int).ModInverse(k, pMinus1)
	r := new(big.Int).Exp(g, k, p)
	xr := new(big.Int).Mul(x, r)
	mMinusXr := new(big.Int).Sub(m, xr)
	mMinusXr.Mod(mMinusXr, pMinus1)

	s := new(big.Int).Mul(kInv, mMinusXr)
	s.Mod(s, pMinus1)
	return r, s
}

func randomK(pMinus1 *big.Int) *big.Int {
	for {
		k, _ := rand.Int(rand.Reader, pMinus1)
		if k.Cmp(big.NewInt(1)) > 0 && new(big.Int).GCD(nil, nil, k, pMinus1).Cmp(big.NewInt(1)) == 0 {
			return k
		}
	}
}

func verify(p, g, y, m, r, s *big.Int) bool {
	right := new(big.Int).Exp(g, m, p)

	left := new(big.Int).Mul(y.Exp(y, r, nil), r.Exp(r, s, nil))
	left.Mod(left, p)

	return left.Cmp(right) == 0
}

func h(s string) big.Int {
	var hashInt big.Int

	hash := sha256.Sum256([]byte(s))
	hashInt.SetBytes(hash[:])

	return hashInt
}
