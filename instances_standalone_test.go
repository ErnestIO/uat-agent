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

func TestStandAloneInstances(t *testing.T) {
	var service = "inst"

	service = service + strconv.Itoa(rand.Intn(9999999))

	inCreateSub := make(chan *nats.Msg, 1)
	inUpdateSub := make(chan *nats.Msg, 1)
	inDeleteSub := make(chan *nats.Msg, 1)
	basicSetup("vcloud")

	Convey("Given I have a configured ernest instance", t, func() {
		Convey("When I apply a valid inst1.yml definition", func() {
			f := getDefinitionPath("inst1.yml", service)
			sub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fake-" + service + "-stg-1"
					checkLines[7] = "Updating instances:"
					checkLines[8] = "\t - fake-" + service + "-stg-1"
					checkLines[10] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				event := instanceEvent{}
				msg, err := waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.InstanceName, ShouldEqual, "fake-"+service+"-stg-1")
				So(event.Resource.CPU, ShouldEqual, 1)
				So(len(event.Resource.Disks), ShouldEqual, 0)
				So(event.Resource.IP, ShouldEqual, "10.2.0.90")
				So(event.Resource.RAM, ShouldEqual, 1024)
				So(event.Resource.Catalog, ShouldEqual, "r3")
				So(event.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(event.InstanceType, ShouldEqual, "vcloud-fake")
				So(event.NetworkName, ShouldEqual, "r3-dc2-r3vse1-db")
				So(event.RouterIP, ShouldEqual, "")
				So(event.RouterName, ShouldEqual, "")
				So(event.RouterType, ShouldEqual, "")
				So(event.Service, ShouldNotEqual, "")
			})

			sub.Unsubscribe()
		})

		Convey("When I add an extra instance with inst2.yml definition", func() {
			f := getDefinitionPath("inst2.yml", service)

			csub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			usub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fake-" + service + "-stg-2"
					checkLines[7] = "Updating instances:"
					checkLines[8] = "\t - fake-" + service + "-stg-2"
					checkLines[10] = "SUCCESS: rules successfully applied"

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

				Info("And it will create stg-2 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-stg-2")
				So(i.Resource.CPU, ShouldEqual, 1)
				So(len(i.Resource.Disks), ShouldEqual, 0)
				So(i.Resource.IP, ShouldEqual, "10.2.0.91")
				So(i.Resource.RAM, ShouldEqual, 1024)
				So(i.Resource.Catalog, ShouldEqual, "r3")
				So(i.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "r3-dc2-r3vse1-db")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")
				So(i.Service, ShouldNotEqual, "")

				Info("And it will update stg-2 instance", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-stg-2")
				So(iu.Resource.CPU, ShouldEqual, 1)
				So(len(iu.Resource.Disks), ShouldEqual, 0)
				So(iu.Resource.IP, ShouldEqual, "10.2.0.91")
				So(iu.Resource.RAM, ShouldEqual, 1024)
				So(iu.Resource.Catalog, ShouldEqual, "r3")
				So(iu.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "r3-dc2-r3vse1-db")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
				So(iu.Service, ShouldNotEqual, "")

			})

			csub.Unsubscribe()
			usub.Unsubscribe()
		})
		//time.Sleep(time.Second)

		Convey("When I add an extra instance and modifies the existing one with inst3.yml definition", func() {
			f := getDefinitionPath("inst3.yml", service)

			csub, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			usub, _ := n.ChanSubscribe("instance.update.vcloud-fake", inUpdateSub)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fake-" + service + "-dev-1"
					checkLines[7] = "Updating instances:"
					checkLines[8] = "\t - fake-" + service + "-dev-1"
					checkLines[10] = "SUCCESS: rules successfully applied"
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

				Info("And it will create dev-1 instance", " ", 8)
				So(i.DatacenterName, ShouldEqual, "fake")
				So(i.DatacenterPassword, ShouldEqual, default_pwd)
				So(i.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(i.DatacenterType, ShouldEqual, "vcloud-fake")
				So(i.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(i.InstanceName, ShouldEqual, "fake-"+service+"-dev-1")
				So(i.Resource.CPU, ShouldEqual, 1)
				So(len(i.Resource.Disks), ShouldEqual, 0)
				So(i.Resource.IP, ShouldEqual, "10.1.0.90")
				So(i.Resource.RAM, ShouldEqual, 1024)
				So(i.Resource.Catalog, ShouldEqual, "r3")
				So(i.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(i.InstanceType, ShouldEqual, "vcloud-fake")
				So(i.NetworkName, ShouldEqual, "r3-dc2-r3vse1-web")
				So(i.RouterIP, ShouldEqual, "")
				So(i.RouterName, ShouldEqual, "")
				So(i.RouterType, ShouldEqual, "")
				So(i.Service, ShouldNotEqual, "")

				Info("And it will update dev-1 instance", " ", 8)
				So(iu.DatacenterName, ShouldEqual, "fake")
				So(iu.DatacenterPassword, ShouldEqual, default_pwd)
				So(iu.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(iu.DatacenterType, ShouldEqual, "vcloud-fake")
				So(iu.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(iu.InstanceName, ShouldEqual, "fake-"+service+"-dev-1")
				So(iu.Resource.CPU, ShouldEqual, 1)
				So(len(iu.Resource.Disks), ShouldEqual, 0)
				So(iu.Resource.IP, ShouldEqual, "10.1.0.90")
				So(iu.Resource.RAM, ShouldEqual, 1024)
				So(iu.Resource.Catalog, ShouldEqual, "r3")
				So(iu.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(iu.InstanceType, ShouldEqual, "vcloud-fake")
				So(iu.NetworkName, ShouldEqual, "r3-dc2-r3vse1-web")
				So(iu.RouterIP, ShouldEqual, "")
				So(iu.RouterName, ShouldEqual, "")
				So(iu.RouterType, ShouldEqual, "")
				So(iu.Service, ShouldNotEqual, "")
			})
			csub.Unsubscribe()
			usub.Unsubscribe()
		})

		Convey("When I delete stg-2 from  inst4.yml definition", func() {
			f := getDefinitionPath("inst4.yml", service)

			dsub, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 8)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Deleting instances:"
					checkLines[5] = "\t - fake-" + service + "-stg-2"
					checkLines[7] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				Info("And it will delete stg-2 instance", " ", 8)
				event := instanceEvent{}
				msg, err := waitMsg(inDeleteSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
				So(event.DatacenterType, ShouldEqual, "vcloud-fake")
				So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
				So(event.InstanceName, ShouldEqual, "fake-"+service+"-stg-2")
				So(event.Resource.CPU, ShouldEqual, 1)
				So(len(event.Resource.Disks), ShouldEqual, 0)
				So(event.Resource.IP, ShouldEqual, "10.2.0.91")
				So(event.Resource.RAM, ShouldEqual, 1024)
				So(event.Resource.Catalog, ShouldEqual, "r3")
				So(event.Resource.Image, ShouldEqual, "ubuntu-1404")
				So(event.InstanceType, ShouldEqual, "vcloud-fake")
				So(event.NetworkName, ShouldEqual, "r3-dc2-r3vse1-db")
				So(event.RouterIP, ShouldEqual, "")
				So(event.RouterName, ShouldEqual, "")
				So(event.RouterType, ShouldEqual, "")
				So(event.Service, ShouldNotEqual, "")
			})

			dsub.Unsubscribe()
		})

		Convey("When I delete stg-1 instance from  inst5.yml definition", func() {
			f := getDefinitionPath("inst5.yml", service)

			dsub, _ := n.ChanSubscribe("instance.delete.vcloud-fake", inDeleteSub)

			o, err := ernest("service", "apply", f)
			Convey("Then I should get a valid output for a processed service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 8)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Deleting instances:"
					checkLines[5] = "\t - fake-" + service + "-stg-1"
					checkLines[7] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}
			})

			event := instanceEvent{}
			msg, err := waitMsg(inDeleteSub)
			So(err, ShouldBeNil)
			json.Unmarshal(msg.Data, &event)

			Info("And it will delete stg-2 instance", " ", 8)
			So(event.DatacenterName, ShouldEqual, "fake")
			So(event.DatacenterPassword, ShouldEqual, default_pwd)
			So(event.DatacenterRegion, ShouldEqual, "$(datacenters.items.0.region)")
			So(event.DatacenterType, ShouldEqual, "vcloud-fake")
			So(event.DatacenterUsername, ShouldEqual, default_usr+"@"+default_org)
			So(event.InstanceName, ShouldEqual, "fake-"+service+"-stg-1")
			So(event.Resource.CPU, ShouldEqual, 1)
			So(len(event.Resource.Disks), ShouldEqual, 0)
			So(event.Resource.IP, ShouldEqual, "10.2.0.90")
			So(event.Resource.RAM, ShouldEqual, 1024)
			So(event.Resource.Catalog, ShouldEqual, "r3")
			So(event.Resource.Image, ShouldEqual, "ubuntu-1404")
			So(event.InstanceType, ShouldEqual, "vcloud-fake")
			So(event.NetworkName, ShouldEqual, "r3-dc2-r3vse1-db")
			So(event.RouterIP, ShouldEqual, "")
			So(event.RouterName, ShouldEqual, "")
			So(event.RouterType, ShouldEqual, "")
			So(event.Service, ShouldNotEqual, "")

			dsub.Unsubscribe()
		})

	})

}
