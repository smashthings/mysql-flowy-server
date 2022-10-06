package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Log(format string, a ...any){
	fmt.Println(fmt.Sprintf(format, a...))
}

func main() {
	Log("Starting standalone flowy server...")

	Log("Establishing DB...")
	EstablishDB()

	Log("Initialising router...")
	r := mux.NewRouter()

	r.HandleFunc("/set", func(res http.ResponseWriter, req *http.Request) {
		addCORSHeaders(res)
		set(res, req)
	})

	r.HandleFunc("/{id}", func(res http.ResponseWriter, req *http.Request) {
		addCORSHeaders(res)
		vars := mux.Vars(req)
		id := vars["id"]
		getOrDelete(id, res, req)
	})

	r.HandleFunc("/", index)

	http.Handle("/", r)
	Log("Starting service...")
	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		Log("Exiting main application with returned error - %s", err.Error())
	} else {
		Log("Exiting main application")
	}
}

func index(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "ok")
	// res.WriteHeader(http.StatusOK)
}

// Task is a single item.
type Task struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Checked  bool     `json:"checked"`
	Children []string `json:"children"`
}

func set(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		Log("200 - Responding to OPTIONS request")
		res.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		Log("405 - Bad HTTP Method")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if req.Header.Get("X-API-Key") == "" {
		Log("401 - No header with X-API-Key")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	apiKey := req.Header.Get("X-API-Key")

	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		Log("500 - Could not read body")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		Log("400 - Could not unmarshal, error: %v", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	Log("Parsed task {ID: %s, Text: %s, Checked: %b, Children: %s}", task.ID, task.Text, task.Checked, strings.Join(task.Children, ","))

	_, err = FetchKeyDB(task.ID, apiKey)
	if err == nil {
		err = DeleteKeyDB(task.ID, apiKey)
		if err != nil {
			Log("500 - Found existing key %s but not able to delete it before updating with error: %v", task.ID, err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = AddKeyDB(&task, apiKey)
	if err != nil {
		Log("500 - Failed to add task {ID: %s, Text: %s, Checked: %b, Children: %s} with error: %v", task.ID, task.Text, task.Checked, strings.Join(task.Children, ","), err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(res, "{ \"ok\": true }")
	res.Header().Set("Content-Type", "application/json")
}

func getOrDelete(id string, res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		Log("200 - Responding to OPTIONS request")
		res.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodGet && req.Method != http.MethodDelete {
		Log("405 - Bad HTTP Method")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if req.Header.Get("X-API-Key") == "" {
		Log("401 - No header with X-API-Key")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	apiKey := req.Header.Get("X-API-Key")

	res.Header().Set("Content-Type", "application/json")

	if req.Method == http.MethodDelete {
		if err := DeleteKeyDB(id, apiKey); err != nil {
			Log("500 - Could not remove ID %s from datastore, error: %v", id, err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(res, "\"ok\": true")
	} else {
		t, err := FetchKeyDB(id, apiKey)
		if err != nil {
			Log("500 - Could not get ID %s from datastore, error: %v", id, err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(t)
		if err != nil {
			Log("500 - Could not encode ID %s to json, error: %v", id, err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Write(data)
		res.WriteHeader(http.StatusOK)
	}
}

func addCORSHeaders(res http.ResponseWriter) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, X-API-Key")
	res.Header().Set("Access-Control-Allow-Credentials", "true")
}
