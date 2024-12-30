package utils

import (
	"C"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)
import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Singleton struct {
	db *sql.DB
}

var (
	instance *Singleton
	once     sync.Once
)

func (s *Singleton) execSql(sql string) {
	if _, err := s.db.Exec(sql); err != nil {
		panic(err)
	}
}

func (s *Singleton) UpdateLastOKTime(ip net.IP, t time.Time) {
	formatTime := t.Format("2006-01-02 15:04:05")
	s.execSql(fmt.Sprintf(`insert or ignore into foo (ip_address, last_ok_time) values ('%s', '%s')`, ip, formatTime))
	s.execSql(fmt.Sprintf(`update foo set last_ok_time = '%s' where ip_address = '%s'`, formatTime, ip))
}

func (s *Singleton) GetRecords() []Record {
	var records []Record
	rows, err := s.db.Query(`select * from foo`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var record Record
		err := rows.Scan(&record.ID, &record.IPAddress, &record.LastOkTime)
		if err != nil {
			panic(err)
		}
		records = append(records, record)
	}
	return records
}

func (s *Singleton) Close() {
	s.db.Close()
}

func GetInstance() *Singleton {
	once.Do(func() {
		var err error
		db, err := sql.Open("sqlite3", "./foo.db")
		if err != nil {
			panic(err)
		}
		instance = &Singleton{db: db}
		instance.execSql(`create table if not exists foo (id integer not null primary key, ip_address text unique, last_ok_time datetime);`)
	})
	return instance
}
