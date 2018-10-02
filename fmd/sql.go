package fmd

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	fmdDB *sql.DB
)

const createTable = `CREATE TABLE IF NOT EXISTS markets (
  uid           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  unixtime      INT     NOT NULL,
  poloVol       FLOAT   NOT NULL,
  poloBtcFlo    FLOAT   NOT NULL,
  bittrexVol    FLOAT   NOT NULL,
  bittrexBtcFlo FLOAT   NOT NULL,
  cmcBtcUsd     FLOAT   NOT NULL,
  cmcFloUsd     FLOAT   NOT NULL,
  cmcLtcUsd     FLOAT   NOT NULL,
  volume        FLOAT   NOT NULL,
  weightedBtc   FLOAT   NOT NULL,
  weightedUsd   FLOAT   NOT NULL,
  mrrLast10     FLOAT   NOT NULL,
  mrrLast24hr   FLOAT   NOT NULL,
  nh24hr        FLOAT   NOT NULL,
  nhLast        FLOAT   NOT NULL
)`

const insertStatement = `INSERT INTO markets (unixtime, poloVol, poloBtcFlo, bittrexVol, bittrexBtcFlo, cmcBtcUsd, cmcFloUsd, cmcLtcUsd, volume, weightedBtc, weightedUsd, mrrLast10, mrrLast24hr, nhLast, nh24hr)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

func InitDB(path string) error {
	var err error
	fmdDB, err = sql.Open("sqlite3", "file:"+path+"?cache=shared&mode=rwc")
	if err != nil {
		return err
	}
	_, err = fmdDB.Exec(createTable)
	if err != nil {
		return err
	}

	return nil
}

func insertToDb(unixtime int64, poloVol, poloBtcFlo, bittrexVol, bittrexBtcFlo, cmcBtcUsd, cmcFloUsd, cmcLtcUsd, volume, weightedBtc, weightedUsd, mrrLast10, mrrLast24hr, nhLast, nh24hr float64) error {
	insertPrepared, err := fmdDB.Prepare(insertStatement)
	if err != nil {
		return err
	}
	_, err = insertPrepared.Exec(unixtime, poloVol, poloBtcFlo, bittrexVol, bittrexBtcFlo, cmcBtcUsd, cmcFloUsd, cmcLtcUsd, volume, weightedBtc, weightedUsd, mrrLast10, mrrLast24hr, nh24hr, nhLast)
	return err
}

func fetchDataPoint(from, to, limit int64) ([]DataPoint, error) {
	stmt, err := fmdDB.Prepare(`SELECT
  unixtime,
  poloVol,
  poloBtcFlo,
  bittrexVol,
  bittrexBtcFlo,
  cmcBtcUsd,
  cmcFloUsd,
  cmcLtcUsd,
  volume,
  weightedBtc,
  weightedUsd,
  mrrLast10,
  mrrLast24hr,
  nhLast,
  nh24hr
FROM markets
WHERE unixtime >= ? AND unixtime <= ?
ORDER BY unixtime
  DESC
LIMIT ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []DataPoint
	for rows.Next() {
		dp := DataPoint{}
		rows.Scan(&dp.Unixtime, &dp.PoloVol, &dp.PoloBtcFlo, &dp.BittrexVol, &dp.BittrexBtcFlo, &dp.CmcBtcUsd,
			&dp.CmcFloUsd, &dp.CmcLtcUsd, &dp.Volume, &dp.WeightedBtc, &dp.WeightedUsd, &dp.MrrLast10, &dp.MrrLast24hr,
			&dp.NhLast, &dp.Nh24hr)
		res = append(res, dp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if res == nil {
		// return an empty slice rather than nil for later json.Marshal
		res = []DataPoint{}
	}
	return res, nil
}
