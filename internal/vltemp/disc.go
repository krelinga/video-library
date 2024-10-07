package vltemp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/krelinga/video-library/internal/vlcontext"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

const (
	DiscWFUpdateNameBootstrap = "DiscWFUpdateBootstrap"
)

type DiscWFState struct {
	Videos []VideoWfId `json:"videos"`
}

func DiscPath(ctx context.Context, discWfId DiscWfId) (string, error) {
	volumePath, err := VolumePath(ctx, discWfId.VolumeWfId())
	if err != nil {
		return "", err
	}
	return filepath.Join(volumePath, discWfId.Name()), nil
}

func DiscReadVideoNames(ctx context.Context, discWfId DiscWfId) ([]string, error) {
	dir, err := DiscPath(ctx, discWfId)
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

func DiscBootstrapVideo(ctx context.Context, discWfId DiscWfId, videoFilename string) (VideoWfId, error) {
	temporalClient := vlcontext.GetTemporalClient(ctx)
	videoWfId, err := NewVideoWfIdFromDisc(discWfId, videoFilename)
	if err != nil {
		return "", err
	}

	opts := client.StartWorkflowOptions{
		ID: string(videoWfId),
	}
	wf, err := temporalClient.ExecuteWorkflow(ctx, opts, DiscWF, nil)
	if err != nil {
		return "", err
	}
	updateHandle, err := temporalClient.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
		UpdateID:            uuid.New().String(),
		UpdateName:          VideoUpdateBootstrap,
		WorkflowID:          string(videoWfId),
		WaitForStage:        client.WorkflowUpdateStageCompleted,
		FirstExecutionRunID: wf.GetRunID(),
	})
	if err != nil {
		return "", err
	}
	if err := updateHandle.Get(ctx, nil); err != nil {
		return "", err
	}
	return videoWfId, nil
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
			var videoWfId VideoWfId
			err = workflow.ExecuteActivity(
				workflow.WithActivityOptions(ctx, DiscBootstrapVideoOptions),
				DiscBootstrapVideo, discId, videoFile).Get(ctx, &videoWfId)
			if err != nil {
				return err
			}
			state.Videos = append(state.Videos, videoWfId)
		}
		return
	}

	err := workflow.SetUpdateHandler(ctx, DiscWFUpdateNameBootstrap, bootstrap)
	if err != nil {
		return err
	}

	err = workflow.Await(ctx, wt.AwaitFunc())
	if err != nil {
		return err
	}

	return workflow.NewContinueAsNewError(ctx, DiscWF, state)
}

// A more-refined string to handle Temporal Workflow IDs for Disc workflows.
//
// Use NewDiscWfId() to create a new DiscWfId.  You can also directly case from a string
// with `DiscWfId("my-disc")`, but this will not validate the ID.  You can validate the ID
// with the Validate() method.  Any other methods called on an invalid DiscWfId will panic.
type DiscWfId string

func (id DiscWfId) parse() (volumeWfId VolumeWfId, name string, err error) {
	parts := strings.Split(string(id), "/")
	if len(parts) != 2 || !nameIsValid(parts[1]) {
		err = ErrInvalidWorkflowId
		return
	}
	name = parts[1]
	volumeWfId, err = NewVolumeWfId(parts[0])
	return
}

// Validates the DiscWfId.
func (id DiscWfId) Validate() error {
	_, _, err := id.parse()
	return err
}

// Returns the VolumeWfId of the DiscWfId.
//
// Panics if the DiscWfId is invalid.
func (id DiscWfId) VolumeWfId() VolumeWfId {
	volumeWfId, _, err := id.parse()
	if err != nil {
		panic(err)
	}
	return volumeWfId
}

// Returns the name of the Disc.
//
// Panics if the DiscWfId is invalid.
func (id DiscWfId) Name() string {
	_, ame, err := id.parse()
	if err != nil {
		panic(err)
	}
	return ame
}

func NewDiscWfId(volumeWfId VolumeWfId, discFilename string) (DiscWfId, error) {
	id := DiscWfId(fmt.Sprintf("%s/%s", volumeWfId, discFilename))
	return id, id.Validate()
}
