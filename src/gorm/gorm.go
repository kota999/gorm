package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	GORM_EXTENDED = ".gorm"
	HOME          = os.Getenv("HOME")
	CFG_FILE      = HOME + "/" + GORM_EXTENDED
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = ".prefix"
)

func make_trash_box(dirName string) {
	os.Mkdir(dirName, 0777)
	os.Mkdir(dirName+"/"+FILE_PATH_DIR, 0777)
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

	writeFile, _ := os.OpenFile(CFG_FILE, os.O_WRONLY|os.O_CREATE, 0600)
	writer = bufio.NewWriter(writeFile)
	writer.Write(dirVec)
	writer.Flush()
}

func write_filePath(fullPath, filePathCfg string) {
	var writer *bufio.Writer
	fullPathVec := []byte(fullPath)

	writeFile, _ := os.OpenFile(filePathCfg+GORM_EXTENDED, os.O_WRONLY|os.O_CREATE, 0600)
	writer = bufio.NewWriter(writeFile)
	writer.Write(fullPathVec)
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

func show_path(path string, gormFlagv bool) {
	if gormFlagv {
		fmt.Println(path)
	}
}

func remove(path, newPath string) {
	if exist_file(newPath) {
		Reremove(path, newPath, 1)
	} else {
		if err := os.Rename(path, newPath); err != nil {
			fmt.Println(newPath)
			fmt.Println(err)
		} else {
			fullPath, _ := filepath.Abs(path)
			write_filePath(fullPath, newPath)
		}
	}
}

func Reremove(path string, newPath string, i int) {
	if exist_file(newPath + "." + strconv.Itoa(i)) {
		Reremove(path, newPath, i+1)
	} else {
		if err := os.Rename(path, newPath+"."+strconv.Itoa(i)); err != nil {
			fmt.Println(err)
		} else {
			fullPath, _ := filepath.Abs(path)
			write_filePath(fullPath, newPath+"."+strconv.Itoa(i))
		}
	}
}

func operation_of_remove(path string, trashBoxName string, gormFlagr bool) {
	filename := filepath.Base(path)
	newPath := trashBoxName + "/" + filename
	if exist_file(path) {
		if isDirectory(path) == true {
			if gormFlagr {
				remove(path, newPath)
			} else {
				fmt.Println(path, "is directory. If you need throw away, use option -r or -R")
			}
		} else {
			remove(path, newPath)
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
		gormFlagr     bool
		gormFlagv     bool
	)

	flag.StringVar(&trashBoxName, "box", "", "trash box name")
	flag.BoolVar(&trashBoxClear, "c", false, "clear trash box")
	flag.BoolVar(&trashBoxClear, "C", false, "clear trash box")
	flag.BoolVar(&gormFlagr, "r", false, "throw away directory, recursively")
	flag.BoolVar(&gormFlagr, "R", false, "throw away directory, recursively")
	flag.BoolVar(&gormFlagv, "v", false, "show file name before throw away")
	flag.BoolVar(&gormFlagv, "V", false, "show file name before throw away")

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
			show_path(flag.Args()[i], gormFlagv)
			operation_of_remove(flag.Args()[i], trashBoxName, gormFlagr)
		}
	}

	if trashBoxClear && trashBoxCfg == false {
		clear_trash(trashBoxName)
	}
}
