package src

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const errConnectPattern = "&#9940;Unable to connect to database: %v\n"

func (dbConfig *Config) DBCreateTables() string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	sqlCommand := `
		CREATE TABLE CodesRight(
			ID		integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
			Code    varchar(300) UNIQUE NOT NULL,
			Danger	varchar(50) NOT NULL,
			Sector	varchar(100) NOT NULL,
			TimeBonus	integer NOT NULL,
			TaskID	integer NOT NULL);
		CREATE INDEX ON CodesRight(Code text_pattern_ops);
		
		CREATE TABLE Tasks(
			ID		integer UNIQUE NOT NULL,
			Text	varchar(1000) NOT NULL);
		CREATE INDEX ON Tasks(ID);
			
		CREATE TABLE Teams(
			ID			integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
			Time		timestamp NOT NULL,
			Team		varchar(100) UNIQUE NOT NULL,
			Hash   	    varchar(100) UNIQUE NOT NULL,
			NickName	varchar(100),
			UserID		integer);
		CREATE INDEX ON Teams(Team text_pattern_ops);
	
		CREATE TABLE CodesUser(
			ID			integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
			Time		timestamp NOT NULL,
			NickName	varchar(100) NOT NULL,
			UserID		integer NOT NULL,
			Code		varchar(300) REFERENCES CodesRight (Code) ON DELETE CASCADE ON UPDATE CASCADE,
			Team		varchar(100) REFERENCES Teams (Team) ON DELETE SET NULL);
		CREATE INDEX ON CodesUser(NickName text_pattern_ops);
	
		CREATE TABLE Users(
			ID			integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
			NickName	varchar(100) NOT NULL,
			UserID		integer UNIQUE NOT NULL,
			Time		timestamp NOT NULL,
			Team		varchar(100) REFERENCES Teams (Team) ON DELETE SET NULL,
			Login		varchar(100),
			Password	varchar(100));
		CREATE INDEX ON Users(UserID);
	`

	_, err = db.Exec(sqlCommand)
	if err != nil {
		return fmt.Sprintf("&#9940;ERROR create tables: %s", err)
	}

	return "&#10071;Таблицы созданы заново"
}
func (dbConfig *Config) DBTruncTables(name string) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	maps := make(map[string]string)
	maps["teams"] = `truncate table Teams CASCADE;`
	maps["tasks"] = `truncate table Tasks CASCADE;`
	maps["codes"] = `truncate table CodesRight CASCADE;`

	_, err = db.Exec(maps[name])
	if err != nil {
		return fmt.Sprintf("&#9940;ERROR truncate table %s: %s", name, err)
	}

	return fmt.Sprintf("&#10071;Таблица <b>%s</b> удалена", name)
}
func (dbConfig *Config) DBSelectCodesUser(condition string) []Codes {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT ID, Time, NickName, UserID, Code, Team FROM CodesUser %s ORDER BY ID", condition)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("&#9940;Unable to SELECT CodesUser: %v\n", err)
	}
	defer rows.Close()

	var data []Codes
	for rows.Next() {
		d := Codes{}
		err := rows.Scan(&d.ID, &d.Time, &d.NickName, &d.UserID, &d.Code, &d.Team)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertCodesUsers(codes *Codes) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO CodesUser (Time, NickName, UserID, Code, Team) VALUES ((now() at time zone 'UTC-4'), $1, $2, lower($3), $4)",
		codes.NickName, codes.UserID, codes.Code, codes.Team)
	if err != nil {
		log.Println(err)
		return fmt.Sprintf("&#9940;ERROR insert CodesUser: %s", err)
	}
	return "&#9989;Код пользователя успешно добавлен"
}
func (dbConfig *Config) DBSelectCodesRight() []Codes {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, Code, Danger, Sector, TimeBonus, TaskID FROM CodesRight ORDER BY ID")
	if err != nil {
		log.Printf("&#9940;Unable to SELECT CodesRight: %v\n", err)
	}
	defer rows.Close()

	var data []Codes
	for rows.Next() {
		d := Codes{}
		err := rows.Scan(&d.ID, &d.Code, &d.Danger, &d.Sector, &d.TimeBonus, &d.TaskID)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertCodesRight(addData string) string {
	strArr := strings.Split(addData, ",")
	if len(strArr) < 4 {
		return "&#10071;Нет всех аргументов: <code>/add Code,Danger,Sector,TimeBonus,Tasks</code>"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	ArrTrimSpace(strArr)

	_, err = db.Exec("INSERT INTO CodesRight (Code, Danger, Sector, TimeBonus, TaskID) VALUES (lower($1), $2, $3, $4, $5)",
		strArr[0], strArr[1], strArr[2], strArr[3], strArr[4])
	if err != nil {
		return fmt.Sprintf("&#9940;Unable to INSERT INTO CodesRight: %v\n", err)
	}

	return "&#9989;Данные <b>добавлены</b> в БД."
}
func (dbConfig *Config) DBDeleteCodesRight(deleteStr string) string {
	if len(deleteStr) < 2 {
		return "&#10071;Нет всех аргументов: <code>/delete CodeOld</code>"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM CodesRight WHERE Code = lower($1)", strings.TrimSpace(deleteStr))
	if err != nil {
		return fmt.Sprintf("&#10071;Unable to DELETE CodesRight: %v\n", err)
	}

	return "&#9989;Данные <b>удалены</b> в БД=" + deleteStr
}
func (dbConfig *Config) DBUpdateCodesRight(updateData string) string {
	strArr := strings.Split(updateData, ",")
	if len(strArr) < 5 {
		return "&#10071;Нет всех аргументов: <code>/update CodeNew,Danger,Sector,TimeBonus,TaskID,CodeOld</code>"
	}

	ArrTrimSpace(strArr)

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE CodesRight SET Code = lower($1), Danger = $2, Sector=$3, TimeBonus=$4, TaskID=$5 WHERE Code = lower($6)",
		strArr[0], strArr[1], strArr[2], strArr[3], strArr[4], strArr[5])
	if err != nil {
		return fmt.Sprintf("&#10071;Unable to UPDATE CodesRight: %v\n", err)
	}

	return "&#9989;Данные <b>обновлены</b> в БД."
}

func (dbConfig *Config) DBSelectUsers(condition string) []Users {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT NickName, UserID, Team FROM Users %s", condition)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("&#10071;Unable to SELECT Users: %v\n", err)
	}
	defer rows.Close()

	var data []Users
	for rows.Next() {
		d := Users{}
		err := rows.Scan(&d.NickName, &d.UserID, &d.Team)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertUser(users *Users) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()
	dbConfig.DBDeleteUser(users.UserID)

	_, err = db.Exec("INSERT INTO Users (Time, NickName, UserID, Team) VALUES ((now() at time zone 'UTC-4'), $1, $2, $3)",
		users.NickName, users.UserID, users.Team)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	return fmt.Sprintf("&#9989;Игрок %s <b>добавлен</b> в команду %s.", users.NickName, users.Team)
}
func (dbConfig *Config) DBDeleteUser(UserID int) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	user := Users{}
	row := db.QueryRow("SELECT Team FROM Users WHERE UserID = $1", UserID)
	err = row.Scan(&user.Team)
	if err != nil {
		return fmt.Sprintf("&#10071;Вы не состоите в команде <b>%s</b>.", user.Team)
	}

	_, err = db.Exec("DELETE FROM Users WHERE UserID = $1", UserID)
	if err != nil {
		return fmt.Sprintf("&#9940;Unable to DELETE %d from %s: %v\n", UserID, user.Team, err)
	}

	return fmt.Sprintf("&#9745;Вы покинули команду <b>%s</b>.", user.Team)
}
func (dbConfig *Config) DBSelectTeam(condition string) []Teams {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT ID, Time, Team, Hash, NickName, UserID FROM Teams %s", condition)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("&#9940;Unable to SELECT Teams: %v\n", err)
	}
	defer rows.Close()

	var data []Teams
	for rows.Next() {
		d := Teams{}
		err := rows.Scan(&d.ID, &d.Time, &d.Team, &d.Hash, &d.NickName, &d.UserID)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertTeam(teams *Teams) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()
	// leave now team
	dbConfig.DBDeleteUser(teams.UserID)

	// create team
	_, err = db.Exec("INSERT INTO Teams (Time, Team,  Hash, NickName, UserID) VALUES ((now() at time zone 'UTC-4'), $1, $2, $3, $4)",
		teams.Team, teams.Hash, teams.NickName, teams.UserID)
	if err != nil {
		log.Println(err)
		return "&#10071; Такая команда уже есть"
	}
	// add owner in team
	user := Users{}
	user.NickName = teams.NickName
	user.Team = teams.Team
	user.UserID = teams.UserID
	dbConfig.DBInsertUser(&user)

	return fmt.Sprintf("&#9989;Команда <b>%s</b> создана, для вступления в неё введите: <code>/join %s, %s </code>", teams.Team, teams.Team, teams.Hash)
}

func (dbConfig *Config) DBSelectTask(condition string) []Tasks {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT id, Text FROM Tasks %s ORDER BY ID", condition)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("&#9940;Unable to SELECT Tasks: %v\n", err)
	}
	defer rows.Close()

	var data []Tasks
	for rows.Next() {
		d := Tasks{}
		err := rows.Scan(&d.ID, &d.Text)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertTask(tasks *Tasks) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	// create tasks
	_, err = db.Exec("INSERT INTO Tasks (ID, Text) VALUES ($1, $2)",
		tasks.ID, tasks.Text)
	if err != nil {
		log.Println(err)
		return "&#10071; Такое задание уже есть"
	}
	return "&#9989;Задание <b>добавлено</b> в БД."
}
func (dbConfig *Config) DBDeleteTask(deleteStr string) string {
	if len(deleteStr) == 0 {
		return "&#10071;Нет всех аргументов: <code>/deletetask taskID</code>"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM Tasks WHERE ID = $1", deleteStr)
	if err != nil {
		return fmt.Sprintf("&#10071;Unable to DELETE Tasks: %v\n", err)
	}

	return "&#9989;Задание <b>удалено</b> в БД=" + deleteStr
}
func (dbConfig *Config) DBUpdateTask(updateData string) string {
	strArr := strings.Split(updateData, ",")
	if len(strArr) < 2 {
		return "&#10071;Нет всех аргументов: <code>/updatetask TaskID, Task</code>"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	ArrTrimSpace(strArr)

	_, err = db.Exec("UPDATE Tasks SET Text = $1 WHERE ID = $2",
		strArr[1], strArr[0])
	if err != nil {
		return fmt.Sprintf("&#10071;Unable to UPDATE Tasks: %v\n", err)
	}

	return "&#9989;Задание <b>обновлено</b> в БД."
}
