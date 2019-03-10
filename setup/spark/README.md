# Spark 2.4.0 Setup Steps

# Login as hduser@master

To switch user from admin to hduser:
```sh
sudo -u hduser -i
```

## Download and Install Spark 2.4.0
```sh
cd /opt
wget http://apache.01link.hk/spark/spark-2.4.0/spark-2.4.0-bin-hadoop2.7.tgz
tar zxvf spark-2.4.0-bin-hadoop2.7.tgz
```

## Setup Spark Parameters
```sh
cd /opt/spark-2.4.0-bin-hadoop2.7/conf
cp -a spark-env.sh.template spark-env.sh
cp -a spark-defaults.conf.template spark-defaults.conf
vim spark-env.sh
vim spark-defaults.conf
```
`/opt/spark-2.4.0-bin-hadoop2.7/conf/spark-env.sh`:
```sh
HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
LD_LIBRARY_PATH=/opt/hadoop-2.7.7/lib/native:$LD_LIBRARY_PATH
```
`/opt/spark-2.4.0-bin-hadoop2.7/conf/spark-defaults.conf`:
```
# (pending)
```

## Setup Environment Variables for Spark
```sh
vim ~/.bash_aliases
```
`~/.bash_aliases` (append the following lines to the end):
```sh
export SPARK_HOME=/opt/spark-2.4.0-bin-hadoop2.7
export PATH=$PATH:$PSARK_HOME/sbin:$SPARK_HOME/bin
```

## Copy the Spark Installation to Slaves
```sh
scp /opt/spark-2.4.0-bin-hadoop2.7 slaveX:/opt/
```
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Copy the Spark installation to each slave.

## Copy the Environment Variable Settings to Slaves
```sh
scp ~/.bash_aliases slaveX:
```
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Copy the environment variable settings to each slave.

## Create HDFS folder for event Log
```sh
hdfs dfs -mkdir /tmp/sparkLog
```
