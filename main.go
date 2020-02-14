package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	ascii "./ascii"
)

var temp *template.Template

func main() {

	http.HandleFunc("/", AsciiServer)
	//use - style, js file
	static := http.FileServer(http.Dir("public"))
	//secure, not access another files
	http.Handle("/public/", http.StripPrefix("/public/", static))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)

	}
}

func AsciiServer(w http.ResponseWriter, r *http.Request) {

	temp = template.Must(template.ParseGlob("template/*.gohtml"))

	if r.Method == "GET" {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err := temp.ExecuteTemplate(w, "index.gohtml", nil)
		if err != nil {
			errorHandler(w, r, 500)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)

			return
		}

	}

	if r.Method == "POST" {
		// http.Error()
		input := r.FormValue("input")
		fonts := r.FormValue("fonts")
		//newline
		newInput := strings.Replace(input, "\r\n", "\\n", -1)

		res, status := ascii.FontAscii(newInput, fonts)
		//	w.Write([]byte(res))

		if status == 500 || status == 404 {
			errorHandler(w, r, status)
			//http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		temp.ExecuteTemplate(w, "index.gohtml", res)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	temp = template.Must(template.ParseGlob("template/*.gohtml"))
	if status == http.StatusNotFound {
		temp.ExecuteTemplate(w, "error.gohtml", nil)
		//fmt.Fprint(w, "not found page, error 404")
	}
	if status == 500 {
		temp.ExecuteTemplate(w, "error.gohtml", nil)
	}
}

// 1 request client , html - input data -> request -> to backend
// 2 star os.Exec , output, take args- input, and return data client -> return data Html
