package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/jeromelesaux/ym"
)

var (
	ErrorCheckstringDiffers = errors.New("Checkstring LeOnArD! differs")
	ErrorEndidDiffers       = errors.New("EndID End! differs")
	ErrorFileidDiffers      = errors.New("FileID YM6! differs")
)

func Unmarshall(data []byte, y *ym.Ym) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &y.FileID); err != nil {
		return err
	}

	if err := binary.Read(r, binary.BigEndian, &y.CheckString); err != nil {
		return err
	}
	if string(y.CheckString[:]) != "LeOnArD!" {
		return ErrorCheckstringDiffers
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
				return errors.New("size read differs from sample size")
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
			y.SongName = make([]byte, len(songName))
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
			y.AuthorName = make([]byte, len(authorName))
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
			y.SongComment = make([]byte, len(songComment))
			copy(y.SongComment, songComment)
			break
		}
	}

	for i := 0; i <= 15; i++ {
		y.Data[i] = make([]byte, y.NbFrames)
	}

	for j := 0; j < 16; j++ {
		for i := 0; i < int(y.NbFrames); i++ {
			v, err := r.ReadByte()
			y.Data[j][i] = v // writeRegister(v, j)
			if err != nil {
				return err
			}
		}
	}
	if err := binary.Read(r, binary.BigEndian, &y.EndID); err != nil {
		fmt.Fprintf(os.Stderr, "Warning no Endid found \n")
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
			if err := binary.Write(&b, binary.BigEndian, &y.Digidrums[i].SampleSize); err != nil {
				return b.Bytes(), err
			}
			if err := binary.Write(&b, binary.BigEndian, &y.Digidrums[i].SampleData); err != nil {
				return b.Bytes(), err
			}
		}
	}

	if err := binary.Write(&b, binary.BigEndian, &y.SongName); err != nil {
		return b.Bytes(), err
	}
	var eos byte = 0 // end of string (c compliant)
	if err := binary.Write(&b, binary.BigEndian, eos); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.AuthorName); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, eos); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, &y.SongComment); err != nil {
		return b.Bytes(), err
	}
	if err := binary.Write(&b, binary.BigEndian, eos); err != nil {
		return b.Bytes(), err
	}
	var err error

	for j := 0; j < 16; j++ {
		for i := 0; i < int(y.NbFrames); i++ {
			err = b.WriteByte(y.Data[j][i])
			fmt.Fprintf(os.Stderr, "j:%d,i:%d:%d\n", j, i, y.Data[j][i])
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

func writeRegister(v byte, index int) byte {
	switch index {
	case 0:
		return v & 255
	case 1:
		return v & 15
	case 2:
		return v & 255
	case 3:
		return v & 15
	case 4:
		return v & 255
	case 5:
		return v & 15
	case 6:
		return v & 0x1f
	case 7:
		return v & 255
	case 8:
		return v & 31
	case 9:
		return v & 31
	case 10:
		return v & 31
	case 11:
		return v & 255
	case 12:
		return v & 255
	case 13:
		return v & 0xF
	}
	return v
}
