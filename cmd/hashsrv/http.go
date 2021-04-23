package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ancientlore/hashsrv/engine"
)

const (
	prefix    = "hashsrv-"
	prefixLen = len(prefix)
)

func root(w http.ResponseWriter, r *http.Request) {
	// read body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize engine and initial value
	eng := engine.New()
	if r.Method == "POST" || r.Method == "PUT" {
		eng.PushStack(b)
	}

	// initialize variables map
	eng.SetVariable("body", b)
	for k, v := range r.Header {
		if strings.HasPrefix(strings.ToLower(k), prefix) && len(v) > 0 {
			eng.SetVariable(k[prefixLen:], []byte(v[0]))
		}
	}

	// check for debug mode
	if r.URL.Query().Get("debug") != "" {
		eng.DebugMode = true
	}

	// process commands
	p := strings.TrimPrefix(r.URL.Path, "/")
	var arr []string
	if p != "" {
		arr = strings.Split(p, "/")
	} else {
		arr = make([]string, 0)
	}
	var rb []byte
	if len(arr) == 0 && r.Method == "GET" {
		rb = eng.Help()
	} else {
		rb, err = eng.Run(arr)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	_, err = w.Write(rb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
