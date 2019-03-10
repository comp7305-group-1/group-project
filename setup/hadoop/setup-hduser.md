# Setup the Hadoop Service Account

# Login as admin@master, and admin@slaveX
Note: The following steps have to be done once on master, and once on each slave.

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
