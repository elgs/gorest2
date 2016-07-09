// daemon
package gorest2

import (
	"github.com/elgs/cron"
)

type Job struct {
	Action  func()
	Cron    string
	Handler int
}

var Sched *cron.Cron

func StartDaemons() {
	Sched = cron.New()
	//	for _, job := range JobRegistry {
	//		h, err := Sched.AddFunc(job.Cron, job.Action)
	//		if err != nil {
	//			log.Println(err)
	//			continue
	//		}
	//		job.Handler = h
	//	}
	Sched.Start()
}
