# bird-bgp-unnumbered-helper
Helper tool to assist in generating BGP unnumbered configurations for BIRD

## Usage

```sh
$ export IFACE=eth0
$ sudo bash -c "./bird-bgp-unnumbered-helper $IFACE > /etc/bird/neighbours.conf"
$ cat /etc/bird/neighbours.conf
define eth0_neighbor = fe80::xxxx:xxxx:xxxx:xxxx;

$ sudo vim /etc/bird/bird.conf
include "neighbours.conf";

protocol bgp {
  local as <local_as>;
  neighbor eth0_neighbor%eth0 as <remote_as>;
  ...
}

$ sudo birdc configure
```
