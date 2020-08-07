// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/suisrc/zgo/app/model/ent/oauth2third"
)

// Oauth2ThirdCreate is the builder for creating a Oauth2Third entity.
type Oauth2ThirdCreate struct {
	config
	mutation *Oauth2ThirdMutation
	hooks    []Hook
}

// SetPlatform sets the platform field.
func (oc *Oauth2ThirdCreate) SetPlatform(s string) *Oauth2ThirdCreate {
	oc.mutation.SetPlatform(s)
	return oc
}

// SetAgentID sets the agent_id field.
func (oc *Oauth2ThirdCreate) SetAgentID(s string) *Oauth2ThirdCreate {
	oc.mutation.SetAgentID(s)
	return oc
}

// SetSuiteID sets the suite_id field.
func (oc *Oauth2ThirdCreate) SetSuiteID(s string) *Oauth2ThirdCreate {
	oc.mutation.SetSuiteID(s)
	return oc
}

// SetAppID sets the app_id field.
func (oc *Oauth2ThirdCreate) SetAppID(s string) *Oauth2ThirdCreate {
	oc.mutation.SetAppID(s)
	return oc
}

// SetAppSecret sets the app_secret field.
func (oc *Oauth2ThirdCreate) SetAppSecret(s string) *Oauth2ThirdCreate {
	oc.mutation.SetAppSecret(s)
	return oc
}

// SetAuthorizeURL sets the authorize_url field.
func (oc *Oauth2ThirdCreate) SetAuthorizeURL(s string) *Oauth2ThirdCreate {
	oc.mutation.SetAuthorizeURL(s)
	return oc
}

// SetTokenURL sets the token_url field.
func (oc *Oauth2ThirdCreate) SetTokenURL(s string) *Oauth2ThirdCreate {
	oc.mutation.SetTokenURL(s)
	return oc
}

// SetProfileURL sets the profile_url field.
func (oc *Oauth2ThirdCreate) SetProfileURL(s string) *Oauth2ThirdCreate {
	oc.mutation.SetProfileURL(s)
	return oc
}

// SetDomainDef sets the domain_def field.
func (oc *Oauth2ThirdCreate) SetDomainDef(s string) *Oauth2ThirdCreate {
	oc.mutation.SetDomainDef(s)
	return oc
}

// SetDomainCheck sets the domain_check field.
func (oc *Oauth2ThirdCreate) SetDomainCheck(s string) *Oauth2ThirdCreate {
	oc.mutation.SetDomainCheck(s)
	return oc
}

// SetJsSecret sets the js_secret field.
func (oc *Oauth2ThirdCreate) SetJsSecret(s string) *Oauth2ThirdCreate {
	oc.mutation.SetJsSecret(s)
	return oc
}

// SetStateSecret sets the state_secret field.
func (oc *Oauth2ThirdCreate) SetStateSecret(s string) *Oauth2ThirdCreate {
	oc.mutation.SetStateSecret(s)
	return oc
}

// SetCallback sets the callback field.
func (oc *Oauth2ThirdCreate) SetCallback(i int) *Oauth2ThirdCreate {
	oc.mutation.SetCallback(i)
	return oc
}

// SetCbEncrypt sets the cb_encrypt field.
func (oc *Oauth2ThirdCreate) SetCbEncrypt(i int) *Oauth2ThirdCreate {
	oc.mutation.SetCbEncrypt(i)
	return oc
}

// SetCbToken sets the cb_token field.
func (oc *Oauth2ThirdCreate) SetCbToken(s string) *Oauth2ThirdCreate {
	oc.mutation.SetCbToken(s)
	return oc
}

// SetCbEncoding sets the cb_encoding field.
func (oc *Oauth2ThirdCreate) SetCbEncoding(s string) *Oauth2ThirdCreate {
	oc.mutation.SetCbEncoding(s)
	return oc
}

