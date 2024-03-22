package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func WalkDir(dirname string) {
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func WriteByteData(buf *bytes.Buffer, dir string, fileName string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("unable to mkdir: %s", dir)
		}
	}
	dFile, err := os.Create(dir + "/" + fileName)
	if err != nil {
		fmt.Printf("unable to create/open file: %s", fileName)
	}
	defer dFile.Close()

	fmt.Printf("Writing %s...\n", fileName)

	if _, err = dFile.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("unable to write instructions to: %s", dir+"/"+fileName)
	}
	return nil
}
