// Copyright 2017 Pulumi, Inc. All rights reserved.

package resource

import (
	"github.com/golang/glog"

	"github.com/pulumi/coconut/pkg/compiler/core"
	"github.com/pulumi/coconut/pkg/compiler/errors"
	"github.com/pulumi/coconut/pkg/diag/colors"
	"github.com/pulumi/coconut/pkg/graph"
	"github.com/pulumi/coconut/pkg/pack"
	"github.com/pulumi/coconut/pkg/tokens"
	"github.com/pulumi/coconut/pkg/util/contract"
)

// TODO: concurrency.
// TODO: handle output dependencies

// Plan is the output of analyzing resource graphs and contains the steps necessary to perform an infrastructure
// deployment.  A plan can be generated out of whole cloth from a resource graph -- in the case of new deployments --
// however, it can alternatively be generated by diffing two resource graphs -- in the case of updates to existing
// environments (presumably more common).  The plan contains step objects that can be used to drive a deployment.
type Plan interface {
	Empty() bool                                        // true if the plan is empty.
	Steps() Step                                        // the first step to perform, linked to the rest.
	Replaces() map[URN][]PropertyKey                    // resources being replaced and their properties.
	Unchanged() map[Resource]Resource                   // the resources untouched by this plan.
	Apply(prog Progress) (Snapshot, Step, State, error) // performs the operations specified in this plan.
}

// Progress can be used for progress reporting.
type Progress interface {
	Before(step Step)
	After(step Step, state State, err error)
}

// Step is a specification for a deployment operation.
type Step interface {
	Plan() Plan            // the plan this step belongs to.
	Op() StepOp            // the operation that will be performed.
	Logical() bool         // true if this is a logical step, rather than a physical one.
	Old() Resource         // the old resource state, if any, before performing this step.
	New() Resource         // the new resource state, if any, after performing this step.
	NewProps() PropertyMap // the projected new resource state, factoring in dependency updates.
	Next() Step            // the next step to perform, or nil if none.
	Apply() (State, error) // performs the operation specified by this step.
}

// StepOp represents the kind of operation performed by this step.
type StepOp string

const (
	OpCreate        StepOp = "create"         // creating a new resource.
	OpRead          StepOp = "read"           // reading from an existing resource.
	OpUpdate        StepOp = "update"         // updating an existing resource.
	OpDelete        StepOp = "delete"         // deleting an existing resource.
	OpReplace       StepOp = "replace"        // replacing a resource with a new one (logically).
	OpReplaceCreate StepOp = "replace-create" // the fine-grained replacement step to create the new resource.
	OpReplaceDelete StepOp = "replace-delete" // the fine-grained replacement step to delete the old resource.
)

// Color returns a suggested color for lines of this op type.
func (op StepOp) Color() string {
	switch op {
	case OpCreate, OpReplaceCreate:
		return colors.SpecAdded
	case OpDelete, OpReplaceDelete:
		return colors.SpecDeleted
	case OpUpdate:
		return colors.SpecChanged
	case OpReplace:
		return colors.SpecReplaced
	default:
		contract.Failf("Unrecognized resource step op: %v", op)
		return ""
	}
}

// Prefix returns a suggested prefix for lines of this op type.
func (op StepOp) Prefix() string {
	switch op {
	case OpCreate:
		return op.Color() + "+ "
	case OpDelete:
		return op.Color() + "- "
	case OpUpdate:
		return op.Color() + "  "
	case OpReplace:
		return op.Color() + "-+"
	case OpReplaceCreate:
		return op.Color() + "~+"
	case OpReplaceDelete:
		return op.Color() + "~-"
	default:
		contract.Failf("Unrecognized resource step op: %v", op)
		return ""
	}
}

// Suffix returns a suggested suffix for lines of this op type.
func (op StepOp) Suffix() string {
	if op == OpUpdate || op == OpReplace {
		return colors.Reset // updates and replacements colorize individual lines
	}
	return ""
}

// StepOps returns the full set of step operation types.
func StepOps() []StepOp {
	return []StepOp{
		OpCreate,
		OpUpdate,
		OpDelete,
		OpReplace,
		OpReplaceCreate,
		OpReplaceDelete,
	}
}

