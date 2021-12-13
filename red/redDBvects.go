package red

import "github.com/prometheus/client_golang/prometheus"

const (
	labelRequestDB = "request"
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
	prometheus.MustRegister(requestsDbTotal, errorsDBTotal, durationDB)
	//prometheus.Register(requestsDbTotal)
	//prometheus.Register(errorsDBTotal)
	//prometheus.Register(durationDB)
}

//var (
//	durationDB                     *prometheus.SummaryVec
//	errorsDBTotal, requestsDbTotal *prometheus.CounterVec
//)
//
//func initMetricsDB() {
//	durationDB = prometheus.NewSummaryVec(
//		prometheus.SummaryOpts{
//			Name:       "duration_seconds_request_db",
//			Help:       "Summary of request duration in seconds",
//			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
//		},
//		[]string{labelRequestDB, labelMethod},
//	)
//
//	errorsDBTotal = prometheus.NewCounterVec(
//		prometheus.CounterOpts{
//			Name: "errors_request_db_total",
//			Help: "Total number of errors",
//		},
//		[]string{labelRequestDB, labelMethod},
//	)
//
//	requestsDbTotal = prometheus.NewCounterVec(
//		prometheus.CounterOpts{
//			Name: "request_db_total",
//			Help: "The total number of requests on write to the database",
//		},
//		[]string{labelRequestDB, labelMethod},
//	)
//}
