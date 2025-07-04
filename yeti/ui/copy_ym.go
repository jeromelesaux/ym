package ui

import (
	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/bit"
)

func copyCPCYm(y *ym.Ym) *ym.Ym {
	n := ym.NewYm()
	n.FileID = y.FileID
	n.NbFrames = y.NbFrames
	n.SongAttributes = y.SongAttributes
	n.YmMasterClock = y.YmMasterClock
	n.FrameHz = y.FrameHz
	n.LoopFrame = y.LoopFrame
	n.Size = y.Size

	n.DigidrumNb = y.DigidrumNb
	n.Digidrums = make([]ym.Digidrum, y.DigidrumNb)
	for i := range int(y.DigidrumNb) {
		n.Digidrums[i].SampleSize = y.Digidrums[i].SampleSize
		n.Digidrums[i].RepLen = y.Digidrums[i].RepLen
		n.Digidrums[i].SampleData = make([]byte, n.Digidrums[i].SampleSize)
		copy(n.Digidrums[i].SampleData, y.Digidrums[i].SampleData)
	}

	n.SongName = append(n.SongName, y.SongName...)
	n.SongComment = append(n.SongComment, y.SongComment...)
	n.AuthorName = append(n.AuthorName, y.AuthorName...)

	n.Size = y.Size
	for j := range 16 {
		n.Data[j] = make([]byte, len(y.Data[j]))
		copy(n.Data[j][:], y.Data[j][:])
	}

	for i := range len(n.Data[0]) {
		r1 := n.Data[1][i]
		r3 := n.Data[3][i]
		r8 := n.Data[8][i]
		r9 := n.Data[9][i]
		r5 := n.Data[5][i]
		r6 := n.Data[6][i]
		r10 := n.Data[10][i]

		r1r8 := (r8 << 4) + r1

		r3r9 := bit.Set(bit.Get((r9), bit.B0), 4) +
			bit.Set(bit.Get((r9), bit.B1), 5) +
			bit.Set(bit.Get((r9), bit.B2), 6) +
			bit.Set(bit.Get(r9, bit.B3), 7) +
			bit.Set(bit.Get((r3), bit.B0), 0) +
			bit.Set(bit.Get((r3), bit.B1), 1) +
			bit.Set(bit.Get((r3), bit.B2), 2) +
			bit.Set(bit.Get(r3, bit.B3), 3)

		r5r10 := bit.Set(bit.Get((r10), bit.B0), 4) +
			bit.Set(bit.Get((r10), bit.B1), 5) +
			bit.Set(bit.Get((r10), bit.B2), 6) +
			bit.Set(bit.Get(r10, bit.B3), 7) + r5
		r6r8r9r10 := bit.Set(bit.Get((r10), bit.B4), 7) +
			bit.Set(bit.Get((r9), bit.B4), 6) +
			bit.Set(bit.Get((r8), bit.B4), 5) + r6

		// set new merged values
		n.Data[1][i] = r1r8
		n.Data[3][i] = r3r9
		n.Data[5][i] = r5r10
		n.Data[6][i] = r6r8r9r10

		// set to 0 useless values
		n.Data[8][i] = 0
		n.Data[9][i] = 0
		n.Data[10][i] = 0
	}

	n.NbMixBlock = y.NbMixBlock
	n.MixBlock = make([]ym.MixBlock, y.NbMixBlock)
	for i := range int(y.NbMixBlock) {
		n.MixBlock[i].NbRepeat = y.MixBlock[i].NbRepeat
		n.MixBlock[i].ReplayFreq = y.MixBlock[i].ReplayFreq
		n.MixBlock[i].SampleLength = y.MixBlock[i].SampleLength
		n.MixBlock[i].SampleStart = y.MixBlock[i].SampleStart
	}

	n.NbTimeKey = y.NbTimeKey
	n.TimeInfo = make([]ym.TimeKey, y.NbTimeKey)
	for i := range int(y.NbTimeKey) {
		n.TimeInfo[i].Time = y.TimeInfo[i].Time
		n.TimeInfo[i].NRepeat = y.TimeInfo[i].NRepeat
		n.TimeInfo[i].NBlock = y.TimeInfo[i].NBlock
	}
	n.MusicLenInMs = y.MusicLenInMs
	n.NbVoice = y.NbVoice
	n.EndID = y.EndID

	return n
}
