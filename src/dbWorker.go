package src

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const errConnectPattern = "Unable to connect to database: %v\n"

func (dbConfig *Config) DBResetTable(table Tables) string {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	// delete table
	_, err = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table.name))
	if err != nil {
		return fmt.Sprintf("ERROR delete table %s: %s", table.name, err)
	}

	// create table
	_, err = db.Exec(table.ddl)
	if err != nil {
		return fmt.Sprintf("ERROR create table %s: %s", table.name, err)
	}

	return fmt.Sprintf("\n\n&#9940;БД %s удалена и создана заново!", table.name)
}
func (dbConfig *Config) DBSelectCodesUser(condition string) []Codes {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT ID, Time, NickName, UserID, Code, Danger, Sector, Team FROM CodesUser %s", condition)

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Unable to SELECT CodesUser: %v\n", err)
	}
	defer rows.Close()

	var data []Codes
	for rows.Next() {
		d := Codes{}
		err := rows.Scan(&d.ID, &d.Time, &d.NickName, &d.UserID, &d.Code, &d.Danger, &d.Sector, &d.Team)
		if err != nil {
			log.Println(err)
			continue
		}
		data = append(data, d)
	}

	return data
}
func (dbConfig *Config) DBInsertCodesUsers(codes *Codes) {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO CodesUser (Time, NickName, UserID, Code, Danger, Sector, Team) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		codes.Time, codes.NickName, codes.UserID, codes.Code, codes.Danger, codes.Sector, codes.Team)
	if err != nil {
		log.Println(err)
	}
}
func (dbConfig *Config) DBSelectCodesRight() []Codes {
	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		log.Printf(errConnectPattern, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, Code, Danger, Sector FROM CodesRight")
	if err != nil {
		log.Printf("Unable to SELECT CodesRight: %v\n", err)
	}
	defer rows.Close()

	var data []Codes
	for rows.Next() {
		d := Codes{}
		err := rows.Scan(&d.ID, &d.Code, &d.Danger, &d.Sector)
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
	if len(strArr) < 3 {
		return "&#10071;Нет всех аргументов: /add Code,Danger,Sector"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO CodesRight (Code, Danger, Sector) VALUES ($1, $2, $3)",
		strings.TrimSpace(strArr[0]), strings.TrimSpace(strArr[1]), strings.TrimSpace(strArr[2]))
	if err != nil {
		return fmt.Sprintf("Unable to INSERT INTO CodesRight: %v\n", err)
	}

	return "&#10004;Данные <b>добавлены</b> в БД."
}
func (dbConfig *Config) DBDeleteCodesRight(deleteStr string) string {
	if len(deleteStr) < 2 {
		return "&#10071;Нет всех аргументов: /delete CodeOld"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM CodesRight WHERE Code = $1", deleteStr)
	if err != nil {
		return fmt.Sprintf("Unable to DELETE CodesRight: %v\n", err)
	}

	return "&#8252;Данные <b>удалены</b> в БД=" + deleteStr
}
func (dbConfig *Config) DBUpdateCodesRight(updateData string) string {
	strArr := strings.Split(updateData, ",")
	if len(strArr) < 4 {
		return "&#10071;Нет всех аргументов: /update CodeNew,Danger,Sector,CodeOld"
	}

	db, err := sql.Open(dbConfig.DriverNameDB, dbConfig.DBURL)
	if err != nil {
		return fmt.Sprintf(errConnectPattern, err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE CodesRight SET Code = $1, Danger = $2, Sector=$3 WHERE Code = $4",
		strings.TrimSpace(strArr[0]), strArr[1], strArr[2], strings.TrimSpace(strArr[3]))
	if err != nil {
		return fmt.Sprintf("Unable to UPDATE CodesRight: %v\n", err)
	}

	return "&#10071;Данные <b>обновлены</b> в БД."
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
		log.Printf("Unable to SELECT Users: %v\n", err)
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

	_, err = db.Exec("INSERT INTO Users (NickName, UserID, Time, Team) VALUES ($1, $2, $3, $4)",
		users.NickName, users.UserID, users.Time, users.Team)
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
		return fmt.Sprintf("&#8252;Вы не состоите в команде <b>%s</b>.", user.Team)
	}

	_, err = db.Exec("DELETE FROM Users WHERE UserID = $1", UserID)
	if err != nil {
		log.Printf("Unable to DELETE %d from %s: %v\n", UserID, user.Team, err)
		return fmt.Sprintf("Невозможно удалить <b>%s</b> из команды <b>%s</b>\n", user.NickName, user.Team)
	}

	return fmt.Sprintf("&#8252;Вы покинули команду <b>%s</b>.", user.Team)
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
		log.Printf("Unable to SELECT Teams: %v\n", err)
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
	_, err = db.Exec("INSERT INTO Teams (Time, Team,  Hash, NickName, UserID) VALUES ($1, $2, $3, $4, $5)",
		teams.Time, teams.Team, teams.Hash, teams.NickName, teams.UserID)
	if err != nil {
		log.Println(err)
		return "&#10071; Такая команда уже есть"
	}
	// add owner in team
	user := Users{}
	user.NickName = teams.NickName
	user.Time = teams.Time
	user.Team = teams.Team
	user.UserID = teams.UserID
	dbConfig.DBInsertUser(&user)

	return fmt.Sprintf("&#9989;Команда <b>%s</b> создана, для вступления в неё введите: <code>/join %s, %s </code>", teams.Team, teams.Team, teams.Hash)
}
