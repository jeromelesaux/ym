package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/jeromelesaux/ym"
)

func Unmarshall(data []byte, y *ym.Ym) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &y.FileID); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.CheckString); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.NbFrames); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.SongAttributes); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.DigidrumNb); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.YmMasterClock); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.FrameHz); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.LoopFrame); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &y.Size); err != nil {
		return err
	}
	if y.DigidrumNb > 0 {
		y.Digidrums = make([]ym.Digidrum, y.DigidrumNb)
		for i := 0; i < int(y.DigidrumNb); i++ {
			d := ym.Digidrum{}

			if err := binary.Read(r, binary.BigEndian, &d.SampleSize); err != nil {
				return err
			}
			d.SampleData = make([]byte, d.SampleSize)
			n, err := r.Read(d.SampleData)
			if err != nil {
				return err
			}
			if n != int(d.SampleSize) {
				return errors.New("size read differs from sample size.")
			}
			y.Digidrums[i] = d

		}
	}
	songName := []byte{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		songName = append(songName, b)
		if b == 0 {
			copy(y.SongName, songName)
			break
		}
	}
	authorName := []byte{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		authorName = append(authorName, b)
		if b == 0 {
			copy(y.AuthorName, authorName)
			break
		}
	}
	songComment := []byte{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		songComment = append(songComment, b)
		if b == 0 {
			copy(y.SongComment, songComment)
			break
		}
	}

	for i := 0; i <= 15; i++ {
		y.Data[i] = make([]byte, y.NbFrames)
	}
	var err error
	for i := 0; i < int(y.NbFrames); i++ {
		for j := 0; j < 16; j++ {
			y.Data[j][i], err = r.ReadByte()
			if err != nil {
				return err
			}
		}
	}
	if err := binary.Read(r, binary.BigEndian, &y.EndID); err != nil {
		return err
	}
	return nil
}

func Marshall(y *ym.Ym) ([]byte, error) {
	var b bytes.Buffer
	if err := binary.Write(&b, binary.BigEndian, &y.FileID); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.CheckString); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.NbFrames); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.SongAttributes); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.DigidrumNb); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.YmMasterClock); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.FrameHz); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.LoopFrame); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.Size); err != nil {
		return b.Bytes(), err
	}
	if y.DigidrumNb > 0 {
		for i := 0; i < int(y.DigidrumNb); i++ {
			if err := binary.Write(&b, binary.BigEndian, &y.Digidrums[i]); err != nil {
				return b.Bytes(), err
			}
		}
	}

	if err := binary.Write(&b, binary.BigEndian, &y.SongName); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.AuthorName); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.SongComment); err != nil {
		return b.Bytes(), err
	}
	var err error
	for i := 0; i < int(y.NbFrames); i++ {
		for j := 0; j < 16; j++ {
			err = b.WriteByte(y.Data[j][i])
			if err != nil {
				return b.Bytes(), err
			}
		}
	}
	if err := binary.Write(&b, binary.BigEndian, &y.EndID); err != nil {
		return b.Bytes(), err
	}
	return b.Bytes(), err
}