// NewPlan analyzes a resource graph new compared to an optional old resource graph old, and creates a plan
// that will carry out operations necessary to bring the old resource graph in line with the new one.  It is possible
// for old, new, or both to be nil; combinations of these can be used to create different kinds of plans: (1) a creation
// plan from a new snapshot when old doesn't exist (nil), (2) an update plan when both old and new exist, and (3) a
// deletion plan when old exists, but not new, and (4) an "empty plan" when both are nil.
func NewPlan(ctx *Context, old Snapshot, new Snapshot, analyzers []tokens.QName) (Plan, error) {
	return newPlan(ctx, old, new, analyzers)
}

type plan struct {
	ctx       *Context              // this plan's context.
	ns        tokens.QName          // the namespace target being deployed into.
	pkg       tokens.Package        // the package from which this snapshot came.
	args      core.Args             // the arguments used to compile this package.
	first     *step                 // the first step to take.
	replaces  map[URN][]PropertyKey // resources being replaced and their properties.
	unchanged map[Resource]Resource // the resources that are remaining the same without modification.
}

var _ Plan = (*plan)(nil)

func (p *plan) Replaces() map[URN][]PropertyKey  { return p.replaces }
func (p *plan) Unchanged() map[Resource]Resource { return p.unchanged }
func (p *plan) Empty() bool                      { return p.Steps() == nil }

func (p *plan) Steps() Step {
	if p.first == nil {
		return nil
	}
	return p.first
}

// Provider fetches the provider for a given resource, possibly lazily allocating the plugins for it.  If a provider
// could not be found, or an error occurred while creating it, a non-nil error is returned.
func (p *plan) Provider(res Resource) (Provider, error) {
	t := res.Type()
	pkg := t.Package()
	return p.ctx.Provider(pkg)
}

// Apply performs all steps in the plan, calling out to the progress reporting functions as desired.  It returns four
// things: the resulting Snapshot, no matter whether an error occurs or not; an error, if something went wrong; the step
// that failed, if the error is non-nil; and finally the state of the resource modified in the failing step.
func (p *plan) Apply(prog Progress) (Snapshot, Step, State, error) {
	// First go ahead and propagate IDs for unchanged resources.
	for old, new := range p.unchanged {
		contract.Assert(old.HasID())
		contract.Assert(!new.HasID())
		new.SetID(old.ID())
	}

	// Next, walk the plan linked list and apply each step.
	var res []Resource
	var rst State
	var err error

	stepno := 1
	step := p.Steps()
	for step != nil {
		if prog != nil {
			prog.Before(step)
		}

		rst, err = step.Apply()
		if prog != nil {
			prog.After(step, rst, err)
		}

		// If an error occurred, append the old step to the list (and all subsequent steps).  Else, the new one.
		if err != nil {
			old := step.Old()
			glog.V(7).Infof("Plan step #%v failed [%v]; hasold = %v: %v", stepno, step.Op(), old != nil, err)
			if old != nil {
				res = append(res, old)
			}
			rest := step.Next()
			for rest != nil {
				restres := rest.Old()
				glog.V(7).Infof("Plan step #%v rest.old=%v", restres != nil)
				if restres != nil && !step.Logical() {
					res = append(res, restres) // track all remaining physical resources
				}
				rest = rest.Next()
			}
			break
		} else {
			new := step.New()
			glog.V(7).Infof("Plan step #%v succeeded [%v]; hasnew = %v", stepno, step.Op(), new != nil)
			if new != nil && !step.Logical() {
				res = append(res, new) // track all new physical resources
			}
		}

		step = step.Next()
		stepno++
	}

	// Append all the resources that aren't getting modified.
	glog.V(7).Infof("Adding %v unchanged resource(s) to checkpoint", len(p.unchanged))
	for _, unres := range p.unchanged {
		res = append(res, unres)
	}

	// Finally, produce a new snapshot and return the resulting information.
	return p.checkpoint(res), step, rst, err
}

// checkpoint takes the outputs from a plan application and returns it so that it's suitable for persistence.
func (p *plan) checkpoint(resources []Resource) Snapshot {
	// Produce a resource graph and then topsort it.  Store the result of that.
	g := newResourceGraph(resources)
	topverts, err := graph.Topsort(g)
	contract.Assertf(err == nil, "Fatal inability to topsort plan's output resources; checkpoint impossible")
	var tops []Resource
	for _, topvert := range topverts {
		tops = append(tops, topvert.Data().(Resource))
	}
	glog.V(7).Infof("Checkpointing plan application: %v total resources", len(tops))
	return NewSnapshot(p.ctx, p.ns, p.pkg, p.args, tops)
}

