// daemon
package gorest2

import (
	"log"

	"github.com/elgs/cron"
)

type Job struct {
	Action  func()
	Cron    string
	Handler int
}

var Sched *cron.Cron
var JobRegistry = make(map[string]*Job)

func RegisterJob(id string, job *Job) {
	JobRegistry[id] = job
}

func GetJob(id string) *Job {
	return JobRegistry[id]
}

func StartDaemons() {
	Sched = cron.New()
	for _, job := range JobRegistry {
		h, err := Sched.AddFunc(job.Cron, job.Action)
		if err != nil {
			log.Println(err)
			continue
		}
		job.Handler = h
	}
	Sched.Start()
}
