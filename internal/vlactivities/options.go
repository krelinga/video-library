package vlactivities

import (
	"time"

	"github.com/krelinga/video-library/internal/vltemp"
	"go.temporal.io/sdk/workflow"
)

var lightOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           vltemp.Light,
}