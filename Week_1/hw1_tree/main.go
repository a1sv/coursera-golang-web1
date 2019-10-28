package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type byFilename []os.FileInfo

func (finfo byFilename) Len() int {
	return len(finfo)
}

func (finfo byFilename) Swap(i, j int) {
	finfo[i], finfo[j] = finfo[j], finfo[i]
}

func (finfo byFilename) Less(i, j int) bool {
	return finfo[i].Name() < finfo[j].Name()
}

func filterDirs(fileList []os.FileInfo) []os.FileInfo {
	var filteredList []os.FileInfo

	for _, file := range fileList {
		if file.IsDir() {
			filteredList = append(filteredList, file)
		}
	}

	return filteredList
}

var ignoredFiles = []string{
	".DS_Store",
	".gitignore",
	".directory",
	".vscode",
}

func checkIgnoredFile(file string) bool {
	for _, i := range ignoredFiles {
		if i == file {
			return true
		}
	}
	return false
}

func printDir(output io.Writer, path string, openedDirs map[int]bool, depth int, printFiles bool) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, _ := file.Readdir(-1)
	if !printFiles {
		fileInfo = filterDirs(fileInfo)
	}
	openedDirs[depth] = true
	sort.Sort(byFilename(fileInfo))

	for i, file := range fileInfo {
		fileName := file.Name()
		if checkIgnoredFile(fileName) {
			continue
		}

		isLast := i+1 == len(fileInfo)
		openedDirs[depth] = !isLast
		symbol := addTab(depth, isLast, openedDirs)

		if file.IsDir() {
			fmt.Fprintf(output, "%s%s\n", symbol, fileName)
			nextPath := fmt.Sprintf("%s/%s", path, fileName)

			printDir(output, nextPath, openedDirs, depth+1, printFiles)
		} else if printFiles {
			fileSize := printSize(file.Size())

			fmt.Fprintf(output, "%s%s%s\n", symbol, fileName, fileSize)
		}
	}
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	openedDirs := map[int]bool{}

	stats, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	if mode := stats.Mode(); mode.IsDir() {
		printDir(output, path, openedDirs, 0, printFiles)
	}

	return nil
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
}

func printSize(size int64) string {
	if size == 0 {
		return " (empty)"
	}

	return fmt.Sprintf(" (%vb)", size)
}

func addTab(depth int, isLast bool, openedDirs map[int]bool) string {
	symbol := "├───"
	tab := ""

	if isLast {
		symbol = "└───"
	}

	for level := 0; level < depth; level++ {
		dirOpened, ok := openedDirs[level]
		if ok {
			nextTab := "\t"
			if dirOpened {
				nextTab = "│\t"
			}

			tab = tab + nextTab
		}
	}

	return tab + symbol
}
