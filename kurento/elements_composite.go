package kurento

import "fmt"

type IComposite interface {
}

/*A `Hub` that mixes the :rom:attr:`MediaType.AUDIO` stream of its connected sources and constructs a grid with the :rom:attr:`MediaType.VIDEO` streams of its connected sources into its sink*/
type Composite struct {
	Hub
}

// Return Constructor Params to be called by "Create".
func (elem *Composite) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}
