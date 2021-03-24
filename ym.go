package ym

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
