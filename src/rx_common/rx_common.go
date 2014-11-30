package rx_common

//
// This source file is rx_common package source for rx command series utility.
//
// This package is provide utility function and structure (working likely class).
//  - Basic enviroments
//  - File checker
//  - show file or path function
//  - Read/Write configure function
//  - Trashbox structure

//{{{ import
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

//}}}

//{{{ vars
//
// Basic environments for rx commnad series.
//
var (
	// Your HOME directory.
	HOME = os.Getenv("HOME")
	// File extended for rx command series.
	RX_EXTENDED = ".rx"
	// Configure filename for rx command series.
	CFG_FILE = HOME + "/" + RX_EXTENDED
	// Default trashbox directory path.
	DEFAULT_DIR = HOME + "/.trashbox"
	// This is directory name, for storing remove + backup logging data.
	FILE_PATH_DIR = ".prefix"
	// New line character for Linux/Unix.
	// rx commnad series is support Linux/Unix os system, so this is setted \n
	newLineChar = "\n"
)

//}}}

//{{{ func Exist_file(name string) bool
// Check exisit 'name' in your file system.
func Exist_file(name string) bool {

	// Existing name in your file system, err = nil.
	// Not existing name in your file system, err != nill.
	_, err := os.Stat(name)
	// If err == nill is true, name is exisit in your file system.
	// And err == nill is false, name is not exisit in your file system.
	return err == nil
}

//}}}

//{{{ func IsDirectory(name string) bool {
// Check file or directory type
func IsDirectory(name string) bool {

	// Get file infomation
	fInfo, err := os.Stat(name)
	// If file is not exisit, name is not directory.
	// Thus return false.
	if err != nil {
		return false
	}

	// Return if directory type from file info of name.
	return fInfo.IsDir()
}

//}}}

//{{{ func Show_path(path string, rxFlagv bool)
// Standard out file or directory name of target in condition of selected -v or -V.
func Show_path(path string, rxFlagv bool) {

	// If -v or -V option is selected, standard out file or directory name of target.
	if rxFlagv {
		fmt.Println(path)
	}
}

//}}}

//{{{ func Get_filename_version(filename string, index int) string
// Get file or directory version name.
func Get_filename_version(filename string, index int) string {

	// Declare variable of file or directory version name.
	var filenameVersion string
	// Check specificated version index is 0.
	if index == 0 {

		// If version is 0, file or directory name is default.
		filenameVersion = filename
	} else {

		// If version is not 0, file or directory name is FILENAME.i
		filenameVersion = filename + "." + strconv.Itoa(index)
	}
	return filenameVersion
}

//}}}

//{{{ func Get_prefix_filename_version(filename string, index int) string
// Get configure file version name of file or directory.
func Get_prefix_filename_version(filename string, index int) string {
	return Get_filename_version(filename, index) + RX_EXTENDED
}

//}}}

//{{{ func Get_prefix_filename(filename string) string
// Get configure filename of file or directory.
func Get_prefix_filename(filename string) string {

	// It is get configure file version of file or directory in condition of version 0.
	// Use this function for rx command series configure file.
	return Get_filename_version(filename, 0) + RX_EXTENDED
}

//}}}

//{{{ func Generate_reader(filename string) *bufio.Reader
// Generate instance for reading contents in file.
func Generate_reader(filename string) *bufio.Reader {

	// Open file of read file, read only and mode 600.
	readFile, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
	// Generate instance for reading contents in file from read file.
	reader := bufio.NewReader(readFile)
	return reader
}

//}}}

//{{{ func Read_rx_cfg() string
// Read trashbox name from configure file for rx command series. ($HOME.rx)
func Read_rx_cfg() string {

	// Generate instance for reading contents in configure file for rx command series.
	reader := Generate_reader(CFG_FILE)
	// Read your trashbox directory path.
	contents, _, _ := reader.ReadLine()
	// ReadLine 's return is []byte format, so cast string format for this cunftion return.
	return string(contents)
}

//}}}

//{{{ func ReadLine(reader *bufio.Reader) string
// Read one-line contents from instance for reading contents, ordering file top.
func ReadLine(reader *bufio.Reader) string {

	contents, _, _ := reader.ReadLine()
	// ReadLine 's return is []byte format, so cast string format for this function return.
	return string(contents)
}

