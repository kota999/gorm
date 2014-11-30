package main

//
// This source file is rxbox command source.
//
// rxbox command is operation for your trashbox directory.
//
// Description of rxbox command.
// rxbox command function is
// 1. initializing your trashbox
// 2. clearing all contents from your trashbox
// 3. change trashbox directory
//
// Options
// rxbox [-d] [-c|-C] [-box] TRASHBOX_NAME
//
// Usages
// rxbox -d : initialize your trashbox configuration (your trashbox directory -> $/.trashbox)
//            remove your trashbox directory
// rxbox -c (-C) : remove all contents in your trashbox
// rxbox -box TRASHBOX_NAME : exchange trashbox configuration
//                              rewrite new trashbox directory path in $HOME/.rx .

//{{{ import
import (
	"flag"
	"fmt"
	"os"
	// Utilty packages for rx command series.
	"rx_common"
	"strings"
)

//}}}

//{{{ func clear_trash(trashBoxName string)
// Clear your trashbox.
func clear_trash(trashBoxName string) {

	// Check trashbox name is root directory ("/").
	if trashBoxName == "/" {

		// Standard output message: Fatal and unexpected process.
		fmt.Println("!!!!! ileagl option !!!!!")
		// rxbox command is sttoped by fatal Error.
		os.Exit(1)
	}

	// Clear your trashbox.
	if err := os.RemoveAll(trashBoxName); err != nil {

		// Standard out error message in condition of falied clear trashbox case.
		fmt.Println(err)
	}
	os.Mkdir(trashBoxName, 0777)
}

//}}}

//{{{ func remove_trash_box(tarshBoxName string)
// Remove your trashbox.
func remove_trash_box(trashBoxName string) {

	// Check trashbox name is root directory ("/").
	if trashBoxName == "/" {

		// Standard output message: Fatal and unexpected process.
		fmt.Println("!!!!! ilegal option !!!!!")
		// rxbox command is sttoped by fatal Error.
		os.Exit(1)
	}

	// Remove your trashbox.
	if err := os.RemoveAll(trashBoxName); err != nil {

		// Standard out error message in condition of failed remove trashbox case.
		fmt.Println(err)
	}
}

//}}}

//{{{ func init_default()
func init_default() {

	// Standard out messages: your options is initializing your trashbox configure.
	fmt.Println("initializing trash-box configure and clear trash-box.")
	fmt.Println("INFO: initialized trash-box directory is $HOME/.trashbox.")

	// Declare TrashBox structure for initializing your trashbox.
	var t *rx_common.TrashBox
	// Call TrashBox structure by recent cunfiguration.
	t = rx_common.Get_trashBox_cfg()
	// Remove recent TrashBox structure.
	remove_trash_box(t.Get_trashBoxName())
	// Regenerate default TrashBox structure.
	t = rx_common.Set_trashBox_cfg(rx_common.DEFAULT_DIR)
}

//}}}

//{{{ main
func main() {

	// Declare option flag.
	// trashBoxName is new trashbox name.
	// trashBoxCfg is change trashbox directory option flag (-box option).
	// trashBoxClear is -c or -C option flag for clearing all contents in your trashbox.
	// defaultCfg is -d option flag for initialize your trashbox.
	var (
		trashBoxName  = ""
		trashBoxCfg   bool
		trashBoxClear bool
		defaultCfg    bool
		t             *rx_common.TrashBox
	)

	// Declare option name, initial value, description for help.
	// View help is selecting -help option.
	// ie.) rxbox -help
	flag.BoolVar(&trashBoxCfg, "box", false, "set trash box name")
	flag.BoolVar(&trashBoxClear, "C", false, "clear trash box")
	flag.BoolVar(&trashBoxClear, "c", false, "it is same option, -C")
	flag.BoolVar(&defaultCfg, "d", false, "initialize default setting (default dir: $HOME/.trashbox)")
	// Analyze option flag from standard input, and set flag.
	// And set not flag standard input.
	flag.Parse()

	// Initialize trashobox case ( -d option case)
	if defaultCfg {

		// Initialize trashbox
		init_default()
		// Exit process.
		os.Exit(1)
	}

	// Check -box option or exist standard input.
	if trashBoxCfg == false || flag.NArg() == 0 {

		// Not exist -box option and standard input case.
		// trashBoxName is not setted.
		trashBoxName = ""
	} else {

		// -box option flag is true and exist TRASHBOX_NAME case.
		// Set trashBoxName to TRASHBOX_NAME
		trashBoxName = flag.Args()[0]
		if strings.HasPrefix(trashBoxName, "-") || trashBoxName == "/" {

			// Inputted invlid option.
			// trashBoxName is initialized.
			trashBoxName = ""
		}
	}

	// Get trashBox structure.
	// In the case of changing trashBoxName, change trashBoxName in trashBox structure.
	if trashBoxName == "" {

		t = rx_common.Get_trashBox_cfg()
	} else {

		t = rx_common.Set_trashBox_cfg(trashBoxName)
	}

	// In the case of -box option is false, -C or -c option is true, trashBoxName is not root directory/
	// Clear all contents in your trashbox.
	if trashBoxClear && trashBoxCfg == false && t.Get_trashBoxName() != "/" {

		clear_trash(t.Get_trashBoxName())
	}
}

//}}}
