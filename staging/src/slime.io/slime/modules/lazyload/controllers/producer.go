package controllers

import (
	"context"
	stderrors "errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	envoy_config_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	data_accesslog "github.com/envoyproxy/go-control-plane/envoy/data/accesslog/v3"
	prometheusApi "github.com/prometheus/client_golang/api"
	prometheusV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	"slime.io/slime/framework/apis/config/v1alpha1"
	"slime.io/slime/framework/bootstrap"
	"slime.io/slime/framework/model/metric"
	"slime.io/slime/framework/model/trigger"
	"slime.io/slime/modules/lazyload/api/config"
	lazyloadapiv1alpha1 "slime.io/slime/modules/lazyload/api/v1alpha1"
)

const (
	AccessLogConvertorName     = "lazyload-accesslog-convertor"
	MetricSourceTypePrometheus = "prometheus"
	MetricSourceTypeAccesslog  = "accesslog"
	SvfResource                = "servicefences"
)

// call back function for watcher producer
func (r *ServicefenceReconciler) handleWatcherEvent(event trigger.WatcherEvent) metric.QueryMap {
	log := log.WithField("func", "handleWatcherEvent")
	log.Debugf("watcher event: %+v", event)
	// check event
	gvks := []schema.GroupVersionKind{
		{Group: "networking.istio.io", Version: "v1alpha3", Kind: "Sidecar"},
	}
	invalidEvent := false
	for _, gvk := range gvks {
		if event.GVK == gvk && r.getInterestMeta()[event.NN.String()] {
			invalidEvent = true
		}
	}
	if !invalidEvent {
		log.Debugf("invalid Event, InterestMest: %+v", r.getInterestMeta())
		return nil
	}
	// generate query map for producer
	qm := make(map[string][]metric.Handler)
	var hs []metric.Handler

	// check metric source type
	switch r.cfg.MetricSourceType {
	case MetricSourceTypePrometheus:
		for pName, pHandler := range r.env.Config.Metric.Prometheus.Handlers {
			hs = append(hs, generateHandler(event.NN.Name, event.NN.Namespace, pName, pHandler))
		}
	case MetricSourceTypeAccesslog:
		hs = []metric.Handler{
			{
				Name:  AccessLogConvertorName,
				Query: "",
			},
		}
	}

	qm[event.NN.String()] = hs
	log.Debugf("qm: +v", qm)

	return qm
}

// call back function for ticker producer
func (r *ServicefenceReconciler) handleTickerEvent(_ trigger.TickerEvent) metric.QueryMap {
	// no need to check time duration

	// generate query map for producer
	// check metric source type
	qm := make(map[string][]metric.Handler)
	log := log.WithField("func", "handleTickerEvent")

	switch r.cfg.MetricSourceType {
	case MetricSourceTypePrometheus:
		for meta := range r.getInterestMeta() {
			namespace, name := strings.Split(meta, "/")[0], strings.Split(meta, "/")[1]
			var hs []metric.Handler
			for pName, pHandler := range r.env.Config.Metric.Prometheus.Handlers {
				hs = append(hs, generateHandler(name, namespace, pName, pHandler))
			}
			qm[meta] = hs
		}
	case MetricSourceTypeAccesslog:
		for meta := range r.getInterestMeta() {
			qm[meta] = []metric.Handler{
				{
					Name:  AccessLogConvertorName,
					Query: "",
				},
			}
		}
	}

	log.Infof("qm: %+v", qm)
	// map[
	//	istio-system/istio-egressgateway:[{Name:lazyload-accesslog-convertor Query:}]
	//	istio-system/istio-ingressgateway:[{Name:lazyload-accesslog-convertor Query:}]
	//	istio-system/istiod:[{Name:lazyload-accesslog-convertor Query:}]
	//	kube-system/kube-dns:[{Name:lazyload-accesslog-convertor Query:}]
	//	mesh-operator/lazyload:[{Name:lazyload-accesslog-convertor Query:}]]

	return qm
}

