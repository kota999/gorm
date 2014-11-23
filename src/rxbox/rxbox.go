package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	HOME          = os.Getenv("HOME")
	CFG_FILE      = HOME + "/.gorm"
	DEFAULT_DIR   = HOME + "/.trashbox"
	FILE_PATH_DIR = ".prefix"
	GORM_EXTENDED = ".gorm"
)

func make_trash_box(dirName string) {
	os.Mkdir(dirName, 0777)
	os.Mkdir(dirName+"/"+FILE_PATH_DIR, 0777)
}

func exist_file(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func isDirectory(name string) bool {
	fInfo, err := os.Stat(name)
	if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func read_cfg() string {
	contents, _ := ioutil.ReadFile(CFG_FILE)
	return string(contents)
}

func write_cfg(dirName string) {
	os.Remove(CFG_FILE)
	var writer *bufio.Writer
	dirVec := []byte(dirName)

	writeFile, _ := os.OpenFile(CFG_FILE, os.O_WRONLY|os.O_CREATE, 0600)
	writer = bufio.NewWriter(writeFile)
	writer.Write(dirVec)
	writer.Flush()
}

func set_trashBox_cfg(trashBoxName string) string {
	if trashBoxName != "" {
		fmt.Println("INFO: you setted trash-box directory option, so trash-box clear options is not effective.")
		if exist_file(trashBoxName) && isDirectory(trashBoxName) == false {
			fmt.Println("!! Error !!")
			fmt.Println("your option directory is existed as file type,")
			fmt.Println("trash box dir is setted as ./.trash.")
			trashBoxName = DEFAULT_DIR
		}
		write_cfg(trashBoxName)
	} else if exist_file(CFG_FILE) {
		trashBoxName = read_cfg()
	} else {
		trashBoxName = DEFAULT_DIR
		write_cfg(trashBoxName)
	}

	make_trash_box(trashBoxName)

	return trashBoxName
}

func clear_trash(trashBoxName string) {
	if err := os.RemoveAll(trashBoxName + "/"); err != nil {
		fmt.Println(err)
	}
	os.Mkdir(trashBoxName, 0777)
}

func remove_trash_box(trashBoxName string) {
	if err := os.RemoveAll(trashBoxName + "/"); err != nil {
		fmt.Println(err)
	}
}

func init_default() {
	fmt.Println("initializing trash-box configure and clear trash-box.")
	fmt.Println("INFO: initialized trash-box directory is $HOME/.trashbox.")
	trashBoxName := set_trashBox_cfg("")
	remove_trash_box(trashBoxName)
	trashBoxName = set_trashBox_cfg(DEFAULT_DIR)
}

func main() {

	var (
		trashBoxName  string
		trashBoxCfg   bool
		trashBoxClear bool
		defaultCfg    bool
	)

	flag.BoolVar(&trashBoxCfg, "box", false, "set trash box name")
	flag.BoolVar(&trashBoxClear, "C", false, "clear trash box")
	flag.BoolVar(&trashBoxClear, "c", false, "it is same option, -C")
	flag.BoolVar(&defaultCfg, "d", false, "initialize default setting (default dir: $HOME/.trashbox)")
	flag.Parse()

	if defaultCfg {
		init_default()
		os.Exit(1)
	}

	if trashBoxCfg == false || flag.NArg() == 0 {
		trashBoxName = ""
	} else {
		trashBoxName = flag.Args()[0]
		if strings.HasPrefix(trashBoxName, "-") {
			trashBoxName = ""
		}
	}

	trashBoxName = set_trashBox_cfg(trashBoxName)

	if trashBoxClear && trashBoxCfg == false {
		clear_trash(trashBoxName)
	}
}
