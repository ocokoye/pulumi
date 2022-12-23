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
	"reflect"
	"sync"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// ResourceTransformationArgs is the argument bag passed to a resource transformation.
type ResourceTransformationArgs struct {
	// The resource instance that is being transformed.
	Resource Resource
	// The type of the resource.
	Type string
	// The name of the resource.
	Name string
	// The original properties passed to the resource constructor.
	Props Input
	// The original resource options passed to the resource constructor.
	Opts []ResourceOption
}

// ResourceTransformationResult is the result that must be returned by a resource transformation
// callback.  It includes new values to use for the `props` and `opts` of the `Resource` in place of
// the originally provided values.
type ResourceTransformationResult struct {
	// The new properties to use in place of the original `props`.
	Props Input
	// The new resource options to use in place of the original `opts`.
	Opts []ResourceOption
}

// ResourceTransformation is the callback signature for the `transformations` resource option.  A
// transformation is passed the same set of inputs provided to the `Resource` constructor, and can
// optionally return back alternate values for the `props` and/or `opts` prior to the resource
// actually being created.  The effect will be as though those props and opts were passed in place
// of the original call to the `Resource` constructor.  If the transformation returns nil,
// this indicates that the resource will not be transformed.
type ResourceTransformation func(*ResourceTransformationArgs) *ResourceTransformationResult

func unmarshalMap(m resource.PropertyMap) Map {
	props := Map{}
	for k, v := range m {
		if v.IsString() {
			props[string(k)] = String(v.StringValue())
		} else if v.IsNumber() {
			props[string(k)] = Float64(v.NumberValue())
		} else if v.IsObject() {
			props[string(k)] = unmarshalMap(v.ObjectValue())
		}
	}
	return props
}

// func unmarshalResourceTransformationArgs(args resource.PropertyMap) ResourceTransformationArgs {
// 	return ResourceTransformationArgs{
// 		Type:  args["type"].StringValue(),
// 		Name:  args["name"].StringValue(),
// 		Props: unmarshalMap(args["props"].ObjectValue()),
// 		// TODO Opts
// 	}
// }

func (args *ResourceTransformationArgs) marshal() (resource.PropertyMap, error) {
	propsValue := resource.NewNullProperty()
	if args.Props != nil {
		marshaledProps, _, _, err := marshalInputs(args.Props)
		if err != nil {
			return nil, fmt.Errorf("marshaling resource transformation props: %w", err)
		}
		propsValue = resource.NewObjectProperty(marshaledProps)
	}

	// Note: We currently explicitly do not pass Resource.
	return resource.PropertyMap{
		"type":  resource.NewStringProperty(args.Type),
		"name":  resource.NewStringProperty(args.Name),
		"props": propsValue,
		// TODO opts
	}, nil
}

var typeTokenToArgsType sync.Map // map[string]reflect.Type

// RegisterInputType registers an Input type with the Pulumi runtime. This allows the input type to be instantiated
// for a given input interface.
func RegisterResourceArgsType(typ string, input Input) {
	concreteType := reflect.TypeOf(input)
	existing, hasExisting := typeTokenToArgsType.LoadOrStore(typ, concreteType)
	if hasExisting {
		panic(fmt.Errorf("an args type for %q is already registered: %v", typ, existing))
	}
}