func generateHandler(name, namespace, pName string, pHandler *v1alpha1.Prometheus_Source_Handler) metric.Handler {
	query := strings.ReplaceAll(pHandler.Query, "$namespace", namespace)
	query = strings.ReplaceAll(query, "$source_app", name)
	return metric.Handler{Name: pName, Query: query}
}

func NewProducerConfig(env bootstrap.Environment, cfg config.Fence) (*metric.ProducerConfig, error) {
	// init metric source
	var enablePrometheusSource bool
	var prometheusSourceConfig metric.PrometheusSourceConfig
	var accessLogSourceConfig metric.AccessLogSourceConfig
	var err error

	log.Infof("tklog env:%+v", env)
	// env:{
	//	Config:global:
	//	{service:"app" istioNamespace:"istio-system" slimeNamespace:"mesh-operator"
	//	log:{logLevel:"info" klogLevel:5}
	//	misc:{key:"aux-addr" value:":8081"}
	//	misc:{key:"enableLeaderElection" value:"off"}
	//	misc:{key:"logSourcePort" value:":8082"}
	//	misc:{key:"metrics-addr" value:":8080"}
	//	misc:{key:"pathRedirect" value:""}
	//	misc:{key:"seLabelSelectorKeys" value:"app"}
	//	misc:{key:"xdsSourceEnableIncPush" value:"true"}
	//	}
	//	name:"lazyload" enable:true general:{}
	//	kind:"lazyload" K8SClient:0xc0007896c0 DynamicClient:0xc00043a4a0
	//	HttpPathHandler:{Prefix:lazyload PathHandler:0xc000953f20}
	//	ReadyManager:0x2224c20 Stop:0xc0001a1880 ConfigController:<nil>
	//	IstioConfigController:<nil>} module=lazyload pkg=controllers

	log.Infof("tklog cfg:%+v", cfg)
	// {state:{NoUnkeyedLiterals:{} DoNotCompare:[] DoNotCopy:[] atomicMessageInfo:0xc00031c808}
	//sizeCache:0 unknownFields:[] WormholePort:[9080] AutoFence:true Namespace:[] Dispatches:[] DomainAliases:[]
	//DefaultFence:true AutoPort:true ClusterGsNamespace: FenceLabelKeyAlias: EnableShortDomain:false
	//DisableIpv4Passthrough:false PassthroughByDefault:false SupportH2:false AddEnvHeaderViaLua:false
	//GlobalSidecarMode:cluster Render: MetricSourceType:accesslog CleanupWormholePort:false ManagementSelectors:[]
	//NamespaceList:<nil> ProxyVersion: StableHost:[]} module=lazyload pkg=controllers

	switch cfg.MetricSourceType {
	case MetricSourceTypePrometheus:
		enablePrometheusSource = true
		prometheusSourceConfig, err = newPrometheusSourceConfig(env)
		if err != nil {
			return nil, err
		}
	case MetricSourceTypeAccesslog:
		enablePrometheusSource = false
		// init log source port
		port := env.Config.Global.Misc["logSourcePort"]

		// init accessLog source config
		accessLogSourceConfig = metric.AccessLogSourceConfig{
			ServePort: port,
			AccessLogConvertorConfigs: []metric.AccessLogConvertorConfig{
				{
					Name:    AccessLogConvertorName,
					Handler: nil,
				},
			},
		}
	default:
		return nil, stderrors.New("wrong metricSourceType")
	}

	log.Infof("tklog prometheusSlurceConfig: %+v", prometheusSourceConfig)
	log.Infof("tklog accessLogSourceConfig: %+v", accessLogSourceConfig)
	// {ServePort::8082 AccessLogConvertorConfigs:[{Name:lazyload-accesslog-convertor Handler:<nil>}]} module=lazyload pkg=controller

	// init whole producer config
	pc := &metric.ProducerConfig{
		EnablePrometheusSource: enablePrometheusSource,
		PrometheusSourceConfig: prometheusSourceConfig,
		AccessLogSourceConfig:  accessLogSourceConfig,
		EnableWatcherProducer:  true,
		WatcherProducerConfig: metric.WatcherProducerConfig{
			Name:       "lazyload-watcher",
			MetricChan: make(chan metric.Metric),
			WatcherTriggerConfig: trigger.WatcherTriggerConfig{
				Kinds: []schema.GroupVersionKind{
					{
						Group:   "networking.istio.io",
						Version: "v1alpha3",
						Kind:    "Sidecar",
					},
				},
				EventChan:     make(chan trigger.WatcherEvent),
				DynamicClient: env.DynamicClient,
			},
		},
		EnableTickerProducer: true,
		TickerProducerConfig: metric.TickerProducerConfig{
			Name:       "lazyload-ticker",
			MetricChan: make(chan metric.Metric),
			TickerTriggerConfig: trigger.TickerTriggerConfig{
				Durations: []time.Duration{
					10 * time.Second,
				},
				EventChan: make(chan trigger.TickerEvent),
			},
		},
		StopChan: env.Stop,
	}

	log.Infof("tklog whole producer config: %+v", pc)
	// {
	//	EnablePrometheusSource:false PrometheusSourceConfig:{Api:<nil> Convertor:<nil>}
	// 	AccessLogSourceConfig:{ServePort::8082
	//	AccessLogConvertorConfigs:[{Name:lazyload-accesslog-convertor Handler:<nil>}]}
	//	EnableMockSource:false EnableWatcherProducer:true
	//	WatcherProducerConfig:{
	//		Name:lazyload-watcher NeedUpdateMetricHandler:<nil>
	//		MetricChan:0xc0001a18f0 WatcherTriggerConfig:{Kinds:[networking.istio.io/v1alpha3, Kind=Sidecar]
	//		DynamicClient:0xc00043a4a0 EventChan:0xc0001a1960}}
	//	EnableTickerProducer:true
	//	TickerProducerConfig:{
	//		Name:lazyload-ticker NeedUpdateMetricHandler:<nil>
	//		MetricChan:0xc0001a19d0 TickerTriggerConfig:{Durations:[10s] EventChan:0xc0001a1a40}}
	//		StopChan:0xc0001a1880}
	//	module=lazyload pkg=controllers

	return pc, nil
}

