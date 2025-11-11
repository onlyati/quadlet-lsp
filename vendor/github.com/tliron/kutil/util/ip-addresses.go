package util

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

// Returns true if the IP address string representation is IPv6.
//
// Note that this function does no validation and assumes the argument is already
// a valid IPv6 or IPv4 address.
func IsIPv6(address string) bool {
	// See: https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6
	return strings.Contains(address, ":")
}

// Returns "[address]:port" for IPv6 and "address:port" for IPv4.
//
// Meant to satisfy the unfortunate requirement of many APIs to provide
// an address (or hostname) and port with a single string argument.
func JoinIPAddressPort(address string, port int) string {
	if IsIPv6(address) {
		return "[" + address + "]:" + strconv.FormatInt(int64(port), 10)
	} else {
		return address + ":" + strconv.FormatInt(int64(port), 10)
	}
}

func SplitIPAddressPort(addressPort string) (string, int, bool) {
	if p := strings.LastIndex(addressPort, ":"); p != -1 {
		address := addressPort[:p]
		port := addressPort[p+1:]
		if strings.HasPrefix(address, "[") {
			// IPv6
			var ok bool
			if address, _, ok = strings.Cut(address[1:], "]"); ok {
				if port_, err := strconv.Atoi(port); err == nil {
					return address, port_, true
				}
			}
		} else {
			// IPv4
			if port_, err := strconv.Atoi(port); err == nil {
				return address, port_, true
			}
		}
	}

	return "", 0, false
}

// If the zone is not empty returns "address%zone".
// It is expected that the argument does not already have a zone.
//
// For IPv6 address string representations only, see:
// https://en.wikipedia.org/wiki/IPv6_address#Scoped_literal_IPv6_addresses_(with_zone_index)
func JoinIPAddressZone(address string, zone string) string {
	if zone != "" {
		return address + "%" + zone
	} else {
		return address
	}
}

// Returns true if the two UDP addresses are equal.
func IsUDPAddrEqual(a *net.UDPAddr, b *net.UDPAddr) bool {
	return a.IP.Equal(b.IP) && (a.Port == b.Port) && (a.Zone == b.Zone)
}

// Always returns a specified address. If the argument is already a specified address,
// returns it as is. Otherwise (when it's "::" or "0.0.0.0") will attempt to find a specified
// address by enumerating the active local interfaces, chosing one arbitrarily, with a
// preference for a global unicast address.
//
// The IP version of the returned address will match that of the argument, IPv6 for "::"
// and IPv4 for "0.0.0.0".
//
// Note that a returned IPv6 address may include a zone (when not a global unicast).
func ToReachableIPAddress(address string) (string, error) {
	// (Note: net.ParseIP can't parse IPv6 with zone, but netip.ParseAddr can)
	if addr, err := netip.ParseAddr(address); err == nil {
		if addr.IsUnspecified() {
			// Try to find a global unicast first

			collector := IPAddressCollector{
				IPv6:     IsIPv6(address),
				WithZone: true,
				FilterInterface: func(interface_ net.Interface) bool {
					return (interface_.Flags&net.FlagLoopback == 0) && (interface_.Flags&net.FlagUp != 0)
				},
				FilterIP: func(ip net.IP) bool {
					return ip.IsGlobalUnicast()
				},
			}

			if addresses, err := collector.Collect(); err == nil {
				if len(addresses) > 0 {
					return addresses[0], nil
				}
			} else {
				return "", err
			}

			// Otherwise, just use the first address
			collector.FilterIP = nil
			if addresses, err := collector.Collect(); err == nil {
				if len(addresses) > 0 {
					return addresses[0], nil
				}
			} else {
				return "", err
			}

			return "", fmt.Errorf("cannot find an equivalent reachable address for: %s", address)
		}
	} else {
		return "", err
	}

	return address, nil
}

