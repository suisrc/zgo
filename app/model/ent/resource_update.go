// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/suisrc/zgo/app/model/ent/predicate"
	"github.com/suisrc/zgo/app/model/ent/resource"
)

// ResourceUpdate is the builder for updating Resource entities.
type ResourceUpdate struct {
	config
	hooks      []Hook
	mutation   *ResourceMutation
	predicates []predicate.Resource
}

// Where adds a new predicate for the builder.
func (ru *ResourceUpdate) Where(ps ...predicate.Resource) *ResourceUpdate {
	ru.predicates = append(ru.predicates, ps...)
	return ru
}

// SetResource sets the resource field.
func (ru *ResourceUpdate) SetResource(s string) *ResourceUpdate {
	ru.mutation.SetResource(s)
	return ru
}

// SetPath sets the path field.
func (ru *ResourceUpdate) SetPath(s string) *ResourceUpdate {
	ru.mutation.SetPath(s)
	return ru
}

// SetNetmask sets the netmask field.
func (ru *ResourceUpdate) SetNetmask(s string) *ResourceUpdate {
	ru.mutation.SetNetmask(s)
	return ru
}

// SetAllow sets the allow field.
func (ru *ResourceUpdate) SetAllow(i int) *ResourceUpdate {
	ru.mutation.ResetAllow()
	ru.mutation.SetAllow(i)
	return ru
}

// AddAllow adds i to allow.
func (ru *ResourceUpdate) AddAllow(i int) *ResourceUpdate {
	ru.mutation.AddAllow(i)
	return ru
}

// SetDesc sets the desc field.
func (ru *ResourceUpdate) SetDesc(s string) *ResourceUpdate {
	ru.mutation.SetDesc(s)
	return ru
}

// SetCreator sets the creator field.
func (ru *ResourceUpdate) SetCreator(s string) *ResourceUpdate {
	ru.mutation.SetCreator(s)
	return ru
}

