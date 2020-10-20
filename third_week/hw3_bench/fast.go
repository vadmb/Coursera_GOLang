package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type User struct {
	Email    string
	Name     string
	Browsers []string
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]bool)
	uniqueBrowsers := 0
	index := 0
	user := &User{}
	fmt.Fprintln(out, "found users:")
	decoder := json.NewDecoder(file)
	for {
		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		browsers := user.Browsers

		isAndroid, uniqueValue := readBrowser(browsers, "Android", seenBrowsers)
		uniqueBrowsers += uniqueValue
		isMSIE, uniqueValue := readBrowser(browsers, "MSIE", seenBrowsers)
		uniqueBrowsers += uniqueValue

		if !(isAndroid && isMSIE) {
			index++
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", 1)
		fmt.Fprintf(out, "[%d] %s <%s>\n", index, user.Name, email)
		index++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}

func readBrowser(browsers []string, neededBrowser string, seen map[string]bool) (bool, int) {
	unique := 0
	is := false

	for _, browser := range browsers {
		if ok := strings.Contains(browser, neededBrowser); ok {
			is = true
			if _, exists := seen[browser]; !exists {
				seen[browser] = true
				unique++
			}
		}
	}
	return is, unique
}
