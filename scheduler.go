/**
 * File              : scheduler.go
 * Author            : Alexandre Saison <alexandre.saison@inarix.com>
 * Date              : 08.11.2021
 * Last Modified Date: 14.11.2021
 * Last Modified By  : Alexandre Saison <alexandre.saison@inarix.com>
 */
package simba

import (
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
)

func funcHandler(dbClient *gorm.DB, client *slack.Client, config *Config) error {
	threadTs, err := SendSlackBlocks(client, config, nil, dbClient)
	if err != nil {
		log.Printf("Error => %s", err)
		return err
	}
	//Sending threadTS
	config.SLACK_MESSAGE_CHANNEL <- threadTs
	return nil
}

func InitScheduler(dbClient *gorm.DB, client *slack.Client, config *Config) (*gocron.Scheduler, *gocron.Job, error) {
	scheduler := gocron.NewScheduler(time.Local)
	if os.Getenv("APP_ENV") == "production" {
		scheduler.CronWithSeconds(config.CRON_EXPRESSION)
	} else if os.Getenv("APP_ENV") == "test" {
		scheduler.Every(10).Minute()
	} else {
		scheduler.Every(10).Minute()
	}

	job, err := scheduler.Do(funcHandler, dbClient, client, config)
	if err != nil {
		return scheduler, nil, err
	} else if job.Error() != nil {
		return scheduler, job, err
	}

	return scheduler, job, nil
}

func watcherNewBadMoodUser(client *slack.Client, config *Config) error {
	for {
		select {
		case newBadMoodUser := <-config.LAST_BAD_MOOD_USER:
			log.Println("newBadMoodUser=", newBadMoodUser)
			//TODO: handle new Bad Mood user to update message
		}
	}
}