// nolint: lll
func (r *ServicefenceReconciler) LogHandler(logEntry []*data_accesslog.HTTPAccessLogEntry) (map[string]map[string]string, error) {
	return accessLogHandler(logEntry, r.ipToSvcCache, r.svcToIpsCache, r.ipTofence, r.fenceToIp, r.cfg.EnableShortDomain)
}

func newPrometheusSourceConfig(env bootstrap.Environment) (metric.PrometheusSourceConfig, error) {
	ps := env.Config.Metric.Prometheus
	if ps == nil {
		return metric.PrometheusSourceConfig{}, stderrors.New("failure create prometheus client, empty prometheus config")
	}
	promClient, err := prometheusApi.NewClient(prometheusApi.Config{
		Address:      ps.Address,
		RoundTripper: nil,
	})
	if err != nil {
		return metric.PrometheusSourceConfig{}, err
	}

	return metric.PrometheusSourceConfig{
		Api: prometheusV1.NewAPI(promClient),
	}, nil
}

func NewCache(env bootstrap.Environment) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	svfGvr := schema.GroupVersionResource{
		Group:    lazyloadapiv1alpha1.GroupVersion.Group,
		Version:  lazyloadapiv1alpha1.GroupVersion.Version,
		Resource: SvfResource,
	}

	dc := env.DynamicClient
	svfList, err := dc.Resource(svfGvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list servicefence error: %v", err)
	}
	ServicefenceLoads.Record(float64(len(svfList.Items)))
	log.Infof("tklog NewCache svfList: %+v", svfList)
	// {Object:
	//	map[apiVersion:microservice.slime.io/v1alpha1 kind:ServiceFenceList
	//	metadata:map[continue: resourceVersion:29444]]
	//	Items:[
	//	{Object:
	//	map[apiVersion:microservice.slime.io/v1alpha1
	//	kind:ServiceFence
	//	metadata:map[creationTimestamp:2024-11-27T02:00:52Z generation:1
	//		labels:map[app.kubernetes.io/created-by:fence-controller]
	//		managedFields:[map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1
	//		fieldsV1:map[f:metadata:map[f:labels:map[.:map[] f:app.kubernetes.io/created-by:map[]]]
	//		f:spec:map[.:map[] f:enable:map[] f:workloadSelector:map[.:map[]
	//		f:fromService:map[]]]] manager:manager
	//		operation:Update time:2024-11-27T02:00:52Z]
	//	map[apiVersion:microservice.slime.io/v1alpha1
	//	fieldsType:FieldsV1 fieldsV1:map[f:status:map[]] manager:manager operation:Update
	//	subresource:status time:2024-11-27T02:00:52Z]]
	//	name:istio-egressgateway namespace:istio-system resourceVersion:1476 uid:8e6d862e-05eb-43cb-8382-71cba70a6796] spec:map[enable:true workloadSelector:map[fromService:true]] status:map[]]}
	//	{Object:map[apiVersion:microservice.slime.io/v1alpha1 kind:ServiceFence metadata:map[creationTimestamp:2024-11-27T02:00:52Z generation:1 labels:map[app.kubernetes.io/created-by:fence-controller] managedFields:[map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:metadata:map[f:labels:map[.:map[] f:app.kubernetes.io/created-by:map[]]] f:spec:map[.:map[] f:enable:map[] f:workloadSelector:map[.:map[] f:fromService:map[]]]] manager:manager operation:Update time:2024-11-27T02:00:52Z] map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:status:map[]] manager:manager operation:Update subresource:status time:2024-11-27T02:00:52Z]]
	//	name:istio-ingressgateway namespace:istio-system resourceVersion:1478 uid:68dd279d-8b68-40f9-9211-8afc95835166] spec:map[enable:true workloadSelector:map[fromService:true]] status:map[]]}
	//	{Object:map[apiVersion:microservice.slime.io/v1alpha1 kind:ServiceFence metadata:map[creationTimestamp:2024-11-27T02:00:52Z generation:1 labels:map[app.kubernetes.io/created-by:fence-controller] managedFields:[map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:metadata:map[f:labels:map[.:map[] f:app.kubernetes.io/created-by:map[]]] f:spec:map[.:map[] f:enable:map[] f:workloadSelector:map[.:map[] f:fromService:map[]]]] manager:manager operation:Update time:2024-11-27T02:00:52Z] map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:status:map[]] manager:manager operation:Update subresource:status time:2024-11-27T02:00:52Z]]
	//	name:istiod namespace:istio-system resourceVersion:1480 uid:cd9a3925-a4e5-4e9f-81ae-33b5b27dc749] spec:map[enable:true workloadSelector:map[fromService:true]] status:map[]]}
	//	{Object:map[apiVersion:microservice.slime.io/v1alpha1 kind:ServiceFence metadata:map[creationTimestamp:2024-11-27T02:00:52Z generation:1 labels:map[app.kubernetes.io/created-by:fence-controller] managedFields:[map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:metadata:map[f:labels:map[.:map[] f:app.kubernetes.io/created-by:map[]]] f:spec:map[.:map[] f:enable:map[] f:workloadSelector:map[.:map[] f:fromService:map[]]]] manager:manager operation:Update time:2024-11-27T02:00:52Z] map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:status:map[]] manager:manager operation:Update subresource:status time:2024-11-27T02:00:52Z]]
	//	name:kube-dns namespace:kube-system resourceVersion:1463 uid:54f7b0d5-810f-4c46-94da-7246043d4d1d] spec:map[enable:true workloadSelector:map[fromService:true]] status:map[]]}
	//	{Object:map[apiVersion:microservice.slime.io/v1alpha1 kind:ServiceFence metadata:map[creationTimestamp:2024-11-27T02:00:52Z generation:1 labels:map[app.kubernetes.io/created-by:fence-controller] managedFields:[map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:metadata:map[f:labels:map[.:map[] f:app.kubernetes.io/created-by:map[]]] f:spec:map[.:map[] f:enable:map[] f:workloadSelector:map[.:map[] f:fromService:map[]]]] manager:manager operation:Update time:2024-11-27T02:00:52Z] map[apiVersion:microservice.slime.io/v1alpha1 fieldsType:FieldsV1 fieldsV1:map[f:status:map[]] manager:manager operation:Update subresource:status time:2024-11-27T02:00:52Z]]
	//	name:lazyload namespace:mesh-operator resourceVersion:1474 uid:b487dc51-0cd9-45d8-bcbf-44640dab6791] spec:map[enable:true workloadSelector:map[fromService:true]] status:map[]]}]}

	for _, svf := range svfList.Items {
		meta := svf.GetNamespace() + "/" + svf.GetName()
		value := make(map[string]string)
		log.Infof("tklog NewCache meta:%+v, value:%+v", meta, value)
		ms, existed, err := unstructured.NestedMap(svf.Object, "status", "metricStatus")
		if err != nil {
			log.Errorf("tklog got servicefence %s status.metricStatus error: %v", meta, err)
			continue
		}
		if existed {
			for k, v := range ms {
				ok := false
				if value[k], ok = v.(string); !ok {
					log.Errorf("servicefence %s metricStatus value is not string, value: %+v", meta, v)
					continue
				}
			}
		}
		result[meta] = value
	}

	return result, nil
}

