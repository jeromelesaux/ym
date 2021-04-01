package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/jeromelesaux/lha"
	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/encoding"
)

type YmSpecialEffect struct {
	bDrum    bool
	drumSize uint32
	drumData []byte
	drumPos  uint32
	drumStep uint32

	bSid    bool
	sidPos  uint32
	sidStep uint32
	sidVol  int32
}

type MixBlock struct {
	SampleStart  uint32
	SampleLength uint32
	NbRepeat     uint16
	ReplayFreq   uint16
}

type TimeKey struct {
	time    uint32
	nRepeat uint16
	nBlock  uint16
}

type ymTrackerVoice struct {
	pSample      []byte
	sampleSize   uint32
	samplePos    uint32
	repLen       uint32
	sampleVolume int32
	sampleFreq   uint32
	bLoop        bool
	bRunning     bool
}

type ymTrackerLine struct {
	noteOn   byte
	volume   byte
	freqHigh byte
	freqLow  byte
}

type ymFileType int

type YMMusic struct {
	SongName         []byte
	SongAuthor       []byte
	SongComment      []byte
	SongType         ymFileType
	SongPlayer       []byte
	MusicTimeInSec   int32
	MusicTimeInMs    int32
	bLoop            bool
	replayRate       int
	playerRate       int
	pBigSampleBuffer []byte
	pMixBlock        []MixBlock
	innerSamplePos   int
	currentPos       uint32
	nbDrum           int
	pDrumTab         []ym.Digidrum
	ymChip           *CYm2149Ex
	currentFrame     int
	loopFrame        int
	bMusicOver       bool
	//streamInc        int

	pDataStream [][]byte // structure deinterleave frame by frame
	bMusicOk    bool
	bPause      bool
	//nbTimerKey              int32
	pTimeInfo               *TimeKey
	musicLenInMs            uint32
	iMusicPosAccurateSample uint32
	iMusicPosInMs           uint32
	mixPos                  int32
	nbRepeat                int32
	pCurrentMixSample       []byte
	currentSampleLength     uint32
	currentPente            uint32
	nbMixBlock              int32
	ymTrackerNbSampleBefore int
	ymTrackerVoice          []ymTrackerVoice
	nbVoice                 int
	nbFrame                 int
	ymTrackerFreqShift      int
	ymTrackerVolumeTable    []byte
}

type WAVEHeader struct {
	RIFFMagic     uint32
	FileLength    uint32
	FileType      uint32
	FormMagic     uint32
	FormLength    uint32
	SampleFormat  uint16
	NumChannels   uint16
	PlayRate      uint32
	BytesPerSec   uint32
	Pad           uint16
	BitsPerSample uint16
	DataMagic     uint32
	DataLength    uint32
}

func NewYMMusic() *YMMusic {
	return &YMMusic{
		SongName:             make([]byte, 0),
		SongAuthor:           make([]byte, 0),
		SongComment:          make([]byte, 0),
		SongPlayer:           make([]byte, 0),
		pBigSampleBuffer:     make([]byte, 0),
		pDataStream:          make([][]byte, 16),
		pMixBlock:            make([]MixBlock, 0),
		pDrumTab:             make([]ym.Digidrum, 0),
		ymTrackerVolumeTable: make([]byte, 256*64),
		playerRate:           50,
		replayRate:           44100,
		pTimeInfo:            &TimeKey{},
		ymTrackerVoice:       make([]ymTrackerVoice, MAX_VOICE),
		ymChip:               NewCYm2149Ex(encoding.ATARI_CLOCK, 1, 44100),
	}
}

func (y *YMMusic) trackerInit(volMaxPercent int32) {
	for i := 0; i < MAX_VOICE; i++ {
		y.ymTrackerVoice[i].bRunning = false
	}
	y.ymTrackerNbSampleBefore = 0
	scale := (256 * volMaxPercent) / int32(y.nbVoice*100)
	// Construit la table de volume.
	index := 0
	for vol := 0; vol < 64; vol++ {
		for s := -128; s < 128; s++ {
			y.ymTrackerVolumeTable[index] = byte((s * int(scale) * vol) / 64)
			index++
		}
	}
}

