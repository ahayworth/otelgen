package metrics

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	apiMetric "go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
)

// Histogram demonstrates how to record a distribution of individual values
func SimulateHistogram(mp *sdkMetric.MeterProvider, conf *Config, logger *zap.Logger) {
	c := *conf
	err := run(conf, logger, histogram(mp, c, logger))
	if err != nil {
		logger.Error("failed to run histogram", zap.Error(err))
	}
}

// histogram generates a histogram metric
func histogram(mp *sdkMetric.MeterProvider, c Config, logger *zap.Logger) WorkerFunc {
	return func(ctx context.Context) {
		name := fmt.Sprintf("%v.metrics.histogram", c.ServiceName)
		durRecorder, _ := mp.Meter(c.ServiceName).Int64Histogram(
			name,
			apiMetric.WithUnit("microseconds"),
			apiMetric.WithDescription("Histogram demonstrates how to record a distribution of individual values"),
		)

		if c.TotalDuration > 0 {
			logger.Info("generation duration", zap.Float64("seconds", c.TotalDuration.Seconds()))

		loop:
			for timeout := time.After(c.TotalDuration); ; {
				select {
				case <-timeout:
					break loop
				default:
				}
				logger.Info("generating", zap.String("name", name))
				dur := time.Duration(rand.NormFloat64()*5000000) * time.Microsecond
				durRecorder.Record(ctx, dur.Microseconds())
				time.Sleep(time.Duration(c.Rate) * time.Second)
			}
		} else {
			for {
				logger.Info("generating", zap.String("name", name))
				dur := time.Duration(rand.NormFloat64()*5000000) * time.Microsecond
				durRecorder.Record(ctx, dur.Microseconds())
				time.Sleep(time.Duration(c.Rate) * time.Second)
			}
		}
	}
}
