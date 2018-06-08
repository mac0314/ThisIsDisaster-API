package controllers

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"strings"

	"github.com/go-gorp/gorp"
	"github.com/revel/revel"
)

var DB *sql.DB

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func getConnectionString() string {
	host := getParamString("db.host", "")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "")
	pass := getParamString("db.password", "")
	dbname := getParamString("db.name", "")
	protocol := getParamString("db.protocol", "tcp")
	dbargs := getParamString("dbargs", " ")

	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs)
}

func defineTable(dbm *gorp.DbMap) {
	defineAchievementTable(dbm)
	defineAuthorizeTable(dbm)
	defineBodyCostumeTable(dbm)
	defineCharacterTable(dbm)
	defineErrorTable(dbm)
	defineEventTable(dbm)
	defineFeedbackTable(dbm)
	defineHeadCostumeTable(dbm)
	defineItemTable(dbm)
	defineMonsterTable(dbm)
	defineNoticeTable(dbm)
	defineStageTable(dbm)
	defineUserTable(dbm)
	defineUserSettingTable(dbm)
	defineHaveHeadCostumeTable(dbm)
	defineHaveBodyCostumeTable(dbm)
	defineDisasterTable(dbm)
}

func initDB() {
	connectionString := getConnectionString()
	if db, err := sql.Open("mysql", connectionString); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		Dbm = &gorp.DbMap{
			Db:      db,
			Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	}
	// Defines the table for use by GORP
	// This is a function we will create soon.
	defineTable(Dbm)

	if err := Dbm.CreateTablesIfNotExists(); err != nil {
		revel.ERROR.Fatal(err)
	}
}

func init() {
	revel.OnAppStart(initDB)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)
}
