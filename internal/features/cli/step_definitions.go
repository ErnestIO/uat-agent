package cli

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
	. "github.com/gucumber/gucumber"
	"github.com/nats-io/nats"
)

var lastOutput string
var lastError error
var cfg *ecc.Config
var n *nats.Conn

func init() {
	cfg = ecc.NewConfig(os.Getenv("NATS_URI"))
	n = cfg.Nats()

	Before("@login", func() {
		// runs before every feature or scenario tagged with @login
	})

	Given(`^I setup ernest with target "(.+?)"$`, func(target string) {
		if os.Getenv("CURRENT_INSTANCE") != "" {
			target = os.Getenv("CURRENT_INSTANCE")
		}

		ernest("target", target)
	})

	Given(`^I'm logged in as "(.+?)" / "(.+?)"$`, func(u, p string) {
		ernest("login", "--user", u, "--password", p)
	})

	When(`^I run ernest with "(.+?)"$`, func(args string) {
		cmdArgs := strings.Split(args, " ")
		ernest(cmdArgs...)
	})

	Then(`^The output should contain "(.+?)"$`, func(needle string) {
		if strings.Contains(lastOutput, needle) == false {
			T.Errorf(`Last output string does not contain "` + needle + `": ` + "\n" + lastOutput)
		}
	})

	Then(`^The output should not contain "(.+?)"$`, func(needle string) {
		if strings.Contains(lastOutput, needle) == true {
			T.Errorf(`Last output string does contains "` + needle + `" but it shouldn't: ` + "\n" + lastOutput)
		}
	})

	When(`^I logout$`, func() {
		ernest("logout")
	})

	When(`^I enter text "(.+?)"$`, func(input string) {
		cmd := exec.Command("ernest-cli", input)
		o, err := cmd.CombinedOutput()
		lastOutput = string(o)
		lastError = err
	})

	And(`^the group "(.+?)" does not exist$`, func(group string) {
		msg := []byte(`{"name":"` + group + `"}`)
		n.Request("group.del", msg, time.Second*3)
	})

	And(`^the user "(.+?)" does not exist$`, func(user string) {
		msg := []byte(`{"username":"` + user + `"}`)
		n.Request("user.del", msg, time.Second*3)
	})

	And(`^the datacenter "(.+?)" does not exist$`, func(d string) {
		msg := []byte(`{"name":"` + d + `", "type":"aws"}`)
		n.Request("datacenter.del", msg, time.Second*3)
	})

	And(`^the service "(.+?)" does not exist$`, func(d string) {
		msg := []byte(`{"name":"` + d + `", "type":"aws"}`)
		n.Request("service.del", msg, time.Second*3)
	})

	And(`^the group "(.+?)" exists$`, func(group string) {
		msg := []byte(`{"name":"` + group + `"}`)
		n.Request("group.del", msg, time.Second*3)
		msg = []byte(`{"name":"` + group + `"}`)
		n.Request("group.set", msg, time.Second*3)
	})

	And(`^the user "(.+?)" exists$`, func(user string) {
		msg := []byte(`{"username":"` + user + `"}`)
		n.Request("user.del", msg, time.Second*3)
		msg = []byte(`{"username":"` + user + `","password":"pwd"}`)
		n.Request("user.set", msg, time.Second*3)
	})

	And(`^the datacenter "(.+?)" exists$`, func(d string) {
		msg := []byte(`{"name":"` + d + `", "type":"aws"}`)
		n.Request("datacenter.del", msg, time.Second*3)
		msg = []byte(`{"name":"` + d + `"}`)
		n.Request("datacenter.set", msg, time.Second*3)
	})

	Then(`^The output users table should contain "(.+?)" assigned to "(.+?)" group$`, func(user string, group string) {
		lines := strings.Split(lastOutput, "\n")
		for _, l := range lines {
			if strings.Contains(l, user) {
				if !strings.Contains(l, "| "+group) {
					T.Errorf(`User doesn't seem to belong to specified group: \n` + l)
				}
			}
		}
	})

	Then(`^The output datacenters table should contain "(.+?)" assigned to "(.+?)" group$`, func(datacenter string, group string) {
		lines := strings.Split(lastOutput, "\n")
		for _, l := range lines {
			if strings.Contains(l, datacenter) {
				if !strings.Contains(l, "| "+group) {
					T.Errorf(`Datacenter doesn't seem to belong to specified group: \n` + l)
				}
			}
		}
	})

	Then(`^The output line number "(.+?)" should contain "(.+?)"$`, func(number int, needle string) {
		lines := strings.Split(lastOutput, "\n")
		n := strconv.Itoa(number)

		if len(lines) < number {
			T.Errorf(`Last output has less than "` + n + `" lines : ` + "\n" + lastOutput)
			return
		}

		if strings.Contains(lines[number], needle) == false {
			T.Errorf(`Line "` + n + `" should contain "` + needle + `" but it doesn't: ` + "\n" + lastOutput)
		}

	})
}

func ernest(cmdArgs ...string) {
	cmd := exec.Command("ernest-cli", cmdArgs...)
	o, err := cmd.CombinedOutput()
	lastOutput = string(o)
	lastError = err
}
