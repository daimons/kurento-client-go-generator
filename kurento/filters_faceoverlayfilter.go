package kurento

import "fmt"

type IFaceOverlayFilter interface {
	UnsetOverlayedImage() error

	SetOverlayedImage(uri string, offsetXPercent float64, offsetYPercent float64, widthPercent float64, heightPercent float64) error
}

/*FaceOverlayFilter interface. This type of `Filter` detects faces in a video feed. The face is then overlaid with an image.*/
type FaceOverlayFilter struct {
	Filter
}

// Return Constructor Params to be called by "Create".
func (elem *FaceOverlayFilter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

/*Clear the image to be shown over each detected face. Stops overlaying the faces.*/

func (elem *FaceOverlayFilter) UnsetOverlayedImage() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "unsetOverlayedImage",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Sets the image to use as overlay on the detected faces.*/

func (elem *FaceOverlayFilter) SetOverlayedImage(uri string, offsetXPercent float64, offsetYPercent float64, widthPercent float64, heightPercent float64) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "uri", uri)

	setIfNotEmpty(params, "offsetXPercent", offsetXPercent)

	setIfNotEmpty(params, "offsetYPercent", offsetYPercent)

	setIfNotEmpty(params, "widthPercent", widthPercent)

	setIfNotEmpty(params, "heightPercent", heightPercent)

	req["params"] = map[string]interface{}{
		"operation": "setOverlayedImage",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}
