package main

//{{{ import
import (
	"flag"
	"fmt"
	"rx_common"
	"rx_pv"
)

//}}}

//{{{ func show_location(name, contents_str, date_str string, rxFlagl bool)
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

//}}}

//{{{ func operation_of_ls(filename string, t *rx_common.TrashBox, rxFlagl bool)
func operation_of_ls(filename string, t *rx_common.TrashBox, rxFlagl bool) {
	var (
		i    int
		name string
	)
	fmt.Println("--> ls", filename+"*")
	fileNames, fileNamesLen := rx_pv.Get_match(t.Get_trashBoxName(), filename)
	filePrefixNames, _ := rx_pv.Get_match(t.Get_filePrefixDir(), filename)
	if fileNamesLen == 0 {
		fmt.Println(filename+"*", "is not match")
	} else {
		for i, name = range fileNames {
			fmt.Printf("(%d) ", i)
			prefixName := rx_pv.Get_name(filePrefixNames, name)
			contents_str, date_str := rx_pv.Get_location(name, prefixName, t)
			show_location(name, contents_str, date_str, rxFlagl)
		}
	}
}

//}}}

//{{{ func main()
func main() {

	var (
		rxFlagl bool
	)
	flag.BoolVar(&rxFlagl, "l", false, "show file name before throw away")

	t := rx_common.Get_trashBox_cfg()

	flag.Parse()
	if flag.NArg() == 0 {
		operation_of_ls("", t, rxFlagl)
	}
	for i := 0; i < flag.NArg(); i++ {
		operation_of_ls(flag.Args()[i], t, rxFlagl)
	}
}

//}}}
