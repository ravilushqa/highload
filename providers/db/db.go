package db

import (
	"fmt"
	"time"

	"github.com/linxGnu/mssqlx"
)

func New(masterURL string, slaveURLs []string) (*mssqlx.DBs, error) {
	dsns := make([]string, 0, len(slaveURLs)+1)
	dsns = append(dsns, masterURL)
	dsns = append(dsns, slaveURLs...)
	for i := range dsns {
		dsns[i] = dsns[i] + "?parseTime=true"
	}
	db, errs := mssqlx.ConnectMasterSlaves("mysql", dsns[:1], dsns[1:])

	for _, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("failed init db connection: %v", errs)
		}
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetHealthCheckPeriod(1000)
	errs = db.Ping()
	for _, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("database is unreachable: %v", errs)
		}
	}

	return db, nil
}
