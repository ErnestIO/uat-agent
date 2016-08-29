# Build Ernest on vCloud

**add ernest.local to /etc/hosts**

```
apt-get update
apt-get upgrade
apt-get install -y wget git make ruby
wget https://opscode-omnibus-packages.s3.amazonaws.com/ubuntu/12.04/x86_64/chefdk_0.10.0-1_amd64.deb
dpkg -i chefdk_0.10.0-1_amd64.deb
git clone https://github.com/ernestio/ernest-vagrant.git
cd ernest-vagrant
git checkout master
nohup make deploy &
```
