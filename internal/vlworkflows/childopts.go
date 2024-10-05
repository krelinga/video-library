package vlworkflows

import (
	temporalenums "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

func childOptions(workflowID string) workflow.ChildWorkflowOptions {
	return workflow.ChildWorkflowOptions{
		WorkflowID:            workflowID,
		WorkflowIDReusePolicy: temporalenums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		ParentClosePolicy:     temporalenums.PARENT_CLOSE_POLICY_ABANDON,
	}
}