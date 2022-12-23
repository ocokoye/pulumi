// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pulumi

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func makeRemoteResourceTransformation(ctx *Context, reference string) ResourceTransformation {
	return func(args *ResourceTransformationArgs) *ResourceTransformationResult {
		marshaledArgs, err := args.marshal()
		if err != nil {
			panic(err)
		}

		invokeArgs := resource.PropertyMap{
			"reference": resource.NewStringProperty(reference),
			"args":      resource.NewObjectProperty(marshaledArgs),
		}

		rpcArgs, err := plugin.MarshalProperties(
			invokeArgs,
			ctx.withKeepOrRejectUnknowns(plugin.MarshalOptions{
				KeepSecrets:      true,
				KeepResources:    true,
				KeepOutputValues: true,
			}),
		)
		if err != nil {
			panic(fmt.Errorf("marshaling arguments: %w", err))
		}

		tok := "pulumi:pulumi:invokeCallback"
		resp, err := ctx.monitor.Invoke(ctx.ctx, &pulumirpc.ResourceInvokeRequest{
			Tok:             tok,
			Args:            rpcArgs,
			AcceptResources: true,
		})
		if err != nil {
			panic(fmt.Errorf("invoke(%s, ...): error: %v", tok, err))
		}

		unmarshalledReturn, err := plugin.UnmarshalProperties(resp.Return, plugin.MarshalOptions{
			KeepUnknowns:     true,
			KeepSecrets:      true,
			KeepResources:    true,
			KeepOutputValues: true,
		})
		if err != nil {
			panic(fmt.Errorf("unmarshaling return: %w", err))
		}
		if unmarshalledReturn["result"].IsNull() {
			return nil
		}
		result := unmarshalledReturn["result"].ObjectValue()

		var props Input
		if result["props"].IsObject() {
			props = unmarshalMap(result["props"].ObjectValue())
		}
		return &ResourceTransformationResult{
			Props: props,
			Opts:  args.Opts, // TODO make this go over the wire
		}
	}
}
