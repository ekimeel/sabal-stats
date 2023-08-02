package stat

import (
	"github.com/ekimeel/sabal-pb/pb"
	"math"
	"time"
)

type stat struct {
	pointId           uint32
	lastUpdated       time.Time
	min               float64
	max               float64
	count             int
	sum               float64
	mean              float64
	m2                float64 // sum of squares of differences from the mean
	earliestTimestamp time.Time
	earliestValue     float64
	latestTimestamp   time.Time
	latestValue       float64
}

func newStat(pointId uint32) *stat {
	return &stat{
		pointId:           pointId,
		min:               math.MaxFloat64,
		max:               -math.MaxFloat64,
		lastUpdated:       time.Unix(0, 0),
		earliestTimestamp: time.Unix(math.MaxInt64, 0),
		latestTimestamp:   time.Unix(0, 0),
	}
}

func (st *stat) update(value float64, timestamp time.Time) {
	st.count++
	st.sum += value

	if value < st.min {
		st.min = value
	}

	if value > st.max {
		st.max = value
	}

	if st.earliestTimestamp.IsZero() || timestamp.Before(st.earliestTimestamp) {
		st.earliestTimestamp = timestamp
		st.earliestValue = value
	}

	if timestamp.After(st.latestTimestamp) {
		st.latestTimestamp = timestamp
		st.latestValue = value
	}

	delta := value - st.mean
	st.mean += delta / float64(st.count)
	st.m2 += delta * (value - st.mean)
}

func (st *stat) stdDev() float64 {
	if st.count < 2 {
		return 0
	}
	variance := st.m2 / float64(st.count-1)
	return math.Sqrt(variance)
}

func (st *stat) avg() float64 {
	if st.count == 0 {
		return 0
	}
	return st.sum / float64(st.count)
}

func (st *stat) calc(metrics []*pb.Metric) {
	for _, metric := range metrics {
		if metric.GetPointId() == st.pointId {
			st.update(metric.GetValue(), metric.Timestamp.AsTime())
		}
	}
}
