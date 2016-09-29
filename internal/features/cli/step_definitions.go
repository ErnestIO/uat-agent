package cli

import (
	"os"
	"os/exec"
	"strings"

	. "github.com/gucumber/gucumber"
)

var lastOutput string
var lastError error

func init() {
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
}

func ernest(cmdArgs ...string) {
	cmd := exec.Command("ernest-cli", cmdArgs...)
	o, err := cmd.CombinedOutput()
	lastOutput = string(o)
	lastError = err
}