func accessLogHandler(logEntry []*data_accesslog.HTTPAccessLogEntry, ipToSvcCache *IpToSvcCache,
	svcToIpsCache *SvcToIpsCache, ipTofenceCache *IpTofence, _ *FenceToIp, enableShortDomain bool,
) (map[string]map[string]string, error) {
	log = log.WithField("reporter", "AccessLogConvertor").WithField("function", "accessLogHandler")
	// map sourceSvc to destinationSvc
	result := make(map[string]map[string]string)
	tmpResult := make(map[string]map[string]int)

	for _, entry := range logEntry {
		// fetch source ip
		log.Debugf("entry: %+v", entry)
		sourceIp, err := fetchSourceIp(entry)
		log.Debugf("sourceIp: %+v, err: %+v", sourceIp, err)
		if err != nil {
			return nil, err
		}
		if sourceIp == "" {
			continue
		}

		// fetch all the services which is the source ip belongs to
		sourceSvcs, err := spliceSourceSvc(sourceIp, ipToSvcCache)
		log.Debugf("sourceSvcs: %+v, err: %+v", sourceSvcs, err)
		if err != nil {
			return nil, err
		}

		// fetch workload fence which is the source ip associated with
		fenceNN, err := spliceSourcefence(sourceIp, ipTofenceCache)
		log.Debugf("fenceNN: %+v, err: %+v", fenceNN, err)
		if err != nil {
			fenceNN = nil
			return nil, err
		}

		// the source ip is not a pod ip or the pod is not selected by any service and the workloadfence is not enabled, skip
		if len(sourceSvcs) == 0 && fenceNN == nil {
			continue
		}

		// fetch all destination services like:
		// []string{`{destination_service="foo.default.svc.cluster.local"`}
		destinationSvcs := spliceDestinationSvc(entry, sourceSvcs, svcToIpsCache, fenceNN, enableShortDomain)
		log.Debugf("destinationSvcs: %+v", destinationSvcs)
		if len(destinationSvcs) == 0 {
			continue
		}

		// record source service -> dest
		for _, sourceSvc := range sourceSvcs {
			if _, ok := tmpResult[sourceSvc]; !ok {
				tmpResult[sourceSvc] = make(map[string]int)
			}
			for _, dest := range destinationSvcs {
				tmpResult[sourceSvc][dest] += 1
			}
		}

		// record workload servicefence -> dest
		if fenceNN != nil {
			nn := fenceNN.String()
			if _, ok := tmpResult[nn]; !ok {
				tmpResult[nn] = make(map[string]int)
			}
			for _, dest := range destinationSvcs {
				tmpResult[nn][dest] += 1
			}
		}
	}

	for sourceSvc, dstSvcMappings := range tmpResult {
		result[sourceSvc] = make(map[string]string)
		for dstSvc, count := range dstSvcMappings {
			result[sourceSvc][dstSvc] = strconv.Itoa(count)
		}
	}

	return result, nil
}

