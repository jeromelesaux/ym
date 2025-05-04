package xls

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type XlsFile struct {
}

var (
	sheet = "Sheet1"
)

func (x XlsFile) New(filepath string, d [16][]byte) error {

	// new file created
	f := excelize.NewFile()
	defer f.Close()

	// set headers
	for i := range len(d) {
		if err := f.SetCellValue(
			sheet,
			fmt.Sprintf("%s%d", string(rune('A')+int32(i)), 1),
			fmt.Sprintf("Register %d", i),
		); err != nil {
			return err
		}
	}

	// set values

	for i := range d[0] {
		index := rune('A')
		for j := range 16 {
			v := int(d[j][i])
			if err := f.SetCellValue(
				sheet,
				fmt.Sprintf("%s%d", string(index+int32(j)), i+2),
				fmt.Sprintf("%X", v),
			); err != nil {
				return err
			}
		}
	}

	return f.SaveAs(filepath)

}

func (x XlsFile) Get(filepath string) ([16][]byte, error) {
	data := [16][]byte{}
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return data, err
	}

	raw, err := f.GetRows(sheet)
	if err != nil {
		return data, err
	}

	// skip header
	for i := 1; i < len(raw); i++ {
		values := raw[i]
		for j := 0; j < len(values) || j < 16; j++ {
			v, err := strconv.ParseUint(values[j], 16, 32)
			if err != nil {
				log.Printf("error in row [%d], col [%d], value not integer [%s]", i, j, values[j])
				return data, err
			}
			data[j] = append(data[j], byte(v))
		}
	}

	return data, nil
}
