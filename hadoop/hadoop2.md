# Hadoop 2.7.7 - Part 2

# Login as hduser@master

To switch user from admin to hduser:
```sh
sudo -u hduser -i
```

## Protect Home Directory
```sh
chmod 700 ~
```

## Download and Extract Hadoop 2.7
```sh
cd /opt
wget http://apache.01link.hk/hadoop/common/hadoop-2.7.7/hadoop-2.7.7.tar.gz
tar zxvf hadoop-2.7.7.tar.gz
```

## Edit Haddop Environment
```sh
cd etc/hadoop
vim hadoop-env.sh
vim core-site.xml
vim hdfs-site.xml
vim mapred-site.xml
vim yarn-site.xml
vim masters
vim slaves
```
`/opt/hadoop-2.7.7/etc/hadoop/hadoop-env.sh`:
```sh
JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64  # For OpenJDK
#JAVA_HOME=/usr/lib/jvm/java-8-oracle        # For Oracke JDK
HADOOP_HOME=/opt/hadoop-2.7.7
HADOOP_CONF_DIR=/opt/hadoop-2.7.7/etc/hadoop
```
`/opt/hadoop-2.7.7/etc/hadoop/core-site.xml`:
```xml
<!-- pending -->
```
`/opt/hadoop-2.7.7/etc/hadoop/hdfs-site.xml`:
```xml
<!-- pending -->
```
`/opt/hadoop-2.7.7/etc/hadoop/mapred-site.xml`:
```xml
<!-- pending -->
```
`/opt/hadoop-2.7.7/etc/hadoop/yarn-site.xml`:
```xml
<!-- pending -->
```
`/opt/hadoop-2.7.7/etc/hadoop/masters`:
```
master
```
`/opt/hadoop-2.7.7/etc/hadoop/slaves`:
```
slave1
slave2
```

## Setup Environment Variables for JDK
```sh
export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64  # For OpenJDK
#export JAVA_HOME=/usr/lib/jvm/java-8-oracle        # For Oracke JDK
export JRE_HOME=$JAVA_HOME/jre
export PATH=$PATH:$JAVA_HOME/bin:$JRE_HOME/bin
```

## Setup environment variables for Hadoop
```sh
export HADOOP_HOME=/opt/hadoop-2.7.7
export CLASSPATH=$HADOOP_HOME/lib
export PATH=$PATH:$HADOOP_HOME/sbin:$HADOOP_HOME/bin
```

## Generate SSH Key Pair
```sh
ssh-keygen -t ecdsda -b 521 -f ~/.ssh/id_ecdsa -N ''
```

## Install the SSH Public Key
```sh
ssh-copy-id -i id_ecdsa.pub hduser@master
ssh-copy-id -i id_ecdsa.pub hduser@slaveX
```

## Copy the base directory to slaves
(pending)

## Format the namenode:
```sh
hdfs namenode -format
```
