// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ServerType int

const (
	Grpc ServerType = iota + 1
)

var serverTypes = map[string]ServerType{
	"grpc": Grpc,
}

func parseServerType(str string) (ServerType, bool) {
	c, ok := serverTypes[strings.ToLower(str)]
	return c, ok
}

type (
	Config struct {
		Command            *Command         `yaml:"-" mapstructure:"command"`
		Log                *Log             `yaml:"-" mapstructure:"log"`
		DataPlaneConfig    *DataPlaneConfig `yaml:"-" mapstructure:"data_plane_config"`
		Client             *Client          `yaml:"-" mapstructure:"client"`
		Collector          *Collector       `yaml:"-" mapstructure:"collector"`
		File               *File            `yaml:"-" mapstructure:"file"`
		Common             *CommonSettings  `yaml:"-"`
		Watchers           *Watchers        `yaml:"-"`
		Version            string           `yaml:"-"`
		Path               string           `yaml:"-"`
		UUID               string           `yaml:"-"`
		AllowedDirectories []string         `yaml:"-" mapstructure:"allowed_directories"`
		Features           []string         `yaml:"-"`
	}

	Log struct {
		Level string `yaml:"-" mapstructure:"level"`
		Path  string `yaml:"-" mapstructure:"path"`
	}

	DataPlaneConfig struct {
		Nginx *NginxDataPlaneConfig `yaml:"-" mapstructure:"nginx"`
	}

	NginxDataPlaneConfig struct {
		ExcludeLogs            []string      `yaml:"-" mapstructure:"exclude_logs"`
		ReloadMonitoringPeriod time.Duration `yaml:"-" mapstructure:"reload_monitoring_period"`
		TreatWarningsAsErrors  bool          `yaml:"-" mapstructure:"treat_warnings_as_errors"`
	}

	Client struct {
		Timeout             time.Duration `yaml:"-" mapstructure:"timeout"`
		Time                time.Duration `yaml:"-" mapstructure:"time"`
		PermitWithoutStream bool          `yaml:"-" mapstructure:"permit_without_stream"`
		// if MaxMessageSize is size set then we use that value,
		// otherwise MaxMessageRecieveSize and MaxMessageSendSize for individual settings
		MaxMessageSize        int `yaml:"-" mapstructure:"max_message_size"`
		MaxMessageRecieveSize int `yaml:"-" mapstructure:"max_message_receive_size"`
		MaxMessageSendSize    int `yaml:"-" mapstructure:"max_message_send_size"`
	}

	Collector struct {
		ConfigPath string     `yaml:"-" mapstructure:"config_path"`
		Log        *Log       `yaml:"-" mapstructure:"log"`
		Exporters  Exporters  `yaml:"-" mapstructure:"exporters"`
		Extensions Extensions `yaml:"-" mapstructure:"extensions"`
		Processors Processors `yaml:"-" mapstructure:"processors"`
		Receivers  Receivers  `yaml:"-" mapstructure:"receivers"`
	}

	Exporters struct {
		Debug              *DebugExporter      `yaml:"-" mapstructure:"debug"`
		PrometheusExporter *PrometheusExporter `yaml:"-" mapstructure:"prometheus_exporter"`
		OtlpExporters      []OtlpExporter      `yaml:"-" mapstructure:"otlp_exporters"`
	}

	OtlpExporter struct {
		Server        *ServerConfig `yaml:"-" mapstructure:"server"`
		TLS           *TLSConfig    `yaml:"-" mapstructure:"tls"`
		Compression   string        `yaml:"-" mapstructure:"compression"`
		Authenticator string        `yaml:"-" mapstructure:"authenticator"`
	}

	Extensions struct {
		Health        *Health        `yaml:"-" mapstructure:"health"`
		HeadersSetter *HeadersSetter `yaml:"-" mapstructure:"headers_setter"`
	}

	Health struct {
		Server *ServerConfig `yaml:"-" mapstructure:"server"`
		TLS    *TLSConfig    `yaml:"-" mapstructure:"tls"`
		Path   string        `yaml:"-" mapstructure:"path"`
	}

	HeadersSetter struct {
		Headers []Header `yaml:"-" mapstructure:"headers"`
	}

	Header struct {
		Action       string `yaml:"-" mapstructure:"action"`
		Key          string `yaml:"-" mapstructure:"key"`
		Value        string `yaml:"-" mapstructure:"value"`
		DefaultValue string `yaml:"-" mapstructure:"default_value"`
		FromContext  string `yaml:"-" mapstructure:"from_context"`
	}

	DebugExporter struct{}

	PrometheusExporter struct {
		Server *ServerConfig `yaml:"-" mapstructure:"server"`
		TLS    *TLSConfig    `yaml:"-" mapstructure:"tls"`
	}

	// OTel Collector Processors configuration.
	Processors struct {
		Attribute *Attribute `yaml:"-" mapstructure:"attribute"`
		Resource  *Resource  `yaml:"-" mapstructure:"resource"`
		Batch     *Batch     `yaml:"-" mapstructure:"batch"`
	}

	Attribute struct {
		Actions []Action `yaml:"-" mapstructure:"actions"`
	}

	Action struct {
		Key    string `yaml:"key"    mapstructure:"key"`
		Action string `yaml:"action" mapstructure:"action"`
		Value  string `yaml:"value"  mapstructure:"value"`
	}

	Resource struct {
		Attributes []ResourceAttribute `yaml:"-" mapstructure:"attributes"`
	}

	ResourceAttribute struct {
		Key    string `yaml:"key"    mapstructure:"key"`
		Action string `yaml:"action" mapstructure:"action"`
		Value  string `yaml:"value"  mapstructure:"value"`
	}

	Batch struct {
		SendBatchSize    uint32        `yaml:"-" mapstructure:"send_batch_size"`
		SendBatchMaxSize uint32        `yaml:"-" mapstructure:"send_batch_max_size"`
		Timeout          time.Duration `yaml:"-" mapstructure:"timeout"`
	}

	// OTel Collector Receiver configuration.
	Receivers struct {
		HostMetrics        *HostMetrics        `yaml:"-" mapstructure:"host_metrics"`
		OtlpReceivers      []OtlpReceiver      `yaml:"-" mapstructure:"otlp_receivers"`
		NginxReceivers     []NginxReceiver     `yaml:"-" mapstructure:"nginx_receivers"`
		NginxPlusReceivers []NginxPlusReceiver `yaml:"-" mapstructure:"nginx_plus_receivers"`
	}

	OtlpReceiver struct {
		Server        *ServerConfig  `yaml:"-" mapstructure:"server"`
		Auth          *AuthConfig    `yaml:"-" mapstructure:"auth"`
		OtlpTLSConfig *OtlpTLSConfig `yaml:"-" mapstructure:"tls"`
	}

	NginxReceiver struct {
		InstanceID string      `yaml:"-" mapstructure:"instance_id"`
		StubStatus APIDetails  `yaml:"-" mapstructure:"api_details"`
		AccessLogs []AccessLog `yaml:"-" mapstructure:"access_logs"`
	}

	APIDetails struct {
		URL      string `yaml:"-" mapstructure:"url"`
		Listen   string `yaml:"-" mapstructure:"listen"`
		Location string `yaml:"-" mapstructure:"location"`
	}

	AccessLog struct {
		FilePath  string `yaml:"-" mapstructure:"file_path"`
		LogFormat string `yaml:"-" mapstructure:"log_format"`
	}

	NginxPlusReceiver struct {
		InstanceID string     `yaml:"-" mapstructure:"instance_id"`
		PlusAPI    APIDetails `yaml:"-" mapstructure:"api_details"`
	}

	HostMetrics struct {
		Scrapers           *HostMetricsScrapers `yaml:"-" mapstructure:"scrapers"`
		CollectionInterval time.Duration        `yaml:"-" mapstructure:"collection_interval"`
		InitialDelay       time.Duration        `yaml:"-" mapstructure:"initial_delay"`
	}

	HostMetricsScrapers struct {
		CPU        *CPUScraper        `yaml:"-" mapstructure:"cpu"`
		Disk       *DiskScraper       `yaml:"-" mapstructure:"disk"`
		Filesystem *FilesystemScraper `yaml:"-" mapstructure:"filesystem"`
		Memory     *MemoryScraper     `yaml:"-" mapstructure:"memory"`
		Network    *NetworkScraper    `yaml:"-" mapstructure:"network"`
	}
	CPUScraper        struct{}
	DiskScraper       struct{}
	FilesystemScraper struct{}
	MemoryScraper     struct{}
	NetworkScraper    struct{}

	GRPC struct {
		Target         string        `yaml:"-" mapstructure:"target"`
		ConnTimeout    time.Duration `yaml:"-" mapstructure:"connection_timeout"`
		MinConnTimeout time.Duration `yaml:"-" mapstructure:"minimum_connection_timeout"`
		BackoffDelay   time.Duration `yaml:"-" mapstructure:"backoff_delay"`
	}

	Command struct {
		Server *ServerConfig `yaml:"-" mapstructure:"server"`
		Auth   *AuthConfig   `yaml:"-" mapstructure:"auth"`
		TLS    *TLSConfig    `yaml:"-" mapstructure:"tls"`
	}

	ServerConfig struct {
		Host string     `yaml:"-" mapstructure:"host"`
		Port int        `yaml:"-" mapstructure:"port"`
		Type ServerType `yaml:"-" mapstructure:"type"`
	}

	AuthConfig struct {
		Token     string `yaml:"-" mapstructure:"token"`      // literal token value, keeping for backwards-compatibility, not recommended
		TokenPath string `yaml:"-" mapstructure:"token-path"` // path to token file
	}

	TLSConfig struct {
		Cert       string `yaml:"-" mapstructure:"cert"`
		Key        string `yaml:"-" mapstructure:"key"`
		Ca         string `yaml:"-" mapstructure:"ca"`
		ServerName string `yaml:"-" mapstructure:"server_name"`
		SkipVerify bool   `yaml:"-" mapstructure:"skip_verify"`
	}

	// Specialized TLS configuration for OtlpReceiver with self-signed cert generation.
	OtlpTLSConfig struct {
		Cert                   string `yaml:"-" mapstructure:"cert"`
		Key                    string `yaml:"-" mapstructure:"key"`
		Ca                     string `yaml:"-" mapstructure:"ca"`
		ServerName             string `yaml:"-" mapstructure:"server_name"`
		ExistingCert           bool   `yaml:"-"`
		SkipVerify             bool   `yaml:"-" mapstructure:"skip_verify"`
		GenerateSelfSignedCert bool   `yaml:"-" mapstructure:"generate_self_signed_cert"`
	}

	File struct {
		Location string `yaml:"-" mapstructure:"location"`
	}

	CommonSettings struct {
		InitialInterval     time.Duration `yaml:"-" mapstructure:"initial_interval"`
		MaxInterval         time.Duration `yaml:"-" mapstructure:"max_interval"`
		MaxElapsedTime      time.Duration `yaml:"-" mapstructure:"max_elapsed_time"`
		RandomizationFactor float64       `yaml:"-" mapstructure:"randomization_factor"`
		Multiplier          float64       `yaml:"-" mapstructure:"multiplier"`
	}

	Watchers struct {
		InstanceWatcher       InstanceWatcher       `yaml:"-" mapstructure:"instance_watcher"`
		InstanceHealthWatcher InstanceHealthWatcher `yaml:"-" mapstructure:"instance_health_watcher"`
		FileWatcher           FileWatcher           `yaml:"-" mapstructure:"file_watcher"`
	}

	InstanceWatcher struct {
		MonitoringFrequency time.Duration `yaml:"-" mapstructure:"monitoring_frequency"`
	}

	InstanceHealthWatcher struct {
		MonitoringFrequency time.Duration `yaml:"-" mapstructure:"monitoring_frequency"`
	}

	FileWatcher struct {
		MonitoringFrequency time.Duration `yaml:"-" mapstructure:"monitoring_frequency"`
	}
)

