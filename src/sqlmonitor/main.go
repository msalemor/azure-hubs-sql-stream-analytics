package main

import (
	"context"
	"database/sql"
	"fmt"
	"hubs/common"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/vutran/ansi"
	"github.com/vutran/ansi/colors"
	"github.com/vutran/ansi/styles"
	"github.com/wonderivan/logger"
)

var db *sql.DB
var wg *sync.WaitGroup

func getRowCount(table string) (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		logger.Error(err)
		return -1, err
	}

	// tsql := fmt.Sprintf("SELECT Id FROM %s;", table)
	// //fmt.Println(tsql)

	// // Execute query
	// rows, err := db.QueryContext(ctx, tsql)
	// if err != nil {
	// 	logger.Error(err)
	// 	return -1, err
	// }

	// defer rows.Close()

	// //fmt.Println(rows)
	// var count int

	// // Iterate through the result set.
	// for rows.Next() {
	// 	// var name, location string
	// 	// var id int

	// 	// // Get values from row.
	// 	// err := rows.Scan(&id)
	// 	// if err != nil {
	// 	// 	return -1, err
	// 	// }

	// 	// fmt.Printf("ID: %d, Name: %s, Location: %s\n", id, name, location)
	// 	count++
	// }

	var count int
	query := fmt.Sprintf("Select count(*) from %s;", table)
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		logger.Error(err)
		return -1, err
	}

	return count, nil
}

func getRowTotals() {

	for {
		acRows, _ := getRowCount("ACEvents")
		genRows, _ := getRowCount("GeneratorEvents")
		motorRows, _ := getRowCount("MotorEvents")
		//logger.Info(acRows, genRows, motorRows)
		fmt.Println(styles.Bold(colors.Green("AC Events:")), acRows)
		fmt.Println(styles.Bold(colors.Green("Generator Events:")), genRows)
		fmt.Println(styles.Bold(colors.Green("Motor Events:")), motorRows)
		time.Sleep(250 * time.Millisecond)
		ansi.HideCursor()
		fmt.Print(ansi.CursorUp(3))
		ansi.ShowCursor()
		//fmt.Print(ansi.EraseLine(0))
	}
	wg.Done()
}

// Server=tcp:sql-alemorhubs-demo-eus.database.windows.net,1433;Initial Catalog=hubdb;Persist Security Info=False;User ID=dbadmin;Password={your_password};MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;
func main() {

	wg = new(sync.WaitGroup)

	connString := common.MustPassEvn("DB_CONNECTION_STRING")

	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		logger.Error("Error creating connection pool: ", err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Debug("Connected!")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Debug("Terminating program & closing database connection")
		db.Close()
		os.Exit(1)
	}()

	logger.Info("Getting row counts:")
	wg.Add(1)
	go getRowTotals()
	wg.Wait()
}
