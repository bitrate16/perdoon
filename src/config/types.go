package config

import (
	"fmt"
)

type ValueRange struct {
	Start int
	End   int
}

func (pr *ValueRange) String() string {
	return fmt.Sprintf("%d-%d", pr.Start, pr.End)
}

type ResponseConfig struct {
	Sizes    []*ValueRange `yaml:"sizes"`
	Bytes    string        `yaml:"bytes"`
	Strategy string        `yaml:"strategy"`
}

type TCPConfig struct {
	Ports     []*ValueRange   `yaml:"ports"`
	Exclude   []int           `yaml:"exclude"`
	ChunkSize int             `yaml:"chunk-size"`
	Response  *ResponseConfig `yaml":response"`
}

type UDPConfig struct {
	Ports     []*ValueRange   `yaml:"ports"`
	Exclude   []int           `yaml:"exclude"`
	ChunkSize int             `yaml:"chunk-size"`
	Response  *ResponseConfig `yaml":response"`
}

type RecordConfig struct {
	RequestPayload  bool `yaml:"request-payload"`
	ResponsePayload bool `yaml:"response-payload"`
}

type DatabaseConfig struct {
	Path   string       `yaml:"path"`
	Table  string       `yaml:"table"`
	Record RecordConfig `yaml:"record"`
}

type Config struct {
	TCP      TCPConfig      `yaml:"tcp"`
	UDP      UDPConfig      `yaml:"udp"`
	Database DatabaseConfig `yaml:"database"`
	Debug    bool           `yaml:"debug"`
}

func defaultConfig() *Config {
	return &Config{
		TCP: TCPConfig{
			Ports: []*ValueRange{
				{
					Start: 80,
					End:   80,
				},
			},
			ChunkSize: 1024,
			Response: &ResponseConfig{
				Sizes: []*ValueRange{
					{
						Start: 8,
						End:   1 * 1024 * 1024,
					},
				},
				Strategy: "random",
			},
		},
		UDP: UDPConfig{
			Ports: []*ValueRange{
				{
					Start: 1000,
					End:   2000,
				},
			},
			Exclude: []int{
				1500,
			},
			ChunkSize: 1024,
			Response: &ResponseConfig{
				Sizes: []*ValueRange{
					{
						Start: 100,
						End:   200,
					},
				},
				Bytes:    "3a33",
				Strategy: "bytes",
			},
		},
		Database: DatabaseConfig{
			Path:  "perdoon.db",
			Table: "track",
			Record: RecordConfig{
				RequestPayload:  true,
				ResponsePayload: true,
			},
		},
		Debug: false,
	}
}
