// schema - describes table(s) in database
package database

import "time"

const ctxTimeTableCreate = 10 * time.Second

const (
	Schema = `
CREATE TABLE IF NOT EXISTS scheduler
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date VARCHAR(8) NOT NULL,
    title VARCHAR(255) UNIQUE NOT NULL,
    comment VARCHAR(2048) NULL,
    repeat VARCHAR(128) NOT NULL CHECK (LENGTH(repeat) <= 128)
);
CREATE INDEX IF NOT EXISTS date_id ON scheduler (date);`
)
