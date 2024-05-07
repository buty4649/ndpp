# ndpp

A CLI tool for automating tasks based on Neighbor Discovery Protocol events

## Usage

```sh
$ export IFACE=eth0
$ sudo ./ndpp $IFACE
- local:
    addr: fe80::xxxx:xxxx:xxxx:xxxx
    interface: eth0
    lladdr: xx:xx:xx:xx:xx:xx
  router:
    addr: fe80::yyyy:yyyy:yyyy:yyyy
    addr_with_zone: fe80::yyyy:yyyy:yyyy:yyyy%eth0
    lladdr: yy:yy:yy:yy:yy:yy
```
### Generate configuration file for BGP Unnumbered in BIRD

ndpp can generate configuration files for BGP Unnumbered in BIRD. The method to achieve this is simple: just specify -T bird. Additionally, using the -o option allows you to output to a specified file path.

```sh
$ export IFACE=eth0
$ sudo ./ndpp -T bird $IFACE
define eth0_neighbor_address = fe80::xxxx:xxxx:xxxx:xxxx;
define eth0_local_address = fe80::yyyy:yyyy:yyyy:yyyy;
define eth0_lladdr = hex:yy:yy:yy:yy:yy:yy;

#
# The following is a sample configuration. Please uncomment to use.
#
#define remote_asn = <Please set remote ASN>;
#define local_asn = <Please set local ASN>;
#
#protocol radv {
#  interface "eth0" {
#    custom option type 1 value eth0_lladdr;
#  };
#}
#protocol bgp eth0 {
#  neighbor eth0_neighbor_address}} as remote_asn;
#  interface "eth0";
#  local as local_asn;
#  ipv4 { extended next hop; import all; export all; };
#  ipv6 { import all; export all; };
#}
```
