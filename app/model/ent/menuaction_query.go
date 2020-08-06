// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/suisrc/zgo/app/model/ent/menuaction"
	"github.com/suisrc/zgo/app/model/ent/predicate"
)

// MenuActionQuery is the builder for querying MenuAction entities.
type MenuActionQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.MenuAction
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (maq *MenuActionQuery) Where(ps ...predicate.MenuAction) *MenuActionQuery {
	maq.predicates = append(maq.predicates, ps...)
	return maq
}

// Limit adds a limit step to the query.
func (maq *MenuActionQuery) Limit(limit int) *MenuActionQuery {
	maq.limit = &limit
	return maq
}

// Offset adds an offset step to the query.
func (maq *MenuActionQuery) Offset(offset int) *MenuActionQuery {
	maq.offset = &offset
	return maq
}

// Order adds an order step to the query.
func (maq *MenuActionQuery) Order(o ...OrderFunc) *MenuActionQuery {
	maq.order = append(maq.order, o...)
	return maq
}

// First returns the first MenuAction entity in the query. Returns *NotFoundError when no menuaction was found.
func (maq *MenuActionQuery) First(ctx context.Context) (*MenuAction, error) {
	mas, err := maq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(mas) == 0 {
		return nil, &NotFoundError{menuaction.Label}
	}
	return mas[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (maq *MenuActionQuery) FirstX(ctx context.Context) *MenuAction {
	ma, err := maq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return ma
}

// FirstID returns the first MenuAction id in the query. Returns *NotFoundError when no id was found.
func (maq *MenuActionQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = maq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{menuaction.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (maq *MenuActionQuery) FirstXID(ctx context.Context) int {
	id, err := maq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only MenuAction entity in the query, returns an error if not exactly one entity was returned.
func (maq *MenuActionQuery) Only(ctx context.Context) (*MenuAction, error) {
	mas, err := maq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(mas) {
	case 1:
		return mas[0], nil
	case 0:
		return nil, &NotFoundError{menuaction.Label}
	default:
		return nil, &NotSingularError{menuaction.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (maq *MenuActionQuery) OnlyX(ctx context.Context) *MenuAction {
	ma, err := maq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return ma
}

// OnlyID returns the only MenuAction id in the query, returns an error if not exactly one id was returned.
func (maq *MenuActionQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = maq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = &NotSingularError{menuaction.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (maq *MenuActionQuery) OnlyIDX(ctx context.Context) int {
	id, err := maq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of MenuActions.
func (maq *MenuActionQuery) All(ctx context.Context) ([]*MenuAction, error) {
	if err := maq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return maq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (maq *MenuActionQuery) AllX(ctx context.Context) []*MenuAction {
	mas, err := maq.All(ctx)
	if err != nil {
		panic(err)
	}
	return mas
}

// IDs executes the query and returns a list of MenuAction ids.
func (maq *MenuActionQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := maq.Select(menuaction.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (maq *MenuActionQuery) IDsX(ctx context.Context) []int {
	ids, err := maq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (maq *MenuActionQuery) Count(ctx context.Context) (int, error) {
	if err := maq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return maq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (maq *MenuActionQuery) CountX(ctx context.Context) int {
	count, err := maq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (maq *MenuActionQuery) Exist(ctx context.Context) (bool, error) {
	if err := maq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return maq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (maq *MenuActionQuery) ExistX(ctx context.Context) bool {
	exist, err := maq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (maq *MenuActionQuery) Clone() *MenuActionQuery {
	return &MenuActionQuery{
		config:     maq.config,
		limit:      maq.limit,
		offset:     maq.offset,
		order:      append([]OrderFunc{}, maq.order...),
		unique:     append([]string{}, maq.unique...),
		predicates: append([]predicate.MenuAction{}, maq.predicates...),
		// clone intermediate query.
		sql:  maq.sql.Clone(),
		path: maq.path,
	}
}

// GroupBy used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		MenuID int `json:"menu_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.MenuAction.Query().
//		GroupBy(menuaction.FieldMenuID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (maq *MenuActionQuery) GroupBy(field string, fields ...string) *MenuActionGroupBy {
	group := &MenuActionGroupBy{config: maq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := maq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return maq.sqlQuery(), nil
	}
	return group
}

// Select one or more fields from the given query.
//
// Example:
//
//	var v []struct {
//		MenuID int `json:"menu_id,omitempty"`
//	}
//
//	client.MenuAction.Query().
//		Select(menuaction.FieldMenuID).
//		Scan(ctx, &v)
//
func (maq *MenuActionQuery) Select(field string, fields ...string) *MenuActionSelect {
	selector := &MenuActionSelect{config: maq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := maq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return maq.sqlQuery(), nil
	}
	return selector
}

func (maq *MenuActionQuery) prepareQuery(ctx context.Context) error {
	if maq.path != nil {
		prev, err := maq.path(ctx)
		if err != nil {
			return err
		}
		maq.sql = prev
	}
	return nil
}

func (maq *MenuActionQuery) sqlAll(ctx context.Context) ([]*MenuAction, error) {
	var (
		nodes = []*MenuAction{}
		_spec = maq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &MenuAction{config: maq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, maq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (maq *MenuActionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := maq.querySpec()
	return sqlgraph.CountNodes(ctx, maq.driver, _spec)
}

func (maq *MenuActionQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := maq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (maq *MenuActionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   menuaction.Table,
			Columns: menuaction.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: menuaction.FieldID,
			},
		},
		From:   maq.sql,
		Unique: true,
	}
	if ps := maq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := maq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := maq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := maq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (maq *MenuActionQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(maq.driver.Dialect())
	t1 := builder.Table(menuaction.Table)
	selector := builder.Select(t1.Columns(menuaction.Columns...)...).From(t1)
	if maq.sql != nil {
		selector = maq.sql
		selector.Select(selector.Columns(menuaction.Columns...)...)
	}
	for _, p := range maq.predicates {
		p(selector)
	}
	for _, p := range maq.order {
		p(selector)
	}
	if offset := maq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := maq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// MenuActionGroupBy is the builder for group-by MenuAction entities.
type MenuActionGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (magb *MenuActionGroupBy) Aggregate(fns ...AggregateFunc) *MenuActionGroupBy {
	magb.fns = append(magb.fns, fns...)
	return magb
}

// Scan applies the group-by query and scan the result into the given value.
func (magb *MenuActionGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := magb.path(ctx)
	if err != nil {
		return err
	}
	magb.sql = query
	return magb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (magb *MenuActionGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := magb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(magb.fields) > 1 {
		return nil, errors.New("ent: MenuActionGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := magb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (magb *MenuActionGroupBy) StringsX(ctx context.Context) []string {
	v, err := magb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = magb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (magb *MenuActionGroupBy) StringX(ctx context.Context) string {
	v, err := magb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(magb.fields) > 1 {
		return nil, errors.New("ent: MenuActionGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := magb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (magb *MenuActionGroupBy) IntsX(ctx context.Context) []int {
	v, err := magb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = magb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (magb *MenuActionGroupBy) IntX(ctx context.Context) int {
	v, err := magb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(magb.fields) > 1 {
		return nil, errors.New("ent: MenuActionGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := magb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (magb *MenuActionGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := magb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = magb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (magb *MenuActionGroupBy) Float64X(ctx context.Context) float64 {
	v, err := magb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(magb.fields) > 1 {
		return nil, errors.New("ent: MenuActionGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := magb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (magb *MenuActionGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := magb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from group-by. It is only allowed when querying group-by with one field.
func (magb *MenuActionGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = magb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (magb *MenuActionGroupBy) BoolX(ctx context.Context) bool {
	v, err := magb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (magb *MenuActionGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := magb.sqlQuery().Query()
	if err := magb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (magb *MenuActionGroupBy) sqlQuery() *sql.Selector {
	selector := magb.sql
	columns := make([]string, 0, len(magb.fields)+len(magb.fns))
	columns = append(columns, magb.fields...)
	for _, fn := range magb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(magb.fields...)
}

// MenuActionSelect is the builder for select fields of MenuAction entities.
type MenuActionSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (mas *MenuActionSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := mas.path(ctx)
	if err != nil {
		return err
	}
	mas.sql = query
	return mas.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (mas *MenuActionSelect) ScanX(ctx context.Context, v interface{}) {
	if err := mas.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Strings(ctx context.Context) ([]string, error) {
	if len(mas.fields) > 1 {
		return nil, errors.New("ent: MenuActionSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := mas.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (mas *MenuActionSelect) StringsX(ctx context.Context) []string {
	v, err := mas.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = mas.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (mas *MenuActionSelect) StringX(ctx context.Context) string {
	v, err := mas.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Ints(ctx context.Context) ([]int, error) {
	if len(mas.fields) > 1 {
		return nil, errors.New("ent: MenuActionSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := mas.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (mas *MenuActionSelect) IntsX(ctx context.Context) []int {
	v, err := mas.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = mas.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (mas *MenuActionSelect) IntX(ctx context.Context) int {
	v, err := mas.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(mas.fields) > 1 {
		return nil, errors.New("ent: MenuActionSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := mas.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (mas *MenuActionSelect) Float64sX(ctx context.Context) []float64 {
	v, err := mas.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = mas.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (mas *MenuActionSelect) Float64X(ctx context.Context) float64 {
	v, err := mas.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(mas.fields) > 1 {
		return nil, errors.New("ent: MenuActionSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := mas.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (mas *MenuActionSelect) BoolsX(ctx context.Context) []bool {
	v, err := mas.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from selector. It is only allowed when selecting one field.
func (mas *MenuActionSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = mas.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{menuaction.Label}
	default:
		err = fmt.Errorf("ent: MenuActionSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (mas *MenuActionSelect) BoolX(ctx context.Context) bool {
	v, err := mas.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (mas *MenuActionSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := mas.sqlQuery().Query()
	if err := mas.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (mas *MenuActionSelect) sqlQuery() sql.Querier {
	selector := mas.sql
	selector.Select(selector.Columns(mas.fields...)...)
	return selector
}
