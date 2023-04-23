package clickhouse

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

type Statement struct {
	execSQL    strings.Builder
	sqlBuilder StatementSQLBuilder
	conditions []*wrappedSQLBuilder
}

func (s *Statement) AddConditionClause(buildType SQLBuilderType, args ...interface{}) {
	switch buildType {
	case PrepareQueryBuilder:
		s.conditions = append(s.conditions, &wrappedSQLBuilder{
			args:           args,
			SQLBuilderType: PrepareQueryBuilder,
		})
		return
	case ArrayJoinBuilder:
		s.conditions = append(s.conditions, &wrappedSQLBuilder{
			args:           args,
			SQLBuilderType: ArrayJoinBuilder,
		})
		return
	case WhereBuilder:
		s.conditions = append(s.conditions, &wrappedSQLBuilder{
			args:           args,
			SQLBuilderType: WhereBuilder,
		})
		return
	case InBuilder:
		s.conditions = append(s.conditions, &wrappedSQLBuilder{
			args:           args,
			SQLBuilderType: InBuilder,
		})
		return
	case GroupByBuilder:
		s.conditions = append(s.conditions, &wrappedSQLBuilder{
			args:           args,
			SQLBuilderType: GroupByBuilder,
		})
		return
	default:

	}
}

func (s *Statement) BuildCondition() error {
	if len(s.conditions) == 0 {
		return fmt.Errorf("build conditions is nil")
	}

	for _, condition := range s.conditions {
		if err := s.buildCondition(condition.SQLBuilderType, condition.args...); err != nil {
			return err
		}
	}

	return nil
}

/*
[WITH expr_list|(subquery)]
SELECT [DISTINCT] expr_list
[FROM [db.]table | (subquery) | table_function] [FINAL]
[SAMPLE sample_coeff]
[ARRAY JOIN ...]
[GLOBAL] [ANY|ALL|ASOF] [INNER|LEFT|RIGHT|FULL|CROSS] [OUTER|SEMI|ANTI] JOIN (subquery)|table (ON <expr_list>)|(USING <column_list>)
[PREWHERE expr]
[WHERE expr]
[GROUP BY expr_list] [WITH TOTALS]
[HAVING expr]
[ORDER BY expr_list] [WITH FILL] [FROM expr] [TO expr] [STEP expr]
[LIMIT [offset_value, ]n BY columns]
[LIMIT [n, ]m] [WITH TIES]
[UNION ALL ...]
[INTO OUTFILE filename]
[FORMAT format]
*/
func (s *Statement) buildCondition(buildType SQLBuilderType, args ...interface{}) error {
	switch buildType {
	case PrepareQueryBuilder:
		prepareQuerySQl, err := s.sqlBuilder.PrepareQuery(args...)
		if err != nil {
			return err
		}
		s.execSQL.WriteString(prepareQuerySQl)
		if !strings.HasSuffix(s.execSQL.String(), " ") {
			s.execSQL.WriteString(" ")
		}
		return nil
	case ArrayJoinBuilder:
		arrayJoinSQL, err := s.sqlBuilder.ArrayJoin(args...)
		if err != nil {
			return err
		}
		s.execSQL.WriteString(arrayJoinSQL)
		if !strings.HasSuffix(s.execSQL.String(), " ") {
			s.execSQL.WriteString(" ")
		}
		return nil
	case WhereBuilder:
		whereSQL, err := s.sqlBuilder.Where(args...)
		if err != nil {
			return err
		}

		if !strings.Contains(s.execSQL.String(), "where") {
			s.execSQL.WriteString("where ")
		} else {
			s.execSQL.WriteString("and ")
		}
		s.execSQL.WriteString(whereSQL)
		if !strings.HasSuffix(s.execSQL.String(), " ") {
			s.execSQL.WriteString(" ")
		}
		return nil
	case InBuilder:
		inSQL, err := s.sqlBuilder.In(args...)
		if err != nil {
			return err
		}

		if !strings.Contains(s.execSQL.String(), "where") {
			s.execSQL.WriteString("where ")
		} else {
			s.execSQL.WriteString("and ")
		}
		s.execSQL.WriteString(inSQL)
		if !strings.HasSuffix(s.execSQL.String(), " ") {
			s.execSQL.WriteString(" ")
		}
		return nil
	case GroupByBuilder:
		groupBySQL, err := s.sqlBuilder.GroupByBuilder(args...)
		if err != nil {
			return err
		}

		s.execSQL.WriteString(groupBySQL)
		if !strings.HasSuffix(s.execSQL.String(), " ") {
			s.execSQL.WriteString(" ")
		}
		return nil
	default:
		return fmt.Errorf("unsupport sql condition type")
	}
}

type SQLBuilderType int

const (
	PrepareQueryBuilder SQLBuilderType = iota
	WhereBuilder
	ArrayJoinBuilder
	InBuilder
	GroupByBuilder
)

type StatementSQLBuilder interface {
	PrepareQuery(args ...interface{}) (string, error)
	Where(args ...interface{}) (string, error)
	ArrayJoin(args ...interface{}) (string, error)
	In(args ...interface{}) (string, error)
	GroupByBuilder(args ...interface{}) (string, error)
}

type chStatementSQLBuilder struct {
	prepareQuerySQLBuilder
	whereSQLBuilder
	arrayJoinSQLBuilder
	inSQLBuilder
	groupBySQLBuilder
}

func (c *chStatementSQLBuilder) PrepareQuery(args ...interface{}) (string, error) {
	return c.prepareQuerySQLBuilder.Build(args...)
}

func (c *chStatementSQLBuilder) Where(args ...interface{}) (string, error) {
	return c.whereSQLBuilder.Build(args...)
}

