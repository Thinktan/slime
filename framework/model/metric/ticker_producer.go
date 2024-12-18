package metric

import (
	log "github.com/sirupsen/logrus"

	"slime.io/slime/framework/model/trigger"
)

type TickerProducer struct {
	name                    string
	needUpdateMetricHandler func(event trigger.TickerEvent) QueryMap
	tickerTrigger           *trigger.TickerTrigger
	source                  Source
	MetricChan              chan Metric
	StopChan                chan struct{}
}

func NewTickerProducer(config TickerProducerConfig, source Source) *TickerProducer {
	tp := &TickerProducer{
		name:                    config.Name,
		needUpdateMetricHandler: config.NeedUpdateMetricHandler,
		tickerTrigger:           trigger.NewTickerTrigger(config.TickerTriggerConfig),
		source:                  source,
		MetricChan:              config.MetricChan,
		StopChan:                make(chan struct{}),
	}

	return tp
}

func (p *TickerProducer) HandleTickerEvent() {
	l := log.WithField("reporter", "TickerProducer").WithField("function", "HandleTriggerEvent")
	for {
		select {
		case <-p.StopChan:
			l.Infof("ticker producer exited")
			return
		// handle ticker event
		case event, ok := <-p.tickerTrigger.EventChan():
			if !ok {
				l.Warningf("ticker event channel closed, break process loop")
				return
			}

			// reconciler callback
			queryMap := p.needUpdateMetricHandler(event)
			if queryMap == nil {
				continue
			}
			l.Debugf("queryMap: %+v", queryMap)

			// get metric
			metric, err := p.source.QueryMetric(queryMap)
			if err != nil {
				l.Errorf("%v", err)
				continue
			}

			if len(metric) == 0 {
				continue
			}
			// produce metric event
			p.MetricChan <- metric
		}
	}
}

func (p *TickerProducer) Start() {
	p.tickerTrigger.Start()
	_ = p.source.Start()
}

func (p *TickerProducer) Stop() {
	p.tickerTrigger.Stop()
	p.StopChan <- struct{}{}
}
