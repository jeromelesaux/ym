package wav

type CDcAdjuster struct {
	buffer []int32
	pos    int32
	sum    int32
}

func NewCDcAdjuster() *CDcAdjuster {
	return &CDcAdjuster{
		buffer: make([]int32, DC_ADJUST_BUFFERLEN),
	}
}

func (c *CDcAdjuster) AddSample(sample int32) {
	c.sum -= c.buffer[c.pos]
	c.sum += sample

	c.buffer[c.pos] = sample
	c.pos = (c.pos + 1) & (DC_ADJUST_BUFFERLEN - 1)

}
func (c *CDcAdjuster) Reset() {
	for i := range DC_ADJUST_BUFFERLEN {
		c.buffer[i] = 0
	}
	c.pos = 0
	c.sum = 0
}

func (c *CDcAdjuster) GetDcLevel() int32 {
	return c.sum / DC_ADJUST_BUFFERLEN
}
