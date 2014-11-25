package rx_common

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
var (
	HOME          = os.Getenv("HOME")
	RX_EXTENDED   = ".rx"
	CFG_FILE      = HOME + "/" + RX_EXTENDED
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = ".prefix"
	newLineChar   = "\n"
)

//}}}

//{{{ func Exist_file(name string) bool
func Exist_file(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

//}}}

//{{{ func IsDirectory(name string) bool {
func IsDirectory(name string) bool {
	fInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fInfo.IsDir()
}

//}}}

//{{{ func Show_path(path string, rxFlagv bool)
func Show_path(path string, rxFlagv bool) {
	if rxFlagv {
		fmt.Println(path)
	}
}

//}}}

//{{{ func Get_filename_version(filename string, index int) string
func Get_filename_version(filename string, index int) string {
	var filenameVersion string
	if index == 0 {
		filenameVersion = filename
	} else {
		filenameVersion = filename + "." + strconv.Itoa(index)
	}
	return filenameVersion
}

//}}}

//{{{ func Get_prefix_filename_version(filename string, index int) string
func Get_prefix_filename_version(filename string, index int) string {
	return Get_filename_version(filename, index) + RX_EXTENDED
}

//}}}

//{{{ func Get_prefix_filename(filename string) string
func Get_prefix_filename(filename string) string {
	return Get_filename_version(filename, 0) + RX_EXTENDED
}

//}}}

//{{{ func Generate_reader(filename string) *bufio.Reader
func Generate_reader(filename string) *bufio.Reader {
	readFile, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readFile)
	return reader
}

//}}}

//{{{ func Read_rx_cfg() string
func Read_rx_cfg() string {
	reader := Generate_reader(CFG_FILE)
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

//}}}

//{{{ func ReadLine(reader *bufio.Reader) string
func ReadLine(reader *bufio.Reader) string {
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

//}}}

//{{{ func Read_file(filename string) string
func Read_file(filename string) string {
	contents, _ := ioutil.ReadFile(filename)
	return string(contents)
}

//}}}

//{{{ func Generate_writer(filename string) *bufio.Writer
func Generate_writer(filename string) *bufio.Writer {
	writeFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	return writer
}

//}}}

//{{{ func Write_rx_cfg(dirName string)
func Write_rx_cfg(dirName string) {
	os.Remove(CFG_FILE)
	dirVec := []byte(dirName)
	writer := Generate_writer(CFG_FILE)
	writer.Write(dirVec)
	writer.Flush()
}

//}}}

//{{{ func Write_file_cfg(path, trashBoxName string, index int)
func Write_file_cfg(path string, t *TrashBox, index int) {
	filename := filepath.Base(path)
	fullPath, _ := filepath.Abs(path)
	now := time.Now().String()[:19]
	contentsVec := []byte(fullPath + newLineChar + now)
	fileCfg := t.Get_filePrefixDir() + Get_prefix_filename_version(filename, index)
	writer := Generate_writer(fileCfg)
	writer.Write(contentsVec)
	writer.Flush()
}

//}}}

//{{{ type TrashBox struct
type TrashBox struct {
	TrashBoxName string
}

//}}}

//{{{ func (t *TrashBox) Get_trashBoxName() string
func (t *TrashBox) Get_trashBoxName() string {
	return t.TrashBoxName
}

//}}}

//{{{ func (t *TrashBox) Set_trashBoxName(trashBoxName string)
func (t *TrashBox) Set_trashBoxName(trashBoxName string) {
	t.TrashBoxName = trashBoxName
}

//}}}

//{{{ func (t *TrashBox) Get_filePrefixDir() string
func (t *TrashBox) Get_filePrefixDir() string {
	return t.Get_trashBoxName() + "/" + FILE_PATH_DIR + "/"
}

//}}}

//{{{ func (t *TrashBox) Get_fullPath_trash(filename string) string
func (t *TrashBox) Get_fullPath_trash(filename string) string {
	return t.Get_trashBoxName() + "/" + filename
}

//}}}

//{{{ func Get_fullPath_prefix(filename string) string
func (t *TrashBox) Get_fullPath_prefix(filename string) string {
	return t.Get_filePrefixDir() + filename
}

//}}}

//{{{ func (t *TrashBox) Make_trashBox()
func (t *TrashBox) Make_trashBox() {
	os.Mkdir(t.Get_trashBoxName(), 0777)
	os.Mkdir(t.Get_trashBoxName()+"/"+FILE_PATH_DIR, 0777)
}

//}}}

//{{{ func Get_trashBox_cfg() *TrashBox
func Get_trashBox_cfg() *TrashBox {
	var t *TrashBox
	if Exist_file(CFG_FILE) {
		t = &TrashBox{Read_rx_cfg()}
		if t.Get_trashBoxName() == "" {
			t.Set_trashBoxName(DEFAULT_DIR)
			Write_rx_cfg(t.Get_trashBoxName())
		}
	} else {
		t.Set_trashBoxName(DEFAULT_DIR)
		Write_rx_cfg(t.Get_trashBoxName())
	}
	t.Make_trashBox()
	return t
}

//}}}

//{{{ func Set_trashBox_cfg(trashBoxName string) *TrashBox
func Set_trashBox_cfg(trashBoxName string) *TrashBox {
	var t *TrashBox = &TrashBox{trashBoxName}
	if t.Get_trashBoxName() != "" {
		fmt.Println("INFO: you setted trash-box directory option, so trash-box clear options is not effective.")
		if Exist_file(t.Get_trashBoxName()) && IsDirectory(t.Get_trashBoxName()) == false {
			fmt.Println("Error: your option directory is exist as file type,")
			fmt.Println("       trash-box dir is setted as ./.trashbox.")
			t.Set_trashBoxName(DEFAULT_DIR)
		}
		Write_rx_cfg(t.Get_trashBoxName())
	} else if Exist_file(CFG_FILE) {
		t.Set_trashBoxName(Read_rx_cfg())
	} else {
		t.Set_trashBoxName(DEFAULT_DIR)
		Write_rx_cfg(t.Get_trashBoxName())
	}
	t.Make_trashBox()
	return t
}

//}}}
