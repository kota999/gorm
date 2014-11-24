package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"rx_common"
	"rx_pv"
	"strconv"
)

func show_location(name, contents_str, date_str string, rxFlagv, fileLenFlag bool) {
	if fileLenFlag {
		fmt.Print(name, " 's original location : ")
		if contents_str == "" {
			fmt.Println("do not know original location, you will recover manually")
		} else {
			fmt.Println(contents_str)
		}
		if date_str == "" {
			date_str = "did not take the date log of rx executed"
		}
		if rxFlagv {
			fmt.Println("    date of rx executed :", date_str)
		}
	}
}

func check_can_recov(contents_str string) bool {
	if contents_str == "" {
		return false
	} else {
		return true
	}
}

func get_can_recov_num(canRecovs []bool) int {
	var sum = 0
	for _, canRecov := range canRecovs {
		if canRecov {
			sum++
		}
	}
	return sum
}

func show_selector(canRecovs []bool) {
	var index = 0
	canRecov_num := get_can_recov_num(canRecovs)
	fmt.Print("you will select one from (")
	for i, canRecov := range canRecovs {
		if canRecov {
			fmt.Print(i)
			index++
			if index < canRecov_num {
				fmt.Print(", ")
			}
		}
	}
	fmt.Print(") > ")
}

func get_index_from_selector(canRecovs []bool, fileNamesLen int) int {
	show_selector(canRecovs)
	var (
		str   string
		index int
		err   error
	)

	for {
		if _, err := fmt.Scanf("%s", &str); err != nil {
			fmt.Println(err)
		}
		index, err = strconv.Atoi(str)
		if err != nil {
			fmt.Println(err)
		} else if index < 0 || index >= fileNamesLen {
			fmt.Println("Error: this number is invalid range")
		} else if canRecovs[index] {
			break
		} else {
			fmt.Println("Error: this option do not know original location")
		}
		show_selector(canRecovs)
	}
	return index
}

func undo(name, contents_str, trashBoxName string) {
	if rx_common.Exist_file(contents_str) {
		fmt.Println("Error:", contents_str, "is already exist")
	} else {
		if err := os.Rename(rx_common.Get_fullPath_trash(name, trashBoxName), contents_str); err != nil {
			if contents_str == "" {
				contents_str = "you will input file name"
			}
			fmt.Println("target file or directory :", contents_str)
			fmt.Println(err)
		} else {
			os.Remove(rx_common.Get_fullPath_prefix(rx_common.Get_prefix_filename(name), trashBoxName))
			fmt.Println("--> finish recovering to", contents_str)
		}
	}
}

func operation_of_undo(filename, trashBoxName string, rxFlagv bool) {
	var (
		i            int
		index        int
		name         string
		contents_str string
	)
	fmt.Println("--> recovering", filename+"*")
	fileInfo, _ := ioutil.ReadDir(trashBoxName)
	fileInfoPrefix, _ := ioutil.ReadDir(rx_common.Get_filePrefixDir(trashBoxName))
	fileNames, fileNamesLen := rx_pv.Check_match(fileInfo, filename)
	filePrefixNames, filePrefixNamesLen := rx_pv.Check_match(fileInfoPrefix, filename)
	if fileNamesLen == 0 {
		fmt.Println("Error:", filename, "is not backuped")
	} else if filePrefixNamesLen == 0 {
		fmt.Println("Error: do not know", filename, "'s original location")
		fmt.Println("you will recover manually")
	} else {
		var canRecovs = make([]bool, fileNamesLen)
		var prefixNames = make([]string, fileNamesLen)
		for i, name = range fileNames {
			if fileNamesLen != 1 {
				fmt.Printf("(%d) ", i)
			}
			prefixName := rx_pv.Check_name(filePrefixNames, name)
			contents_str, date_str := rx_pv.Check_location(name, prefixName, trashBoxName)
			show_location(name, contents_str, date_str, rxFlagv, fileNamesLen != 1)
			prefixNames[i] = contents_str
			canRecovs[i] = check_can_recov(contents_str)
		}
		if fileNamesLen != 1 {
			index = get_index_from_selector(canRecovs, fileNamesLen)
		}
		name = fileNames[index]
		contents_str = prefixNames[index]
		undo(name, contents_str, trashBoxName)
	}
}

func main() {

	var (
		trashBoxName string
		rxFlagv      bool
	)

	flag.BoolVar(&rxFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&rxFlagv, "v", false, "it is same option, -V")
	flag.Parse()

	trashBoxName = rx_common.Get_trashBox_cfg()

	if rxFlagv {
		fmt.Print("target files : ")
		fmt.Println(flag.Args())
	}
	for i := 0; i < flag.NArg(); i++ {
		operation_of_undo(flag.Args()[i], trashBoxName, rxFlagv)
	}

}
