package kurento

import "fmt"

type IHttpPostEndpoint interface {
}

/*An `HttpPostEndpoint` contains SINK pads for AUDIO and VIDEO, which provide access to an HTTP file upload function*/
/**/
/*This type of endpoint provide unidirectional communications. Its `MediaSources <MediaSource>` are accessed through the `HTTP` POST method.*/
type HttpPostEndpoint struct {
	HttpEndpoint
}

// Return Constructor Params to be called by "Create".
func (elem *HttpPostEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":        fmt.Sprintf("%s", from),
		"disconnectionTimeout": 2,
		"useEncodedMedia":      fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

type IHttpEndpoint interface {
	GetUrl() (string, error)
}

/*Endpoint that enables Kurento to work as an HTTP server, allowing peer HTTP clients to access media.*/
type HttpEndpoint struct {
	SessionEndpoint
}

// Return Constructor Params to be called by "Create".
func (elem *HttpEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Obtains the URL associated to this endpoint*/

// Returns
/*The url as a String*/

func (elem *HttpEndpoint) GetUrl() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getUrl",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The url as a String*/

	return response.Result["value"], response.Error

}