func fetchSourceIp(entry *data_accesslog.HTTPAccessLogEntry) (string, error) {
	log := log.WithField("reporter", "accesslog convertor").WithField("function", "fetchSourceIp")
	if entry.CommonProperties.DownstreamDirectRemoteAddress == nil {
		log.Debugf("DownstreamDirectRemoteAddress is nil, skip")
		return "", nil
	}
	addr := entry.CommonProperties.DownstreamDirectRemoteAddress.Address
	downstreamSock, ok := addr.(*envoy_config_core.Address_SocketAddress)
	if !ok {
		return "", stderrors.New("wrong type of DownstreamDirectRemoteAddress")
	}
	if downstreamSock == nil || downstreamSock.SocketAddress == nil {
		return "", stderrors.New("downstream socket address is nil")
	}
	log.Infof("SourceEp is: %s", downstreamSock.SocketAddress.Address)
	return downstreamSock.SocketAddress.Address, nil
}

func spliceSourceSvc(sourceIp string, ipToSvcCache *IpToSvcCache) ([]string, error) {
	log = log.WithField("function", "spliceSourceSvc")
	ipToSvcCache.RLock()
	defer ipToSvcCache.RUnlock()

	log.Debugf("sourceIp: %+v, iptoSvcCache: %+v", sourceIp, ipToSvcCache)

	if svc, ok := ipToSvcCache.Data[sourceIp]; ok {
		keys := make([]string, 0, len(svc))
		for key := range svc {
			keys = append(keys, key)
		}
		return keys, nil
	}

	log.Infof("svc not found base on sourceIp %s", sourceIp)
	return []string{}, nil
}

