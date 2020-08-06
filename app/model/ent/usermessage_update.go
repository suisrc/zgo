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
	"github.com/suisrc/zgo/app/model/ent/usermessage"
)

// UserMessageUpdate is the builder for updating UserMessage entities.
type UserMessageUpdate struct {
	config
	hooks      []Hook
	mutation   *UserMessageMutation
	predicates []predicate.UserMessage
}

// Where adds a new predicate for the builder.
func (umu *UserMessageUpdate) Where(ps ...predicate.UserMessage) *UserMessageUpdate {
	umu.predicates = append(umu.predicates, ps...)
	return umu
}

// SetUID sets the uid field.
func (umu *UserMessageUpdate) SetUID(s string) *UserMessageUpdate {
	umu.mutation.SetUID(s)
	return umu
}

// SetAvatar sets the avatar field.
func (umu *UserMessageUpdate) SetAvatar(s string) *UserMessageUpdate {
	umu.mutation.SetAvatar(s)
	return umu
}

// SetTitle sets the title field.
func (umu *UserMessageUpdate) SetTitle(s string) *UserMessageUpdate {
	umu.mutation.SetTitle(s)
	return umu
}

// SetDatetime sets the datetime field.
func (umu *UserMessageUpdate) SetDatetime(s string) *UserMessageUpdate {
	umu.mutation.SetDatetime(s)
	return umu
}

// SetType sets the type field.
func (umu *UserMessageUpdate) SetType(s string) *UserMessageUpdate {
	umu.mutation.SetType(s)
	return umu
}

// SetRead sets the read field.
func (umu *UserMessageUpdate) SetRead(i int) *UserMessageUpdate {
	umu.mutation.ResetRead()
	umu.mutation.SetRead(i)
	return umu
}

// AddRead adds i to read.
func (umu *UserMessageUpdate) AddRead(i int) *UserMessageUpdate {
	umu.mutation.AddRead(i)
	return umu
}

// SetDescription sets the description field.
func (umu *UserMessageUpdate) SetDescription(s string) *UserMessageUpdate {
	umu.mutation.SetDescription(s)
	return umu
}

// SetClickClose sets the clickClose field.
func (umu *UserMessageUpdate) SetClickClose(i int) *UserMessageUpdate {
	umu.mutation.ResetClickClose()
	umu.mutation.SetClickClose(i)
	return umu
}

// AddClickClose adds i to clickClose.
func (umu *UserMessageUpdate) AddClickClose(i int) *UserMessageUpdate {
	umu.mutation.AddClickClose(i)
	return umu
}

// SetStatus sets the status field.
func (umu *UserMessageUpdate) SetStatus(i int) *UserMessageUpdate {
	umu.mutation.ResetStatus()
	umu.mutation.SetStatus(i)
	return umu
}

// AddStatus adds i to status.
func (umu *UserMessageUpdate) AddStatus(i int) *UserMessageUpdate {
	umu.mutation.AddStatus(i)
	return umu
}

// SetCreator sets the creator field.
func (umu *UserMessageUpdate) SetCreator(s string) *UserMessageUpdate {
	umu.mutation.SetCreator(s)
	return umu
}

