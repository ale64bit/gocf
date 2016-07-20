package main

import (
	"fmt"
	"io"
	"os"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FileNotExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// @see http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func Yes(msg string) bool {
	fmt.Print(msg, " (yN): ")
	var ans string
	fmt.Scanf("%s\n", &ans)
	for ans != "y" && ans != "N" {
		fmt.Println("Please, reply with 'y' or 'N' only.")
		fmt.Print(msg, "(yN): ")
		fmt.Scanf("%s\n", &ans)
	}
	return ans == "y"
}

func ReadDefault(msg, def string) string {
	fmt.Print(msg + " [default=" + def + "]: ")
	var ret string
	fmt.Scanf("%s\n", &ret)
	if len(ret) == 0 {
		ret = def
	}
	return ret
}