func (y *YMMusic) LoadMemory(v *ym.Ym) error {
	if v == nil {
		return fmt.Errorf("YM is nil")
	}
	y.SongName = make([]byte, len(v.SongName))
	copy(y.SongName, v.SongName)
	y.SongAuthor = make([]byte, len(v.AuthorName))
	copy(y.SongAuthor, v.AuthorName)
	y.SongComment = make([]byte, len(v.SongComment))
	copy(y.SongComment, v.SongComment)
	y.MusicTimeInMs = getMusicTime(v)
	y.MusicTimeInSec = getMusicTime(v) / 1000
	y.nbFrame = int(v.NbFrames)
	y.nbDrum = int(v.DigidrumNb)
	y.pDrumTab = make([]ym.Digidrum, v.DigidrumNb)
	for i := 0; i < y.nbDrum; i++ {
		copy(y.pDrumTab[i].SampleData, v.Digidrums[i].SampleData)
		y.pDrumTab[i].SampleSize = v.Digidrums[i].SampleSize
	}
	switch v.FileID {
	case ym.YM2:
		y.SongType = YM_V2
	case ym.YM3:
		y.SongType = YM_V3
	case ym.YM4:
		y.SongType = YM_V4
	case ym.YM5:
		y.SongType = YM_V5
	case ym.YM6:
		y.SongType = YM_V6
	case ym.YMT1:
		y.SongType = YM_MIX1
	//	y.trackerInit(100)
	case ym.YMT2:
		y.SongType = YM_MIX2
		//	y.trackerInit(100)
	}

	for i := 0; i < 16; i++ {
		y.pDataStream[i] = make([]byte, y.nbFrame)
		for j := 0; j < y.nbFrame; j++ {
			y.pDataStream[i][j] = v.Data[i][j]
		}
	}

	y.ymChip.reset()
	y.MusicTimeInMs = getMusicTime(v)
	y.MusicTimeInSec = getMusicTime(v) / 1000
	y.bMusicOk = true
	y.bPause = false

	return nil
}

func (y *YMMusic) Load(filePath string) error {
	archive := lha.NewLha(filePath)
	headers, err := archive.Headers()
	if err != nil {
		return err
	}

	if len(headers) == 0 {
		return fmt.Errorf("no headers found in archive " + filePath)

	}
	content, err := archive.DecompresBytes(headers[0])
	if err != nil {
		return err
	}
	v := ym.NewYm()
	err = encoding.Unmarshall(content, v)
	if err != nil && err != io.EOF {
		return err
	}

	return y.LoadMemory(v)
}

func getMusicTime(v *ym.Ym) int32 {
	if v.NbFrames > 0 && v.FrameHz > 0 {
		return int32(v.NbFrames) * 1000 / int32(v.FrameHz)
	}
	return 0
}

func (y *YMMusic) Wave() ([]byte, error) {
	head := &WAVEHeader{}

	wavWriter := &bytes.Buffer{}

	if err := binary.Write(wavWriter, binary.LittleEndian, head); err != nil {
		return wavWriter.Bytes(), err
	}
	y.bLoop = false
	totalNbSample := 0
	convertBuffer := make([]int16, 1024)

	for {
		ok := y.MusicCompute(&convertBuffer, NBSAMPLEPERBUFFER)

		if err := binary.Write(wavWriter, binary.LittleEndian, convertBuffer); err != nil {
			return wavWriter.Bytes(), err
		}
		totalNbSample += NBSAMPLEPERBUFFER

		if !ok {
			break
		}
	}

	finalHeader := &bytes.Buffer{}
	head.RIFFMagic = ID_RIFF
	head.FileType = ID_WAVE
	head.FormMagic = ID_FMT
	head.DataMagic = ID_DATA
	head.FormLength = 0x10
	head.SampleFormat = 1
	head.NumChannels = 1
	head.PlayRate = 44100
	head.BitsPerSample = 16
	head.BytesPerSec = 44100 * (16 / 8)
	head.Pad = (16 / 8)
	head.DataLength = uint32(totalNbSample) * (16 / 8)
	head.FileLength = head.DataLength + 44 - 8 // 44 sizeof waveheader

	if err := binary.Write(finalHeader, binary.LittleEndian, head); err != nil {
		return wavWriter.Bytes(), err
	}
	wavContent := wavWriter.Bytes()
	copy(wavContent[0:], finalHeader.Bytes())
	return wavContent, nil
}

