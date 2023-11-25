package mail

import (
	"backend/src/common"
	"backend/src/db"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/smtp"
)

var (
	SmtpHost string
	SmtpPort string
	SmtpUser string
	SmtpPass string
)

func SendKeystoneToAllUsers(dBase *gorm.DB, subject string, keystone *common.Keystone, replyToID string) error {
	users, err := db.RetrieveAllUsers(dBase)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}

	var receivers []string
	for _, user := range users {
		receivers = append(receivers, user.EMail)
	}

	receivers = []string{"setor.blagogee@gmx.de", "capital.snb@gmail.com", "setor@snblago.com"}
	err = SendKeystoneEmail(dBase, receivers, subject, keystone, replyToID)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	return err
}

func SendKeystoneEmail(dBase *gorm.DB, receiverEmails []string, subject string, keystone *common.Keystone, replyToID string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", SmtpUser, SmtpPass, SmtpHost)

	messageId := fmt.Sprintf("%s@snblago.com", uuid.New().String())

	// Email header
	header := make(map[string]string)
	header["From"] = SmtpUser
	header["To"] = SmtpUser
	header["Message-ID"] = messageId
	header["Subject"] = subject
	header["In-Reply-To"] = replyToID
	header["References"] = replyToID

	// Build raw email buffer
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + keystone.Content + "\r\n\r\n" + "Author: " + keystone.User.Nickname

	rawEmail := []byte(message)

	// Send email
	err := smtp.SendMail(SmtpHost+":"+SmtpPort, auth, SmtpUser, receiverEmails, rawEmail)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	} else {
		err = AppendToSentFolder(rawEmail)
		if err != nil {
			common.ErrorLogger.Println(err)
			return err
		}
	}

	err = db.InsertMailDetails(dBase, &common.MailingDetails{
		Subject:    subject,
		MailID:     messageId,
		KeystoneID: keystone.ID,
	})
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	return err
}

func SendReflectionToAllUsers(dBase *gorm.DB, subject string, reflection *common.Reflection, replyToID string) error {
	users, err := db.RetrieveAllUsers(dBase)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}

	var receivers []string
	for _, user := range users {
		receivers = append(receivers, user.EMail)
	}

	receivers = []string{"setor.blagogee@gmx.de", "capital.snb@gmail.com", "setor@snblago.com"}
	err = SendReflectionEmail(dBase, receivers, subject, reflection, replyToID)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	return err
}

func SendReflectionEmail(dBase *gorm.DB, receiverEmails []string, subject string, reflection *common.Reflection, replyToID string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", SmtpUser, SmtpPass, SmtpHost)

	messageId := fmt.Sprintf("%s@snblago.com", uuid.New().String())

	// Email header
	header := make(map[string]string)
	header["From"] = SmtpUser
	header["To"] = SmtpUser
	header["Message-ID"] = messageId
	header["Subject"] = subject
	header["In-Reply-To"] = replyToID
	header["References"] = replyToID

	// Build raw email buffer
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + reflection.Content + "\r\n\r\n" + "Author: " + reflection.User.Nickname

	rawEmail := []byte(message)

	// Send email todo outsource
	err := smtp.SendMail(SmtpHost+":"+SmtpPort, auth, SmtpUser, receiverEmails, rawEmail)
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	} else {
		err = AppendToSentFolder(rawEmail)
		if err != nil {
			common.ErrorLogger.Println(err)
			return err
		}
	}

	err = db.InsertMailDetails(dBase, &common.MailingDetails{
		Subject:      subject,
		MailID:       messageId,
		KeystoneID:   reflection.KeystoneID,
		ReflectionID: &reflection.ID,
	})
	if err != nil {
		common.ErrorLogger.Println(err)
		return err
	}
	return err
}
