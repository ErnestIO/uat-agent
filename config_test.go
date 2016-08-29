/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErnestConfig(t *testing.T) {
	Convey("Given I have an unconfigured ernest instance", t, func() {
		deleteConfig()
		Convey("When I try to get ernest's info", func() {
			o, err := ernest("info")
			Convey("Then I should get a warning message", func() {
				if err != nil {
					log.Println(err.Error())
				} else {
					So(o, ShouldEqual, "Environment not configured, please use target command\nCurrent target : \nCurrent user : \n")
				}
			})
		})

		Convey("When I set ernest target", func() {
			o, err := ernest("target", ernest_instance)
			Convey("When I try to get ernest's info", func() {
				Convey("Then I will get the configured target", func() {
					if err != nil {
						log.Println(err.Error())
					} else {
						So(o, ShouldEndWith, "Target set\n")
					}
				})
			})
		})
	})
}
