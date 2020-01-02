package server

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Personal_Website/src/pages"
	"golang.org/x/crypto/acme/autocert"
)

var templates *template.Template

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeArticleHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			renderTemplate(w, "page_not_found", nil)
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.LoadPage(title)
	if err != nil {
		log.Print(err)
		renderTemplate(w, "page_not_found", nil)
		return
	}
	renderTemplate(w, "viewPage", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.LoadPage(title)
	if err != nil {
		p = &pages.Page{Title: title}
	}
	renderTemplate(w, "editPage", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &pages.Page{Title: title, Body: []byte(body)}
	err := p.SavePage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, templateString string, data interface{}) {
	err := templates.ExecuteTemplate(w, templateString+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getTemplate(templateString string) string {
	return os.Getenv("TEMPLATES_PATH") + templateString + ".html"
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil
}
func cacheTemplates() {
	templates = template.Must(template.ParseFiles(getTemplate("editPage"), getTemplate("viewPage"), getTemplate("page_not_found")))
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(os.Getenv("PAGES_PATH") + "/frontPage/index.html")
	templateVars := map[string]interface{}{
		"Groupdesk":    getModalTemplate("groupdesk"),
		"isidore":      getModalTemplate("isidore"),
		"FSAE":         getModalTemplate("FSAE"),
		"The6ix":       getModalTemplate("the6ix"),
		"pchacks":      getModalTemplate("pchacks"),
		"warriors":     getModalTemplate("warriors"),
		"lego_printer": getModalTemplate("lego_printer"),
		"gokart":       getModalTemplate("gokart"),
		"eve":          getModalTemplate("eve"),
		"sudoku":       getModalTemplate("sudoku"),
		"compEng":      getModalTemplate("compEng"),
		"arcade":       getModalTemplate("arcade"),
		"Path":         "/static/pages/frontPage/",
		"ExtPath":      "/static/external_files/",
	}
	t.Execute(w, templateVars)
}
func getModalTemplate(templateString string) string {
	body, err := ioutil.ReadFile(os.Getenv("TEMPLATES_PATH") + "portfolio_modals/" + templateString + ".html")
	if err != nil {
		log.Print("Failed to load template " + templateString)
		return ""
	}
	bodyStr := string(body)
	bodyStr = strings.Replace(bodyStr, "{{.Path}}", "/static/pages/frontPage/", -1)
	bodyStr = strings.Replace(bodyStr, "{{.ExtPath}}", "/static/external_files/", -1)
	return bodyStr
}

func handleRoutes(mux *http.ServeMux) {
	cacheTemplates()
	staticFs := http.FileServer(http.Dir(os.Getenv("STATIC")))
	mux.Handle("/static/", http.StripPrefix("/static/", staticFs))
	rootFs := http.FileServer(http.Dir(os.Getenv("ROOT_PATH")))
	mux.Handle("/sitemap.xml", rootFs)
	mux.HandleFunc("/", frontPage)
	/*mux.HandleFunc("/view/", makeArticleHandler(viewHandler))
	mux.HandleFunc("/edit/", makeArticleHandler(editHandler))
	mux.HandleFunc("/save/", makeArticleHandler(saveHandler))*/
}

func loadProductionEnvironment(server *http.Server, certManager *autocert.Manager) {
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
func loadDevEnvironment(server *http.Server) {
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func loadMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}
func loadServer(mux *http.ServeMux, https bool, port string) (*http.Server, *autocert.Manager) {
	port = ":" + port
	if https == true {
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("certs"),
		}
		return &http.Server{
			Addr:    port,
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}, &certManager
	}
	return &http.Server{
		Addr:    port,
		Handler: mux,
	}, &autocert.Manager{}
}

func RunServer() {
	mux := loadMux()
	handleRoutes(mux)
	if os.Getenv("PRODUCTION") == "TRUE" {
		server, certManager := loadServer(mux, true, "443")
		loadProductionEnvironment(server, certManager)
	} else {
		server, _ := loadServer(mux, false, "8080")
		loadDevEnvironment(server)
	}
}
