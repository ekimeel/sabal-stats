package stat

import (
	"context"
	"github.com/ekimeel/sabal-pb/pb"
	"github.com/ekimeel/sabal-plugin/pkg/metric_utils"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (
	PluginName    = "stat"
	PluginVersion = "v1.0"
)

var (
	singletonService *Service
	onceService      sync.Once
	Logger           *log.Logger
)

type Service struct {
	dao *dao
}

func GetService() *Service {

	onceService.Do(func() {
		singletonService = &Service{}
		singletonService.dao = getDao()
	})

	return singletonService
}

func (s *Service) Run(ctx context.Context, metrics []*pb.Metric) {
	unitOfWork := metric_utils.GroupMetricsByPointId(metrics)
	var wg sync.WaitGroup

	for pointId, items := range unitOfWork {
		wg.Add(1)
		go func(pointId uint32, items []*pb.Metric) {
			defer wg.Done()
			s.compute(pointId, items)
		}(pointId, items)
	}

	wg.Wait()
}

func (s *Service) compute(pointId uint32, metrics []*pb.Metric) {
	Logger.WithField("plugin", PluginName).Tracef("computing point: %d", pointId)

	st, err := s.dao.selectByPointId(pointId)

	if err != nil {
		Logger.WithField("plugin", PluginName).Errorf("dao error: %s", err)
		return
	}

	if st == nil {
		Logger.WithField("plugin", PluginName).Tracef("no stat found, creating new stat for point: %d", pointId)
		st = newStat(pointId)
		st.calc(metrics)
		Logger.WithField("plugin", PluginName).Tracef("attempting to insert stat for point: %d", pointId)
		_, err = s.dao.insert(st)
	} else {
		st.calc(metrics)
		Logger.WithField("plugin", PluginName).Tracef("updating stat for point: %d", pointId)
		_, err = s.dao.update(st)
	}

}
