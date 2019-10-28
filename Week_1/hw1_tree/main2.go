package main

// FIXIT
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type dir struct {
	name  string
	files []string
	depth int
	last  bool
}

func (d *dir) setName(name string) {
	d.name = name
}
func (d dir) String() string {
	var filesymbol string

	tab := getTabs(d.depth, d.last)
	return fmt.Sprintf("%v%v%v // last:%v", tab, filesymbol, path.Base(d.name), d.last)
}

func getTabs(depth int, isLast bool) string {
	prefix := "├───"
	tabulation := ""

	if isLast {
		prefix = "└───"
	}

	for level := 0; level < depth-1; level++ {
		nextTab := "\t"
		if !isLast {
			nextTab = "│\t"

		}
		tabulation = tabulation + nextTab
	}

	return tabulation + prefix
}

var dirs []dir
var openedDirs map[int]bool

// type byName *dir

func dirTree(w io.Writer, p string, shouldprint bool) error {

	dirInfo(p, shouldprint)

	if shouldprint {
		err := fileInfo()
		if err != nil {
			fmt.Println("error occured")
		}
	}
	return nil
}

func dirInfo(p string, shouldprint bool) dir {
	info := readDir(p)

	// filtering out non dir entries from info
	if !shouldprint {
		n := 0
		for _, x := range info {
			if x.IsDir() {
				info[n] = x
				n++
			}
		}
		info = info[:n]
	}

	d := dir{}

	for i, e := range info {
		subdir := filepath.Join(p, e.Name())
		numberOfParentDirectories := len(strings.Split(subdir, "/"))
		if e.IsDir() {
			d.setName(subdir + "/")
			d.depth = numberOfParentDirectories
			if i == len(info)-1 {
				d.last = true
			}
			dirs = append(dirs, d)
			dirInfo(subdir, shouldprint)
		}
	}
	return d
}
func fileInfo() error {
	for i, d := range dirs {
		res := readDir(d.name)
		for _, v := range res {
			var file string
			if !v.IsDir() && v.Name() != ".directory" {
				fSize := strconv.FormatInt(v.Size(), 10)
				if fSize == "0" {
					fSize = "(empty)"
					file = v.Name() + " " + fSize
				} else {
					file = v.Name() + " (" + fSize + "b)"
				}
				dirs[i].files = append(dirs[i].files, file)
			}
		}
		sort.Strings(d.files)
	}
	return nil
}
func readDir(p string) []os.FileInfo {
	d, err := ioutil.ReadDir(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return d
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")

	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	for _, v := range dirs {
		fmt.Println(v)
	}
}
