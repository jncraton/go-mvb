package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var root string

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found")
}

func pathToFiles(path string) []string {
	// dir/01/slug-name/ => dir/01-slug-name/
	re := regexp.MustCompile("(\\d+)/[a-z\\-]*/")
	path = re.ReplaceAllString(path, "$1-*/")

	var files []string

	files = append(files, root+path)
	files = append(files, root+path+"content.html")
	files = append(files, root+path+"content.md")

	return files
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Redirect paths without trailing slashes that aren't files
	if !strings.Contains(r.URL.Path, ".") && !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", 301)
	}

	var content []byte

	for _, file := range pathToFiles(r.URL.Path[1:]) {
		fmt.Printf("%s\n", file)

		content, _ = ioutil.ReadFile(file)

		if strings.HasSuffix(file, ".md") {
			content = blackfriday.MarkdownCommon(content)
		}

		if len(content) > 0 {
			break
		}
	}

	if len(content) > 0 {
		fmt.Fprintf(w, string(content))
	} else {
		notFound(w)
	}
}

func main() {
	root = "pages/"

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
