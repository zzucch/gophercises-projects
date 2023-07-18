package urlshort

import (
	"gopkg.in/yaml.v3"
	"net/http"
)

// MapHandler will return a http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if originalURL, ok := pathsToUrls[request.URL.Path]; ok {
			http.Redirect(writer, request, originalURL, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappers []URLMapper
	if err := yaml.Unmarshal(yml, &mappers); err != nil {
		return nil, err
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		for _, mapper := range mappers {
			if mapper.Path == request.URL.Path {
				http.Redirect(writer, request, mapper.URL, http.StatusMovedPermanently)
				return
			}
		}
		fallback.ServeHTTP(writer, request)
	}, nil
}

type URLMapper struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
