package kurento

import "fmt"

type IGStreamerFilter interface {
}

/*This is a generic filter interface, that creates GStreamer filters in the media server.*/
type GStreamerFilter struct {
	Filter

	/*GStreamer command.*/
	Command string
}

// Return Constructor Params to be called by "Create".
func (elem *GStreamerFilter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
		"command":       "",
		"filterType":    fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}
