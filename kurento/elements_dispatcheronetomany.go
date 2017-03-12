package kurento

import "fmt"

type IDispatcherOneToMany interface {
	SetSource(source HubPort) error

	RemoveSource() error
}

/*A `Hub` that sends a given source to all the connected sinks*/
type DispatcherOneToMany struct {
	Hub
}

// Return Constructor Params to be called by "Create".
func (elem *DispatcherOneToMany) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

/*Sets the source port that will be connected to the sinks of every `HubPort` of the dispatcher*/

func (elem *DispatcherOneToMany) SetSource(source HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "source", source)

	req["params"] = map[string]interface{}{
		"operation": "setSource",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Remove the source port and stop the media pipeline.*/

func (elem *DispatcherOneToMany) RemoveSource() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "removeSource",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}