//}}}

//{{{ func Read_file(filename string) string
// Read all conntens from file.
func Read_file(filename string) string {

	contents, _ := ioutil.ReadFile(filename)
	// ReadFile 's return is []byte format, so cast to string format for this function return.
	return string(contents)
}

//}}}

//{{{ func Generate_writer(filename string) *bufio.Writer
// Generate instance for writing contents in file.
func Generate_writer(filename string) *bufio.Writer {

	// Open file of write file, permission is Read and Write and 0600 mode,
	// If filename is not exist in your file system, create file.
	writeFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	// Generate instance of writing contents in file from write file.
	writer := bufio.NewWriter(writeFile)
	return writer
}

//}}}

//{{{ func Write_rx_cfg(dirName string)
// Write trashbox name in configure file for rx command series ($HOME.rx)
func Write_rx_cfg(dirName string) {

	// Remove configure file for rx command series,
	// before writing trashbox name.
	os.Remove(CFG_FILE)
	// Cast to []bytes format from string format, so instance of writing contents function input is []byte.
	dirVec := []byte(dirName)
	// Generate instance of write contents in file, this file is configure file for rx commnad series.
	writer := Generate_writer(CFG_FILE)
	// Write conntents in configure file buffer.
	writer.Write(dirVec)
	// Write configure file from buffer.
	writer.Flush()
}

//}}}

//{{{ func Write_file_cfg(path, trashBoxName string, index int)
// Write logging data in logging file, rx command.
func Write_file_cfg(path string, t *TrashBox, index int) {

	// Set your logging file fullpath.
	// Get filename from file or directory path of target.
	filename := filepath.Base(path)
	// Get fullpath from file or directory path of target.
	fullPath, _ := filepath.Abs(path)
	// Get executed date, format is bellow.
	// YYYY-mm-dd hh:mm:ss
	now := time.Now().String()[:19]
	// Format input conntents like bellow, tyoe format is []byte for write function.
	// 1: original filepath
	// 2: date of executed rx command
	contentsVec := []byte(fullPath + newLineChar + now)
	// Get fullpath of logging data of backuped file.
	// Fullpath of logging data is
	// Logging data path + File version name
	fileCfg := t.Get_filePrefixDir() + Get_prefix_filename_version(filename, index)
	writer := Generate_writer(fileCfg)
	// Write conntents in logging file buffer.
	writer.Write(contentsVec)
	// Write logging file from buffer.
	writer.Flush()
}

//}}}

//{{{ type TrashBox struct
//
// For operating trashbox file, structure format. (working like class)
//
type TrashBox struct {
	// Your trashbox name (default: $HOME/.trashbox
	TrashBoxName string
}

//}}}

//{{{ func (t *TrashBox) Get_trashBoxName() string
// Get your trashbox name
func (t *TrashBox) Get_trashBoxName() string {
	return t.TrashBoxName
}

//}}}

//{{{ func (t *TrashBox) Set_trashBoxName(trashBoxName string)
// Set your trashbox name, member variable of structure of TrashBox.
func (t *TrashBox) Set_trashBoxName(trashBoxName string) {
	t.TrashBoxName = trashBoxName
}

//}}}

//{{{ func (t *TrashBox) Get_filePrefixDir() string
// Get your logging data directory path from your trashbox directory.
// The logging data directory path is [your trashbox directory]/FILE_PATH_DIR
func (t *TrashBox) Get_filePrefixDir() string {
	return t.Get_trashBoxName() + "/" + FILE_PATH_DIR + "/"
}

//}}}

//{{{ func (t *TrashBox) Get_fullPath_trash(filename string) string
// Get filename in your trashbox path.
// This fullpath is path of backupping to.
func (t *TrashBox) Get_fullPath_trash(filename string) string {
	return t.Get_trashBoxName() + "/" + filename
}

//}}}

//{{{ func Get_fullPath_prefix(filename string) string
// Get logging data in your trashbox path.
func (t *TrashBox) Get_fullPath_prefix(filename string) string {
	return t.Get_filePrefixDir() + filename
}

//}}}