// newPlan handles all three cases: (1) a creation plan from a new snapshot when old doesn't exist (nil), (2) an update
// plan when both old and new exist, and (3) a deletion plan when old exists, but not new.
func newPlan(ctx *Context, old Snapshot, new Snapshot, analyzers []tokens.QName) (*plan, error) {
	// Create a new plan builder and then proceed to do some building.
	pb := newPlanBuilder(ctx, old, new)
	if glog.V(7) {
		glog.V(7).Infof("Creating plan with #old=%v #new=%v #analyzers=%v\n",
			len(pb.OldRes), len(pb.NewRes), len(analyzers))
	}

	// Give analyzers a chance to inspect the overall deployment.
	for _, aname := range analyzers {
		analyzer, err := ctx.Analyzer(aname)
		if err != nil {
			return nil, err
		}
		// TODO: we want to use the full package URL, including its SHA1 hash/version/etc.
		failures, err := analyzer.Analyze(pack.PackageURL{Name: new.Pkg().Name()})
		if err != nil {
			return nil, err
		}
		for _, failure := range failures {
			ctx.Diag.Errorf(errors.ErrorAnalyzeFailure, aname, failure.Reason)
		}
	}

	// Initialize the builder's maps used by everything below (olds, news, dependencies).
	for _, old := range pb.OldRes {
		m := old.URN()
		pb.Olds[m] = old
		contract.Assert(old.HasID())
		// Keep track of which dependents exist for all resources.
		for dep := range old.Properties().AllResources() {
			pb.Depends[dep] = append(pb.Depends[dep], m)
		}
	}
	for _, new := range pb.NewRes {
		pb.News[new.URN()] = new
	}

	// Do a quick pass over the new resources and make sure properties pass muster.
	for _, new := range pb.NewRes {
		t := new.Type()
		props := new.Properties()
		urn := new.URN()

		// First ensure that the provider is okay with this resource.
		prov, err := pb.P.Provider(new)
		if err != nil {
			return nil, err
		}
		failures, err := prov.Check(t, props)
		if err != nil {
			return nil, err
		}
		for _, failure := range failures {
			ctx.Diag.Errorf(errors.ErrorResourcePropertyValueInvalid, urn, failure.Property, failure.Reason)
		}

		// Next, give each analyzer -- if any -- a chance to inspect the reosurce too.
		for _, aname := range analyzers {
			analyzer, err := ctx.Analyzer(aname)
			if err != nil {
				return nil, err
			}
			failures, err := analyzer.AnalyzeResource(t, props)
			if err != nil {
				return nil, err
			}
			for _, failure := range failures {
				ctx.Diag.Errorf(errors.ErrorAnalyzeResourceFailure, aname, urn, failure.Property, failure.Reason)
			}
		}
	}

	// Here's the real meat of the process: diffing the snapshots, looking for:
	//
	//     - Anything in old but not new is a delete
	//     - Anything in new but not old is a create
	//     - For those things in both new and old, any changed properties imply an update
	//
	// Any property changes that require replacement are applied, recursively, in a cascading manner.

	// First, those things in old but not new, and add them to the delete queue.
	for _, old := range pb.OldRes {
		m := old.URN()
		if _, hasnew := pb.News[m]; !hasnew {
			step := newDeleteStep(pb.P, old)
			pb.Deletes[m] = newPlanVertex(step)
			glog.V(7).Infof("Update plan decided to delete '%v'", m)
		}
	}

	// Next, creates and updates: creates are those in new but not old, and updates are those in both.
	for _, new := range pb.NewRes {
		m := new.URN()
		if old, hasold := pb.Olds[m]; hasold {
			// The resource exists in both new and old; it could be an update.  This resource is an update if one of
			// these two conditions exist: 1) either the old and new properties don't match or 2) the update impact
			// is assessed as having to replace the resource, in which case the ID will change.  This might have a
			// cascading impact on subsequent updates too, since those IDs must trigger recreations, etc.
			contract.Assert(old.Type() == new.Type())
			computed := new.Properties().ReplaceResources(func(r URN) URN {
				if pb.Replace(r) {
					// If the resource is being replaced, simply mangle the URN so that it's different; this value
					// won't actually be used for anything other than the diffing algorithms below.
					r = r.Replace()
					glog.V(7).Infof("Patched resource '%v's URN property: %v", m, r)
				}
				return r
			})
			if !old.Properties().DeepEquals(computed) {
				// See if this update has the effect of deleting and recreating the resource.  If so, we need to make
				// sure to insert the right replacement steps into the graph (a create, replace, and delete).
				// TODO[pulumi/coconut#90]: this should generalize to any property changes.
				prov, err := pb.P.Provider(old)
				if err != nil {
					return nil, err
				}
				replaces, _, err := prov.UpdateImpact(old.ID(), old.Type(), old.Properties(), computed)
				if err != nil {
					return nil, err
				}

				// Now create a step and vertex of the right kind.
				if len(replaces) > 0 {
					// To perform a replacement, create a creation, deletion, and add the appropriate edges.  Namely:
					//
					//     - Replacement depends on creation
					//     - Deletion depends on replacement
					//     - Existing dependencies depend on replacement (ensured through usual update logic)
					//     - Deletion depends on updating all existing dependencies (ensured through usual update logic)
					//
					// This ensures the right sequencing, with the replacement node acting as a juncture in the graph.
					replkeys := make([]PropertyKey, len(replaces))
					for i, repl := range replaces {
						replkeys[i] = PropertyKey(repl)
					}
					pb.Replaces[m] = replkeys
					create := newReplaceCreateStep(pb.P, new)
					pb.Creates[m] = newPlanVertex(create)
					replace := newReplaceStep(pb.P, old, new, computed)
					pb.Updates[m] = newPlanVertex(replace)
					pb.Updates[m].connectTo(pb.Creates[m]) // replacement depends on creation
					delete := newReplaceDeleteStep(pb.P, old)
					pb.Deletes[m] = newPlanVertex(delete)
					pb.Deletes[m].connectTo(pb.Updates[m]) // deletion depends on replacement
					glog.V(7).Infof("Update plan decided to update '%v'; necessitates a replacement", m)
				} else {
					// An update is simple: just create a single update step and associated node in the graph.
					step := newUpdateStep(pb.P, old, new, computed)
					pb.Updates[m] = newPlanVertex(step)
					glog.V(7).Infof("Update plan decided to update '%v'", m)
				}
			} else {
				pb.Unchanged[old] = new
				glog.V(7).Infof("Update plan decided not to update '%v'", m)
			}
		} else {
			// The resource isn't in the old map, so it must be a resource creation.
			step := newCreateStep(pb.P, new)
			pb.Creates[m] = newPlanVertex(step)
			glog.V(7).Infof("Update plan decided to create '%v'", m)
		}
	}

	// Finally, we need to sequence the overall set of changes to create the final plan.  To do this, we create a DAG
	// of the above operations, so that inherent dependencies between operations are respected; specifically:
	//
	//     - Deleting a resource depends on deletes of dependents and updates whose olds refer to it
	//     - Creating a resource depends on creates of dependencies
	//     - Updating a resource depends on creates or updates of news
	//
	// Clearly we must prohibit cycles in this overall graph of resource operations (hence the DAG part).  To ensure
	// this ordering, we will produce a plan graph whose vertices are operations and whose edges encode dependencies.
	for _, old := range pb.OldRes {
		m := old.URN()
		if delete, isdelete := pb.Deletes[m]; isdelete {
			pb.ConnectDelete(m, delete) // connect this delete so it happens before dependencies.
		} else if update, isupdate := pb.Updates[m]; isupdate {
			pb.ConnectUpdate(m, update) // connect this delete so it happens after dependencies are created/updated.
		}
	}
	for _, new := range pb.NewRes {
		m := new.URN()
		if create, iscreate := pb.Creates[m]; iscreate {
			pb.ConnectCreate(m, create) // connect this create so it happens after dependencies are created/updated.
		}
	}

	// Finally, finish the creation of the plan, and return it.
	return pb.Plan()
}

