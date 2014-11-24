package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"rx_common"
)

func check_match_num(infos []os.FileInfo, pattern string) int {
	var sum = 0
	for _, info := range infos {
		var name = info.Name()
		if name != rx_common.FILE_PATH_DIR {
			if matched, _ := path.Match(pattern+"*", name); matched {
				sum += 1
			}
		}
	}
	return sum
}

func check_match_names(infos []os.FileInfo, pattern string, num int) []string {
	var i = 0
	var names = make([]string, num)
	for _, info := range infos {
		var name = info.Name()
		if name != rx_common.FILE_PATH_DIR {
			if matched, _ := path.Match(pattern+"*", name); matched {
				names[i] = name
				i += 1
			}
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
		if matched, _ := path.Match(rx_common.Get_prefix_filename(pattern), name); matched {
			return name
		}
	}
	return ""
}

func check_location(name, prefixName, trashBoxName string, rxFlagl bool) (string, string) {
	reader := rx_common.Generate_reader(rx_common.Get_filePrefixDir(trashBoxName) + prefixName)
	contents_str := rx_common.ReadLine(reader)
	date_str := rx_common.ReadLine(reader)
	fmt.Print(name, " 's original location : ")
	if contents_str == "" {
		fmt.Println("do not know original location, you will recover manually")
	} else {
		fmt.Println(contents_str)
	}

	if date_str == "" {
		date_str = "did not take the date log of rx executed"
	}
	if rxFlagl {
		fmt.Println("    date of rx executed :", date_str)
	}

	return contents_str, date_str
}

func operation_of_ls(filename, trashBoxName string, rxFlagl bool) {
	var (
		i    int
		name string
	)
	fmt.Println("--> ls", filename+"*")
	fileInfo, _ := ioutil.ReadDir(trashBoxName)
	fileInfoPrefix, _ := ioutil.ReadDir(rx_common.Get_filePrefixDir(trashBoxName))
	fileNames, fileNamesLen := check_match(fileInfo, filename)
	filePrefixNames, _ := check_match(fileInfoPrefix, filename)
	if fileNamesLen == 0 {
		fmt.Println(filename+"*", "is not match")
	} else {
		for i, name = range fileNames {
			fmt.Printf("(%d) ", i)
			prefixName := check_name(filePrefixNames, name)
			check_location(name, prefixName, trashBoxName, rxFlagl)
		}
	}
}

func main() {

	var (
		trashBoxName string
		rxFlagl      bool
	)
	flag.BoolVar(&rxFlagl, "l", false, "show file name before throw away")

	trashBoxName = rx_common.Get_trashBox_cfg()

	flag.Parse()
	if flag.NArg() == 0 {
		operation_of_ls("", trashBoxName, rxFlagl)
	}
	for i := 0; i < flag.NArg(); i++ {
		operation_of_ls(flag.Args()[i], trashBoxName, rxFlagl)
	}
}
