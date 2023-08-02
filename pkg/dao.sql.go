package stat

const (
	sqlCreateTable = `CREATE TABLE IF NOT EXISTS pn_stat (
    point_id INTEGER PRIMARY KEY,
    last_updated TIMESTAMP,
    min FLOAT,
    max FLOAT,
    count INT,
    sum FLOAT,
    mean FLOAT,
    m2 FLOAT,
    std_dev FLOAT,
    earliest_time TIMESTAMP,
    latest_time TIMESTAMP,
    earliest_value FLOAT,
    latest_value FLOAT,
    FOREIGN KEY (point_id) REFERENCES point(id),
    UNIQUE(point_id));`
	sqlInsert = `INSERT INTO pn_stat (
                     point_id,
                     last_updated,
                     min,
                     max,
                     count,
                     sum,
                     mean,
                     m2,
                     std_dev,
                     earliest_time,
                     latest_time,
                     earliest_value,
                     latest_value) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlUpdate = `UPDATE pn_stat SET 
                   last_updated = ?,
                   min = ?,
                   max = ?, 
                   count = ?,
                   sum = ?,
                   mean = ?,
                   m2 = ?,
                   std_dev = ?,
                   earliest_time = ?,
                   latest_time = ?,
                   earliest_value = ?,
                   latest_value = ? WHERE point_id = ?`

	sqlSelectByPointId = `SELECT 
    point_id, 
    last_updated, 
    min,
    max,
    count,
    sum,
    mean,
    m2,
    earliest_time,
    latest_time,
    earliest_value,
    latest_value
    FROM 
        pn_stat 
    WHERE 
        point_id = ?`
)
