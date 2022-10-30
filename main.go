package main

import (
	"errors"
	"log"
	"service-update-accounts/config"
	"strconv"
	"strings"
	"time"

	pq "github.com/pasiol/gopq"
)

var (
	// Version for build
	Version string
	// Build for build
	Build               string
	jobName             = "service-update-accounts"
	Env                 string
	debugState          bool
	updateAccountConfig = ""
)

func getAccounts() ([]string, error) {
	pq.Debug = debugState
	c := config.GetPrimusConfig()
	query := config.Accounts()
	query.Host = c.PrimusHost
	query.Port = c.PrimusPort
	query.User = c.PrimusUser
	query.Pass = c.PrimusPassword

	output, err := pq.ExecuteAndRead(query, 30)
	if err != nil {
		return nil, err
	}
	rows := strings.Fields(output)
	return rows, nil
}

func updateAccount(row string) (string, error) {
	c := config.GetPrimusConfig()
	fields := strings.Split(row, ";")
	account := config.MapData(fields)
	filename, err := config.UpdateAccountXML(account)
	if err != nil {
		return fields[0], err
	}
	_, errorCount, err := pq.ExecuteAtomicImportQuery(filename, c.PrimusHost, c.PrimusPort, c.PrimusUser, c.PrimusPassword, updateAccountConfig)
	if err != nil {
		log.Printf("error: %s", err)
		return fields[0], errors.New("update was not executed")
	}
	if errorCount > 0 {
		return fields[0], errors.New("updating card causes errors")
	}
	return fields[0], nil
}

func main() {

	start := time.Now()

	log.Print("Service started.")

	rows, err := getAccounts()
	if err != nil {
		log.Fatalf("getting accounts failed")
	}
	count := len(rows)
	log.Printf("Founded %s accounts.", strconv.Itoa(count))
	if count > 0 {
		for _, row := range rows {
			id, err := updateAccount(row)
			if err != nil {
				log.Printf("Updating user account failed: %s", id)
			}
			log.Printf("Updating user account succeed: %s", id)
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Print("Ending service in a controlled manner.")
	log.Printf("Elapsed processing time %d.", elapsed)
}
