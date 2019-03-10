# Hadoop 2.7.7 - Part 3

# Login as hduser@master, hduser@slaveX

## Setup Environment Variables for JDK
```sh
vim ~/.bash_aliases
```
`~/.bash_aliases` (append the following lines to the end):
```sh
export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64  # For OpenJDK
#export JAVA_HOME=/usr/lib/jvm/java-8-oracle        # For Oracke JDK
export JRE_HOME=$JAVA_HOME/jre
export PATH=$PATH:$JAVA_HOME/bin:$JRE_HOME/bin
```

## Setup Environment Variables for Hadoop
```sh
vim ~/.bash_aliases
```
`~/.bash_aliases` (append the following lines to the end):
```sh
export HADOOP_HOME=/opt/hadoop-2.7.7
export CLASSPATH=$HADOOP_HOME/lib
export PATH=$PATH:$HADOOP_HOME/sbin:$HADOOP_HOME/bin
```