func (c *chStatementSQLBuilder) ArrayJoin(args ...interface{}) (string, error) {
	return c.arrayJoinSQLBuilder.Build(args...)
}

func (c *chStatementSQLBuilder) In(args ...interface{}) (string, error) {
	return c.inSQLBuilder.Build(args...)
}

func (c *chStatementSQLBuilder) GroupByBuilder(args ...interface{}) (string, error) {
	return c.groupBySQLBuilder.Build(args...)
}

type SQLBuilder interface {
	Build(args ...interface{}) (string, error)
	Type() SQLBuilderType
}

type wrappedSQLBuilder struct {
	args []interface{}
	SQLBuilderType
}

type prepareQuerySQLBuilder struct{}

func (*prepareQuerySQLBuilder) Build(args ...interface{}) (string, error) {
	if args == nil {
		return "", fmt.Errorf("prepareQuery sql is nil")
	}

	if len(args) != 1 {
		return "", fmt.Errorf("prepareQuery sql nums nust be one")
	}

	prepareQuerySQl := args[0]
	k := reflect.TypeOf(prepareQuerySQl).Kind()
	if k != reflect.String {
		return "", fmt.Errorf("prepareQuery arg query kind must be  Stinrg")
	}

	return cast.ToString(prepareQuerySQl), nil
}

func (*prepareQuerySQLBuilder) Type() SQLBuilderType {
	return PrepareQueryBuilder
}

type whereSQLBuilder struct{}

func (*whereSQLBuilder) Build(args ...interface{}) (string, error) {
	if args == nil {
		return "", fmt.Errorf("where args is nil")
	}

	//  a = ? and b <= ?
	whereTemplateSQL := cast.ToString(args[0])
	var whereSQL strings.Builder
	for _, v := range args[1:] {
		i := strings.Index(whereTemplateSQL, "?")
		if i < 0 {
			break
		}
		whereSQL.WriteString(fmt.Sprintf("%s'%s'", whereTemplateSQL[:i], cast.ToString(v)))
		if i < len(whereTemplateSQL) {
			whereTemplateSQL = whereTemplateSQL[i+1:]
		}
	}
	whereSQL.WriteString(whereTemplateSQL)
	return whereSQL.String(), nil
}

func (*whereSQLBuilder) Type() SQLBuilderType {
	return WhereBuilder
}

type arrayJoinSQLBuilder struct {
}

func (*arrayJoinSQLBuilder) Build(args ...interface{}) (string, error) {
	if args == nil {
		return "", fmt.Errorf("arrayJoin filed is nil")
	}

	if len(args) != 1 {
		return "", fmt.Errorf("arrayJoin fileds nums nust be one")
	}

	filed := args[0]
	if filed == nil {
		return "", fmt.Errorf("arrayJoin filed is nil")
	}

	k := reflect.TypeOf(filed).Kind()
	if k != reflect.String {
		return "", fmt.Errorf("build arg filed kind  must be  Stinrg")
	}
	return fmt.Sprintf("array join `%s` ", cast.ToString(filed)), nil
}

func (*arrayJoinSQLBuilder) Type() SQLBuilderType {
	return ArrayJoinBuilder
}

type inSQLBuilder struct{}

func (*inSQLBuilder) Build(args ...interface{}) (string, error) {
	if args == nil {
		return "", fmt.Errorf("in args is nil")
	}

	l := len(args)
	if l&(l-1) != 0 {
		return "", fmt.Errorf("in args nums must power of two")
	}

	parseFun := func(filed, values interface{}) (string, error) {
		f := cast.ToString(filed)
		if reflect.TypeOf(values).Kind() != reflect.Slice {
			return fmt.Sprintf("`%s` in ('%s')", f, values), nil
		}

		var inSubSQL strings.Builder
		inSubSQL.WriteString(fmt.Sprintf("`%s` in (", f))
		vales := convertToInterfaceSlice(values)
		for i, v := range vales {
			inSubSQL.WriteString(fmt.Sprintf("'%v'", v))
			if i != len(vales)-1 {
				inSubSQL.WriteString(",")
			}
		}

		inSubSQL.WriteString(")")
		return inSubSQL.String(), nil
	}
	var inSQL strings.Builder
	for i := 0; i < l/2; i++ {
		filed := args[i]
		values := args[l/2+i]
		inSubSQL, err := parseFun(filed, values)
		if err != nil {
			return "", err
		}
		inSQL.WriteString(inSubSQL)
		if i != l/2-1 {
			inSQL.WriteString(" and ")
		}
	}

	return inSQL.String(), nil
}

func (*inSQLBuilder) Type() SQLBuilderType {
	return InBuilder
}

type groupBySQLBuilder struct{}

func (*groupBySQLBuilder) Build(args ...interface{}) (string, error) {
	if args == nil {
		return "", fmt.Errorf("groupBy fileds is nil")
	}

	if reflect.TypeOf(args).Kind() != reflect.Slice {
		return fmt.Sprintf("group by `%s` ", cast.ToString(args)), nil
	}

	var groupBySQL strings.Builder
	groupBySQL.WriteString("group by ")
	for i, v := range args {
		groupBySQL.WriteString(fmt.Sprintf("`%v`", v))
		if i != len(args)-1 {
			groupBySQL.WriteString(",")
		}
	}
	return groupBySQL.String(), nil
}

func (*groupBySQLBuilder) Type() SQLBuilderType {
	return GroupByBuilder
}

func convertToInterfaceSlice(slice interface{}) []interface{} {
	sliceValue := reflect.ValueOf(slice)

	result := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		elemValue := sliceValue.Index(i)
		elemInterfaceValue := elemValue.Interface()
		result[i] = elemInterfaceValue
	}

	return result
}
