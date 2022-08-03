package tools

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/kpango/glg"
)

// GetTestDB
//
// get the test database running in DoltDB, which is the in memory db
func GetTestDB(dbName string) (db *sql.DB) {

	port := 33060

	engine := sqle.NewDefault()

	engine.AddDatabase(memory.NewDatabase(dbName))
	engine.AddDatabase(information_schema.NewInformationSchemaDatabase(engine.Catalog))

	config := server.Config{
		Protocol:       "tcp",
		Address:        fmt.Sprintf("localhost:%d", port),
		Auth:           auth.NewNativeSingle("root", "", auth.AllPermissions),
		MaxConnections: 1000,
	}
	var err error

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	go s.Start()

	time.Sleep(100 * time.Millisecond)

	connStr := "root@tcp(localhost:33060)/lemon?parseTime=true&tls=false&parseTime=true"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	//
	sqlstatement := `some create statement`

	if strings.Contains(sqlstatement, ";") == false {
		sqlstatement += ";"
	}

	stmt, err := db.Prepare(sqlstatement)
	if err != nil {
		glg.Errorf("sqlstatement:: %s", sqlstatement)
		panic(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		if len(sqlstatement) > 200 {
			sqlstatement = sqlstatement[0:200]
		}
		panic(fmt.Sprintf("%s - %s", sqlstatement, err.Error()))
	}

	return
}
