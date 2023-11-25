package main

import (
	"backend/src/common"
	"backend/src/db"
	"backend/src/mail"
	"backend/src/server"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func init() {
	err := os.Mkdir("./logs", 0750)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		common.ErrorLogger.Println(err.Error())
		return
	}

	file, err := os.OpenFile(fmt.Sprintf("./logs/logs-%s.txt", time.Now()), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	multi := io.MultiWriter(file, os.Stdout)

	common.DebugLogger = log.New(multi, "[DEBUG] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	common.InfoLogger = log.New(multi, "[INFO] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	common.WarningLogger = log.New(multi, "[WARNING] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	common.ErrorLogger = log.New(multi, "[ERROR] ", log.Ldate|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)

	mail.SmtpHost = os.Getenv("SMTP_HOST")
	mail.SmtpPort = os.Getenv("SMTP_PORT")
	mail.SmtpUser = os.Getenv("SMTP_USER")
	mail.SmtpPass = os.Getenv("SMTP_PASS")

	mail.ImapHost = os.Getenv("IMAP_HOST")
	mail.ImapPort = os.Getenv("IMAP_PORT")
	mail.ImapUser = mail.SmtpUser
	mail.ImapPass = mail.SmtpPass
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	common.InfoLogger.Println("Program Started")

	postgres, err := db.ConnectToPostgres()
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	err = db.Migrate(postgres)
	if err != nil {
		common.ErrorLogger.Println(err)
		return
	}

	api := server.Daemon{
		DB: postgres,
	}

	go server.RunServer(&api)

	for true {
		select {
		case <-interrupt:
			return
		}
	}
}
