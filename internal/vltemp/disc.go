package vltemp

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vlcontext"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

const (
	Disc                = "disc"
	DiscUpdateBootstrap = "disc-update-bootstrap"
)

type DiscWFState struct {
	Videos []string `json:"videos"`
}

var ErrInvalidDiscID = errors.New("invalid discID")
var ErrInvalidDiscBase = errors.New("invalid discBase")

func DiscParseID(discID string, volumeID, discBase *string) error {
	parts := strings.Split(discID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ErrInvalidDiscID
	}
	if volumeID != nil {
		*volumeID = parts[0]
	}
	if discBase != nil {
		*discBase = parts[1]
	}
	return nil
}

func DiscID(volumeID, discBase string) (string, error) {
	if err := ValidateVolumeID(volumeID); err != nil {
		return "", err
	}
	if discBase == "" {
		return "", ErrInvalidDiscBase
	}
	return filepath.Join(volumeID, discBase), nil
}

func DiscPath(ctx context.Context, discID string) (string, error) {
	var volumeID, discBase string
	err := DiscParseID(discID, &volumeID, &discBase)
	if err != nil {
		return "", err
	}
	volumePath, err := VolumePath(ctx, volumeID)
	if err != nil {
		return "", err
	}
	return filepath.Join(volumePath, discBase), nil
}

func DiscReadVideoNames(ctx context.Context, discID string) ([]string, error) {
	dir, err := DiscPath(ctx, discID)
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
	lineage := &VideoLineage{
		FromDisc: &VideoFromDisc{
			DiscID:   discID,
			Filename: videoFilename,
		},
	}
	videoID, err := VideoID(lineage)
	if err != nil {
		return "", err
	}

	request := &VideoUpdateBootstrapRequest{
		Lineage: lineage,
	}
	opts := client.StartWorkflowOptions{
		ID: videoID,
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, Disc, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID:            uuid.New().String(),
		UpdateName:          VideoUpdateBootstrap,
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

func DiscWF(ctx workflow.Context, state *DiscWFState) error {
	discId := workflow.GetInfo(ctx).WorkflowExecution.ID
	wt := workTracker{}

	bootstrap := func(ctx workflow.Context) (err error) {
		defer wt.WorkIfNoError(err)

		state = &DiscWFState{}
		var videoFiles []string
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, DiscReadVideoNamesOptions),
			DiscReadVideoNames, discId).Get(ctx, &videoFiles)
		if err != nil {
			return
		}
		for _, videoFile := range state.Videos {
			var videoId string
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, DiscBootstrapVideoOptions),
				DiscBootstrapVideo, discId, videoFile).Get(ctx, &videoId)
			if err != nil {
				return err
			}
			state.Videos = append(state.Videos, videoId)
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, DiscUpdateBootstrap, bootstrap)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, DiscWF, state)
}
