/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package db_test

import (
	"database/sql"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/suite"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/sqorc"
)

func TestQuery(t *testing.T) {
	suite.Run(t, &QueryTestSuite{})
}

type QueryTestSuite struct {
	suite.Suite
	resourceManager dbtest.ResourceManager
}

func (s *QueryTestSuite) SetupSuite() {
	builder := sqorc.GetSqlBuilder()
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.resourceManager = dbtest.NewResourceManager(s.T(), database, builder)
	err = s.resourceManager.CreateTables(&someModel{}, &otherModel{}, &anotherModel{}, &modelWithUniqueFields{})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TearDownTest() {
	err := s.resourceManager.DropResources(&someModel{}, &otherModel{}, &anotherModel{})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TestCreate() {
	testCases := []struct {
		name      string
		fieldMask db.FieldMask
		input     db.Model
		expected  db.Model
	}{{
		name:      "Should create resource with required fields",
		fieldMask: db.NewExcludeMask(),
		input:     getSomeModel(),
		expected:  getSomeModel(),
	}, {
		name:      "Should create resource with nullable fields",
		fieldMask: db.NewExcludeMask(),
		input:     getOtherModel(),
		expected:  getOtherModel(),
	}, {
		name:      "Should create resource with null fields",
		fieldMask: db.NewExcludeMask(),
		input:     &otherModel{id: db.MakeInt(id * 2)},
		expected:  &otherModel{id: db.MakeInt(id * 2)},
	}, {
		name:      "Should create resource with default value",
		fieldMask: db.NewExcludeMask("default_value"),
		input:     &anotherModel{id: db.MakeInt(id * 3)},
		expected:  &anotherModel{id: db.MakeInt(id * 3), defaultValue: db.MakeInt(defaultValue)},
	}, {
		name:      "Should create resource with unique fields",
		fieldMask: db.NewExcludeMask(),
		input:     getModelWithUniqueFields(),
		expected:  getModelWithUniqueFields(),
	}}
	for _, tt := range testCases {
		s.Run(tt.name, s.inTransaction(func() {
			id := s.whenModelIsInserted(tt.fieldMask, tt.input)

			actual := s.whenSingleModelIsFetched(id, tt.input)
			s.Assert().Equal(tt.expected, actual)
		}))
	}
}

func (s *QueryTestSuite) TestCreateResourceWithAutogeneratedId() {
	err := s.resourceManager.InTransaction(func() {
		data := getSomeModel()
		data.id = sql.NullInt64{}

		id, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(data).
			Select(db.NewExcludeMask("id")).
			Insert()
		s.Require().NoError(err)

		actual := s.whenSingleModelIsFetched(id, &someModel{})
		expected := getSomeModel().withId(id)
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TestCount() {
	err := s.resourceManager.InsertResources(
		db.NewExcludeMask(),
		getSomeModel(),
		getSomeDifferentModel(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		query    *db.Query
		expected int64
	}{{
		name: "Should count all resources from given table",
		query: db.NewQuery().
			From(&someModel{}),
		expected: 2,
	}, {
		name: "Should count resources with filter",
		query: db.NewQuery().
			From(&someModel{}).
			Where(sq.Eq{"id": id}),
		expected: 1,
	}, {
		name: "Should count resources from empty table",
		query: db.NewQuery().
			From(&otherModel{}),
		expected: 0,
	}}
	for _, tt := range testCases {
		s.Run(tt.name, s.inTransaction(func() {
			actual, err := tt.query.
				WithBuilder(s.resourceManager.GetBuilder()).
				Count()
			s.Require().NoError(err)
			s.Assert().Equal(tt.expected, actual)
		}))
	}
}

func (s *QueryTestSuite) TestUpdate() {
	err := s.resourceManager.InTransaction(func() {
		id := s.whenModelIsInserted(db.NewExcludeMask(), getSomeModel())

		updateData := getSomeDifferentModel()
		err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(updateData).
			Select(db.NewExcludeMask("id")).
			Where(sq.Eq{"id": id}).
			Update()
		s.Require().NoError(err)

		actual := s.whenSingleModelIsFetched(id, &someModel{})
		expected := updateData.withId(id)
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TestUpdateOnlySelectedFields() {
	err := s.resourceManager.InTransaction(func() {
		id := s.whenModelIsInserted(db.NewExcludeMask(), getSomeModel())

		updateData := getSomeDifferentModel()
		err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(updateData).
			Select(db.NewIncludeMask("name", "value")).
			Where(sq.Eq{"id": id}).
			Update()
		s.Require().NoError(err)

		actual := s.whenSingleModelIsFetched(id, &someModel{})
		expected := getSomeModel()
		expected.value = updateData.value
		expected.name = updateData.name
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TestDelete() {
	err := s.resourceManager.InTransaction(func() {
		id := s.whenModelIsInserted(db.NewExcludeMask(), getSomeModel())

		err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&someModel{}).
			Where(sq.Eq{"id": id}).
			Delete()
		s.Require().NoError(err)

		_, err = db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&someModel{}).
			Select(db.NewIncludeMask("id")).
			Where(sq.Eq{"id": id}).
			Fetch()
		s.Assert().ErrorIs(err, sql.ErrNoRows)
	})
	s.Require().NoError(err)
}

func (s *QueryTestSuite) TestFetch() {
	err := s.resourceManager.InsertResources(
		db.NewExcludeMask(),
		getSomeModel(),
		getSomeDifferentModel(),
		getOtherModel(),
		getAnotherModel(),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name     string
		query    *db.Query
		expected []db.Model
	}{{
		name: "Should get only selected fields",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewIncludeMask("id")).
			Where(sq.Eq{"id": id}),
		expected: []db.Model{&someModel{
			id: db.MakeInt(id),
		}},
	}, {
		name: "Should use inner join",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask())),
		expected: []db.Model{
			getSomeModel(),
			getOtherModel(),
		},
	}, {
		name: "Should use left join",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": 2 * id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask()).
				Nullable()),
		expected: []db.Model{
			getSomeDifferentModel(),
			&otherModel{},
		},
	}, {
		name: "Should use filter as join condition",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask()).
				Where(sq.Eq{otherTable + ".id": id - 1}).
				Nullable()),
		expected: []db.Model{
			getSomeModel(),
			&otherModel{},
		},
	}, {
		name: "Should allow nested joins",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask()).
				Join(db.NewQuery().
					From(&anotherModel{}).
					Select(db.NewExcludeMask()))),
		expected: []db.Model{
			getSomeModel(),
			getOtherModel(),
			getAnotherModel(),
		},
	}, {
		name: "Should join nested joins first",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask()).
				Join(db.NewQuery().
					From(&anotherModel{}).
					Select(db.NewExcludeMask()).
					Where(sq.Eq{anotherTable + ".id": id - 1})).
				Nullable()),
		expected: []db.Model{
			getSomeModel(),
			&otherModel{},
			&anotherModel{},
		},
	}, {
		name: "Should use nested left joins",
		query: db.NewQuery().
			From(&someModel{}).
			Select(db.NewExcludeMask()).
			Where(sq.Eq{someTable + ".id": id}).
			Join(db.NewQuery().
				From(&otherModel{}).
				Select(db.NewExcludeMask()).
				Join(db.NewQuery().
					From(&anotherModel{}).
					Select(db.NewExcludeMask()).
					Where(sq.Eq{anotherTable + ".id": id - 1}).
					Nullable()).
				Nullable()),
		expected: []db.Model{
			getSomeModel(),
			getOtherModel(),
			&anotherModel{},
		},
	}}

	for _, tt := range testCases {
		s.Run(tt.name, s.inTransaction(func() {
			actual, err := tt.query.
				WithBuilder(s.resourceManager.GetBuilder()).
				Fetch()
			s.Require().NoError(err)
			s.Assert().Equal(tt.expected, actual)
		}))
	}
}

