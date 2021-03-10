package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
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
	pathUrlsMap, err := parseYaml(yml)

	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(pathUrlsMap), fallback), nil
}

func parseYaml(yml []byte) ([]pathURL, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(yml, &pathUrls)

	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//     {
//        "path":"/google",
//        "url":"https://google.com/"
//     },
//     {
//        "path":"/github",
//        "url":"https://github.com/"
//     }
//  ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrlsMap, err := parseJSON(jsonData)

	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(pathUrlsMap), fallback), nil
}

// DBHandler will parse the provided DB data and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the data, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//     {
//        "path":"/google",
//        "url":"https://google.com/"
//     },
//     {
//        "path":"/github",
//        "url":"https://github.com/"
//     }
//  ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func DBHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	pathUrlsMap, err := parseDbData(db)

	if err != nil {
		return nil, err
	}

	return MapHandler(buildMap(pathUrlsMap), fallback), nil
}

func parseDbData(db *bolt.DB) ([]pathURL, error) {
	var pathUrls []pathURL

	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("PATH_URL"))
		bkt.ForEach(func(k []byte, v []byte) error {
			pathUrls = append(pathUrls, pathURL{Path: string(k), URL: string(v)})
			return nil
		})

		return nil
	})

	return pathUrls, err
}

func parseJSON(jsonData []byte) ([]pathURL, error) {
	var pathUrls []pathURL

	err := json.Unmarshal(jsonData, &pathUrls)

	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}

func buildMap(pathUrls []pathURL) map[string]string {
	pathUrlsMap := make(map[string]string)

	for _, v := range pathUrls {
		pathUrlsMap[v.Path] = v.URL
	}

	return pathUrlsMap
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}
