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

func DiscPath(ctx context.Context, discWfId DiscWfId) string {
	volumePath := VolumePath(ctx, discWfId.VolumeWfId())
	return filepath.Join(volumePath, discWfId.Name())
}

func actDiscReadVideoNames(ctx context.Context, discWfId DiscWfId) ([]string, error) {
	dir := DiscPath(ctx, discWfId)
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

var actDiscReadVideoNamesOptions = lightOptions

func actDiscNewVideo(ctx context.Context, discWfId DiscWfId, videoFilename string) (VideoWfId, error) {
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

var actDiscNewVideoOptions = lightOptions

func discWfNew(ctx workflow.Context, discWfId DiscWfId, state *DiscWFState) error {
	var videoFiles []string
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, actDiscReadVideoNamesOptions),
		actDiscReadVideoNames, discWfId).Get(ctx, &videoFiles)
	if err != nil {
		return err
	}
	for _, videoFile := range state.Videos {
		var videoWfId VideoWfId
		err = workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, actDiscNewVideoOptions),
			actDiscNewVideo, discWfId, videoFile).Get(ctx, &videoWfId)
		if err != nil {
			return err
		}
		state.Videos = append(state.Videos, videoWfId)
	}
	return nil
}

func DiscWF(ctx workflow.Context, state *DiscWFState) error {
	discWfId := DiscWfId(workflow.GetInfo(ctx).WorkflowExecution.ID)
	if err := discWfId.Validate(); err != nil {
		return err
	}
	wt := workTracker{}

	// TODO: don't rely on an update to do this, just do it as soon as the workflow starts
	bootstrap := func(ctx workflow.Context) error {
		state = &DiscWFState{}
		err := discWfNew(ctx, discWfId, state)
		wt.WorkIfNoError(err)
		return err

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
