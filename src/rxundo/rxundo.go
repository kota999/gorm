package main

//{{{ import
import (
	"flag"
	"fmt"
	"os"
	"rx_common"
	"rx_pv"
	"strconv"
)

//}}}

//{{{ func show_location(name, contents_str, date_str string, index int, rxFlagv bool, fileLenFlag bool)
func show_location(name, contents_str, date_str string, index int, rxFlagv bool, fileLenFlag bool) {
	if fileLenFlag {
		fmt.Printf("(%d) ", index)
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

//}}}

//{{{ type Selector struct
type Selector struct {
	CanRecovs    []bool
	PrefixNames  []string
	CanRecovsNum int
	Length       int
}

//}}}

//{{{ func (s *Selector) Check_can_recov(contents_str string, index int)
func (s *Selector) Check_can_recov(contents_str string, index int) {
	if contents_str == "" {
		s.CanRecovs[index] = false
	} else {
		s.CanRecovs[index] = true
	}
}

//}}}

//{{{ func (s *Selector) Get_can_recov_num()
func (s *Selector) Get_can_recov_num() {
	var sum = 0
	for _, canRecov := range s.CanRecovs {
		if canRecov {
			sum++
		}
	}
	s.CanRecovsNum = sum
}

//}}}

//{{{ func (s *Selector) Show_selector()
func (s *Selector) Show_selector() {
	s.Get_can_recov_num()
	var index = 0
	fmt.Print("you will select one from (")
	for i, canRecov := range s.CanRecovs {
		if canRecov {
			fmt.Print(i)
			index++
			if index < s.CanRecovsNum {
				fmt.Print(", ")
			}
		}
	}
	fmt.Print(") > ")
}

//}}}

//{{{ func (s *Selector) Get_index_from_selector() int
func (s *Selector) Get_index_from_selector() int {
	s.Show_selector()
	var (
		str   string
		index int
		err   error
	)

	if s.Length == 1 {
		return 0
	}
	for {
		if _, err := fmt.Scanf("%s", &str); err != nil {
			fmt.Println(err)
		}
		index, err = strconv.Atoi(str)
		if err != nil {
			fmt.Println(err)
		} else if index < 0 || index >= s.Length {
			fmt.Println("Error: this number is invalid range")
		} else if s.CanRecovs[index] {
			break
		} else {
			fmt.Println("Error: this option do not know original location")
		}
		s.Show_selector()
	}
	return index
}

//}}}

//{{{ func (s *Selector) Locate(filePrefixNames []string, pattern string, t *rx_common.TrashBox, rxFlagv bool)
func (s *Selector) Locate(filePrefixNames []string, pattern string, index int, t *rx_common.TrashBox, rxFlagv bool) {
	contents_str, date_str := rx_pv.Get_prefix(filePrefixNames, pattern, t)
	show_location(pattern, contents_str, date_str, index, rxFlagv, s.Length != 1)
	s.Check_can_recov(contents_str, index)
	s.PrefixNames[index] = contents_str
}

//}}}

//{{{ func undo(name string, contents_str string, t *rx_common.TrashBox)
func undo(name string, contents_str string, t *rx_common.TrashBox) {
	if rx_common.Exist_file(contents_str) {
		fmt.Println("Error:", contents_str, "is already exist")
	} else {
		if err := os.Rename(t.Get_fullPath_trash(name), contents_str); err != nil {
			if contents_str == "" {
				contents_str = "you will input file name"
			}
			fmt.Println("target file or directory :", contents_str)
			fmt.Println(err)
		} else {
			os.Remove(t.Get_fullPath_prefix(rx_common.Get_prefix_filename(name)))
			fmt.Println("--> finish recovering to", contents_str)
		}
	}
}

//}}}

//{{{ func operation_of_undo(filename string, t *rx_common.TrashBox, rxFlagv bool)
func operation_of_undo(filename string, t *rx_common.TrashBox, rxFlagv bool) {
	var (
		i    int
		name string
	)
	fmt.Println("--> recovering", filename+"*")
	fileNames, fileNamesLen := rx_pv.Get_match(t.Get_trashBoxName(), filename)
	filePrefixNames, filePrefixNamesLen := rx_pv.Get_match(t.Get_filePrefixDir(), filename)
	if fileNamesLen == 0 {
		fmt.Println("Error:", filename, "is not backuped")
	} else if filePrefixNamesLen == 0 {
		fmt.Println("Error: do not know", filename, "'s original location")
		fmt.Println("you will recover manually")
	} else {
		s := &Selector{make([]bool, fileNamesLen), make([]string, fileNamesLen), 0, fileNamesLen}
		for i, name = range fileNames {
			s.Locate(filePrefixNames, name, i, t, rxFlagv)
		}
		index := s.Get_index_from_selector()
		undo(fileNames[index], s.PrefixNames[index], t)
	}
}

//}}}

//{{{ func main()
func main() {

	var (
		rxFlagv bool
	)

	flag.BoolVar(&rxFlagv, "V", false, "show file name before throw away")
	flag.BoolVar(&rxFlagv, "v", false, "it is same option, -V")
	flag.Parse()

	t := rx_common.Get_trashBox_cfg()

	if rxFlagv {
		fmt.Print("target files : ")
		fmt.Println(flag.Args())
	}
	for i := 0; i < flag.NArg(); i++ {
		operation_of_undo(flag.Args()[i], t, rxFlagv)
	}

}

//}}
