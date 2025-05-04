package xls_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jeromelesaux/ym/yeti/ui/xls"
	"github.com/stretchr/testify/assert"
)

func TestSaveExcel(t *testing.T) {
	filepath := "./test.xlsx"
	os.Remove(filepath)

	t.Run("ok_saving", func(t *testing.T) {
		d := data(10)
		x := xls.XlsFile{}

		assert.Nil(t, x.New(filepath, d))
	})

	t.Run("ok_parsing", func(t *testing.T) {
		d, err := xls.XlsFile{}.Get(filepath)
		assert.Nil(t, err)
		assert.Equal(t, 16, len(d))
		assert.Equal(t, 10, len(d[0]))
		fmt.Printf("%v", d)
	})

}

func data(length int) [16][]byte {
	var d [16][]byte

	for i := range 16 {
		d[i] = make([]byte, length)
	}

	for i := range length {
		for j := range 16 {
			d[j][i] = byte(i)
		}
	}
	return d
}

func TestIterateChar(t *testing.T) {
	for i := range 10 {
		fmt.Printf("%s", string(rune('A')+int32(i)))
	}
}
