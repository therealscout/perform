package main

import (
	"log"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/mg"
)

type ScheduledEmail struct {
	Id   string `json:"id"`
	Time int64  `json:"time,omitempty"`
	//Data          []Data                 `json:"data,omitempty"`

	Data          map[string]Data        `json:"data,omitempty"`
	Vals          map[string]interface{} `json:"vals,omitempty"`
	Sent          bool                   `json:"sent"`
	Email         mg.Email               `json:"email,omitempty"`
	Template      string                 `json:"template,omitempty"`
	Reschedule    bool                   `json:"reschedule,omitempty"`
	IntervalMonth int                    `json:"intervalMonth"`
	IntervalYear  int                    `json:"intervalYear"`
	GroupId       string                 `json:"groupId,omitempty"`
}

type GroupedEmail struct {
	Id      string
	Time    int64
	DataId  string
	DataKey string
	Vals    map[string]interface{}
	GroupId string
	Sent    bool
}

func (scheduledEmail *ScheduledEmail) Send() (string, error) {
	var groupedEmails []GroupedEmail
	if scheduledEmail.GroupId != "" {
		begM, endM := ThisMonth()
		db.TestQuery("grouped-email", &groupedEmails, adb.Eq("sent", "false"), adb.Gt("time", strconv.Itoa(int(begM))), adb.Lt("time", strconv.Itoa(int(endM))), adb.Eq("groupId", `"`+scheduledEmail.GroupId+`"`))
		if len(groupedEmails) < 1 {
			return "", nil
		}
		for _, ge := range groupedEmails {
			d, ok := scheduledEmail.Data[ge.DataKey]
			if !ok {
				d = Data{Ids: []string{}}
			}
			d.Ids = append(d.Ids, ge.DataId)
			scheduledEmail.Data[ge.DataKey] = d
			scheduledEmail.Vals[ge.DataId] = ge.Vals
		}
	}

	for _, data := range scheduledEmail.Data {
		key, val := data.Get()
		scheduledEmail.Vals[key] = val
	}
	body, err := mg.Body(scheduledEmail.Template, scheduledEmail.Vals, nil)

	if err != nil {
		log.Printf("main.go >> scheduledEmail.Send() >> mg.Body() >> %v\n\n", err)
		return "", err
	}
	scheduledEmail.Email.HTML = body
	resp, err := mg.SendEmail(scheduledEmail.Email)

	if err == nil {
		scheduledEmail.Sent = true
		scheduledEmail.Email.HTML = ""
		for _, ge := range groupedEmails {
			ge.Sent = true
			db.Set("grouped-email", ge.Id, ge)
		}
		if scheduledEmail.Reschedule {
			scheduledEmail.Sent = false
			rt := time.Unix(scheduledEmail.Time, 0)
			rt.AddDate(scheduledEmail.IntervalYear, scheduledEmail.IntervalMonth, 0)
			scheduledEmail.Time = rt.Unix()
		}
		db.Set("scheduled-email", scheduledEmail.Id, scheduledEmail)
	}

	return resp, err
}

type Data struct {
	Store string
	Key   string
	Ids   []string
}

func (data *Data) Get() (string, interface{}) {
	if len(data.Ids) == 1 {
		var val map[string]interface{}
		db.Get(data.Store, data.Ids[0], &val)
		return data.Key, val
	}
	var vals []map[string]interface{}
	for _, id := range data.Ids {
		var val map[string]interface{}
		db.Get(data.Store, id, &val)
		vals = append(vals, val)
	}
	return data.Key, vals
}

func Today() (int64, int64) {
	loc, _ := time.LoadLocation("")
	now := time.Now()
	beg := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := beg.AddDate(0, 0, 1)
	return beg.Unix() - 1, end.Unix()
}

func ThisMonth() (int64, int64) {
	loc, _ := time.LoadLocation("")
	now := time.Now()
	beg := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	end := beg.AddDate(0, 1, 0)
	return beg.Unix() - 1, end.Unix()
}

func Scrape() []ScheduledEmail {
	var scheduledEmail []ScheduledEmail
	beg, end := Today()
	db.TestQuery("scheduled-email", &scheduledEmail, adb.Eq("sent", "false"), adb.Gt("time", strconv.Itoa(int(beg))), adb.Lt("time", strconv.Itoa(int(end))), adb.Eq("groupId", ""))
	return scheduledEmail
}
