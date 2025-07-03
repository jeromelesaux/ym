package bit

type bit byte

const (
	B0 bit = 0b0001
	B1 bit = 0b0010
	B2 bit = 0b0100
	B3 bit = 0b1000
	B4 bit = 0b10000
	B5 bit = 0b100000
	B6 bit = 0b1000000
	B7 bit = 0b10000000
)

func Get(value byte, pos bit) byte {
	//fmt.Printf("org[%.4b][%d]  [%b] bit [%v] \n", value, pos, value&byte(pos), value&byte(pos) > 0)
	if value&byte(pos) > 0 {
		return 1
	}
	return 0
}

func Set(b byte, n bit) byte {
	return b | byte(n)
}