func (y *YMMusic) WaveFile(wavFilepath string) error {
	head := &WAVEHeader{}
	fw, err := os.Create(wavFilepath)
	if err != nil {
		return err
	}
	defer fw.Close()
	if err := binary.Write(fw, binary.LittleEndian, head); err != nil {
		return err
	}
	y.bLoop = false
	totalNbSample := 0
	nbTotal := y.MusicTimeInSec * 44100
	oldRatio := -1
	convertBuffer := make([]int16, 1024)

	for {
		ok := y.MusicCompute(&convertBuffer, NBSAMPLEPERBUFFER)

		if err := binary.Write(fw, binary.LittleEndian, convertBuffer); err != nil {
			return err
		}
		totalNbSample += NBSAMPLEPERBUFFER
		ratio := (totalNbSample * 100) / int(nbTotal)
		if ratio != oldRatio {
			fmt.Printf("Rendering... (%d%%)\r", ratio)
			oldRatio = ratio
		}
		if !ok {
			break
		}
	}
	fmt.Printf("\n")
	fw.Seek(0, io.SeekStart)
	head.RIFFMagic = ID_RIFF
	head.FileType = ID_WAVE
	head.FormMagic = ID_FMT
	head.DataMagic = ID_DATA
	head.FormLength = 0x10
	head.SampleFormat = 1
	head.NumChannels = 1
	head.PlayRate = 44100
	head.BitsPerSample = 16
	head.BytesPerSec = 44100 * (16 / 8)
	head.Pad = (16 / 8)
	head.DataLength = uint32(totalNbSample) * (16 / 8)
	head.FileLength = head.DataLength + 44 - 8 // 44 sizeof waveheader
	if err := binary.Write(fw, binary.LittleEndian, head); err != nil {
		return err
	}
	return nil
}

func (y *YMMusic) MusicCompute(buffer *[]int16, nbSample int) bool {
	clearBuffer := make([]int16, nbSample)
	copy(*buffer, clearBuffer)
	var nbs int = nbSample
	vblNbSample := y.replayRate / y.playerRate
	var sampleToCompute int
	var pOut int
	if !y.bMusicOk || y.bPause || y.bMusicOver {
		if y.bMusicOver {
			return false
		}
	}

	if y.SongType >= YM_MIX1 && y.SongType < YM_MIXMAX {
		y.stDigitMix(buffer, int32(nbSample))
	} else {
		if y.SongType >= YM_TRACKER1 && y.SongType < YM_TRACKERMAX {
			//to implement	y.TrackerUpdate(buffer, int32(nbSample))
		} else {

			for {
				// Nb de sample � calculer avant l'appel de Player
				sampleToCompute = vblNbSample - y.innerSamplePos
				// Test si la fin du buffer arrive avant la fin de sampleToCompute
				if sampleToCompute > nbs {
					sampleToCompute = nbs
				}
				y.innerSamplePos += sampleToCompute
				if y.innerSamplePos >= vblNbSample {
					y.player() // Lecture de la partition (playerRate Hz)
					y.innerSamplePos -= vblNbSample
				}
				if sampleToCompute > 0 {
					y.ymChip.update(buffer, int32(pOut), int32(sampleToCompute)) // YM Emulation.
					pOut += sampleToCompute
				}
				nbs -= sampleToCompute
				if nbs <= 0 {
					break
				}
			}
		}
	}
	return true
}

