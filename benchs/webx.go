package benchs

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/postgresql"
)

var webxDB db.Database

func init() {
	st := NewSuite("webx")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, WebxInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, WebxInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, WebxUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, WebxRead)
		st.AddBenchmark("MultiRead limit 2000", 2000*ORM_MULTI, 1000, WebxReadSlice)

		cfg, _ := postgresql.ParseURL(ORM_SOURCE)
		webxDB, _ = postgresql.Open(cfg)
		webxDB.SetMaxIdleConns(ORM_MAX_IDLE)
		webxDB.SetMaxOpenConns(ORM_MAX_CONN)
	}
}

func WebxInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Id = 0
		if _, err := webxDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func WebxInsertMulti(b *B) {
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
		inserter := webxDB.(sqlbuilder.SQLBuilder).InsertInto("models").Columns("id", "name", "title", "fax", "web", "age", "counter").Batch(b.N)
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

func WebxUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := webxDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := webxDB.Collection("models").Find(`id = ?`, m.Id).Update(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func WebxRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		if _, err := webxDB.Collection("models").Insert(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := webxDB.Collection("models").Find().One(m); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
}

func WebxReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			if _, err := webxDB.Collection("models").Insert(m); err != nil {
				fmt.Println(err)
				b.FailNow()
			}
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var models []*Model
		if err := webxDB.Collection("models").Find("id > ?", 0).OrderBy("id").Limit(b.L).All(&models); err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}

}
