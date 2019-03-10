# Install Hadoop

# Login as hduser@master

To switch user from admin to hduser:
```sh
sudo -u hduser -i
```

## Download and Install Hadoop 2.7.7
```sh
cd /opt
wget http://apache.01link.hk/hadoop/common/hadoop-2.7.7/hadoop-2.7.7.tar.gz
tar zxvf hadoop-2.7.7.tar.gz
```

## Setup Hadoop Parameters
```sh
cd /opt/hadoop-2.7.7etc/hadoop
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
slaveX
```
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Add an entry for each slave.

## Setup Environment Variables for JDK and Hadoop
```sh
vim ~/.bash_aliases
```
`~/.bash_aliases` (append the following lines to the end):
```sh
export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64  # For OpenJDK
#export JAVA_HOME=/usr/lib/jvm/java-8-oracle        # For Oracke JDK
export JRE_HOME=$JAVA_HOME/jre
export PATH=$PATH:$JAVA_HOME/bin:$JRE_HOME/bin

export HADOOP_HOME=/opt/hadoop-2.7.7
export CLASSPATH=$HADOOP_HOME/lib
export PATH=$PATH:$HADOOP_HOME/sbin:$HADOOP_HOME/bin
```

## Copy the Hadoop Installation to Slaves
```sh
scp /opt/hadoop-2.7.7 slaveX:/opt/
```
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Copy the Hadoop installation to each slave.

## Copy the Environment Variable Settings to Slaves
```sh
scp ~/.bash_aliases slaveX:
```
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Copy the environment variable settings to each slave.
