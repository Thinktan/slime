// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: config.proto

package v1alpha1

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Limiter_RateLimitBackend int32

const (
	Limiter_netEaseLocalFlowControl Limiter_RateLimitBackend = 0
	Limiter_envoyLocalRateLimit     Limiter_RateLimitBackend = 1
)

var Limiter_RateLimitBackend_name = map[int32]string{
	0: "netEaseLocalFlowControl",
	1: "envoyLocalRateLimit",
}

var Limiter_RateLimitBackend_value = map[string]int32{
	"netEaseLocalFlowControl": 0,
	"envoyLocalRateLimit":     1,
}

func (x Limiter_RateLimitBackend) String() string {
	return proto.EnumName(Limiter_RateLimitBackend_name, int32(x))
}

func (Limiter_RateLimitBackend) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3, 0}
}

type Prometheus_Source_Type int32

const (
	Prometheus_Source_Value Prometheus_Source_Type = 0
	Prometheus_Source_Group Prometheus_Source_Type = 1
)

var Prometheus_Source_Type_name = map[int32]string{
	0: "Value",
	1: "Group",
}

var Prometheus_Source_Type_value = map[string]int32{
	"Value": 0,
	"Group": 1,
}

func (x Prometheus_Source_Type) String() string {
	return proto.EnumName(Prometheus_Source_Type_name, int32(x))
}

func (Prometheus_Source_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{6, 0}
}

type LocalSource struct {
	Mount                string   `protobuf:"bytes,1,opt,name=mount,proto3" json:"mount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LocalSource) Reset()         { *m = LocalSource{} }
func (m *LocalSource) String() string { return proto.CompactTextString(m) }
func (*LocalSource) ProtoMessage()    {}
func (*LocalSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}
func (m *LocalSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LocalSource.Unmarshal(m, b)
}
func (m *LocalSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LocalSource.Marshal(b, m, deterministic)
}
func (m *LocalSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LocalSource.Merge(m, src)
}
func (m *LocalSource) XXX_Size() int {
	return xxx_messageInfo_LocalSource.Size(m)
}
func (m *LocalSource) XXX_DiscardUnknown() {
	xxx_messageInfo_LocalSource.DiscardUnknown(m)
}

var xxx_messageInfo_LocalSource proto.InternalMessageInfo

func (m *LocalSource) GetMount() string {
	if m != nil {
		return m.Mount
	}
	return ""
}

type RemoteSource struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoteSource) Reset()         { *m = RemoteSource{} }
func (m *RemoteSource) String() string { return proto.CompactTextString(m) }
func (*RemoteSource) ProtoMessage()    {}
func (*RemoteSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}
func (m *RemoteSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoteSource.Unmarshal(m, b)
}
func (m *RemoteSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoteSource.Marshal(b, m, deterministic)
}
func (m *RemoteSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteSource.Merge(m, src)
}
func (m *RemoteSource) XXX_Size() int {
	return xxx_messageInfo_RemoteSource.Size(m)
}
func (m *RemoteSource) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteSource.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteSource proto.InternalMessageInfo

func (m *RemoteSource) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type Plugin struct {
	// Types that are valid to be assigned to WasmSource:
	//	*Plugin_Local
	//	*Plugin_Remote
	WasmSource           isPlugin_WasmSource `protobuf_oneof:"wasm_source"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Plugin) Reset()         { *m = Plugin{} }
func (m *Plugin) String() string { return proto.CompactTextString(m) }
func (*Plugin) ProtoMessage()    {}
func (*Plugin) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{2}
}
func (m *Plugin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Plugin.Unmarshal(m, b)
}
func (m *Plugin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Plugin.Marshal(b, m, deterministic)
}
func (m *Plugin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Plugin.Merge(m, src)
}
func (m *Plugin) XXX_Size() int {
	return xxx_messageInfo_Plugin.Size(m)
}
func (m *Plugin) XXX_DiscardUnknown() {
	xxx_messageInfo_Plugin.DiscardUnknown(m)
}

