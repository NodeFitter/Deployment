package measurements

import (
	"errors"
	"math"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrorIncrement           = errors.New("error while incrementing counter. More details: ")
	ErrorDecrement           = errors.New("error while decrementing counter. More details: ")
	ErrorNoInitializedConfig = errors.New("error while performing the action: meter has not been initialized. Invoke Init.")
)

type Meter struct {
	requestCounter uint
	prCounter      prometheus.Counter

	hasBeenInit bool
}

func (m *Meter) Init() error {

	// Create and register Prometheus counters
	m.prCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Help: "Total number of times the function Increment() was called",
			Name: "total_request_qt_total",
		},
	)

	// Create new Prometheus register to avoid seeing default metrics
	reg := prometheus.NewRegistry()

	reg.MustRegister(m.prCounter)

	// Initialize internal counter
	m.requestCounter = 0

	// Register metrics handler for the new register
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{},
	))

	m.hasBeenInit = true

	return nil
}

func (m *Meter) Increment() error {

	if !m.hasBeenInit {
		return ErrorNoInitializedConfig
	}

	if m.requestCounter == math.MaxUint {
		return errors.New(ErrorIncrement.Error() + "overflow maximum integer value")
	}

	m.requestCounter += 1

	m.prCounter.Add(1)

	return nil
}

func (m *Meter) Get() (uint, error) {

	if !m.hasBeenInit {
		return 0, ErrorNoInitializedConfig
	}

	return m.requestCounter, nil
}
