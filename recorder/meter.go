package recorder

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xh3b4sd/tracer"
	exporter "go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MeterConfig struct {
	Env string
	Reg prometheus.Registerer // optional
	Sco string
	Ver string
}

func NewMeter(c MeterConfig) otelmetric.Meter {
	if c.Env == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Env must not be empty", c)))
	}
	if c.Sco == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Sco must not be empty", c)))
	}
	if c.Ver == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Ver must not be empty", c)))
	}

	var err error

	var opt []exporter.Option
	{
		opt = append(opt, exporter.WithoutCounterSuffixes())
	}

	if c.Reg != nil {
		opt = append(opt, exporter.WithRegisterer(c.Reg))
	}

	var exp *exporter.Exporter
	{
		exp, err = exporter.New(opt...)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exp),
	).Meter(
		fmt.Sprintf("%s.%s.splits.org", c.Sco, c.Env), // e.g. otel_scope_name="kayron.production.splits.org"
		otelmetric.WithInstrumentationVersion(c.Ver),  // e.g. otel_scope_version="v0.1.0"
	)
}