var xxx_messageInfo_Plugin proto.InternalMessageInfo

type isPlugin_WasmSource interface {
	isPlugin_WasmSource()
}

type Plugin_Local struct {
	Local *LocalSource `protobuf:"bytes,2,opt,name=local,proto3,oneof" json:"local,omitempty"`
}
type Plugin_Remote struct {
	Remote *RemoteSource `protobuf:"bytes,3,opt,name=remote,proto3,oneof" json:"remote,omitempty"`
}

func (*Plugin_Local) isPlugin_WasmSource()  {}
func (*Plugin_Remote) isPlugin_WasmSource() {}

func (m *Plugin) GetWasmSource() isPlugin_WasmSource {
	if m != nil {
		return m.WasmSource
	}
	return nil
}

func (m *Plugin) GetLocal() *LocalSource {
	if x, ok := m.GetWasmSource().(*Plugin_Local); ok {
		return x.Local
	}
	return nil
}

func (m *Plugin) GetRemote() *RemoteSource {
	if x, ok := m.GetWasmSource().(*Plugin_Remote); ok {
		return x.Remote
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Plugin) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Plugin_Local)(nil),
		(*Plugin_Remote)(nil),
	}
}

type Limiter struct {
	Multicluster         bool                     `protobuf:"varint,2,opt,name=Multicluster,proto3" json:"Multicluster,omitempty"`
	Backend              Limiter_RateLimitBackend `protobuf:"varint,3,opt,name=backend,proto3,enum=slime.config.v1alpha1.Limiter_RateLimitBackend" json:"backend,omitempty"`
	Refresh              *time.Duration           `protobuf:"bytes,4,opt,name=refresh,proto3,stdduration" json:"refresh,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Limiter) Reset()         { *m = Limiter{} }
func (m *Limiter) String() string { return proto.CompactTextString(m) }
func (*Limiter) ProtoMessage()    {}
func (*Limiter) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
}
func (m *Limiter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Limiter.Unmarshal(m, b)
}
func (m *Limiter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Limiter.Marshal(b, m, deterministic)
}
func (m *Limiter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Limiter.Merge(m, src)
}
func (m *Limiter) XXX_Size() int {
	return xxx_messageInfo_Limiter.Size(m)
}
func (m *Limiter) XXX_DiscardUnknown() {
	xxx_messageInfo_Limiter.DiscardUnknown(m)
}

var xxx_messageInfo_Limiter proto.InternalMessageInfo

func (m *Limiter) GetMulticluster() bool {
	if m != nil {
		return m.Multicluster
	}
	return false
}

func (m *Limiter) GetBackend() Limiter_RateLimitBackend {
	if m != nil {
		return m.Backend
	}
	return Limiter_netEaseLocalFlowControl
}

func (m *Limiter) GetRefresh() *time.Duration {
	if m != nil {
		return m.Refresh
	}
	return nil
}

type Global struct {
	Service              string   `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Multicluster         string   `protobuf:"bytes,2,opt,name=multicluster,proto3" json:"multicluster,omitempty"`
	IstioNamespace       string   `protobuf:"bytes,3,opt,name=istioNamespace,proto3" json:"istioNamespace,omitempty"`
	SlimeNamespace       string   `protobuf:"bytes,4,opt,name=slimeNamespace,proto3" json:"slimeNamespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Global) Reset()         { *m = Global{} }
func (m *Global) String() string { return proto.CompactTextString(m) }
func (*Global) ProtoMessage()    {}
func (*Global) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{4}
}
func (m *Global) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Global.Unmarshal(m, b)
}
func (m *Global) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Global.Marshal(b, m, deterministic)
}
func (m *Global) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Global.Merge(m, src)
}
func (m *Global) XXX_Size() int {
	return xxx_messageInfo_Global.Size(m)
}
func (m *Global) XXX_DiscardUnknown() {
	xxx_messageInfo_Global.DiscardUnknown(m)
}

var xxx_messageInfo_Global proto.InternalMessageInfo

func (m *Global) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Global) GetMulticluster() string {
	if m != nil {
		return m.Multicluster
	}
	return ""
}

func (m *Global) GetIstioNamespace() string {
	if m != nil {
		return m.IstioNamespace
	}
	return ""
}

func (m *Global) GetSlimeNamespace() string {
	if m != nil {
		return m.SlimeNamespace
	}
	return ""
}

type Fence struct {
	WormholePort         []string `protobuf:"bytes,2,rep,name=wormholePort,proto3" json:"wormholePort,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Fence) Reset()         { *m = Fence{} }
func (m *Fence) String() string { return proto.CompactTextString(m) }
func (*Fence) ProtoMessage()    {}
func (*Fence) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{5}
}
func (m *Fence) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Fence.Unmarshal(m, b)
}
func (m *Fence) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Fence.Marshal(b, m, deterministic)
}
func (m *Fence) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Fence.Merge(m, src)
}
func (m *Fence) XXX_Size() int {
	return xxx_messageInfo_Fence.Size(m)
}
func (m *Fence) XXX_DiscardUnknown() {
	xxx_messageInfo_Fence.DiscardUnknown(m)
}

