package kurento

import "fmt"

type IPlayerEndpoint interface {
	Play() error
}

/*<p>       Retrieves content from seekable or non-seekable sources, and injects them into `KMS`, so they can be delivered to any Filter or Endpoint in the same MediaPipeline. Following URI schemas are supported:       <ul>         <li>           Files: Mounted in the local file system.           <ul><li>file:///path/to/file</li></ul>         </li>         <li>           RTSP: Those of IP cameras would be a good example.           <ul>             <li>rtsp://<server-ip></li>             <li>rtsp://username:password@<server-ip></li>           </ul>         </li>         <li>           HTTP: Any file available in an HTTP server           <ul>             <li>http(s)://<server-ip>/path/to/file</li>             <li>http(s)://username:password@<server-ip>/path/to/file</li>           </ul>         </li>       </ul>       </p>       <p>       For the player to stream the contents of the file, the server must have access to the resource. In case of local files, the user running the process must have read permissions over the file. For network resources, the path to the resource must be accessible: IP and port access not blocked, correct credentials, etc.The resource location can’t be changed after the player is created, and a new player should be created for streaming a different resource.       </p>       <p>       The list of valid operations is       <ul>         <li>*play*: starts streaming media. If invoked after pause, it will resume playback.</li>         <li>*stop*: stops streaming media. If play is invoked afterwards, the file will be streamed from the beginning.</li>         <li>*pause*: pauses media streaming. Play must be invoked in order to resume playback.</li>         <li>*seek*: If the source supports “jumps” in the timeline, then the PlayerEndpoint can           <ul>             <li>*setPosition*: allows to set the position in the file.</li>             <li>*getPosition*: returns the current position being streamed.</li>           </ul>         </li>       </ul>       </p>       <p>       <h2>Events fired:</h2>       <ul><li>EndOfStreamEvent: If the file is streamed completely.</li></ul>       </p>*/
type PlayerEndpoint struct {
	UriEndpoint

	/*Returns info about the source being played*/
	VideoInfo *VideoInfo

	/*Get or set the actual position of the video in ms. .. note:: Setting the position only works for seekable videos*/
	Position Iint64
}

// Return Constructor Params to be called by "Create".
func (elem *PlayerEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":   fmt.Sprintf("%s", from),
		"uri":             "",
		"useEncodedMedia": fmt.Sprintf("%s", from),
		"networkCache":    2000,
	}

	mergeOptions(ret, options)
	return ret

}

/*Starts reproducing the media, sending it to the `MediaSource`. If the endpoint*/
/*has been connected to other endpoints, those will start receiving media.*/

func (elem *PlayerEndpoint) Play() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "play",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}