// planBuilder records a lot of the necessary information during the creation of a plan.
type planBuilder struct {
	P         *plan                 // the plan under construction.
	Olds      map[URN]Resource      // a map of URN to old resource.
	OldRes    []Resource            // a flat list of old resources (in topological order).
	News      map[URN]Resource      // a map of URN to new resource.
	NewRes    []Resource            // a flat list of new resources (in topological order).
	Depends   map[URN][]URN         // a map of URN to all existing (old) dependencies.
	Creates   map[URN]*planVertex   // a map of pending creates to their associated vertex.
	Updates   map[URN]*planVertex   // a map of pending updates to their associated vertex.
	Deletes   map[URN]*planVertex   // a map of pending deletes to their associated vertex.
	Replaces  map[URN][]PropertyKey // a map of URNs scheduled for replacement to properties being replaced.
	Unchanged map[Resource]Resource // a map of unchanged resources to their ID-stamped state.
}

// newPlanBuilder initializes a fresh plan state instance, ready to use for planning.
func newPlanBuilder(ctx *Context, old Snapshot, new Snapshot) *planBuilder {
	// These variables are read from either snapshot (preferred new, since it may have updated args).
	var ns tokens.QName
	var pkg tokens.Package
	var args core.Args

	// Now extract the resources and settings from the old and/or new snapshots.
	var oldres []Resource
	if old != nil {
		oldres = old.Resources()
		if new == nil {
			ns = old.Namespace()
			pkg = old.Pkg()
			args = old.Args()
		}
	}
	var newres []Resource
	if new != nil {
		newres = new.Resources()
		ns = new.Namespace()
		pkg = new.Pkg()
		args = new.Args()
	}

	// Create a new, unfinished plan; it will be completed later on after the builder is done.
	p := &plan{
		ctx:  ctx,
		ns:   ns,
		pkg:  pkg,
		args: args,
	}

	return &planBuilder{
		P:         p,
		Olds:      make(map[URN]Resource),
		OldRes:    oldres,
		News:      make(map[URN]Resource),
		NewRes:    newres,
		Depends:   make(map[URN][]URN),
		Creates:   make(map[URN]*planVertex),
		Updates:   make(map[URN]*planVertex),
		Deletes:   make(map[URN]*planVertex),
		Replaces:  make(map[URN][]PropertyKey),
		Unchanged: make(map[Resource]Resource),
	}
}

