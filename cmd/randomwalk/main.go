package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/monomonedula/randomwalk/pkg/bot"
	"github.com/monomonedula/randomwalk/pkg/state"
	txtr "github.com/monomonedula/randomwalk/pkg/texter"

	"github.com/go-redis/redis/v8"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken(),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	opt := redisOptions()
	rdb := redis.NewClient(&opt)
	rs := state.NewRedisState(rdb, ctx)
	texter := txtr.BasicTexter()

	eh := bot.ErrorHandler(b, texter, *log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile))
	b.Handle("/start", bot.Help(b, texter))
	b.Handle("/help", bot.Help(b, texter))
	b.Handle(tb.OnLocation, bot.AcceptLocation(b, &rs, texter, eh))
	b.Handle(tb.OnText, bot.AcceptText(b, &rs, texter, eh))

	b.Start()
}

func redisOptions() redis.Options {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PSWD")
	db := os.Getenv("DATABASE")
	if db == "" {
		db = "0"
	}
	dbnum, err := strconv.Atoi(db)
	if err != nil {
		panic("Should not happen")
	}

	return redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbnum,
	}
}

func botToken() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Specify BOT_TOKEN env varialbe")
	}
	return token
}
