package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

/* 1. Parse l'url pour trouver le path
path := r.URL.Path[:]
2. Cherche dans la map l'url correspondante
url, ok := pathsToUrls[path]
3. Redirige le user vers la page
http.Redirect(w, r, newUrl, http.StatusSeeOther)
*/

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[:]
		url, ok := pathsToUrls[path]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		log.Fatalln(err)
	}
	pathToUrls := convertYamlIntoMap(pathUrls)
	return MapHandler(pathToUrls, fallback), err
}

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func convertYamlIntoMap(pathUrls []pathUrl) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pathUrl := range pathUrls {
		pathToUrls[pathUrl.Path] = pathUrl.URL
	}
	return pathToUrls
}
