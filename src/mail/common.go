package mail

import (
	"backend/src/common"
	"backend/src/db"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"gorm.io/gorm"
)

var (
	ImapHost string
	ImapPort string
	ImapUser string
	ImapPass string
)

func FetchSentEmailsMessageID(dBase *gorm.DB) error {
	c, err := ConnectToIMAPServer()
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	//uid := uint32(42)
	fetchOptions := &imap.FetchOptions{Envelope: true}

	// Login, select and fetch a message in a single roundtrip
	selectCmd := c.Select("Sent", nil)

	selectData, err := selectCmd.Wait()
	if err != nil {
		common.ErrorLogger.Printf("failed to select Sent: %v", err)
		return err
	}

	// Calculate the starting message number for the last ten messages
	start := int(selectData.NumMessages) - 3
	if start < 0 {
		start = 0 // In case there are fewer than 10 messages
	}

	// Create a sequence set for the last ten messages
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(uint32(start+1), selectData.NumMessages)

	// Fetch the last ten messages
	fetchCmd := c.UIDFetch(*seqSet, fetchOptions)

	messages, err := fetchCmd.Collect()
	if err != nil {
		common.ErrorLogger.Printf("failed to fetch messages: %v", err)
		return err
	}

	// Process the fetched messages
	for _, msg := range messages {
		common.InfoLogger.Printf("Subject: %v", msg.Envelope.Subject)
		common.InfoLogger.Printf("Message-ID: %v", msg.Envelope.MessageID)
		err = db.InsertMailDetails(dBase, &common.MailingDetails{
			Subject: msg.Envelope.Subject,
			MailID:  msg.Envelope.MessageID,
		})
		if err != nil {
			common.ErrorLogger.Println(err)
			return err
		}
	}

	return nil
}

func ConnectToIMAPServer() (*imapclient.Client, error) {
	// Connect to the IMAP server
	c, err := imapclient.DialTLS(ImapHost+":"+ImapPort, nil)
	if err != nil {
		common.ErrorLogger.Println(err)
		return nil, err
	}

	loginCMD := c.Login(ImapUser, ImapPass)
	err = loginCMD.Wait()
	if err != nil {
		common.ErrorLogger.Println(err)
		return nil, err
	}

	return c, nil
}

func AppendToSentFolder(rawEmail []byte) error {
	c, err := ConnectToIMAPServer()
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}

	size := int64(len(rawEmail))
	appendCmd := c.Append("Sent", size, nil)
	_, err = appendCmd.Write(rawEmail)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	if err = appendCmd.Close(); err != nil {
		common.ErrorLogger.Printf("failed to close message: %v", err)
		return err
	}
	if _, err = appendCmd.Wait(); err != nil {
		common.ErrorLogger.Printf("APPEND command failed: %v", err)
		return err
	}
	// Append the email to the 'Sent' folder
	return nil
}
