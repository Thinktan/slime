package metric

import (
	log "github.com/sirupsen/logrus"

	"slime.io/slime/framework/model/trigger"
)

type WatcherProducer struct {
	name                    string
	needUpdateMetricHandler func(event trigger.WatcherEvent) QueryMap
	watcherTrigger          *trigger.WatcherTrigger
	source                  Source
	MetricChan              chan Metric
	StopChan                chan struct{}
}

func NewWatcherProducer(config WatcherProducerConfig, source Source) *WatcherProducer {
	wp := &WatcherProducer{
		name:                    config.Name,
		needUpdateMetricHandler: config.NeedUpdateMetricHandler,
		watcherTrigger:          trigger.NewWatcherTrigger(config.WatcherTriggerConfig),
		source:                  source,
		MetricChan:              config.MetricChan,
		StopChan:                make(chan struct{}),
	}

	return wp
}

func (p *WatcherProducer) HandleWatcherEvent() {
	l := log.WithField("reporter", "WatcherProducer").WithField("function", "HandleWatcherEvent")
	for {
		select {
		case <-p.StopChan:
			l.Infof("watcher producer exited")
			return
		// handle watcher event
		case event, ok := <-p.watcherTrigger.EventChan():
			if !ok {
				l.Warningf("watcher event channel closed, break process loop")
				return
			}
			l.Debugf("got watcher trigger event [%+v]", event)
			// reconciler callback
			queryMap := p.needUpdateMetricHandler(event)
			if queryMap == nil {
				l.Debugf("queryMap is nil, finish")
				continue
			}
			// queryMap: map[istio-system/istio-egressgateway:[{lazyload-accesslog-convertor }]]
			// key对wathcer监听到的sidecar资源对应的service；

			// get metric
			metric, err := p.source.QueryMetric(queryMap)
			if err != nil {
				l.Errorf("err %v", err)
				continue
			}

			log.Debugf("metric: %+v", metric)
			// map[istio-system/istio-egressgateway:[{Name:lazyload-accesslog-convertor Value:map[]}]]

			// produce metric event
			p.MetricChan <- metric
		}
	}
}

func (p *WatcherProducer) Start() {
	p.watcherTrigger.Start()
	_ = p.source.Start()
}

func (p *WatcherProducer) Stop() {
	p.watcherTrigger.Stop()
	p.StopChan <- struct{}{}
}
