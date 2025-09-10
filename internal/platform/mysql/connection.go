package person

import (
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/config"
	_ "github.com/go-sql-driver/mysql"
)

func GetDB(dbConfig config.Database) (*sql.DB, error) {
	var dsn string

	
	if dbConfig.URL != "" {
		dsn = dbConfig.URL
	} else {
		
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name, 
		)

		
		if dbConfig.SSL != "" {
			dsn += "&tls=" + dbConfig.SSL
		}
	}

	db, err := sql.Open(dbConfig.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("Error to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %w", err)
	}

	return db, nil
}
