package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	HOME          = os.Getenv("HOME")
	CFG_FILE      = HOME + "/.gorm"
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = "/.prefix/"
	GORM_EXTENDED = ".gorm"
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

//func read_cfg() string {
//contents, _ := ioutil.ReadFile(CFG_FILE)
//return string(contents)
//}

func read_cfg() string {
	readFile, _ := os.OpenFile(CFG_FILE, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readFile)
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

func read_file(filename string) string {
	contents, _ := ioutil.ReadFile(filename)
	return string(contents)
}

func write_cfg(name, cfgName string) {
	var nameVec []byte
	if cfgName == "" {
		os.Remove(CFG_FILE)
		cfgName = CFG_FILE
	}
	read_str := read_file(cfgName)
	if read_str == "" {
		nameVec = []byte(name)
	} else {
		nameVec = []byte(read_str + "\n" + name)
	}

	var writer *bufio.Writer
	writeFile, _ := os.OpenFile(cfgName, os.O_WRONLY|os.O_CREATE, 0600)
	writer = bufio.NewWriter(writeFile)
	writer.Write(nameVec)
	writer.Flush()
}

func set_trashBox_cfg() string {
	var trashBoxName string
	if exist_file(CFG_FILE) {
		if trashBoxName = read_cfg(); trashBoxName == "" {
			trashBoxName = DEFAULT_DIR
			write_cfg(trashBoxName, "")
		}
	} else {
		trashBoxName = DEFAULT_DIR
		write_cfg(trashBoxName, "")
	}

	make_trash_box(trashBoxName)

	return trashBoxName
}

func write_file_cfg(path, trashBoxName string, i int) {
	var fileCfg string
	filename := filepath.Base(path)
	if i == 0 {
		fileCfg = trashBoxName + FILE_PATH_DIR + filename + GORM_EXTENDED
	} else {
		fileCfg = trashBoxName + FILE_PATH_DIR + filename + "." + strconv.Itoa(i) + GORM_EXTENDED
	}

	fullPath, _ := filepath.Abs(path)
	write_cfg(fullPath, fileCfg)
	now := time.Now().String()[:19]
	write_cfg(now, fileCfg)
}

func show_path(path string, gormFlagv bool) {
	if gormFlagv {
		fmt.Println(path)
	}
}

func remove(path string, trashBoxName string, i int) {
	var name string
	filename := filepath.Base(path)
	newPath := trashBoxName + "/" + filename
	if i == 0 {
		name = newPath
	} else {
		name = newPath + "." + strconv.Itoa(i)
	}

	if exist_file(name) {
		remove(path, trashBoxName, i+1)
	} else {
		if err := os.Rename(path, name); err != nil {
			fmt.Println(newPath)
			fmt.Println(err)
		} else {
			write_file_cfg(path, trashBoxName, i)
		}
	}
}

func operation_of_remove(path string, trashBoxName string, gormFlagr bool) {
	if exist_file(path) {
		if isDirectory(path) == true {
			if gormFlagr {
				remove(path, trashBoxName, 0)
			} else {
				fmt.Println(path, "is directory. If you need throw away, use option -r or -R")
			}
		} else {
			remove(path, trashBoxName, 0)
		}
	} else {
		fmt.Println(path, "is not exist.")
	}
}

func main() {

	var (
		trashBoxName string
		gormFlagr    bool
		gormFlagv    bool
	)

	flag.BoolVar(&gormFlagr, "R", false, "throw away directory, recursively")
	flag.BoolVar(&gormFlagr, "r", false, "it is same option, -R")
	flag.BoolVar(&gormFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&gormFlagv, "v", false, "it is same option, -V")
	flag.Parse()

	trashBoxName = set_trashBox_cfg()

	for i := 0; i < flag.NArg(); i++ {
		show_path(flag.Args()[i], gormFlagv)
		operation_of_remove(flag.Args()[i], trashBoxName, gormFlagr)
	}

}
