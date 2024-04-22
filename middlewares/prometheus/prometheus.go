package prometheus

import (
	"earnth"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (m *MiddlewareBuilder) Build() earnth.MiddlewareFunc {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        m.Name,
		Subsystem:   m.Subsystem,
		ConstLabels: m.ConstLabels,
		Help:        m.Help,
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(summaryVec)
	return func(next earnth.HandleFunc) earnth.HandleFunc {
		return func(ctx *earnth.Context) {
			startTime := time.Now()
			next(ctx)
			endTime := time.Now()
			go report(endTime.Sub(startTime), ctx, summaryVec)
		}
	}
}

func report(dur time.Duration, ctx *earnth.Context, vec prometheus.ObserverVec) {
	status := ctx.RespStatusCode
	route := "unknown"
	//if ctx.MatchedRoute != "" {
	//	route = ctx.MatchedRoute
	//}
	ms := dur / time.Millisecond
	vec.WithLabelValues(route, ctx.Req.Method, strconv.Itoa(status)).Observe(float64(ms))
}
