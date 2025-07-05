package ws

import (
	"fmt"
	"lesson_server/constants"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type TimeData struct {
	Minute int
	Second int
	Flag   bool
}

func CreateNewTime(minute, second int, flag bool) *TimeData {
	return &TimeData{
		Minute: minute,
		Second: second,
		Flag:   flag,
	}
}

func (t *TimeData) GetDataTime() (*int, *int, *bool) {
	return &t.Minute, &t.Second, &t.Flag
}

func (t *TimeData) CountdownTime() {
	if t.Flag {
		if t.Minute > 0 || t.Second > 0 {
			if t.Second == 0 {
				t.Minute--
				t.Second = 59
			} else {
				t.Second--
			}
		}
	}
}

func GetTime(c *websocket.Conn, timeData *TimeData) {
	defer func() {
		c.Close()
		log.Println(constants.SuccCloseWS)
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		minute, second, flag := timeData.GetDataTime()

		send := map[string]interface{}{
			"time": fmt.Sprintf("%02d:%02d", *minute, *second),
			"flag": *flag,
		}

		if err := c.WriteJSON(send); err != nil {
			log.Println(constants.ErrInternalServer, err)
			return
		}
	}
}

var TimeLesson = CreateNewTime(45, 0, false)
var TimeOnly = CreateNewTime(2, 0, false)
var TimeTeam = CreateNewTime(1, 0, false)

func InitializeTimers() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			TimeLesson.CountdownTime()
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			TimeOnly.CountdownTime()
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			TimeTeam.CountdownTime()
		}
	}()
}

func GetLessonTime(c *websocket.Conn) {
	GetTime(c, TimeLesson)
}

func GetOnlyTime(c *websocket.Conn) {
	GetTime(c, TimeOnly)
}

func GetTeamTime(c *websocket.Conn) {
	GetTime(c, TimeTeam)
}