var xxx_messageInfo_Fence proto.InternalMessageInfo

func (m *Fence) GetWormholePort() []string {
	if m != nil {
		return m.WormholePort
	}
	return nil
}

type Prometheus_Source struct {
	Address              string                                `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Handlers             map[string]*Prometheus_Source_Handler `protobuf:"bytes,2,rep,name=handlers,proto3" json:"handlers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                              `json:"-"`
	XXX_unrecognized     []byte                                `json:"-"`
	XXX_sizecache        int32                                 `json:"-"`
}

func (m *Prometheus_Source) Reset()         { *m = Prometheus_Source{} }
func (m *Prometheus_Source) String() string { return proto.CompactTextString(m) }
func (*Prometheus_Source) ProtoMessage()    {}
func (*Prometheus_Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{6}
}
func (m *Prometheus_Source) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Prometheus_Source.Unmarshal(m, b)
}
func (m *Prometheus_Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Prometheus_Source.Marshal(b, m, deterministic)
}
func (m *Prometheus_Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Prometheus_Source.Merge(m, src)
}
func (m *Prometheus_Source) XXX_Size() int {
	return xxx_messageInfo_Prometheus_Source.Size(m)
}
func (m *Prometheus_Source) XXX_DiscardUnknown() {
	xxx_messageInfo_Prometheus_Source.DiscardUnknown(m)
}

var xxx_messageInfo_Prometheus_Source proto.InternalMessageInfo

func (m *Prometheus_Source) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Prometheus_Source) GetHandlers() map[string]*Prometheus_Source_Handler {
	if m != nil {
		return m.Handlers
	}
	return nil
}

