package v3

import (
	"net/http"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccversion"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/translatableerror"
	"code.cloudfoundry.org/cli/command/v3/shared"
)

//go:generate counterfeiter . ShareServiceActor

type ShareServiceActor interface {
	ShareServiceInstanceByOrganizationAndSpaceName(serviceInstanceName string, orgGUID string, spaceName string) (v3action.Warnings, error)
}

type V3ShareServiceCommand struct {
	RequiredArgs flag.ServiceInstance `positional-args:"yes"`
	//TODO flag.Space does not capture the command line value
	OrgName         string      `short:"o" required:"false" description:"Org of the other space (Default: targeted org)"`
	SpaceName       string      `short:"s" required:"true" description:"Space to share the service instance into"`
	usage           interface{} `usage:"cf v3-share-service SERVICE_INSTANCE -s OTHER_SPACE [-o OTHER_ORG]"`
	relatedCommands interface{} `related_commands:"bind-service, service, services"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       ShareServiceActor
}

func (cmd *V3ShareServiceCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor(config)

	client, _, err := shared.NewClients(config, ui, true)
	if err != nil {
		if v3Err, ok := err.(ccerror.V3UnexpectedResponseError); ok && v3Err.ResponseCode == http.StatusNotFound {
			return translatableerror.MinimumAPIVersionNotMetError{MinimumVersion: ccversion.MinVersionRunTaskV3}
		}
		return err
	}
	cmd.Actor = v3action.NewActor(client, config, nil, nil)

	return nil
}

func (cmd V3ShareServiceCommand) Execute(args []string) error {
	cmd.UI.DisplayText(command.ExperimentalWarning)
	cmd.UI.DisplayNewline()

	err := cmd.SharedActor.CheckTarget(true, true)
	if err != nil {
		return err
	}

	user, _ := cmd.Config.CurrentUser()
	// if err != nil {
	// 	return err
	// }

	cmd.UI.DisplayTextWithFlavor("Sharing service instance {{.ServiceInstanceName}} into org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...", map[string]interface{}{
		"ServiceInstanceName": cmd.RequiredArgs.ServiceInstance,
		"OrgName":             cmd.Config.TargetedOrganization().Name,
		"SpaceName":           cmd.SpaceName,
		"Username":            user.Name,
	})

	warnings, err := cmd.Actor.ShareServiceInstanceByOrganizationAndSpaceName(cmd.RequiredArgs.ServiceInstance, cmd.Config.TargetedOrganization().GUID, cmd.SpaceName)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return err
	}

	cmd.UI.DisplayOK()

	return nil
}
