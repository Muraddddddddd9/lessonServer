package main

import (
	"fmt"
	"lesson_server/config"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/api"
	"lesson_server/handlers/ws"
	"lesson_server/utils"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApiGroup(apiG fiber.Router, db *db_core.DatabaseStruct) {
	apiG.Post("/entry_system", func(c *fiber.Ctx) error {
		return api.LoginSystem(c, db)
	})

	apiG.Post("/entry", AccessToken(), func(c *fiber.Ctx) error {
		return api.Entry(c, db)
	})

	apiG.Post("/exit", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.Exit(c, db)
	})

	apiG.Post("/lesson", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.LessonChange(c, db)
	})

	apiG.Get("/questions", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.GetQuestionsOnly(c, db)
	})

	apiG.Post("/answer", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.CheckAnswerOnly(c, db)
	})

	apiG.Post("/answer/team", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.CheckAnswerTeam(c, db)
	})

	apiG.Post("/score", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.ScoreTeam(c, db)
	})

	apiG.Post("/time", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ChangeTime(c, db)
	})

	apiG.Post("/time/edit", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.EditTime(c, db)
	})

	apiG.Post("/presentation", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.EditPresentation(c, db)
	})

	apiG.Get("/presentation", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.GetPresentation(c, db)
	})

	apiG.Post("/clear", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ClearData(c, db)
	})

	apiG.Post("/clear/lesson", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ClearStageLesson(c, db)
	})

	apiG.Post("/clear/student", AccessToken(), Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.DeleteStudent(c, db)
	})

	apiG.Get("/health", func(c *fiber.Ctx) error {
		utils.LogginAPI(c.Path(), c.Method(), fiber.StatusAccepted, c.IP(), nil, "health")
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "health",
		})
	})

	apiG.Get("/fast_test", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"fast_test": api.FastTest,
		})
	})

	apiG.Get("/get_file/:group", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		group := c.Params("group")
		return c.SendFile(fmt.Sprintf("pr2/file_%s.docx", group))
	})

	apiG.Get("/file/download/:group", AccessToken(), Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		group := c.Params("group")
		filename := fmt.Sprintf("pr2/file_%s.docx", group)

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Файл не существует",
			})
		}

		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=file_%s.docx", group))
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

		return c.SendFile(filename)
	})
}

func WsGroup(wsG fiber.Router, db *db_core.DatabaseStruct) {
	wsG.Get("/lesson", websocket.New(func(c *websocket.Conn) {
		ws.GetStageLesson(c, db)
	}))

	wsG.Get("/lesson/time", websocket.New(func(c *websocket.Conn) {
		ws.GetLessonTime(c)
	}))

	wsG.Get("/users", websocket.New(func(c *websocket.Conn) {
		ws.GetUsers(c, db)
	}))

	wsG.Get("/time/test/only", websocket.New(func(c *websocket.Conn) {
		ws.GetOnlyTime(c)
	}))

	wsG.Get("/time/test/team", websocket.New(func(c *websocket.Conn) {
		ws.GetTeamTime(c)
	}))
}

func main() {
	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
		BodyLimit:   50 * 1024 * 1024,
	})
	app.Static("/pr2", "./pr2")

	cfg, err := config.ConfigLoad()
	if err != nil {
		panic(err)
	}

	db, err := db_core.NewConnectDB()
	if err != nil {
		panic(err)
	}

	err = InitTeacher(db)
	if err != nil {
		panic(err)
	}

	err = InitSetting(db)
	if err != nil {
		panic(err)
	}

	ws.InitializeTimers()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     cfg.ORIGIN,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST",
	}))

	app.Use(swagger.New(swagger.Config{
		Next:     nil,
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
		Title:    "Fiber API documentation",
		CacheAge: 0,
	}))

	apiG := app.Group("/api")
	wsG := app.Group("/ws")

	ApiGroup(apiG, db)
	WsGroup(wsG, db)

	app.Listen(":8080")
}
