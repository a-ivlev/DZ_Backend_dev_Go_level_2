package red

import (
	"database/sql"
	"time"
)

func ExecDBmetrics(db *sql.DB, sqlInsertEntity string, arg ...string) (*sql.Result, error) {
	p := "db.Exec"
	m := "write"

	requestsDbTotal.WithLabelValues(p, m).Inc()
	startTime := time.Now()
	result, err := db.Exec(sqlInsertEntity, arg)
	if err != nil {
		errorsDBTotal.
			WithLabelValues(p, m).
			Inc()
		return &result, err
	}
	durationDB.WithLabelValues(p, m).Observe(time.Since(startTime).Seconds())

	return &result, err
}

func QueryDBmetrics(db *sql.DB, sqlSelectEntities string) (*sql.Rows, error) {
	p := "db.Query"
	m := "read"

	requestsDbTotal.WithLabelValues(p, m).Inc()
	start := time.Now()
	rr, err := db.Query(sqlSelectEntities)
	if err != nil {
		errorsDBTotal.
			WithLabelValues(p, m).
			Inc()
		return rr, err
	}

	durationDB.WithLabelValues(p, m).Observe(time.Since(start).Seconds())

	return rr, err
}
