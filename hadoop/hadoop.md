# Disable IPv6

# Login as admin@dom0

## Disable IPv6 on Dom0 / Bare Metal
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

## Disable IPv6 on Dom0 for DomU
```sh
sudo xl destroy vmX
sudo vim /etc/xen/vmX.cfg
sudo xl create /etc/xen/vmX.cfg
```
`/etc/xen/vmX.cfg` (append the following line to the end):
```
extra = 'ipv6.disable=1'
```
- Note: `vmX` is a placeholder for the hostname of a VM.
- Note: Repeat this step for each VM.
