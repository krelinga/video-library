package vlapi

import (
	"context"

	vlproto "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/library/v1"
	"connectrpc.com/connect"
)

type Service struct {}

func (s *Service) NewVolume(ctx context.Context, req *connect.Request[vlproto.NewVolumeRequest]) (*connect.Response[vlproto.NewVolumeResponse], error) {
	return nil, nil
}

func (s *Service) FindVolumes(ctx context.Context, req *connect.Request[vlproto.FindVolumesRequest]) (*connect.Response[vlproto.FindVolumesResponse], error) {
	return nil, nil
}

func (s *Service) DiscoverNewDiscs(ctx context.Context, req *connect.Request[vlproto.DiscoverNewDiscsRequest]) (*connect.Response[vlproto.DiscoverNewDiscsResponse], error) {
	return nil, nil
}

func (s *Service) FindDiscs(ctx context.Context, req *connect.Request[vlproto.FindDiscsRequest]) (*connect.Response[vlproto.FindDiscsResponse], error) {
	return nil, nil
}