// SetCreatedAt sets the created_at field.
func (ru *ResourceUpdate) SetCreatedAt(t time.Time) *ResourceUpdate {
	ru.mutation.SetCreatedAt(t)
	return ru
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableCreatedAt(t *time.Time) *ResourceUpdate {
	if t != nil {
		ru.SetCreatedAt(*t)
	}
	return ru
}

// SetUpdatedAt sets the updated_at field.
func (ru *ResourceUpdate) SetUpdatedAt(t time.Time) *ResourceUpdate {
	ru.mutation.SetUpdatedAt(t)
	return ru
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableUpdatedAt(t *time.Time) *ResourceUpdate {
	if t != nil {
		ru.SetUpdatedAt(*t)
	}
	return ru
}

// SetVersion sets the version field.
func (ru *ResourceUpdate) SetVersion(i int) *ResourceUpdate {
	ru.mutation.ResetVersion()
	ru.mutation.SetVersion(i)
	return ru
}

// SetNillableVersion sets the version field if the given value is not nil.
func (ru *ResourceUpdate) SetNillableVersion(i *int) *ResourceUpdate {
	if i != nil {
		ru.SetVersion(*i)
	}
	return ru
}

// AddVersion adds i to version.
func (ru *ResourceUpdate) AddVersion(i int) *ResourceUpdate {
	ru.mutation.AddVersion(i)
	return ru
}

// Mutation returns the ResourceMutation object of the builder.
func (ru *ResourceUpdate) Mutation() *ResourceMutation {
	return ru.mutation
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (ru *ResourceUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(ru.hooks) == 0 {
		affected, err = ru.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ResourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ru.mutation = mutation
			affected, err = ru.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(ru.hooks) - 1; i >= 0; i-- {
			mut = ru.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ru.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (ru *ResourceUpdate) SaveX(ctx context.Context) int {
	affected, err := ru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ru *ResourceUpdate) Exec(ctx context.Context) error {
	_, err := ru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ru *ResourceUpdate) ExecX(ctx context.Context) {
	if err := ru.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ru *ResourceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   resource.Table,
			Columns: resource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: resource.FieldID,
			},
		},
	}
	if ps := ru.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ru.mutation.Resource(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldResource,
		})
	}
	if value, ok := ru.mutation.Path(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldPath,
		})
	}
	if value, ok := ru.mutation.Netmask(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldNetmask,
		})
	}
	if value, ok := ru.mutation.Allow(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldAllow,
		})
	}
	if value, ok := ru.mutation.AddedAllow(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldAllow,
		})
	}
	if value, ok := ru.mutation.Desc(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldDesc,
		})
	}
	if value, ok := ru.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldCreator,
		})
	}
	if value, ok := ru.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: resource.FieldCreatedAt,
		})
	}
	if value, ok := ru.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: resource.FieldUpdatedAt,
		})
	}
	if value, ok := ru.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldVersion,
		})
	}
	if value, ok := ru.mutation.AddedVersion(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldVersion,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resource.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// ResourceUpdateOne is the builder for updating a single Resource entity.
type ResourceUpdateOne struct {
	config
	hooks    []Hook
	mutation *ResourceMutation
}

// SetResource sets the resource field.
func (ruo *ResourceUpdateOne) SetResource(s string) *ResourceUpdateOne {
	ruo.mutation.SetResource(s)
	return ruo
}

// SetPath sets the path field.
func (ruo *ResourceUpdateOne) SetPath(s string) *ResourceUpdateOne {
	ruo.mutation.SetPath(s)
	return ruo
}

// SetNetmask sets the netmask field.
func (ruo *ResourceUpdateOne) SetNetmask(s string) *ResourceUpdateOne {
	ruo.mutation.SetNetmask(s)
	return ruo
}

// SetAllow sets the allow field.
func (ruo *ResourceUpdateOne) SetAllow(i int) *ResourceUpdateOne {
	ruo.mutation.ResetAllow()
	ruo.mutation.SetAllow(i)
	return ruo
}

// AddAllow adds i to allow.
func (ruo *ResourceUpdateOne) AddAllow(i int) *ResourceUpdateOne {
	ruo.mutation.AddAllow(i)
	return ruo
}

// SetDesc sets the desc field.
func (ruo *ResourceUpdateOne) SetDesc(s string) *ResourceUpdateOne {
	ruo.mutation.SetDesc(s)
	return ruo
}

// SetCreator sets the creator field.
func (ruo *ResourceUpdateOne) SetCreator(s string) *ResourceUpdateOne {
	ruo.mutation.SetCreator(s)
	return ruo
}

// SetCreatedAt sets the created_at field.
func (ruo *ResourceUpdateOne) SetCreatedAt(t time.Time) *ResourceUpdateOne {
	ruo.mutation.SetCreatedAt(t)
	return ruo
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableCreatedAt(t *time.Time) *ResourceUpdateOne {
	if t != nil {
		ruo.SetCreatedAt(*t)
	}
	return ruo
}

// SetUpdatedAt sets the updated_at field.
func (ruo *ResourceUpdateOne) SetUpdatedAt(t time.Time) *ResourceUpdateOne {
	ruo.mutation.SetUpdatedAt(t)
	return ruo
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableUpdatedAt(t *time.Time) *ResourceUpdateOne {
	if t != nil {
		ruo.SetUpdatedAt(*t)
	}
	return ruo
}

// SetVersion sets the version field.
func (ruo *ResourceUpdateOne) SetVersion(i int) *ResourceUpdateOne {
	ruo.mutation.ResetVersion()
	ruo.mutation.SetVersion(i)
	return ruo
}

// SetNillableVersion sets the version field if the given value is not nil.
func (ruo *ResourceUpdateOne) SetNillableVersion(i *int) *ResourceUpdateOne {
	if i != nil {
		ruo.SetVersion(*i)
	}
	return ruo
}

// AddVersion adds i to version.
func (ruo *ResourceUpdateOne) AddVersion(i int) *ResourceUpdateOne {
	ruo.mutation.AddVersion(i)
	return ruo
}

// Mutation returns the ResourceMutation object of the builder.
func (ruo *ResourceUpdateOne) Mutation() *ResourceMutation {
	return ruo.mutation
}

// Save executes the query and returns the updated entity.
func (ruo *ResourceUpdateOne) Save(ctx context.Context) (*Resource, error) {
	var (
		err  error
		node *Resource
	)
	if len(ruo.hooks) == 0 {
		node, err = ruo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ResourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			ruo.mutation = mutation
			node, err = ruo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(ruo.hooks) - 1; i >= 0; i-- {
			mut = ruo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, ruo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (ruo *ResourceUpdateOne) SaveX(ctx context.Context) *Resource {
	r, err := ruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return r
}

// Exec executes the query on the entity.
func (ruo *ResourceUpdateOne) Exec(ctx context.Context) error {
	_, err := ruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ruo *ResourceUpdateOne) ExecX(ctx context.Context) {
	if err := ruo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (ruo *ResourceUpdateOne) sqlSave(ctx context.Context) (r *Resource, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   resource.Table,
			Columns: resource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: resource.FieldID,
			},
		},
	}
	id, ok := ruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Resource.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := ruo.mutation.Resource(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldResource,
		})
	}
	if value, ok := ruo.mutation.Path(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldPath,
		})
	}
	if value, ok := ruo.mutation.Netmask(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldNetmask,
		})
	}
	if value, ok := ruo.mutation.Allow(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldAllow,
		})
	}
	if value, ok := ruo.mutation.AddedAllow(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldAllow,
		})
	}
	if value, ok := ruo.mutation.Desc(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldDesc,
		})
	}
	if value, ok := ruo.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: resource.FieldCreator,
		})
	}
	if value, ok := ruo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: resource.FieldCreatedAt,
		})
	}
	if value, ok := ruo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: resource.FieldUpdatedAt,
		})
	}
	if value, ok := ruo.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldVersion,
		})
	}
	if value, ok := ruo.mutation.AddedVersion(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: resource.FieldVersion,
		})
	}
	r = &Resource{config: ruo.config}
	_spec.Assign = r.assignValues
	_spec.ScanValues = r.scanValues()
	if err = sqlgraph.UpdateNode(ctx, ruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{resource.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return r, nil
}
