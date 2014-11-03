package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const CFG_FILE = ".thrw"
const DEFAULT_DIR = ".trash"

func make_trash_box(dirName string) {
	//if err := os.Mkdir(dirName, 0777); err != nil {
	//fmt.Println(" directory is already exist !")
	//}
	os.Mkdir(dirName, 0777)
}

func exist_file(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func isDirectory(name string) bool {
	fInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func remove_cfg() {
	os.Remove(CFG_FILE)
}

func read_cfg() string {
	contents, _ := ioutil.ReadFile(CFG_FILE)
	return string(contents)
}

func write_cfg(dirName string) {
	remove_cfg()
	var writer *bufio.Writer
	dirVec := []byte(dirName)

	write_file, _ := os.OpenFile(CFG_FILE, os.O_WRONLY|os.O_CREATE, 0600)
	writer = bufio.NewWriter(write_file)
	writer.Write(dirVec)
	writer.Flush()
}

func set_newname(oldName, dirName string) string {
	var buffer bytes.Buffer
	buffer.WriteString(dirName)
	buffer.WriteString("/")
	buffer.WriteString(oldName)
	return buffer.String()
}

func main() {

	var (
		trashBoxName string
		thrwFlagf    bool
		thrwFlagr    bool
		thrwFlagv    bool
	)

	flag.StringVar(&trashBoxName, "box", "blank", "trash box name")
	flag.BoolVar(&thrwFlagf, "f", false, "ignore warning")
	flag.BoolVar(&thrwFlagr, "r", false, "throw away directory, recursively")
	if thrwFlagr == false {
		flag.BoolVar(&thrwFlagr, "R", false, "throw away directory, recursively")
	}
	flag.BoolVar(&thrwFlagv, "v", false, "show file name before throw away")
	flag.Parse()

	if trashBoxName != "blank" {
		write_cfg(trashBoxName)
	} else if exist_file(CFG_FILE) {
		trashBoxName = read_cfg()
	} else {
		trashBoxName = DEFAULT_DIR
		write_cfg(trashBoxName)
	}
	make_trash_box(trashBoxName)

	for i := 0; i < flag.NArg(); i++ {
		if thrwFlagv {
			fmt.Println(flag.Args()[i])
		}
		newname := set_newname(flag.Args()[i], trashBoxName)
		if exist_file(flag.Args()[i]) {
			if thrwFlagr {
				fmt.Println(newname)
				err := os.Rename(flag.Args()[i], newname)
				if err != nil {
					fmt.Println(err)
				}
			} else if isDirectory(flag.Args()[i]) == false {
				os.Rename(flag.Args()[i], newname)
			} else {
				fmt.Println(flag.Args()[i], "is directory. If you need throw away, use option -r or -R")
			}
		} else {
			fmt.Println(flag.Args()[i], "is not exist.")
		}
	}

}
