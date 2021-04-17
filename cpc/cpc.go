package cpc

import (
	"fmt"

	"github.com/jeromelesaux/ym"
)

var (
	Register16bitsMaxIndice = 3
	Register8bitsMaxIndice  = 12
)

type CpcYM struct {
	*ym.Ym
}

func NewCpcYM() *CpcYM {
	y := &CpcYM{
		Ym: ym.NewYm(),
	}
	return y
}

func NewCpcYMFromYM(y0 *ym.Ym) *CpcYM {
	y := NewCpcYM()
	y.Ym = ym.CopyYm(y0)
	for i := 0; i < 3; i++ {
		register := i * 2
		for frame := 0; frame < int(y.NbFrames); frame++ {
			reg0 := y.Data[register][frame] & 0xf0 / 2
			reg1 := y.Data[register+1][frame] & 0xF0
			y.Data[register][frame] = reg0
			y.Data[register+1][frame] = reg1
		}

	}
	return y
}

// Entered    returned
// 0          1 - 0
// 1          3 - 2
// 2          5 - 4
func (c *CpcYM) GetRegister16bits(registerNumber, frame int) (uint16, error) {
	if frame > int(c.NbFrames) {
		return 0, fmt.Errorf("the frame is overload the frame number")
	}
	if registerNumber >= Register16bitsMaxIndice {
		return 0, fmt.Errorf("register 16 bits exceed register allowed")
	}
	registerNumber *= 2
	reg0 := c.Data[registerNumber][frame]
	reg1 := c.Data[registerNumber+1][frame]
	reg1b := uint16(reg1)<<8 + uint16(reg0)
	return reg1b, nil
}

// Entered    returned
// 3          6
// 4          7
// 5          8
// ............
// 12           15
func (c *CpcYM) GetRegister8bits(registerNumber, frame int) (byte, error) {
	if frame > int(c.NbFrames) {
		return 0, fmt.Errorf("the frame is overload the frame number")
	}
	if registerNumber < Register16bitsMaxIndice {
		return 0, fmt.Errorf("register 8 bits is not allowed")
	}
	return c.Data[registerNumber+3][frame], nil
}

func (c *CpcYM) SetRegister16bits(registerNumber, frame int, value uint16) error {
	if frame > int(c.NbFrames) {
		return fmt.Errorf("the frame is overload the frame number")
	}
	if registerNumber >= Register16bitsMaxIndice {
		return fmt.Errorf("register 16 bits exceed register allowed")
	}
	registerNumber *= 2
	reg1 := byte(value >> 8)
	reg0 := byte(value)
	c.Data[registerNumber][frame] = reg0
	c.Data[registerNumber+1][frame] = reg1
	return nil
}

func (c *CpcYM) SetRegister8bits(registerNumber, frame int, value byte) error {
	if frame > int(c.NbFrames) {
		return fmt.Errorf("the frame is overload the frame number")
	}
	if registerNumber < Register16bitsMaxIndice {
		return fmt.Errorf("register 8 bits is not allowed for this indice")
	}
	c.Data[registerNumber+3][frame] = value
	return nil
}
