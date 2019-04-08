package model

import (
	"encoding/json"
	"log"
	"time"
)

//Timers slice hold all current timers
var Timers = make(map[int]*time.Timer)

//GetTimer returns timer at timerID
func GetTimer(timerID int) *time.Timer {
	return Timers[timerID]
}

//StopTimer stops a timer and removes it from map
func StopTimer(timerID int) {
	stop := Timers[timerID].Stop()
	if stop {
		log.Printf("Timer %v stopped\n", timerID)
	}
	delete(Timers, timerID) //deletes timer from map so it may be re-assigned
}

//UpdateTimer should update the time of an existing timer (delete and create new timer)
func UpdateTimer(newData json.RawMessage, newTime time.Time, payload Payload) bool {

	payload.Data = newData

	//update time in payload object
	payload.ScheduledTime = newTime

	StopTimer(payload.ID)

	timeNow := time.Now().UTC().Add(time.Hour) // Get Norwegian time now  TODO time-schedulerservice
	Duration := newTime.Sub(timeNow)           // subtract now's time from target time to get time until trigger

	if Duration < 0 { //scheduled time has to be in the future
		return false
	}

	switch payload.Task {
	case "peer":
		task, err := payload.GetPeerTask()
		if err != nil {
			log.Println("Something went wrong fetching peertask from payload")
			return false
		}
		log.Printf("Timer %v started with duration %v\n", payload.ID, Duration)
		Timers[payload.ID] = time.AfterFunc(Duration, task.Trigger)

	default:
		return false
	}

	if !payload.UpdatePayload() {
		log.Println("Something went wrong updating the payload")
		return false
	}

	return true
}

//InitializeTimers fetches timers from database on startup
func InitializeTimers() {

	time.Sleep(10 * time.Second) //wait 10 seconds to let db setup properly, not best fix but it does the trick

	payloads := GetPayloads()

	for _, payload := range payloads {
		if payload.ScheduledTime.Sub(time.Now().UTC().Add(time.Hour)) < 0 { // Trigger tasks that has dinged when service was down  TODO time-schedulerservice
			task, err := payload.GetPeerTask()
			if err != nil {
				log.Printf("asdla") //todo
				return
			}

			task.Trigger()
			return
		} else if !ScheduleTask(payload) { //schedule task
			log.Printf("Could not initialize timer for AssignmentID: %v\n", payload.AssignmentID)
			return
		}
	}
}
