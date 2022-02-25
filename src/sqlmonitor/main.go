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
	count := -1

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		logger.Error(err)
		return -1, err
	}

	query := fmt.Sprintf("Select count(*) from %s;", table)
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		logger.Error(err)
		return -1, err
	}

	return count, nil
}

func getRowTotals() {

	fmt.Print(ansi.HideCursor())
	for {
		acRows, _ := getRowCount("ACEvents")
		genRows, _ := getRowCount("GeneratorEvents")
		motorRows, _ := getRowCount("MotorEvents")
		fmt.Println(styles.Bold(colors.Green("AC Events:")), acRows)
		fmt.Println(styles.Bold(colors.Green("Generator Events:")), genRows)
		fmt.Println(styles.Bold(colors.Green("Motor Events:")), motorRows)
		time.Sleep(250 * time.Millisecond)
		ansi.HideCursor()
		fmt.Print(ansi.CursorUp(3))
	}
	wg.Done()
}

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
		<-c // Block until message received
		fmt.Print(ansi.ShowCursor())
		logger.Debug("Terminating program & closing database connection")
		db.Close()
		os.Exit(1)
	}()

	logger.Info("Getting row counts:")
	wg.Add(1)
	go getRowTotals()
	wg.Wait()
}