// The argument is validated as being a multicast address, e.g. "ff02::1" (IPv6) or
// "239.0.0.1" (IPv4). For IPv6, if it does not include a zone, a valid zone will be
// added by enumerating the active local interfaces, chosing one arbitrarily.
func ToBroadcastIPAddress(address string) (string, error) {
	// (Note: net.ParseIP can't parse IPv6 with zone, but netip.ParseAddr can)
	if addr, err := netip.ParseAddr(address); err == nil {
		if !addr.IsMulticast() {
			return "", fmt.Errorf("not a multicast address: %s", address)
		}

		if IsIPv6(address) && (addr.Zone() == "") {
			if interfaces, err := net.Interfaces(); err == nil {
				for _, interface_ := range interfaces {
					if (interface_.Flags&net.FlagLoopback == 0) && (interface_.Flags&net.FlagUp != 0) &&
						(interface_.Flags&net.FlagBroadcast != 0) && (interface_.Flags&net.FlagMulticast != 0) {
						// The IPv6 zone is usually the interface name
						return JoinIPAddressZone(address, interface_.Name), nil
					}
				}
			} else {
				return "", err
			}

			return "", fmt.Errorf("cannot find IPv6 zone for: %s", address)
		}

		return address, nil
	} else {
		return "", err
	}
}

func IPAddressPortWithoutZone(address string) string {
	if strings.Contains(address, "%") {
		var port string
		if colon := strings.LastIndex(address, ":"); colon != -1 {
			port = address[colon+1:]
			address = address[:colon]
		}

		var ipv6 bool
		if strings.HasPrefix(address, "[") {
			// This should always be the case
			ipv6 = true
			address = address[1 : len(address)-1]
		}

		address, _, _ = strings.Cut(address, "%")

		if ipv6 {
			address = "[" + address + "]"
		}

		if port == "" {
			return address
		} else {
			return address + ":" + port
		}
	} else {
		return address
	}
}

func DumpIPAddress(address any) {
	// Note: net.ParseIP can't parse IPv6 zone
	ip := netip.MustParseAddr(ToString(address))
	fmt.Printf("address: %s\n", ip)
	fmt.Printf("  global unicast:            %t\n", ip.IsGlobalUnicast())
	fmt.Printf("  interface local multicast: %t\n", ip.IsInterfaceLocalMulticast())
	fmt.Printf("  link local multicast:      %t\n", ip.IsLinkLocalMulticast())
	fmt.Printf("  link local unicast:        %t\n", ip.IsLinkLocalUnicast())
	fmt.Printf("  loopback:                  %t\n", ip.IsLoopback())
	fmt.Printf("  multicast:                 %t\n", ip.IsMulticast())
	fmt.Printf("  private:                   %t\n", ip.IsPrivate())
	fmt.Printf("  unspecified:               %t\n", ip.IsUnspecified())
}

//
// IPAddressCollector
//

type FilterInterfaceFunc func(interface_ net.Interface) bool

type FilterIPFunc func(ip net.IP) bool

type IPAddressCollector struct {
	// If nil, will call net.Interfaces().
	Interfaces []net.Interface

	// Which IP version to accept.
	IPv6 bool

	// Include IPv6 zone in returned addresses.
	WithZone bool

	// Return true to accept an interface (can be nil).
	FilterInterface FilterInterfaceFunc

	// Return true to accept an IP (can be nil).
	// Note that the argument's address (ip.String()) does not include the IPv6 zone.
	FilterIP FilterIPFunc
}

func (self *IPAddressCollector) Collect() ([]string, error) {
	var addresses []string

	if interfaces, err := self.interfaces(); err == nil {
		for _, interface_ := range interfaces {
			if (self.FilterInterface == nil) || self.FilterInterface(interface_) {
				if addrs, err := interface_.Addrs(); err == nil {
					for _, addr := range addrs {
						if ipNet, ok := addr.(*net.IPNet); ok {
							ip := ipNet.IP
							address := ip.String()
							isIpV6 := IsIPv6(address)
							if (isIpV6 == self.IPv6) && ((self.FilterIP == nil) || self.FilterIP(ip)) {
								if self.WithZone && isIpV6 {
									// The IPv6 zone is usually the interface name
									address = JoinIPAddressZone(address, interface_.Name)
								}
								addresses = append(addresses, address)
							}
						}
					}
				} else {
					return nil, err
				}
			}
		}
	} else {
		return nil, err
	}

	return addresses, nil
}

func (self *IPAddressCollector) interfaces() ([]net.Interface, error) {
	if self.Interfaces == nil {
		var err error
		if self.Interfaces, err = net.Interfaces(); err != nil {
			return nil, err
		}
	}

	return self.Interfaces, nil
}