func (s *QueryTestSuite) TestList() {
	resources := make([]db.Model, 4)
	for i := range resources {
		resources[i] = &otherModel{
			id:    db.MakeInt(int64(i)),
			value: db.MakeFloat(float64(i)),
		}
	}
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), resources...)
	s.Require().NoError(err)

	const column = "value"
	testCases := []struct {
		name     string
		query    *db.Query
		expected []int
	}{{
		name: "Should apply filter",
		query: db.NewQuery().
			Where(sq.GtOrEq{column: 2}).
			OrderBy(column, db.OrderAsc),
		expected: []int{2, 3},
	}, {
		name: "Should apply limit",
		query: db.NewQuery().
			Limit(2).
			OrderBy(column, db.OrderAsc),
		expected: []int{0, 1},
	}, {
		name: "Should apply order",
		query: db.NewQuery().
			OrderBy(column, db.OrderDesc),
		expected: []int{3, 2, 1, 0},
	}, {
		name: "Should apply offset",
		query: db.NewQuery().
			Limit(2).
			Offset(2).
			OrderBy(column, db.OrderAsc),
		expected: []int{2, 3},
	}, {
		name: "Should apply filtering and pagination",
		query: db.NewQuery().
			Where(sq.Lt{column: 3}).
			Limit(1).
			Offset(1).
			OrderBy(column, db.OrderDesc),
		expected: []int{1},
	}}

	for _, tt := range testCases {
		s.Run(tt.name, s.inTransaction(func() {
			res, err := tt.query.
				WithBuilder(s.resourceManager.GetBuilder()).
				Select(db.NewIncludeMask("id", column)).
				From(&otherModel{}).
				List()
			s.Require().NoError(err)

			actual := make([]db.Model, len(res))
			for i, model := range res {
				s.Require().Len(model, 1)
				actual[i] = model[0]
			}
			expected := make([]db.Model, len(tt.expected))
			for i, index := range tt.expected {
				expected[i] = resources[index]
			}
			s.Assert().Equal(expected, actual)
		}))
	}
	s.Require().NoError(err)
}

