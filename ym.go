package ym

type Ym struct {
	FileID         uint32
	CheckString    [8]byte
	NbFrames       uint32
	SongAttributes uint32
	DigidrumNb     uint16
	YmMasterClock  uint32
	FrameHz        uint32
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
	return &Ym{}
}
