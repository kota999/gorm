package main

//{{{ import
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"rx_common"
)

//}}}

//{{{ func remove(path string, t *rx_common.TrashBox, index int)
func remove(path string, t *rx_common.TrashBox, index int) {
	filename := filepath.Base(path)
	newPath := t.Get_fullPath_trash(filename)
	newPath = rx_common.Get_filename_version(newPath, index)

	if rx_common.Exist_file(newPath) {
		remove(path, t, index+1)
	} else {
		if err := os.Rename(path, newPath); err != nil {
			fmt.Println("backup to", newPath, "is failed")
			fmt.Println(err)
		} else {
			rx_common.Write_file_cfg(path, t, index)
		}
	}
}

//}}}

//{{{ func operation_of_remove(path string, t *rx_common.TrashBox, rxFlagr bool)
func operation_of_remove(path string, t *rx_common.TrashBox, rxFlagr bool) {
	if rx_common.Exist_file(path) {
		if rx_common.IsDirectory(path) == true {
			if rxFlagr {
				remove(path, t, 0)
			} else {
				fmt.Println(path, "is directory. If you need throw away, use option -r or -R")
			}
		} else {
			remove(path, t, 0)
		}
	} else {
		fmt.Println(path, "is not exist.")
	}
}

//}}}

//{{{ func main()
func main() {

	var (
		rxFlagr bool
		rxFlagv bool
	)

	flag.BoolVar(&rxFlagr, "R", false, "throw away directory, recursively")
	flag.BoolVar(&rxFlagr, "r", false, "it is same option, -R")
	flag.BoolVar(&rxFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&rxFlagv, "v", false, "it is same option, -V")
	flag.Parse()

	t := rx_common.Get_trashBox_cfg()

	for i := 0; i < flag.NArg(); i++ {
		rx_common.Show_path(flag.Args()[i], rxFlagv)
		operation_of_remove(flag.Args()[i], t, rxFlagr)
	}

}

//}}}
