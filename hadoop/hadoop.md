# Hadoop on Dom0

# Login as admin@dom0

## Disable IPv6 on Dom0
```sh
sudo vim /etc/default/grub
sudo update-grub
sudo reboot
```
`/etc/default/grub` (modify the following lines):
- Before:
  ```
  GRUB_CMDLINE_LINUX_DEFAULT="splash quiet"
  GRUB_CMDLINE_LINUX=""
  ```
- After:
  ```
  GRUB_CMDLINE_LINUX_DEFAULT="splash quiet ipv6.disable=1"
  GRUB_CMDLINE_LINUX="ipv6.disable=1"
  ```

## Disable IPv6 on DomU
```sh
sudo xl destroy slaveX
sudo vim /etc/xen/slaveX.cfg
sudo xl create /etc/xen/slaveX.cfg
```
`/etc/xen/slaveX.cfg` (append the following line to the end):
```
extra = 'ipv6.disable=1'
```
- Note: `slaveX` is a placeholder for the hostname of a slave.
- Note: Repeat this step for each slave.



# Login as admin@master, and admin@slaveX

Note: The following steps have to be done once on master, and once on each slave.

## Check CPU, RAM, Swap
```sh
free -h
```

## Check and update /etc/hosts
```sh
127.0.0.1 localhost
4.3.2.1 master
4.3.2.2 slaveX
```
- Note: The IPs are just example IPs.
- Note: `slaveX` is a placeholder for the hostname of a slave.
- Note: Add an entry for each slave.

## Update Packages
```sh
sudo apt-get update
sudo apt-get -y upgrade
```

## Install JDK (Choose Either One)

### OpenJDK
```sh
sudo apt-get -y install openjdk-8-jdk-headless
```

### Oracle JDK
```sh
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get update
sudo apt-get -y install oracle-java8-installer
```

## Setup the Hadoop Service Account
```sh
sudo addgroup hadoop
sudo adduser --ingroup hadoop hduser
```

## Setup Directories for Hadoop
```sh
sudo chmod 777 /opt
sudo mkdir /var/hadoop
sudo chown hduser:hadoop /var/hadoop
```