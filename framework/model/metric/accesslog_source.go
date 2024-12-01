package metric

import (
	"net"
	"strings"
	"sync"

	service_accesslog "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"crypto/rand"
	"encoding/hex"
)

func randomHexString(length int) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type AccessLogSource struct {
	sync.Once                        // 确保某些操作只执行一次
	servePort  string                // gRPC服务监听的端口
	convertors []*AccessLogConvertor // 用于将访问日志转换为业务指标
}

func NewAccessLogSource(config AccessLogSourceConfig) *AccessLogSource {
	source := &AccessLogSource{
		servePort: config.ServePort,
	}
	for _, convertorConfig := range config.AccessLogConvertorConfigs {
		source.convertors = append(source.convertors, NewAccessLogConvertor(convertorConfig))
	}
	return source
}

// StreamAccessLogs accept access log from lazy xds egress gateway
func (s *AccessLogSource) StreamAccessLogs(logServer service_accesslog.AccessLogService_StreamAccessLogsServer) error {
	log := log.WithField("reporter", "AccessLogSource").WithField("function", "StreamAccessLogs")
	for {
		// 从gRPC流中持续接受访问日志
		message, err := logServer.Recv()
		if err != nil {
			return err
		}

		randomStr, err := randomHexString(16)
		log2 := log.WithField("traceid", randomStr)

		httpLogEntries := message.GetHttpLogs()
		log2.Infof("recv accesslog %+v", httpLogEntries)
		if httpLogEntries != nil {
			for _, convertor := range s.convertors {
				if err = convertor.Convert(httpLogEntries.LogEntry); err != nil {
					log2.Errorf("convertor [%s] converted error: %+v", convertor.Name(), err)
				} else {
					log2.Debugf("convertor [%s] converts successfully", convertor.Name())
				}
			}
		}
		log2.Infof("recv deal done")
	}
}

// Start grpc server
func (s *AccessLogSource) Start() error {
	var err error
	s.Do(func() {
		log := log.WithField("reporter", "AccessLogSource").WithField("function", "Start")
		var lis net.Listener
		lis, err = net.Listen("tcp", s.servePort) // 8082
		if err != nil {
			return
		}

		server := grpc.NewServer()
		service_accesslog.RegisterAccessLogServiceServer(server, s)

		go func() {
			log.Infof("accesslog grpc server starts on %s", s.servePort)
			if errL := server.Serve(lis); errL != nil {
				log.Errorf("accesslog grpc server error: %+v", errL)
			}
		}()
	})
	return err
}

func (s *AccessLogSource) QueryMetric(queryMap QueryMap) (Metric, error) {
	log := log.WithField("reporter", "AccessLogSource").WithField("function", "QueryMetric")

	metric := make(map[string][]Result)

	log.Debugf("queryMap: %+v", queryMap)
	// meta: istio-system/istio-egressgateway:
	// handlers: [ {lazyload-accesslog-convertor } ]

	for meta, handlers := range queryMap {
		if len(handlers) == 0 {
			continue
		}

		for _, handler := range handlers {
			for _, convertor := range s.convertors {
				if convertor.Name() != handler.Name {
					continue
				}
				result := Result{
					Name:  handler.Name,
					Value: convertor.CacheResultCopy()[meta],
				}
				metric[meta] = append(metric[meta], result)
				log.Infof("%s add metric from accesslog %+v", meta, result)
			}
		}
	}

	log.Debugf("successfully get metric from accesslog, metric: [%+v]", metric)

	return metric, nil
}

func (s *AccessLogSource) Reset(info string) error {
	parts := strings.Split(info, "/")
	ns, name := parts[0], parts[1]

	for _, convertor := range s.convertors {
		convertor.convertorLock.Lock()

		// reset ns/name
		for k := range convertor.cacheResult {
			// it will reset all svf in ns if svc is empty
			if name == "" {
				if ns == strings.Split(k, "/")[0] {
					convertor.cacheResult[k] = map[string]string{}
				}
			} else {
				if k == info {
					convertor.cacheResult[k] = map[string]string{}
				}
			}
		}

		// sync to cacheResultCopy
		newCacheResultCopy := make(map[string]map[string]string)
		for meta, value := range convertor.cacheResult {
			tmpValue := make(map[string]string)
			for k, v := range value {
				tmpValue[k] = v
			}
			newCacheResultCopy[meta] = tmpValue
		}
		convertor.cacheResultCopy = newCacheResultCopy

		convertor.convertorLock.Unlock()
	}

	return nil
}

func (s *AccessLogSource) Fullfill(cache map[string]map[string]string) error {
	for _, convertor := range s.convertors {
		convertor.convertorLock.Lock()
		for meta, value := range cache {
			convertor.cacheResult[meta] = value
			tmpValue := make(map[string]string)
			for k, v := range value {
				tmpValue[k] = v
			}

			convertor.cacheResultCopy[meta] = tmpValue
		}

		convertor.convertorLock.Unlock()
	}
	return nil
}