func (s *QueryTestSuite) inTransaction(f func()) func() {
	return func() {
		s.Require().NoError(s.resourceManager.InTransaction(f))
	}
}

func (s *QueryTestSuite) whenModelIsInserted(mask db.FieldMask, model db.Model) int64 {
	id, err := db.NewQuery().
		WithBuilder(s.resourceManager.GetBuilder()).
		From(model).
		Select(mask).
		Insert()
	s.Require().NoError(err)
	return id
}

func (s *QueryTestSuite) whenSingleModelIsFetched(id int64, model db.Model) db.Model {
	res, err := db.NewQuery().
		WithBuilder(s.resourceManager.GetBuilder()).
		From(model).
		Select(db.NewExcludeMask()).
		Where(sq.Eq{"id": id}).
		Fetch()
	s.Require().NoError(err)
	s.Require().Len(res, 1)
	return res[0]
}

const (
	someTable    = "some"
	otherTable   = "other"
	anotherTable = "another"

	id           = 100
	defaultValue = 12345
)

func getSomeModel() *someModel {
	return &someModel{
		id:    db.MakeInt(id),
		value: db.MakeFloat(123),
		name:  db.MakeString("abc"),
		flag:  db.MakeBool(true),
		date:  db.MakeTime(time.Unix(1e6, 0).UTC()),
	}
}

func getSomeDifferentModel() *someModel {
	return &someModel{
		id:    db.MakeInt(2 * id),
		value: db.MakeFloat(789),
		name:  db.MakeString("xyz"),
		flag:  db.MakeBool(false),
		date:  db.MakeTime(time.Unix(9e6, 0).UTC()),
	}
}

type someModel struct {
	id    sql.NullInt64
	value sql.NullFloat64
	name  sql.NullString
	flag  sql.NullBool
	date  sql.NullTime
}

func (s *someModel) withId(id int64) *someModel {
	s.id = db.MakeInt(id)
	return s
}

func (s *someModel) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "some",
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"value": {
				SqlType: sqorc.ColumnTypeReal,
			},
			"name": {
				SqlType: sqorc.ColumnTypeText,
			},
			"flag": {
				SqlType: sqorc.ColumnTypeBool,
			},
			"date": {
				SqlType: sqorc.ColumnTypeDatetime,
			},
		},
		Relations: nil,
		CreateObject: func() db.Model {
			return &someModel{}
		},
	}
}

