package constants

const (
	ErrLoadEnv           = "Ошибка загрузки env"
	ErrInternalServer    = "Ошибка сервера"
	ErrConnectDatabase   = "Ошибка в подключении к БД: %v"
	ErrCheckPingDatabase = "Ошибка пинга БД: %v"
	ErrCloseDatabase     = "Ошибка в закрытии БД: %v"

	ErrTeacherAlreadyExist = "Учитель уже создан"
	ErrSettingAlreadyExist = "Настройки уже созданы"

	ErrInputValue     = "Не верный ввод данных"
	ErrUserNotFound   = "Пользователь не найден"
	ErrUserExit       = "Ошибка выхода"
	ErrNotSendTest    = "Нельзя отправиьт тест"
	ErrNoFullAnser    = "Ответьте на все вопросы"
	ErrAlreadyReplied = "Вы уже ответили на все вопросы"
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
)
