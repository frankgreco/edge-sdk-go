package types

type DHCPOptions struct {
	DefaultRoute         string `json:"default-route,omitempty"`
	DefaultRouteDistance int    `json:"-"`
	NameServer           string `json:"name-server,omitempty"`
}

type IP struct {
	EnableProxyARP interface{} `json:"enable-proxy-arp,omitempty"` // don't know what this is yet and api respose has it as null
}

type FirewallAttachment struct {
	Interface string  `json:"-" tfsdk:"interface"`
	In        *string `json:"in,omitempty" tfsdk:"in"`
	Out       *string `json:"out,omitempty" tfsdk:"out"`
	Local     *string `json:"local,omitempty" tfsdk:"local"`
}

type Ethernet struct {
	ID          string              `json:"-" tfsdk:"id"`
	Addresses   []string            `json:"address,omitempty"`
	Description string              `json:"description,omitempty"`
	DHCPOptions *DHCPOptions        `json:"dhcp-options,omitempty"`
	Duplex      string              `json:"duplex,omitempty"`
	Speed       string              `json:"speed,omitempty"`
	IP          *IP                 `json:"ip,omitempty"`
	Firewall    *FirewallAttachment `json:"firewall"`
}
