package ym_test

import (
	"fmt"
	"testing"
)

func TestYMConstants(t *testing.T) {
	fmt.Printf("%X\n", ('Y'<<24)|('M'<<16)|('2'<<8)|('!'))
	fmt.Printf("%X\n", ('Y'<<24)|('M'<<16)|('3'<<8)|('!'))
}
