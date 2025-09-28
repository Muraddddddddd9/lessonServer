package db

import (
	"context"
	"database/sql"
	"fmt"
	"lesson_server/config"
	"lesson_server/constants"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseStruct struct {
	db *sql.DB
}

var (
	TableUsers         = "users"
	TableSetting       = "setting"
	TableAnswerTestTwo = "answer_test_two"
	TableSuperUser     = "superuser"
)

// Подключение к БД
func NewConnectDB() (*DatabaseStruct, error) {
	cfg, err := config.ConfigLoad()
	if err != nil {
		log.Printf(constants.ErrLoadEnv)
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf(constants.ErrConnectDatabase, err)
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf(constants.ErrCheckPingDatabase, err)
		return nil, err
	}

	return &DatabaseStruct{db: db}, nil
}

// Закрытие БД
func (d *DatabaseStruct) Close() error {
	return d.db.Close()
}

// GET Получение всех пользователей со всеми данными
func (d *DatabaseStruct) GetDataUsers(arg ...any) ([]UserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE status = ?", TableUsers)
	rows, _ := d.db.QueryContext(ctx, query, constants.StudentStatus)
	defer rows.Close()

	var users []UserStruct

	for rows.Next() {
		var user UserStruct
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Password,
			&user.Status,
			&user.BimCoin1,
			&user.BimCoin2,
			&user.Team,
			&user.TestFirst,
			&user.TimeTest,
			&user.BimTotal,
		)

		users = append(users, user)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

// GET Получение всех пользователей
func (d *DatabaseStruct) GetUsers() ([]SendUserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT id, name, bim_coin1, bim_coin2, time_test, team, bim_total FROM %s WHERE status = ? ORDER BY bim_total DESC", TableUsers)
	rows, err := d.db.QueryContext(ctx, query, constants.StudentStatus)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var users []SendUserStruct
	for rows.Next() {
		var user SendUserStruct
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.BimCoin1,
			&user.BimCoin2,
			&user.TimeTest,
			&user.Team,
			&user.BimTotal,
		)
		if err != nil {
			log.Printf("Функция GetUsers(), ошибка данных: %v", user)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

func (d *DatabaseStruct) GetTeamUser(team int) ([]SendUserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT name FROM %s WHERE team = ? ", TableUsers)
	rows, err := d.db.QueryContext(ctx, query, team)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var users []SendUserStruct
	for rows.Next() {
		var user SendUserStruct
		err = rows.Scan(
			&user.Name,
		)
		if err != nil {
			log.Printf("Функция GetTeamUser(), ошибка данных: %v", user)
			continue
		}

		users = append(users, user)
	}

	return users, nil
}

// GET Получение пользователя
func (d *DatabaseStruct) GetOneUser(where string, arg ...any) (*UserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", TableUsers, where)
	var user UserStruct
	err := d.db.QueryRowContext(ctx, query, arg...).Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.Status,
		&user.BimCoin1,
		&user.BimCoin2,
		&user.Team,
		&user.TestFirst,
		&user.TimeTest,
		&user.BimTotal,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GET Получение настроек
func (d *DatabaseStruct) GetSetting() (*SettingStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT now_stage_lesson, id_presentation, test_team_first, test_team_second FROM %s", TableSetting)
	var setting SettingStruct
	err := d.db.QueryRowContext(ctx, query).Scan(&setting.NowStageLesson, &setting.IdPresentation, &setting.TestTeamFirst, &setting.TestTeamSecond)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// GET Получение Ответов на второй тест
func (d *DatabaseStruct) GetAnswerTaskTwo(team int) (*AnswerTaskTwo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM %s WHERE team = ?", TableAnswerTestTwo)
	var setting AnswerTaskTwo
	err := d.db.QueryRowContext(ctx, query, team).Scan(
		&setting.Team,
		&setting.ConflictElements,
		&setting.Visualization1,
		&setting.Visualization2,
		&setting.Visualization3,
		&setting.MatrixFile,
		&setting.CollisionsCount,
		&setting.Priority,
		&setting.RiskAssessment,
		&setting.Consequences,
		&setting.Problems,
		&setting.Specialists,
		&setting.Stage,
	)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}

// GET Получение суперпользователя
func (d *DatabaseStruct) GetSuperUser(username string) (*SuperUserSturct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT password FROM %s WHERE username = ?", TableSuperUser)
	var user SuperUserSturct
	err := d.db.QueryRowContext(ctx, query, username).Scan(
		&user.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// INSERT Добавление новых пользователей
func (d *DatabaseStruct) InsertUser(newUser UserStruct) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (name, password, status, bim_coin1, bim_coin2, team, test_first, bim_total) VALUES (?,?,?,?,?,?,?,?)", TableUsers)
	result, err := d.db.ExecContext(
		ctx,
		query,
		newUser.Name,
		newUser.Password,
		newUser.Status,
		newUser.BimCoin1,
		newUser.BimCoin2,
		newUser.Team,
		newUser.TestFirst,
		newUser.BimTotal,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// INSERT Добавление настройки
func (d *DatabaseStruct) InsertSetting(newSetting SettingStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO %s (now_stage_lesson, id_presentation, test_team_first, test_team_second) VALUES (?,?,?,?)", TableSetting)
	_, err := d.db.ExecContext(
		ctx,
		query,
		newSetting.NowStageLesson,
		newSetting.IdPresentation,
		newSetting.TestTeamFirst,
		newSetting.TestTeamSecond,
	)
	if err != nil {
		return err
	}

	return nil
}

// INSERT/UPDATE Добавление ответа студента
func (d *DatabaseStruct) UpsertAnswerTaskTwo(team int, keys []string, values ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	qs := make([]string, len(keys))
	for i := range qs {
		qs[i] = "?"
	}

	updateParts := make([]string, len(keys))
	for i, k := range keys {
		updateParts[i] = fmt.Sprintf("%s=VALUES(%s)", k, k)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s, team) VALUES (%s, ?) ON DUPLICATE KEY UPDATE %s",
		TableAnswerTestTwo,
		strings.Join(keys, ", "),
		strings.Join(qs, ", "),
		strings.Join(updateParts, ", "),
	)

	args := append(values, team)

	_, err := d.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// UPDATE Обновление данных
func (d *DatabaseStruct) UpdateData(table, column string, where string, arg ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var query string
	if strings.TrimSpace(where) == "" {
		query = fmt.Sprintf("UPDATE %s SET %s", table, column)
	} else {
		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, column, where)
	}

	_, err := d.db.ExecContext(ctx, query, arg...)
	if err != nil {
		return err
	}

	return nil
}

// DELETE Удаление только студентов
func (d *DatabaseStruct) DeleteStudent() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE status = ?", TableUsers)
	_, err := d.db.ExecContext(ctx, query, constants.StudentStatus)
	if err != nil {
		return err
	}

	return nil
}

// DELETE Удаление студентов по id
func (d *DatabaseStruct) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", TableUsers)
	_, err := d.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// DELETE Удаление ответа на второй вопрос
func (d *DatabaseStruct) DeleteAllAnswerTestTwo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM %s", TableAnswerTestTwo)
	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// COUNT Кол-во данных в таблице
func (d *DatabaseStruct) CountTable(table string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var count int
	err := d.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(id) FROM %s", table)).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
