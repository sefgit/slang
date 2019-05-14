package daemon

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func proxyRequestToOperator(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if len(path) == 0 || path[0] != '/' {
		w.WriteHeader(400)
		return
	}

	path = path[1:]
	sp := strings.Index(path, "/")
	handle := ""
	newPath := ""

	if sp == -1 {
		handle = path
		newPath = "/"
	} else {
		handle = path[:sp]
		newPath = path[sp:]
	}

	operator, ok := runningInstances[handle]
	if !ok {
		w.WriteHeader(404)
		return
	}

	newPort := operator.Port
	newURL := url.URL{}
	newURL.Scheme = "http"
	newURL.Host = "localhost:" + strconv.Itoa(newPort)
	newURL.Path = newPath
	newURL.RawQuery = r.URL.RawQuery

	req, _ := http.NewRequest(r.Method, newURL.String(), r.Body)
	for key := range r.Header {
		req.Header.Set(key, r.Header.Get(key))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	for key := range resp.Header {
		w.Header().Set(key, resp.Header.Get(key))
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