func (pb *planBuilder) Replace(m URN) bool {
	return len(pb.Replaces[m]) > 0
}

func (pb *planBuilder) ConnectCreate(m URN, v *planVertex) {
	pb.connectCreateUpdate(m, v, false)
}

func (pb *planBuilder) ConnectUpdate(m URN, v *planVertex) {
	pb.connectCreateUpdate(m, v, true)
}

func (pb *planBuilder) connectCreateUpdate(m URN, v *planVertex, update bool) {
	var label string
	if update {
		label = "Updating"
	} else {
		label = "Creating"
	}

	// Add edges to:
	//     - new resources this step depends on
	//     - updated resources that this step depends on
	new := v.Step().New()
	for dep := range new.Properties().AllResources() {
		tov, has := pb.Creates[dep] // see if we're creating the dependency.
		if !has {
			tov, has = pb.Updates[dep] // see if the dependency is being updated.
		}
		if has {
			contract.Assert(tov != nil)
			v.connectTo(tov)
			glog.V(7).Infof("%v '%v' depends on resource '%v'; edge created", label, m, dep)
		} else {
			// A missing entry is ok; it means the resource isn't changing.
			old := pb.Olds[dep]
			contract.Assertf(old != nil, "Expected '%v' to be an existing resource", dep)
			contract.Assertf(pb.Unchanged[old] != nil, "Expected '%v' to be unchanged", dep)
			glog.V(7).Infof("%v '%v' depends on resource '%v'; unchanged, so ignoring", label, m, dep)
		}
	}
}

func (pb *planBuilder) ConnectDelete(m URN, v *planVertex) {
	// Add edges to:
	//     - any dependents that used to refer to this (and are necessarily being deleted or updated)
	for _, dep := range pb.Depends[m] {
		tov, has := pb.Deletes[dep] // see if dependents are being deleted
		if !has {
			tov, has = pb.Updates[dep] // else, they should be updated, otherwise there is a problem
		}
		contract.Assertf(has, "Resource '%v' depends on '%v' (scheduled for deletion)", m, dep)
		contract.Assert(tov != nil)
		v.connectTo(tov)
		glog.V(7).Infof("Deletion '%v' depends on resource '%v'; edge created", m, dep)
	}
}

// Plan finishes the plan building and returns the resulting, completed plan (or non-nil error if it fails).
func (pb *planBuilder) Plan() (*plan, error) {
	// For all plan vertices with no ins, make them root nodes.
	var roots []*planEdge
	for _, vs := range []map[URN]*planVertex{pb.Creates, pb.Updates, pb.Deletes} {
		for _, v := range vs {
			if len(v.Ins()) == 0 {
				roots = append(roots, &planEdge{to: v})
			}
		}
	}

	// Now topologically sort the steps in the order they must execute, thread the plan together, and return it.
	g := newPlanGraph(roots)
	topdag, err := graph.Topsort(g)
	if err != nil {
		return nil, err
	}
	var prev *step
	for _, v := range topdag {
		insertStep(&prev, v.Data().(*step))
	}

	// Remember extra information useful for plan consumers.
	pb.P.replaces = pb.Replaces
	pb.P.unchanged = pb.Unchanged

	return pb.P, nil
}

