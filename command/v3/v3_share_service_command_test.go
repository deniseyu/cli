package v3_test

import (
	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/command/commandfakes"
	"code.cloudfoundry.org/cli/command/v3"
	"code.cloudfoundry.org/cli/command/v3/v3fakes"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("share-service Command", func() {
	var (
		cmd             v3.ShareServiceCommand
		testUI          *ui.UI
		fakeConfig      *commandfakes.FakeConfig
		fakeSharedActor *commandfakes.FakeSharedActor
		fakeActor       *v3fakes.FakeShareServiceActor
		binaryName      string
		executeErr      error
	)

	BeforeEach(func() {
		testUI = ui.NewTestUI(nil, NewBuffer(), NewBuffer())
		fakeConfig = new(commandfakes.FakeConfig)
		fakeSharedActor = new(commandfakes.FakeSharedActor)
		fakeActor = new(v3fakes.FakeShareServiceActor)

		cmd = v3.ShareServiceCommand{
			UI:          testUI,
			Config:      fakeConfig,
			SharedActor: fakeSharedActor,
			Actor:       fakeActor,
		}

		cmd.RequiredArgs.ServiceInstance = "some-service-instance"

		binaryName = "faceman"
		fakeConfig.BinaryNameReturns(binaryName)

		// TODO: test minimum version requirement
	})

	JustBeforeEach(func() {
		executeErr = cmd.Execute(nil)
	})

	// Context("when the API version is below the minimum", func() {
	// 	BeforeEach(func() {
	// 		fakeActor.CloudControllerAPIVersionReturns("0.0.0")
	// 	})

	// 	It("returns a MinimumAPIVersionNotMetError", func() {
	// 		Expect(executeErr).To(MatchError(translatableerror.MinimumAPIVersionNotMetError{
	// 			CurrentVersion: "0.0.0",
	// 			MinimumVersion: ccversion.MinVersionRunTaskV3,
	// 		}))
	// 	})
	// })

	Context("when checking target fails", func() {
		BeforeEach(func() {
			fakeSharedActor.CheckTargetReturns(actionerror.NotLoggedInError{BinaryName: binaryName})
		})

		It("returns an error", func() {
			// TODO: return translatable error
			Expect(executeErr).To(MatchError(actionerror.NotLoggedInError{BinaryName: binaryName}))

			Expect(fakeSharedActor.CheckTargetCallCount()).To(Equal(1))
			checkTargetedOrg, checkTargetedSpace := fakeSharedActor.CheckTargetArgsForCall(0)
			Expect(checkTargetedOrg).To(BeTrue())
			Expect(checkTargetedSpace).To(BeTrue())
		})
	})

	Context("when the user is logged in, and a space and org are targeted", func() {
		BeforeEach(func() {
			fakeConfig.HasTargetedOrganizationReturns(true)
			fakeConfig.TargetedOrganizationReturns(configv3.Organization{
				GUID: "some-org-guid",
				Name: "some-org",
			})
			fakeConfig.HasTargetedSpaceReturns(true)
			fakeConfig.TargetedSpaceReturns(configv3.Space{
				GUID: "some-space-guid",
				Name: "some-space",
			})
		})

		// Context("when getting the current user returns an error", func() {
		// 	var expectedErr error

		// 	BeforeEach(func() {
		// 		expectedErr = errors.New("get current user error")
		// 		fakeConfig.CurrentUserReturns(
		// 			configv3.User{},
		// 			expectedErr)
		// 	})

		// 	It("returns the error", func() {
		// 		Expect(executeErr).To(MatchError(expectedErr))
		// 	})
		// })

		Context("when getting the current user does not return an error", func() {
			BeforeEach(func() {
				fakeConfig.CurrentUserReturns(
					configv3.User{Name: "some-user"},
					nil)
			})

			Context("when provided a valid service instance", func() {
				Context("when the -o flag is not provided", func() {
					Context("when the share to space exists", func() {
						BeforeEach(func() {
							cmd.SpaceName.Space = "some-space"
							fakeActor.ShareServiceInstanceByOrganizationAndSpaceNameReturns(
								v3action.Warnings{"share-service-warning"},
								nil)
						})

						It("shares the service instance with the provided space and displays all warnings", func() {
							Expect(executeErr).ToNot(HaveOccurred())

							Expect(testUI.Out).To(Say("Sharing service instance some-service-instance into org some-org / space some-space as some-user\\.\\.\\."))
							Expect(testUI.Out).To(Say("OK"))
							Expect(testUI.Err).To(Say("share-service-warning"))

							Expect(fakeActor.ShareServiceInstanceByOrganizationAndSpaceNameCallCount()).To(Equal(1))
							serviceInstanceNameArg, orgGUIDArg, spaceNameArg := fakeActor.ShareServiceInstanceByOrganizationAndSpaceNameArgsForCall(0)
							Expect(serviceInstanceNameArg).To(Equal("some-service-instance"))
							Expect(orgGUIDArg).To(Equal("some-org-guid"))
							Expect(spaceNameArg).To(Equal("some-space"))
						})
					})
				})
			})
		})
	})
})

// cf share-service myservice some-other-space
// cf curl -X POST /v3/service_instances/c962853a-e329-4cb4-a578-a5cb878f5d16/relationships/shared_spaces \
// 	-d "{\"data\": [{\"guid\": \"4e25314e-c458-4a88-be25-0917e089b945\"}]}"
