# Testing Ernest Commands

## Pre-configured vShield Edges

| vDC | vShield | IP-1 | IP-2 |
|:---:|:-------:|:----:|:----:|
| r3-dc2 | vse1 | 172.16.186.42 | 172.16.186.60 |
| r3-dc2 | vse2 | 172.16.186.44 | 172.16.186.61 |

## Setup Accounts & Datacenters

```
ernest target https://172.17.241.221

ernest user create --user admin --password w4rmR3d user2 654321
ernest login --user user2 --password 654321
ernest datacenter create --datacenter-user r3test --datacenter-password 123456abcdef --datacenter-org r3-org1 --vse-url https://vse-creator-service.example.net r3-dc1 --url https://vcloud.example.net r3-net-ext

ernest user create --user admin --password w4rmR3d user1 123456
ernest login --user user1 --password 123456
ernest datacenter create --datacenter-user r3user --datacenter-password abcdef123456 --datacenter-org r3-org2 --vse-url https://vse-creator-service.example.net r3-dc2 https://vcloud.example.net r3-net-ext
```

## Test No-VSE Case

Non-Bootstrapped:
```
ernest service apply novse1.yml # initial creation
ernest service apply novse2.yml # firewall change
ernest service apply novse3.yml # port-forward change
ernest service apply novse4.yml # increase instance count
ernest service apply novse5.yml # add cpu
ernest service apply novse6.yml # add disk
ernest service apply novse7.yml # add ram
ernest service apply novse8.yml # add network
ernest service apply novse9.yml # add instance
ernest service apply novse10.yml # reduce instance count
ernest service apply novse11.yml # remove instance
```

Bootstrapped:
```
ernest service apply novse12.yml # initial creation
ernest service apply novse13.yml # increase instance count
ernest service apply novse14.yml # add command
ernest service apply novse15.yml # add instance
ernest service apply novse16.yml # reduce instance count
```

Run both tests in parallel. Do not proceed to testing the Optional Case until all tests have completed.

## Test Optional Case

Non-Bootstrapped:
```
ernest service apply inst1.yml # initial creation
ernest service apply inst2.yml # increase instance count
ernest service apply inst3.yml # add instance
ernest service apply inst4.yml # reduce instance count
ernest service apply inst5.yml # remove instance
```

The optional case does not support SALT bootstrapping.

## Test VSE Case

Non-Bootstrapped:
```
ernest service apply vse1.yml # initial creation
ernest service apply vse2.yml # firewall change
ernest service apply vse3.yml # port-forward change
ernest service apply vse4.yml # increase instance count
ernest service apply vse5.yml # add cpu
ernest service apply vse6.yml # add disk
ernest service apply vse7.yml # add ram
ernest service apply vse8.yml # add network
ernest service apply vse9.yml # add instance
ernest service apply vse10.yml # reduce instance count
ernest service apply vse11.yml # remove instance
```

Bootstrapped:
```
ernest service apply vse12.yml # initial creation
ernest service apply vse13.yml # increase instance count
ernest service apply vse14.yml # add command
ernest service apply vse15.yml # add instance
ernest service apply vse16.yml # reduce instance count
```

Run both tests in parallel.
