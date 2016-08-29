/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func Test2VSE(t *testing.T) {
	var service = "vse"

	service = service + strconv.Itoa(rand.Intn(9999999))

	inCreateSub := make(chan *nats.Msg, 1)
	fiCreateSub := make(chan *nats.Msg, 1)
	roCreateSub := make(chan *nats.Msg, 1)
	neCreateSub := make(chan *nats.Msg, 1)
	naCreateSub := make(chan *nats.Msg, 1)
	inUpdateSub := make(chan *nats.Msg, 1)
	basicSetup("vcloud")

	Convey("Given I have a configuraed ernest instance", t, func() {
		Convey("When I apply a valid vse12.yml definition", func() {
			subIn, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subRo, _ := n.ChanSubscribe("router.create.vcloud-fake", roCreateSub)
			subNe, _ := n.ChanSubscribe("network.create.vcloud-fake", neCreateSub)
			subFi, _ := n.ChanSubscribe("firewall.create.vcloud-fake", fiCreateSub)
			subNa, _ := n.ChanSubscribe("nat.create.vcloud-fake", naCreateSub)

			f := getDefinitionPath("vse12.yml", service)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 29)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating routers:"
					checkLines[5] = "\t1.1.1.1"
					checkLines[6] = "Routers successfully created"
					checkLines[7] = "Creating networks:"
					checkLines[8] = "\t- 10.254.254.0/24"
					checkLines[9] = "\t- 10.1.0.0/24"
					checkLines[10] = "Networks successfully created"
					checkLines[11] = "Creating instances:"
					checkLines[12] = "\t - fake-" + service + "-salt-master"
					checkLines[13] = "\t - fake-" + service + "-web-1"
					checkLines[14] = "Instances successfully created"
					checkLines[15] = "Updating instances:"
					checkLines[16] = "\t - fake-" + service + "-salt-master"
					checkLines[17] = "\t - fake-" + service + "-web-1"
					checkLines[18] = "Instances successfully updated"
					checkLines[19] = "Setting up firewalls:"
					checkLines[20] = "Firewalls Created"
					checkLines[21] = "Configuring nats"
					checkLines[22] = "Nats Created"
					checkLines[23] = "Bootstrapping"
					checkLines[24] = "Instances bootstrapped"
					checkLines[25] = "Running executions"
					checkLines[26] = "Executions ran"
					checkLines[27] = "SUCCESS: rules successfully applied"
					checkLines[28] = "Your environment endpoint is: 1.1.1.1"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				r := routerEvent{}
				msg, err := waitMsg(roCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &r)
				n := networkEvent{}
				msg, err = waitMsg(neCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &n)
				i := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				f := firewallEvent{}
				msg, err = waitMsg(fiCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &f)
				na := natEvent{}
				msg, err = waitMsg(naCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &na)

				Info("And it creates router vse5", " ", 8)
				So(r.DatacenterName, ShouldEqual, "fake")
				So(r.DatacenterPassword, ShouldEqual, default_pwd)
				So(r.DatacenterRegion, ShouldEqual, "")
				So(r.DatacenterType, ShouldEqual, "vcloud-fake")
				So(r.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(r.RouterName, ShouldEqual, "vse5")
				So(r.RouterType, ShouldEqual, "vcloud-fake")
				So(r.Service, ShouldNotEqual, "")
				So(r.ClientName, ShouldNotEqual, "")
				So(r.VCloudURL, ShouldNotEqual, "")
				So(r.VseURL, ShouldNotEqual, "")
				So(r.Status, ShouldEqual, "processing")

				Info("And it creates network *-salt", " ", 8)
				So(n.DatacenterName, ShouldEqual, "fake")
				So(n.DatacenterPassword, ShouldEqual, default_pwd)
				So(n.DatacenterType, ShouldEqual, "vcloud-fake")
				So(n.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(n.NetworkType, ShouldEqual, "vcloud-fake")
				So(n.NetworkName, ShouldEqual, "fake-"+service+"-salt")
				So(n.NetworkGateway, ShouldEqual, "10.254.254.1")
				So(n.NetworkNetmask, ShouldEqual, "255.255.255.0")
				So(n.NetworkStartAddress, ShouldEqual, "10.254.254.5")
				So(n.NetworkEndAddress, ShouldEqual, "10.254.254.250")

				Info("And it creates instance *-salt-master", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-salt-master")
				So(i.Resource.CPU, ShouldEqual, 1)
				So(len(i.Resource.Disks), ShouldEqual, 0)
				So(i.Resource.IP, ShouldEqual, "10.254.254.100")
				So(i.Resource.RAM, ShouldEqual, 2048)
				So(i.Resource.Catalog, ShouldEqual, "r3")
				So(i.Resource.Image, ShouldEqual, "r3-salt-master")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-salt")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")
				So(i.Service, ShouldNotEqual, "")

				Info("Then it configures ACLs on router vse5", " ", 8)
				So(f.DatacenterName, ShouldEqual, "fake")
				So(f.DatacenterName, ShouldEqual, "fake")
				So(f.DatacenterPassword, ShouldEqual, default_pwd)
				So(f.DatacenterType, ShouldEqual, "vcloud-fake")
				So(f.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(f.FirewallType, ShouldEqual, "vcloud")
				So(len(f.Rules), ShouldEqual, 9)
				So(f.RouterIP, ShouldEqual, "1.1.1.1")
				So(f.RouterName, ShouldEqual, "vse5")
				So(f.RouterType, ShouldEqual, "vcloud-fake")
				So(f.Service, ShouldNotEqual, "")
				Printf("\n        And it will allow 10.254.254.0/24:any to 22:tcp ")
				So(f.Rules[0].SourcePort, ShouldEqual, "any")
				So(f.Rules[0].SourceIP, ShouldEqual, "10.254.254.0/24")
				So(f.Rules[0].DestinationIP, ShouldEqual, "any")
				So(f.Rules[0].DestinationPort, ShouldEqual, "22")
				So(f.Rules[0].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 10.254.254.0/24:any to internal:any ")
				So(f.Rules[1].SourcePort, ShouldEqual, "any")
				So(f.Rules[1].SourceIP, ShouldEqual, "10.254.254.0/24")
				So(f.Rules[1].DestinationIP, ShouldEqual, "any")
				So(f.Rules[1].DestinationPort, ShouldEqual, "5985")
				So(f.Rules[1].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow iternal:any to external:any ")
				So(f.Rules[2].SourcePort, ShouldEqual, "any")
				So(f.Rules[2].SourceIP, ShouldEqual, "internal")
				So(f.Rules[2].DestinationIP, ShouldEqual, "external")
				So(f.Rules[2].DestinationPort, ShouldEqual, "any")
				So(f.Rules[2].Protocol, ShouldEqual, "any")
				Printf("\n        And it will allow 172.17.241.221:any to 1.1.1.1:8000 ")
				So(f.Rules[3].SourcePort, ShouldEqual, "any")
				So(f.Rules[3].SourceIP, ShouldEqual, "172.17.241.221")
				So(f.Rules[3].DestinationIP, ShouldEqual, "1.1.1.1")
				So(f.Rules[3].DestinationPort, ShouldEqual, "8000")
				So(f.Rules[3].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 172.17.240.161:any to 1.1.1.1:8000 ")
				So(f.Rules[4].SourcePort, ShouldEqual, "any")
				So(f.Rules[4].SourceIP, ShouldEqual, "172.17.240.161")
				So(f.Rules[4].DestinationIP, ShouldEqual, "1.1.1.1")
				So(f.Rules[4].DestinationPort, ShouldEqual, "8000")
				So(f.Rules[4].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 10.1.0.0/24:any to 10.254.254.100:4505 ")
				So(f.Rules[5].SourcePort, ShouldEqual, "any")
				So(f.Rules[5].SourceIP, ShouldEqual, "10.1.0.0/24")
				So(f.Rules[5].DestinationIP, ShouldEqual, "10.254.254.100")
				So(f.Rules[5].DestinationPort, ShouldEqual, "4505")
				So(f.Rules[5].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow 10.1.0.0/24:any to 10.254.254.100:4506 ")
				So(f.Rules[6].SourcePort, ShouldEqual, "any")
				So(f.Rules[6].SourceIP, ShouldEqual, "10.1.0.0/24")
				So(f.Rules[6].DestinationIP, ShouldEqual, "10.254.254.100")
				So(f.Rules[6].DestinationPort, ShouldEqual, "4506")
				So(f.Rules[6].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it will allow internal:any to internal:any ")
				So(f.Rules[7].SourcePort, ShouldEqual, "any")
				So(f.Rules[7].SourceIP, ShouldEqual, "internal")
				So(f.Rules[7].DestinationIP, ShouldEqual, "internal")
				So(f.Rules[7].DestinationPort, ShouldEqual, "any")
				So(f.Rules[7].Protocol, ShouldEqual, "any")
				Printf("\n        And it will allow internal:any to external:any ")
				So(f.Rules[8].SourcePort, ShouldEqual, "any")
				So(f.Rules[8].SourceIP, ShouldEqual, "internal")
				So(f.Rules[8].DestinationIP, ShouldEqual, "external")
				So(f.Rules[8].DestinationPort, ShouldEqual, "any")
				So(f.Rules[8].Protocol, ShouldEqual, "any")

				Info("Then it configures NATs on router vse5", " ", 8)
				So(na.DatacenterName, ShouldEqual, "fake")
				So(na.DatacenterPassword, ShouldEqual, default_pwd)
				So(na.DatacenterRegion, ShouldEqual, "")
				So(na.DatacenterType, ShouldEqual, "vcloud-fake")
				So(na.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(na.NatName, ShouldEqual, "fake-"+service+"-vse5")
				So(len(na.NatRules), ShouldEqual, 4)
				So(na.RouterIP, ShouldEqual, "1.1.1.1")
				So(na.RouterName, ShouldEqual, "vse5")
				So(na.RouterType, ShouldEqual, "vcloud-fake")
				Printf("\n        And it configures from 1.1.1.1/24:8000 to 10.254.254.100:8000")
				So(na.NatRules[0].Network, ShouldEqual, "fake")
				So(na.NatRules[0].OriginIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[0].OriginPort, ShouldEqual, "8000")
				So(na.NatRules[0].Type, ShouldEqual, "dnat")
				So(na.NatRules[0].TranslationIP, ShouldEqual, "10.254.254.100")
				So(na.NatRules[0].TranslationPort, ShouldEqual, "8000")
				So(na.NatRules[0].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it configures from 1.1.1.1/24:22 to 10.254.254.100:22")
				So(na.NatRules[1].Network, ShouldEqual, "fake")
				So(na.NatRules[1].OriginIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[1].OriginPort, ShouldEqual, "22")
				So(na.NatRules[1].Type, ShouldEqual, "dnat")
				So(na.NatRules[1].TranslationIP, ShouldEqual, "10.254.254.100")
				So(na.NatRules[1].TranslationPort, ShouldEqual, "22")
				So(na.NatRules[1].Protocol, ShouldEqual, "tcp")
				Printf("\n        And it configures from 10.254.254.0/24:any to 1.1.1.1:any")
				So(na.NatRules[2].Network, ShouldEqual, "fake")
				So(na.NatRules[2].OriginIP, ShouldEqual, "10.254.254.0/24")
				So(na.NatRules[2].OriginPort, ShouldEqual, "any")
				So(na.NatRules[2].Type, ShouldEqual, "snat")
				So(na.NatRules[2].TranslationIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[2].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[2].Protocol, ShouldEqual, "any")
				Printf("\n        And it configures from 10.1.0.0/24:any to 1.1.1.1:any")
				So(na.NatRules[3].Network, ShouldEqual, "fake")
				So(na.NatRules[3].OriginIP, ShouldEqual, "10.1.0.0/24")
				So(na.NatRules[3].OriginPort, ShouldEqual, "any")
				So(na.NatRules[3].Type, ShouldEqual, "snat")
				So(na.NatRules[3].TranslationIP, ShouldEqual, "1.1.1.1")
				So(na.NatRules[3].TranslationPort, ShouldEqual, "any")
				So(na.NatRules[3].Protocol, ShouldEqual, "any")
			})

			subIn.Unsubscribe()
			subRo.Unsubscribe()
			subNe.Unsubscribe()
			subFi.Unsubscribe()
			subNa.Unsubscribe()
		})

		Convey("When I apply a valid vse13.yml definition", func() {
			subInCreate, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("vse13.yml", service)
			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 9)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fake-" + service + "-web-2"
					checkLines[7] = "Updating instances:"
					checkLines[8] = "\t - fake-" + service + "-web-2"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("And it will create web-2 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(i.Resource.CPU, ShouldEqual, 1)
				So(len(i.Resource.Disks), ShouldEqual, 0)
				So(i.Resource.IP, ShouldEqual, "10.1.0.12")
				So(i.Resource.RAM, ShouldEqual, 1024)
				So(i.Resource.Catalog, ShouldEqual, "r3")
				So(i.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")
				So(i.Service, ShouldNotEqual, "")

				Info("Then it will update web-2 instance", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterRegion, ShouldEqual, "")
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-web-2")
				So(iu.Resource.CPU, ShouldEqual, 1)
				So(len(iu.Resource.Disks), ShouldEqual, 0)
				So(iu.Resource.IP, ShouldEqual, "10.1.0.12")
				So(iu.Resource.RAM, ShouldEqual, 1024)
				So(iu.Resource.Catalog, ShouldEqual, "r3")
				So(iu.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
				So(iu.Service, ShouldNotEqual, "")

			})

			subInCreate.Unsubscribe()
			subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse14.yml definition", func() {

			f := getDefinitionPath("vse14.yml", service)
			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 5)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Running executions"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}
			})
			//TODO : we may need to check executions here
		})

		Convey("When I apply a valid vse15.yml definition", func() {
			subInCreate, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			subInUpdate, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			f := getDefinitionPath("vse15.yml", service)
			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 9)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fake-" + service + "-db-1"
					checkLines[7] = "Updating instances:"
					checkLines[8] = "\t - fake-" + service + "-db-1"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				i := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &i)
				iu := instanceEvent{}
				msg, err = waitMsg(inUpdateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &iu)

				Info("Then it will create db-1 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-db-1")
				So(i.Resource.CPU, ShouldEqual, 1)
				So(len(i.Resource.Disks), ShouldEqual, 0)
				So(i.Resource.IP, ShouldEqual, "10.1.0.21")
				So(i.Resource.RAM, ShouldEqual, 1024)
				So(i.Resource.Catalog, ShouldEqual, "r3")
				So(i.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")
				So(i.Service, ShouldNotEqual, "")

				Info("Then it will update db-1 instance", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterRegion, ShouldEqual, "")
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-db-1")
				So(iu.Resource.CPU, ShouldEqual, 1)
				So(len(iu.Resource.Disks), ShouldEqual, 0)
				So(iu.Resource.IP, ShouldEqual, "10.1.0.21")
				So(iu.Resource.RAM, ShouldEqual, 1024)
				So(iu.Resource.Catalog, ShouldEqual, "r3")
				So(iu.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "fake-"+service+"-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
				So(iu.Service, ShouldNotEqual, "")
			})

			subInCreate.Unsubscribe()
			subInUpdate.Unsubscribe()
		})

		Convey("When I apply a valid vse16.yml definition", func() {

			f := getDefinitionPath("vse16.yml", service)
			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 7)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Deleting instances:"
					checkLines[5] = "\t - fake-" + service + "-web-2"
					checkLines[6] = "Instances deleted"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}
			})
		})

	})
}
