package server

import (
	"backend/src/common"
	"backend/src/db"
	"backend/src/mail"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetCheckToken is only there to send a success message to the client if the basic auth will work for other endpoints
func (d *Daemon) GetCheckToken(c *gin.Context) {
	c.Status(200)
	return
}

func (d *Daemon) GetKeystones(c *gin.Context) {
	keystones, err := db.RetrieveAllKeystones(d.DB)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch keystones"})
		return
	}
	var result []*common.KeystoneTransfer
	for _, keystone := range keystones {
		result = append(result, keystone.ToTransfer())
	}
	c.JSON(http.StatusOK, result)
	return
}

func (d *Daemon) GetKeystone(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse id"})
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse id"})
		return
	}
	keystone, err := db.RetrieveKeystoneById(d.DB, uint(id))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch keystone"})
		return
	}
	c.JSON(http.StatusOK, keystone.ToTransfer())
	return
}

func (d *Daemon) GetReflections(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse id"})
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not parse id"})
		return
	}
	reflections, err := db.RetrieveReflectionsByKeystoneID(d.DB, uint(id))
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch reflections"})
		return
	}

	var result []*common.ReflectionTransfer
	for _, reflection := range reflections {
		result = append(result, reflection.ToTransfer())
	}
	c.JSON(http.StatusOK, result)
	return
}

func (d *Daemon) PostPublishKeystone(c *gin.Context) {
	user, exists := c.Get(ContextUserKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}
	authenticatedUser, ok := user.(*common.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User found in context is of incorrect type"})
		return
	}

	var keystone common.Keystone
	err := c.BindJSON(&keystone)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind body"})
		return
	}

	if keystone.Content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no content"})
		return
	}

	if keystone.Title == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no title"})
		return
	}

	keystone.Timestamp = time.Now()
	keystone.UserID = authenticatedUser.ID
	keystone.User = *authenticatedUser

	err = db.InsertKeystone(d.DB, &keystone)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not insert new keystone"})
		return
	}

	subject := fmt.Sprintf("ID: %d - %s", keystone.ID, keystone.Title)
	go func() {
		err = mail.SendKeystoneToAllUsers(d.DB, subject, &keystone, "")
		if err != nil {
			common.ErrorLogger.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		common.InfoLogger.Println("Sent mails for keystone:", keystone.ID)
	}()
	c.JSON(http.StatusOK, gin.H{"result": "success"})
	return
}

func (d *Daemon) PostPublishReflection(c *gin.Context) {
	user, exists := c.Get(ContextUserKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}
	authenticatedUser, ok := user.(*common.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User found in context is of incorrect type"})
		return
	}

	var reflection common.Reflection
	err := c.BindJSON(&reflection)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not bind body"})
		return
	}

	if reflection.Content == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no content"})
		return
	}

	reflection.Timestamp = time.Now()
	reflection.UserID = authenticatedUser.ID
	reflection.User = *authenticatedUser

	err = db.InsertReflection(d.DB, &reflection)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not insert new reflection"})
		return
	}
	keystone, err := db.RetrieveKeystoneByIdFull(d.DB, reflection.KeystoneID)
	if err != nil {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch related keystone"})
		return
	}
	mailingDetails, err := db.RetrieveLastMailingDetailsByKeystoneID(d.DB, reflection.KeystoneID)
	if err != nil && err.Error() != "record not found" {
		common.ErrorLogger.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not fetch last mailing details"})
		return
	}

	var replyId string
	if mailingDetails != nil {
		replyId = mailingDetails.MailID
	}

	subject := fmt.Sprintf("REF: %d ID: %d - %s", reflection.ID, keystone.ID, keystone.Title)
	go func() {
		err = mail.SendReflectionToAllUsers(d.DB, subject, &reflection, replyId)
		if err != nil {
			common.ErrorLogger.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		common.InfoLogger.Println("Sent mails for", keystone.Title)
	}()

	c.JSON(http.StatusOK, gin.H{"result": "success"})
	return
}
