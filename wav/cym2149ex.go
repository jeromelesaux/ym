package wav

import "fmt"

type CYm2149Ex struct {
	dcAdjust *CDcAdjuster

	frameCycle     uint32
	cyclePerSample uint32

	replayFrequency int32
	internalClock   uint32
	registers       [14]byte

	cycleSample uint32
	stepA       uint32
	stepB       uint32
	stepC       uint32
	posA        uint32
	posB        uint32
	posC        uint32
	volA        int32
	volB        int32
	volC        int32
	volE        int32
	mixerTA     uint32
	mixerTB     uint32
	mixerTC     uint32
	mixerNA     uint32
	mixerNB     uint32
	mixerNC     uint32
	pVolA       *int32
	pVolB       *int32
	pVolC       *int32

	noiseStep    uint32
	noisePos     uint32
	rndRack      uint32
	currentNoise uint32
	bWrite13     uint32

	envStep      uint32
	envPos       uint32
	envPhase     int32
	envShape     int32
	envData      [][][]byte //[16][2][16 * 2]byte
	globalVolume int32

	specialEffect   [3]YmSpecialEffect
	bSyncBuzzer     bool
	syncBuzzerStep  uint32
	syncBuzzerPhase uint32
	syncBuzzerShape int32

	lowPassFilter [2]int
	bFilter       bool
}

func NewCYm2149Ex(masterClock uint32, prediv int32, playRate uint32) *CYm2149Ex {
	c := &CYm2149Ex{
		bFilter:  true,
		dcAdjust: NewCDcAdjuster(),
	}
	c.envData = make([][][]byte, 16)
	for i := 0; i < 16; i++ {
		c.envData[i] = make([][]byte, 2)
		for j := 0; j < 2; j++ {
			c.envData[i][j] = make([]byte, 16*2)
		}

	}

	if ymVolumeTable[15] == 32767 {
		for i := 0; i < 16; i++ {
			ymVolumeTable[i] = (ymVolumeTable[i] * 2) / 6
		}
	}
	for env := 0; env < 16; env++ {
		pse := EnvWave[env]
		for phase := 0; phase < 2; phase++ {
			c.ym2149EnvInit(env, 0, phase, pse[phase*2+0], pse[phase*2+1])
		}
		for phase := 2; phase < 4; phase++ {
			c.ym2149EnvInit(env, 1, phase-2, pse[phase*2+0], pse[phase*2+1])
		}
	}
	c.internalClock = masterClock / uint32(prediv) // YM at 2Mhz on ATARI ST
	c.replayFrequency = int32(playRate)            // DAC at 44.1Khz on PC
	c.cycleSample = 0
	c.reset()
	return c
}

func (c *CYm2149Ex) reset() {
	for i := 0; i < 14; i++ {
		c.writeRegister(int32(i), 0)
	}
	c.writeRegister(7, 0xff)
	c.currentNoise = 0xffff
	c.rndRack = 1
	c.sidStop(0)
	c.sidStop(1)
	c.sidStop(2)
	c.envShape = 0
	c.envPhase = 0
	c.envPos = 0

	c.dcAdjust.Reset()

	c.syncBuzzerStop()

	c.lowPassFilter[0] = 0
	c.lowPassFilter[1] = 0
}

func (c *CYm2149Ex) syncBuzzerStop() {

	c.bSyncBuzzer = false
	c.syncBuzzerShape = 0
	c.syncBuzzerStep = 0
}

func (c *CYm2149Ex) sidStop(voice int32) {

	c.specialEffect[voice].bSid = false
}

