package main

import (
	"flag"
	"fmt"
	"os"
	"rx_common"
	"strings"
)

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

	trashBoxName := rx_common.Set_trashBox_cfg("")
	remove_trash_box(trashBoxName)
	trashBoxName = rx_common.Set_trashBox_cfg(rx_common.DEFAULT_DIR)
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

	trashBoxName = rx_common.Set_trashBox_cfg(trashBoxName)

	if trashBoxClear && trashBoxCfg == false {
		clear_trash(trashBoxName)
	}
}
