package DBHandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "root"
	hostname = "mysql"
	dbname   = "url_shortener"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func CreateDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	// err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

func CreateURLShortenerTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS url_shortener(id int primary key, long_url text, 
        short_url text, created_at int default 0)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating url_shortener table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Rows affected when creating table: %d", rows)
	return nil
}

func InsertRow(db *sql.DB, id int64, longUrl string, shortUrl string) error {
	currentTime := time.Now().Unix()
	query := `INSERT INTO url_shortener(id, long_url, short_url, created_at) VALUES (?, ?, ?, ?)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Print("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, id, longUrl, shortUrl, currentTime)
	if err != nil {
		log.Printf("Error %s when inserting row into url_shortener table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d url mapping created ", rows)
	return nil
}

func SearchLongUrl(db *sql.DB, longUrl string) (int64, string, int64, error) {
	query := `SELECT id, short_url,created_at  FROM url_shortener where long_url = ?`
	rows, err := db.Query(query, longUrl)
    if err != nil {
        return 0, "", 0, err
    }
    defer rows.Close()
	var id int64
	var shortUrl string
	var timeStamp int64
	for rows.Next() {
        if err := rows.Scan(&id, &shortUrl, &timeStamp); err != nil {
			log.Printf("Error %s when finding rows ", err)
            return 0, "", 0, err
        }
    }
    if err = rows.Err(); err != nil {
        return 0, "", 0, err
    }
    return id, shortUrl, timeStamp, nil
}

func DeleteRow(db *sql.DB, id int64) error {
	query := `DELETE FROM url_shortener WHERE id = ?`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Print("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Printf("Error %s when deleting row from url_shortener table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d url mapping created ", rows)
	return nil
}

// func DBHandle( {
// 	db, err := dbConnectio()
// 	if err = nil{
// 		og.Printf("Error %s when getting db connection", er)
// 		retrn
//	}
// 	defer db.Clos()
// 	log.Printf("Sucessfully connected to database)
// 	err = createURLShortenerTable(b)
// 	if err = nil{
// 		og.Printf("Create product table failed with error %s", er)
// 	retun
//	}
// }
