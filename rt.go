package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	HOME        = os.Getenv("HOME")
	CFG_FILE    = HOME + "/.rt"
	DEFAULT_DIR = HOME + "/.trashbox"
)

func make_trash_box(dirName string) {
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

func set_trashBox_cfg(trashBoxName string) string {
	if trashBoxName != "" {
		fmt.Println("you setted trash-box directory option, so throw-away and trash-box clear options are not effective.")
		fmt.Println("")
		if exist_file(trashBoxName) && isDirectory(trashBoxName) == false {
			fmt.Println("your option directory is existed as file type,")
			fmt.Println("trash box dir is setted .trash .")
			fmt.Println("")
			trashBoxName = DEFAULT_DIR
		}
		write_cfg(trashBoxName)

	} else if exist_file(CFG_FILE) {
		trashBoxName = read_cfg()

	} else {
		trashBoxName = DEFAULT_DIR
		write_cfg(trashBoxName)

	}

	make_trash_box(trashBoxName)

	return trashBoxName
}

func show_path(path string, rtFlagv bool) {
	if rtFlagv {
		fmt.Println(path)
	}
}

func remove(path, newpath string) {
	if exist_file(newpath) {
		Reremove(path, newpath, 1)
	} else {
		if err := os.Rename(path, newpath); err != nil {
			fmt.Println(newpath)
			fmt.Println(err)
		}
	}
}

func Reremove(path string, newpath string, i int) {
	if exist_file(newpath + "." + strconv.Itoa(i)) {
		Reremove(path, newpath, i+1)
	} else {
		if err := os.Rename(path, newpath+"."+strconv.Itoa(i)); err != nil {
			fmt.Println(err)
		}
	}
}

func operation_of_remove(path string, trashBoxName string, rtFlagr bool) {
	newpath := trashBoxName + "/" + path
	if exist_file(path) {
		if isDirectory(path) == true {
			if rtFlagr {
				remove(path, newpath)
			} else {
				fmt.Println(path, "is directory. If you need throw away, use option -r or -R")
			}
		} else {
			remove(path, newpath)
		}
	} else {
		fmt.Println(path, "is not exist.")
	}
}

func clear_trash(trashBoxName string) {
	if err := os.RemoveAll(trashBoxName + "/"); err != nil {
		fmt.Println(err)
	}
	os.Mkdir(trashBoxName, 0777)
}

func main() {

	var (
		trashBoxName  string
		trashBoxCfg   bool
		trashBoxClear bool
		rtFlagr       bool
		rtFlagv       bool
	)

	flag.StringVar(&trashBoxName, "box", "", "trash box name")
	flag.BoolVar(&trashBoxClear, "c", false, "clear trash box")
	flag.BoolVar(&trashBoxClear, "C", false, "clear trash box")
	flag.BoolVar(&rtFlagr, "r", false, "throw away directory, recursively")
	flag.BoolVar(&rtFlagr, "R", false, "throw away directory, recursively")
	flag.BoolVar(&rtFlagv, "v", false, "show file name before throw away")

	flag.Parse()

	if trashBoxName == "" {
		trashBoxCfg = false
	} else if strings.HasPrefix(trashBoxName, "-") {
		trashBoxName = ""
		trashBoxCfg = false
	} else {
		trashBoxCfg = true
	}

	trashBoxName = set_trashBox_cfg(trashBoxName)

	if trashBoxCfg == false {
		for i := 0; i < flag.NArg(); i++ {
			show_path(flag.Args()[i], rtFlagv)
			operation_of_remove(flag.Args()[i], trashBoxName, rtFlagr)
		}
	}

	if trashBoxClear && trashBoxCfg == false {
		clear_trash(trashBoxName)
	}
}
