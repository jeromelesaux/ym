package ym

const (
	YM1     = 0x594D3121
	YM2     = 0x594D3221 // ('Y' << 24) | ('M' << 16) | ('2' << 8) | ('!'),
	YM3     = 0x594D3321 //('Y' << 24) | ('M' << 16) | ('3' << 8) | ('b') ('Y' << 24) | ('M' << 16) | ('3' << 8) | ('!')
	YM4     = 0x594D3421 //('Y' << 24) | ('M' << 16) | ('4' << 8) | ('!')
	YM5     = 0x594D3521 //('Y' << 24) | ('M' << 16) | ('5' << 8) | ('!')
	YM6     = 0x594D3621 // ('Y' << 24) | ('M' << 16) | ('6' << 8) | ('!')
	YM_MIX1 = ('M' << 24) | ('I' << 16) | ('X' << 8) | ('1')
	YMT1    = ('Y' << 24) | ('M' << 16) | ('T' << 8) | ('1')
	YMT2    = ('Y' << 24) | ('M' << 16) | ('T' << 8) | ('2')
)

type Ym struct {
	FileID         uint32
	CheckString    [8]byte
	NbFrames       uint32
	SongAttributes uint32
	DigidrumNb     uint16
	YmMasterClock  uint32
	FrameHz        uint16
	LoopFrame      uint32
	Size           uint16
	Digidrums      []Digidrum
	SongName       []byte
	AuthorName     []byte
	SongComment    []byte
	Data           [16][]byte
	EndID          uint32
}

type Digidrum struct {
	SampleSize uint32
	SampleData []byte
}

func NewYm() *Ym {
	y := &Ym{
		Digidrums:   make([]Digidrum, 0),
		SongName:    make([]byte, 0),
		AuthorName:  make([]byte, 0),
		SongComment: make([]byte, 0),
	}
	for i := 0; i < 16; i++ {
		y.Data[i] = make([]byte, 0)
	}
	copy(y.CheckString[:], []byte("LeOnArD!"))

	return y
}

func CopyYm(ym *Ym) *Ym {
	n := NewYm()
	n.FileID = ym.FileID
	n.NbFrames = ym.NbFrames
	n.SongAttributes = ym.SongAttributes
	n.YmMasterClock = ym.YmMasterClock
	n.FrameHz = ym.FrameHz
	n.LoopFrame = ym.LoopFrame
	n.Size = ym.Size

	n.DigidrumNb = ym.DigidrumNb
	n.Digidrums = make([]Digidrum, ym.DigidrumNb)
	for i := 0; i < int(ym.DigidrumNb); i++ {
		n.Digidrums[i].SampleSize = ym.Digidrums[i].SampleSize
		n.Digidrums[i].SampleData = make([]byte, n.Digidrums[i].SampleSize)
		copy(n.Digidrums[i].SampleData, ym.Digidrums[i].SampleData)
	}

	n.SongName = append(n.SongName, ym.SongName...)
	n.SongComment = append(n.SongComment, ym.SongComment...)
	n.AuthorName = append(n.AuthorName, ym.AuthorName...)

	n.Size = ym.Size
	for j := 0; j < 16; j++ {
		n.Data[j] = make([]byte, ym.NbFrames)
		for i := 0; i < int(ym.NbFrames); i++ {
			n.Data[j][i] = ym.Data[j][i]
		}
	}
	n.EndID = ym.EndID
	return n
}

func (y *Ym) FormatType() string {
	switch y.FileID {
	case YM1:
		return "YM1"
	case YM2:
		return "YM2"
	case YM3:
		return "YM3"
	case YM4:
		return "YM4"
	case YM5:
		return "YM5"
	case YM6:
		return "YM6"
	default:
		return "Unknown"
	}
}

func (y *Ym) Extract(startFrame, endFrame int) *Ym {
	n := NewYm()
	n.FileID = y.FileID
	n.NbFrames = y.NbFrames
	n.SongAttributes = y.SongAttributes
	n.YmMasterClock = y.YmMasterClock
	n.FrameHz = y.FrameHz
	n.LoopFrame = y.LoopFrame
	n.Size = y.Size

	n.DigidrumNb = y.DigidrumNb
	n.Digidrums = make([]Digidrum, y.DigidrumNb)
	for i := 0; i < int(y.DigidrumNb); i++ {
		n.Digidrums[i].SampleSize = y.Digidrums[i].SampleSize
		n.Digidrums[i].SampleData = make([]byte, n.Digidrums[i].SampleSize)
		copy(n.Digidrums[i].SampleData, y.Digidrums[i].SampleData)
	}

	n.SongName = append(n.SongName, y.SongName...)
	n.SongComment = append(n.SongComment, y.SongComment...)
	n.AuthorName = append(n.AuthorName, y.AuthorName...)

	for j := 0; j < 16; j++ {
		n.Data[j] = make([]byte, endFrame-startFrame)
		index := 0
		for i := startFrame; i < endFrame; i++ {
			n.Data[j][index] = y.Data[j][i]
			index++
		}
	}
	n.EndID = y.EndID
	n.NbFrames = uint32(endFrame) - uint32(startFrame)
	return n
}
