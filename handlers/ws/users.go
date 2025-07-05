package ws

import (
	"encoding/json"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type SendUserDataStruct struct {
	Username   string `json:"username"`
	BimCoin    uint64 `json:"bim_coin"`
	Team       int    `json:"team"`
	TeamLeader bool   `json:"team_leader"`
}

var users []db_core.SendUserStruct

func GetUsers(c *websocket.Conn, db *db_core.DatabaseStruct) {
	defer func() {
		c.Close()
		log.Println(constants.SuccCloseWS)
	}()

	users, _ = db.GetUsers()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if constants.NewUsers {
			users, _ = db.GetUsers()
		}

		data, err := json.Marshal(users)
		if err != nil {
			log.Println(constants.ErrInternalServer)
			return
		}

		if err := c.WriteMessage(1, data); err != nil {
			log.Println(constants.ErrInternalServer, err)
			return
		}

		constants.NewUsers = false
	}
}
