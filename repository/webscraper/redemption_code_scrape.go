package webscraper

import (
	env "go_stream_api/environment"
	"go_stream_api/notification"
	db "go_stream_api/repository/database"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gocolly/colly"
)

type redemptionCodeScheduler struct {
	sched          *gocron.Scheduler
	cronExpression string
	urlToVisit     string
	codes          []string
}

func (r *redemptionCodeScheduler) runSchedulerAsync() {
	_, err := r.sched.Cron(r.cronExpression).Do(r.scrapeRedemptionCode)
	if err != nil {
		log.Fatal(err)
	}
	r.sched.StartAsync()
}

func (r *redemptionCodeScheduler) stopScheduler() {
	// Will not stop ongoing job (if there's any)
	r.sched.Clear()
}

func (r *redemptionCodeScheduler) scrapeRedemptionCode() {
	codes := []string{}

	c := colly.NewCollector()

	c.OnHTML(".codes > div > .code", func(e *colly.HTMLElement) {
		codes = append(codes, e.Text)
	})

	c.OnError(errorCallback)

	c.Visit(r.urlToVisit)

	r.codes = codes

	r.handleCodesDistribution()
}

func (r *redemptionCodeScheduler) handleCodesDistribution() {
	codesCollection, err := db.Conn.Pg.GetCodes()
	if err != nil {
		log.Fatal(err)
	}

	if len(codesCollection) != 0 {
		newCodes := []string{}

		for _, code := range r.codes {
			_, codeMatchesWithDB := codesCollection[code]
			if !codeMatchesWithDB {
				newCodes = append(newCodes, code)
			}
		}

		shouldSendNotification := len(newCodes) != 0
		if shouldSendNotification {
			err = notification.SendCodesNotificationToUsers(newCodes)
			if err != nil {
				log.Fatal(err)
			}

			err = db.Conn.Pg.DeleteCodes()
			if err != nil {
				log.Fatal(err)
			}

			err = db.Conn.Pg.InsertCodes(r.codes)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		err = db.Conn.Pg.InsertCodes(r.codes)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func initRedemptionCodeScheduler() *redemptionCodeScheduler {
	return &redemptionCodeScheduler{
		sched:          gocron.NewScheduler(time.Local),
		cronExpression: "0 * * * *",
		urlToVisit:     env.RedemptionCodeURL,
		codes:          nil,
	}
}
