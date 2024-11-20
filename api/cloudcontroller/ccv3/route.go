package ccv3

import (
	ccv3internal "code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/internal"
	"code.cloudfoundry.org/cli/api/internal"
	"code.cloudfoundry.org/cli/resources"
)

func (client Client) CreateRoute(route resources.Route) (resources.Route, Warnings, error) {
	var responseBody resources.Route

	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.PostRouteRequest,
		RequestBody:  route,
		ResponseBody: &responseBody,
	})

	return responseBody, warnings, err
}

func (client Client) DeleteOrphanedRoutes(spaceGUID string) (JobURL, Warnings, error) {
	jobURL, warnings, err := client.MakeRequest(RequestParams{
		RequestName: ccv3internal.DeleteOrphanedRoutesRequest,
		URIParams:   internal.Params{"space_guid": spaceGUID},
		Query:       []Query{{Key: UnmappedFilter, Values: []string{"true"}}},
	})

	return jobURL, warnings, err
}

func (client Client) DeleteRoute(routeGUID string) (JobURL, Warnings, error) {
	jobURL, warnings, err := client.MakeRequest(RequestParams{
		RequestName: ccv3internal.DeleteRouteRequest,
		URIParams:   internal.Params{"route_guid": routeGUID},
	})

	return jobURL, warnings, err
}

func (client Client) GetApplicationRoutes(appGUID string) ([]resources.Route, Warnings, error) {
	var routes []resources.Route

	_, warnings, err := client.MakeListRequest(RequestParams{
		RequestName:  ccv3internal.GetApplicationRoutesRequest,
		URIParams:    internal.Params{"app_guid": appGUID},
		ResponseBody: resources.Route{},
		AppendToList: func(item interface{}) error {
			routes = append(routes, item.(resources.Route))
			return nil
		},
	})

	return routes, warnings, err
}

func (client Client) GetRouteDestinations(routeGUID string) ([]resources.RouteDestination, Warnings, error) {
	var responseBody struct {
		Destinations []resources.RouteDestination `json:"destinations"`
	}

	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.GetRouteDestinationsRequest,
		URIParams:    internal.Params{"route_guid": routeGUID},
		ResponseBody: &responseBody,
	})

	return responseBody.Destinations, warnings, err
}

func (client Client) GetRoutes(query ...Query) ([]resources.Route, Warnings, error) {
	var routes []resources.Route

	_, warnings, err := client.MakeListRequest(RequestParams{
		RequestName:  ccv3internal.GetRoutesRequest,
		Query:        query,
		ResponseBody: resources.Route{},
		AppendToList: func(item interface{}) error {
			routes = append(routes, item.(resources.Route))
			return nil
		},
	})

	return routes, warnings, err
}

func (client Client) MapRoute(routeGUID string, appGUID string, destinationProtocol string) (Warnings, error) {
	type destinationProcess struct {
		ProcessType string `json:"process_type"`
	}

	type destinationApp struct {
		GUID    string              `json:"guid"`
		Process *destinationProcess `json:"process,omitempty"`
	}
	type destination struct {
		App      destinationApp `json:"app"`
		Protocol string         `json:"protocol,omitempty"`
	}

	type body struct {
		Destinations []destination `json:"destinations"`
	}

	requestBody := body{
		Destinations: []destination{
			{
				App: destinationApp{GUID: appGUID},
			},
		},
	}
	if destinationProtocol != "" {
		requestBody.Destinations[0].Protocol = destinationProtocol
	}

	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName: ccv3internal.MapRouteRequest,
		URIParams:   internal.Params{"route_guid": routeGUID},
		RequestBody: &requestBody,
	})

	return warnings, err
}

func (client Client) UnmapRoute(routeGUID string, destinationGUID string) (Warnings, error) {
	var responseBody resources.Build

	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.UnmapRouteRequest,
		URIParams:    internal.Params{"route_guid": routeGUID, "destination_guid": destinationGUID},
		ResponseBody: &responseBody,
	})

	return warnings, err
}

func (client Client) UnshareRoute(routeGUID string, spaceGUID string) (Warnings, error) {
	var responseBody resources.Build

	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.UnshareRouteRequest,
		URIParams:    internal.Params{"route_guid": routeGUID, "space_guid": spaceGUID},
		ResponseBody: &responseBody,
	})
	return warnings, err
}

func (client Client) UpdateDestination(routeGUID string, destinationGUID string, protocol string) (Warnings, error) {
	type body struct {
		Protocol string `json:"protocol"`
	}
	requestBody := body{
		Protocol: protocol,
	}
	var responseBody resources.Build
	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.PatchDestinationRequest,
		URIParams:    internal.Params{"route_guid": routeGUID, "destination_guid": destinationGUID},
		RequestBody:  &requestBody,
		ResponseBody: &responseBody,
	})
	return warnings, err
}

func (client Client) ShareRoute(routeGUID string, spaceGUID string) (Warnings, error) {
	type space struct {
		GUID string `json:"guid"`
	}

	type body struct {
		Data []space `json:"data"`
	}

	requestBody := body{
		Data: []space{
			{GUID: spaceGUID},
		},
	}

	var responseBody resources.Build
	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.ShareRouteRequest,
		URIParams:    internal.Params{"route_guid": routeGUID},
		RequestBody:  &requestBody,
		ResponseBody: &responseBody,
	})
	return warnings, err
}

func (client Client) MoveRoute(routeGUID string, spaceGUID string) (Warnings, error) {
	type space struct {
		GUID string `json:"guid"`
	}

	type body struct {
		Data space `json:"data"`
	}

	requestBody := body{
		Data: space{
			GUID: spaceGUID,
		},
	}

	var responseBody resources.Build
	_, warnings, err := client.MakeRequest(RequestParams{
		RequestName:  ccv3internal.PatchMoveRouteRequest,
		URIParams:    internal.Params{"route_guid": routeGUID},
		RequestBody:  &requestBody,
		ResponseBody: &responseBody,
	})
	return warnings, err
}
