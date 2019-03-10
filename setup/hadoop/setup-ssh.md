# Setup Public Key Login for the Hadoop Service Account

# Login as hduser@master

To switch user from admin to hduser:
```sh
sudo -u hduser -i
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
- Note: `slaveX` is just a placeholder for the hostname of a slave.
- Note: Install the public key on each slave.
