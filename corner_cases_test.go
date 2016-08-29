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

func TestPatchService(t *testing.T) {
	var service = "corn"
	type ServiceCreate struct {
		ID string `json:"id"`
	}

	service = service + strconv.Itoa(rand.Intn(9999999))

	createEvent := ServiceCreate{}
	patchEvent := ServiceCreate{}

	inCreateSub := make(chan *nats.Msg, 1)
	inCreateServiceSub := make(chan *nats.Msg, 1)
	basicSetup("vcloud")

	Convey("Given I have a configuraed ernest instance", t, func() {
		Convey("When I apply a valid inst1.yml definition", func() {
			subIn, _ := n.ChanSubscribe("instance.create.vcloud-fake", inCreateSub)
			sub, _ := n.ChanSubscribe("service.create", inCreateServiceSub)

			f := getDefinitionPath("inst1.yml", service)

			Convey("Then I should successfully create a valid service", func() {

				Info("And user output should be correct", " ", 6)
				o, err := ernest("service", "apply", f)
				if err != nil {
					log.Println(err.Error())
				} else {
					lines := strings.Split(o, "\n")
					checkLines := make([]string, 21)

					checkLines[0] = "Environment creation requested"

					vo := CheckOutput(lines, checkLines)
					if os.Getenv("CHECK_OUTPUT") != "" {
						So(vo, ShouldEqual, true)
					}
				}

				msg, err := waitMsg(inCreateServiceSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &createEvent)

				event := instanceEvent{}
				msg, err = waitMsg(inCreateSub)
				So(err, ShouldBeNil)
				json.Unmarshal(msg.Data, &event)

				Info("And I should receive a valid instance.create.vcloud-fake", " ", 8)
				So(event.DatacenterName, ShouldEqual, "fake")
				So(event.DatacenterPassword, ShouldEqual, default_pwd)
				So(event.DatacenterRegion, ShouldEqual, "")
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
				So(event.NetworkName, ShouldEqual, "fake-"+service+"-r3-dc2-r3vse1-db")
				So(event.RouterIP, ShouldEqual, "")
				So(event.RouterName, ShouldEqual, "")
				So(event.RouterType, ShouldEqual, "")
			})

			subIn.Unsubscribe()
			sub.Unsubscribe()

			waitToDone()
		})

		Convey("When this service is marked as errored", func() {
			n.Publish("service.set", []byte(`{"id":"`+createEvent.ID+`","status":"errored"}`))
			Convey("And I re-apply the same service", func() {
				sub, _ := n.ChanSubscribe("service.patch", inCreateSub)
				f := getDefinitionPath("inst1.yml", service)
				_, err := ernest("service", "apply", f)
				So(err, ShouldBeNil)

				msg, err := waitMsg(inCreateSub)
				json.Unmarshal(msg.Data, &patchEvent)

				Info("And I should receive an event to re-create the service", " ", 8)
				So(patchEvent.ID, ShouldEqual, createEvent.ID)
				So(strings.Contains(string(msg.Data), `"service.create"`), ShouldBeFalse)

				So(err, ShouldBeNil)

				sub.Unsubscribe()
				// waitToDone()
			})
		})
	})
}
