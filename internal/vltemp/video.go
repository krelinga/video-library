package vltemp

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

const (
	VideoUpdateBootstrap = "video-update-bootstrap"
)

type VideoLineage struct {
	FromDisc *VideoFromDisc `json:"from_disc"`
	// TODO: eventually support other options here.
}

type VideoFromDisc struct {
	DiscID   string `json:"disc_id"`
	Filename string `json:"filename"`
}

type VideoUpdateBootstrapRequest struct {
	Lineage *VideoLineage `json:"lineage"`
}

var ErrCorruptVideoLineage = errors.New("corrupt video lineage")

func VideoPath(ctx context.Context, videoLineage *VideoLineage) (string, error) {
	switch {
	case videoLineage.FromDisc != nil:
		discPath, err := DiscPath(ctx, videoLineage.FromDisc.DiscID)
		if err != nil {
			return "", errors.Join(ErrCorruptVideoLineage, err)
		}
		if videoLineage.FromDisc.Filename == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing filename"))
		}
		return filepath.Join(discPath, videoLineage.FromDisc.Filename), nil
	default:
		return "", errors.Join(ErrCorruptVideoLineage, errors.New("unknown lineage"))
	}
}

func LegacyVideoID(videoLineage *VideoLineage) (string, error) {
	switch {
	case videoLineage.FromDisc != nil:
		if videoLineage.FromDisc.DiscID == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing DiscID"))
		}
		if videoLineage.FromDisc.Filename == "" {
			return "", errors.Join(ErrCorruptVideoLineage, errors.New("missing filename"))
		}
		return filepath.Join(videoLineage.FromDisc.DiscID, videoLineage.FromDisc.Filename), nil
	default:
		return "", errors.Join(ErrCorruptVideoLineage, errors.New("unknown lineage"))
	}
}

type VideoWfId string

type parsedVideoWfId struct {
	protocol string

	// Only one of these will be set depending on protocol.
	discWfId DiscWfId
	// other types here in the future
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
		p.discWfId = DiscWfId(other)
		err = p.discWfId.Validate()
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

func (id VideoWfId) DiscWfId() (discWfId DiscWfId, ok bool) {
	parsed, err := id.parse()
	if err != nil{
		panic(err)
	}
	if parsed.protocol != "disc" {
		return
	}
	discWfId = parsed.discWfId
	ok = true
	return
}