package main

import (
	"bufio"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var temp *template.Template

func posthandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/asciiart" {
			http.Error(w, "ERROR-404\nPage not found(", http.StatusNotFound)
			return
		}
		temp.ExecuteTemplate(w, "Template.html", nil)
	}

	if r.Method == "POST" {
		text := r.FormValue("text")
		banner := r.FormValue("banner")
		//check if user type in proper ascii art
		for _, v := range text {
			if !(v >= 32 && v <= 126) {
				http.Error(w, "ERROR-400\nBad request!", http.StatusBadRequest)
				return
			}
		}

		file, err := os.Open(banner + ".txt")

		if err != nil {
			http.Error(w, "ERROR-500\nInternal Server Error!!! \nPlease make sure you select a banner.", http.StatusInternalServerError)
			return
		}

		defer file.Close()
		//read the file
		Scanner := bufio.NewScanner(file)

		//identify the letters with ascii code
		var lines []string
		for Scanner.Scan() {
			lines = append(lines, Scanner.Text())
		}
		asciiChrs := make(map[int][]string)
		dec := 31

		for _, line := range lines {
			if line == "" {
				dec++
			} else {
				asciiChrs[dec] = append(asciiChrs[dec], line)

			}
		}
		var c = ""
		for i := 0; i < len(text); i++ {
			if text[i] == 92 && text[i+1] == 110 {
				c = PrintArt(text[:i], asciiChrs) + PrintArt(text[i+2:], asciiChrs)
			}
		}
		if !strings.Contains(text, "\\n") {
			c = PrintArt(text, asciiChrs)
		}
		pin := os.WriteFile("download.doc", []byte(c), 0644)
		if pin != nil {
			panic(pin)
		}
		pin1 := os.WriteFile("download.txt", []byte(c), 0644)
		if pin1 != nil {
			panic(pin1)
		}
		temp.ExecuteTemplate(w, "Template.html", c)
	}

}

func PrintArt(n string, y map[int][]string) string {
	//prints horizontally
	a := []string{}
	// prints horizontally
	for j := 0; j < len(y[32]); j++ {
		for _, letter := range n {
			a = append(a, y[int(letter)][j])
		}
		a = append(a, "\n")
	}
	b := strings.Join(a, "")
	return b
}

func download(w http.ResponseWriter, r *http.Request) {

	formatType := r.FormValue("fileformat")

	f, _ := os.Open("download." + formatType)
	defer f.Close()

	file, _ := f.Stat()
	fsize := file.Size()

	sfSize := strconv.Itoa(int(fsize))

	w.Header().Set("Content-Disposition", "attachment; filename=asciiresults."+formatType)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", sfSize)

	io.Copy(w, f)
}

func main() {
	temp = template.Must(template.ParseGlob("Template.html"))

	http.HandleFunc("/asciiart", posthandler)
	http.HandleFunc("/down", download)
	http.ListenAndServe(":8080", nil)
}
