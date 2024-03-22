package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v3"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var entries []yamlObj
	err := yaml.Unmarshal(yml, &entries)
	if err != nil {
		return nil, err
	}
	m := buildMap(entries)
	return MapHandler(m, fallback), nil
}

type yamlObj struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func buildMap(entries []yamlObj) map[string]string {
	m := make(map[string]string)
	for _, entry := range entries {
		m[entry.Path] = entry.URL
	}
	return m
}
