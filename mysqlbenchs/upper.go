package benchs

import (
	"fmt"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

var upperDB db.Session

func init() {
	st := NewSuite("upper")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, UpperInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, UpperInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, UpperUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, UpperRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, UpperReadSlice)

		cfg, _ := mysql.ParseURL(ORM_SOURCE)
		upperDB, _ = mysql.Open(cfg)
		upperDB.SetMaxIdleConns(ORM_MAX_IDLE)
		upperDB.SetMaxOpenConns(ORM_MAX_CONN)
	}
}

func UpperInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := upperDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func UpperInsertMulti(b *B) {
	var ms []*Model
	wrapExecute(b, func() {
		initDB()
		ms = make([]*Model, 0, 100)
		for i := 0; i < 100; i++ {
			ms = append(ms, NewModel())
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inserter := upperDB.SQL().InsertInto("models").Columns("id", "name", "title", "fax", "web", "age", "counter").Batch(b.N)
		go func() {
			for _, m := range ms {
				inserter.Values(m.Id, m.Name, m.Title, m.Fax, m.Web, m.Age, m.Counter)
			}
			inserter.Done()
		}()
		err := inserter.Wait()
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func UpperUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := upperDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := upperDB.Collection("models").Find(`id = ?`, m.Id).Update(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func UpperRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := upperDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := upperDB.Collection("models").Find().One(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func UpperReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := upperDB.Collection("models").Insert(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := upperDB.Collection("models").Find("id > ?", 0).OrderBy("id").Limit(b.L).All(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}
