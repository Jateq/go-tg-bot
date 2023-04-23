package main

import (
	"context"
	"encoding/json"
	"github.com/Jateq/go-tg-bot/constants"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var Wg sync.WaitGroup
var Mutex sync.Mutex

func init() {
	Wg = sync.WaitGroup{}
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDebug(),
	}

	b, _ := bot.New(constants.BotAccess, opts...)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, Greetings)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/image", bot.MatchTypeExact, GetImage)
	b.Start(ctx)
}
func Greetings(ctx context.Context, b *bot.Bot, update *models.Update) {
	params := &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to my bot! To get started, type /image to receive a random photo.",
	}
	b.SendMessage(ctx, params)
}

func GetImage(ctx context.Context, b *bot.Bot, update *models.Update) {
	Wg.Add(1)
	url := make(chan string)
	go GetUrlFromAPI(&Wg, &Mutex, url)

	params := &bot.SendPhotoParams{
		ChatID: update.Message.Chat.ID,
		Photo:  &models.InputFileString{Data: <-url},
	}
	Wg.Wait()
	b.SendPhoto(ctx, params)
}

func GetUrlFromAPI(wg *sync.WaitGroup, m *sync.Mutex, url chan string) {
	Mutex.Lock()
	response, err := http.Get("https://api.unsplash.com/photos/random?client_id=" + constants.UnsplashAccess)
	if err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var photoLinks map[string]interface{}
	err = json.Unmarshal(body, &photoLinks)
	if err != nil {
		log.Fatal(err)
	}
	photoUrl := photoLinks["urls"].(map[string]interface{})["small"].(string)
	url <- photoUrl
	Mutex.Unlock()
	wg.Done()

}
