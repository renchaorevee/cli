package space_test

import (
	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("space-ssh-allowed command", func() {
	var (
		ui                  *testterm.FakeUI
		requirementsFactory *testreq.FakeReqFactory
		deps                commandregistry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.UI = ui
		commandregistry.Commands.SetCommand(commandregistry.Commands.FindCommand("space-ssh-allowed").SetDependency(deps, pluginCall))
	}

	runCommand := func(args ...string) bool {
		return testcmd.RunCLICommand("space-ssh-allowed", args, requirementsFactory, updateCommandDependency, false)
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	Describe("requirements", func() {
		It("fails with usage when called without enough arguments", func() {
			requirementsFactory.LoginSuccess = true

			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Incorrect Usage", "Requires", "argument"},
			))

		})

		It("fails requirements when not logged in", func() {
			Expect(runCommand("my-space")).To(BeFalse())
		})

		It("does not pass requirements if org is not targeted", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = false

			Expect(runCommand("my-space")).To(BeFalse())
		})

		It("does not pass requirements if space does not exist", func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = true
			requirementsFactory.SpaceRequirementFails = true

			Expect(runCommand("my-space")).To(BeFalse())
		})
	})

	Describe("space-ssh-allowed", func() {
		var space models.Space

		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			requirementsFactory.TargetedOrgSuccess = true

			space = models.Space{}
			space.Name = "the-space-name"
			space.GUID = "the-space-guid"
		})

		Context("when SSH is enabled for the space", func() {
			It("notifies the user", func() {
				space.AllowSSH = true
				requirementsFactory.Space = space

				runCommand("the-space-name")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"ssh support is enabled in space 'the-space-name'"}))
			})
		})

		Context("when SSH is disabled for the space", func() {
			It("notifies the user", func() {
				requirementsFactory.Space = space

				runCommand("the-space-name")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"ssh support is disabled in space 'the-space-name'"}))
			})
		})
	})
})
