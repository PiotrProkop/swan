package mutilate

import (
	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"os"
	"time"
)

const (
	NAME    = "mutilate"
	VERSION = 1
	TYPE    = plugin.CollectorPluginType
)

type mutilate struct {
	now time.Time
}

func NewMutilate(now time.Time) *mutilate {
	return &mutilate{now}
}

func (mutilate *mutilate) GetMetricTypes(configType plugin.ConfigType) ([]plugin.MetricType, error) {
	var metrics []plugin.MetricType
	phaseName, _ := config.GetConfigItem(configType, "phase_name")
	tags := map[string]string{"phase_name": phaseName.(string)}
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("avg"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("std"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("min"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "5th"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "10th"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "90th"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "95th"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "99th"), Tags_: tags, Unit_: "ns", Version_: VERSION})
	metrics = append(metrics, plugin.MetricType{Namespace_: createNewMetricNamespace("percentile", "99.999th"), Tags_: tags, Unit_: "ns", Version_: VERSION})

	return metrics, nil
}

func createNewMetricNamespace(metricName ...string) core.Namespace {
	namespace := core.NewNamespace([]string{"intel", "swan", "mutilate"})
	namespace = namespace.AddDynamicElement("hostname", "Name of the host that reports the metric")
	for _, value := range metricName {
		namespace = namespace.AddStaticElement(value)
	}

	return namespace
}

func (mutilate *mutilate) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	var metrics []plugin.MetricType
	rawMetrics, _ := parseMutilateStdout("/home/developer/go/src/github.com/intelsdi-x/swan/misc/snap-plugin-collector-mutilate/mutilate/mutilate.stdout", mutilate.now)
	hostname, _ := os.Hostname()
	for key, metricType := range metricTypes {
		metricType.Data_ = rawMetrics[key].value
		metricType.Namespace_[3].Value = hostname
		metricType.Timestamp_ = rawMetrics[key].time
		metrics = append(metrics, metricType)
	}

	return metrics, nil
}

func (mutilate *mutilate) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	phaseName, _ := cpolicy.NewStringRule("phase_name", true)
	tags := cpolicy.NewPolicyNode()
	tags.Add(phaseName)
	policy := cpolicy.New()
	policy.Add([]string{"experiment"}, tags)

	return policy, nil
}

func Meta() *plugin.PluginMeta {
	meta := plugin.NewPluginMeta(
		NAME,
		VERSION,
		TYPE,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.RoutingStrategy(plugin.DefaultRouting),
		plugin.CacheTTL(1100*time.Millisecond),
	)
	meta.RPCType = plugin.JSONRPC

	return meta
}
