package ccv3

type ServiceInstance struct {
	GUID string
	Name string
}

func (client Client) GetServiceInstanceByName(serviceInstanceName string) (ServiceInstance, Warnings, error) {
	return ServiceInstance{}, nil, nil
}
