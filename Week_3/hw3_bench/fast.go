package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Последний результат без easyjson
// BenchmarkSlow-4                1        1134139239 ns/op        18888584 B/op     195820 allocs/op
// BenchmarkFast-4                3         433348389 ns/op         3304042 B/op      60533 allocs/op
// PASS
// ok      coursera/Week_3/hw3_bench       5.549s

// Последний результат c easyjson
// BenchmarkSlow-4                1        1134967491 ns/op        18872920 B/op     195816 allocs/op
// BenchmarkFast-4               18          67973535 ns/op         2258911 B/op      13609 allocs/op
// PASS
// ok      coursera/Week_3/hw3_bench       3.825s

/*
   go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 Тестим производительность и сохраняем рузультат
   go tool pprof --web cpu.out -> смотрим результаты как svg построенный graphviz в браузере
*/
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	users := make([]User, 0)
	for scanner.Scan() {
		user := User{}
		_ = user.UnmarshalJSON([]byte(scanner.Text()))
		// fmt.Printf("%v %v\n", err, user)
		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	seenBrowsers := []string{}

	uniqueBrowsers := 0
	foundUsers := ""

	for i, user := range users {

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browserRaw := range browsers {
			browser := browserRaw
			if strings.Contains(browser, "Android") {
				isAndroid = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}

		}

		for _, browserRaw := range browsers {
			browser := browserRaw
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				notSeenBefore := true
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
				if notSeenBefore {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}
		email := strings.Replace(user.Email, "@", " [at] ", 1)
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
