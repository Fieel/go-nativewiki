# go-nativewebserver

Simple implementation of a web server using native Go components, no external frameworks, no frills or whatsoever.

Run the webserver as such:

```bash
go run wiki.go
```

You'll see the output in your terminale:
```bash
$ go run wiki.go

Server started... listening on post 8080
URL: http://localhost:8080/view/

*************************************
Create a new page by visiting this url with the desired page name:
localhost:8080/view/{newPageName}
*************************************

List of existing pages:
http://localhost:8080/view/FirstPage
```

#

## Dependencies

Of course you'll need the [Go cli](https://golang.org/dl/) to run the project.

This wiki.go file uses the following imports, which are necessary to run this web server:

```go
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
```

These packages should be found in your local ~/go/bin folder (or wherever your GOPATH points to).

#

## Pages 

This simple implementation uses two html templates and .txt data sources to create and show individual pages.

Pages are stored in the /pages subdirectory.

This project links to the [bulma.io](https://bulma.io) css library for quick styling.

### Page templates

Templates are stored in the /templates subdirectory.

Two layouts are used, view.html and edit.html and are respectively used by Go whenever accessing the /view/ or /edit/ subfolder from the browser.

### Page contents

The pages temselves are stored as .txt files.

The page name will become the HTML page title and the body of the file will populate the HTML body.