// SetCreator sets the creator field.
func (oc *Oauth2ThirdCreate) SetCreator(s string) *Oauth2ThirdCreate {
	oc.mutation.SetCreator(s)
	return oc
}

// SetCreatedAt sets the created_at field.
func (oc *Oauth2ThirdCreate) SetCreatedAt(t time.Time) *Oauth2ThirdCreate {
	oc.mutation.SetCreatedAt(t)
	return oc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (oc *Oauth2ThirdCreate) SetNillableCreatedAt(t *time.Time) *Oauth2ThirdCreate {
	if t != nil {
		oc.SetCreatedAt(*t)
	}
	return oc
}

// SetUpdatedAt sets the updated_at field.
func (oc *Oauth2ThirdCreate) SetUpdatedAt(t time.Time) *Oauth2ThirdCreate {
	oc.mutation.SetUpdatedAt(t)
	return oc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (oc *Oauth2ThirdCreate) SetNillableUpdatedAt(t *time.Time) *Oauth2ThirdCreate {
	if t != nil {
		oc.SetUpdatedAt(*t)
	}
	return oc
}

// SetVersion sets the version field.
func (oc *Oauth2ThirdCreate) SetVersion(i int) *Oauth2ThirdCreate {
	oc.mutation.SetVersion(i)
	return oc
}

// SetNillableVersion sets the version field if the given value is not nil.
func (oc *Oauth2ThirdCreate) SetNillableVersion(i *int) *Oauth2ThirdCreate {
	if i != nil {
		oc.SetVersion(*i)
	}
	return oc
}

// SetString1 sets the string_1 field.
func (oc *Oauth2ThirdCreate) SetString1(s string) *Oauth2ThirdCreate {
	oc.mutation.SetString1(s)
	return oc
}

// SetString2 sets the string_2 field.
func (oc *Oauth2ThirdCreate) SetString2(s string) *Oauth2ThirdCreate {
	oc.mutation.SetString2(s)
	return oc
}

// SetString3 sets the string_3 field.
func (oc *Oauth2ThirdCreate) SetString3(s string) *Oauth2ThirdCreate {
	oc.mutation.SetString3(s)
	return oc
}

// SetNumber1 sets the number_1 field.
func (oc *Oauth2ThirdCreate) SetNumber1(i int) *Oauth2ThirdCreate {
	oc.mutation.SetNumber1(i)
	return oc
}

// SetNumber2 sets the number_2 field.
func (oc *Oauth2ThirdCreate) SetNumber2(i int) *Oauth2ThirdCreate {
	oc.mutation.SetNumber2(i)
	return oc
}

// SetNumber3 sets the number_3 field.
func (oc *Oauth2ThirdCreate) SetNumber3(i int) *Oauth2ThirdCreate {
	oc.mutation.SetNumber3(i)
	return oc
}

// Mutation returns the Oauth2ThirdMutation object of the builder.
func (oc *Oauth2ThirdCreate) Mutation() *Oauth2ThirdMutation {
	return oc.mutation
}

// Save creates the Oauth2Third in the database.
func (oc *Oauth2ThirdCreate) Save(ctx context.Context) (*Oauth2Third, error) {
	if _, ok := oc.mutation.Platform(); !ok {
		return nil, &ValidationError{Name: "platform", err: errors.New("ent: missing required field \"platform\"")}
	}
	if _, ok := oc.mutation.AgentID(); !ok {
		return nil, &ValidationError{Name: "agent_id", err: errors.New("ent: missing required field \"agent_id\"")}
	}
	if _, ok := oc.mutation.SuiteID(); !ok {
		return nil, &ValidationError{Name: "suite_id", err: errors.New("ent: missing required field \"suite_id\"")}
	}
	if _, ok := oc.mutation.AppID(); !ok {
		return nil, &ValidationError{Name: "app_id", err: errors.New("ent: missing required field \"app_id\"")}
	}
	if _, ok := oc.mutation.AppSecret(); !ok {
		return nil, &ValidationError{Name: "app_secret", err: errors.New("ent: missing required field \"app_secret\"")}
	}
	if _, ok := oc.mutation.AuthorizeURL(); !ok {
		return nil, &ValidationError{Name: "authorize_url", err: errors.New("ent: missing required field \"authorize_url\"")}
	}
	if _, ok := oc.mutation.TokenURL(); !ok {
		return nil, &ValidationError{Name: "token_url", err: errors.New("ent: missing required field \"token_url\"")}
	}
	if _, ok := oc.mutation.ProfileURL(); !ok {
		return nil, &ValidationError{Name: "profile_url", err: errors.New("ent: missing required field \"profile_url\"")}
	}
	if _, ok := oc.mutation.DomainDef(); !ok {
		return nil, &ValidationError{Name: "domain_def", err: errors.New("ent: missing required field \"domain_def\"")}
	}
	if _, ok := oc.mutation.DomainCheck(); !ok {
		return nil, &ValidationError{Name: "domain_check", err: errors.New("ent: missing required field \"domain_check\"")}
	}
	if _, ok := oc.mutation.JsSecret(); !ok {
		return nil, &ValidationError{Name: "js_secret", err: errors.New("ent: missing required field \"js_secret\"")}
	}
	if _, ok := oc.mutation.StateSecret(); !ok {
		return nil, &ValidationError{Name: "state_secret", err: errors.New("ent: missing required field \"state_secret\"")}
	}
	if _, ok := oc.mutation.Callback(); !ok {
		return nil, &ValidationError{Name: "callback", err: errors.New("ent: missing required field \"callback\"")}
	}
	if _, ok := oc.mutation.CbEncrypt(); !ok {
		return nil, &ValidationError{Name: "cb_encrypt", err: errors.New("ent: missing required field \"cb_encrypt\"")}
	}
	if _, ok := oc.mutation.CbToken(); !ok {
		return nil, &ValidationError{Name: "cb_token", err: errors.New("ent: missing required field \"cb_token\"")}
	}
	if _, ok := oc.mutation.CbEncoding(); !ok {
		return nil, &ValidationError{Name: "cb_encoding", err: errors.New("ent: missing required field \"cb_encoding\"")}
	}
	if _, ok := oc.mutation.Creator(); !ok {
		return nil, &ValidationError{Name: "creator", err: errors.New("ent: missing required field \"creator\"")}
	}
	if _, ok := oc.mutation.CreatedAt(); !ok {
		v := oauth2third.DefaultCreatedAt()
		oc.mutation.SetCreatedAt(v)
	}
	if _, ok := oc.mutation.UpdatedAt(); !ok {
		v := oauth2third.DefaultUpdatedAt()
		oc.mutation.SetUpdatedAt(v)
	}
	if _, ok := oc.mutation.Version(); !ok {
		v := oauth2third.DefaultVersion
		oc.mutation.SetVersion(v)
	}
	if _, ok := oc.mutation.String1(); !ok {
		return nil, &ValidationError{Name: "string_1", err: errors.New("ent: missing required field \"string_1\"")}
	}
	if _, ok := oc.mutation.String2(); !ok {
		return nil, &ValidationError{Name: "string_2", err: errors.New("ent: missing required field \"string_2\"")}
	}
	if _, ok := oc.mutation.String3(); !ok {
		return nil, &ValidationError{Name: "string_3", err: errors.New("ent: missing required field \"string_3\"")}
	}
	if _, ok := oc.mutation.Number1(); !ok {
		return nil, &ValidationError{Name: "number_1", err: errors.New("ent: missing required field \"number_1\"")}
	}
	if _, ok := oc.mutation.Number2(); !ok {
		return nil, &ValidationError{Name: "number_2", err: errors.New("ent: missing required field \"number_2\"")}
	}
	if _, ok := oc.mutation.Number3(); !ok {
		return nil, &ValidationError{Name: "number_3", err: errors.New("ent: missing required field \"number_3\"")}
	}
	var (
		err  error
		node *Oauth2Third
	)
	if len(oc.hooks) == 0 {
		node, err = oc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*Oauth2ThirdMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			oc.mutation = mutation
			node, err = oc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(oc.hooks) - 1; i >= 0; i-- {
			mut = oc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, oc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (oc *Oauth2ThirdCreate) SaveX(ctx context.Context) *Oauth2Third {
	v, err := oc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (oc *Oauth2ThirdCreate) sqlSave(ctx context.Context) (*Oauth2Third, error) {
	o, _spec := oc.createSpec()
	if err := sqlgraph.CreateNode(ctx, oc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	o.ID = int(id)
	return o, nil
}

func (oc *Oauth2ThirdCreate) createSpec() (*Oauth2Third, *sqlgraph.CreateSpec) {
	var (
		o     = &Oauth2Third{config: oc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: oauth2third.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: oauth2third.FieldID,
			},
		}
	)
	if value, ok := oc.mutation.Platform(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldPlatform,
		})
		o.Platform = value
	}
	if value, ok := oc.mutation.AgentID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldAgentID,
		})
		o.AgentID = value
	}
	if value, ok := oc.mutation.SuiteID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldSuiteID,
		})
		o.SuiteID = value
	}
	if value, ok := oc.mutation.AppID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldAppID,
		})
		o.AppID = value
	}
	if value, ok := oc.mutation.AppSecret(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldAppSecret,
		})
		o.AppSecret = value
	}
	if value, ok := oc.mutation.AuthorizeURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldAuthorizeURL,
		})
		o.AuthorizeURL = value
	}
	if value, ok := oc.mutation.TokenURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldTokenURL,
		})
		o.TokenURL = value
	}
	if value, ok := oc.mutation.ProfileURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldProfileURL,
		})
		o.ProfileURL = value
	}
	if value, ok := oc.mutation.DomainDef(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldDomainDef,
		})
		o.DomainDef = value
	}
	if value, ok := oc.mutation.DomainCheck(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldDomainCheck,
		})
		o.DomainCheck = value
	}
	if value, ok := oc.mutation.JsSecret(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldJsSecret,
		})
		o.JsSecret = value
	}
	if value, ok := oc.mutation.StateSecret(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldStateSecret,
		})
		o.StateSecret = value
	}
	if value, ok := oc.mutation.Callback(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldCallback,
		})
		o.Callback = value
	}
	if value, ok := oc.mutation.CbEncrypt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldCbEncrypt,
		})
		o.CbEncrypt = value
	}
	if value, ok := oc.mutation.CbToken(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldCbToken,
		})
		o.CbToken = value
	}
	if value, ok := oc.mutation.CbEncoding(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldCbEncoding,
		})
		o.CbEncoding = value
	}
	if value, ok := oc.mutation.Creator(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldCreator,
		})
		o.Creator = value
	}
	if value, ok := oc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: oauth2third.FieldCreatedAt,
		})
		o.CreatedAt = value
	}
	if value, ok := oc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: oauth2third.FieldUpdatedAt,
		})
		o.UpdatedAt = value
	}
	if value, ok := oc.mutation.Version(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldVersion,
		})
		o.Version = value
	}
	if value, ok := oc.mutation.String1(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldString1,
		})
		o.String1 = value
	}
	if value, ok := oc.mutation.String2(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldString2,
		})
		o.String2 = value
	}
	if value, ok := oc.mutation.String3(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: oauth2third.FieldString3,
		})
		o.String3 = value
	}
	if value, ok := oc.mutation.Number1(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldNumber1,
		})
		o.Number1 = value
	}
	if value, ok := oc.mutation.Number2(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldNumber2,
		})
		o.Number2 = value
	}
	if value, ok := oc.mutation.Number3(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: oauth2third.FieldNumber3,
		})
		o.Number3 = value
	}
	return o, _spec
}