type Prometheus_Source_Handler struct {
	Query                string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Type                 Prometheus_Source_Type `protobuf:"varint,2,opt,name=type,proto3,enum=slime.config.v1alpha1.Prometheus_Source_Type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *Prometheus_Source_Handler) Reset()         { *m = Prometheus_Source_Handler{} }
func (m *Prometheus_Source_Handler) String() string { return proto.CompactTextString(m) }
func (*Prometheus_Source_Handler) ProtoMessage()    {}
func (*Prometheus_Source_Handler) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{6, 0}
}
func (m *Prometheus_Source_Handler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Prometheus_Source_Handler.Unmarshal(m, b)
}
func (m *Prometheus_Source_Handler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Prometheus_Source_Handler.Marshal(b, m, deterministic)
}
func (m *Prometheus_Source_Handler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Prometheus_Source_Handler.Merge(m, src)
}
func (m *Prometheus_Source_Handler) XXX_Size() int {
	return xxx_messageInfo_Prometheus_Source_Handler.Size(m)
}
func (m *Prometheus_Source_Handler) XXX_DiscardUnknown() {
	xxx_messageInfo_Prometheus_Source_Handler.DiscardUnknown(m)
}

var xxx_messageInfo_Prometheus_Source_Handler proto.InternalMessageInfo

func (m *Prometheus_Source_Handler) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *Prometheus_Source_Handler) GetType() Prometheus_Source_Type {
	if m != nil {
		return m.Type
	}
	return Prometheus_Source_Value
}

type K8S_Source struct {
	Handlers             []string `protobuf:"bytes,1,rep,name=handlers,proto3" json:"handlers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *K8S_Source) Reset()         { *m = K8S_Source{} }
func (m *K8S_Source) String() string { return proto.CompactTextString(m) }
func (*K8S_Source) ProtoMessage()    {}
func (*K8S_Source) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{7}
}
func (m *K8S_Source) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_K8S_Source.Unmarshal(m, b)
}
func (m *K8S_Source) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_K8S_Source.Marshal(b, m, deterministic)
}
func (m *K8S_Source) XXX_Merge(src proto.Message) {
	xxx_messageInfo_K8S_Source.Merge(m, src)
}
func (m *K8S_Source) XXX_Size() int {
	return xxx_messageInfo_K8S_Source.Size(m)
}
func (m *K8S_Source) XXX_DiscardUnknown() {
	xxx_messageInfo_K8S_Source.DiscardUnknown(m)
}

var xxx_messageInfo_K8S_Source proto.InternalMessageInfo

func (m *K8S_Source) GetHandlers() []string {
	if m != nil {
		return m.Handlers
	}
	return nil
}

type Metric struct {
	Prometheus           *Prometheus_Source `protobuf:"bytes,1,opt,name=prometheus,proto3" json:"prometheus,omitempty"`
	K8S                  *K8S_Source        `protobuf:"bytes,2,opt,name=k8s,proto3" json:"k8s,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{8}
}
func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (m *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(m, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetPrometheus() *Prometheus_Source {
	if m != nil {
		return m.Prometheus
	}
	return nil
}

func (m *Metric) GetK8S() *K8S_Source {
	if m != nil {
		return m.K8S
	}
	return nil
}

type Config struct {
	Plugin               *Plugin  `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
	Limiter              *Limiter `protobuf:"bytes,2,opt,name=limiter,proto3" json:"limiter,omitempty"`
	Global               *Global  `protobuf:"bytes,3,opt,name=global,proto3" json:"global,omitempty"`
	Fence                *Fence   `protobuf:"bytes,4,opt,name=fence,proto3" json:"fence,omitempty"`
	Metric               *Metric  `protobuf:"bytes,6,opt,name=metric,proto3" json:"metric,omitempty"`
	Name                 string   `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Enable               bool     `protobuf:"varint,7,opt,name=enable,proto3" json:"enable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{9}
}
func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetPlugin() *Plugin {
	if m != nil {
		return m.Plugin
	}
	return nil
}

func (m *Config) GetLimiter() *Limiter {
	if m != nil {
		return m.Limiter
	}
	return nil
}

func (m *Config) GetGlobal() *Global {
	if m != nil {
		return m.Global
	}
	return nil
}

func (m *Config) GetFence() *Fence {
	if m != nil {
		return m.Fence
	}
	return nil
}

func (m *Config) GetMetric() *Metric {
	if m != nil {
		return m.Metric
	}
	return nil
}

