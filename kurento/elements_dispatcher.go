package kurento

import "fmt"

type IDispatcher interface {
	Connect(source HubPort, sink HubPort) error
}

/*A `Hub` that allows routing between arbitrary port pairs*/
type Dispatcher struct {
	Hub
}

// Return Constructor Params to be called by "Create".
func (elem *Dispatcher) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

/*Connects each corresponding :rom:enum:`MediaType` of the given source port with the sink port.*/

func (elem *Dispatcher) Connect(source HubPort, sink HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "source", source)

	setIfNotEmpty(params, "sink", sink)

	req["params"] = map[string]interface{}{
		"operation": "connect",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}
