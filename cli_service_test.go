/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCLIServiceInfo(t *testing.T) {
	var service = "service-info"
	service = service + strconv.Itoa(rand.Intn(9999999))
	basicSetup("aws")

	f := getDefinitionPathAWS("aws1.yml", service)
	ernest("service", "apply", f)
	f = getDefinitionPathAWS("aws2.yml", service)
	ernest("service", "apply", f)
	f = getDefinitionPathAWS("aws3.yml", service)
	ernest("service", "apply", f)

	Convey("Scenario : service info", t, func() {
		Convey("Given I’m logged in as a plain user", func() {
			Info("When <service_name> does not exist", " ", 6)
			Info("And I run “ernest service info <service_name>”", " ", 8)
			Info("Then I should see a warning message", " ", 8)
			o, err := ernest("service", "info", "unexisting")
			So(err, ShouldBeNil)
			So(o, ShouldEqual, "Specified service not found\n")

			Info("When <service_name> does exist", " ", 6)
			Info("And I run “ernest service info <service_name>”", " ", 6)
			o, err = ernest("service", "info", service)
			Info("Then I should see information about the last build of provided service", " ", 8)
			So(err, ShouldBeNil)
			lines := strings.Split(o, "\n")
			So(lines[0], ShouldEqual, "Name : "+service)
			So(lines[6], ShouldEqual, "| fakeaws-"+service+"-web | foo |")
			So(lines[13], ShouldEqual, "| fakeaws-"+service+"-web-1 |    |           | 10.1.0.11  |")
			So(lines[23], ShouldEqual, "| fakeaws-"+service+"-web-sg-1 | foo      |")

			Info("When I run “ernest service info <service_name> --build <non_existing_build>”", " ", 8)
			Info("And non_existing_build does not exist", " ", 8)
			o, err = ernest("service", "info", service, "--build", "unexisting")
			Info("Then I should see a warning message", " ", 8)
			So(o, ShouldEqual, "Specified build not found\n")

			Info("When I run “ernest service info <service_name> --build <my_build>”", " ", 8)
			o, err = ernest("service", "history", service)
			lines = strings.Split(o, "\n")
			cols := strings.Split(lines[1], "\t")

			Info("And my_build does exist", " ", 8)
			o, err = ernest("service", "info", service, "--build", cols[1])
			Info("Then I should see information about the specified build of provided service", " ", 8)
			lines = strings.Split(o, "\n")
			So(lines[0], ShouldEqual, "Name : "+service)
			So(lines[6], ShouldEqual, "| fakeaws-"+service+"-web | foo |")
			So(lines[13], ShouldEqual, "| fakeaws-"+service+"-web-1 |    |           | 10.1.0.11  |")
			So(lines[23], ShouldEqual, "| fakeaws-"+service+"-web-sg-1 | foo      |")
		})

	})
}
