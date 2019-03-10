# Xen 4.9 Setup Steps

1. [Install Xen](install-xen.md)
   - Once on Bare Metal
   - Login as admin@bare-metal
   - Turn Bare Metal to Dom0
1. [Setup Dom0](setup-dom0.md)
   - Once on Dom0
   - Login as admin@dom0
1. [Setup a New VM](setup-domu.md)
   - Every time a new DomU is created
   - Login as admin@dom0
   - Then login as root@domu

---

# Extra Setup Steps

## CPU Pinning
```sh
sudo vim /etc/default/grub.d/xen.cfg
sudo reboot
```
`/etc/default/grub.d/xen.cfg` (modify the appropriate lines):
- Before:
  ```
  GRUB_CMDLINE_XEN_DEFAULT=""
  GRUB_CMDLINE_XEN=""
  ```
- After:
  ```
  GRUB_CMDLINE_XEN_DEFAULT="dom0_max_vcpus=2 dom0_vcpus_pin"
  GRUB_CMDLINE_XEN="dom0_max_vcpus=2 dom0_vcpus_pin"
  ```