func spliceSourcefence(sourceIp string, ipTofence *IpTofence) (*types.NamespacedName, error) {
	ipTofence.RLock()
	defer ipTofence.RUnlock()
	log.WithField("function", "spliceSourcefence")

	log.Debugf("sourceIp: %+v, ipTofence: %+v", sourceIp, ipTofence)

	if nn, ok := ipTofence.Data[sourceIp]; ok {
		return &nn, nil
	}
	log.Debugf("fence not found base on sourceIp %s", sourceIp)
	return nil, nil
}

// spliceDestinationSvc splices destination service from entry with below rules
// - only handle inbound access log
// - only handle request whih non-IP authority
// - try resolve shortname to FQDN:
//   - if the dest only contains one label, try add ns suffix from sourceSvc's namespace or workload fence's namespace.
//     because this domain can only be resolved when the source client and destination service is in the same namespace.
//   - if the dest contains two labels,like `label[0].lable[1]`, and `label[1]/label[0]` match a service, we append k8s
//     k8s local domain as suffix. otherwise, we treat it as a fqdn.
//   - if the dest contains three labels, like `label[0].lable[1].label[2]`, and `label[1]/label[0]` match a service,
//     and the `label[2]` equals `svc`, we append k8s local domain as suffix. otherwise, we treat it as a fqdn.
//   - if the dest contains more than three labels, we treat it as a fqdn.
func spliceDestinationSvc(
	entry *data_accesslog.HTTPAccessLogEntry,
	sourceSvcs []string,
	svcToIpsCache *SvcToIpsCache,
	fenceNN *types.NamespacedName,
	enableShortDomain bool,
) []string {
	log = log.WithField("reporter", "accesslog convertor").WithField("function", "spliceDestinationSvc")
	var destSvcs []string
	log.Debugf("sourceSvcs: %+v, svcToIpsCache: %+v, fenceNN: %+v, enableShortDomain: %+v", sourceSvcs, svcToIpsCache, fenceNN, enableShortDomain)

	upstreamCluster := entry.CommonProperties.UpstreamCluster
	log.Debugf("upstreamCluster: %+v", upstreamCluster)
	parts := strings.Split(upstreamCluster, "|")
	if len(parts) != 4 {
		log.Warnf("UpstreamCluster is wrong: parts number is not 4, skip")
		return destSvcs
	}
	// only handle inbound access log
	if parts[0] != "inbound" {
		log.Debugf("This log is not inbound, skip")
		return destSvcs
	}
	// get destination service info from request.authority
	auth := entry.Request.Authority
	dest := strings.Split(auth, ":")[0]

	// dest is ip address, skip
	if net.ParseIP(dest) != nil {
		log.Debugf("Accesslog is %+v -> %s, in which the destination is not domain, skip", sourceSvcs, dest)
		return destSvcs
	}

	// both short name and k8s fqdn will be added as following

	// expand domain as step2
	var destSvc string
	destParts := strings.Split(dest, ".")

	switch len(destParts) {
	case 1:
		destSvc = completeDestSvcWithDestName(dest, sourceSvcs, fenceNN)
	case 2:
		destSvc = completeDestSvcName(destParts, dest, "svc.cluster.local", svcToIpsCache)
	case 3:
		if destParts[2] == "svc" {
			destSvc = completeDestSvcName(destParts, dest, "cluster.local", svcToIpsCache)
		} else {
			destSvc = dest
		}
	default:
		destSvc = dest
	}

	destSvcs = append(destSvcs, destSvc)
	// We don't know if it is a short domain，
	// So we add the original dest to slice if it is not same with the parsed
	if enableShortDomain {
		if dest != destSvc {
			destSvcs = append(destSvcs, dest)
		}
	}

	result := make([]string, 0)
	for _, svc := range destSvcs {
		result = append(result, fmt.Sprintf("{destination_service=\"%s\"}", svc))
	}
	log.Infof("DestinationSvc is: %+v", result)
	return result
}

func completeDestSvcName(destParts []string, dest, suffix string, svcToIpsCache *SvcToIpsCache) (destSvc string) {
	svcToIpsCache.RLock()
	defer svcToIpsCache.RUnlock()

	svc := destParts[1] + "/" + destParts[0]
	if _, ok := svcToIpsCache.Data[svc]; ok {
		// dest is abbreviation of service, add suffix
		destSvc = dest + "." + suffix
	} else {
		// not abbreviation of service, no suffix
		destSvc = dest
	}
	return
}

// exact dest ns from sourceSvc and fence，otherwise return the original value
func completeDestSvcWithDestName(dest string, sourceSvcs []string, fenceNN *types.NamespacedName) (destSvc string) {
	destSvc = dest
	if len(sourceSvcs) > 0 {
		srcParts := strings.Split(sourceSvcs[0], "/")
		if len(srcParts) == 2 {
			destSvc = dest + "." + srcParts[0] + ".svc.cluster.local"
		}
	} else if fenceNN != nil {
		destSvc = dest + "." + fenceNN.Namespace + ".svc.cluster.local"
	}
	return
}
