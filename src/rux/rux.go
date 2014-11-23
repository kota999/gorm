package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var (
	HOME          = os.Getenv("HOME")
	CFG_FILE      = HOME + "/.gorm"
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = "/.prefix"
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

func read_gorm_cfg() string {
	readerFile, _ := os.OpenFile(CFG_FILE, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readerFile)
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

func write_cfg(dirName string) {
	os.Remove(CFG_FILE)

	dirVec := []byte(dirName)

	writeFile, _ := os.OpenFile(CFG_FILE, os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	writer.Write(dirVec)
	writer.Flush()
}

func set_trashBox_cfg() string {
	var trashBoxName string
	if exist_file(CFG_FILE) {
		if trashBoxName = read_gorm_cfg(); trashBoxName == "" {
			trashBoxName = DEFAULT_DIR
			write_cfg(trashBoxName)
		}
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

func check_match_num(infos []os.FileInfo, pattern string) int {
	var sum = 0
	for _, info := range infos {
		var name = info.Name()
		if matched, _ := path.Match(pattern+"*", name); matched {
			sum += 1
		}
	}
	return sum
}

func check_match_names(infos []os.FileInfo, pattern string, num int) []string {
	var i = 0
	var names = make([]string, num)
	for _, info := range infos {
		var name = info.Name()
		if matched, _ := path.Match(pattern+"*", name); matched {
			names[i] = name
			i += 1
		}
	}
	return names
}

func check_match(infos []os.FileInfo, pattern string) ([]string, int) {
	num := check_match_num(infos, pattern)
	names := check_match_names(infos, pattern, num)
	return names, num
}

func check_name(names []string, pattern string) string {
	for _, name := range names {
		if matched, _ := path.Match(pattern+GORM_EXTENDED, name); matched {
			return name
		}
	}
	return ""
}

func check_locate(name, prefixName, trashBoxName string) string {
	fmt.Print(name, " 's original location : ")
	readerFile, _ := os.OpenFile(trashBoxName+FILE_PATH_DIR+"/"+prefixName, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readerFile)
	contents, _, _ := reader.ReadLine()
	contents_str := string(contents)
	if contents_str == "" {
		fmt.Println("do not know original location")
	} else {
		fmt.Println(contents_str)
	}
	return contents_str
}

func undo(name, contents_str, trashBoxName string) {
	if exist_file(contents_str) {
		fmt.Println("Error:", contents_str, "is already exist")
	} else {
		if err := os.Rename(trashBoxName+"/"+name, contents_str); err != nil {
			fmt.Println(contents_str)
			fmt.Println(err)
		} else {
			os.Remove(trashBoxName + FILE_PATH_DIR + "/" + name + GORM_EXTENDED)
		}
	}
}

func operation_of_undo(filename, trashBoxName string, gormFlagv bool) {
	fmt.Println("***", filename, "recovering ***")
	fileInfo, _ := ioutil.ReadDir(trashBoxName)
	fileInfoPrefix, _ := ioutil.ReadDir(trashBoxName + FILE_PATH_DIR)
	fileNames, fileNamesLen := check_match(fileInfo, filename)
	filePrefixNames, filePrefixNamesLen := check_match(fileInfoPrefix, filename)
	if fileNamesLen == 0 {
		fmt.Println("Error:", filename, "is not backup")
	} else if filePrefixNamesLen == 0 {
		fmt.Println("Error: do not know", filename, "'s original location")
		fmt.Println("you will manuary recover")
	} else {
		for _, name := range fileNames {
			prefixName := check_name(filePrefixNames, name)
			if fileNamesLen == 1 {
				contents_str := check_locate(name, prefixName, trashBoxName)
				undo(name, contents_str, trashBoxName)
			} else {
				contents_str := check_locate(name, prefixName, trashBoxName)
				fmt.Println(contents_str)
			}
		}
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
		operation_of_undo(flag.Args()[i], trashBoxName, gormFlagv)
	}

}
