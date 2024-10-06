package vlactivities

import (
	"context"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vltemp"
	"github.com/krelinga/video-library/internal/vlcontext"
	"github.com/krelinga/video-library/internal/vllib"
	"github.com/krelinga/video-library/internal/vltypes"
	"go.temporal.io/sdk/client"
)

func DiscReadVideoNames(ctx context.Context, discID string) ([]string, error) {
	dir, err := vllib.DiscPath(ctx, discID)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var videos []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".mkv") || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		videos = append(videos, entry.Name())
	}
	return videos, nil
}

var DiscReadVideoNamesOptions = lightOptions

func DiscBootstrapVideo(ctx context.Context, discID, videoFilename string) (string, error) {
	temporalClient := vlcontext.GetTemporalClient(ctx)
	lineage := &vltypes.VideoLineage{
		FromDisc: &vltypes.VideoFromDisc{
			DiscID:   discID,
			Filename: videoFilename,
		},
	}
	videoID, err := vllib.VideoID(lineage)
	if err != nil {
		return "", err
	}

	request := &vltypes.VideoUpdateBootstrapRequest{
		Lineage: lineage,
	}
	opts := client.StartWorkflowOptions{
		ID: videoID,
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, vltemp.Disc, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID:            uuid.New().String(),
		UpdateName:          vltemp.VideoUpdateBootstrap,
		WorkflowID:          videoID,
		WaitForStage:        client.WorkflowUpdateStageCompleted,
		FirstExecutionRunID: wf.GetRunID(),
		Args:                []interface{}{request},
	})
	if err != nil {
		return "", err
	}
	if err := updateHandle.Get(ctx, nil); err != nil {
		return "", err
	}
	return videoID, nil
}

var DiscBootstrapVideoOptions = lightOptions
