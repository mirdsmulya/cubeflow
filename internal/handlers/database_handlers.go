package handlers

import (
	"cubeflow/internal/backup"
	"cubeflow/internal/slack"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DatabaseBackupV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		dbName := c.Param("dbname")
		c.JSON(http.StatusOK, "Database backup initiated!")
		go func() {
			_, err := backup.BackupDB(dbName)
			if err != nil {
				log.Println("Failed to backup DB: ", err)
				return
			}
			log.Println("DB bakcup completed!")
		}()

	}
}

func DatabaseBackupV2(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticated, dbName, triggerUserName, incomingSlackChannel := slack.AuthChecker(c)
		ChannelExist := slack.IsChannelExist(incomingSlackChannel)
		if !authenticated || !ChannelExist {
			log.Printf("Error message: failed to authenticated or wrong channel or wrong application name")
			c.JSON(http.StatusOK, "Failed, you might be running this command in wrong channel or wrong application name")
			return
		}

		c.JSON(http.StatusOK, "Database backup initiated")
		slack.SendSlackMessage(incomingSlackChannel, "Database backup started! Please wait for completion :hourglass_flowing_sand: \nDB name: "+dbName+" \nEnv: "+env+" \nTriggered by: "+triggerUserName)

		go func() {
			bucketName, err := backup.BackupDB(dbName)
			if err != nil {
				log.Println("Failed to backup DB: ", err)
				slack.SendSlackMessage(incomingSlackChannel, "Oops DB Backup failed! :x:\n Error: "+err.Error())
				return
			}
			slack.SendSlackMessage(incomingSlackChannel, "Database backup completed! :white_check_mark: \nBackup stored: https://console.cloud.google.com/storage/browser/"+bucketName)
		}()

	}
}
