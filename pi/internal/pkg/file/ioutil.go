package file

import (
	"VGO/pkg/fun"
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Read(path string) []byte {
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("ioutil open err : %v\n", err)
	}
	return byt
}
func Write(path, str string) {

	err := ioutil.WriteFile(path, []byte(str), os.ModePerm)
	if err != nil {
		fmt.Printf("ioutil write err : %v\n", err)
	}
}

func Unzip(zipFile, dest string, exclude []string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	toDist := func(f *zip.File) error {
		//log.Println(f.Name)
		if fun.InSliceString(f.Name, exclude) {
			return nil
		}
		fPath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fPath, 0755); err != nil {
				return err
			}
			return nil
		}
		if err = os.MkdirAll(filepath.Dir(fPath), 0755); err != nil {
			return err
		}
		w, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			return err
		}
		defer w.Close()
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		_, err = io.Copy(w, rc)
		return err
	}
	for _, f := range reader.File {
		err := toDist(f)
		if err != nil {
			return err
		}
	}
	return nil
}
