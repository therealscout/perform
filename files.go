package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/cagnosolutions/web"
)

type Node struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Children bool   `json:"children"`
	Type     string `json:"type"`
	State    string `json:"state"`
}

var files = web.Route{"GET", "/files", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "files.tmpl", nil)
	return
}}

var filesApi = web.Route{"GET", "/api/files", func(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("id")
	files, err := ioutil.ReadDir("./upload/company" + path)
	if err != nil {
		ajaxResponse(w, `{"id":"#","children":false}`)
		return
	}
	// var nodes []Node
	nodes := []Node{}
	for _, file := range files {
		if file.Name()[0] != '.' {
			n := Node{}
			n.Id = path + "/" + file.Name()
			n.Text = file.Name()
			n.Type = "file"
			if file.IsDir() {
				n.Type = "dir"
				n.Children = true
				n.State = "closed"
			}
			nodes = append(nodes, n)
		}
	}
	b, err := json.Marshal(nodes)
	respString := string(b)
	if err != nil {
		respString = `{"id":"#","Children":false}`
	}
	ajaxResponse(w, respString)
	return
}}

var newFolder = web.Route{"POST", "/api/mkdir", func(w http.ResponseWriter, r *http.Request) {
	path := "." + r.FormValue("id")
	if !strings.HasPrefix(path, "./upload/") {
		web.SetErrorRedirect(w, r, r.FormValue("redirect"), "Error creating new folder")
		return
	}
	folder := r.FormValue("folder")
	if folder == "" {
		web.SetErrorRedirect(w, r, r.FormValue("redirect"), "Error creating new folder")
		return
	}
	if err := os.MkdirAll(path+"/"+folder, 0755); err != nil {
		web.SetErrorRedirect(w, r, r.FormValue("redirect"), "Error creating new folder")
		return
	}
	web.SetSuccessRedirect(w, r, r.FormValue("redirect"), "Successfully created new folder")
	return
}}
