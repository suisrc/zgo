// Code generated by entc, DO NOT EDIT.

package menurole

import (
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/suisrc/zgo/app/model/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// RoleID applies equality check predicate on the "role_id" field. It's identical to RoleIDEQ.
func RoleID(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRoleID), v))
	})
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUserID), v))
	})
}

// MenuID applies equality check predicate on the "menu_id" field. It's identical to MenuIDEQ.
func MenuID(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMenuID), v))
	})
}

// Creator applies equality check predicate on the "creator" field. It's identical to CreatorEQ.
func Creator(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreator), v))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// RoleIDEQ applies the EQ predicate on the "role_id" field.
func RoleIDEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRoleID), v))
	})
}

// RoleIDNEQ applies the NEQ predicate on the "role_id" field.
func RoleIDNEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRoleID), v))
	})
}

// RoleIDIn applies the In predicate on the "role_id" field.
func RoleIDIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRoleID), v...))
	})
}

// RoleIDNotIn applies the NotIn predicate on the "role_id" field.
func RoleIDNotIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRoleID), v...))
	})
}

// RoleIDGT applies the GT predicate on the "role_id" field.
func RoleIDGT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRoleID), v))
	})
}

// RoleIDGTE applies the GTE predicate on the "role_id" field.
func RoleIDGTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRoleID), v))
	})
}

// RoleIDLT applies the LT predicate on the "role_id" field.
func RoleIDLT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRoleID), v))
	})
}

// RoleIDLTE applies the LTE predicate on the "role_id" field.
func RoleIDLTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRoleID), v))
	})
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUserID), v))
	})
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUserID), v))
	})
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUserID), v...))
	})
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUserID), v...))
	})
}

// UserIDGT applies the GT predicate on the "user_id" field.
func UserIDGT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUserID), v))
	})
}

// UserIDGTE applies the GTE predicate on the "user_id" field.
func UserIDGTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUserID), v))
	})
}

// UserIDLT applies the LT predicate on the "user_id" field.
func UserIDLT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUserID), v))
	})
}

// UserIDLTE applies the LTE predicate on the "user_id" field.
func UserIDLTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUserID), v))
	})
}

// MenuIDEQ applies the EQ predicate on the "menu_id" field.
func MenuIDEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldMenuID), v))
	})
}

// MenuIDNEQ applies the NEQ predicate on the "menu_id" field.
func MenuIDNEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldMenuID), v))
	})
}

// MenuIDIn applies the In predicate on the "menu_id" field.
func MenuIDIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldMenuID), v...))
	})
}

// MenuIDNotIn applies the NotIn predicate on the "menu_id" field.
func MenuIDNotIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldMenuID), v...))
	})
}

// MenuIDGT applies the GT predicate on the "menu_id" field.
func MenuIDGT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldMenuID), v))
	})
}

// MenuIDGTE applies the GTE predicate on the "menu_id" field.
func MenuIDGTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldMenuID), v))
	})
}

// MenuIDLT applies the LT predicate on the "menu_id" field.
func MenuIDLT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldMenuID), v))
	})
}

// MenuIDLTE applies the LTE predicate on the "menu_id" field.
func MenuIDLTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldMenuID), v))
	})
}

// CreatorEQ applies the EQ predicate on the "creator" field.
func CreatorEQ(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreator), v))
	})
}

// CreatorNEQ applies the NEQ predicate on the "creator" field.
func CreatorNEQ(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreator), v))
	})
}

// CreatorIn applies the In predicate on the "creator" field.
func CreatorIn(vs ...string) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreator), v...))
	})
}

// CreatorNotIn applies the NotIn predicate on the "creator" field.
func CreatorNotIn(vs ...string) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreator), v...))
	})
}

// CreatorGT applies the GT predicate on the "creator" field.
func CreatorGT(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreator), v))
	})
}

// CreatorGTE applies the GTE predicate on the "creator" field.
func CreatorGTE(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreator), v))
	})
}

// CreatorLT applies the LT predicate on the "creator" field.
func CreatorLT(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreator), v))
	})
}

// CreatorLTE applies the LTE predicate on the "creator" field.
func CreatorLTE(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreator), v))
	})
}

// CreatorContains applies the Contains predicate on the "creator" field.
func CreatorContains(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldCreator), v))
	})
}

// CreatorHasPrefix applies the HasPrefix predicate on the "creator" field.
func CreatorHasPrefix(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldCreator), v))
	})
}

// CreatorHasSuffix applies the HasSuffix predicate on the "creator" field.
func CreatorHasSuffix(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldCreator), v))
	})
}

// CreatorEqualFold applies the EqualFold predicate on the "creator" field.
func CreatorEqualFold(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldCreator), v))
	})
}

// CreatorContainsFold applies the ContainsFold predicate on the "creator" field.
func CreatorContainsFold(v string) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldCreator), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldVersion), v))
	})
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldVersion), v))
	})
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldVersion), v...))
	})
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...int) predicate.MenuRole {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.MenuRole(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldVersion), v...))
	})
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldVersion), v))
	})
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldVersion), v))
	})
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldVersion), v))
	})
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v int) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldVersion), v))
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.MenuRole) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.MenuRole) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.MenuRole) predicate.MenuRole {
	return predicate.MenuRole(func(s *sql.Selector) {
		p(s.Not())
	})
}