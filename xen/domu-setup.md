# Setup a New VM

# On Dom0

## Create the Image of the VM
```sh
sudo xen-create-image --hostname=vmX
```
- Note: `vmX` is just a placeholder for the hostname of a VM.
- Note: Save the command output

## (Optional) Monitor the Creation of the Image
```sh
sudo tail -f /var/log/xen-tools/vmX.log
```
- Note: `vmX` is just a placeholder for the hostname of a VM.
- Note: Execute this on another terminal

## Configure the Parameters of the VM
```sh
sudo vim /etc/xen/vmX.cfg
```
- Note: `vmX` is just a placeholder for the hostname of a VM.
- Note: Change the MAC address

## Start the VM
```sh
sudo xl create /etc/xen/vmX.cfg -c
```
- Note: `vmX` is just a placeholder for the hostname of a VM.
- Note: Use `root` for username, and the password saved from the previous step

# On DomU

## Change `root` password
```sh
passwd
```

## Create Netplan Config for DomU
```sh
vi /etc/netplan/01-netcfg.yaml
```
`01-netcfg.yaml` [[download](01-netcfg.yaml)]:
```yaml
network:
  version: 2
    renderer: networkd
    ethernets:
      eth0:
        dhcp4: yes
```

## Apply the Netplan for DomU
```sh
netplan generate
netplan apply
```

## Configure DNS
```sh
rm /etc/resolv.conf
vi /etc/resolv.conf
```
`resolv.conf`:
```
nameserver 192.168.1.1
```
- Note: The `rm` step is essential, since the original file is a symlink, and it is generated and managed by Netplan.

## Add Admin Account
```sh
adduser admin
usermod -a -G sudo admin
```

## Protect Home Directory
```sh
chmod 700 /home/admin
```

## Update Packages
```sh
apt-get update
apt-get -y upgrade
```

## Install Essential Utilities
```sh
apt-get install vim net-tools
```

## Reboot the VM
```sh
reboot
```

# On Dom0

## Install the SSH Public Key
```sh
ssh-copy-id -i ~/.ssh/id_ecdsa.pub admin@vmX
```
- Note: `vmX` is just a placeholder for the hostname of a VM.
