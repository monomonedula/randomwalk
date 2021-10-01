package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/monomonedula/randomwalk/pkg/state"
	txtr "github.com/monomonedula/randomwalk/pkg/texter"

	"github.com/go-redis/redis/v8"
	"github.com/tidwall/geodesic"
	tb "gopkg.in/tucnak/telebot.v2"
)

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
		Password: password, // no password set
		DB:       dbnum,    // use default DB
	}
}

func botToken() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Specify BOT_TOKEN env varialbe")
	}
	return token
}

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

	eh := errorHandler(b, texter, *log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile))
	b.Handle("/help", helpCommand(b, texter))
	b.Handle(tb.OnLocation, acceptLocation(b, &rs, texter, eh, 1))
	b.Handle(tb.OnText, acceptText(b, &rs, texter, eh))

	b.Start()
}

func errorHandler(b *tb.Bot, texter txtr.Texter, log log.Logger) func(m *tb.Message, err error) {
	return func(m *tb.Message, err error) {
		msg, _ := json.MarshalIndent(m, "", "\t")
		log.Println(err, "msg: ", msg)
		b.Send(m.Sender, texter.TextFor(*m.Sender, "whoops"))
	}
}

func helpCommand(b *tb.Bot, texter txtr.Texter) func(m *tb.Message) {
	return func(m *tb.Message) {
		txt := texter.TextFor(*m.Sender, "help")
		b.Send(m.Sender, txt)
	}
}

func acceptText(b *tb.Bot, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error)) func(m *tb.Message) {
	return func(m *tb.Message) {
		acceptRange(b, state, texter, errorHandler, m)
	}
}

func acceptLocation(b *tb.Bot, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error), nextState int) func(m *tb.Message) {
	return func(m *tb.Message) {
		data, err := state.UserData(m.Sender.ID)
		if err != nil {
			errorHandler(m, err)
			return
		}

		data.Location = *m.Location
		err = state.UpdateUserData(m.Sender.ID, data)
		if err != nil {
			errorHandler(m, err)
			return
		}

		txt := texter.TextFor(*m.Sender, "input_range")
		b.Send(m.Sender, txt)
	}
}

func acceptRange(b *tb.Bot, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error), m *tb.Message) {
	user, err := state.UserData(m.Sender.ID)
	if err != nil {
		errorHandler(m, err)
		return
	}

	parts := strings.Split(m.Text, " ")
	if len(parts) != 2 {
		b.Send(m.Sender, texter.TextFor(*m.Sender, "bad_range_format"))
	} else {
		minWalk, maxWalk := parts[0], parts[1]
		minWalkInt, err := strconv.Atoi(minWalk)
		if err != nil {
			_, err := b.Send(m.Sender, texter.TextFor(*m.Sender, "bad_range_format"))
			if err != nil {
				errorHandler(m, err)
			}
		}
		maxWalkInt, err := strconv.Atoi(maxWalk)
		if err != nil {
			_, err = b.Send(m.Sender, texter.TextFor(*m.Sender, "bad_range_format"))
			if err != nil {
				errorHandler(m, err)
			}
		}

		emptyLoc := tb.Location{}
		if user.Location == emptyLoc {
			b.Send(m.Sender, texter.TextFor(*m.Sender, "send_you_location_first"))
		} else {
			location, err := RandomLocation(float64(user.Location.Lat), float64(user.Location.Lng), minWalkInt, maxWalkInt)
			if err != nil {
				errorHandler(m, err)
				return
			}

			_, err = b.Send(m.Sender, &location)
			if err != nil {
				errorHandler(m, err)
			}
		}
	}

}

func RandomLocation(lat float64, lng float64, minRadius int, maxRadius int) (tb.Location, error) {
	if minRadius > maxRadius {
		return tb.Location{}, errors.New("minRadius cannot be greater that maxRadius")
	}
	var destLat, destLng float64

	distance := minRadius + rand.Intn(maxRadius-minRadius)
	azimuth := rand.Float64() * 360.0
	geodesic.WGS84.Direct(
		lat, lng,
		azimuth, float64(distance), &destLat, &destLng, nil,
	)
	return tb.Location{Lat: float32(destLat), Lng: float32(destLng)}, nil
}
