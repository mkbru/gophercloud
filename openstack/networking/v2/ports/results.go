package ports

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// IP is a sub-struct that represents an individual IP.
type IP struct {
	SubnetID  string `mapstructure:"subnet_id" json:"subnet_id"`
	IPAddress string `mapstructure:"ip_address" json:"ip_address,omitempty"`
}

// Port represents a Neutron port. See package documentation for a top-level
// description of what this is.
type Port struct {
	// UUID for the port.
	ID string `mapstructure:"id" json:"id"`
	// Network that this port is associated with.
	NetworkID string `mapstructure:"network_id" json:"network_id"`
	// Human-readable name for the port. Might not be unique.
	Name string `mapstructure:"name" json:"name"`
	// Administrative state of port. If false (down), port does not forward packets.
	AdminStateUp bool `mapstructure:"admin_state_up" json:"admin_state_up"`
	// Indicates whether network is currently operational. Possible values include
	// `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.
	Status string `mapstructure:"status" json:"status"`
	// Mac address to use on this port.
	MACAddress string `mapstructure:"mac_address" json:"mac_address"`
	// Specifies IP addresses for the port thus associating the port itself with
	// the subnets where the IP addresses are picked from
	FixedIPs []IP `mapstructure:"fixed_ips" json:"fixed_ips"`
	// Owner of network. Only admin users can specify a tenant_id other than its own.
	TenantID string `mapstructure:"tenant_id" json:"tenant_id"`
	// Identifies the entity (e.g.: dhcp agent) using this port.
	DeviceOwner string `mapstructure:"device_owner" json:"device_owner"`
	// Specifies the IDs of any security groups associated with a port.
	SecurityGroups []string `mapstructure:"security_groups" json:"security_groups"`
	// Identifies the device (e.g., virtual server) using this port.
	DeviceID string `mapstructure:"device_id" json:"device_id"`
}

// PortPage is the page returned by a pager when traversing over a collection
// of network ports.
type PortPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of ports has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (p PortPage) NextPageURL() (string, error) {
	type resp struct {
		Links []struct {
			Href string `mapstructure:"href"`
			Rel  string `mapstructure:"rel"`
		} `mapstructure:"ports_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

// IsEmpty checks whether a PortPage struct is empty.
func (p PortPage) IsEmpty() (bool, error) {
	is, err := ExtractPorts(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractPorts accepts a Page struct, specifically a PortPage struct,
// and extracts the elements into a slice of Port structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPorts(page pagination.Page) ([]Port, error) {
	var resp struct {
		Ports []Port `mapstructure:"ports" json:"ports"`
	}

	err := mapstructure.Decode(page.(PortPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Ports, nil
}
