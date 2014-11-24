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

func check_location(name, prefixName, trashBoxName string, gormFlagl bool) (string, string) {
	readerFile, _ := os.OpenFile(trashBoxName+FILE_PATH_DIR+"/"+prefixName, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readerFile)
	contents, _, _ := reader.ReadLine()
	contents_str := string(contents)
	contents, _, _ = reader.ReadLine()
	date_str := string(contents)
	fmt.Print(name, " 's original location : ")
	if contents_str == "" {
		fmt.Println("do not know original location, you will recover manually")
	} else {
		fmt.Println(contents_str)
	}

	if date_str == "" {
		date_str = "did not take the date log of rx executed"
	}
	if gormFlagl {
		fmt.Println("    date of rx executed :", date_str)
	}

	return contents_str, date_str
}

func operation_of_ls(filename, trashBoxName string, gormFlagl bool) {
	var (
		i    int
		name string
	)
	fmt.Println("--> ls", filename+"*")
	fileInfo, _ := ioutil.ReadDir(trashBoxName)
	fileInfoPrefix, _ := ioutil.ReadDir(trashBoxName + FILE_PATH_DIR)
	fileNames, fileNamesLen := check_match(fileInfo, filename)
	filePrefixNames, _ := check_match(fileInfoPrefix, filename)
	if fileNamesLen == 0 {
		fmt.Println(filename+"*", "is not match")
	} else {
		for i, name = range fileNames {
			fmt.Printf("(%d) ", i)
			prefixName := check_name(filePrefixNames, name)
			check_location(name, prefixName, trashBoxName, gormFlagl)
		}
	}
}

func main() {

	var (
		trashBoxName string
		gormFlagl    bool
	)
	flag.BoolVar(&gormFlagl, "l", false, "show file name before throw away")

	trashBoxName = set_trashBox_cfg()

	flag.Parse()
	if flag.NArg() == 0 {
		operation_of_ls("", trashBoxName, gormFlagl)
	}
	for i := 0; i < flag.NArg(); i++ {
		operation_of_ls(flag.Args()[i], trashBoxName, gormFlagl)
	}
}
