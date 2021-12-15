package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// The banner struct contains the input from the user (selecting banner), and output (returned ascii, strings)
type Banner struct {
	Ban1    string
	Ban2    string
	Ban3    string
	String1 string
	String2 string
}

var tpl *template.Template

// indexHandler writes the index template
func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := Banner{
		Ban1: "Shadow",
		Ban2: "Standard",
		Ban3: "Thinkertoy",
	}

	// handling any pages that are not the index or ascii-art (404 error)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		fmt.Fprintf(w, "Status 404: Page Not Found")
		return
	}
	tpl.ExecuteTemplate(w, "index.html", p)
}

// processHandler receives data from user and writes the ascii version of a string
func processHandler(w http.ResponseWriter, r *http.Request) {

	getban1 := r.FormValue("banner")
	getban2 := r.FormValue("banner")
	getban3 := r.FormValue("banner")
	tbox := r.FormValue("textbox")

	// handling bad request status code
	if len(getban1) == 0 && len(getban2) == 0 && len(getban3) == 0 || len(tbox) == 0 || strings.Contains(tbox, "£") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Status 400: Bad Request")
		return
	}

	testReturn := struct {
		ban1    string
		ban2    string
		ban3    string
		textbox string
	}{
		ban1:    getban1,
		ban2:    getban2,
		ban3:    getban3,
		textbox: tbox,
	}

	file, err := os.Open(testReturn.ban1 + ".txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "Status 500: Internal Server Error")
		return

	}

	scanned := bufio.NewScanner(file) // reading file
	scanned.Split(bufio.ScanLines)

	var lines []string

	for scanned.Scan() {
		lines = append(lines, scanned.Text())
	}

	file.Close()

	asciiChrs := make(map[int][]string)
	id := 31

	for _, line := range lines {
		if string(line) == "" {
			id++
		} else {
			asciiChrs[id] = append(asciiChrs[id], line)
		}
	}

	// convert textbox to bytes to figure out where linebreak is (10)
	b := []byte(testReturn.textbox)
	count := 0
	for _, num := range b {
		count++
		if num == 10 {
			break
		}
	}

	// checking if there is linebreak in string, returning the string seperated on 2 lines if there is
	// 2nd line is an empty string if there isnt a line break
	if strings.Contains(testReturn.textbox, "\n") {
		p := Banner{
			Ban1:    "Shadow",
			Ban2:    "Standard",
			Ban3:    "Thinkertoy",
			String1: Newline(testReturn.textbox[:count-2], asciiChrs),
			String2: Newline(testReturn.textbox[count:], asciiChrs),
		}
		if err := tpl.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		p := Banner{
			Ban1:    "Shadow",
			Ban2:    "Standard",
			Ban3:    "Thinkertoy",
			String1: Newline(testReturn.textbox, asciiChrs),
			String2: "",
		}
		if err := tpl.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Newline function returns the ascii art string horizontally
func Newline(n string, y map[int][]string) string {

	var empty string
	for j := 0; j < len(y[32]); j++ {
		var line string
		for _, letter := range n {
			line = line + string((y[int(letter)][j]))
		}
		empty += line + "\n"
		line = ""
	}
	return empty
}

// main runs the api(server) and its respective handlers
func main() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ascii-art", processHandler)
	http.ListenAndServe(":8080", nil)
}
