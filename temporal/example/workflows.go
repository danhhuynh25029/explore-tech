package iplocate

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

func WorkflowA(ctx workflow.Context, name string) (string, error) {
	//Define the activity options, including the retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		//RetryPolicy: &temporal.RetryPolicy{
		//	InitialInterval:    time.Second, //amount of time that must elapse before the first retry occurs
		//	MaximumInterval:    time.Minute, //maximum interval between retries
		//	BackoffCoefficient: 2,           //how much the retry interval increases
		//	MaximumAttempts:    5,           // Uncomment this if you want to limit attempts
		//},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var Activities *Activities

	err := workflow.ExecuteActivity(ctx, Activities.ActivityA).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	err = workflow.ExecuteActivity(ctx, Activities.ActivityB).Get(ctx, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("success"), nil
}

func WorkflowB(ctx workflow.Context, name string) (string, error) {
	return "", nil
}
