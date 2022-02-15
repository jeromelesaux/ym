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

type MixBlock struct {
	SampleStart  uint32
	SampleLength uint32
	NbRepeat     uint16
	ReplayFreq   uint16
}

type TimeKey struct {
	Time    uint32
	NRepeat uint16
	NBlock  uint16
}
type Ym struct {
	FileID           uint32
	CheckString      [8]byte
	NbVoice          int16
	TrackerFreqShift int
	NbFrames         uint32
	SongAttributes   uint32
	DigidrumNb       uint16
	YmMasterClock    uint32
	FrameHz          uint16
	LoopFrame        uint32
	Size             uint16
	Digidrums        []Digidrum
	NbMixBlock       uint32
	TimeInfo         []TimeKey
	NbTimeKey        int32
	MusicLenInMs     int32
	MixBlock         []MixBlock
	SongName         []byte
	AuthorName       []byte
	SongComment      []byte
	Data             [16][]byte
	EndID            uint32
}

type Digidrum struct {
	SampleSize uint32
	SampleData []byte
	RepLen     uint32
}

func NewYm() *Ym {
	y := &Ym{
		Digidrums:   make([]Digidrum, 0),
		MixBlock:    make([]MixBlock, 0),
		TimeInfo:    make([]TimeKey, 0),
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
		n.Digidrums[i].RepLen = ym.Digidrums[i].RepLen
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

	n.NbMixBlock = ym.NbMixBlock
	n.MixBlock = make([]MixBlock, ym.NbMixBlock)
	for i := 0; i < int(ym.NbMixBlock); i++ {
		n.MixBlock[i].NbRepeat = ym.MixBlock[i].NbRepeat
		n.MixBlock[i].ReplayFreq = ym.MixBlock[i].ReplayFreq
		n.MixBlock[i].SampleLength = ym.MixBlock[i].SampleLength
		n.MixBlock[i].SampleStart = ym.MixBlock[i].SampleStart
	}

	n.NbTimeKey = ym.NbTimeKey
	n.TimeInfo = make([]TimeKey, ym.NbTimeKey)
	for i := 0; i < int(ym.NbTimeKey); i++ {
		n.TimeInfo[i].Time = ym.TimeInfo[i].Time
		n.TimeInfo[i].NRepeat = ym.TimeInfo[i].NRepeat
		n.TimeInfo[i].NBlock = ym.TimeInfo[i].NBlock
	}
	n.MusicLenInMs = ym.MusicLenInMs
	n.NbVoice = ym.NbVoice
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
	case YMT1:
		return "YM Tracker 1"
	case YMT2:
		return "YM Tracker 2"
	case YM_MIX1:
		return "YM Mix"
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
		n.Digidrums[i].RepLen = y.Digidrums[i].RepLen
	}

	n.NbMixBlock = y.NbMixBlock
	n.MixBlock = make([]MixBlock, y.NbMixBlock)
	for i := 0; i < int(y.NbMixBlock); i++ {
		n.MixBlock[i].NbRepeat = y.MixBlock[i].NbRepeat
		n.MixBlock[i].ReplayFreq = y.MixBlock[i].ReplayFreq
		n.MixBlock[i].SampleLength = y.MixBlock[i].SampleLength
		n.MixBlock[i].SampleStart = y.MixBlock[i].SampleStart
	}

	n.NbTimeKey = y.NbTimeKey
	n.TimeInfo = make([]TimeKey, y.NbTimeKey)
	for i := 0; i < int(y.NbTimeKey); i++ {
		n.TimeInfo[i].Time = y.TimeInfo[i].Time
		n.TimeInfo[i].NRepeat = y.TimeInfo[i].NRepeat
		n.TimeInfo[i].NBlock = y.TimeInfo[i].NBlock
	}

	n.SongName = append(n.SongName, y.SongName...)
	n.SongComment = append(n.SongComment, y.SongComment...)
	n.AuthorName = append(n.AuthorName, y.AuthorName...)

	for j := 0; j < 16; j++ {
		n.Data[j] = make([]byte, endFrame-startFrame)
		index := 0
		for i := startFrame; i < endFrame && i < len(y.Data[j]); i++ {
			n.Data[j][index] = y.Data[j][i]
			index++
		}
	}
	n.MusicLenInMs = y.MusicLenInMs
	n.NbVoice = y.NbVoice
	n.EndID = y.EndID
	n.NbFrames = uint32(endFrame) - uint32(startFrame)
	return n
}

/*

	x x x x x x x x	r0
	0 0 0 0 x x x x	r1		// Special FX 1a
	x x x x x x x x	r2
	0 0 0 0 x x x x r3		// Special FX 2a
	x x x x x x x x r4
	0 0 0 0 x x x x r5
	0 0 0 x x x x x r6		// Special FX 1b
	0 0 x x x x x x r7
	0 0 0 x x x x x r8		// Special FX 2b
	0 0 0 x x x x x r9
	0 0 0 x x x x x r10
	x x x x x x x x r11
	x x x x x x x x r12
	0 0 0 0 x x x x r13
	0 0 0 0 0 0 0 0 r14		// Special FX 1c
	0 0 0 0 0 0 0 0 r15		// Special FX 2c


  Special Fx ?a
	0 0 0 0	: No special FX running
	0 0 0 1 : Sid Voice A
	0 0 1 0 : Sid Voice B
	0 0 1 1 : Sid Voice C
	0 1 0 0 : Extended Fx voice A
	0 1 0 1 : Digidrum voice A
	0 1 1 0 : Digidrum voice B
	0 1 1 1 : Digidrum voice C
	1 0 0 0 : Extended Fx voice B
	1 0 0 1 : Sinus SID voice A
	1 0 1 0 : Sinus SID voice B
	1 0 1 1 : Sinus SID voice C
	1 1 0 0 : Extended Fx voice C
	1 1 0 1 : Sync Buzzer voice A
	1 1 1 0 : Sync Buzzer voice B
	1 1 1 1 : Sync Buzzer voice C



*/

func (y *Ym) ComputeTime() {
	//-------------------------------------------
	// Compute nb of mixblock
	//-------------------------------------------
	y.NbTimeKey = 0

	for i := 0; i < int(y.NbMixBlock); i++ {
		if y.MixBlock[i].NbRepeat >= 32 {
			y.MixBlock[i].NbRepeat = 32
		}

		y.NbTimeKey += int32(y.MixBlock[i].NbRepeat)
	}

	//-------------------------------------------
	// Parse all mixblock keys
	//-------------------------------------------
	y.TimeInfo = make([]TimeKey, y.NbTimeKey)

	var time uint32

	for i := 0; i < int(y.NbMixBlock); i++ {
		for j := 0; j < int(y.MixBlock[i].NbRepeat); j++ {
			y.TimeInfo[j].Time = time
			y.TimeInfo[j].NRepeat = (y.MixBlock[i].NbRepeat) - uint16(j)
			y.TimeInfo[j].NBlock = uint16(i)

			time += (y.MixBlock[i].SampleLength * 1000) / uint32(y.MixBlock[i].ReplayFreq)
		}
	}
	y.MusicLenInMs = int32(time)

}
