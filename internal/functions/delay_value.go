package functions

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = DelayValue{}

func NewDelayValue() function.Function {
	return DelayValue{}
}

type DelayValue struct{}

func (d DelayValue) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "delay_value"
}

func (d DelayValue) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Return the provided value after waiting for a specified time.",
		Description:         "Wait for a specified delay (in seconds) before returning the provided value.",
		MarkdownDescription: "Wait for a specified delay (in seconds) before returning the provided value.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "delay",
				Description:         "The delay time in seconds before returning the value.",
				MarkdownDescription: "The delay time in seconds before returning the value.",
			},
			function.StringParameter{
				Name:                "value",
				Description:         "The value to return after the delay.",
				MarkdownDescription: "The value to return after the delay.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (d DelayValue) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var delayStr, value string

	if err := req.Arguments.Get(ctx, &delayStr); err != nil {
		resp.Error = err
		return
	}
	if err := req.Arguments.Get(ctx, &value); err != nil {
		resp.Error = err
		return
	}

	delaySec, err := strconv.Atoi(delayStr)
	if err != nil {
		resp.Error = function.NewFuncError("Invalid delay value, expected an integer. Error: " + err.Error())
		return
	}

	time.Sleep(time.Duration(delaySec) * time.Second)

	if err := resp.Result.Set(ctx, value); err != nil {
		resp.Error = err
	}
}
