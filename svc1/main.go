package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"lesson04/red"
	"net/http"
	"time"
)

var (
	db *sql.DB

	measurable = red.MeasurableHandler

	router = mux.NewRouter()
	web    = http.Server{
		Handler: router,
	}
)

const (
	labelRequestDB = "request"
	labelMethod    = "method"
)

var (
	durationDB = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "duration_seconds_request_db",
			Help:       "Summary of request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{labelRequestDB, labelMethod},
	)

	errorsDBTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_request_db_total",
			Help: "Total number of errors",
		},
		[]string{labelRequestDB, labelMethod},
	)

	requestsDbTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_db_total",
			Help: "The total number of requests on write to the database",
		},
		[]string{labelRequestDB, labelMethod},
	)
)

func init() {
	router.
		HandleFunc("/entities", measurable(ListEntitiesHandler)).
		Methods(http.MethodGet)
	router.
		HandleFunc("/entity", measurable(AddEntityHandler)).
		Methods(http.MethodPost)

	var err error
	db, err = sql.Open("mysql", "root:test@tcp(mysql:3306)/test")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test.entities (
		id INT PRIMARY KEY,
		data VARCHAR(32)
	);`)
	if err != nil {
		db.Close()
		return
	}

	prometheus.MustRegister(requestsDbTotal)
	prometheus.MustRegister(errorsDBTotal)
	prometheus.MustRegister(durationDB)
}

func main() {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

const sqlInsertEntity = `
	  INSERT INTO test.entities(id, data) VALUES (?, ?)
	  `

func AddEntityHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get(fmt.Sprintf("http://acl/identity?token=%s", r.FormValue("token")))

	switch {
	case err != nil:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case res.StatusCode != http.StatusOK:
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res.Body.Close()

	p := "db.Exec"
	m := "write"

	requestsDbTotal.WithLabelValues(p, m).Inc()
	t := time.Now()
	_, err = db.Exec(sqlInsertEntity, r.FormValue("id"), r.FormValue("data"))
	durationDB.WithLabelValues(p, m).Observe(time.Since(t).Seconds())
	if err != nil {
		errorsDBTotal.
			WithLabelValues(p, m).
			Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const sqlSelectEntities = `
	  SELECT id, data FROM test.entities
	  `

type ListEntityItemResponse struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func ListEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	p := "db.Query"
	m := "read"

	requestsDbTotal.WithLabelValues(p, m).Inc()

	rr, err := db.Query(sqlSelectEntities)
	durationDB.WithLabelValues(p, m).Observe(time.Since(t).Seconds())
	if err != nil {
		errorsDBTotal.
			WithLabelValues(p, m).
			Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rr.Close()

	ii := []*ListEntityItemResponse{}
	for rr.Next() {
		i := &ListEntityItemResponse{}
		err = rr.Scan(&i.Id, &i.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ii = append(ii, i)
	}
	bb, err := json.Marshal(ii)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
