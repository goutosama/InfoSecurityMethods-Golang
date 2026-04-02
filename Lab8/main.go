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
	// -enc для шифрования, -dec - расшифрования
	// для шифрования используется закрытый ключ (передается через -private) и общий ключ (параметр -n)
	// для дешифрования используется открытый ключ (передается через -open) и общий ключ (параметр -n)
	// -keys - генерация открытого и закрытого ключа
	// последний параметр - шифруемая фраза
	// если присутствует параметр -file шифруется файл по заданному в параметре пути
	isEncode := flag.Bool("enc", false, "Encode message with a RSA signature, using the private key with parameters -private and -n")
	isDecode := flag.Bool("dec", false, "Decode message with RSA by providing the open key with parameters -openkey and -n")
	isKeys := flag.Bool("keys", false, "Generate a pair of private and open key")
	openkey := flag.String("open", "", "Open key for decoding")
	private := flag.String("private", "", "Private key for encoding")
	n := flag.String("n", "", "Second part of the key, shared with open and private")
	isFile := flag.String("file", "", "A path to the file, if not specified works with a message string")

	flag.Parse()

	if len(os.Args) > 1 {
		e := big.NewInt(65537)
		if *isKeys {
			var p *big.Int
			var q *big.Int

			for {
				var err error
				p, err = rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil))
				if err != nil {
					fmt.Println(err)
					continue
				}
				if p.ProbablyPrime(3) {
					break
				}
			}
			for {
				var err error
				q, err = rand.Int(rand.Reader, new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil))
				if err != nil {
					fmt.Println(err)
					continue
				}
				if q.ProbablyPrime(3) {
					break
				}
			}

			var n *big.Int = new(big.Int).Mul(p, q)
			var fn *big.Int = euler(p, q)

			d := new(big.Int).ModInverse(e, fn)
			if d == nil {
				fmt.Println("D was not found")
			}
			fmt.Println("Open key: ", e, n)
			fmt.Println("Private key: ", d, n)
		} else {

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

			var m *big.Int = h(M)

			if !(*isEncode && *isDecode) && (*n != "") {
				if *isEncode && *private != "" {
					n1, err := strconv.ParseInt(*n, 10, 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					privateN, err := strconv.ParseInt(*private, 10, 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(new(big.Int).Exp(m, big.NewInt(privateN), big.NewInt(n1)))
				}
				if *isDecode && *openkey != "" {
					n1, err := strconv.ParseInt(*n, 10, 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					openkeyN, err := strconv.ParseInt(*openkey, 10, 64)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(new(big.Int).Exp(m, big.NewInt(openkeyN), big.NewInt(n1)))
				}
			}
		}
	}
}

func euler(p, q *big.Int) *big.Int {
	return new(big.Int).Mul(p.Sub(p, big.NewInt(1)), q.Sub(q, big.NewInt(1)))
}

func h(s string) *big.Int {
	var hashInt big.Int

	hash := sha256.Sum256([]byte(s))
	hashInt.SetBytes(hash[:])

	return &hashInt
}

// func euclidian(av, bv *big.Int) *big.Int {
// 	var a *big.Int
// 	var b *big.Int
// 	if av.Cmp(bv) >= 0 {
// 		a = av
// 		b = bv
// 	} else {
// 		a = bv
// 		b = av
// 	}
// 	for {
// 		r := new(big.Int).Mod(a, b)
// 		if r == big.NewInt(0) {
// 			return b
// 		} else {
// 			b = r
// 			b = a
// 		}
// 	}
// }
