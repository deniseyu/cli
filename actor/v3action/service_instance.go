package v3action

import (
	"fmt"
	"net/url"

	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
)

type ServiceInstance ccv3.ServiceInstance

func (actor Actor) ShareServiceInstanceByOrganizationAndSpaceName(serviceInstanceName string, orgGUID string, spaceName string) (Warnings, error) {
	// get the service instnace guid
	serviceInstance, _, err := actor.GetServiceInstanceByName(serviceInstanceName)
	if err != nil {
		fmt.Println(err)
	}

	space, _, err := actor.GetSpaceByName(spaceName)
	if err != nil {
		fmt.Println(err)
	}

	//Think about name of this.
	relationship, _, err := actor.CloudControllerClient.PostServiceInstanceSharedSpaces(serviceInstance.GUID, []string{space.GUID})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(relationship)

	return nil, nil
}

func (actor Actor) GetServiceInstanceByName(serviceInstanceName string) (ServiceInstance, Warnings, error) {
	serviceInstances, warnings, err := actor.CloudControllerClient.GetServiceInstances(url.Values{
		ccv3.NameFilter: []string{serviceInstanceName},
	})

	if err != nil {
		return ServiceInstance{}, Warnings(warnings), err
	}

	if len(serviceInstances) == 0 {
		return ServiceInstance{}, Warnings(warnings), actionerror.ServiceInstanceNotFoundError{Name: serviceInstanceName}
	}

	//Handle multiple serviceInstances being returned as GetServiceInstances arnt filtered by space
	return ServiceInstance(serviceInstances[0]), Warnings(warnings), nil
}
