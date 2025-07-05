package ws

import (
	"lesson_server/constants"
	db_core "lesson_server/database"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

var lesson string

func GetStageLesson(c *websocket.Conn, db *db_core.DatabaseStruct) {
	defer func() {
		c.Close()
		log.Println(constants.SuccCloseWS)
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	setting, err := db.GetSetting()
	if err != nil {
		log.Println(constants.ErrInternalServer, err.Error())
	}
	lesson = setting.NowStageLesson

	for range ticker.C {
		if constants.NewLesson {
			setting, err := db.GetSetting()
			if err != nil {
				log.Println(constants.ErrInternalServer, err.Error())
			}
			lesson = setting.NowStageLesson
		}

		if err := c.WriteJSON(map[string]string{"lesson": lesson}); err != nil {
			log.Println(constants.ErrInternalServer, err.Error())
			return
		}

		constants.NewLesson = false
	}
}
