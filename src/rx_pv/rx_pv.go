package rx_pv

import (
	"os"
	"path"
	"rx_common"
)

func Check_match_num(infos []os.FileInfo, pattern string) int {
	var sum = 0
	for _, info := range infos {
		var name = info.Name()
		if name != rx_common.FILE_PATH_DIR {
			if matched, _ := path.Match(pattern+"*", name); matched {
				sum += 1
			}
		}
	}
	return sum
}

func Check_match_names(infos []os.FileInfo, pattern string, num int) []string {
	var i = 0
	var names = make([]string, num)
	for _, info := range infos {
		var name = info.Name()
		if name != rx_common.FILE_PATH_DIR {
			if matched, _ := path.Match(pattern+"*", name); matched {
				names[i] = name
				i += 1
			}
		}
	}
	return names
}

func Check_match(infos []os.FileInfo, pattern string) ([]string, int) {
	num := Check_match_num(infos, pattern)
	names := Check_match_names(infos, pattern, num)
	return names, num
}

func Check_name(names []string, pattern string) string {
	for _, name := range names {
		if matched, _ := path.Match(rx_common.Get_prefix_filename(pattern), name); matched {
			return name
		}
	}
	return ""
}

func Check_location(name, prefixName, trashBoxName string) (string, string) {
	reader := rx_common.Generate_reader(rx_common.Get_fullPath_prefix(prefixName, trashBoxName))
	contents_str := rx_common.ReadLine(reader)
	date_str := rx_common.ReadLine(reader)
	return contents_str, date_str
}
