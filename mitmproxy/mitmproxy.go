package mitmproxy

import (
	"net"
	"net/http"
	"strconv"

	"github.com/AdguardTeam/golibs/log"
	"github.com/AdguardTeam/gomitmproxy"
)

// MITMProxy - MITM proxy structure
type MITMProxy struct {
	proxy *gomitmproxy.Proxy
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
	c := gomitmproxy.Config{}
	addr, port, err := net.SplitHostPort(p.conf.ListenAddr)
	if err != nil {
		log.Error("net.SplitHostPort: %s", err)
		return false
	}

	c.ListenAddr = &net.TCPAddr{}
	c.ListenAddr.IP = net.ParseIP(addr)
	if c.ListenAddr.IP == nil {
		log.Error("Invalid IP: %s", addr)
		return false
	}
	c.ListenAddr.Port, err = strconv.Atoi(port)
	if c.ListenAddr.IP == nil {
		log.Error("Invalid port number: %s", port)
		return false
	}

	p.proxy = gomitmproxy.NewProxy(c)
	return true
}
