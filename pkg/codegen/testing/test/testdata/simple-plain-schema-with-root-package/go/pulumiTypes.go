// Code generated by test DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package different

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Foo struct {
	A bool    `pulumi:"a"`
	B *bool   `pulumi:"b"`
	C int     `pulumi:"c"`
	D *int    `pulumi:"d"`
	E string  `pulumi:"e"`
	F *string `pulumi:"f"`
}

// FooInput is an input type that accepts FooArgs and FooOutput values.
// You can construct a concrete instance of `FooInput` via:
//
//	FooArgs{...}
type FooInput interface {
	pulumi.Input

	ToFooOutput() FooOutput
	ToFooOutputWithContext(context.Context) FooOutput
}

type FooArgs struct {
	A bool    `pulumi:"a"`
	B *bool   `pulumi:"b"`
	C int     `pulumi:"c"`
	D *int    `pulumi:"d"`
	E string  `pulumi:"e"`
	F *string `pulumi:"f"`
}

func (FooArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*Foo)(nil)).Elem()
}

func (i FooArgs) ToFooOutput() FooOutput {
	return i.ToFooOutputWithContext(context.Background())
}

func (i FooArgs) ToFooOutputWithContext(ctx context.Context) FooOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FooOutput)
}

func (i FooArgs) ToFooPtrOutput() FooPtrOutput {
	return i.ToFooPtrOutputWithContext(context.Background())
}

func (i FooArgs) ToFooPtrOutputWithContext(ctx context.Context) FooPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FooOutput).ToFooPtrOutputWithContext(ctx)
}

// FooPtrInput is an input type that accepts FooArgs, FooPtr and FooPtrOutput values.
// You can construct a concrete instance of `FooPtrInput` via:
//
//	        FooArgs{...}
//
//	or:
//
//	        nil
type FooPtrInput interface {
	pulumi.Input

	ToFooPtrOutput() FooPtrOutput
	ToFooPtrOutputWithContext(context.Context) FooPtrOutput
}

type fooPtrType FooArgs

func FooPtr(v *FooArgs) FooPtrInput {
	return (*fooPtrType)(v)
}

func (*fooPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**Foo)(nil)).Elem()
}

func (i *fooPtrType) ToFooPtrOutput() FooPtrOutput {
	return i.ToFooPtrOutputWithContext(context.Background())
}

func (i *fooPtrType) ToFooPtrOutputWithContext(ctx context.Context) FooPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FooPtrOutput)
}

// FooArrayInput is an input type that accepts FooArray and FooArrayOutput values.
// You can construct a concrete instance of `FooArrayInput` via:
//
//	FooArray{ FooArgs{...} }
type FooArrayInput interface {
	pulumi.Input

	ToFooArrayOutput() FooArrayOutput
	ToFooArrayOutputWithContext(context.Context) FooArrayOutput
}

type FooArray []FooInput

func (FooArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Foo)(nil)).Elem()
}

func (i FooArray) ToFooArrayOutput() FooArrayOutput {
	return i.ToFooArrayOutputWithContext(context.Background())
}

func (i FooArray) ToFooArrayOutputWithContext(ctx context.Context) FooArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(FooArrayOutput)
}

type FooOutput struct{ *pulumi.OutputState }

func (FooOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Foo)(nil)).Elem()
}

func (o FooOutput) ToFooOutput() FooOutput {
	return o
}

func (o FooOutput) ToFooOutputWithContext(ctx context.Context) FooOutput {
	return o
}

func (o FooOutput) ToFooPtrOutput() FooPtrOutput {
	return o.ToFooPtrOutputWithContext(context.Background())
}

func (o FooOutput) ToFooPtrOutputWithContext(ctx context.Context) FooPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v Foo) *Foo {
		return &v
	}).(FooPtrOutput)
}

func (o FooOutput) A() pulumi.BoolOutput {
	return o.ApplyT(func(v Foo) bool { return v.A }).(pulumi.BoolOutput)
}

func (o FooOutput) B() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v Foo) *bool { return v.B }).(pulumi.BoolPtrOutput)
}

func (o FooOutput) C() pulumi.IntOutput {
	return o.ApplyT(func(v Foo) int { return v.C }).(pulumi.IntOutput)
}

func (o FooOutput) D() pulumi.IntPtrOutput {
	return o.ApplyT(func(v Foo) *int { return v.D }).(pulumi.IntPtrOutput)
}

func (o FooOutput) E() pulumi.StringOutput {
	return o.ApplyT(func(v Foo) string { return v.E }).(pulumi.StringOutput)
}

func (o FooOutput) F() pulumi.StringPtrOutput {
	return o.ApplyT(func(v Foo) *string { return v.F }).(pulumi.StringPtrOutput)
}

type FooPtrOutput struct{ *pulumi.OutputState }

func (FooPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Foo)(nil)).Elem()
}

func (o FooPtrOutput) ToFooPtrOutput() FooPtrOutput {
	return o
}

func (o FooPtrOutput) ToFooPtrOutputWithContext(ctx context.Context) FooPtrOutput {
	return o
}

func (o FooPtrOutput) Elem() FooOutput {
	return o.ApplyT(func(v *Foo) Foo {
		if v != nil {
			return *v
		}
		var ret Foo
		return ret
	}).(FooOutput)
}

func (o FooPtrOutput) A() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Foo) *bool {
		if v == nil {
			return nil
		}
		return &v.A
	}).(pulumi.BoolPtrOutput)
}

func (o FooPtrOutput) B() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Foo) *bool {
		if v == nil {
			return nil
		}
		return v.B
	}).(pulumi.BoolPtrOutput)
}

func (o FooPtrOutput) C() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Foo) *int {
		if v == nil {
			return nil
		}
		return &v.C
	}).(pulumi.IntPtrOutput)
}

func (o FooPtrOutput) D() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Foo) *int {
		if v == nil {
			return nil
		}
		return v.D
	}).(pulumi.IntPtrOutput)
}

func (o FooPtrOutput) E() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Foo) *string {
		if v == nil {
			return nil
		}
		return &v.E
	}).(pulumi.StringPtrOutput)
}

func (o FooPtrOutput) F() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Foo) *string {
		if v == nil {
			return nil
		}
		return v.F
	}).(pulumi.StringPtrOutput)
}

type FooArrayOutput struct{ *pulumi.OutputState }

func (FooArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Foo)(nil)).Elem()
}

func (o FooArrayOutput) ToFooArrayOutput() FooArrayOutput {
	return o
}

func (o FooArrayOutput) ToFooArrayOutputWithContext(ctx context.Context) FooArrayOutput {
	return o
}

func (o FooArrayOutput) Index(i pulumi.IntInput) FooOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Foo {
		return vs[0].([]Foo)[vs[1].(int)]
	}).(FooOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*FooInput)(nil)).Elem(), FooArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*FooPtrInput)(nil)).Elem(), FooArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*FooArrayInput)(nil)).Elem(), FooArray{})
	pulumi.RegisterOutputType(FooOutput{})
	pulumi.RegisterOutputType(FooPtrOutput{})
	pulumi.RegisterOutputType(FooArrayOutput{})
}
