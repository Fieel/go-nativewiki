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
	List  []string
	Body  []byte
}

// Available html templates
var templates = template.Must(template.ParseFiles(
	"./templates/home.html",
	"./templates/edit.html",
	"./templates/view.html",
	"./templates/header.html",
	"./templates/footer.html",
))

// Accept only following paths
// /
// /view/{pagename}
// /edit/{pagename}
// /save/{pagename}
var validPath = regexp.MustCompile("(^/(edit|save|view)/([a-zA-Z0-9]+))|(^/)$")

// Main function, program starts here
func main() {

	// Dynamic port (used by Heroku for example)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT must be set, using 5000 as default...")
		port = "5000"
	}

	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// For static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets/"))))

	fmt.Println("")
	fmt.Println("Server started... listening on port " + port)
	fmt.Println("URL: http://localhost:" + port + "/")
	fmt.Println("")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Create handler, based on given url path
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil && !(r.URL.Path == "//") {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[3])
	}
}

// Loads the homepage at /
func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage("Home")
	if err != nil {
		fmt.Println("homeHandler() error: ", err)
		return
	}
	renderTemplate(w, "home", p)
}

// Returns the HTML view template filled with data
// If there is not data or matching txt file, auto redirect to /edit
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		fmt.Println("viewHandler() error: ", err)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Returns the HTML edit template filled with data, if available
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		fmt.Println("editHandler() error: ", err)
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
		fmt.Println("saveHandler() error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// Inject data in the HTML template and render it
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println("renderTemplate() error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Used in the handlers to load pages
func loadPage(title string) (*Page, error) {
	filename := "./pages/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	list := fetchPageList()
	if err != nil {
		fmt.Println("loadPage() error: ", err)
		return nil, err
	}
	return &Page{Title: title, Body: body, List: list}, nil
}

// Used to write files on disk, in this case txt files
func (p *Page) save() error {
	filename := "./pages/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// reads the content of the /pages directory and returns a Slice of strings representing the page names
func fetchPageList() []string {
	dirname := "./pages"

	f, err := os.Open(dirname)
	if err != nil {
		fmt.Println("listExistingPages() error: ", err)
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Println("listExistingPages() error: ", err)
		log.Fatal(err)
	}

	// Stores the list of available pages to show on the homepage
	var pagesSlice []string
	for _, file := range files {
		pagesSlice = append(pagesSlice, strings.Split(file.Name(), ".")[0])
	}

	return pagesSlice
}

// Check for files in the /pages subfolder and prints their URLs
func listExistingPages(files []string) {
	// cycle all files in the folder
	for _, file := range files {
		// extract name without extension and print it
		fmt.Println("http://localhost:8080/view/" + strings.Split(file, ".")[0])
	}
}
