package model

type Subscription struct {
	Port               int                     `yaml:"port,omitempty"`
	SocksPort          int                     `yaml:"socks-port,omitempty"`
	AllowLan           bool                    `yaml:"allow-lan"`
	Mode               string                  `yaml:"mode,omitempty"`
	LogLevel           string                  `yaml:"logger-level,omitempty"`
	ExternalController string                  `yaml:"external-controller,omitempty"`
	Proxies            []Proxy                 `yaml:"proxies,omitempty"`
	ProxyGroups        []ProxyGroup            `yaml:"proxy-groups,omitempty"`
	Rules              []string                `yaml:"rules,omitempty"`
	RuleProviders      map[string]RuleProvider `yaml:"rule-providers,omitempty,omitempty"`
}

type ProxyGroup struct {
	Name          string   `yaml:"name,omitempty"`
	Type          string   `yaml:"type,omitempty"`
	Proxies       []string `yaml:"proxies,omitempty"`
	IsCountryGrop bool     `yaml:"-"`
	Url           string   `yaml:"url,omitempty"`
	Interval      int      `yaml:"interval,omitempty"`
	Tolerance     int      `yaml:"tolerance,omitempty"`
	Lazy          bool     `yaml:"lazy"`
}

type RuleProvider struct {
	Type     string `yaml:"type,omitempty"`
	Behavior string `yaml:"behavior,omitempty"`
	Url      string `yaml:"url,omitempty"`
	Path     string `yaml:"path,omitempty"`
	Interval int    `yaml:"interval,omitempty"`
}

type Payload struct {
	Rules []string `yaml:"payload,omitempty"`
}
