package vltemp

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

var lightOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 5 * time.Second,
	TaskQueue:           Light,
}