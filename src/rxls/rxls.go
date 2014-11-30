package main

//
// This source file is rxls command source.
//
// Description of rxls command.
// rx command function is view items in your trashbox.
// Viewing files are patial mattched from NAME in your trashbox.
//
// Options
// rxls [-l] NAME
//
// Usage
// rxls NAME : veiw contents name and original path from your trashbox
// rxls -l NAME : view backup date additionaly

//{{{ import
import (
	"flag"
	"fmt"
	// For rx command utilty.
	"rx_common"
	// Utilty for rxls and rxundo command.
	"rx_pv"
)

//}}}

//{{{ func show_location(name, contents_str, date_str string, rxFlagl bool)
// Show original path of backuped file or directory.
func show_location(name, contents_str, date_str string, rxFlagl bool) {

	// View origininal path of backuped file or directory.
	fmt.Print(name, " 's original location : ")
	if contents_str == "" {

		// Not logging original path caes.
		fmt.Println("do not know original location, you will recover manually")
	} else {

		// Logging original path case.
		fmt.Println(contents_str)
	}

	if date_str == "" {

		// Not logging backup date.
		date_str = "did not take the date log of rx executed"
	}

	// If add -l option, show backup date.
	if rxFlagl {

		fmt.Println("    date of rx executed :", date_str)
	}
}

//}}}

//{{{ func operation_of_ls(filename string, t *rx_common.TrashBox, rxFlagl bool)
// Controller view contents from your trashbox.
func operation_of_ls(filename string, t *rx_common.TrashBox, rxFlagl bool) {

	var (
		// For searching contents version.
		i int
		// For searching contents name in your trashbox.
		name string
	)

	fmt.Println("--> ls", filename+"*")
	// Get file name array and array length by patial mattched from filename.
	fileNames, fileNamesLen := rx_pv.Get_match(t.Get_trashBoxName(), filename)
	// Get logging data of file name array and array length by patial mattched from filename.
	filePrefixNames, _ := rx_pv.Get_match(t.Get_filePrefixDir(), filename)
	// Check if pattern mattch is successed.
	if fileNamesLen == 0 {

		// Not mattched case.
		fmt.Println(filename+"*", "is not match")
	} else {

		// Successed patial mattch case.
		// Show contents infomation recursively.
		for i, name = range fileNames {

			fmt.Printf("(%d) ", i)
			// Get logging data path from your trashbox.
			prefixName := rx_pv.Get_name(filePrefixNames, name)
			// Get original path and executed date.
			contents_str, date_str := rx_pv.Get_location(name, prefixName, t)
			// Show original path and executed date (selected -l option).
			show_location(name, contents_str, date_str, rxFlagl)
		}
	}
}

//}}}

//{{{ func main()
func main() {

	// declare option flag
	// rxFlagl is flag for -l option.
	// If this flag is true, show executed (backuped) date additionaly.
	var (
		rxFlagl bool
	)

	// Declare option name, initial value, description for help.
	// View help is selecting -help option.
	// is.) rxls -help
	flag.BoolVar(&rxFlagl, "l", false, "show file name before throw away")

	// Get TrashBox structure from rx command series configure file. (configure filename is $HOME/.rx)
	t := rx_common.Get_trashBox_cfg()

	// Analyze option flag from standard input, and set flag.
	// And set not flag standard input.
	flag.Parse()

	// Check if flag is setted.
	if flag.NArg() == 0 {

		// Show failed partial mattch.
		operation_of_ls("", t, rxFlagl)
	}

	// Show name and original path and backuped date (-l option) by patial mattched recursively.
	for i := 0; i < flag.NArg(); i++ {

		operation_of_ls(flag.Args()[i], t, rxFlagl)
	}
}

//}}}
