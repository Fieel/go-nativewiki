package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Page struct (aka interface)
type Page struct {
	Title string
	Body  []byte
}

// Used to write files on disk, in this case txt files
func (p *Page) save() error {
	filename := "./pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// Used in the handlers to load pages
func loadPage(title string) (*Page, error) {
	filename := "./pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Returns the HTML view template filled with data
// If there is not data or matching txt file, auto redirect to /edit
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Returns the HTML edit template filled with data, if available
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// Creates a new txt file in the /pages subdirectory
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("./templates/edit.html", "./templates/view.html"))

// Inject data in the HTML template and render it
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Create handler, based on given url path
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Main function, program starts here
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	fmt.Println("")
	fmt.Println("Server started... listening on post 8080")
	fmt.Println("URL: http://localhost:8080/view/")
	fmt.Println("")
	fmt.Println("*************************************")
	fmt.Println("Create a new page by visiting this url with the desired page name:")
	fmt.Println("localhost:8080/view/{newPageName}")
	fmt.Println("*************************************")
	fmt.Println("")
	fmt.Println("List of existing pages:")
	listExistingPages()
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Check for files in the /pages subfolder and prints their URLs
func listExistingPages() {
	dirname := "./pages"

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// cycle all files in the folder
	for _, file := range files {
		// extract name without extension and print it
		fmt.Println("http://localhost:8080/view/" + strings.Split(file.Name(), ".")[0])
	}
}
