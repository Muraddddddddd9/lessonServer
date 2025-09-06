package main

import (
	"lesson_server/config"
	"lesson_server/constants"
	db_core "lesson_server/database"
	"lesson_server/handlers/api"
	"lesson_server/handlers/ws"
	"lesson_server/utils"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApiGroup(apiG fiber.Router, db *db_core.DatabaseStruct) {
	apiG.Post("/entry", func(c *fiber.Ctx) error {
		return api.Entry(c, db)
	})

	apiG.Post("/exit", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.Exit(c, db)
	})

	apiG.Post("/lesson", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.LessonChange(c, db)
	})

	apiG.Get("/questions", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.GetQuestionsOnly(c)
	})

	apiG.Get("/questions/team", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.GetQuestionsTeam(c)
	})

	apiG.Post("/answer", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.CheckAnswerOnly(c, db)
	})

	apiG.Post("/answer/team", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.CheckAnswerTeam(c, db)
	})

	apiG.Post("/time", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ChangeTime(c, db)
	})

	apiG.Post("/time/edit", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.EditTime(c, db)
	})

	apiG.Post("/presentation", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.EditPresentation(c, db)
	})

	apiG.Get("/presentation", Access(db, []string{constants.TeacherStatus, constants.StudentStatus}), func(c *fiber.Ctx) error {
		return api.GetPresentation(c, db)
	})

	apiG.Post("/clear", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ClearData(c, db)
	})

	apiG.Post("/clear/lesson", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.ClearStageLesson(c, db)
	})

	apiG.Post("/clear/student", Access(db, []string{constants.TeacherStatus}), func(c *fiber.Ctx) error {
		return api.DeleteStudent(c, db)
	})

	apiG.Get("/health", func(c *fiber.Ctx) error {
		utils.LogginAPI(c.Path(), c.Method(), fiber.StatusAccepted, c.IP(), nil, "health")
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "health",
		})
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
	app := fiber.New()
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
