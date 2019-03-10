# Hadoop 2.7.7 Setup Steps

1. [Disable IPv6](disable-ipv6.md)
   - Once on Dom0 or Bare Metal
   - Once on Dom0 for each DomU
   - Login as admin@dom0
1. [Install JDK](install-jdk.md)
   - Once on Master
   - Once on each Slave
   - Login as admin@master on Master
   - Login as admin@slaveX on each Slave
1. [Install Hadoop](install-hadoop.md)
    - Once on Master
    - Login as hduser@master
    - Copy Installation and Settings to each Slave
1. [Format the Namenode](format-namenode.md)
    - Once on Master
    - Login as hduser@master on Master
