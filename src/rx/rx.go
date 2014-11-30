package main

//
// This source file is rx command source.
//
// rx command is remove + backup command.
//
// Description of rx command.
// Throw away NAME into trashbox directory, default trashbox name is $HOME/.trahbox.
// If you will change trahbox path, rename trashbox path in rx configure file ($HOME/.rx)
// or execute rxbox command as below.
// rxbox -box TRASHBOX_NAME
//
// Options
// rx [-r|-R] [-v|-V] NAME
// [-r|-R] : execute rx command recursively.
// [-v|-V] : show rx command process.

//{{{ import
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	// For rx command series utilty.
	"rx_common"
)

//}}}

//{{{ func remove(path string, t *rx_common.TrashBox, index int)
// remove file or directory of target
func remove(path string, t *rx_common.TrashBox, index int) {

	// get filename from path of file or directory of target
	filename := filepath.Base(path)
	// get fullpath to remove into trashbox
	newPath := t.Get_fullPath_trash(filename)
	// get file or directory version from file or directory on trashbox
	newPath = rx_common.Get_filename_version(newPath, index)

	// Check exist file or directory
	// If exsiting this file or directory version name into trashbox.
	if rx_common.Exist_file(newPath) {

		// Existing this file or directory version name into trashbox.
		// Try checking whether there is next file or directory version name.
		remove(path, t, index+1)
	} else {

		// Not existing this file or directory version name into trashbox.
		// Try remove file or directory of traget file or directory,
		// backup file or directory to trashbox in considerly.
		if err := os.Rename(path, newPath); err != nil {

			// If it is an error in executing remove + backup command,
			// standard out file or directory name and fail messages.
			fmt.Println("backup to", newPath, "is failed")
			fmt.Println(err)
		} else {

			// Success executing remove + backup command,
			// write configure file.
			// - orignal file or directory path (for Unremove)
			// - executing time stamp (for version control)
			rx_common.Write_file_cfg(path, t, index)
		}
	}
}

//}}}

//{{{ func operation_of_remove(path string, t *rx_common.TrashBox, rxFlagr bool)
// Controller remove + backup command.
func operation_of_remove(path string, t *rx_common.TrashBox, rxFlagr bool) {

	// Check if exsiting this file or directory name of target.
	if rx_common.Exist_file(path) {

		// Existing this file or directory name of traget,
		// check this name of traget, file type or directory type.
		if rx_common.IsDirectory(path) == true {

			// This name of traget is directory type.
			// Need option -r or -R,
			// check if selected r or R option.
			if rxFlagr {

				// Option -r or -R is selected, execute remove + backup command.
				remove(path, t, 0)
			} else {

				// Option -r or -R is not selected, not execute remove + backup command.
				// Standard out message: need -r or -R option.
				fmt.Println(path, "is directory. If you need throw away, use option -r or -R")
			}
		} else {

			// This name of target is file type, execute remove + backup command.
			remove(path, t, 0)
		}
	} else {

		// Not existing this file or directory name of traget,
		// Standard out error messages, file or directory of traget is not exist.
		fmt.Println(path, "is not exist.")
	}
}

//}}}

//{{{ func main()
// Main function os rx command
func main() {

	// declare option flag
	// rxFlagr is flag for -r, -R option.
	// If this option is true, remove + backup recursively into directory of target.
	// rxFlagv is flag for -v, -V option.
	// If this option is true, standard out message: file or directory name, executing message.
	var (
		rxFlagr bool
		rxFlagv bool
	)

	// Declare option name, initial value, description for help.
	// View help is selecting -help option.
	// ie.) rx -help
	flag.BoolVar(&rxFlagr, "R", false, "throw away directory, recursively")
	flag.BoolVar(&rxFlagr, "r", false, "it is same option, -R")
	flag.BoolVar(&rxFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&rxFlagv, "v", false, "it is same option, -V")
	// Analyze option flag from standard input, and set flag.
	// And set not flag standard input.
	flag.Parse()

	// Get TrashBox structure from rx comannd series configure file. (configure filename is $HOME/.rx)
	// This TrashBox structure has trashbox name, get path name into trashbox, operation trashbox path,,,
	t := rx_common.Get_trashBox_cfg()

	// Execute remove + backup command recursively.
	for i := 0; i < flag.NArg(); i++ {

		// If view option flag is true, standard out file and directory name.
		rx_common.Show_path(flag.Args()[i], rxFlagv)
		// Call command remove + backup controller.
		operation_of_remove(flag.Args()[i], t, rxFlagr)
	}

}

//}}}
