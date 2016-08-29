# Ernest Testing

## vCloud Provider

### Everything Optional

**prep**

* `ernest service apply novse11.yml`
 * builds the networks we will use to test instX.yml

**inst1.yml**

* initial service apply
* in datacenter r3-dc2
* creates instance r3-dc2-r3test3-stg-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-db
 * has IP 10.2.0.90

**inst2.yml**

* modifies service from inst1.yml
* creates instance r3-dc2-r3test3-stg-2:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-db
 * has IP 10.2.0.91

**inst3.yml**

* modifies service from inst2.yml
* creates instance r3-dc2-r3test3-dev-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-web
 * has IP 10.1.0.90

**inst4.yml**

* modifies service from inst3.yml
* deletes instance r3-dc2-r3test3-stg-2

**inst5.yml**

* modifies service from inst4.yml
* deletes instance r3-dc2-r3test3-stg-1

**cleanup**

* `ernest service destroy r3test3`
* `ernest service destroy r3vse1`

### No VSE Creator

**novse1.yml**

* initial service apply
* in datacenter r3-dc2
* creates network r3-dc2-r3vse1-web
 * connected to router vse2
 * with network 10.1.0.0/24
* creates instance r3-dc2-r3vse1-web-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-web
 * has IP 10.1.0.11
* configures ACLs on router vse2
 * allow internal:any to internal:any
 * allow 172.18.143.3:any to internal:22
 * allow 172.17.240.0/24:any to internal:22
 * allow 172.19.186.30/24:any to internal:22
* configures NATs on router vse2
 * from 10.1.0.0/24:any to 172.16.186.44:any
 * from 172.16.186.44:22 to 10.1.0.11:22

**novse2.yml**

* modifies service from novse1.yml
* adds ACLs on router vse2
 * allow internal:any to external:any

**novse3.yml**

* modifies service from novse2.yml
* adds NATs on router vse2
 * from 172.16.186.61:22 to 10.1.0.12:22

**novse4.yml**

* modifies service from novse3.yml
* creates instance r3-dc2-r3vse1-web-2:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-web
 * has IP 10.1.0.12

**novse5.yml**

* modifies service from novse4.yml
* modifies instances r3-dc2-r3vse1-web-1 and r3-dc2-r3vse1-web-2:
 * change from 1CPU to 2CPU
 * instances are modified one at a time

**novse6.yml**

* modifies service from novse5.yml
* modifies instances r3-dc2-r3vse1-web-1 and r3-dc2-r3vse1-web-2:
 * adds 10GB of disk
 * instances are modified one at a time

**novse7.yml**

* modifies service from novse6.yml
* modifies instances r3-dc2-r3vse1-web-1 and r3-dc2-r3vse1-web-2:
 * changes from 1GB RAM to 2GB RAM
 * instances are modified one at a time

**novse8.yml**

* modifies service from novse7.yml
* creates network r3-dc2-r3vse1-db
 * connected to router vse2
 * with network 10.2.0.0/24

**novse9.yml**

* modifies service from novse8.yml
* creates instance r3-dc2-r3vse1-db-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse1-db
 * has IP 10.2.0.11

**novse10.yml**

* modifies service from novse9.yml
* deletes instance r3-dc2-r3vse1-web-2:

**novse11.yml**

* modifies service from novse10.yml
* deletes instance r3-dc2-r3vse1-db-1:

**cleanup**

* `ernest service destroy r3vse1`
* needed so we can re-use router vse2 for the next batch of tests

**novse12.yml**

* initial service apply
* in datacenter r3-dc2
* creates network r3-dc2-r3vse2-web
 * connected to router vse2
 * with network 10.1.0.0/24
* creates network r3-dc2-r3vse2-salt
 * connected to router vse2
 * with network 10.254.254.0/24
* creates instance r3-dc2-r3vse2-web-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse2-web
 * has IP 10.1.0.11
 * is bootstrapped
 * `date` command has run
* creates instance r3-dc2-r3vse2-salt-master:
 * from template r3/r3-salt-master
 * it is running
 * has 1CPU
 * has 2GB RAM
 * is on network r3-dc2-r3vse2-salt
 * has IP 10.254.254.100
* configures ACLs on router vse2
 * allow internal:any to internal:any
 * allow internal:any to external:any
 * allow 10.254.254.0/24:any to any:22
 * allow 10.254.254.0/24:any to any:5985
 * allow 172.17.241.95:any to 172.16.186.44:8000
 * allow 172.17.241.95:any to 172.16.186.44:22
 * allow 10.1.0.0/24:any to 10.254.254.100:4505
 * allow 10.1.0.0/24:any to 10.254.254.100:4506
* configures NATs on router vse2
 * from 10.1.0.0/24:any to 172.16.186.44:any
 * from 10.254.254.0/24:any to 172.16.186.44:any
 * from 172.16.186.44:22 to 10.254.254.100:22
 * from 172.16.186.44:8000 to 10.254.254.100:8000

**novse13.yml**

* modifies service from novse12.yml
* creates instance r3-dc2-r3vse2-web-2:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse2-web
 * has IP 10.1.0.12
 * is bootstrapped
 * `date` command has run

**novse14.yml**

* modifies service from novse13.yml
* modifies instances r3-dc2-r3vse2-web-1 and r3-dc2-r3vse2-web-2:
 * `date` and `uptime` commands have run

**novse15.yml**

