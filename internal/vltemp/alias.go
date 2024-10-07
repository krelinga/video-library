package vltemp

import "github.com/krelinga/video-library/internal/vltemp/internal/ids"

type VolumeWfId = ids.VolumeWfId

func NewVolumeWfId(name string) (VolumeWfId, error) {
	return ids.NewVolumeWfId(name)
}

type DiscWfId = ids.DiscWfId

func NewDiscWfId(volumeWfId VolumeWfId, discName string) (DiscWfId, error) {
	return ids.NewDiscWfId(volumeWfId, discName)
}
func NewDiscWfIdFromString(asString string) (DiscWfId, error) {
	return ids.NewDiscWfIdFromString(asString)
}
