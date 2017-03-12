package kurento

import "fmt"

type IImageOverlayFilter interface {
	RemoveImage(id string) error

	AddImage(id string, uri string, offsetXPercent float64, offsetYPercent float64, widthPercent float64, heightPercent float64, keepAspectRatio bool, center bool) error
}

/*ImageOverlayFilter interface. This type of `Filter` draws an image in a configured position over a video feed.*/
type ImageOverlayFilter struct {
	Filter
}

// Return Constructor Params to be called by "Create".
func (elem *ImageOverlayFilter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

/*Remove the image with the given ID.*/

func (elem *ImageOverlayFilter) RemoveImage(id string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "id", id)

	req["params"] = map[string]interface{}{
		"operation": "removeImage",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Add an image to be used as overlay.*/

func (elem *ImageOverlayFilter) AddImage(id string, uri string, offsetXPercent float64, offsetYPercent float64, widthPercent float64, heightPercent float64, keepAspectRatio bool, center bool) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "id", id)

	setIfNotEmpty(params, "uri", uri)

	setIfNotEmpty(params, "offsetXPercent", offsetXPercent)

	setIfNotEmpty(params, "offsetYPercent", offsetYPercent)

	setIfNotEmpty(params, "widthPercent", widthPercent)

	setIfNotEmpty(params, "heightPercent", heightPercent)

	setIfNotEmpty(params, "keepAspectRatio", keepAspectRatio)

	setIfNotEmpty(params, "center", center)

	req["params"] = map[string]interface{}{
		"operation": "addImage",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}
