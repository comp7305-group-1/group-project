# Preparations

# Login as admin@master, and admin@slaveX

Note: The following steps have to be done once on master, and once on each slave.

## Check CPU, RAM, Swap
```sh
free -h
```

## Check and Update Local DNS Table
```sh
sudo vim /etc/hosts
```
`/etc/hosts`:
```sh
127.0.0.1 localhost
192.168.1.2 master
192.168.1.3 slaveX
```
- Note: `192.168.1.2` and `192.168.1.3` are just placeholders for the IP of the master and a slave.
- Note: `master` and `slaveX` are just placeholders for the hostname of the master and a slave.
- Note: Add an entry for each slave.