* modifies service from novse14.yml
* creates instance r3-dc2-r3vse2-db-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3vse2-web
 * has IP 10.1.0.21
 * is bootstrapped
 * `date` command has run

**novse16.yml**

* modifies service from novse15.yml
* modifies instance r3-dc2-r3vse2-web-2:
 * has been unregistered with salt master

**cleanup**

* `ernest service destroy r3vse2`

### VSE Creator

**vse1.yml**

* initial service apply
* in datacenter r3-dc2
* creates router vse4
 * connects external network r3-net-ext
 * assigns external IP for NAT (NAT-IP)
* creates network r3-dc2-r3test4-web
 * connected to router vse4
 * with network 10.1.0.0/24
* creates instance r3-dc2-r3test4-web-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test4-web
 * has IP 10.1.0.11
* configures ACLs on router vse4
 * allow internal:any to internal:any
 * allow 172.18.143.3:any to internal:22
 * allow 172.17.240.0/24:any to internal:22
 * allow 172.19.186.30/24:any to internal:22
* configures NATs on router vse4
 * from 10.1.0.0/24:any to NAT-IP:any
 * from NAT-IP:22 to 10.1.0.11:22

**vse2.yml**

* modifies service from vse1.yml
* adds ACLs on router vse4
 * allow internal:any to external:any

**vse3.yml**

* modifies service from vse2.yml
* adds NATs on router vse4
 * from NAT-IP:23 to 10.1.0.12:23

**vse4.yml**

* modifies service from vse3.yml
* creates instance r3-dc2-r3test4-web-2:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test4-web
 * has IP 10.1.0.12

**vse5.yml**

* modifies service from vse4.yml
* modifies instances r3-dc2-r3test4-web-1 and r3-dc2-r3test4-web-2:
 * change from 1CPU to 2CPU
 * instances are modified one at a time

**vse6.yml**

* modifies service from vse5.yml
* modifies instances r3-dc2-r3test4-web-1 and r3-dc2-r3test4-web-2:
 * adds 10GB of disk
 * instances are modified one at a time

**vse7.yml**

* modifies service from vse6.yml
* modifies instances r3-dc2-r3test4-web-1 and r3-dc2-r3test4-web-2:
 * changes from 1GB RAM to 2GB RAM
 * instances are modified one at a time

**vse8.yml**

* modifies service from vse7.yml
* creates network r3-dc2-r3test4-db
 * connected to router vse4
 * with network 10.2.0.0/24

**vse9.yml**

* modifies service from vse8.yml
* creates instance r3-dc2-r3test4-db-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test4-db
 * has IP 10.2.0.11

**vse10.yml**

* modifies service from vse9.yml
* deletes instance r3-dc2-r3test4-web-2:

**vse11.yml**

* modifies service from vse10.yml
* deletes instance r3-dc2-r3test4-db-1:

**cleanup**

* `ernest service destroy r3test4`

**vse12.yml**

* initial service apply
* in datacenter r3-dc2
* creates router vse5
 * connects external network r3-net-ext
 * assigns external IP for NAT (NAT-IP)
* creates network r3-dc2-r3test5-web
 * connected to router vse5
 * with network 10.1.0.0/24
* creates network r3-dc2-r3test5-salt
 * connected to router vse5
 * with network 10.254.254.0/24
* creates instance r3-dc2-r3test5-web-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test5-web
 * has IP 10.1.0.11
 * is bootstrapped
 * `date` command has run
* creates instance r3-dc2-r3test5-salt-master:
 * from template r3/r3-salt-master
 * it is running
 * has 1CPU
 * has 2GB RAM
 * is on network r3-dc2-r3test5-salt
 * has IP 10.254.254.100
* configures ACLs on router vse5
 * allow internal:any to internal:any
 * allow internal:any to external:any
 * allow 10.254.254.0/24:any to any:22
 * allow 10.254.254.0/24:any to any:5985
 * allow 172.17.241.95:any to NAT-IP:8000
 * allow 172.17.241.95:any to NAT-IP:22
 * allow 10.1.0.0/24:any to 10.254.254.100:4505
 * allow 10.1.0.0/24:any to 10.254.254.100:4506
* configures NATs on router vse5
 * from 10.1.0.0/24:any to NAT-IP:any
 * from 10.254.254.0/24:any to NAT-IP:any
 * from NAT-IP:22 to 10.254.254.100:22
 * from NAT-IP:8000 to 10.254.254.100:8000

**vse13.yml**

* modifies service from vse12.yml
* creates instance r3-dc2-r3test5-web-2:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test5-web
 * has IP 10.1.0.12
 * is bootstrapped
 * `date` command has run

**vse14.yml**

* modifies service from vse13.yml
* modifies instances r3-dc2-r3test5-web-1 and r3-dc2-r3test5-web-2:
 * `date` and `uptime` commands have run

**vse15.yml**

* modifies service from vse14.yml
* creates instance r3-dc2-r3test5-db-1:
 * from template r3/ubuntu-1404
 * it is running
 * has 1CPU
 * has 1GB RAM
 * is on network r3-dc2-r3test5-web
 * has IP 10.1.0.21
 * is bootstrapped
 * `date` command has run

**vse16.yml**

* modifies service from vse15.yml
* modifies instance r3-dc2-r3test5-web-2:
 * has been unregistered with salt master

**cleanup**

* `ernest service destroy r3test5`
