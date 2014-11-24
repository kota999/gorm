package rx_common

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	HOME          = os.Getenv("HOME")
	RX_EXTENDED   = ".rx"
	CFG_FILE      = HOME + "/" + RX_EXTENDED
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = ".prefix"
	newLineChar   = "\n"
)

func Make_trashBox(dirName string) {
	os.Mkdir(dirName, 0777)
	os.Mkdir(dirName+"/"+FILE_PATH_DIR, 0777)
}

func Exist_file(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func IsDirectory(name string) bool {
	fInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func Generate_reader(filename string) *bufio.Reader {
	readFile, _ := os.OpenFile(filename, os.O_RDONLY, 0600)
	reader := bufio.NewReader(readFile)
	return reader
}

func Read_rx_cfg() string {
	reader := Generate_reader(CFG_FILE)
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

func ReadLine(reader *bufio.Reader) string {
	contents, _, _ := reader.ReadLine()
	return string(contents)
}

func Read_file(filename string) string {
	contents, _ := ioutil.ReadFile(filename)
	return string(contents)
}

func Generate_writer(filename string) *bufio.Writer {
	writeFile, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	writer := bufio.NewWriter(writeFile)
	return writer
}

func Write_rx_cfg(dirName string) {
	os.Remove(CFG_FILE)
	dirVec := []byte(dirName)
	writer := Generate_writer(CFG_FILE)
	writer.Write(dirVec)
	writer.Flush()
}

func Get_trashBox_cfg() string {
	var trashBoxName string
	if Exist_file(CFG_FILE) {
		if trashBoxName = Read_rx_cfg(); trashBoxName == "" {
			trashBoxName = DEFAULT_DIR
			Write_rx_cfg(trashBoxName)
		}
	} else {
		trashBoxName = DEFAULT_DIR
		Write_rx_cfg(trashBoxName)
	}
	Make_trashBox(trashBoxName)
	return trashBoxName
}

func Set_trashBox_cfg(trashBoxName string) string {
	if trashBoxName != "" {
		fmt.Println("INFO: you setted trash-box directory option, so trash-box clear options is not effective.")
		if Exist_file(trashBoxName) && IsDirectory(trashBoxName) == false {
			fmt.Println("Error: your option directory is exist as file type,")
			fmt.Println("       trash-box dir is setted as ./.trashbox.")
			trashBoxName = DEFAULT_DIR
		}
		Write_rx_cfg(trashBoxName)
	} else if Exist_file(CFG_FILE) {
		trashBoxName = Read_rx_cfg()
	} else {
		trashBoxName = DEFAULT_DIR
		Write_rx_cfg(trashBoxName)
	}
	Make_trashBox(trashBoxName)
	return trashBoxName
}

func Get_filePrefixDir(trashBoxName string) string {
	return trashBoxName + "/" + FILE_PATH_DIR + "/"
}

func Get_filename_version(filename string, index int) string {
	var filenameVersion string
	if index == 0 {
		filenameVersion = filename
	} else {
		filenameVersion = filename + "." + strconv.Itoa(index)
	}
	return filenameVersion
}

func Get_prefix_filename_version(filename string, index int) string {
	return Get_filename_version(filename, index) + RX_EXTENDED
}

func Get_prefix_filename(filename string) string {
	return Get_filename_version(filename, 0) + RX_EXTENDED
}

func Get_fullPath_trash(filename, trashBoxName string) string {
	return trashBoxName + "/" + filename
}

func Get_fullPath_prefix(filename, trashBoxName string) string {
	return Get_filePrefixDir(trashBoxName) + filename
}

func Write_file_cfg(path, trashBoxName string, index int) {
	filename := filepath.Base(path)
	fullPath, _ := filepath.Abs(path)
	now := time.Now().String()[:19]
	contentsVec := []byte(fullPath + newLineChar + now)
	fileCfg := Get_filePrefixDir(trashBoxName) + Get_prefix_filename_version(filename, index)
	writer := Generate_writer(fileCfg)
	writer.Write(contentsVec)
	writer.Flush()
}

func Show_path(path string, rxFlagv bool) {
	if rxFlagv {
		fmt.Println(path)
	}
}
