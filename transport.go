package azuretls

import (
	"net/url"

	http "github.com/Noooste/fhttp"
	"github.com/Noooste/fhttp/http2"
)

func (s *Session) initTransport(browser string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Transport == nil {
		s.initHTTP1()
	}

	if s.HTTP2Transport == nil {
		if err = s.initHTTP2(browser); err != nil {
			return
		}
	}

	return
}

func (s *Session) initHTTP1() {
	s.Transport = &http.Transport{
		TLSHandshakeTimeout:   s.TimeOut,
		ResponseHeaderTimeout: s.TimeOut,
		Proxy: func(req *http.Request) (*url.URL, error) {
			if s.ProxyDialer == nil {
				return nil, nil
			}
			return s.ProxyDialer.ProxyURL, nil
		},
	}
}

type Http2Config struct {
	Priorities     []http2.Priority
	Settings       map[http2.SettingID]uint32
	SettingsOrder  []http2.SettingID
	ConnectionFlow uint32
	HeaderPriority *http2.PriorityParam
}

func (s *Session) initHTTP2(browser string) error {
	tr, err := http2.ConfigureTransports(s.Transport) // upgrade to HTTP2, while keeping http.Transport

	if err != nil {
		return err
	}

	h2Config := s.H2Config
	if h2Config == nil {
		s, so := defaultHeaderSettings(browser)
		h2Config = &Http2Config{
			Priorities:     defaultStreamPriorities(browser),
			Settings:       s,
			SettingsOrder:  so,
			ConnectionFlow: defaultWindowsUpdate(browser),
			HeaderPriority: defaultHeaderPriorities(browser),
		}
	}

	tr.Priorities = h2Config.Priorities
	tr.Settings, tr.SettingsOrder = h2Config.Settings, h2Config.SettingsOrder
	tr.ConnectionFlow = h2Config.ConnectionFlow
	tr.HeaderPriority = h2Config.HeaderPriority
	tr.StrictMaxConcurrentStreams = true

	tr.PushHandler = &http2.DefaultPushHandler{}

	for k, v := range tr.Settings {
		switch k {
		case http2.SettingInitialWindowSize:
			tr.InitialWindowSize = v

		case http2.SettingHeaderTableSize:
			tr.HeaderTableSize = v
		}
	}

	s.HTTP2Transport = tr

	return nil
}
