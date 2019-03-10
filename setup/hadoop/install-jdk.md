# Install JDK

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
4.3.2.1 master
4.3.2.2 slaveX
```
- Note: The IPs are just example IPs.
- Note: `slaveX` is a just placeholder for the hostname of a slave.
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

## Protect Home Directory
```sh
chmod 700 /home/hduser
```

## Setup Directories for Hadoop
```sh
sudo chmod 777 /opt
sudo mkdir /var/hadoop
sudo chown hduser:hadoop /var/hadoop
```
