package rx_pv

//
// For rxls and rxundo utilty.
// This package is operation pattern match, search logging file, get file original prefix.
// This package 's function is used rxls, rxundo command.
// For use summary, ushow list of ndo remove and backup files and directories.
//

//{{{ import
import (
	"io/ioutil"
	"os"
	"path"
	// For rx command series utilty package
	"rx_common"
)

//}}}

//{{{ type FileInfo struct
//
// Operation file path, name structure (working like class)
// os.FileInfo 's wrapper
//
type FileInfo struct {
	infos []os.FileInfo
}

//}}}

//{{{ func (f *FileInfo) Get_match_num(pattern string) int
// Get number of the pattern match file or directory in the trashbox.
func (f *FileInfo) Get_match_num(pattern string) int {

	// Initial number of the pattern match file or directory in your trashbox.
	var sum = 0
	// Get each FileInfo from items in your trashbox.
	for _, info := range f.infos {

		// Get file or directory name from each item.
		var name = info.Name()
		// Ignore executed rx command logging data directory.
		if name != rx_common.FILE_PATH_DIR {

			// Check if pattern match file or directory name to your items in your trashbox.
			if matched, _ := path.Match(pattern+"*", name); matched {

				// Add 1 to number of the pattern match file or directory in your trashbox.
				sum += 1
			}
		}
	}

	// Return number of the pattern match file or directory in your trashbox.
	return sum
}

//}}}

//{{{ func (f *FileInfo) Get_match_names(pattern string, num int) []string
// Get file or directory names (string) of pattern match to items in your trashbox.
func (f *FileInfo) Get_match_names(pattern string, num int) []string {

	// Initial vaule of index of names.
	var i = 0
	// Allocate memories for names.
	var names = make([]string, num)
	// Get each FileInfo from items in your trahbox.
	for _, info := range f.infos {

		// Get file or directory name.
		var name = info.Name()
		// Ignore executed rx command logging data directory.
		if name != rx_common.FILE_PATH_DIR {

			// Check if pattern match file or directory name to your items in your trashbox.
			if matched, _ := path.Match(pattern+"*", name); matched {

				// Store pattern match file or directory name to name memory area.
				names[i] = name
				// Add 1 to index of the pattern match file or directory in your trashbox.
				i += 1
			}
		}
	}

	// Return names of the pattern match file or directory in your trashbox.
	return names
}

//}}}

//{{{ func Get_match(infos []os.FileInfo, pattern string) ([]string, int)
// Generate FileInfo structure, and call Get_match_num, Get_match_names
func Get_match(dirName, pattern string) ([]string, int) {

	// Get FileInfo s in your trashbox.
	infos, _ := ioutil.ReadDir(dirName)
	// Init New FileInfo structure from FileInfo s in your trashbox.
	var f *FileInfo = &FileInfo{infos}
	// Call Get_match_num function for getting pattern matched number of items.
	num := f.Get_match_num(pattern)
	// Call Get_match_names function for getting pattern matched names of items.
	names := f.Get_match_names(pattern, num)
	// Return number and names of the pattern matched items in your trashbox.
	return names, num
}

//}}}

//{{{ func Get_name(names []string, pattern string) string
// Get pattern matched name from inputted names.
func Get_name(names []string, pattern string) string {

	// Get each file or directory name from inputted names.
	for _, name := range names {

		// Check if pattern matched name to file or directory name.
		if matched, _ := path.Match(rx_common.Get_prefix_filename(pattern), name); matched {

			// Return pattern matched name (string)
			return name
		}
	}

	// If not pattern match, return empty string.
	return ""
}

//}}}

//{{{ func Get_location(name, prefixName string, t *rx_common.TrashBox) (string, string)
// Get location (file oringinal prefix) and executed rx date from executed rx logging file.
func Get_location(name, prefixName string, t *rx_common.TrashBox) (string, string) {

	// Generate Reader structure, used logging file.
	reader := rx_common.Generate_reader(t.Get_fullPath_prefix(prefixName))
	// Read fullpathes of matched file or directory.
	contents_str := rx_common.ReadLine(reader)
	// Read executed rx dates of matched file or directory.
	date_str := rx_common.ReadLine(reader)
	// Return fullpathes, dates.
	return contents_str, date_str
}

//}}}

//{{{ func Get_prefix(prefixNames, pattern string, t *rx_common.TrashBox) (string, string)
// Get location and executed rx command date of file or directory
// from existed file or directory pattern matched name.
func Get_prefix(prefixNames []string, pattern string, t *rx_common.TrashBox) (string, string) {
	prefixName := Get_name(prefixNames, pattern)
	return Get_location(pattern, prefixName, t)
}

//}}}
