package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"math/rand"
	"time"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)


type Payload struct {
	Token       string
	TeamID      string
	TeamDomain  string
	ChannelID   string
	ChannelName string
	UserID      string
	UserName    string
	Command     string
	Text        string
}

func newPayloadByForm(form url.Values) *Payload {
	return &Payload{
		Token:       form.Get("token"),
		TeamID:      form.Get("team_id"),
		TeamDomain:  form.Get("team_domain"),
		ChannelID:   form.Get("channel_id"),
		ChannelName: form.Get("channel_name"),
		UserID:      form.Get("user_id"),
		UserName:    form.Get("user_name"),
		Command:     form.Get("command"),
		Text:        form.Get("text"),
	}
}

type ImageAttachment struct {
	ImageURL string `json:"image_url"`
}

type CommandReply struct {
	ResponseType string            `json:"response_type"`
	Text         string            `json:"text"`
	Attachments  []ImageAttachment `json:"attachments"`
}


func main() {
    godotenv.Load()

    algoliaCliAppId, envSet := os.LookupEnv("ALGOLIA_APP_ID")
    algoliaCliApiKey, envSet := os.LookupEnv("ALGOLIA_API_KEY")
    SlackCommandToken, envSet := os.LookupEnv("SLACK_CMD_TOKEN")
    if (!envSet) {
        log.Fatal("Please set the ALGOLIA_APP_ID, ALGOLIA_API_KEY and SLACK_CMD_TOKEN environment variables")
    }

	searchClient := algoliasearch.NewClient(algoliaCliAppId, algoliaCliApiKey)
	searchIndex := searchClient.InitIndex("classicprogrammerpaintings")
	searchParams := algoliasearch.Map{
		"typoTolerance": "min",
		"hitsPerPage": 1000,
	}

	seed := rand.NewSource(time.Now().Unix())
	rnd := rand.New(seed)
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		payload := newPayloadByForm(req.Form)
		if payload.Token != SlackCommandToken {
			log.Println("Token validation failed, received", payload.Token, "instead of", SlackCommandToken)
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			res, _ := searchIndex.Search(payload.Text, searchParams)
			if len(res.Hits) < 1 {
				io.WriteString(w, "No results for " + payload.Text)
			} else {
				w.Header().Set("Content-Type", "application/json")
				index := rnd.Intn(len(res.Hits))
				hit := res.Hits[index]
				imageUrl := fmt.Sprintf("%v", hit["ImageUrl"])
				title := fmt.Sprintf("%v", hit["Description"])
				a := ImageAttachment{
					ImageURL: imageUrl,
				}
				resp := CommandReply{
					Text:         title,
					Attachments:  []ImageAttachment{a},
					ResponseType: "in_channel",
				}
				jsonReply, _ := json.Marshal(resp)
				io.WriteString(w, string(jsonReply))
			}

		}
	})

    log.Println("Listening on 127.0.0.1:8000")
	http.ListenAndServe("127.0.0.1:8000", router)
}
