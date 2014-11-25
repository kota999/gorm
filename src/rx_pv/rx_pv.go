package rx_pv

//{{{ import
import (
	"io/ioutil"
	"os"
	"path"
	"rx_common"
)

//}}}

//{{{ type FileInfo struct
type FileInfo struct {
	infos []os.FileInfo
}

//}}}

//{{{ func (f *FileInfo) Get_match_num(pattern string) int
func (f *FileInfo) Get_match_num(pattern string) int {
	var sum = 0
	for _, info := range f.infos {
		var name = info.Name()
		if name != rx_common.FILE_PATH_DIR {
			if matched, _ := path.Match(pattern+"*", name); matched {
				sum += 1
			}
		}
	}
	return sum
}

//}}}

//{{{ func (f *FileInfo) Get_match_names(pattern string, num int) []string
func (f *FileInfo) Get_match_names(pattern string, num int) []string {
	var i = 0
	var names = make([]string, num)
	for _, info := range f.infos {
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

//}}}

//{{{ func Get_match(infos []os.FileInfo, pattern string) ([]string, int)
func Get_match(dirName, pattern string) ([]string, int) {
	infos, _ := ioutil.ReadDir(dirName)
	var f *FileInfo = &FileInfo{infos}
	num := f.Get_match_num(pattern)
	names := f.Get_match_names(pattern, num)
	return names, num
}

//}}}

//{{{ func Get_name(names []string, pattern string) string
func Get_name(names []string, pattern string) string {
	for _, name := range names {
		if matched, _ := path.Match(rx_common.Get_prefix_filename(pattern), name); matched {
			return name
		}
	}
	return ""
}

//}}}

//{{{ func Get_location(name, prefixName string, t *rx_common.TrashBox) (string, string)
func Get_location(name, prefixName string, t *rx_common.TrashBox) (string, string) {
	reader := rx_common.Generate_reader(t.Get_fullPath_prefix(prefixName))
	contents_str := rx_common.ReadLine(reader)
	date_str := rx_common.ReadLine(reader)
	return contents_str, date_str
}

//}}}

//{{{ func Get_prefix(prefixNames, pattern string, t *rx_common.TrashBox) (string, string)
func Get_prefix(prefixNames []string, pattern string, t *rx_common.TrashBox) (string, string) {
	prefixName := Get_name(prefixNames, pattern)
	return Get_location(pattern, prefixName, t)
}

//}}}