func (c *CYm2149Ex) writeRegister(reg, data int32) {
	switch reg {
	case 0:
		c.registers[0] = byte(data & 255)
		c.stepA = c.toneStepCompute(c.registers[1], c.registers[0])
		if c.stepA == 0 {
			c.posA = (1 << 31) // Assume output always 1 if 0 period (for Digi-sample !)
		}
	case 2:
		c.registers[2] = byte(data & 255)
		c.stepB = c.toneStepCompute(c.registers[3], c.registers[2])
		if c.stepB == 0 {
			c.posB = (1 << 31)
		} // Assume output always 1 if 0 period (for Digi-sample !)

	case 4:
		c.registers[4] = byte(data & 255)
		c.stepC = c.toneStepCompute(c.registers[5], c.registers[4])
		if c.stepC == 0 {
			c.posC = (1 << 31) // Assume output always 1 if 0 period (for Digi-sample !)
		}

	case 1:
		c.registers[1] = byte(data & 15)
		c.stepA = c.toneStepCompute(c.registers[1], c.registers[0])
		if c.stepA == 0 {
			c.posA = (1 << 31)
		} // Assume output always 1 if 0 period (for Digi-sample !)
	case 3:
		c.registers[3] = byte(data & 15)
		c.stepB = c.toneStepCompute(c.registers[3], c.registers[2])
		if c.stepB == 0 {
			c.posB = (1 << 31)
		} // Assume output always 1 if 0 period (for Digi-sample !)

	case 5:
		c.registers[5] = byte(data & 15)
		c.stepC = c.toneStepCompute(c.registers[5], c.registers[4])
		if c.stepC == 0 {
			c.posC = (1 << 31)
		} // Assume output always 1 if 0 period (for Digi-sample !)

	case 6:
		c.registers[6] = byte(data & 0x1f)
		c.noiseStep = c.noiseStepCompute(c.registers[6])
		if c.noiseStep == 0 {
			c.noisePos = 0
			c.currentNoise = 0xffff
		}

	case 7:
		c.registers[7] = byte(data & 255)
		if data&(1<<0) != 0 {
			c.mixerTA = 0xffff
		}
		if data&(1<<1) != 0 {
			c.mixerTB = 0xffff
		}
		if data&(1<<2) != 0 {
			c.mixerTC = 0xffff
		}
		if data&(1<<3) != 0 {
			c.mixerNA = 0xffff
		}
		if data&(1<<4) != 0 {
			c.mixerNB = 0xffff
		}
		if data&(1<<5) != 0 {
			c.mixerNC = 0xffff
		}

	case 8:
		c.registers[8] = byte(data & 31)
		c.volA = ymVolumeTable[data&15]
		if (data & 0x10) != 0 {
			c.pVolA = &c.volE
		} else {
			c.pVolA = &c.volA
		}

	case 9:
		c.registers[9] = byte(data & 31)
		c.volB = ymVolumeTable[data&15]
		if (data & 0x10) != 0 {
			c.pVolB = &c.volE
		} else {
			c.pVolB = &c.volB
		}
	case 10:
		c.registers[10] = byte(data & 31)
		c.volC = ymVolumeTable[data&15]
		if (data & 0x10) != 0 {
			c.pVolC = &c.volE
		} else {
			c.pVolC = &c.volC
		}

	case 11:
		c.registers[11] = byte(data & 255)
		c.envStep = c.envStepCompute(c.registers[12], c.registers[11])

	case 12:
		c.registers[12] = byte(data & 255)
		c.envStep = c.envStepCompute(c.registers[12], c.registers[11])

	case 13:
		c.registers[13] = byte(data & 0xf)
		c.envPos = 0
		c.envPhase = 0
		c.envShape = data & 0xf

	default:
		fmt.Printf("this register %d does not exist.\n", reg)

	}
}

func (c *CYm2149Ex) toneStepCompute(rHigh, rLow byte) uint32 {
	var per uint32 = uint32(rHigh & 15)
	per = (per << 8) + uint32(rLow)
	if per <= 5 {
		return 0
	}
	var step int64 = int64(c.internalClock)
	step <<= (15 + 16 - 3)
	step /= (int64(per) * int64(c.replayFrequency))
	var istep uint32 = uint32(step)
	return istep
}

func (c *CYm2149Ex) envStepCompute(rHigh, rLow byte) uint32 {
	var per uint32 = uint32(rHigh & 15)
	per = (per << 8) + uint32(rLow)
	if per < 3 {
		return 0
	}
	var step int64 = int64(c.internalClock)
	step <<= (16 + 16 - 9)
	step /= int64(int64(per) * int64(c.replayFrequency))
	return uint32(step)
}

func (c *CYm2149Ex) noiseStepCompute(rNoise byte) uint32 {
	var per int32 = (int32(rNoise) & 0x1f)
	if per < 3 {
		return 0
	}

	step := int64(c.internalClock)
	step <<= (16 - 1 - 3)
	step /= (int64(per) * int64(c.replayFrequency))

	return uint32(step)
}

func (c *CYm2149Ex) sidStart(voice, timerFreq, vol int32) {
	tmp := int(timerFreq) * ((1 << 31) / int(c.replayFrequency))
	c.specialEffect[voice].sidStep = uint32(tmp)
	c.specialEffect[voice].sidVol = vol & 15
	c.specialEffect[voice].bSid = true
}

func (c *CYm2149Ex) drumStart(voice int32, pDrumBuffer []byte, drumSize uint32, drumFreq int32) {
	if (len(pDrumBuffer) > 0) && (drumSize != 0) {
		c.specialEffect[voice].drumData = pDrumBuffer
		c.specialEffect[voice].drumPos = 0
		c.specialEffect[voice].drumSize = drumSize
		c.specialEffect[voice].drumStep = uint32(drumFreq<<DRUM_PREC) / uint32(c.replayFrequency)
		c.specialEffect[voice].bDrum = true
	}
}

func (c *CYm2149Ex) sidSinStart(voice, timerFreq, vol int32) {
	// TODO
}

func (c *CYm2149Ex) readRegister(reg int32) int32 {
	if reg >= 0 && reg <= 13 {
		return int32(c.registers[reg])
	}
	return -1
}

func (c *CYm2149Ex) syncBuzzerStart(timerFreq, envShape int32) {
	var tmp uint32 = uint32(timerFreq) * ((1 << 31) / uint32(c.replayFrequency))
	c.envShape = envShape & 15
	c.syncBuzzerStep = uint32(tmp)
	c.syncBuzzerPhase = 0
	c.bSyncBuzzer = true
}

