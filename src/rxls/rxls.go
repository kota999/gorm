package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"rx_common"
	"rx_pv"
)

func show_location(name, contents_str, date_str string, rxFlagl bool) {
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

}

func operation_of_ls(filename, trashBoxName string, rxFlagl bool) {
	var (
		i    int
		name string
	)
	fmt.Println("--> ls", filename+"*")
	fileInfo, _ := ioutil.ReadDir(trashBoxName)
	fileInfoPrefix, _ := ioutil.ReadDir(rx_common.Get_filePrefixDir(trashBoxName))
	fileNames, fileNamesLen := rx_pv.Check_match(fileInfo, filename)
	filePrefixNames, _ := rx_pv.Check_match(fileInfoPrefix, filename)
	if fileNamesLen == 0 {
		fmt.Println(filename+"*", "is not match")
	} else {
		for i, name = range fileNames {
			fmt.Printf("(%d) ", i)
			prefixName := rx_pv.Check_name(filePrefixNames, name)
			contents_str, date_str := rx_pv.Check_location(name, prefixName, trashBoxName)
			show_location(name, contents_str, date_str, rxFlagl)
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