/*
func (y *YMMusic) ymTrackerPlayer(pVoice []ymTrackerVoice) {

	pLine := &ymTrackerLine{}
	pDataStream[y.currentFrame][y.nbVoice]
	binary.Read(y.pDataStreamReader, binary.BigEndian, pline)
	//	pLine = (ymTrackerLine_t*)pDataStream;
	pLine += (y.currentFrame * y.nbVoice)
	for i := 0; i < y.nbVoice; i++ {
		var n int32
		pVoice[i].sampleFreq = (uint32(pLine.freqHigh) << 8) | uint32(pLine.freqLow)
		if pVoice[i].sampleFreq != 0 {
			pVoice[i].sampleVolume = int32(pLine.volume) & 63
			pVoice[i].bLoop = false
			if (pLine.volume & 0x40) != 0 {
				pVoice[i].bLoop = true
			}
			n = int32(pLine.noteOn)
			if n != 0xff { // Note ON.

				pVoice[i].bRunning = true
				copy(pVoice[i].pSample, y.pDrumTab[n].SampleData)
				pVoice[i].sampleSize = y.pDrumTab[n].SampleSize
				//pVoice[i].repLen = y.pDrumTab[n].
				pVoice[i].samplePos = 0
			}
		} else {
			pVoice[i].bRunning = false
		}
		pLine++
	}

	y.currentFrame++
	if y.currentFrame >= y.nbFrame {
		if !y.bLoop {
			y.bMusicOver = true
		}
		y.currentFrame = 0
	}
}*/

/*
func (y *YMMusic) ymTrackerVoiceAdd(pVoice *ymTrackerVoice, pBuffer *[]byte, nbs int32) {
	var pVolumeTab *ymsample
	var pSample []byte
	var samplePos uint32
	var sampleEnd uint32
	var sampleInc uint32
	var repLen uint32
	var step float64

	if !(pVoice.bRunning) {
		return
	}

	pVolumeTab = &y.ymTrackerVolumeTable[256*(pVoice.sampleVolume&63)]
	pSample = pVoice.pSample
	samplePos = pVoice.samplePos

	step = float64(pVoice.sampleFreq << YMTPREC)
	step *= float64(int64(1 << y.ymTrackerFreqShift))
	step /= float64(y.replayRate)
	sampleInc = uint32(step)

	sampleEnd = (pVoice.sampleSize << YMTPREC)
	repLen = (pVoice.repLen << YMTPREC)
	if nbs > 0 {
		for {
			var va int32 = pVolumeTab[pSample[samplePos>>YMTPREC]]

			var vb int32 = va
			if samplePos < (sampleEnd - (1 << YMTPREC)) {
				vb = pVolumeTab[pSample[(samplePos>>YMTPREC)+1]]
			}
			var frac int32 = int32(samplePos) & ((1 << YMTPREC) - 1)
			va += (((vb - va) * frac) >> YMTPREC)

			//(*pBuffer++) += va;

			samplePos += sampleInc
			if samplePos >= sampleEnd {
				if pVoice.bLoop {
					samplePos -= repLen
				} else {
					pVoice.bRunning = false
					return
				}
			}
			nbs--
			if nbs <= 0 {
				break
			}
		}

	}
	pVoice.samplePos = samplePos
}
*/
/*
func (y *YMMusic) TrackerUpdate(pBuffer *[]int16, nbSample int32) {
	var nbs int32
	for i := 0; i < int(nbSample); i++ {
		(*pBuffer)[i] = 0
	}
	if y.bMusicOver {
		return
	}
	for {
		if y.ymTrackerNbSampleBefore == 0 {
			// Lit la partition ymTracker
			y.ymTrackerPlayer(y.ymTrackerVoice)
			if y.bMusicOver {
				return
			}
			y.ymTrackerNbSampleBefore = y.replayRate / y.playerRate
		}
		nbs = int32(y.ymTrackerNbSampleBefore) // nb avant playerUpdate.
		if nbs > nbSample {
			nbs = nbSample
		}
		y.ymTrackerNbSampleBefore -= int(nbs)
		if nbs > 0 {
			// Genere les samples.
			for i := 0; i < y.nbVoice; i++ {
				y.ymTrackerVoiceAdd(&y.ymTrackerVoice[i], pBuffer, nbs)
			}
			pBuffer += nbs
			nbSample -= nbs
		}
		if nbSample <= 0 {
			break
		}
	}
}
*/

