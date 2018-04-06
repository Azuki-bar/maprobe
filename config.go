package maprobe

import (
	"io/ioutil"
	"log"
	"os"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey       string         `yaml:"apikey"`
	ProbesConfig []*ProbeConfig `yaml:"probes"`
	ProbeOnly    bool           `yaml:"probe_only"`
}

type ProbeConfig struct {
	Service string           `yaml:"service"`
	Role    string           `yaml:"role"`
	Roles   []string         `yaml:"roles"`
	Ping    *PingProbeConfig `yaml:"ping"`
	TCP     *TCPProbeConfig  `yaml:"tcp"`
	HTTP    *HTTPProbeConfig `yaml:"http"`
}

func (pc *ProbeConfig) GenerateProbes(host *mackerel.Host) []Probe {
	var probes []Probe

	if pingConfig := pc.Ping; pingConfig != nil {
		p, err := pingConfig.GenerateProbe(host)
		if err != nil {
			log.Printf("[error] cannot generate ping probe. HostID:%s Name:%s %s", host.ID, host.Name, err)
		} else {
			probes = append(probes, p)
		}
	}

	if tcpConfig := pc.TCP; tcpConfig != nil {
		p, err := tcpConfig.GenerateProbe(host)
		if err != nil {
			log.Printf("[error] cannot generate tcp probe. HostID:%s Name:%s %s", host.ID, host.Name, err)
		} else {
			probes = append(probes, p)
		}
	}

	if httpConfig := pc.HTTP; httpConfig != nil {
		p, err := httpConfig.GenerateProbe(host)
		if err != nil {
			log.Printf("[error] cannot generate http probe. HostID:%s Name:%s %s", host.ID, host.Name, err)
		} else {
			probes = append(probes, p)
		}
	}

	return probes
}

func LoadConfig(path string) (*Config, error) {
	c := Config{
		APIKey: os.Getenv("MACKEREL_APIKEY"),
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "load config failed")
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	for _, pc := range c.ProbesConfig {
		if pc.Role != "" {
			pc.Roles = append(pc.Roles, pc.Role)
		}
	}

	return &c, c.validate()
}

func (c *Config) validate() error {
	if c.APIKey == "" {
		return errors.New("no API Key")
	}
	return nil
}
