package cyoa

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	// Alternative : tpl = template.Must(template.New("").Parse(defaultHandlerTmpl)) // Il faut créer la variable defaultHandlerTmpl contenant l'ensemble du HTML au format string
	fp := path.Join("templates", "index.html")
	tpl = template.Must(template.ParseFiles(fp)) // .Must permet de déclencher une panic error automatiquement en cas de problème. C'est intéressant ici car si cette partie de code ne fonctionne pas alors il ne sert à rien de compiler l'ensemble du code.
	// .New crée un nouveau template avec le nom donné
	// .Parse permet de parser le string et le transformer en template
}

func JsonStory(file io.Reader) (Story, error) {
	decoder := json.NewDecoder(file)
	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fc func(r *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFc = fc
	}
}

func NewHandler(s Story, options ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFc}
	for _, option := range options {
		option(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFc func(r *http.Request) string
}

func defaultPathFc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}
	return path[1:] // "/intro" => "intro"

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFc(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			checkError(err)
			http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
