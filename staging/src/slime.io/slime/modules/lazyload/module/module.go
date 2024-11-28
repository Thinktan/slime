package module

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
	istionetworkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	basecontroller "slime.io/slime/framework/controllers"
	"slime.io/slime/framework/model/metric"
	"slime.io/slime/framework/model/module"
	"slime.io/slime/modules/lazyload/api/config"
	lazyloadapiv1alpha1 "slime.io/slime/modules/lazyload/api/v1alpha1"
	"slime.io/slime/modules/lazyload/controllers"
	modmodel "slime.io/slime/modules/lazyload/model"
	"slime.io/slime/modules/lazyload/pkg/server"
)

var log = modmodel.ModuleLog

type Module struct {
	config config.Fence
}

func (m *Module) Kind() string {
	return modmodel.ModuleName
}

func (m *Module) Config() proto.Message {
	return &m.config
}

func (m *Module) InitScheme(scheme *runtime.Scheme) error {
	for _, f := range []func(*runtime.Scheme) error{
		clientgoscheme.AddToScheme,          // 注册 Kubernetes 内置资源，如 Pod、Service 等
		lazyloadapiv1alpha1.AddToScheme,     // 注册 Lazyload 的自定义资源 ServiceFence
		istionetworkingv1alpha3.AddToScheme, // 注册 Istio 的资源，如 VirtualService
	} {
		if err := f(scheme); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) Clone() module.Module {
	ret := *m
	return &ret
}

func (m *Module) Setup(opts module.ModuleOptions) error {
	log = log.WithField("addby", "tklog")
	log.Infof("lazyload setup begin")

	env, mgr, le := opts.Env, opts.Manager, opts.LeaderElectionCbs
	pc, err := controllers.NewProducerConfig(env, m.config)
	if err != nil {
		return fmt.Errorf("unable to create ProducerConfig, %+v", err)
	}
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

	// -- env: {
	//	Config:global:{service:"app" istioNamespace:"istio-system" slimeNamespace:"mesh-operator"
	//	log:{logLevel:"info" klogLevel:5} misc:{key:"aux-addr" value:":8081"}
	//	misc:{key:"enableLeaderElection" value:"off"} misc:{key:"logSourcePort" value:":8082"}
	//	misc:{key:"metrics-addr" value:":8080"} misc:{key:"pathRedirect" value:""}
	//	misc:{key:"seLabelSelectorKeys" value:"app"} misc:{key:"xdsSourceEnableIncPush" value:"true"}}
	//	name:"lazyload" enable:true general:{} kind:"lazyload" K8SClient:0xc0002ff040
	//	DynamicClient:0xc00093abd0 HttpPathHandler:{Prefix:lazyload PathHandler:0xc000939a40}
	//	ReadyManager:0x2224c20 Stop:0xc0005b6930 ConfigController:<nil> IstioConfigController:<nil>},

	// 通过sfReconciler注册与ServiceFence资源相关的控制器
	// 设置nsSvcCache等一系列结构题
	// 创建一个 sfReconciler 实例，用于调谐 ServiceFence 资源。
	sfReconciler := controllers.NewReconciler(
		controllers.ReconcilerWithCfg(&m.config),
		controllers.ReconcilerWithEnv(env),
		controllers.ReconcilerWithProducerConfig(pc), // 这里会设置好pc的watcher/ticker/accesslog相关的handler函数NeedUpdateMetricHandler
	)
	sfReconciler.Client = mgr.GetClient() // 获取 Kubernetes 的客户端，用于与 API Server 通信。
	sfReconciler.Scheme = mgr.GetScheme() // 获取当前控制器的资源类型注册表（runtime.Scheme）

	if env.ConfigController != nil {
		sfReconciler.RegisterSeHandler()
	}

	podNs := os.Getenv("WATCH_NAMESPACE")
	podName := os.Getenv("POD_NAME")

	opts.InitCbs.AddStartup(func(ctx context.Context) {
		log.Infof("AddStartup")
		sfReconciler.StartCache(ctx)
		if env.Config.Global != nil && env.Config.Global.Misc["enableLeaderElection"] == "on" {
			log.Infof("delete leader labels before working")
			deleteLeaderLabelUntilSucceed(env.K8SClient, podNs, podName)
		}
	})

	// build metric source
	// 通过 metric.NewSource 构建指标源。
	source := metric.NewSource(pc)

	// 从 ServiceFence 中加载流量信息到缓存，供动态流量控制使用。
	cache, err := controllers.NewCache(env)
	log.Infof("NewCache return: %+v", cache)
	// -- cache
	// istio-system/istio-egressgateway:map[]
	// istio-system/istio-ingressgateway:map[]
	//	istio-system/istiod:map[] kube-system/kube-dns:map[]
	//	mesh-operator/lazyload:map[]
	if err != nil {
		return fmt.Errorf("GetCacheFromServicefence occured err: %s", err)
	}
	_ = source.Fullfill(cache)
	log.Debugf("GetCacheFromServicefence %+v", cache)

	// register svf reset，将source设置为空
	handler := &server.Handler{
		HttpPathHandler: env.HttpPathHandler,
		Source:          source,
	}
	svfResetRegister(handler)

	var builder basecontroller.ObjectReconcilerBuilder

	// auto generate ServiceFence or not
	// 注册控制器逻辑，将 Reconcile 函数绑定到具体的资源。
	if m.config.AutoFence {
		builder = builder.Add(basecontroller.ObjectReconcileItem{
			Name:    "Namespace",
			ApiType: &corev1.Namespace{},
			R:       reconcile.Func(sfReconciler.ReconcileNamespace),
		}).Add(basecontroller.ObjectReconcileItem{
			Name:    "Service",
			ApiType: &corev1.Service{},
			R:       reconcile.Func(sfReconciler.ReconcileService),
		})
		// use FenceLabelKeyAlias ad switch to turn on/off workload fence
		log.Infof("m.config.FenceLabelKeyAlias: ", m.config.FenceLabelKeyAlias)
		if m.config.FenceLabelKeyAlias != "" {
			log.Infof("do NewPodController")
			podController := sfReconciler.NewPodController(env.K8SClient, m.config.FenceLabelKeyAlias)
			// 为模块设置leader Election机制，选举成功后，调用模块的回调逻辑
			le.AddOnStartedLeading(func(ctx context.Context) {
				go podController.Run(ctx.Done())
			})
		} else {
			log.Infof("not NewPodController") // --> do here
		}
	}

	builder = builder.Add(basecontroller.ObjectReconcileItem{
		Name: "ServiceFence",
		R:    sfReconciler,
	}).Add(basecontroller.ObjectReconcileItem{
		Name: "VirtualService",
		R: &basecontroller.VirtualServiceReconciler{
			Env:    &env,
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		},
	})

	if err := builder.Build(mgr); err != nil {
		return fmt.Errorf("unable to create controller,%+v", err)
	}

	le.AddOnStartedLeading(func(ctx context.Context) {
		log.Infof("retrieve metric from svf status.metric")
		cache, err := controllers.NewCache(env)
		if err != nil {
			log.Warnf("GetCacheFromServicefence occured err in StartedLeading: %s", err)
			return
		}
		_ = source.Fullfill(cache)
		log.Infof("GetCacheFromServicefence is %+v", cache)
	})

	le.AddOnStartedLeading(func(ctx context.Context) {
		log.Infof("producers starts")
		metric.NewProducer(pc, source)
	})

	if m.config.AutoPort {
		// 启用自动端口监控
		le.AddOnStartedLeading(func(ctx context.Context) {
			sfReconciler.StartAutoPort(ctx)
		})
	}

	if env.Config.Metric != nil ||
		m.config.MetricSourceType == controllers.MetricSourceTypeAccesslog {
		// 启动指标采集
		// 使用 sfReconciler.WatchMetric 动态监听访问流量
		le.AddOnStartedLeading(func(ctx context.Context) {
			go sfReconciler.WatchMetric(ctx)
		})
	} else {
		log.Warningf("watching metric is not running")
	}

	// 在 Leader Election 场景下，动态管理 Leader Pod 的标签。
	if env.Config.Global != nil && env.Config.Global.Misc["enableLeaderElection"] == "on" {
		log.Infof("add/delete leader label in StartedLeading/stoppedLeading")
		le.AddOnStartedLeading(func(ctx context.Context) {
			log.Infof("Leading management leader label")
			first := make(chan struct{}, 1)
			first <- struct{}{}
			var retry <-chan time.Time

			go func() {
				for {
					select {
					case <-ctx.Done():
						log.Infof("ctx is done, retrun")
						return
					case <-first:
					case <-retry:
						retry = nil
					}
					if err = addPodLabel(ctx, env.K8SClient, podNs, podName); err != nil {
						log.Errorf("add leader labels error %s, retry", err)
						retry = time.After(1 * time.Second)
					} else {
						log.Infof("add leader labels succeed")
						return
					}
				}
			}()
		})

		le.AddOnStoppedLeading(func() {
			go deleteLeaderLabelUntilSucceed(env.K8SClient, podNs, podName)
		})
	}

	le.AddOnStoppedLeading(sfReconciler.Clear)
	return nil
}

func svfResetRegister(handler *server.Handler) {
	handler.HandleFunc("/debug/svfReset", handler.SvfResetSetting)
}

func deleteLeaderLabelUntilSucceed(client *kubernetes.Clientset, podNs, podName string) {
	first := make(chan struct{}, 1)
	first <- struct{}{}
	var retry <-chan time.Time
	for {
		select {
		case <-first:
		case <-retry:
		}

		if err := deletePodLabel(context.TODO(), client, podNs, podName); err != nil {
			log.Errorf("delete leader labels error %s", err)
			retry = time.After(1 * time.Second)
		} else {
			log.Infof("delete leader labels succeed")
			return
		}
	}
}

func addPodLabel(ctx context.Context, client *kubernetes.Clientset, podNs, podName string) error {
	po, err := getPod(ctx, client, podNs, podName)
	if err != nil {
		return err
	}

	po.Labels[modmodel.SlimeLeader] = "true"
	_, err = client.CoreV1().Pods(podNs).Update(ctx, po, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("update pod namespace/name: %s/%s err %s", podNs, podName, err)
	}
	return nil
}

func deletePodLabel(ctx context.Context, client *kubernetes.Clientset, podNs, podName string) error {
	po, err := getPod(ctx, client, podNs, podName)
	if err != nil {
		return err
	}
	// if slime.io/leader not exist, skip
	if _, ok := po.Labels[modmodel.SlimeLeader]; !ok {
		log.Infof("label slime.io/leader is not found, skip")
		return nil
	}

	delete(po.Labels, modmodel.SlimeLeader)
	_, err = client.CoreV1().Pods(podNs).Update(ctx, po, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 通过 Kubernetes API 获取 Pod 的信息。
func getPod(ctx context.Context, client *kubernetes.Clientset, podNs, podName string) (*corev1.Pod, error) {
	pod, err := client.CoreV1().Pods(podNs).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			err = fmt.Errorf("pod %s/%s is not found", podNs, podName)
		} else {
			err = fmt.Errorf("get pod %s/%s err %s", podNs, podName, err)
		}
		return nil, err
	}
	return pod, nil
}
