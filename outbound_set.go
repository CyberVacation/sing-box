package box

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/outboundset"
	"github.com/sagernet/sing-box/option"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/service"
)

// Bootstrap uses an isolated service registry and only the static dependencies
// of the download client and DNS. It never starts inbounds or changes the
// configured transport to work around a dependency on the set being fetched.
func downloadOutboundSet(ctx context.Context, options option.Options, url string, client option.HTTPClientOptions) ([]byte, error) {
	bootstrap := option.Options{
		Log:                  &option.LogOptions{Disabled: true},
		DNS:                  options.DNS,
		Certificate:          options.Certificate,
		CertificateProviders: options.CertificateProviders,
		NetworkNamespaces:    options.NetworkNamespaces,
	}
	if options.Route != nil {
		route := *options.Route
		route.Rules, route.RuleSet = nil, nil
		route.Final, route.DefaultHTTPClient = "", ""
		bootstrap.Route = &route
	}
	outbounds := make(map[string]option.Outbound)
	endpoints := make(map[string]option.Endpoint)
	for i, value := range options.Outbounds {
		if value.Tag == "" {
			value.Tag = strconv.Itoa(i)
		}
		outbounds[value.Tag] = value
	}
	for i, value := range options.Endpoints {
		if value.Tag == "" {
			value.Tag = strconv.Itoa(i)
		}
		endpoints[value.Tag] = value
	}
	state := make(map[string]uint8)
	var add func(string) error
	add = func(tag string) error {
		if state[tag] == 2 {
			return nil
		}
		if state[tag] == 1 {
			return E.New("circular bootstrap outbound dependency: ", tag)
		}
		state[tag] = 1
		if value, found := outbounds[tag]; found {
			switch group := value.Options.(type) {
			case *option.SelectorOutboundOptions:
				if len(group.OutboundSet) > 0 {
					return E.New("download depends on outbound-set group: ", tag, "; provide cache or initial_path")
				}
			case *option.URLTestOutboundOptions:
				if len(group.OutboundSet) > 0 {
					return E.New("download depends on outbound-set group: ", tag, "; provide cache or initial_path")
				}
			}
			for _, dependency := range outboundset.ReferencedOutbounds(value.Options) {
				if err := add(dependency); err != nil {
					return err
				}
			}
			bootstrap.Outbounds = append(bootstrap.Outbounds, value)
		} else if value, found := endpoints[tag]; found {
			for _, dependency := range outboundset.ReferencedOutbounds(value.Options) {
				if err := add(dependency); err != nil {
					return err
				}
			}
			bootstrap.Endpoints = append(bootstrap.Endpoints, value)
		} else {
			return E.New("bootstrap outbound unavailable: ", tag, "; provide cache or initial_path")
		}
		state[tag] = 2
		return nil
	}
	references := outboundset.ReferencedOutbounds(options.DNS)
	if client.Detour != "" {
		references = append(references, client.Detour)
	}
	for _, reference := range references {
		if err := add(reference); err != nil {
			return nil, err
		}
	}
	instance, err := New(Options{Context: service.ExtendContext(ctx), Options: bootstrap})
	if err != nil {
		return nil, E.Cause(err, "create outbound-set bootstrap")
	}
	defer instance.Close()
	if err = instance.PreStart(); err != nil {
		return nil, E.Cause(err, "start outbound-set bootstrap")
	}
	// Endpoints have an additional Start stage after PreStart.
	if err = instance.endpoint.Start(adapter.StartStateStart); err != nil {
		return nil, err
	}
	clients := service.FromContext[adapter.HTTPClientManager](instance.ctx)
	transport, err := clients.ResolveTransport(instance.ctx, instance.logger, client)
	if err != nil {
		return nil, err
	}
	defer transport.CloseIdleConnections()
	return outboundset.Download(ctx, &http.Client{Transport: transport, Timeout: 30 * time.Second}, url)
}
