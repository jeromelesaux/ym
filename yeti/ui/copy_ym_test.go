package ui_test

import (
	"fmt"
	"testing"

	"github.com/jeromelesaux/ym/bit"
	"github.com/jeromelesaux/ym/yeti/ui"
	"github.com/stretchr/testify/assert"
)

func TestConversion(t *testing.T) {
	var r1 byte = 0
	var r8 byte = 0x0D
	fmt.Printf("0x%X %b\n", r8, r8)
	r := bit.Get((r8), bit.B0)
	fmt.Printf("%X", r)
	r84 := bit.Set(bit.Get((r8), bit.B0), 4)
	r85 := bit.Set(bit.Get((r8), bit.B1), 5)
	r86 := bit.Set(bit.Get((r8), bit.B2), 6)
	r87 := bit.Set(bit.Get(r8, bit.B3), 7)

	fmt.Printf("%X,%X,%X,%X %b|%b|%b|%b %X\n", r84, r85, r86, r87, r84, r85, r86, r87, r84+r85+r86+r87)
	res := (r8 << 4) + r1

	fmt.Printf("Ox%X %d %b", res, res, res)

}

func TestIspair(t *testing.T) {
	assert.True(t, ui.IsPair(2))
	assert.False(t, ui.IsPair(3))
	assert.False(t, ui.IsPair(11))
}
