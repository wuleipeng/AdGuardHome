package mitmproxy

import (
	"net"
	"net/http"
	"strconv"

	"github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/urlfilter/proxy"
)

// MITMProxy - MITM proxy structure
type MITMProxy struct {
	proxy *proxy.Server
	conf  Config
}

// Config - module configuration
type Config struct {
	ListenAddr string `yaml:"listen_address"`

	// Called when the configuration is changed by HTTP request
	ConfigModified func() `yaml:"-"`

	// Register an HTTP handler
	HTTPRegister func(string, string, func(http.ResponseWriter, *http.Request)) `yaml:"-"`
}

// New - create a new instance of the query log
func New(conf Config) *MITMProxy {
	p := MITMProxy{}
	p.conf = conf
	if !p.create() {
		return nil
	}
	if p.conf.HTTPRegister != nil {
		p.initWeb()
	}
	return &p
}

// Close - close the object
func (p *MITMProxy) Close() {
	p.proxy.Close()
}

// WriteDiskConfig - write configuration on disk
func (p *MITMProxy) WriteDiskConfig(c *Config) {
	*c = p.conf
}

// Start - start proxy server
func (p *MITMProxy) Start() error {
	return p.proxy.Start()
}

// Create a gomitmproxy object
func (p *MITMProxy) create() bool {
	c := proxy.Config{}
	addr, port, err := net.SplitHostPort(p.conf.ListenAddr)
	if err != nil {
		log.Error("net.SplitHostPort: %s", err)
		return false
	}

	c.CompressContentScript = true
	c.ProxyConfig.ListenAddr = &net.TCPAddr{}
	c.ProxyConfig.ListenAddr.IP = net.ParseIP(addr)
	if c.ProxyConfig.ListenAddr.IP == nil {
		log.Error("Invalid IP: %s", addr)
		return false
	}
	c.ProxyConfig.ListenAddr.Port, err = strconv.Atoi(port)
	if c.ProxyConfig.ListenAddr.IP == nil {
		log.Error("Invalid port number: %s", port)
		return false
	}

	p.proxy, err = proxy.NewServer(c)
	if err != nil {
		log.Error("proxy.NewServer: %s", err)
		return false
	}
	return true
}
