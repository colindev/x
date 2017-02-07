package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	x "log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/capability"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type token int

var (
	tokenKey token
	secret   Secret
)

// Secret is struct of secret file json
type Secret struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"project_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
}

func (secret Secret) getToken(r *http.Request) {

	//conf := jwt.Config{}

}

// Row ...
type Row struct {
	Data string
}

func main() {
	http.HandleFunc("/put", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)

		key := datastore.NewIncompleteKey(ctx, "gctest", nil)
		if _, err := datastore.Put(ctx, key, &Row{time.Now().String()}); err != nil {
			log.Errorf(ctx, "could not put into datastore: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debugf(ctx, r.URL.String())
		w.Write([]byte(r.URL.String()))

	}))

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)

		key := datastore.NewIncompleteKey(ctx, "gctest", nil)

		var row Row
		if err := datastore.Get(ctx, key, &row); err != nil {
			log.Errorf(ctx, "get fail from datastore: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Debugf(ctx, "%#v", row)
		w.Write([]byte(fmt.Sprintf("%#v", row)))

	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)

		query := &log.Query{AppLogs: true}
		if offset := r.FormValue("offset"); offset != "" {
			query.Offset, _ = base64.URLEncoding.DecodeString(strings.TrimSpace(offset))
		}

		res := query.Run(ctx)

		var display struct {
			Records []*log.Record
			Offset  string
		}

		for i := 0; i < 5; i++ {
			rec, err := res.Next()
			if err == log.Done {
				break
			}

			if err != nil {
				log.Errorf(ctx, "Reading log records: %v", err)
				break
			}

			display.Records = append(display.Records, rec)
			if i == 4 {
				display.Offset = base64.URLEncoding.EncodeToString(rec.Offset)
			}

			x.Println(rec.Offset)

		}

		log.Debugf(ctx, ">> %s", time.Now())

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.Encode(display)
	})

	http.HandleFunc("/capabilities", func(w http.ResponseWriter, r *http.Request) {

		ctx := appengine.NewContext(r)
		cap := r.FormValue("capability")
		mode := r.FormValue("mode")

		if !capability.Enabled(ctx, cap, mode) {
			w.Write([]byte(fmt.Sprintf("%s %s unavailable", cap, mode)))
			return
		}

		w.Write([]byte(fmt.Sprintf("%s %s available", cap, mode)))
	})

	appengine.Main()
}
