package bot

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/monomonedula/randomwalk/pkg/state"
	txtr "github.com/monomonedula/randomwalk/pkg/texter"

	"github.com/tidwall/geodesic"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Sender interface {
	Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
}

func ErrorHandler(b Sender, texter txtr.Texter, log log.Logger) func(m *tb.Message, err error) {
	return func(m *tb.Message, err error) {
		msg, _ := json.MarshalIndent(m, "", "\t")
		log.Println(err, "msg: ", msg)
		b.Send(m.Sender, texter.TextFor(*m.Sender, "whoops"))
	}
}

func Help(b Sender, texter txtr.Texter) func(m *tb.Message) {
	return func(m *tb.Message) {
		txt := texter.TextFor(*m.Sender, "help")
		b.Send(m.Sender, txt)
	}
}

func AcceptText(b Sender, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error)) func(m *tb.Message) {
	return func(m *tb.Message) {
		AcceptRange(b, state, texter, errorHandler, m)
	}
}

func AcceptLocation(b Sender, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error)) func(m *tb.Message) {
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

func AcceptRange(b Sender, state state.State, texter txtr.Texter, errorHandler func(m *tb.Message, err error), m *tb.Message) {
	rng, err := WalkRangeFrom(m.Text)

	if err != nil {
		b.Send(m.Sender, texter.TextFor(*m.Sender, "bad_range_format"))
		return
	}

	user, err := state.UserData(m.Sender.ID)
	if err != nil {
		errorHandler(m, err)
		return
	}

	if IsEmptyLocation(user.Location) {
		b.Send(m.Sender, texter.TextFor(*m.Sender, "send_you_location_first"))
		return
	}
	location, err := RandomLocation(float64(user.Location.Lat), float64(user.Location.Lng), rng)
	if err != nil {
		errorHandler(m, err)
		return
	}

	_, err = b.Send(m.Sender, &location)
	if err != nil {
		errorHandler(m, err)
	}
}

func IsEmptyLocation(loc tb.Location) bool {
	return loc == tb.Location{}
}

type WalkRange struct {
	Min int
	Max int
}

func WalkRangeFrom(txt string) (WalkRange, error) {
	parts := strings.Split(txt, " ")
	if len(parts) != 2 {
		return WalkRange{}, errors.New("bad range format: " + txt)
	}
	minWalk, maxWalk := parts[0], parts[1]
	minWalkInt, err := strconv.Atoi(minWalk)
	if err != nil {
		return WalkRange{}, errors.New("bad range format: " + txt)
	}
	maxWalkInt, err := strconv.Atoi(maxWalk)
	if err != nil {
		return WalkRange{}, errors.New("bad range format: " + txt)
	}
	return WalkRange{Max: maxWalkInt, Min: minWalkInt}, nil
}

func RandomLocation(lat float64, lng float64, rng WalkRange) (tb.Location, error) {
	if rng.Min > rng.Max {
		return tb.Location{}, errors.New("minRadius cannot be greater that maxRadius")
	}
	var destLat, destLng float64

	distance := rng.Min + rand.Intn(rng.Max-rng.Min)
	azimuth := rand.Float64() * 360.0
	geodesic.WGS84.Direct(
		lat, lng,
		azimuth, float64(distance), &destLat, &destLng, nil,
	)
	return tb.Location{Lat: float32(destLat), Lng: float32(destLng)}, nil
}