func (m *Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Config) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func init() {
	proto.RegisterEnum("slime.config.v1alpha1.Limiter_RateLimitBackend", Limiter_RateLimitBackend_name, Limiter_RateLimitBackend_value)
	proto.RegisterEnum("slime.config.v1alpha1.Prometheus_Source_Type", Prometheus_Source_Type_name, Prometheus_Source_Type_value)
	proto.RegisterType((*LocalSource)(nil), "slime.config.v1alpha1.LocalSource")
	proto.RegisterType((*RemoteSource)(nil), "slime.config.v1alpha1.RemoteSource")
	proto.RegisterType((*Plugin)(nil), "slime.config.v1alpha1.Plugin")
	proto.RegisterType((*Limiter)(nil), "slime.config.v1alpha1.Limiter")
	proto.RegisterType((*Global)(nil), "slime.config.v1alpha1.Global")
	proto.RegisterType((*Fence)(nil), "slime.config.v1alpha1.Fence")
	proto.RegisterType((*Prometheus_Source)(nil), "slime.config.v1alpha1.Prometheus_Source")
	proto.RegisterMapType((map[string]*Prometheus_Source_Handler)(nil), "slime.config.v1alpha1.Prometheus_Source.HandlersEntry")
	proto.RegisterType((*Prometheus_Source_Handler)(nil), "slime.config.v1alpha1.Prometheus_Source.Handler")
	proto.RegisterType((*K8S_Source)(nil), "slime.config.v1alpha1.K8S_Source")
	proto.RegisterType((*Metric)(nil), "slime.config.v1alpha1.Metric")
	proto.RegisterType((*Config)(nil), "slime.config.v1alpha1.Config")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4) }

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 779 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdd, 0xae, 0xdb, 0x44,
	0x10, 0x8e, 0xf3, 0xe3, 0x34, 0x93, 0xd3, 0xa3, 0xb0, 0xfc, 0xd4, 0x84, 0x52, 0x0e, 0xae, 0x04,
	0x41, 0xa8, 0x36, 0x4d, 0x05, 0x0a, 0x95, 0xb8, 0x20, 0xa5, 0xa7, 0x41, 0xb4, 0xe8, 0x68, 0x8b,
	0xb8, 0xe0, 0xa6, 0xda, 0x38, 0x13, 0xc7, 0xca, 0xda, 0x6b, 0xd6, 0xeb, 0x1c, 0xe5, 0x09, 0xb8,
	0xe7, 0x02, 0xf1, 0x46, 0x3c, 0x06, 0xef, 0xc0, 0x3d, 0x12, 0xf2, 0xee, 0x3a, 0x27, 0x39, 0x90,
	0xd2, 0xde, 0xcd, 0x8c, 0xbf, 0x6f, 0xbe, 0xd9, 0x19, 0xcf, 0xc0, 0x49, 0x24, 0xb2, 0x65, 0x12,
	0x07, 0xb9, 0x14, 0x4a, 0x90, 0xb7, 0x0b, 0x9e, 0xa4, 0x18, 0xd8, 0xd8, 0xe6, 0x3e, 0xe3, 0xf9,
	0x8a, 0xdd, 0x1f, 0xde, 0x8b, 0x13, 0xb5, 0x2a, 0xe7, 0x41, 0x24, 0xd2, 0x30, 0x16, 0xb1, 0x08,
	0x35, 0x7a, 0x5e, 0x2e, 0xb5, 0xa7, 0x1d, 0x6d, 0x99, 0x2c, 0xc3, 0x3b, 0xb1, 0x10, 0x31, 0xc7,
	0x2b, 0xd4, 0xa2, 0x94, 0x4c, 0x25, 0x22, 0x33, 0xdf, 0xfd, 0xbb, 0xd0, 0x7f, 0x2a, 0x22, 0xc6,
	0x9f, 0x8b, 0x52, 0x46, 0x48, 0xde, 0x82, 0x4e, 0x2a, 0xca, 0x4c, 0x79, 0xce, 0x99, 0x33, 0xea,
	0x51, 0xe3, 0xf8, 0x23, 0x38, 0xa1, 0x98, 0x0a, 0x85, 0x16, 0xe5, 0x41, 0x97, 0x2d, 0x16, 0x12,
	0x8b, 0xc2, 0xe2, 0x6a, 0xd7, 0xff, 0xd5, 0x01, 0xf7, 0x82, 0x97, 0x71, 0x92, 0x91, 0x87, 0xd0,
	0xe1, 0x55, 0x66, 0xaf, 0x79, 0xe6, 0x8c, 0xfa, 0x63, 0x3f, 0xf8, 0xcf, 0xf7, 0x04, 0x7b, 0xea,
	0xb3, 0x06, 0x35, 0x14, 0xf2, 0x15, 0xb8, 0x52, 0x0b, 0x7a, 0x2d, 0x4d, 0xbe, 0x7b, 0x84, 0xbc,
	0x5f, 0xd5, 0xac, 0x41, 0x2d, 0x69, 0x7a, 0x13, 0xfa, 0x97, 0xac, 0x48, 0x5f, 0x14, 0xfa, 0x83,
	0xff, 0xb7, 0x03, 0xdd, 0xa7, 0x49, 0x9a, 0x28, 0x94, 0xc4, 0x87, 0x93, 0x67, 0x25, 0x57, 0x49,
	0xc4, 0xcb, 0x42, 0xa1, 0xd4, 0xc5, 0xdd, 0xa0, 0x07, 0x31, 0xf2, 0x2d, 0x74, 0xe7, 0x2c, 0x5a,
	0x63, 0xb6, 0xd0, 0xf2, 0xa7, 0xe3, 0xf0, 0x58, 0xed, 0x26, 0x69, 0x40, 0x99, 0x42, 0x6d, 0x4f,
	0x0d, 0x8d, 0xd6, 0x7c, 0xf2, 0x25, 0x74, 0x25, 0x2e, 0x25, 0x16, 0x2b, 0xaf, 0xad, 0x5f, 0xf2,
	0x6e, 0x60, 0x06, 0x12, 0xd4, 0x03, 0x09, 0xbe, 0xb1, 0x03, 0x99, 0xb6, 0x7f, 0xff, 0xf3, 0x03,
	0x87, 0xd6, 0x78, 0x7f, 0x06, 0x83, 0xeb, 0x79, 0xc9, 0x7b, 0x70, 0x2b, 0x43, 0xf5, 0x98, 0x15,
	0xa8, 0xdb, 0x76, 0xce, 0xc5, 0xe5, 0x23, 0x91, 0x29, 0x29, 0xf8, 0xa0, 0x41, 0x6e, 0xc1, 0x9b,
	0x98, 0x6d, 0xc4, 0x56, 0x7f, 0xda, 0x51, 0x07, 0x8e, 0xff, 0x9b, 0x03, 0xee, 0x13, 0x2e, 0xe6,
	0x8c, 0x57, 0x93, 0x2b, 0x50, 0x6e, 0x92, 0x08, 0xeb, 0xc9, 0x59, 0xb7, 0x6a, 0x4c, 0x7a, 0xbd,
	0x31, 0x3d, 0x7a, 0x10, 0x23, 0x1f, 0xc1, 0x69, 0x52, 0xa8, 0x44, 0x7c, 0xcf, 0x52, 0x2c, 0x72,
	0x16, 0x99, 0xf1, 0xf4, 0xe8, 0xb5, 0x68, 0x85, 0xd3, 0x0d, 0xbb, 0xc2, 0xb5, 0x0d, 0xee, 0x30,
	0xea, 0x7f, 0x0a, 0x9d, 0x73, 0xcc, 0x8c, 0xf8, 0xa5, 0x90, 0xe9, 0x4a, 0x70, 0xbc, 0x10, 0x52,
	0x79, 0xcd, 0xb3, 0x56, 0x25, 0xbe, 0x1f, 0xf3, 0xff, 0x6a, 0xc2, 0x1b, 0x17, 0x52, 0xa4, 0xa8,
	0x56, 0x58, 0x16, 0x2f, 0xfe, 0xef, 0x57, 0x24, 0x14, 0x6e, 0xac, 0x58, 0xb6, 0xe0, 0x28, 0x0b,
	0x9d, 0xaf, 0x3f, 0xfe, 0xe2, 0xc8, 0x18, 0xff, 0x95, 0x35, 0x98, 0x59, 0xe2, 0xe3, 0x4c, 0xc9,
	0x2d, 0xdd, 0xe5, 0x19, 0xce, 0xa1, 0x6b, 0x3f, 0x55, 0x9b, 0xf2, 0x73, 0x89, 0x72, 0x5b, 0x6f,
	0x8a, 0x76, 0xc8, 0xd7, 0xd0, 0x56, 0xdb, 0x1c, 0x75, 0xf7, 0x4e, 0xc7, 0xf7, 0x5e, 0x59, 0xf0,
	0x87, 0x6d, 0x8e, 0x54, 0x53, 0x87, 0x29, 0xdc, 0x3c, 0x90, 0x27, 0x03, 0x68, 0xad, 0xb1, 0xd6,
	0xa9, 0x4c, 0x72, 0x0e, 0x9d, 0x0d, 0xe3, 0x25, 0xda, 0xd5, 0xfa, 0xec, 0x75, 0xdf, 0x45, 0x0d,
	0xfd, 0x61, 0x73, 0xe2, 0xf8, 0xb7, 0xa1, 0x5d, 0x89, 0x93, 0x1e, 0x74, 0x7e, 0xac, 0x82, 0x83,
	0x46, 0x65, 0x3e, 0x91, 0xa2, 0xcc, 0x07, 0x8e, 0x3f, 0x02, 0xf8, 0x6e, 0xf2, 0xbc, 0x6e, 0xf6,
	0x70, 0xaf, 0xa5, 0x8e, 0x1e, 0xd1, 0xce, 0xf7, 0x7f, 0x71, 0xc0, 0x7d, 0x86, 0x4a, 0x26, 0x11,
	0x99, 0x01, 0xe4, 0x3b, 0x69, 0x5d, 0x77, 0x7f, 0x3c, 0x7a, 0xd5, 0x1a, 0xe9, 0x1e, 0x97, 0x3c,
	0x80, 0xd6, 0x7a, 0x52, 0xd8, 0x67, 0x7e, 0x78, 0x24, 0xc5, 0x55, 0x81, 0xb4, 0x42, 0xfb, 0x7f,
	0x34, 0xc1, 0x7d, 0xa4, 0x31, 0xe4, 0x73, 0x70, 0x73, 0x7d, 0x8d, 0x6c, 0x15, 0xef, 0x1f, 0xab,
	0x42, 0x83, 0xa8, 0x05, 0x93, 0x09, 0x74, 0xb9, 0x59, 0x6d, 0x2b, 0x7d, 0xe7, 0xe5, 0x07, 0x80,
	0xd6, 0xf0, 0x4a, 0x30, 0xd6, 0x9b, 0x66, 0x0f, 0xd7, 0x31, 0x41, 0xb3, 0x8e, 0xd4, 0x82, 0xc9,
	0x18, 0x3a, 0xcb, 0x6a, 0x11, 0xec, 0x91, 0xb8, 0x7d, 0x84, 0xa5, 0x97, 0x85, 0x1a, 0x68, 0x25,
	0x95, 0xea, 0x7e, 0x7b, 0xee, 0x4b, 0xa5, 0xcc, 0x50, 0xa8, 0x05, 0x13, 0x02, 0xed, 0x8c, 0xa5,
	0xe8, 0x75, 0xf4, 0xef, 0xa4, 0x6d, 0xf2, 0x0e, 0xb8, 0x98, 0xb1, 0x39, 0x47, 0xaf, 0xab, 0xcf,
	0xa1, 0xf5, 0xa6, 0x9f, 0xfc, 0xf4, 0xb1, 0xc9, 0x99, 0x88, 0x50, 0x1b, 0x61, 0xbe, 0x8e, 0x43,
	0x96, 0x27, 0x45, 0x68, 0x54, 0xc2, 0x5a, 0x65, 0xee, 0xea, 0x7b, 0xf6, 0xe0, 0x9f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x4d, 0x7e, 0x73, 0xec, 0xc4, 0x06, 0x00, 0x00,
}
