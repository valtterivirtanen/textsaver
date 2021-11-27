package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

const serverAddr = ":8080"

const textRoute = "/text"

var id int64
var textStore sync.Map

func main() {
	http.HandleFunc(textRoute, textHandler())

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func textHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getSingleText()(w, r)
		case http.MethodPost:
			createNewText()(w, r)
		}
	}
}

type textReq struct {
	Text string `json:"text"`
}

type textResp struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}

func createNewText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody textReq
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to decode request body;err=%v", err)
			return
		}

		newID := atomic.AddInt64(&id, 1)
		textStore.Store(newID, reqBody.Text)

		var respBody textResp
		respBody.Text = reqBody.Text
		respBody.ID = newID

		writeJSONResponse(w, respBody, http.StatusOK)
	}
}

func getSingleText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		textIDStr := r.URL.Query().Get("id")
		if len(textIDStr) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("invalid text id")
			return
		}

		textID, err := strconv.ParseInt(textIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to parse text id from query parameter;err=%v", err)
			return
		}

		text, ok := textStore.Load(textID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var respBody textResp
		respBody.Text, _ = text.(string)
		respBody.ID = textID

		writeJSONResponse(w, respBody, http.StatusOK)
	}
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, httpStatusCode int) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to encode response body;err=%v", err)
		return
	}

	w.WriteHeader(httpStatusCode)
	w.Header().Add("Content-Type", "application/json")
}