// SetCreatedAt sets the created_at field.
func (umu *UserMessageUpdate) SetCreatedAt(t time.Time) *UserMessageUpdate {
	umu.mutation.SetCreatedAt(t)
	return umu
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (umu *UserMessageUpdate) SetNillableCreatedAt(t *time.Time) *UserMessageUpdate {
	if t != nil {
		umu.SetCreatedAt(*t)
	}
	return umu
}

// SetUpdatedAt sets the updated_at field.
func (umu *UserMessageUpdate) SetUpdatedAt(t time.Time) *UserMessageUpdate {
	umu.mutation.SetUpdatedAt(t)
	return umu
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (umu *UserMessageUpdate) SetNillableUpdatedAt(t *time.Time) *UserMessageUpdate {
	if t != nil {
		umu.SetUpdatedAt(*t)
	}
	return umu
}

// SetVersion sets the version field.
func (umu *UserMessageUpdate) SetVersion(i int) *UserMessageUpdate {
	umu.mutation.ResetVersion()
	umu.mutation.SetVersion(i)
	return umu
}

// SetNillableVersion sets the version field if the given value is not nil.
func (umu *UserMessageUpdate) SetNillableVersion(i *int) *UserMessageUpdate {
	if i != nil {
		umu.SetVersion(*i)
	}
	return umu
}

// AddVersion adds i to version.
func (umu *UserMessageUpdate) AddVersion(i int) *UserMessageUpdate {
	umu.mutation.AddVersion(i)
	return umu
}

// Mutation returns the UserMessageMutation object of the builder.
func (umu *UserMessageUpdate) Mutation() *UserMessageMutation {
	return umu.mutation
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (umu *UserMessageUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(umu.hooks) == 0 {
		affected, err = umu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			umu.mutation = mutation
			affected, err = umu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(umu.hooks) - 1; i >= 0; i-- {
			mut = umu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, umu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (umu *UserMessageUpdate) SaveX(ctx context.Context) int {
	affected, err := umu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (umu *UserMessageUpdate) Exec(ctx context.Context) error {
	_, err := umu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (umu *UserMessageUpdate) ExecX(ctx context.Context) {
	if err := umu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (umu *UserMessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usermessage.Table,
			Columns: usermessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usermessage.FieldID,
			},
		},
	}
	if ps := umu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := umu.mutation.UID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldUID,
		})
	}
	if value, ok := umu.mutation.Avatar(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldAvatar,
		})
	}
	if value, ok := umu.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldTitle,
		})
	}
	if value, ok := umu.mutation.Datetime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldDatetime,
		})
	}
	if value, ok := umu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldType,
		})
	}
	if value, ok := umu.mutation.Read(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldRead,
		})
	}
	if value, ok := umu.mutation.AddedRead(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldRead,
		})
	}
	if value, ok := umu.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldDescription,
		})
	}
	if value, ok := umu.mutation.ClickClose(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldClickClose,
		})
	}
	if value, ok := umu.mutation.AddedClickClose(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldClickClose,
		})
	}
	if value, ok := umu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldStatus,
		})
	}
	if value, ok := umu.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldStatus,
		})
	}
	if value, ok := umu.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldCreator,
		})
	}
	if value, ok := umu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usermessage.FieldCreatedAt,
		})
	}
	if value, ok := umu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usermessage.FieldUpdatedAt,
		})
	}
	if value, ok := umu.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldVersion,
		})
	}
	if value, ok := umu.mutation.AddedVersion(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldVersion,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, umu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usermessage.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// UserMessageUpdateOne is the builder for updating a single UserMessage entity.
type UserMessageUpdateOne struct {
	config
	hooks    []Hook
	mutation *UserMessageMutation
}

// SetUID sets the uid field.
func (umuo *UserMessageUpdateOne) SetUID(s string) *UserMessageUpdateOne {
	umuo.mutation.SetUID(s)
	return umuo
}

// SetAvatar sets the avatar field.
func (umuo *UserMessageUpdateOne) SetAvatar(s string) *UserMessageUpdateOne {
	umuo.mutation.SetAvatar(s)
	return umuo
}

// SetTitle sets the title field.
func (umuo *UserMessageUpdateOne) SetTitle(s string) *UserMessageUpdateOne {
	umuo.mutation.SetTitle(s)
	return umuo
}

// SetDatetime sets the datetime field.
func (umuo *UserMessageUpdateOne) SetDatetime(s string) *UserMessageUpdateOne {
	umuo.mutation.SetDatetime(s)
	return umuo
}

// SetType sets the type field.
func (umuo *UserMessageUpdateOne) SetType(s string) *UserMessageUpdateOne {
	umuo.mutation.SetType(s)
	return umuo
}

// SetRead sets the read field.
func (umuo *UserMessageUpdateOne) SetRead(i int) *UserMessageUpdateOne {
	umuo.mutation.ResetRead()
	umuo.mutation.SetRead(i)
	return umuo
}

// AddRead adds i to read.
func (umuo *UserMessageUpdateOne) AddRead(i int) *UserMessageUpdateOne {
	umuo.mutation.AddRead(i)
	return umuo
}

// SetDescription sets the description field.
func (umuo *UserMessageUpdateOne) SetDescription(s string) *UserMessageUpdateOne {
	umuo.mutation.SetDescription(s)
	return umuo
}

// SetClickClose sets the clickClose field.
func (umuo *UserMessageUpdateOne) SetClickClose(i int) *UserMessageUpdateOne {
	umuo.mutation.ResetClickClose()
	umuo.mutation.SetClickClose(i)
	return umuo
}

// AddClickClose adds i to clickClose.
func (umuo *UserMessageUpdateOne) AddClickClose(i int) *UserMessageUpdateOne {
	umuo.mutation.AddClickClose(i)
	return umuo
}

// SetStatus sets the status field.
func (umuo *UserMessageUpdateOne) SetStatus(i int) *UserMessageUpdateOne {
	umuo.mutation.ResetStatus()
	umuo.mutation.SetStatus(i)
	return umuo
}

// AddStatus adds i to status.
func (umuo *UserMessageUpdateOne) AddStatus(i int) *UserMessageUpdateOne {
	umuo.mutation.AddStatus(i)
	return umuo
}

// SetCreator sets the creator field.
func (umuo *UserMessageUpdateOne) SetCreator(s string) *UserMessageUpdateOne {
	umuo.mutation.SetCreator(s)
	return umuo
}

// SetCreatedAt sets the created_at field.
func (umuo *UserMessageUpdateOne) SetCreatedAt(t time.Time) *UserMessageUpdateOne {
	umuo.mutation.SetCreatedAt(t)
	return umuo
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (umuo *UserMessageUpdateOne) SetNillableCreatedAt(t *time.Time) *UserMessageUpdateOne {
	if t != nil {
		umuo.SetCreatedAt(*t)
	}
	return umuo
}

// SetUpdatedAt sets the updated_at field.
func (umuo *UserMessageUpdateOne) SetUpdatedAt(t time.Time) *UserMessageUpdateOne {
	umuo.mutation.SetUpdatedAt(t)
	return umuo
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (umuo *UserMessageUpdateOne) SetNillableUpdatedAt(t *time.Time) *UserMessageUpdateOne {
	if t != nil {
		umuo.SetUpdatedAt(*t)
	}
	return umuo
}

// SetVersion sets the version field.
func (umuo *UserMessageUpdateOne) SetVersion(i int) *UserMessageUpdateOne {
	umuo.mutation.ResetVersion()
	umuo.mutation.SetVersion(i)
	return umuo
}

// SetNillableVersion sets the version field if the given value is not nil.
func (umuo *UserMessageUpdateOne) SetNillableVersion(i *int) *UserMessageUpdateOne {
	if i != nil {
		umuo.SetVersion(*i)
	}
	return umuo
}

// AddVersion adds i to version.
func (umuo *UserMessageUpdateOne) AddVersion(i int) *UserMessageUpdateOne {
	umuo.mutation.AddVersion(i)
	return umuo
}

// Mutation returns the UserMessageMutation object of the builder.
func (umuo *UserMessageUpdateOne) Mutation() *UserMessageMutation {
	return umuo.mutation
}

// Save executes the query and returns the updated entity.
func (umuo *UserMessageUpdateOne) Save(ctx context.Context) (*UserMessage, error) {
	var (
		err  error
		node *UserMessage
	)
	if len(umuo.hooks) == 0 {
		node, err = umuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			umuo.mutation = mutation
			node, err = umuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(umuo.hooks) - 1; i >= 0; i-- {
			mut = umuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, umuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (umuo *UserMessageUpdateOne) SaveX(ctx context.Context) *UserMessage {
	um, err := umuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return um
}

// Exec executes the query on the entity.
func (umuo *UserMessageUpdateOne) Exec(ctx context.Context) error {
	_, err := umuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (umuo *UserMessageUpdateOne) ExecX(ctx context.Context) {
	if err := umuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (umuo *UserMessageUpdateOne) sqlSave(ctx context.Context) (um *UserMessage, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   usermessage.Table,
			Columns: usermessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: usermessage.FieldID,
			},
		},
	}
	id, ok := umuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing UserMessage.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := umuo.mutation.UID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldUID,
		})
	}
	if value, ok := umuo.mutation.Avatar(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldAvatar,
		})
	}
	if value, ok := umuo.mutation.Title(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldTitle,
		})
	}
	if value, ok := umuo.mutation.Datetime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldDatetime,
		})
	}
	if value, ok := umuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldType,
		})
	}
	if value, ok := umuo.mutation.Read(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldRead,
		})
	}
	if value, ok := umuo.mutation.AddedRead(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldRead,
		})
	}
	if value, ok := umuo.mutation.Description(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldDescription,
		})
	}
	if value, ok := umuo.mutation.ClickClose(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldClickClose,
		})
	}
	if value, ok := umuo.mutation.AddedClickClose(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldClickClose,
		})
	}
	if value, ok := umuo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldStatus,
		})
	}
	if value, ok := umuo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldStatus,
		})
	}
	if value, ok := umuo.mutation.Creator(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: usermessage.FieldCreator,
		})
	}
	if value, ok := umuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usermessage.FieldCreatedAt,
		})
	}
	if value, ok := umuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: usermessage.FieldUpdatedAt,
		})
	}
	if value, ok := umuo.mutation.Version(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldVersion,
		})
	}
	if value, ok := umuo.mutation.AddedVersion(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: usermessage.FieldVersion,
		})
	}
	um = &UserMessage{config: umuo.config}
	_spec.Assign = um.assignValues
	_spec.ScanValues = um.scanValues()
	if err = sqlgraph.UpdateNode(ctx, umuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{usermessage.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return um, nil
}
