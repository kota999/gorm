package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	HOME        = os.Getenv("HOME")
	CFG_FILE    = HOME + "/.thrw"
	DEFAULT_DIR = HOME + "/.trashbox"
)

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

func main() {

	var (
		trashBoxName  string
		trashBoxCfg   bool
		trashBoxClear bool
		thrwFlagf     bool
		thrwFlagr     bool
		thrwFlagv     bool
	)

	flag.StringVar(&trashBoxName, "box", "", "trash box name")
	flag.BoolVar(&trashBoxClear, "c", false, "clear trash box")
	flag.BoolVar(&trashBoxClear, "C", false, "clear trash box")
	flag.BoolVar(&thrwFlagf, "f", false, "ignore warning")
	flag.BoolVar(&thrwFlagr, "r", false, "throw away directory, recursively")
	if thrwFlagr == false {
		flag.BoolVar(&thrwFlagr, "R", false, "throw away directory, recursively")
	}
	flag.BoolVar(&thrwFlagv, "v", false, "show file name before throw away")
	flag.Parse()
	if trashBoxName == "" {
		trashBoxCfg = false
	} else if trashBoxName == "-c" || trashBoxName == "-C" {
		trashBoxName = ""
		trashBoxCfg = false
	} else if trashBoxName == "-r" || trashBoxName == "-R" {
		trashBoxName = ""
		trashBoxCfg = false
	} else if trashBoxName == "-f" || trashBoxName == "-v" {
		trashBoxName = ""
		trashBoxCfg = false
	} else {
		trashBoxCfg = true
	}

	trashBoxName = set_trashBox_cfg(trashBoxName)

	if trashBoxCfg == false {
		for i := 0; i < flag.NArg(); i++ {
			if thrwFlagv {
				fmt.Println(flag.Args()[i])
			}
			newname := set_newname(flag.Args()[i], trashBoxName)
			if exist_file(flag.Args()[i]) {
				if isDirectory(flag.Args()[i]) == true {
					if thrwFlagr {
						remove(flag.Args()[i], newname)
					} else {
						fmt.Println(flag.Args()[i], "is directory. If you need throw away, use option -r or -R")
					}
				} else {
					remove(flag.Args()[i], newname)
				}
			} else {
				fmt.Println(flag.Args()[i], "is not exist.")
			}
		}
	}

	if trashBoxClear && trashBoxCfg == false {
		if err := os.RemoveAll(set_newname("", trashBoxName)); err != nil {
			fmt.Println(err)
		}
		os.Mkdir(trashBoxName, 0777)
	}
}
