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

type TCPConfig struct {
	Ports     []*ValueRange `yaml:"ports"`
	Exclude   []int         `yaml:"exclude"`
	ChunkSize int           `yaml:"chunk-size"`
}

type UDPConfig struct {
	Ports     []*ValueRange `yaml:"ports"`
	Exclude   []int         `yaml:"exclude"`
	ChunkSize int           `yaml:"chunk-size"`
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

type ResponseConfig struct {
	Sizes []*ValueRange `yaml:"sizes"`
	Bytes string        `yaml:"bytes"`
}

type Config struct {
	TCP      TCPConfig      `yaml:"tcp"`
	UDP      UDPConfig      `yaml:"udp"`
	Database DatabaseConfig `yaml:"database"`
	Response ResponseConfig `yaml:"response"`
	Strategy string         `yaml:"strategy"`
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
		},
		Database: DatabaseConfig{
			Path:  "perdoon.db",
			Table: "track",
			Record: RecordConfig{
				RequestPayload:  true,
				ResponsePayload: true,
			},
		},
		Response: ResponseConfig{
			Sizes: []*ValueRange{
				{
					Start: 100,
					End:   200,
				},
			},
			Bytes: ":3",
		},
		Strategy: "random",
		Debug:    false,
	}
}
