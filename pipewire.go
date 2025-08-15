package pwmonitor

import (
	"fmt"
	"time"

	json_v2 "github.com/go-json-experiment/json"
)

type EventType string

const (
	EmptyEvent EventType = ""
	EventNode  EventType = "PipeWire:Interface:Node"
)

type (
	Event struct {
		ID          int        `json:"id"`
		Type        EventType  `json:"type"`
		Version     int        `json:"version"`
		Info        *EventInfo `json:"info"`
		Permissions []string   `json:"permissions"`
		// When the event was received
		CapturedAt time.Time `json:"-"`
	}

	EventInfo struct {
		Direction  string          `json:"direction,omitempty"`
		ChangeMask []string        `json:"change-mask"`
		Props      *EventInfoProps `json:"props,omitempty"`
		Params     *EventParams    `json:"params,omitempty"`
		State      *State          `json:"state,omitempty"`
		Error      *any            `json:"error,omitempty"`
	}

	EventParams struct {
		EnumFormat []ParamEnumFormat `json:"EnumFormat,omitempty"`
		Meta       []ParamMeta       `json:"Meta,omitempty"`
		IO         []ParamIO         `json:"IO,omitempty"`
		Format     []any             `json:"Format,omitempty"`
		Buffers    []any             `json:"Buffers,omitempty"`
		Latency    []ParamLatency    `json:"Latency,omitempty"`
		Tag        []any             `json:"Tag,omitempty"`
	}

	EventInfoProps struct {
		AdaptFollowerSpaNode *string `json:"adapt.follower.spa-node,omitempty"`
		ApplicationIconName  *string `json:"application.icon-name,omitempty"`
		ApplicationID        *string `json:"application.id,omitempty"`
		ApplicationName      *string `json:"application.name,omitempty"`
		ClientID             *int    `json:"client.id,omitempty"`
		ClockQuantumLimit    *int    `json:"clock.quantum-limit,omitempty"`
		FactoryID            *int    `json:"factory.id,omitempty"`
		LibraryName          *string `json:"library.name,omitempty"`
		MediaCategory        *string `json:"media.category,omitempty"`
		MediaClass           *string `json:"media.class,omitempty"`
		MediaName            *string `json:"media.name,omitempty"`
		MediaRole            *string `json:"media.role,omitempty"`
		MediaType            *string `json:"media.type,omitempty"`
		NodeAlwaysProcess    *bool   `json:"node.always-process,omitempty"`
		NodeAutoconnect      *bool   `json:"node.autoconnect,omitempty"`
		NodeDescription      *string `json:"node.description,omitempty"`
		NodeLoopName         *string `json:"node.loop.name,omitempty"`
		NodeName             *string `json:"node.name,omitempty"`
		NodeRate             *string `json:"node.rate,omitempty"`
		NodeWantDriver       *bool   `json:"node.want-driver,omitempty"`
		ObjectID             *int    `json:"object.id,omitempty"`
		ObjectRegister       *bool   `json:"object.register,omitempty"`
		ObjectSerial         *int    `json:"object.serial,omitempty"`
		PortGroup            *string `json:"port.group,omitempty"`
		StreamIsLive         *bool   `json:"stream.is-live,omitempty"`
	}

	ParamEnumFormat struct {
		MediaType    string `json:"mediaType"`
		MediaSubtype string `json:"mediaSubtype"`
		Format       any    `json:"format"`
	}

	ParamMeta struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	ParamIO struct {
		ID   string `json:"id"`
		Size int    `json:"size"`
	}

	ParamLatency struct {
		Direction  string  `json:"direction"`
		MinQuantum float64 `json:"minQuantum"`
		MaxQuantum float64 `json:"maxQuantum"`
		MinRate    int     `json:"minRate"`
		MaxRate    int     `json:"maxRate"`
		MinNs      int     `json:"minNs"`
		MaxNs      int     `json:"maxNs"`
	}

	NodeProps struct {
		Name                     string       `json:"node.name"`
		Description              string       `json:"node.description"`
		Nickname                 string       `json:"node.nick"`
		AudioChannels            int          `json:"audio.channels"`
		AudioPosition            string       `json:"audio.position"`
		ClientID                 int          `json:"client.id"`
		DeviceClass              *DeviceClass `json:"device.class"`
		DeviceID                 int          `json:"device.id"`
		DeviceProfileDescription string       `json:"device.profile.description"`
		DeviceProfileName        string       `json:"device.profile.name"`
		FactoryID                int          `json:"factory.id"`
		FactoryMode              string       `json:"factory.mode"`
		FactoryName              string       `json:"factory.name"`
		LibraryName              string       `json:"library.name"`
		MediaClass               MediaClass   `json:"media.class"`
		ObjectID                 int          `json:"object.id"`
		ObjectPath               string       `json:"object.path"`
		ObjectSerial             int          `json:"object.serial"`
	}
)

type DeviceClass string

const DeviceSound DeviceClass = "sound"

type State string

const (
	StateSuspended State = "suspended"
	StateRunning   State = "running"
	StateIdle      State = "idle"
	StateError     State = "error"
	StateCreating  State = "creating"
)

type MediaClass string

const (
	// A source of audio samples like a microphone
	MediaAudioSource MediaClass = "Audio/Source"
	// A sink for audio samples, like an audio card
	MediaAudioSink MediaClass = "Audio/Sink"
	// A node that is both a sink and a source
	MediaAudioDuplex MediaClass = "Audio/Duplex"
	// A playback stream
	MediaStreamOutputAudio MediaClass = "Stream/Output/Audio"
	// A capture stream
	MediaStreamInputAudio MediaClass = "Stream/Input/Audio"
)

// Example of when an object is removed:
//
//	{
//		"id": 128,
//		"info": null
//	 }
func (e *Event) IsRemovalEvent() bool {
	return e.Info == nil && e.Type == EmptyEvent && e.ID != 0
}

// Tries to convert the JSON info field to NodeProps
func (e *Event) NodeProps() (*NodeProps, error) {
	if e.Type != EventNode {
		return nil, fmt.Errorf("event is not a node event type")
	} else if e.Info == nil {
		return nil, fmt.Errorf("event info is nil")
	}

	var props = &NodeProps{}
	data, err := json_v2.Marshal(e.Info.Props)
	if err != nil {
		return props, err
	}

	err = json_v2.Unmarshal(data, props)
	return props, err
}
