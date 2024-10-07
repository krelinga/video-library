package vltemp

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

const (
	VideoUpdateBootstrap = "video-update-bootstrap"
)

func VideoPath(ctx context.Context, videoWfId VideoWfId) (string, error) {
	if discWfId, discFilepath, ok := videoWfId.FromDisc(); ok {
		discPath, err := DiscPath(ctx, discWfId)
		if err != nil {
			return "", err
		}
		return filepath.Join(discPath, discFilepath), nil
	} else if filepath, ok := videoWfId.FromFilepath(); ok {
		return filepath, nil
	} else {
		panic("unexpected protocol " + videoWfId.Protocol())
	}
}

type VideoWfId string

type parsedVideoWfId struct {
	protocol string

	// These fields will be set if protocol is "disc".
	discWfId     DiscWfId
	discFilename string

	// These fields will be set if protocol is "filepath".
	filepath string
}

func (id VideoWfId) parse() (p parsedVideoWfId, err error) {
	parts := strings.Split(string(id), ":")
	if len(parts) != 2 {
		err = ErrInvalidWorkflowId
		return
	}
	p.protocol = parts[0]

	other := parts[1]
	switch p.protocol {
	case "disc":
		parts := strings.Split(other, "/")
		if len(parts) < 3 {
			err = ErrInvalidWorkflowId
			return
		}
		p.discWfId = DiscWfId(fmt.Sprintf("%s/%s", parts[0], parts[1]))
		err = p.discWfId.Validate()
		if err != nil {
			return
		}
		p.discFilename = strings.Join(parts[2:], "/")
		if !pathIsValid(p.discFilename) {
			err = ErrInvalidWorkflowId
			return
		}
	case "filepath":
		p.filepath = other
		if !rootPathIsValid(p.filepath) {
			err = ErrInvalidWorkflowId
		}
	default:
		err = ErrInvalidWorkflowId
	}
	return
}

func (id VideoWfId) Validate() error {
	_, err := id.parse()
	return err
}

func (id VideoWfId) Protocol() string {
	parsed, err := id.parse()
	if err != nil {
		panic(err)
	}
	return parsed.protocol
}

func (id VideoWfId) FromDisc() (discWfId DiscWfId, videoFilename string, ok bool) {
	parsed, err := id.parse()
	if err != nil {
		panic(err)
	}
	if parsed.protocol != "disc" {
		return
	}
	discWfId = parsed.discWfId
	videoFilename = parsed.discFilename
	ok = true
	return
}

func (id VideoWfId) FromFilepath() (filepath string, ok bool) {
	parsed, err := id.parse()
	if err != nil {
		panic(err)
	}
	if parsed.protocol != "filepath" {
		return
	}
	filepath = parsed.filepath
	ok = true
	return
}

// Create a new VideoWfId for videos contained within a given disc.
func NewVideoWfIdFromDisc(discWfId DiscWfId, videoFileName string) (VideoWfId, error) {
	id := VideoWfId(fmt.Sprintf("disc:%s/%s", discWfId, videoFileName))
	return id, id.Validate()
}

// Create a new VideoWfId for videos represented in the system by a raw file path.
func NewVideoWfIdFromFilepath(filepath string) (VideoWfId, error) {
	id := VideoWfId(fmt.Sprintf("filepath:%s", filepath))
	return id, id.Validate()
}
