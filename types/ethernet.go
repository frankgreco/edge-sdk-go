package types

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type DHCPOptions struct {
	DefaultRoute         string `json:"default-route,omitempty"`
	DefaultRouteDistance int    `json:"-"`
	NameServer           string `json:"name-server,omitempty"`
}

type IP struct {
	EnableProxyARP interface{} `json:"enable-proxy-arp,omitempty"` // don't know what this is yet and api respose has it as null
}

type FirewallAttachment struct {
	ID        tftypes.String `json:"-" tfsdk:"id"`
	Interface string         `json:"-" tfsdk:"interface"`
	In        *string        `json:"in,omitempty" tfsdk:"in"`
	Out       *string        `json:"out,omitempty" tfsdk:"out"`
	Local     *string        `json:"local,omitempty" tfsdk:"local"`
}

type VirtualInterface struct {
	Firewall *FirewallAttachment `json:"firewall,omitempty" tfsdk:"-"`
}

type Ethernet struct {
	ID          string                       `json:"-" tfsdk:"id"`
	Addresses   []string                     `json:"address,omitempty" tfsdk:"-"`
	Description string                       `json:"description,omitempty" tfsdk:"-"`
	DHCPOptions *DHCPOptions                 `json:"dhcp-options,omitempty" tfsdk:"-"`
	Duplex      string                       `json:"duplex,omitempty" tfsdk:"-"`
	Speed       string                       `json:"speed,omitempty" tfsdk:"-"`
	IP          *IP                          `json:"ip,omitempty" tfsdk:"-"`
	Firewall    *FirewallAttachment          `json:"firewall,omitempty" tfsdk:"-"`
	Vif         map[string]*VirtualInterface `json:"vif,omitempty" tfsdk:"-"`
}

func (a *FirewallAttachment) GetID() string {
	return a.Interface
}
