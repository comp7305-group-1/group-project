# Initial Setup (on Dom0)

## Update Packages
```sh
sudo apt-get update
sudo apt-get -y upgrade
```

## Install Essential Utilities
```sh
sudo apt-get install vim net-tools
```

## Install Xen
```sh
sudo apt-get install xen-hypervisor-4.9-amd64
```
- Note: `xen-hypervisor` and `xen-hypervisor-amd64` are virtual packages.

## Install Bridge Utilities
```sh
sudo apt-get install bridge-utils
```

## Create Netplan Config for Dom0
```sh
sudo vim /etc/netplan/01-netcfg.yaml
```
`01-netcfg.yaml` [[download](01-netcfg-dom0-template.yaml)]:
```yaml
network:
  version: 2
    renderer: networkd
    ethernets:
      xxx:
        dhcp4: no
    bridges:
      xenbr0:
        macaddress: xx:xx:xx:xx:xx:xx
          interfaces: [xxx]
          dhcp4: yes
```
- Note: Change `xxx` to the name of the physical interface (e.g. `enp0s10`).
- Note: Change `xx:xx:xx:xx:xx:xx` to the MAC address of the physical interface.

## Apply the Netplan for Dom0
```sh
sudo netplan generate
sudo netplan apply
```

## Reboot the Machine
```sh
sudo reboot
```

## Install Xen Utilities
```sh
sudo apt-get install xen-tools xen-utils xenwatch virt-manager
```

## Edit Xen Daemon Config
```sh
sudo vim /etc/xen/xend-config.sxp
```
Uncomment this line:
```
(xend-unix-server yes)
```

## Create Home for Xen
```sh
sudo mkdir /home/xen
sudo chmod 777 /home/xen
```

## Create Symlink for Xen Library
```sh
cd /usr/lib
sudo ln -s xen-4.9 xen
```

## Generate SSH Key Pair
```sh
ssh-keygen -t ecdsda -b 521 -f ~/.ssh/id_ecdsa -N ''
```
