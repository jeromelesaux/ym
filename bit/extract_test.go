package bit_test

import (
	"fmt"
	"testing"

	"github.com/jeromelesaux/ym/bit"
	"github.com/stretchr/testify/assert"
)

func TestExtractBit(t *testing.T) {
	t.Run("value 1 0001", func(t *testing.T) {
		fmt.Printf("%.4b\n", 1)
		assert.Equal(t, uint8(1), bit.Get(1, bit.B0))
		assert.Equal(t, uint8(0), bit.Get(1, bit.B1))
		assert.Equal(t, uint8(0), bit.Get(1, bit.B2))
		assert.Equal(t, uint8(0), bit.Get(1, bit.B3))
	})
	t.Run("value 2 0010", func(t *testing.T) {
		fmt.Printf("%.4b\n", 2)
		assert.Equal(t, uint8(0), bit.Get(2, bit.B0))
		assert.Equal(t, uint8(1), bit.Get(2, bit.B1))
		assert.Equal(t, uint8(0), bit.Get(2, bit.B2))
		assert.Equal(t, uint8(0), bit.Get(2, bit.B3))
	})
	t.Run("value E 1110", func(t *testing.T) {
		fmt.Printf("%.4b\n", 0xE)
		assert.Equal(t, uint8(0), bit.Get(0xE, bit.B0))
		assert.Equal(t, uint8(1), bit.Get(0xE, bit.B1))
		assert.Equal(t, uint8(1), bit.Get(0xE, bit.B2))
		assert.Equal(t, uint8(1), bit.Get(0xE, bit.B3))
	})

	t.Run("value A 1010", func(t *testing.T) {
		fmt.Printf("%.4b\n", 0xA)
		assert.Equal(t, uint8(0), bit.Get(0xA, bit.B0))
		assert.Equal(t, uint8(1), bit.Get(0xA, bit.B1))
		assert.Equal(t, uint8(0), bit.Get(0xA, bit.B2))
		assert.Equal(t, uint8(1), bit.Get(0xA, bit.B3))
	})

	t.Run("value 9 1001", func(t *testing.T) {
		fmt.Printf("%b\n", 0x9)
		assert.Equal(t, uint8(1), bit.Get(0x9, bit.B0))
		assert.Equal(t, uint8(0), bit.Get(0x9, bit.B1))
		assert.Equal(t, uint8(0), bit.Get(0x9, bit.B2))
		assert.Equal(t, uint8(1), bit.Get(0x9, bit.B3))
	})
}
