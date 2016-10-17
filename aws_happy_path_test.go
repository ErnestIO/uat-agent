/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"math/rand"
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
	naSub := make(chan *nats.Msg, 1)
	lbSub := make(chan *nats.Msg, 1)

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
					expected := `Starting environment creation
Creating Vpc:
 - fakeaws
   Subnet    : 1.1.1.1/24
   Status    : completed
Vpc created
Creating networks:
 - fakeaws-` + service + `-web
   IP     : 10.1.0.0/24
   AWS ID : foo
   Status : completed
Networks successfully created
Creating firewalls:
 - fakeaws-` + service + `-web-sg-1
   Status    : completed
Firewalls created
Creating instances:
 - fakeaws-` + service + `-web-1
   IP        : 10.1.0.11
   Status    : completed
Instances successfully created
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
					expected := `Starting environment creation
Creating instances:
 - fakeaws-` + service + `-web-2
   IP        : 10.1.0.12
   Status    : completed
Instances successfully created
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
			subInD, _ := n.ChanSubscribe("instance.delete.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete xx-web-2 instance", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Deleting instances:
 - fakeaws-` + service + `-web-2
   IP        : 10.1.0.12
   Status    : completed
Instances deleted
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInD.Unsubscribe()

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
					expected := `Starting environment creation
Updating instances:
 - fakeaws-` + service + `-web-1
   IP        : 10.1.0.11
   Status    : completed
Instances successfully updated
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
					expected := `Starting environment creation
Updating firewalls:
 - fakeaws-` + service + `-web-sg-1
   Status    : completed
Firewalls updated
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
			Convey("Then it should add an Egress rule to existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Updating firewalls:
 - fakeaws-` + service + `-web-sg-1
   Status    : completed
Firewalls updated
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
			Convey("Then it should delete previously added egress and ingress rules from  existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Updating firewalls:
 - fakeaws-` + service + `-web-sg-1
   Status    : completed
Firewalls updated
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
					expected := `Starting environment creation
Creating networks:
 - fakeaws-` + service + `-bknd
   IP     : 10.2.0.0/24
   AWS ID : foo
   Status : completed
Networks successfully created
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
					expected := `Starting environment creation
Deleting networks:
 - fakeaws-` + service + `-bknd
   IP     : 10.2.0.0/24
   AWS ID : foo
   Status : completed
Networks deleted
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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

		Convey("When I apply aws10.yml", func() {
			f := getDefinitionPathAWS("aws10.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Creating networks:
 - fakeaws-` + service + `-bknd
   IP     : 10.2.0.0/24
   AWS ID : foo
   Status : completed
Networks successfully created
Creating instances:
 - fakeaws-` + service + `-bknd-1
   IP        : 10.2.0.11
   Status    : completed
Instances successfully created
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
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
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-bknd-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")

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

		Convey("When I apply aws11.yml", func() {
			f := getDefinitionPathAWS("aws11.yml", service)
			subNeD, _ := n.ChanSubscribe("network.delete.aws-fake", neSub)
			subInD, _ := n.ChanSubscribe("instance.delete.aws-fake", inSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete the 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Deleting instances:
 - fakeaws-` + service + `-bknd-1
   IP        : 10.2.0.11
   Status    : completed
Instances deleted
Deleting networks:
 - fakeaws-` + service + `-bknd
   IP     : 10.2.0.0/24
   AWS ID : foo
   Status : completed
Networks deleted
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventI)
				subInD.Unsubscribe()

				event := awsNetworkEvent{}

				msg, err = waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeD.Unsubscribe()

				Info("And should call instance deleter connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				So(eventI.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventI.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-bknd-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "")

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

		Convey("When I apply aws12.yml", func() {
			f := getDefinitionPathAWS("aws12.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subNaC, _ := n.ChanSubscribe("nat.create.aws-fake", naSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Creating networks:
 - fakeaws-` + service + `-db
   IP     : 10.2.0.0/24
   AWS ID : foo
   Status : completed
Networks successfully created
Creating nats:
 - fakeaws-` + service + `-db-nat
   Status    : completed
Nats created
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)
				subNeC.Unsubscribe()

				eventN := awsNatEvent{}

				msg, err = waitMsg(naSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventN)
				subNaC.Unsubscribe()

				Info("And should call nat creator connector with valid fields", " ", 6)
				So(eventN.Type, ShouldEqual, "aws-fake")
				So(eventN.DatacenterRegion, ShouldEqual, "fake")
				So(eventN.DatacenterAccessToken, ShouldEqual, "fake")
				So(eventN.DatacenterAccessKey, ShouldEqual, "secret")
				So(eventN.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventN.PublicNetwork, ShouldEqual, "fakeaws-"+service+"-web")
				So(len(eventN.RoutedNetworks), ShouldEqual, 1)
				So(eventN.RoutedNetworks[0], ShouldEqual, "fakeaws-"+service+"-db")
				So(eventN.Status, ShouldEqual, "")

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				So(event.DatacenterAccessToken, ShouldEqual, "fake")
				So(event.DatacenterAccessKey, ShouldEqual, "secret")
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
				So(event.NetworkIsPublic, ShouldBeFalse)
			})
			waitToDone()
		})

		Convey("When I apply aws13.yml", func() {
			f := getDefinitionPathAWS("aws13.yml", service)
			subLBC, _ := n.ChanSubscribe("elb.create.aws-fake", lbSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should create the new elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Creating ELBs:
 - fakeaws-` + service + `-elb-1
   Status    : completed
   DNS    : fake-dns-name
ELBs created
SUCCESS: rules successfully applied`

					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				eventLB := awsELBEvent{}

				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventLB)
				subLBC.Unsubscribe()

				Info("And should call elb creator connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "fake")
				So(eventLB.DatacenterSecret, ShouldEqual, "secret")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")
				So(len(eventLB.InstanceNames), ShouldEqual, 1)
				So(len(eventLB.InstanceAWSIDs), ShouldEqual, 1)
				So(len(eventLB.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventLB.InstanceNames[0], ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventLB.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(len(eventLB.Listeners), ShouldEqual, 1)
				So(eventLB.Listeners[0].ToPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].FromPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].Protocol, ShouldEqual, "HTTP")
				So(eventLB.Listeners[0].SSLCert, ShouldEqual, "")
			})
			waitToDone()
		})

		Convey("When I apply aws14.yml", func() {
			f := getDefinitionPathAWS("aws14.yml", service)
			subLBU, _ := n.ChanSubscribe("elb.update.aws-fake", lbSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should update the elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Updating ELBs:
 - fakeaws-` + service + `-elb-1
   Status    : completed
   DNS    : fake-dns-name
ELBs updated
SUCCESS: rules successfully applied`

					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				eventLB := awsELBEvent{}

				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventLB)
				subLBU.Unsubscribe()

				Info("And should call elb updater connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "fake")
				So(eventLB.DatacenterSecret, ShouldEqual, "secret")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")
				So(len(eventLB.InstanceNames), ShouldEqual, 1)
				So(len(eventLB.InstanceAWSIDs), ShouldEqual, 1)
				So(len(eventLB.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventLB.InstanceNames[0], ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventLB.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(len(eventLB.Listeners), ShouldEqual, 2)
				So(eventLB.Listeners[0].ToPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].FromPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].Protocol, ShouldEqual, "HTTP")
				So(eventLB.Listeners[0].SSLCert, ShouldEqual, "")
				So(eventLB.Listeners[1].ToPort, ShouldEqual, 443)
				So(eventLB.Listeners[1].FromPort, ShouldEqual, 443)
				So(eventLB.Listeners[1].Protocol, ShouldEqual, "HTTPS")
				So(eventLB.Listeners[1].SSLCert, ShouldEqual, "foo")
			})
			waitToDone()
		})

		Convey("When I apply aws15.yml", func() {
			f := getDefinitionPathAWS("aws15.yml", service)
			subLBD, _ := n.ChanSubscribe("elb.delete.aws-fake", lbSub)
			o, err := ernest("service", "apply", f)
			Convey("Then it should delete the elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					expected := `Starting environment creation
Deleting ELBs:
 - fakeaws-` + service + `-elb-1
   Status    : completed
   DNS    : fake-dns-name
ELBs deleted
SUCCESS: rules successfully applied`
					So(strings.Contains(o, expected), ShouldBeTrue)
				}

				eventLB := awsELBEvent{}

				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &eventLB)
				subLBD.Unsubscribe()

				Info("And should call elb updater connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "fake")
				So(eventLB.DatacenterSecret, ShouldEqual, "secret")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")
			})
			waitToDone()
		})

	})
}
