package wav

const (
	DC_ADJUST_BUFFERLEN = 512
	DRUM_PREC           = 15
	NOISESIZE           = 16384
	MFP_CLOCK           = 2457600
	//---------------------------------------------------------------------
	// To produce a WAV file.
	//---------------------------------------------------------------------
	ID_RIFF           = 0x46464952
	ID_WAVE           = 0x45564157
	ID_FMT            = 0x20746D66
	ID_DATA           = 0x61746164
	NBSAMPLEPERBUFFER = 1024
	MAX_VOICE         = 8
	YMTPREC           = 16
)

var (
	ymVolumeTable = [16]int32{62, 161, 265, 377, 580, 774, 1155, 1575, 2260, 3088, 4570, 6233, 9330, 13187, 21220, 32767}
	Env00xx       = [8]int32{1, 0, 0, 0, 0, 0, 0, 0}
	Env01xx       = [8]int32{0, 1, 0, 0, 0, 0, 0, 0}
	Env1000       = [8]int32{1, 0, 1, 0, 1, 0, 1, 0}
	Env1001       = [8]int32{1, 0, 0, 0, 0, 0, 0, 0}
	Env1010       = [8]int32{1, 0, 0, 1, 1, 0, 0, 1}
	Env1011       = [8]int32{1, 0, 1, 1, 1, 1, 1, 1}
	Env1100       = [8]int32{0, 1, 0, 1, 0, 1, 0, 1}
	Env1101       = [8]int32{0, 1, 1, 1, 1, 1, 1, 1}
	Env1110       = [8]int32{0, 1, 1, 0, 0, 1, 1, 0}
	Env1111       = [8]int32{0, 1, 0, 0, 0, 0, 0, 0}
	EnvWave       = [16][8]int32{
		Env00xx, Env00xx, Env00xx, Env00xx,
		Env01xx, Env01xx, Env01xx, Env01xx,
		Env1000, Env1001, Env1010, Env1011,
		Env1100, Env1101, Env1110, Env1111}
	// ATARI-ST MFP chip predivisor
	mfpPrediv = [8]byte{0, 4, 10, 16, 50, 64, 100, 200}
)

const (
	YM_V2 ymFileType = iota
	YM_V3
	YM_V4
	YM_V5
	YM_V6
	YM_VMAX
)

const (
	YM_TRACKER1   ymFileType = 32
	YM_TRACKER2   ymFileType = 33
	YM_TRACKERMAX ymFileType = 34
)

const (
	YM_MIX1   ymFileType = 64
	YM_MIX2   ymFileType = 65
	YM_MIXMAX ymFileType = 66
)
