package influxdb

import "time"
import (
	"fmt"
	bulkQuerygen "github.com/influxdata/influxdb-comparisons/bulk_query_gen"
	"math/rand"
)

// InfluxDashboardCpuUtilization produces Influx-specific queries for the dashboard single-host case.
type InfluxDashboardCpuUtilization struct {
	InfluxDashboard
	queryTimeRange time.Duration
}

func NewInfluxQLDashboardCpuUtilization(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(InfluxQL, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardCpuUtilization{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func NewFluxDashboardCpuUtilization(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(Flux, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardCpuUtilization{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func (d *InfluxDashboardCpuUtilization) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool

	interval := d.AllInterval.RandWindow(d.queryTimeRange)

	clusterId := fmt.Sprintf("%d", rand.Intn(15))
	var query string
	//SELECT mean("usage_user") FROM "telegraf"."default"."cpu" WHERE time > :dashboardTime: and cluster_id = :Cluster_Id: GROUP BY host, time(1m)
	query = fmt.Sprintf("SELECT last(\"max\") from (SELECT max(\"n_cpus\") FROM system WHERE cluster_id = '%s' and time >= '%s' and time < '%s' group by time(1m))", clusterId, interval.StartString(), interval.EndString())

	humanLabel := fmt.Sprintf("InfluxDB (%s) max n_cpus, rand cluster, %s by 1m", d.language.String(), d.queryTimeRange)

	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
	return q
}
