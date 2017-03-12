package kurento

import "fmt"

type IZBarFilter interface {
}

/*This filter detects `QR` codes in a video feed. When a code is found, the filter raises a :rom:evnt:`CodeFound` event.*/
type ZBarFilter struct {
	Filter
}

// Return Constructor Params to be called by "Create".
func (elem *ZBarFilter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}
