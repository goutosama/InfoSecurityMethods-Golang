package main

import (
	"fmt"
	"math/bits"
)

func main() {
	seed := uint64(0b1001)
	mask := uint64(0b1001)

	plaintext := []byte("HELLO")

	// Шифрование
	cipher := encode(plaintext, seed, mask, 4)

	fmt.Println("Encrypted:", cipher)

	// Дешифрование (ВАЖНО: тот же seed!)
	decrypted := encode(cipher, seed, mask, 4)

	fmt.Println("Decrypted:", string(decrypted))
}

func NextByte(state, seed uint64, mask uint64, size uint) uint8 {
	var result uint8 = 0

	for i := 0; i < 8; i++ {
		_, bit := Step(state, seed, mask, size)
		result |= uint8(bit) << i
	}

	return result
}

func encode(phrase []byte, seed uint64, mask uint64, size uint) []byte {
	state := seed & ((1 << size) - 1)
	result := make([]byte, len(phrase))

	for i := 0; i < len(phrase); i++ {
		keyByte := NextByte(state, seed, mask, size)
		result[i] = phrase[i] ^ keyByte
	}

	return result
}

func Step(state uint64, seed uint64, mask uint64, size uint) (newState uint64, outBit uint64) {
	// Бит, который "выпадает"
	outBit = state & 1

	// Обратная связь
	tapped := state & mask
	newBit := bits.OnesCount64(tapped) % 2

	// Сдвиг
	state >>= 1

	// Вставка нового бита
	state |= uint64(newBit) << (size - 1)

	return state, outBit
}
