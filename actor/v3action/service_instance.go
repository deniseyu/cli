package v3action

func (actor Actor) ShareServiceInstanceByOrganizationAndSpaceName(serviceInstanceName string, orgGUID string, spaceName string) (Warnings, error) {
	// get the service instnace guid
	serviceInstance, _, _ := actor.CloudControllerClient.GetServiceInstanceByName(serviceInstanceName)

	// get the space guid

	// share the space to the space guid
	return nil, nil
}
