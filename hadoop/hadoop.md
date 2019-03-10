# Hadoop 2.7.7 - Part 1

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