func (y *YMMusic) stDigitMix(pWrite16 *[]int16, nbs int32) {
	if y.bMusicOver {
		return
	}

	if y.mixPos == -1 {
		y.nbRepeat = -1
		y.readNextBlockInfo()
	}

	y.iMusicPosAccurateSample += uint32(nbs) * 1000
	y.iMusicPosInMs += ((y.iMusicPosAccurateSample) / uint32(y.replayRate))
	y.iMusicPosAccurateSample %= uint32(y.replayRate)

	if nbs != 0 {
		for {

			var sa int32 = int32(y.pCurrentMixSample[y.currentPos>>12]) << 8
			var sb int32 = sa
			if (y.currentPos >> 12) < ((y.currentSampleLength >> 12) - 1) {
				sb = int32(y.pCurrentMixSample[(y.currentPos>>12)+1]) << 8
			}
			var frac int32 = int32(y.currentPos) & ((1 << 12) - 1)
			sa += (((sb - sa) * frac) >> 12)
			*pWrite16 = append(*pWrite16, int16(sa))
			//*pWrite16++ = sa;

			y.currentPos += y.currentPente
			if y.currentPos >= y.currentSampleLength {
				y.readNextBlockInfo()
				if y.bMusicOver {
					return
				}
			}
			nbs--
			if nbs <= 0 {
				break
			}
		}
	}
}

func (y *YMMusic) readNextBlockInfo() {
	y.nbRepeat--
	if y.nbRepeat <= 0 {
		y.mixPos++
		if y.mixPos >= y.nbMixBlock {
			y.mixPos = 0
			if !y.bLoop {
				y.bMusicOver = true
			}

			y.iMusicPosAccurateSample = 0
			y.iMusicPosInMs = 0
		}
		y.nbRepeat = int32(y.pMixBlock[y.mixPos].NbRepeat)
	}
	copy(y.pCurrentMixSample, y.pBigSampleBuffer[y.pMixBlock[y.mixPos].SampleStart:])
	y.currentSampleLength = (y.pMixBlock[y.mixPos].SampleLength) << 12
	y.currentPente = (uint32(y.pMixBlock[y.mixPos].ReplayFreq) << 12) / uint32(y.replayRate)
	y.currentPos &= ((1 << 12) - 1)
}

