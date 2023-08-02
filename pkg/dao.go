package stat

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"sync"
	"time"
)

var (
	DB           *sql.DB
	singletonDao *dao
	onceDao      sync.Once
)

type dao struct {
	insertStmt          *sql.Stmt
	updateStmt          *sql.Stmt
	selectByPointIdStmt *sql.Stmt
}

func getDao() *dao {
	onceDao.Do(func() {
		singletonDao = &dao{}
		var err error

		singletonDao.createTableIfNotExists()

		singletonDao.insertStmt, err = DB.Prepare(sqlInsert)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

		singletonDao.updateStmt, err = DB.Prepare(sqlUpdate)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

		singletonDao.selectByPointIdStmt, err = DB.Prepare(sqlSelectByPointId)
		if err != nil {
			panic(fmt.Sprintf("failed to prepare statement: %v", err))
		}

	})
	return singletonDao
}
func (dao *dao) createTableIfNotExists() {
	_, err := DB.Exec(sqlCreateTable)
	if err != nil {
		panic(fmt.Sprintf("failed to create table: %v", err))
	}
}

func (dao *dao) insert(stat *stat) (int64, error) {

	result, err := dao.insertStmt.Exec(
		stat.pointId,
		time.Now(),
		stat.min,
		stat.max,
		stat.count,
		stat.sum,
		stat.mean,
		stat.m2,
		stat.stdDev(),
		stat.earliestTimestamp,
		stat.latestTimestamp,
		stat.earliestValue,
		stat.latestValue,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (dao *dao) update(stat *stat) (int64, error) {

	result, err := dao.updateStmt.Exec(
		time.Now(),
		stat.min,
		stat.max,
		stat.count,
		stat.sum,
		stat.mean,
		stat.m2,
		stat.stdDev(),
		stat.earliestTimestamp,
		stat.latestTimestamp,
		stat.earliestValue,
		stat.latestValue,
		stat.pointId,
	)

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (dao *dao) selectByPointId(pointID uint32) (*stat, error) {

	row := dao.selectByPointIdStmt.QueryRow(pointID)
	var s stat
	err := row.Scan(
		&s.pointId,
		&s.lastUpdated,
		&s.min,
		&s.max,
		&s.count,
		&s.sum,
		&s.mean,
		&s.m2,
		&s.earliestTimestamp,
		&s.latestTimestamp,
		&s.earliestValue,
		&s.latestValue,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		Logger.WithField("plugin", PluginName).Errorf("failed to read: %s", err)
	}
	return &s, nil
}
