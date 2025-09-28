package constants

const (
	ErrLoadEnv           = "Ошибка загрузки env"
	ErrInternalServer    = "Ошибка сервера"
	ErrConnectDatabase   = "Ошибка в подключении к БД: %v"
	ErrCheckPingDatabase = "Ошибка пинга БД: %v"
	ErrCloseDatabase     = "Ошибка в закрытии БД: %v"

	ErrTeacherAlreadyExist = "Учитель уже создан"
	ErrSettingAlreadyExist = "Настройки уже созданы"

	ErrInputValue       = "Неверный ввод данных"
	ErrUserNotFound     = "Пользователь не найден"
	ErrUserExit         = "Ошибка выхода"
	ErrNotSendTest      = "Нельзя отправить ответ"
	ErrNoFullAnswer     = "Ответьте на все вопросы"
	ErrAlreadyReplied   = "Вы уже ответили на все вопросы"
	ErrScoreTeam        = "Вы уже отправили баллы"
	ErrClearStageLesson = "Ошибка в очистке этапа урока"
	ErrUserEntry        = "Ошибка входа в систему"
)

const (
	SuccCloseWS              = "WebSocket закрыт"
	SuccEntry                = "Успешный вход"
	SuccExit                 = "Успешный выход"
	SuccUpdateLessonStage    = "Этап урока успешно обновился"
	SuccGetAnswer            = "Ответ получен"
	SuccChangeTime           = "Успешное действие со временем"
	SuccUpdateIdPresentation = "ID презентации обновлена"
	SuccClearData            = "Данные очищены"
	SuccTeacherExist         = "Учитель создан"
	SuccSettingExist         = "Настройка создана"
	SuccCleatStageLesson     = "Этап урока очищен"
	SuccStudentDeleteById    = "Студент удалён"
	SuccScoreTeam            = "Баллы получены"
)
