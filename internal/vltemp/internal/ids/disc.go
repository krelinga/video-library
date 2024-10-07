package ids

import (
	"encoding/json"
	"errors"
	"strings"
)

type DiscWfId interface {
	VolumeWfId() VolumeWfId
	DiscName() string

	String() string
}

type discWfIdImpl struct {
	volumeWfId VolumeWfId
	discName   string
}

func (id *discWfIdImpl) VolumeWfId() VolumeWfId {
	return id.volumeWfId
}

func (id *discWfIdImpl) DiscName() string {
	return id.discName
}

func (id *discWfIdImpl) String() string {
	return strings.Join([]string{id.volumeWfId.String(), id.discName}, "/")
}

func (id *discWfIdImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *discWfIdImpl) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}
	temp, err := newDiscWfIdImplFromString(asString)
	if err != nil {
		return err
	}
	*id = *temp
	return nil
}

var ErrInvalidDiscNameString = errors.New("invalid disc name string")

func newDiscWfIdImpl(volumeWfId VolumeWfId, discName string) (*discWfIdImpl, error) {
	if discName == "" || strings.Contains(discName, "/") {
		return nil, ErrInvalidDiscNameString
	}
	return &discWfIdImpl{volumeWfId: volumeWfId, discName: discName}, nil
}

var ErrInvalidDiscWfIdString = errors.New("invalid disc workflow id string")

func newDiscWfIdImplFromString(asString string) (*discWfIdImpl, error) {
	parts := strings.Split(asString, "/")
	if len(parts) != 2 {
		return nil, ErrInvalidDiscWfIdString
	}
	volumeWfId, err := NewVolumeWfId(parts[0])
	if err != nil {
		return nil, ErrInvalidDiscWfIdString
	}
	return newDiscWfIdImpl(volumeWfId, parts[1])
}

func NewDiscWfId(volumeWfId VolumeWfId, discName string) (DiscWfId, error) {
	return newDiscWfIdImpl(volumeWfId, discName)
}

func NewDiscWfIdFromString(asString string) (DiscWfId, error) {
	return newDiscWfIdImplFromString(asString)
}
