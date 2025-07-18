package registry

import (
	"slices"

	"github.com/0xSplits/otelgo/recorder"
	"github.com/xh3b4sd/tracer"
	"go.opentelemetry.io/otel/attribute"
)

// record validates and executes the data collection for the given metric, if
// that metric is registered in the provided whitelist of recorder interfaces.
// The correct metric names must be registered for the respective counters,
// gauges and histograms. The provided labels must be registered too, by key and
// value respectively.
func (r *Registry) record(wht map[string]recorder.Interface, nam string, val float64, lab map[string]string) error {
	var rec recorder.Interface
	var exi bool

	{
		rec, exi = wht[nam]
		if !exi {
			return tracer.Mask(metricNameWhitelistError, tracer.Context{Key: "metric name", Value: nam})
		}
	}

	for k, v := range lab {
		var lis []string
		{
			lis, exi = rec.Labels()[k]
			if !exi {
				return tracer.Mask(labelKeyWhitelistError, tracer.Context{Key: "label key", Value: k})
			}
		}

		{
			exi = slices.Contains(lis, v)
			if !exi {
				return tracer.Mask(labelValueWhitelistError, tracer.Context{Key: "label value", Value: v})
			}
		}
	}

	// Create the set of labels according to the registry's environment and the
	// provided key-value pairs. Not all use cases are environment specific, so if
	// this registry instance does not have any enviroment configured, then we do
	// not add the "env" label.

	var att []attribute.KeyValue
	if r.env != "" {
		att = append(att, attribute.String("env", r.env))
	}

	for k, v := range lab {
		att = append(att, attribute.String(k, v))
	}

	// Record the given data point, with or without labels.

	{
		rec.Record(val, att...)
	}

	return nil
}
