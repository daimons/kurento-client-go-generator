package kurento

import "fmt"

/*<p>Base interface used to manage capabilities common to all Kurento elements. This includes both: `MediaElement` and `MediaPipeline`</p>       <h4>Properties</h4>       <ul>         <li><b>id</b>: unique identifier assigned to this <code>MediaObject</code> at instantiation time. `MediaPipeline` IDs are generated with a GUID followed by suffix <code>_kurento.MediaPipeline</code>. `MediaElement` IDs are also a GUID with suffix <code>_kurento.elemenType</code> and prefixed by parent's ID.           <blockquote>           <dl>             <dt><i>MediaPipeline ID example</i></dt>             <dd><code>907cac3a-809a-4bbe-a93e-ae7e944c5cae_kurento.MediaPipeline</code></dd>             <dt><i>MediaElement ID example</i></dt> <dd><code>907cac3a-809a-4bbe-a93e-ae7e944c5cae_kurento.MediaPipeline/403da25a-805b-4cf1-8c55-f190588e6c9b_kurento.WebRtcEndpoint</code></dd>           </dl>           </blockquote>         </li>         <li><b>name</b>: free text intended to provide a friendly name for this <code>MediaObject</code>. Its default value is the same as the ID.</li>         <li><b>tags</b>: key-value pairs intended for applications to associate metadata to this <code>MediaObject</code> instance.</li>       </ul>       <p>       <h4>Events</h4>       <ul>         <li>`ErrorEvent`: reports asynchronous error events. It is recommended to always subscribe a listener to this event, as regular error from the pipeline will be notified through it, instead of through an exception when invoking a method.</li>       </ul>*/
type MediaObject struct {
	connection *Connection

	/*`MediaPipeline` to which this <code>MediaObject</code> belongs. It returns itself when invoked for a pipeline object.*/
	MediaPipeline IMediaPipeline

	/*parent of this <code>MediaObject</code>. The parent of a `Hub` or a `MediaElement` is its `MediaPipeline`. A `MediaPipeline` has no parent, so this property will be null.*/
	Parent IMediaObject

	/*unique identifier of this <code>MediaObject</code>. It's a synthetic identifier composed by a GUID and <code>MediaObject</code> type. The ID is prefixed with the parent ID when the object has parent: <i>ID_parent/ID_media-object</i>.*/
	Id string

	/*@deprecated*/
	/*(Use children instead) children of this <code>MediaObject</code>.*/
	Childs []IMediaObject

	/*children of this <code>MediaObject</code>.*/
	Children []IMediaObject

	/*this <code>MediaObject</code>'s name. This is just a comodity to simplify developers' life debugging, it is not used internally for indexing nor idenfiying the objects. By default, it's the object's ID.*/
	Name string

	/*flag activating or deactivating sending the element's tags in fired events.*/
	SendTagsInEvents bool

	/*<code>MediaObject</code> creation time in seconds since Epoch.*/
	CreationTime Iint
}

// Return Constructor Params to be called by "Create".
func (elem *MediaObject) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Adds a new tag to this <code>MediaObject</code>. If the tag is already present, it changes the value.*/