//{{{ func (t *TrashBox) Make_trashBox()
// Make directory your trashbox and logging data directory, permission mode 777.
func (t *TrashBox) Make_trashBox() {
	os.Mkdir(t.Get_trashBoxName(), 0777)
	os.Mkdir(t.Get_trashBoxName()+"/"+FILE_PATH_DIR, 0777)
}

//}}}

//{{{ func Get_trashBox_cfg() *TrashBox
// Get instance of structure of TrashBox.
// This instance of TrashBox is generated from your trashbox name.
// Your trashbox name is readed configure file of rx command series
// or default trashbox name ($HOME/.trashbox).
func Get_trashBox_cfg() *TrashBox {

	// Declare TrashBox structure.
	var t *TrashBox

	// Check if configure file is exist.
	if Exist_file(CFG_FILE) {

		// Exist configure file,
		// call constructer TrashBox structure from reading trashbox name in your configure file.
		t = &TrashBox{Read_rx_cfg()}
		// Check trashbox name is empty.
		if t.Get_trashBoxName() == "" {

			// Trashbox name is empty, set trashbox name to default ($HOME/.trashbox).
			t.Set_trashBoxName(DEFAULT_DIR)
			// Write default trashbox name in your configure file of rx command series ($HOME/.rx).
			Write_rx_cfg(t.Get_trashBoxName())
		}
	} else {

		// Not Exist configure file,
		// call constructer TrashBox structure, input default trashbox directory ($HOME/.trashbox).
		t.Set_trashBoxName(DEFAULT_DIR)
		// Write default trashbox name in your configure file of rx command series ($HOME/.rx).
		Write_rx_cfg(t.Get_trashBoxName())
	}

	// Make trashbox directory.
	// If trashbox directory is already existed, not working.
	t.Make_trashBox()
	return t
}

//}}}

//{{{ func Set_trashBox_cfg(trashBoxName string) *TrashBox
// Get instance of structure of TrashBox.
// This instance of TrashBox is generated from your trashbox name in your command option.
// Write your trashbox name in configure file of rx command series ($HOME/.rx).
func Set_trashBox_cfg(trashBoxName string) *TrashBox {

	// Get instance of structure of TrashBox, this is generated from your trashbox name in your command option.
	var t *TrashBox = &TrashBox{trashBoxName}
	// Check if your trahshbox name is empty.
	if t.Get_trashBoxName() != "" || t.Get_trashBoxName() != "/" {

		//  Standard out message, any trashbox options is not effective for writing trashbox name in configure file.
		fmt.Println("INFO: you setted trash-box directory option, so trash-box clear options is not effective.")
		// Check your if tashbox is already exsit and type is file type.
		if Exist_file(t.Get_trashBoxName()) && IsDirectory(t.Get_trashBoxName()) == false {

			// Already exsit and file type, your trashbox name is selected default trashbox name automaticaly.
			// Standard out message, not use your trashbox as trashbox directory and your trashbox is setted as default trashbox.
			fmt.Println("Error: your option directory is exist as file type,")
			fmt.Println("       trash-box dir is setted as $/.trashbox.")
			t.Set_trashBoxName(DEFAULT_DIR)
		}

		// Write your trashbox name in configure file.
		Write_rx_cfg(t.Get_trashBoxName())

		// Check if your trashbox name is empty and exist configure file.
	} else if Exist_file(CFG_FILE) {

		// Your trashbox name is empty and exist configure file.
		// Your trashbox name is readed from configure file.
		if trashBoxName := Read_rx_cfg(); trashBoxName == "" {

			// rx configure file is empty case, default dir set.
			t.Set_trashBoxName(DEFAULT_DIR)
		} else {

			// Usual case.
			t.Set_trashBoxName(trashBoxName)
		}
	} else {

		// Your trashbox name is empty and not exist configure file.
		// Your trashbox name is setted default trashbox name.
		t.Set_trashBoxName(DEFAULT_DIR)
		// Write default trashbox name in configure file.
		Write_rx_cfg(t.Get_trashBoxName())
	}

	// Make trashbox directory.
	// If your trashbox is already exist, not working.
	t.Make_trashBox()
	return t
}

//}}}
