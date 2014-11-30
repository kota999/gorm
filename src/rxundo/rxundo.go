package main

//
// This source file is rxundo command source.
//
// rxundo command is undo rx command
//
// Description of rxundo command.
// This command function is return original path of contents in your trashbox.
// But this command is working rx command contents ONLY rx command case.
// Item was backuped by rx command, that has /[your-trashbox]/.prefix/NAME.rx.
// If not exist NAME.rx, do not return.
//
// Options
// rxundo [-v] NAME
// rxundo [-v] : view contents detail (backuped date)
//
// Usage
// rxundo NAME
//
// view menu as below,
// --> recovering a*
// (0) a original location : /original/location/a
// (1) aa original location : /original/location/aa
// you will select one from (0, 1) >
//
// you will select number of mercenary.

//{{{ import
import (
	"flag"
	"fmt"
	"os"
	// For rx command series utilty.
	"rx_common"
	// For rxls, rxundo command utilty.
	"rx_pv"
	"strconv"
)

//}}}

//{{{ func show_location(name, contents_str, date_str string, index int, rxFlagv bool, fileLenFlag bool)
// Show original path of backuped file or directory.
func show_location(name, contents_str, date_str string, index int, rxFlagv bool, fileLenFlag bool) {

	// Check number of contents, contents is 1 case not viewing.
	if fileLenFlag {

		// View original path of backuped file or directory.
		fmt.Printf("(%d) ", index)
		fmt.Print(name, " original location : ")
		if contents_str == "" {

			// Not logging original path case.
			fmt.Println("do not know original location, you will recover manually")
		} else {

			// Logging original path case.
			fmt.Println(contents_str)
		}

		// Check if logging backuped data.
		if date_str == "" {

			// Not logging backup date.
			date_str = "did not take the date log of rx executed"
		}
		// If add -v option, show backup date.
		if rxFlagv {

			fmt.Println("    date of rx executed :", date_str)
		}
	}
}

//}}}

//{{{ type Selector struct
// Selector structure (working like class)
type Selector struct {
	// The flag array is can/cannot return.
	CanRecovs []bool
	// The original path array (linkage CanRecovs)
	PrefixNames []string
	// Number of caontents (can recover)
	CanRecovsNum int
	// CanRecovs and PrefixNames array length
	Length int
}

//}}}

//{{{ func (s *Selector) Check_can_recov(contents_str string, index int)
// Store can or cannot recover flag
func (s *Selector) Check_can_recov(contents_str string, index int) {

	// Check can or cannot recover
	if contents_str == "" {

		// Store cannot recover flag
		s.CanRecovs[index] = false
	} else {

		// Store can recover flag
		s.CanRecovs[index] = true
	}
}

//}}}

//{{{ func (s *Selector) Get_can_recov_num()
// Get number of contents (can recover)
func (s *Selector) Get_can_recov_num() {

	// Initial value of can recovered contents
	var sum = 0
	for _, canRecov := range s.CanRecovs {

		// Each contents check can or cannot.
		if canRecov {

			// plus one number of contents (can recover)
			sum++
		}
	}

	// Set number of Can recover contents.
	s.CanRecovsNum = sum
}

//}}}

//{{{ func (s *Selector) Show_selector()
// Show selector
func (s *Selector) Show_selector() {

	// Get number of contents (can recover)
	s.Get_can_recov_num()
	var index = 0
	// Show menu
	fmt.Print("you will select one from (")
	for i, canRecov := range s.CanRecovs {

		// If can recover, standard out index of selection of file or directory.
		if canRecov {

			fmt.Print(i)
			index++

			// If not Last option, standard out cannma.
			if index < s.CanRecovsNum {
				fmt.Print(", ")
			}
		}
	}
	// Show your answer posttion .
	fmt.Print(") > ")
}

//}}}

//{{{ func (s *Selector) Get_index_from_selector() int
// Show index from selector,
// This function is showing selector recursively, in the time of inputting invalid option.
func (s *Selector) Get_index_from_selector() int {

	// show selector menu
	s.Show_selector()
	var (
		str   string
		index int
		err   error
	)

	// If number of contents (can recover) is 1, not showing menu.
	if s.Length == 1 {

		// Not view case:
		return 0
	}
	for {

		// Scanf your option
		if _, err := fmt.Scanf("%s", &str); err != nil {

			// If out Error case, standard out Error message.
			fmt.Println(err)
		}

		// Caset string to number.
		index, err = strconv.Atoi(str)
		if err != nil {

			// If failed cast, standard out Error message.
			fmt.Println(err)
		} else if index < 0 || index >= s.Length {

			// If invalid number case.
			fmt.Println("Error: this number is invalid range")
		} else if s.CanRecovs[index] {

			// If true value range and can recover case.
			break
		} else {

			// If true value range and not logging data option case.
			fmt.Println("Error: this option do not know original location")
		}
		// Review selector
		s.Show_selector()
	}

	// Get index of selector
	return index
}