func (s *someModel) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":    db.IntType{X: &s.id},
		"value": db.FloatType{X: &s.value},
		"name":  db.StringType{X: &s.name},
		"flag":  db.BoolType{X: &s.flag},
		"date":  db.TimeType{X: &s.date},
	}
}

func getOtherModel() *otherModel {
	return &otherModel{
		id:     db.MakeInt(id),
		someId: db.MakeInt(id),
		value:  db.MakeFloat(456),
		name:   db.MakeString("pqr"),
		flag:   db.MakeBool(false),
		date:   db.MakeTime(time.Unix(2e6, 0).UTC()),
	}
}

type otherModel struct {
	id     sql.NullInt64
	someId sql.NullInt64
	value  sql.NullFloat64
	name   sql.NullString
	flag   sql.NullBool
	date   sql.NullTime
}

func (o *otherModel) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "other",
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"some_id": {
				SqlType:  sqorc.ColumnTypeInt,
				Nullable: true,
			},
			"value": {
				SqlType:  sqorc.ColumnTypeReal,
				Nullable: true,
			},
			"name": {
				SqlType:  sqorc.ColumnTypeText,
				Nullable: true,
			},
			"flag": {
				SqlType:  sqorc.ColumnTypeBool,
				Nullable: true,
			},
			"date": {
				SqlType:  sqorc.ColumnTypeDatetime,
				Nullable: true,
			},
		},
		Relations: makeRelationsMap(someTable),
		CreateObject: func() db.Model {
			return &otherModel{}
		},
	}
}

func (o *otherModel) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":      db.IntType{X: &o.id},
		"some_id": db.IntType{X: &o.someId},
		"value":   db.FloatType{X: &o.value},
		"name":    db.StringType{X: &o.name},
		"flag":    db.BoolType{X: &o.flag},
		"date":    db.TimeType{X: &o.date},
	}
}

func getAnotherModel() *anotherModel {
	return &anotherModel{
		id:           db.MakeInt(id),
		otherId:      db.MakeInt(id),
		defaultValue: db.MakeInt(0),
	}
}

type anotherModel struct {
	id           sql.NullInt64
	otherId      sql.NullInt64
	defaultValue sql.NullInt64
}

func (a *anotherModel) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "another",
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"other_id": {
				SqlType:  sqorc.ColumnTypeInt,
				Nullable: true,
			},
			"default_value": {
				SqlType:      sqorc.ColumnTypeInt,
				HasDefault:   true,
				DefaultValue: defaultValue,
			},
		},
		Relations: makeRelationsMap(otherTable),
		CreateObject: func() db.Model {
			return &anotherModel{}
		},
	}
}

func (a *anotherModel) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":            db.IntType{X: &a.id},
		"other_id":      db.IntType{X: &a.otherId},
		"default_value": db.IntType{X: &a.defaultValue},
	}
}

func getModelWithUniqueFields() *modelWithUniqueFields {
	return &modelWithUniqueFields{
		id:                db.MakeInt(id),
		uniqueField:       db.MakeInt(id + 1),
		anotherUniqueFied: db.MakeInt(id + 2),
	}
}

type modelWithUniqueFields struct {
	id                sql.NullInt64
	uniqueField       sql.NullInt64
	anotherUniqueFied sql.NullInt64
}

func (m *modelWithUniqueFields) GetMetadata() *db.ModelMetadata {
	return &db.ModelMetadata{
		Table: "unique_table",
		Properties: map[string]*db.Field{
			"id": {
				SqlType: sqorc.ColumnTypeInt,
			},
			"unique_field": {
				SqlType: sqorc.ColumnTypeInt,
				Unique:  true,
			},
			"another_unique_fied": {
				SqlType: sqorc.ColumnTypeInt,
				Unique:  true,
			},
		},
		CreateObject: func() db.Model {
			return &modelWithUniqueFields{}
		},
	}
}

func (m *modelWithUniqueFields) Fields() map[string]db.BaseType {
	return map[string]db.BaseType{
		"id":                  db.IntType{X: &m.id},
		"unique_field":        db.IntType{X: &m.uniqueField},
		"another_unique_fied": db.IntType{X: &m.anotherUniqueFied},
	}
}

func makeRelationsMap(table string) map[string]string {
	return map[string]string{table: table + "_id"}
}
