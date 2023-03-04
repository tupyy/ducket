package influxdb

import (
	"context"
	"fmt"

	influx "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/tupyy/finance/internal/entity"
)

type InfluxWriter struct {
	Url    string
	Token  string
	Org    string
	Bucket string
}

func (w *InfluxWriter) Write(transactions []*entity.Transaction) error {
	client := influx.NewClient(w.Url, w.Token)

	writeAPI := client.WriteAPIBlocking(w.Org, w.Bucket)

	for _, t := range transactions {
		point := w.createPoint(t)
		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (w *InfluxWriter) createPoint(t *entity.Transaction) *write.Point {
	p := write.NewPointWithMeasurement("transaction").
		AddField("sum", t.Sum).
		AddField("kind", t.Kind).
		AddField("content", t.RawContent).
		SetTime(t.Date)

	for k, v := range t.Labels {
		p.AddTag(k, v)
	}
	return p
}