func (col *Collector) Validate(allowedDirectories []string) error {
	var err error
	cleaned := filepath.Clean(col.ConfigPath)

	if !isAllowedDir(cleaned, allowedDirectories) {
		err = errors.Join(err, fmt.Errorf("collector path %s not allowed", col.ConfigPath))
	}

	for _, nginxReceiver := range col.Receivers.NginxReceivers {
		err = errors.Join(err, nginxReceiver.Validate(allowedDirectories))
	}

	return err
}

func (nr *NginxReceiver) Validate(allowedDirectories []string) error {
	var err error
	if _, uuidErr := uuid.Parse(nr.InstanceID); uuidErr != nil {
		err = errors.Join(err, errors.New("invalid nginx receiver instance ID"))
	}

	for _, al := range nr.AccessLogs {
		if !isAllowedDir(al.FilePath, allowedDirectories) {
			err = errors.Join(err, fmt.Errorf("invalid nginx receiver access log path: %s", al.FilePath))
		}

		if len(al.FilePath) != 0 {
			// The log format's double quotes must be escaped so that
			// valid YAML is produced when executing the Go template.
			al.LogFormat = strings.ReplaceAll(al.LogFormat, `"`, `\"`)
		}
	}

	return err
}

func (c *Config) IsDirectoryAllowed(directory string) bool {
	return isAllowedDir(directory, c.AllowedDirectories)
}

func (c *Config) IsFeatureEnabled(feature string) bool {
	for _, enabledFeature := range c.Features {
		if enabledFeature == feature {
			return true
		}
	}

	return false
}

func (c *Config) IsACollectorExporterConfigured() bool {
	if c.Collector == nil {
		return false
	}

	return c.Collector.Exporters.PrometheusExporter != nil ||
		c.Collector.Exporters.OtlpExporters != nil ||
		c.Collector.Exporters.Debug != nil
}

func (c *Config) AreReceiversConfigured() bool {
	if c.Collector == nil {
		return false
	}

	return c.Collector.Receivers.NginxPlusReceivers != nil ||
		c.Collector.Receivers.OtlpReceivers != nil ||
		c.Collector.Receivers.NginxReceivers != nil ||
		c.Collector.Receivers.HostMetrics != nil
}

func isAllowedDir(dir string, allowedDirs []string) bool {
	for _, allowedDirectory := range allowedDirs {
		if strings.HasPrefix(dir, allowedDirectory) {
			return true
		}
	}

	return false
}
