package main

import (
    "flag"
    "html/template"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "regexp"
)

var (
    addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

type Page struct {
    Title string
    Body  []byte
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    p := &Page{Title: "welcome", Body: []byte("yo")}
    renderTemplate(w, "main", p)
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
    p := &Page{Title: "welcome", Body: []byte("yo")}
    renderTemplate(w, "home", p)
}

func resumeHandler(w http.ResponseWriter, r *http.Request, title string) {
    p := &Page{Title: "resume", Body: []byte("yo")}
    renderTemplate(w, "resume", p)
}

func portfolioHandler(w http.ResponseWriter, r *http.Request, title string) {
    p := &Page{Title: "portfolio", Body: []byte("yo")}
    renderTemplate(w, "portfolio", p)
}

const TMPLDIR = "tmpl/"
var templates = template.Must(template.ParseFiles(TMPLDIR+"main.html", TMPLDIR+"home.html", TMPLDIR+"resume.html", TMPLDIR+"portfolio.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var validPath = regexp.MustCompile("^/(home|resume|portfolio|main)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[1])
    }
}

func main() {
    flag.Parse()

    fs := http.FileServer(http.Dir("resources"))
    http.Handle("/resources/", http.StripPrefix("/resources/", fs))

    http.HandleFunc("/", mainHandler)
    http.HandleFunc("/main/", mainHandler)
    http.HandleFunc("/home/", makeHandler(homeHandler))
    http.HandleFunc("/resume/", makeHandler(resumeHandler))
    http.HandleFunc("/portfolio/", makeHandler(portfolioHandler))

    if *addr {
        l, err := net.Listen("tcp", "127.0.0.1:0")
        if err != nil {
            log.Fatal(err)
        }
        err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
        if err != nil {
            log.Fatal(err)
        }
        s := &http.Server{}
        s.Serve(l)
        return
    }

    http.ListenAndServe(":8080", nil)
}