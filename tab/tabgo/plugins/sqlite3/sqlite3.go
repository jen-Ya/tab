package main

import (
	"database/sql"
	_ "database/sql"

	_ "github.com/mattn/go-sqlite3"
	t "jen-ya.de/tabgo"
)

var Export t.Tab

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func SqliteOpen(arguments t.Tab) t.Tab {
	list := t.ToList(arguments)
	db, err := sql.Open("sqlite3", t.ToString(list[0]))
	panicOnError(err)
	return t.OtherToTab(db)
}

func SqliteClose(arguments t.Tab) t.Tab {
	list := t.ToList(arguments)
	db := t.ToOther(list[0]).(*sql.DB)
	err := db.Close()
	panicOnError(err)
	return t.Tab{}
}

func SqliteQuery(arguments t.Tab) t.Tab {
	list := t.ToList(arguments)
	db := t.ToOther(list[0]).(*sql.DB)
	sqlStmt := t.ToString(list[1])
	rows, err := db.Query(sqlStmt)
	panicOnError(err)
	columns, err := rows.Columns()
	panicOnError(err)
	tabRows := t.TabList{}
	for rows.Next() {
		vals := make([]string, len(columns))
		valPtrs := make([]interface{}, len(columns))
		for i := range columns {
			valPtrs[i] = &vals[i]
		}
		err = rows.Scan(valPtrs...)
		panicOnError(err)
		tabRow := t.TabList{}
		for _, val := range vals {
			tabRow = append(tabRow, t.StringToTab(val))
		}
		tabRows = append(tabRows, t.ListToTab(tabRow))
	}

	return t.ListToTab(tabRows)
}

func init() {
	Export = t.DictToTab(t.TabDict{
		"sqlite-open":  t.NativeFuncToTab(SqliteOpen),
		"sqlite-close": t.NativeFuncToTab(SqliteClose),
		"sqlite-query": t.NativeFuncToTab(SqliteQuery),
	})
}