func (c *CYm2149Ex) update(pSampleBuffer *[]int16, pIndex int32, nbSample int32) {
	if nbSample > 0 {
		for {
			if nbSample == 142 {
				fmt.Printf("debug")
			}
			(*pSampleBuffer)[pIndex] = c.nextSample()
			pIndex++
			nbSample--
			if nbSample == 0 {
				break
			}
		}
	}
}

func (c *CYm2149Ex) rndCompute() uint32 {
	var rBit int32 = int32((c.rndRack)&1) ^ ((int32(c.rndRack) >> 2) & 1)
	c.rndRack = (c.rndRack >> 1) | (uint32(rBit) << 16)
	if rBit != 0 {
		return 0
	}
	return 0xffff
}

func (c *CYm2149Ex) sidVolumeCompute(voice int32, pVol *int32) {
	pVoice := c.specialEffect[voice]

	if pVoice.bSid {
		if pVoice.sidPos&(1<<31) != 0 {
			c.writeRegister(8+voice, pVoice.sidVol)
		} else {
			c.writeRegister(8+voice, 0)
		}
	} else {
		if pVoice.bDrum {
			//			writeRegister(8+voice,pVoice->drumData[pVoice->drumPos>>DRUM_PREC]>>4);

			*pVol = int32(pVoice.drumData[pVoice.drumPos>>DRUM_PREC]*255) / 6

			switch voice {
			case 0:
				*c.pVolA = c.volA
				c.mixerTA = 0xffff
				c.mixerNA = 0xffff

			case 1:
				*c.pVolB = c.volB
				c.mixerTB = 0xffff
				c.mixerNB = 0xffff

			case 2:
				*c.pVolC = c.volC
				c.mixerTC = 0xffff
				c.mixerNC = 0xffff

			}
		}
		pVoice.drumPos += pVoice.drumStep
		if (pVoice.drumPos >> DRUM_PREC) >= pVoice.drumSize {
			pVoice.bDrum = false
		}

	}
}

func (c *CYm2149Ex) nextSample() int16 {

	var vol int32
	var bt, bn int32

	if (c.noisePos & 0xffff0000) != 0 {
		c.currentNoise ^= c.rndCompute()
		c.noisePos &= 0xffff
	}
	bn = int32(c.currentNoise)

	c.volE = ymVolumeTable[c.envData[c.envShape][c.envPhase][c.envPos>>(32-5)]]

	c.sidVolumeCompute(0, &c.volA)
	c.sidVolumeCompute(1, &c.volB)
	c.sidVolumeCompute(2, &c.volC)

	//---------------------------------------------------
	// Tone+noise+env+DAC for three voices !
	//---------------------------------------------------
	bt = ((int32(c.posA) >> 31) | int32(c.mixerTA)) & (bn | int32(c.mixerNA))
	vol = int32(*c.pVolA) & bt
	bt = ((int32(c.posB) >> 31) | int32(c.mixerTB)) & (bn | int32(c.mixerNB))
	vol += int32(*c.pVolB) & bt
	bt = ((int32(c.posC) >> 31) | int32(c.mixerTC)) & (bn | int32(c.mixerNC))
	vol += int32(*c.pVolC) & bt

	//---------------------------------------------------
	// Inc
	//---------------------------------------------------
	c.posA += c.stepA
	c.posB += c.stepB
	c.posC += c.stepC
	c.noisePos += c.noiseStep
	c.envPos += c.envStep
	if c.envPhase == 0 {
		if c.envPos < c.envStep {
			c.envPhase = 1
		}
	}

	c.syncBuzzerPhase += c.syncBuzzerStep
	if c.syncBuzzerPhase&(1<<31) != 0 {
		c.envPos = 0
		c.envPhase = 0
		c.syncBuzzerPhase &= 0x7fffffff
	}

	c.specialEffect[0].sidPos += c.specialEffect[0].sidStep
	c.specialEffect[1].sidPos += c.specialEffect[1].sidStep
	c.specialEffect[2].sidPos += c.specialEffect[2].sidStep

	//---------------------------------------------------
	// Normalize process
	//---------------------------------------------------
	c.dcAdjust.AddSample(vol)
	in := int(vol) - int(c.dcAdjust.GetDcLevel())
	if c.bFilter {
		return int16(c.LowPassFilter(in))
	}
	return int16(in)
}

func (c *CYm2149Ex) LowPassFilter(in int) int {
	out := (c.lowPassFilter[0] >> 2) + (c.lowPassFilter[1] >> 1) + (in >> 2)
	c.lowPassFilter[0] = c.lowPassFilter[1]
	c.lowPassFilter[1] = in
	return out
}

func (c *CYm2149Ex) ym2149EnvInit(index0, index1, indexStart int, a, b int32) byte {
	d := b - a
	a *= 15
	var pEnvIndex int = 16 * indexStart
	for i := 0; i < 16; i++ {
		c.envData[index0][index1][pEnvIndex] = byte(a)
		//(*pEnv)[pEnvIndex] = byte(a)
		pEnvIndex++
		a += d
	}
	return c.envData[index0][index1][pEnvIndex-1]
}