type step struct {
	p        *plan       // this step's plan.
	op       StepOp      // the operation to perform.
	old      Resource    // the state of the resource before this step.
	new      Resource    // the state of the resource after this step.
	newprops PropertyMap // the resource's properties, factoring in dependency updates.
	next     *step       // the next step after this one in the plan.
}

var _ Step = (*step)(nil)

func (s *step) Plan() Plan            { return s.p }
func (s *step) Op() StepOp            { return s.op }
func (s *step) Logical() bool         { return s.op == OpReplace }
func (s *step) Old() Resource         { return s.old }
func (s *step) New() Resource         { return s.new }
func (s *step) NewProps() PropertyMap { return s.newprops }
func (s *step) Next() Step {
	if s.next == nil {
		return nil
	}
	return s.next
}

func (s *step) Provider() (Provider, error) {
	contract.Assert(s.old == nil || s.new == nil || s.old.Type() == s.new.Type())
	if s.old != nil {
		return s.p.Provider(s.old)
	}
	contract.Assert(s.new != nil)
	return s.p.Provider(s.new)
}

func newCreateStep(p *plan, new Resource) *step {
	return &step{p: p, op: OpCreate, new: new, newprops: new.Properties()}
}

func newDeleteStep(p *plan, old Resource) *step {
	return &step{p: p, op: OpDelete, old: old, newprops: nil}
}

func newUpdateStep(p *plan, old Resource, new Resource, newprops PropertyMap) *step {
	return &step{p: p, op: OpUpdate, old: old, new: new, newprops: newprops}
}

func newReplaceStep(p *plan, old Resource, new Resource, newprops PropertyMap) *step {
	return &step{p: p, op: OpReplace, old: old, new: new, newprops: newprops}
}

func newReplaceCreateStep(p *plan, new Resource) *step {
	return &step{p: p, op: OpReplaceCreate, new: new, newprops: new.Properties()}
}

func newReplaceDeleteStep(p *plan, old Resource) *step {
	return &step{p: p, op: OpReplaceDelete, old: old, newprops: nil}
}

func insertStep(prev **step, step *step) {
	contract.Assert(prev != nil)
	if *prev == nil {
		contract.Assert(step.p.first == nil)
		step.p.first = step
		*prev = step
	} else {
		(*prev).next = step
		*prev = step
	}
}

func (s *step) Apply() (State, error) {
	// Fetch the provider.
	prov, err := s.Provider()
	if err != nil {
		return StateOK, err
	}

	// Now simply perform the operation of the right kind.
	switch s.op {
	case OpCreate, OpReplaceCreate:
		// Invoke the Create RPC function for this provider:
		contract.Assert(s.old == nil)
		contract.Assert(s.new != nil)
		contract.Assertf(!s.new.HasID(), "Resources being created must not have IDs already")
		id, rst, err := prov.Create(s.new.Type(), s.new.Properties())
		if err != nil {
			return rst, err
		}
		s.new.SetID(id)

	case OpDelete, OpReplaceDelete:
		// Invoke the Delete RPC function for this provider:
		contract.Assert(s.old != nil)
		contract.Assert(s.new == nil)
		contract.Assertf(s.old.HasID(), "Resources being deleted must have IDs")
		if rst, err := prov.Delete(s.old.ID(), s.old.Type()); err != nil {
			return rst, err
		}

	case OpUpdate:
		// Invoke the Update RPC function for this provider:
		contract.Assert(s.old != nil)
		contract.Assert(s.new != nil)
		contract.Assert(s.old.Type() == s.new.Type())
		contract.Assertf(s.old.HasID(), "Resources being updated must have IDs")
		id := s.old.ID()
		if rst, err := prov.Update(id, s.old.Type(), s.old.Properties(), s.new.Properties()); err != nil {
			return rst, err
		}

		// Propagate the old ID on the new resource, so that the resulting snapshot is correct.
		s.new.SetID(id)

	case OpReplace:
		contract.Assert(s.old != nil)
		contract.Assert(s.new != nil)
		contract.Assert(s.old.Type() == s.new.Type())
		contract.Assertf(s.old.HasID(), "Resources being replaced must have IDs")

		// There is nothing to do for OpReplace nodes; they are here to represent logical steps in the graph, and mostly
		// for visualization purposes; there will be true OpReplaceCreate and OpReplaceDelete nodes in the graph.

	default:
		contract.Failf("Unexpected step operation: %v", s.op)
	}

	return StateOK, nil
}