//}}}

//{{{ func (s *Selector) Locate(filePrefixNames []string, pattern string, t *rx_common.TrashBox, rxFlagv bool)
// View menu
func (s *Selector) Locate(filePrefixNames []string, pattern string, index int, t *rx_common.TrashBox, rxFlagv bool) {

	// Get can recover contents original path and backuped date.
	contents_str, date_str := rx_pv.Get_prefix(filePrefixNames, pattern, t)
	// Show location view
	show_location(pattern, contents_str, date_str, index, rxFlagv, s.Length != 1)
	// Check ans store can or cannot recover contents flag.
	s.Check_can_recov(contents_str, index)
	// Store original path array to undo function
	s.PrefixNames[index] = contents_str
}

//}}}

//{{{ func undo(name string, contents_str string, t *rx_common.TrashBox)
// Controller UNDO rx command
func undo(name string, contents_str string, t *rx_common.TrashBox) {

	// Check original path is already exist
	if rx_common.Exist_file(contents_str) {

		// Standard out Error message.
		fmt.Println("Error:", contents_str, "is already exist")
	} else {

		// If all right case, recover
		if err := os.Rename(t.Get_fullPath_trash(name), contents_str); err != nil {

			// original path is empty case, Error message
			if contents_str == "" {

				contents_str = "you will input file name"
			}

			// View original path and Error message
			fmt.Println("target file or directory :", contents_str)
			fmt.Println(err)
		} else {

			// Remove logging data file  (/your-trashbox/.prefix/NAME.rx)
			os.Remove(t.Get_fullPath_prefix(rx_common.Get_prefix_filename(name)))
			// Standard out finish message
			fmt.Println("--> finish recovering to", contents_str)
		}
	}
}

//}}}

//{{{ func operation_of_undo(filename string, t *rx_common.TrashBox, rxFlagv bool)
// Main undo funciotn
func operation_of_undo(filename string, t *rx_common.TrashBox, rxFlagv bool) {

	var (
		i    int
		name string
	)

	fmt.Println("--> recovering", filename+"*")
	// Get filename array and array length.
	fileNames, fileNamesLen := rx_pv.Get_match(t.Get_trashBoxName(), filename)
	// Get logging filename array and array length.
	filePrefixNames, filePrefixNamesLen := rx_pv.Get_match(t.Get_filePrefixDir(), filename)

	// Not mattecd case.
	if fileNamesLen == 0 {

		// Not mattched case.
		fmt.Println("Error:", filename, "is not backuped")
	} else if filePrefixNamesLen == 0 {

		// Mattch and all contents is cannot recovered case.
		fmt.Println("Error: do not know", filename, " original location")
		fmt.Println("you will recover manually")
	} else {

		// True case, initialize Selector structure.
		s := &Selector{make([]bool, fileNamesLen), make([]string, fileNamesLen), 0, fileNamesLen}
		// Show filename and original path, backuped date (optional) in each file.
		for i, name = range fileNames {

			s.Locate(filePrefixNames, name, i, t, rxFlagv)
		}

		// Get index of recover file and show selector.
		index := s.Get_index_from_selector()
		// Undo function by index of recover file.
		undo(fileNames[index], s.PrefixNames[index], t)
	}
}

//}}}

//{{{ func main()
func main() {

	// Declare optoin flag.
	// rxFlagv is flag for -v option.
	// If this flag is true, show executed (backuped) date additionaly.
	var (
		rxFlagv bool
	)

	// Declare option name, initial value, description for help.
	flag.BoolVar(&rxFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&rxFlagv, "v", false, "it is same option, -V")
	// Analyze option flag from standard input, and set flag.
	// And set not flag standard input.
	flag.Parse()

	// Get TrashBox structure from rx command series configure file. ($HOME/.rx)
	t := rx_common.Get_trashBox_cfg()

	// Check detail veiw mode (select -v option)
	if rxFlagv {

		// Show target files
		fmt.Print("target files : ")
		fmt.Println(flag.Args())
	}

	// Execute UNDO rx command for recover contents to each name.
	for i := 0; i < flag.NArg(); i++ {

		operation_of_undo(flag.Args()[i], t, rxFlagv)
	}

}

//}}
