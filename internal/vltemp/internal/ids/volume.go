package ids

import (
	"encoding/json"
	"errors"
	"strings"
)

type VolumeWfId interface {
	Name() string
}

type volumeWfIdImpl struct {
	name string
}

func (id *volumeWfIdImpl) Name() string {
	return id.name
}

func (id *volumeWfIdImpl) String() string {
	return id.name
}

func (id *volumeWfIdImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.name)
}

func (id *volumeWfIdImpl) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	temp, err := newVolumeWfIdImpl(name)
	if err != nil {
		return err
	}
	*id = *temp
	return nil
}

var ErrInvalidVolumeName = errors.New("invalid volume name")

func newVolumeWfIdImpl(name string) (*volumeWfIdImpl, error) {
	if name == "" || strings.Contains(name, "/") {
		return nil, ErrInvalidVolumeName
	}
	return &volumeWfIdImpl{name: name}, nil
}

func NewVolumeWfId(name string) (VolumeWfId, error) {
	return newVolumeWfIdImpl(name)
}