func (elem *MediaObject) AddTag(key string, value string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	setIfNotEmpty(params, "value", value)

	req["params"] = map[string]interface{}{
		"operation": "addTag",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Removes an existing tag. Exists silently with no error if tag is not defined.*/

func (elem *MediaObject) RemoveTag(key string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	req["params"] = map[string]interface{}{
		"operation": "removeTag",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Returns the value of given tag, or MEDIA_OBJECT_TAG_KEY_NOT_FOUND if tag is not defined.*/

// Returns
/*The value associated to the given key.*/

func (elem *MediaObject) GetTag(key string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	req["params"] = map[string]interface{}{
		"operation": "getTag",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The value associated to the given key.*/

	return response.Result["value"], response.Error

}

/*Returns all tags attached to this <code>MediaObject</code>.*/

// Returns
/*An array containing all key-value pairs associated with this <code>MediaObject</code>.*/

func (elem *MediaObject) GetTags() ([]Tag, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getTags",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*An array containing all key-value pairs associated with this <code>MediaObject</code>.*/

	ret := []Tag{}
	return ret, response.Error

}

type IServerManager interface {
	GetKmd(moduleName string) (string, error)

	GetUsedMemory() (int64, error)
}

/*This is a standalone object for managing the MediaServer*/
type ServerManager struct {
	MediaObject

	/*Server information, version, modules, factories, etc*/
	Info *ServerInfo

	/*All the pipelines available in the server*/
	Pipelines []IMediaPipeline

	/*All active sessions in the server*/
	Sessions []string

	/*Metadata stored in the server*/
	Metadata string
}

// Return Constructor Params to be called by "Create".
func (elem *ServerManager) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Returns the kmd associated to a module*/

// Returns
/*The kmd file*/

func (elem *ServerManager) GetKmd(moduleName string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "moduleName", moduleName)

	req["params"] = map[string]interface{}{
		"operation": "getKmd",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The kmd file*/

	return response.Result["value"], response.Error

}

/*Returns the amount of memory that the server is using in KiB*/

// Returns
/*The amount of KiB of memory being used*/

func (elem *ServerManager) GetUsedMemory() (int64, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getUsedMemory",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The amount of KiB of memory being used*/

	ret := int64{}
	return ret, response.Error

}

type ISessionEndpoint interface {
}

/*All networked Endpoints that require to manage connection sessions with remote peers implement this interface.*/
type SessionEndpoint struct {
	Endpoint
}

// Return Constructor Params to be called by "Create".
func (elem *SessionEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHub interface {
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
}

/*A Hub is a routing `MediaObject`. It connects several `endpoints <Endpoint>` together*/
type Hub struct {
	MediaObject
}

// Return Constructor Params to be called by "Create".
func (elem *Hub) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Returns a string in dot (graphviz) format that represents the gstreamer elements inside the pipeline*/

// Returns
/*The dot graph*/

func (elem *Hub) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	req["params"] = map[string]interface{}{
		"operation": "getGstreamerDot",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The dot graph*/

	return response.Result["value"], response.Error

}

type IFilter interface {
}

/*Base interface for all filters. This is a certain type of `MediaElement`, that processes media injected through its sinks, and delivers the outcome through its sources.*/
type Filter struct {
	MediaElement
}

// Return Constructor Params to be called by "Create".
func (elem *Filter) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IEndpoint interface {
}

/*Base interface for all end points. An Endpoint is a `MediaElement`*/
/*that allow `KMS` to interchange media contents with external systems,*/
/*supporting different transport protocols and mechanisms, such as `RTP`,*/
/*`WebRTC`, `HTTP`, "file:/" URLs... An "Endpoint" may*/
/*contain both sources and sinks for different media types, to provide*/
/*bidirectional communication.*/
type Endpoint struct {
	MediaElement
}

// Return Constructor Params to be called by "Create".
func (elem *Endpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHubPort interface {
}

/*This `MediaElement` specifies a connection with a `Hub`*/
type HubPort struct {
	MediaElement
}

// Return Constructor Params to be called by "Create".
func (elem *HubPort) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"hub": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

type IPassThrough interface {
}

/*This `MediaElement` that just passes media through*/
type PassThrough struct {
	MediaElement
}

// Return Constructor Params to be called by "Create".
func (elem *PassThrough) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	mergeOptions(ret, options)
	return ret

}

type IUriEndpoint interface {
	Pause() error

	Stop() error
}

/*Interface for endpoints the require a URI to work. An example of this, would be a `PlayerEndpoint` whose URI property could be used to locate a file to stream*/
type UriEndpoint struct {
	Endpoint

	/*The uri for this endpoint.*/
	Uri string

	/*State of the endpoint*/
	State *UriEndpointState
}

// Return Constructor Params to be called by "Create".
func (elem *UriEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Pauses the feed*/

func (elem *UriEndpoint) Pause() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "pause",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Stops the feed*/

func (elem *UriEndpoint) Stop() error {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "stop",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

type IMediaPipeline interface {
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
}

/*A pipeline is a container for a collection of `MediaElements<MediaElement>` and `MediaMixers<MediaMixer>`. It offers the methods needed to control the creation and connection of elements inside a certain pipeline.*/
type MediaPipeline struct {
	MediaObject

	/*If statistics about pipeline latency are enabled for all mediaElements*/
	LatencyStats bool
}

// Return Constructor Params to be called by "Create".
func (elem *MediaPipeline) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Returns a string in dot (graphviz) format that represents the gstreamer elements inside the pipeline*/

// Returns
/*The dot graph*/

func (elem *MediaPipeline) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	req["params"] = map[string]interface{}{
		"operation": "getGstreamerDot",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The dot graph*/

	return response.Result["value"], response.Error

}

type ISdpEndpoint interface {
	GenerateOffer() (string, error)

	ProcessOffer(offer string) (string, error)

	ProcessAnswer(answer string) (string, error)

	GetLocalSessionDescriptor() (string, error)

	GetRemoteSessionDescriptor() (string, error)
}

/*This interface is implemented by Endpoints that require an SDP negotiation for the setup of a networked media session with remote peers. The API provides the following functionality:       <ul>         <li>Generate SDP offers.</li>         <li>Process SDP offers.</li>         <li>Configure SDP related params.</li>       </ul>*/
type SdpEndpoint struct {
	SessionEndpoint

	/*Maximum bandwidth for video reception, in kbps. The default value is 500. A value of 0 sets this as unconstrained. .. note:: This has to be set before the SDP is generated.*/
	MaxVideoRecvBandwidth Iint

	/*Maximum bandwidth for audio reception, in kbps. The default value is 500. A value of 0 sets this as leaves this unconstrained. .. note:: This has to be set before the SDP is generated.*/
	MaxAudioRecvBandwidth Iint
}

// Return Constructor Params to be called by "Create".
func (elem *SdpEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Generates an SDP offer with  media capabilities of the Endpoint.           Exceptions           <ul>             <li>               SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.             </li>             <li>               SDP_END_POINT_GENERATE_OFFER_ERROR if the generated offer is empty. This is most likely due to an internal error.             </li>           </ul>*/

// Returns
/*The SDP offer.*/

func (elem *SdpEndpoint) GenerateOffer() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "generateOffer",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The SDP offer.*/

	return response.Result["value"], response.Error

}

/*Processes SDP offer of the remote peer, and generates an SDP answer based on the endpoint's capabilities. If no matching capabilities are found, the SDP will contain no codecs.           Exceptions           <ul>             <li>               SDP_PARSE_ERROR If the offer is empty or has errors.             </li>             <li>               SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.             </li>             <li>               SDP_END_POINT_PROCESS_OFFER_ERROR if the generated offer is empty. This is most likely due to an internal error.             </li>           </ul>*/

// Returns
/*The chosen configuration from the ones stated in the SDP offer*/

func (elem *SdpEndpoint) ProcessOffer(offer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "offer", offer)

	req["params"] = map[string]interface{}{
		"operation": "processOffer",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The chosen configuration from the ones stated in the SDP offer*/

	return response.Result["value"], response.Error

}

/*Generates an SDP offer with  media capabilities of the Endpoint.           Exceptions           <ul>             <li>               SDP_PARSE_ERROR If the offer is empty or has errors.             </li>             <li>               SDP_END_POINT_ALREADY_NEGOTIATED If the endpoint is already negotiated.             </li>             <li>               SDP_END_POINT_PROCESS_ANSWER_ERROR if the result of processing the answer is an empty string. This is most likely due to an internal error.             </li>             <li>               SDP_END_POINT_NOT_OFFER_GENERATED If the method is invoked before the generateOffer method.             </li>           </ul>*/

// Returns
/*Updated SDP offer, based on the answer received.*/

func (elem *SdpEndpoint) ProcessAnswer(answer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "answer", answer)

	req["params"] = map[string]interface{}{
		"operation": "processAnswer",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*Updated SDP offer, based on the answer received.*/

	return response.Result["value"], response.Error

}

/*This method returns the local SDP. The output depends on the negotiation stage:           <ul>             <li>               No offer has been generated: returns null.             </li>             <li>               Offer has been generated: return the SDP offer.             </li>             <li>               Offer has been generated and answer processed: retruns the agreed SDP.             </li>           </ul>*/

// Returns
/*The last agreed SessionSpec*/

func (elem *SdpEndpoint) GetLocalSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getLocalSessionDescriptor",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The last agreed SessionSpec*/

	return response.Result["value"], response.Error

}

/*This method returns the remote SDP. If the negotiation process is not complete, it will return NULL.*/

// Returns
/*The last agreed User Agent session description*/

func (elem *SdpEndpoint) GetRemoteSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	req["params"] = map[string]interface{}{
		"operation": "getRemoteSessionDescriptor",
		"object":    elem.Id,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The last agreed User Agent session description*/

	return response.Result["value"], response.Error

}

type IBaseRtpEndpoint interface {
}

/*This class extends from the SdpEndpoint, and handles RTP communications. All endpoints that rely on this network protocol, like the RTPEndpoint or the WebRtcEndpoint, inherit from this. The endpoint provides information about the connection state and the media state. These can be consulted at any time through the mediaState and the connectionState properties. It is also possible subscribe to events fired when these properties change.       <ul style='list-style-type:circle'>         <li>           ConnectionStateChangedEvent: This event is raised when the connection between two peers changes. It can have two values           <ul>             <li>CONNECTED</li>             <li>DISCONNECTED</li>           </ul>         </li>         <li>           MediaStateChangedEvent: Based on RTCP packet flow, this event provides more reliable information about the state of media flow. Since RTCP packets are not flowing at a constant rate (minimizing a browser with an RTCPeerConnection might affect this interval, for instance), there is a guard period of about 5s. This traduces in a period where there might be no media flowing, but the event hasn't been fired yet. Nevertheless, this is the most reliable and useful way of knowing what the state of media exchange is. Possible values are:           <ul>             <li>CONNECTED: There is an RTCP packet flow between peers.</li>             <li>DISCONNECTED: No RTCP packets have been received, or at least 5s have passed since the last packet arrived.</li>           </ul>         </li>       </ul>       Part of the bandwidth control of the video component of the media session is done here. The values of the properties described are in kbps.       <ul style='list-style-type:circle'>         <li>           Input bandwidth control mechanism: Configuration interval used to inform remote peer the range of bitrates that can be pushed into this BaseRtpEndpoint object.           <ul>             <li>               setMinVideoRecvBandwidth: sets min bitrate limits expected for the received video stream. This value is set to limit the lower value of REMB packages, if supported by the implementing class.             </li>           </ul>           Max values are announced in the SDP, while min values are set to limit the lower value of REMB packages. It follows that min values will only have effect in peers that support this control mechanism, such as Chrome.         </li>         <li>           Output bandwidth control mechanism: Configuration interval used to control bitrate of the output video stream sent to remote peer. It is important to keep in mind that pushed bitrate depends on network and remote peer capabilities. Remote peers can also announce bandwidth limitation in their SDPs (through the b=<modifier>:<value> tag).   Kurento will always enforce bitrate limitations specified by the remote peer over internal configurations.           <ul>             <li>               setMinVideoSendBandwidth: sets the minimum bitrate for video to be sent to remote peer. 0 is considered unconstrained.             </li>             <li>               setMaxVideoSendBandwidth: sets maximum bitrate limits for video sent to remote peer. 0 is considered unconstrained.             </li>           </ul>         </li>       </ul>       All bandwidth control parameters must be changed before the SDP negotiation takes place, and can't be changed afterwards.       </p>*/
type BaseRtpEndpoint struct {
	SdpEndpoint

	/*Minimum bandwidth announced for video reception, in kbps. The default and absolute minimum value is 30 kbps, even if a lower value is set.*/
	MinVideoRecvBandwidth Iint

	/*Minimum bandwidth for video transmission, in kbps. The default value is 100 kbps. 0 is considered unconstrained.*/
	MinVideoSendBandwidth Iint

	/*Maximum bandwidth for video transmission, in kbps. The default value is 500 kbps. 0 is considered unconstrained.*/
	MaxVideoSendBandwidth Iint

	/*Media flow state. Possible values are           <ul>             <li>CONNECTED: There is an RTCP flow.</li>             <li>DISCONNECTED: No RTCP packets have been received for at least 5 sec.</li>           </ul>*/
	MediaState *MediaState

	/*Connection state. Possible values are           <ul>             <li>CONNECTED</li>             <li>DISCONNECTED</li>           </ul>*/
	ConnectionState *ConnectionState

	/*Advanced parameters to configure the congestion control algorithm.*/
	RembParams *RembParams
}

// Return Constructor Params to be called by "Create".
func (elem *BaseRtpEndpoint) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IMediaElement interface {
	GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)

	GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)

	Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error

	Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error

	SetAudioFormat(caps AudioCaps) error

	SetVideoFormat(caps VideoCaps) error

	GetGstreamerDot(details GstreamerDotDetails) (string, error)

	SetOutputBitrate(bitrate int) error

	GetStats(mediaType MediaType) (Stats, error)

	IsMediaFlowingIn(mediaType MediaType, sinkMediaDescription string) (bool, error)

	IsMediaFlowingOut(mediaType MediaType, sourceMediaDescription string) (bool, error)
}

/*<p>This is the basic building block of the media server, that can be interconnected inside a pipeline. A `MediaElement` is a module that encapsulates a specific media capability, and that is able to exchange media with other `MediaElement`s through an internal element called pad.       </p>       <p>       A pad can be defined as an input or output interface. Input pads are called sinks, and it's where the media elements receive media from other media elements. Output interfaces are called sources, and it's the pad used by the media element to feed media to other media elements. There can be only one sink pad per media element. On the other hand, the number of source pads is unconstrained. This means that a certain media element can receive media only from one element at a time, while it can send media to many others. Pads are created on demand, when the connect method is invoked. When two media elements are connected, one media pad is created for each type of media connected. For example, if you connect AUDIO and VIDEO between two media elements, each one will need to create two new pads: one for AUDIO and one for VIDEO.       </p>       <p>       When media elements are connected, it can be case that the encoding used by the elements is not the same, and thus it needs to be transcoded. This is something that is handled transparently by the media elements internals. In practice, the user needs not be aware that the transcodification is taking place. However, this process has a toll in the form of a higher CPU load, so connecting media elements that need media encoded in different formats is something to consider as a high load operation.       </p>*/
type MediaElement struct {
	MediaObject

	/*@deprecated*/
	/*Deprecated due to a typo. Use minOutputBitrate instead of this function. Minimum video bandwidth for transcoding.*/
	/*Unit: bps(bits per second).*/
	/*Default value: 0*/
	MinOuputBitrate Iint

	/*Minimum video bitrate for transcoding.*/
	/*Unit: bps(bits per second).*/
	/*Default value: 0*/
	MinOutputBitrate Iint

	/*@deprecated*/
	/*Deprecated due to a typo. Use maxOutputBitrate instead of this function. Maximum video bandwidth for transcoding. 0 = unlimited.*/
	/*Unit: bps(bits per second).*/
	/*Default value: MAXINT*/
	MaxOuputBitrate Iint

	/*Maximum video bitrate for transcoding. 0 = unlimited.*/
	/*Unit: bps(bits per second).*/
	/*Default value: MAXINT*/
	MaxOutputBitrate Iint
}

// Return Constructor Params to be called by "Create".
func (elem *MediaElement) getConstuctorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

/*Gets information about the sink pads of this media element. Since sink pads are the interface through which a media element gets it's media, whatever is connected to an element's sink pad is formally a source of media. Media can be filtered by type, or by the description given to the pad though which both elements are connected.*/

// Returns
/*A list of the connections information that are sending media to this element. The list will be empty if no sources are found.*/

func (elem *MediaElement) GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "description", description)

	req["params"] = map[string]interface{}{
		"operation": "getSourceConnections",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*A list of the connections information that are sending media to this element. The list will be empty if no sources are found.*/

	ret := []ElementConnectionData{}
	return ret, response.Error

}

/*Gets information about the source pads of this media element. Since source pads connect to other media element's sinks, this is formally the sink of media from the element's perspective. Media can be filtered by type, or by the description given to the pad though which both elements are connected.*/

// Returns
/*A list of the connections information that are receiving media from this element. The list will be empty if no sources are found.*/

func (elem *MediaElement) GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "description", description)

	req["params"] = map[string]interface{}{
		"operation": "getSinkConnections",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*A list of the connections information that are receiving media from this element. The list will be empty if no sources are found.*/

	ret := []ElementConnectionData{}
	return ret, response.Error

}

/*<p>Connects two elements, with the media flowing from left to right: the elements that invokes the connect wil be the source of media, creating one sink pad for each type of media connected. The element given as parameter to the method will be the sink, and it will create one sink pad per media type connected.           </p>           <p>           If otherwise not specified, all types of media are connected by default (AUDIO, VIDEO and DATA). It is recommended to connect the specific types of media if not all of them will be used. For this purpose, the connect method can be invoked more than once on the same two elements, but with different media types.           </p>           <p>           The connection is unidirectional. If a bidirectional connection is desired, the position of the media elements must be inverted. For instance, webrtc1.connect(webrtc2) is connecting webrtc1 as source of webrtc2. In order to create a WebRTC one-2one conversation, the user would need to especify the connection on the other direction with webrtc2.connect(webrtc1).           </p>           <p>           Even though one media element can have one sink pad per type of media, only one media element can be connected to another at a given time. If a media element is connected to another, the former will become the source of the sink media element, regardles whether there was another element connected or not.           </p>*/

func (elem *MediaElement) Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)

	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	req["params"] = map[string]interface{}{
		"operation": "connect",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Disconnectes two media elements. This will release the source pads of the source media element, and the sink pads of the sink media element.*/

func (elem *MediaElement) Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)

	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	req["params"] = map[string]interface{}{
		"operation": "disconnect",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Sets the type of data for the audio stream. MediaElements that do not support configuration of audio capabilities will throw a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR exception.*/

func (elem *MediaElement) SetAudioFormat(caps AudioCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	req["params"] = map[string]interface{}{
		"operation": "setAudioFormat",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Sets the type of data for the video stream. MediaElements that do not support configuration of video capabilities will throw a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR exception*/

func (elem *MediaElement) SetVideoFormat(caps VideoCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	req["params"] = map[string]interface{}{
		"operation": "setVideoFormat",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*This method returns a .dot file describing the topology of the media element. The element can be queried for certain type of data           <ul>             <li>SHOW_ALL: default value</li>             <li>SHOW_CAPS_DETAILS</li>             <li>SHOW_FULL_PARAMS</li>             <li>SHOW_MEDIA_TYPE</li>             <li>SHOW_NON_DEFAULT_PARAMS</li>             <li>SHOW_STATES</li>             <li>SHOW_VERBOSE</li>           </ul>*/

// Returns
/*The dot graph*/

func (elem *MediaElement) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	req["params"] = map[string]interface{}{
		"operation": "getGstreamerDot",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*The dot graph*/

	return response.Result["value"], response.Error

}

/*@deprecated*/
/*Allows change the target bitrate for the media output, if the media is encoded using VP8 or H264. This method only works if it is called before the media starts to flow.*/

func (elem *MediaElement) SetOutputBitrate(bitrate int) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "bitrate", bitrate)

	req["params"] = map[string]interface{}{
		"operation": "setOutputBitrate",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	return response.Error

}

/*Gets the statistics related to an endpoint. If no media type is specified, it returns statistics for all available types.*/

// Returns
/*Delivers a successful result in the form of a RTC stats report. A RTC stats report represents a map between strings, identifying the inspected objects (RTCStats.id), and their corresponding RTCStats objects.*/

func (elem *MediaElement) GetStats(mediaType MediaType) (Stats, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	req["params"] = map[string]interface{}{
		"operation": "getStats",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*Delivers a successful result in the form of a RTC stats report. A RTC stats report represents a map between strings, identifying the inspected objects (RTCStats.id), and their corresponding RTCStats objects.*/

	ret := Stats{}
	return ret, response.Error

}

/*This method indicates whether the media element is receiving media of a certain type. The media sink pad can be identified individually, if needed. It is only supported for AUDIO and VIDEO types, raising a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR otherwise. If the pad indicated does not exist, if will return false.*/

// Returns
/*TRUE if there is media, FALSE in other case*/

func (elem *MediaElement) IsMediaFlowingIn(mediaType MediaType, sinkMediaDescription string) (bool, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	req["params"] = map[string]interface{}{
		"operation": "isMediaFlowingIn",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*TRUE if there is media, FALSE in other case*/

	ret := bool{}
	return ret, response.Error

}

/*This method indicates whether the media element is emitting media of a certain type. The media source pad can be identified individually, if needed. It is only supported for AUDIO and VIDEO types, raising a MEDIA_OBJECT_ILLEGAL_PARAM_ERROR otherwise. If the pad indicated does not exist, if will return false.*/

// Returns
/*TRUE if there is media, FALSE in other case*/

func (elem *MediaElement) IsMediaFlowingOut(mediaType MediaType, sourceMediaDescription string) (bool, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)

	req["params"] = map[string]interface{}{
		"operation": "isMediaFlowingOut",
		"object":    elem.Id,

		"operationParams": params,
	}

	// call server and and wait response
	response := <-elem.connection.Request(req)

	/*TRUE if there is media, FALSE in other case*/

	ret := bool{}
	return ret, response.Error

}
