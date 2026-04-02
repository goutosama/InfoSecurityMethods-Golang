package main

import (
	"fmt"
	"math"
	"math/bits"
)

func main() {
	var state uint64 = 15
	var size int = bits.Len(uint(state))
	var mask uint = 0b1001

	var result uint64 = 0
	for i := 0; i < int(math.Pow(2, float64(size))-1); i++ {
		outBit := state & 1
		result = (result << 1) + outBit
		tapped := state & uint64(mask)
		newBit := bits.OnesCount64(tapped) % 2

		state >>= 1

		state |= uint64(newBit) << (size - 1)
	}

	fmt.Printf("%b\n", result)
}
