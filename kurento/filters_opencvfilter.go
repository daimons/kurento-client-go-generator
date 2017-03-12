package kurento

type IOpenCVFilter interface {
}

/*Generic OpenCV Filter*/
type OpenCVFilter struct {
	Filter
}

// Return Constructor Params to be called by "Create".
func (elem *OpenCVFilter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}