func (y *YMMusic) player() {
	var (
		ptr    int32
		prediv uint32
		voice  int32
		ndrum  int32
	)
	if y.currentFrame < 0 {
		y.currentFrame = 0
	}
	if y.currentFrame >= y.nbFrame {
		if y.bLoop {
			y.currentFrame = y.loopFrame
		} else {
			y.bMusicOver = true
			y.ymChip.reset()
			return
		}
	}
	ptr = int32(y.currentFrame)
	for i := 0; i <= 10; i++ {
		y.ymChip.writeRegister(int32(i), int32(y.pDataStream[int32(i)][ptr]))
	}

	y.ymChip.sidStop(0)
	y.ymChip.sidStop(1)
	y.ymChip.sidStop(2)
	y.ymChip.syncBuzzerStop()

	//---------------------------------------------
	// Check digi-drum
	//---------------------------------------------
	if y.SongType == YM_V2 { // MADMAX specific !

		if y.pDataStream[13][ptr] != 0xff {

			y.ymChip.writeRegister(11, int32(y.pDataStream[11][ptr]))
			y.ymChip.writeRegister(12, 0)
			y.ymChip.writeRegister(13, 10) // MADMAX specific !!
		}
		if y.pDataStream[10][ptr]&0x80 != 0 { // bit 7 volume canal C pour annoncer une digi-drum madmax.

			var sampleNum int32
			var sampleFrq uint32
			y.ymChip.writeRegister(7, y.ymChip.readRegister(7)|0x24) // Coupe TONE + NOISE canal C.
			sampleNum = int32(y.pDataStream[10][ptr] & 0x7f)         // Numero du sample

			if y.pDataStream[12][ptr] != 0 {
				sampleFrq = uint32(MFP_CLOCK / int(y.pDataStream[12][ptr]))
				y.ymChip.drumStart(2, // Voice C
					sampleAdress[sampleNum],
					sampleLen[sampleNum],
					int32(sampleFrq))
			}
		}
	} else {
		if y.SongType >= YM_V3 {

			y.ymChip.writeRegister(11, int32(y.pDataStream[11][ptr]))
			y.ymChip.writeRegister(12, int32(y.pDataStream[12][ptr]))
			if int32(y.pDataStream[13][ptr]) != 0xff {
				y.ymChip.writeRegister(13, int32(y.pDataStream[13][ptr]))
			}
		}
		if y.SongType >= YM_V5 {
			var code int32

			if y.SongType == YM_V6 {
				y.readYm6Effect(y.pDataStream, ptr, 1, 6, 14)
				y.readYm6Effect(y.pDataStream, ptr, 3, 8, 15)
			} else { // YM5 effect decoding

				//------------------------------------------------------
				// Sid Voice !!
				//------------------------------------------------------
				code = (int32(y.pDataStream[1][ptr]) >> 4) & 3
				if code != 0 {
					var tmpFreq uint32
					voice = code - 1
					prediv = uint32(mfpPrediv[(y.pDataStream[6][ptr]>>5)&7])
					prediv *= uint32(y.pDataStream[14][ptr])
					tmpFreq = 0
					if prediv != 0 {
						tmpFreq = 2457600 / prediv
						y.ymChip.sidStart(voice, int32(tmpFreq), int32(y.pDataStream[voice+8][ptr]&15))
					}
				}

				//------------------------------------------------------
				// YM5 Digi Drum.
				//------------------------------------------------------
				code = (int32(y.pDataStream[3][ptr]) >> 4) & 3
				if code != 0 { // Ici un digidrum demarre sur la voie voice.
					voice = code - 1
					ndrum = int32(y.pDataStream[8+voice][ptr]) & 31
					if (ndrum >= 0) && (int(ndrum) < y.nbDrum) {
						var sampleFrq uint32
						prediv = uint32(mfpPrediv[(y.pDataStream[8][ptr]>>5)&7])
						prediv *= uint32(y.pDataStream[15][ptr])
						if prediv != 0 {
							sampleFrq = MFP_CLOCK / prediv
							y.ymChip.drumStart(voice, y.pDrumTab[ndrum].SampleData, y.pDrumTab[ndrum].SampleSize, int32(sampleFrq))
						}
					}
				}
			}
		}
	}
	y.currentFrame++
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

func (y *YMMusic) readYm6Effect(pReg [][]byte, ptr int32, code, prediv, count int32) {
	code = int32(pReg[code][ptr]) & 0xf0
	prediv = int32(pReg[prediv][ptr]>>5) & 7
	count = int32(pReg[count][ptr])
	var voice int32
	var ndrum int32
	if code&0x30 != 0 {
		var tmpFreq uint32
		// Ici il y a un effet sur la voie:

		voice = ((code & 0x30) >> 4) - 1
		switch code & 0xc0 {
		case 0x00: // SID
			prediv = int32(mfpPrediv[prediv])
			prediv *= count
			tmpFreq = 0
			if prediv != 0 {
				tmpFreq = uint32(2457600 / prediv)
				if (code & 0xc0) == 0x00 {
					y.ymChip.sidStart(voice, int32(tmpFreq), int32(pReg[voice+8][ptr]&15))
				} else {
					y.ymChip.sidSinStart(voice, int32(tmpFreq), int32(pReg[voice+8][ptr]&15))
				}
			}
		case 0x80: // Sinus-SID

			prediv = int32(mfpPrediv[prediv])
			prediv *= count
			tmpFreq = 0
			if prediv != 0 {
				tmpFreq = uint32(2457600 / prediv)
				if (code & 0xc0) == 0x00 {
					y.ymChip.sidStart(voice, int32(tmpFreq), int32(pReg[voice+8][ptr]&15))
				} else {
					y.ymChip.sidSinStart(voice, int32(tmpFreq), int32(pReg[voice+8][ptr]&15))
				}
			}

		case 0x40: // DigiDrum
			ndrum = int32(pReg[voice+8][ptr]) & 31
			if (ndrum >= 0) && (int(ndrum) < y.nbDrum) {
				prediv = int32(mfpPrediv[prediv])
				prediv *= count
				if prediv > 0 {
					tmpFreq = uint32(2457600 / prediv)
					y.ymChip.drumStart(voice, y.pDrumTab[ndrum].SampleData, y.pDrumTab[ndrum].SampleSize, int32(tmpFreq))
				}
			}

		case 0xc0: // Sync-Buzzer.

			prediv = int32(mfpPrediv[prediv])
			prediv *= count
			tmpFreq = 0
			if prediv != 0 {
				tmpFreq = uint32(2457600 / prediv)
				y.ymChip.syncBuzzerStart(int32(tmpFreq), int32(pReg[voice+8][ptr]&15))
			}

		}

	}
}
