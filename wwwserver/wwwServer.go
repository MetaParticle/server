package wwwserver

import (
	"github.com/MetaParticle/metaparticle/logger"
	"github.com/MetaParticle/metaparticle/entity"

	"html/template"
	"net/http"

	"fmt"
	"io"
	"os"
	"strings"
	"bytes"
)

const (
	//The standard folder to serve webcontent from.
	//Update if new standard.
	LOCALROOT = "/var/www" 
	STDPORT   = 80

	//NOT USED!
	//404 Error message.
	E404TEXT = `
    <h1>Error 404: File not found!</h1>
    I could not find the file you were looking for.`
)

//Blocking function that serves content from the root folder on the specified port.
//Errors is returned on errchan.
//Negative port numbers will be converted to positive before use.
func ListenAndServe(root string, port int, errchan chan error) {
    
	if len(root) > 1 && strings.HasSuffix(root, "/") {
		root = (root)[:len(root)-1]
	}
	
	port = abs(port)

	//Registration for www handelers
	http.Handle("/res/", http.HandlerFunc( GenResourceServer(root) ))
	http.Handle("/", http.HandlerFunc( GenRootServer(root) ))

	//Listen on www port
	errchan <- http.ListenAndServe(":"+fmt.Sprint(port), nil)
	return
}

//Root server generates a function that serves html content from the root directory and all subdirectories except res/.
//That is the folder dedicated to resources and need special treatment.
//It loads a file from disk and treats it as a template and processes it based on the player object
//it gets from the GetPlayer(...) metod.
func GenRootServer(root string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		filename := getFilename(req, "index.html")
		
		logger.Logf(1, "Serving document to %s: %s", req.RemoteAddr, filename)
	        
		b, err := FileGet(root, filename)
		if len(b) < 1 || err != nil {
			logger.Logf(5, "len(b): %d, err: %s", len(b), err)
			Error(404, root, w, req)
			return
		}
		
		p := entity.GetPlayer(1337, "hunter02")
		err = ApplyTemplate(w, p, b)
		if err != nil {
			logger.Logf(2, "Failed to apply template: %s", err)
			Error(404, root, w, req)
			return
		}
	}
}

//Generates a function that serves resources.
//Resources includes anything not html or plain text.
func GenResourceServer(root string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		filename := root + "/" + req.URL.Path[1:]
		logger.Logf(2, "Serving resource to %s: %s", req.RemoteAddr, filename)
		http.ServeFile(w, req, filename)
	}
}

//Genertor to generate a handler that loads a template and applies a Page loaded from json.
func GenJSONServer(root string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		filename := getFilename(req, "index")
		
		logger.Logf(1, "Serving document to %s: %s", req.RemoteAddr, filename)
		page, err := jsonLoadPage(root, filename)
		p := entity.GetPlayer(1337, "hunter02")
		
		if err == nil {
			
			page.ApplyPageTemplate(p)
			template, err := FileGet(root, "/template.html")
			
			if err != nil {
				Error(404, root, w, req)
				return
			}
			err = ApplyTemplate(w, page, template)
			
			if err != nil {
				Error(404, root, w, req)
				return
			}
		} else {
			logger.Logf(2, "Error: %s", err.Error())
			Error(404, root, w, req)
		}
	}
}


func getFilename(req *http.Request, def string) (filename string) {
	filename = fmt.Sprint(req.URL)
	if strings.HasSuffix(filename, "/") {
		filename += def
	}
	return
}

func ApplyTemplate(w io.Writer, p interface{}, b []byte) (err error) {
	t := template.New("Page Template")
	t, err = t.Parse(string(b))
	if err == nil {
		t.Execute(w, p)
	}
	return err
}

/**
 * Based on the request, it loads a file into a buffer and returns it.
 * If the requested file ends with a / then index.html is appended.
 * 
 */
func FileGet(root string, filename string) (b []byte, err error) {
	buf := bytes.NewBufferString("")
	logger.Logf(4, "Opening file: %s", root + filename)
	r, err := os.Open(root + filename)

	if err == nil {
		defer r.Close()
		defer func() {
			if r := recover(); r != nil {
				logger.Logf(1,"Panic while reading file: %s\n", filename)
				b = []byte("")
			}
		}()
		buf.ReadFrom(r)
	}
	b = buf.Bytes()
	return
}

func Error(code int, root string, w http.ResponseWriter, req *http.Request) {
	filename := fmt.Sprintf("%s/%d.html", root, code)
	logger.Logf(3, "Serving Error to %s: %s", req.RemoteAddr, filename)
	http.ServeFile(w, req, filename)
	//w.Write([]byte(E404TEXT))
}

func abs(number int) int {
    if (number < 0) {
        number = -number
    }
    return number
}
