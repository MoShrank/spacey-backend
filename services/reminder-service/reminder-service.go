package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/moshrank/spacey-backend/config"
)

func task(cfg config.ConfigInterface) {
	deckServiceHostName := cfg.GetDeckServiceHostName()
	//learningServiceHostName := cfg.GetLearningServiceHostName()

	moritzUserID := "62629a030ea54c8763e9694c"

	decksUrl := "http://" + deckServiceHostName + "/decks?userID=" + moritzUserID

	deckRes, err := http.Get(decksUrl)
	if err != nil {
		panic(err)
	}

	type deck struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	decks := []deck{}
	err = json.NewDecoder(deckRes.Body).Decode(&decks)

	if err != nil {
		panic(err)
	}

	fmt.Println(decks)

	// fetch decks
	// fetch learning cards
	// send email
}

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	gocron.Every(10).Second().Do(task, config)

	<-gocron.Start()
}
