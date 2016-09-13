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

func TestAWSHappyPath(t *testing.T) {
	var service = "aws"
	service = service + strconv.Itoa(rand.Intn(9999999))

	neSub := make(chan *nats.Msg, 1)
	inSub := make(chan *nats.Msg, 1)
	fiSub := make(chan *nats.Msg, 1)

	basicSetup("aws")

	Convey("Given I have a non existing aws definition", t, func() {
		Convey("When I apply aws1.yml", func() {
			f := getDefinitionPathAWS("aws1.yml", service)

			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			subFiC, _ := n.ChanSubscribe("firewall.create.aws-fake", fiSub)

			o, err := ernest("service", "apply", f)

			Convey("Then I should create a valid service", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 20)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating networks:"
					checkLines[5] = "\t- 10.1.0.0/24"
					checkLines[6] = "Networks successfully created"
					checkLines[7] = "Setting up firewalls:"
					checkLines[8] = "Firewalls Created"
					checkLines[9] = "Creating instances:"
					checkLines[10] = "\t - fakeaws-" + service + "-web-1"
					checkLines[11] = "Instances successfully created"
					checkLines[12] = "Configuring nats"
					checkLines[13] = "Nats Created"
					checkLines[14] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				event := awsNetworkEvent{}
				eventI := awsInstanceEvent{}
				eventF := awsFirewallEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeC.Unsubscribe()
				msg, err = waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInC.Unsubscribe()
				msg, err = waitMsg(fiSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventF)
				subFiC.Unsubscribe()

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				So(event.DatacenterAccessToken, ShouldEqual, "fake")
				So(event.DatacenterAccessKey, ShouldEqual, "secret")
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.1.0.0/24")

				Info("And should call firewall creator connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				So(eventF.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventF.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "")

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")

			})

			waitToDone()
		})

		Convey("When I apply aws2.yml", func() {
			f := getDefinitionPathAWS("aws2.yml", service)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create a new xx-web-2 instance", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 15)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating instances:"
					checkLines[5] = "\t - fakeaws-" + service + "-web-2"
					checkLines[6] = "Instances successfully created"
					checkLines[7] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-2")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws3.yml", func() {
			f := getDefinitionPathAWS("aws3.yml", service)
			subInC, _ := n.ChanSubscribe("instance.delete.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete xx-web-2 instance", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Deleting instances:"
					checkLines[5] = "\t - fakeaws-" + service + "-web-2"
					checkLines[6] = "Instances deleted"
					checkLines[7] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-2")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws4.yml", func() {
			f := getDefinitionPathAWS("aws4.yml", service)
			subInC, _ := n.ChanSubscribe("instance.update.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should update xx-web-1 instance", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Updating instances:"
					checkLines[5] = "\t - fakeaws-" + service + "-web-1"
					checkLines[6] = "Instances successfully updated"
					checkLines[7] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 0)
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws5.yml", func() {
			f := getDefinitionPathAWS("aws5.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should add an Ingress rule to existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Updating firewalls:"
					checkLines[5] = "Firewalls Updated"
					checkLines[6] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventF)
				subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				So(eventF.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventF.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Ingress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws6.yml", func() {
			f := getDefinitionPathAWS("aws6.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should add an Engress rule to existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Updating firewalls:"
					checkLines[5] = "Firewalls Updated"
					checkLines[6] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventF)
				subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				So(eventF.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventF.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Egress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Egress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Egress[1].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Ingress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws7.yml", func() {
			f := getDefinitionPathAWS("aws7.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete previously added engress and ingress rules from  existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Updating firewalls:"
					checkLines[5] = "Firewalls Updated"
					checkLines[6] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventF)
				subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				So(eventF.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventF.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws8.yml", func() {
			f := getDefinitionPathAWS("aws8.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating networks:"
					checkLines[5] = "\t- 10.2.0.0/24"
					checkLines[6] = "Networks successfully created"
					checkLines[7] = "Configuring nats"
					checkLines[8] = "Nats Created"
					checkLines[9] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeC.Unsubscribe()

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				So(event.DatacenterAccessToken, ShouldEqual, "fake")
				So(event.DatacenterAccessKey, ShouldEqual, "secret")
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
			})
			waitToDone()
		})

		Convey("When I apply aws9.yml", func() {
			f := getDefinitionPathAWS("aws9.yml", service)
			subNeC, _ := n.ChanSubscribe("network.delete.aws-fake", neSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete network 10.2.0.0/24", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 11)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Deleting nats"
					checkLines[5] = "Nats Deleted"
					checkLines[6] = "Deleting networks:"
					checkLines[7] = "\t- 10.2.0.0/24"
					checkLines[8] = "Networks deleted"
					checkLines[9] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeC.Unsubscribe()

				Info("And should call network deleter connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				So(event.DatacenterAccessToken, ShouldEqual, "fake")
				So(event.DatacenterAccessKey, ShouldEqual, "secret")
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")

			})
			waitToDone()
		})

		SkipConvey("When I apply aws10.yml", func() {
			f := getDefinitionPathAWS("aws10.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 13)
					println(o)

					checkLines[0] = "Environment creation requested"
					checkLines[4] = "Creating networks:"
					checkLines[5] = "\t- 10.1.0.0/24"
					checkLines[6] = "Networks successfully created"
					checkLines[7] = "Creating instances:"
					checkLines[8] = "\t - fakeaws-" + service + "-web-2"
					checkLines[9] = "Instances successfully created"
					checkLines[10] = "Configuring nats"
					checkLines[11] = "Nats Created"
					checkLines[12] = "SUCCESS: rules successfully applied"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeC.Unsubscribe()

				eventI := awsInstanceEvent{}

				msg, err = waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-bknd-2")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				So(event.DatacenterAccessToken, ShouldEqual, "fake")
				So(event.DatacenterAccessKey, ShouldEqual, "secret")
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
			})
			waitToDone()
		})

	})
}
