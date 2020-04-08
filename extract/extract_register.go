package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jeromelesaux/ym"
	"github.com/jeromelesaux/ym/encoding"
)

var ErrorIsNotDirectory = errors.New("Is not a directory, Quiting.")

var (
	out  = flag.String("out", "", "folder to save register")
	file = flag.String("ym", "", "ym filepath")
)

func main() {
	flag.Parse()
	if *out == "" || *file == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}
	f, err := os.Open(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open file (%s) error :%v\n", *file, err)
		os.Exit(-1)
	}
	defer f.Close()

	d, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read file (%s) error :%v\n", *file, err)
		os.Exit(-1)
	}
	CheckOutput(*out)
	y := &ym.Ym{}
	if err := encoding.Unmarshall(d, y); err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse file (%s) error :%v\n", *file, err)
		os.Exit(-1)
	}
	for i := 0; i < 16; i++ {
		filename := fmt.Sprintf("r%.2d.bin", i)
		w, err := os.Create(filepath.Join(*out, filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot create file (%s) error :%v\n", filename, err)
			os.Exit(-1)
		}
		defer w.Close()
		_, err = w.Write(y.Data[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot write file (%s) error :%v\n", filename, err)
			os.Exit(-1)
		}
	}
}

func CheckOutput(out string) error {
	infos, err := os.Stat(out)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(out, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", out, err)
			return err
		}
		return nil
	}
	if !infos.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory can not continue\n", out)
		return ErrorIsNotDirectory
	}
	return nil
}
