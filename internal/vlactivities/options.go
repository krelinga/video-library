package vlactivities

import (
	"time"

	"github.com/krelinga/video-library/internal/vlqueues"
	"go.temporal.io/sdk/workflow"
)

var lightOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           vlqueues.Light,
}