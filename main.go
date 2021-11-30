package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	textstore "text_saver/text_store"
)

const textRoute = "/text"

func main() {
	var port string
	flag.StringVar(&port, "port", "80", "Exposed port")
	flag.Parse()

	http.HandleFunc(textRoute, textHandler())

	serverAddr := ":" + port
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func textHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getSingleText()(w, r)
		case http.MethodPost:
			createNewText()(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
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

		newID := textstore.Add(reqBody.Text)

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

		text, err := textstore.GetByID(textID)
		switch err {
		case nil:
		case textstore.ErrItemNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to get text from text store;err=%v", err)
			return
		}

		var respBody textResp
		respBody.Text = text
		respBody.ID = textID

		writeJSONResponse(w, respBody, http.StatusOK)
	}
}

func writeJSONResponse(w http.ResponseWriter, data interface{}, httpStatusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to encode response body;err=%v", err)
		return
	}

}
