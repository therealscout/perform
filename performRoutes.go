package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	http.Redirect(w, r, "/login", 303)
	return
}}

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "login.tmpl", web.Model{})
	return
}}

var loginPost = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	email, pass := r.FormValue("email"), r.FormValue("password")
	var employee Employee
	if !db.Auth("employee", email, pass, &employee) {
		web.SetErrorRedirect(w, r, "/login", "Incorrect username or password")
		return
	}
	sess := web.Login(w, r, employee.Role)
	sess.PutId(w, employee.Id)
	sess["email"] = employee.Email
	sess["collapse"] = false
	web.PutMultiSess(w, r, sess)
	redirect := "/cns"
	if employee.Home != "" {
		redirect = employee.Home
	}
	web.SetSuccessRedirect(w, r, redirect, "Welcome "+employee.FirstName)
	return
}}

var index = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/cns", 303)
	return
}}

var cnsHome = web.Route{"GET", "/cns", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	tc.Render(w, r, "cns-home.tmpl", web.Model{
		"employee": employee,
	})
}}

var cnsAction = web.Route{"GET", "/cns/action", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "cns-action.tmpl", nil)
}}

var cnsTask = web.Route{"GET", "/cns/task", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	beg, end := Today()
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Gt("assignedTime", strconv.Itoa(int(beg))), adb.Lt("assignedTime", strconv.Itoa(int(end))))
	GetTaskEmployeeView(tasks)
	tc.Render(w, r, "cns-task.tmpl", web.Model{
		"employee": employee,
		"tasks":    tasks,
		"page":     "today",
	})
}}

var cnsTaskIncomplete = web.Route{"GET", "/cns/task/incomplete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Eq("complete", "false"))
	GetTaskEmployeeView(tasks)
	tc.Render(w, r, "cns-task.tmpl", web.Model{
		"employee": employee,
		"tasks":    tasks,
		"page":     "incomplete",
	})
}}

var cnsTaskComplete = web.Route{"GET", "/cns/task/complete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Eq("complete", "true"))
	GetTaskEmployeeView(tasks)
	tc.Render(w, r, "cns-task.tmpl", web.Model{
		"employee": employee,
		"tasks":    tasks,
		"page":     "complete",
	})
}}

var cnsTaskMarkStart = web.Route{"POST", "/cns/task/:id/start", func(w http.ResponseWriter, r *http.Request) {
	var task Task
	db.Get("task", r.FormValue(":id"), &task)
	task.StartedTime = time.Now().Unix()
	db.Set("task", task.Id, task)
	web.SetSuccessRedirect(w, r, "/cns/task", "Successfully Started Task")
	return
}}

var cnsTaskMarkComplete = web.Route{"POST", "/cns/task/:id/complete", func(w http.ResponseWriter, r *http.Request) {
	var task Task
	db.Get("task", r.FormValue(":id"), &task)
	task.CompletedTime = time.Now().Unix()
	task.Complete = true
	db.Set("task", task.Id, task)
	web.SetSuccessRedirect(w, r, "/cns/task", "Successfully Started Task")
	return
}}

var cnsTaskMarkNote = web.Route{"POST", "/cns/task/:id/note", func(w http.ResponseWriter, r *http.Request) {
	var task Task
	db.Get("task", r.FormValue(":id"), &task)
	// task.Notes = r.FormValue("notes")
	if r.FormValue("notes") == "" {
		web.SetErrorRedirect(w, r, "/cns/task", "Error notes fields was empty")
		return
	}
	task.Notes += "<li>" + r.FormValue("notes") + "</li>"
	db.Set("task", task.Id, task)
	web.SetSuccessRedirect(w, r, "/cns/task", "Successfully saved notes")
	return
}}

var saveHomePage = web.Route{"POST", "/cns/employee/:id/homepage", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	db.Get("employee", r.FormValue(":id"), &employee)
	employee.Home = r.FormValue("url")
	db.Set("employee", employee.Id, employee)
	ajaxResponse(w, `{"error":false}`)
	return
}}

/* --- Company Management --- */

var companyAll = web.Route{"GET", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "company-all.tmpl", web.Model{
		"companies": companies,
	})
	return
}}

var companyView = web.Route{"GET", "/cns/company/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	compId := r.FormValue(":id")
	if compId != "new" && !db.Get("company", compId, &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var notes NoteSort
	var employees []Employee
	db.TestQuery("note", &notes, adb.Eq("companyId", `"`+company.Id+`"`))
	sort.Stable(sort.Reverse(notes))
	db.All("employee", &employees)
	tc.Render(w, r, "company.tmpl", web.Model{
		"company":    company,
		"notes":      notes,
		"employees":  employees,
		"quickNotes": quickNotes,
		"employeeId": web.GetId(r),
	})
	return
}}

var companySave = web.Route{"POST", "/cns/company", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	db.Get("company", r.FormValue("id"), &company)
	FormToStruct(&company, r.Form, "")

	if company.Id == "" {
		company.Id = strconv.Itoa(int(time.Now().UnixNano()))
		company.RegisteredDate = time.Now().UnixNano()
	}

	db.Set("company", company.Id, company)

	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id, "Successfully saved company")
	return
}}
var companyNoteSave = web.Route{"POST", "/cns/company/:id/note", func(w http.ResponseWriter, r *http.Request) {
	var note Note
	r.ParseForm()
	FormToStruct(&note, r.Form, "")
	if note.Id == "" {
		note.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	dt, err := time.Parse("01/02/2006 3:04 PM", r.FormValue("dateTime"))
	if err != nil {
		log.Printf("cnsRoutes.go >> companySaveNotes >> time.Parse() >> %v\n", err)
	}
	note.StartTime = dt.Unix()
	note.EndTime = dt.Unix()
	note.StartTimePretty = r.FormValue("dateTime")
	note.EndTimePretty = r.FormValue("dateTime")
	db.Set("note", note.Id, note)
	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue(":id"), "Successfully saved note")
	return
}}

var companyFormAll = web.Route{"GET", "/cns/company/:id/form", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var docs []Document
	db.TestQuery("document", &docs, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("stateForm", "true"))

	var vehicles []Vehicle
	db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`))

	tc.Render(w, r, "company-form.tmpl", web.Model{
		"company":  company,
		"docs":     docs,
		"forms":    CompanyForms,
		"vehicles": vehicles,
	})
	return
}}

var companyFormAdd = web.Route{"POST", "/cns/company/:id/form", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company/", "Error finding company")
		return
	}

	id := strconv.Itoa(int(time.Now().UnixNano()))
	docId := r.FormValue("name")
	var vehicleIds []string
	if r.FormValue("vehicleIds") != "" {
		vehicleIds = strings.Split(r.FormValue("vehicleIds"), ",")
	}
	doc := Document{
		Id:         id,
		Name:       docId,
		DocumentId: "st-" + strings.ToLower(strings.Replace(docId, " ", "_", -1)),
		Complete:   false,
		CompanyId:  company.Id,
		VehicleIds: vehicleIds,
		StateForm:  true,
	}
	db.Add("document", id, doc)
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/form", "Successfully added forms")
	return
}}

var companyFormDel = web.Route{"POST", "/cns/company/:compId/form/:formId", func(w http.ResponseWriter, r *http.Request) {
	var form Document
	compId := r.FormValue(":compId")
	if !db.Get("document", r.FormValue(":formId"), &form) || form.CompanyId != compId {
		web.SetErrorRedirect(w, r, "/cns/company/"+compId+"/form", "Error deleting from")
		return
	}
	db.Del("document", form.Id)
	web.SetSuccessRedirect(w, r, "/cns/company/"+compId+"/form", "Successfully deleted form")
	return

}}

var companyFormArchive = web.Route{"POST", "/cns/company/:id/archive", func(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("name") == "" {
		fmt.Println("no name")
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form. No name was specified.<br>Please enter a name and try again."}`)
		return
	}
	name := r.FormValue("name") + ".pdf"
	companyId := r.FormValue(":id")
	if companyId == "" {
		log.Printf("main.go -> uploadCompanyFile() -> os.MkdirAll() -> no company id specified")
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form as `+name+`.<br>Please try again"}`)
		return
	}
	path := "upload/company/" + companyId + "/files/private/archived-forms/"
	if _, err := os.Stat(path + name); err == nil {
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form as `+name+`. <br>A file or archived form aleady exists with that name.<br>Please rename and try again"}`)
		return
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go -> companyFormArchive() -> os.MkdirAll() -> %v\n", err)
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form as `+name+`.<br>Please try again"}`)
		return
	}
	r.ParseForm()
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("main.go -> companyformArchive() -> r.FormFile() -> %v\n", err)
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form as `+name+`.<br>Please try again"}`)
		return
	}

	defer file.Close()
	f, err := os.OpenFile(path+name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("main.go -> companyFormArchive() -> os.OpenFile() -> %v\n", err)
		ajaxResponse(w, `{"error":true,"msg":"Error archiving form as `+name+`.<br>Please try again"}`)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	ajaxResponse(w, `{"error":false,"msg":"Successfully archived form as `+name+`"}`)
	return
}}

var companyFileAll = web.Route{"GET", "/cns/company/:id/file", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	if _, err := os.Stat("upload/company/" + company.Id + "/files/public"); err != nil {
		if err := os.MkdirAll("upload/company/"+company.Id+"/files/public", 0755); err != nil {
			fmt.Fprintf(w, "error making dir \"upload/company/%s/files/public\"\n%v", company.Id, err)
		}
	}

	if _, err := os.Stat("upload/company/" + company.Id + "/files/private"); err != nil {
		if err := os.MkdirAll("upload/company/"+company.Id+"/files/private", 0755); err != nil {
			fmt.Fprintf(w, "error making dir \"upload/company/%s/files/private\"\n%v", company.Id, err)
		}
	}

	tc.Render(w, r, "company-file.tmpl", web.Model{
		"company": company,
	})
	return
}}

var companyNotification = web.Route{"GET", "/cns/company/:id/notification", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var notifications []Notification
	db.TestQuery("notification", &notifications, adb.Eq("type", "COMPANY"), adb.Eq("modelId", `"`+company.Id+`"`))
	tc.Render(w, r, "company-notification.tmpl", web.Model{
		"company":       company,
		"notifications": notifications,
	})
}}

var companyNotificationAdd = web.Route{"POST", "/cns/company/:id/notification", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	nId := strconv.Itoa(int(time.Now().UnixNano()))
	notification := Notification{
		Id:      nId,
		ModelId: company.Id,
		Type:    "COMPANY",
		Title:   r.FormValue("title"),
		Body:    r.FormValue("body"),
		Manual:  true,
	}
	db.Add("notification", nId, notification)
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/notification", "Successfully added notification")
	return
}}

var companyNotificaltionDel = web.Route{"POST", "/cns/company/:id/notification/:notificationId", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var notification Notification
	if !db.Get("notification", r.FormValue(":notificationId"), &notification) {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/notification", "Error finding notification")
		return
	}
	if notification.Type != "COMPANY" || notification.ModelId != company.Id {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/notification", "Error deleting notification")
		return
	}
	db.Del("notification", r.FormValue(":notificationId"))
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/notification", "Successfully deleted notification")
	return
}}

// var companyFeature = web.Route{"GET", "/cns/company/:id/feature", func(w http.ResponseWriter, r *http.Request) {
// 	var company Company
// 	if !db.Get("company", r.FormValue(":id"), &company) {
// 		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
// 		return
// 	}
// 	var companyFeatures CompanyFeatures
// 	db.TestQueryOne("company-features", &companyFeatures, adb.Eq("companyId", `"`+company.Id+`"`))
// 	tc.Render(w, r, "company-features.tmpl", web.Model{
// 		"company":         company,
// 		"companyFeatures": companyFeatures,
// 	})
// }}
//
// var companyFeatureSave = web.Route{"POST", "/cns/company/:id/feature", func(w http.ResponseWriter, r *http.Request) {
// 	var company Company
// 	if !db.Get("company", r.FormValue(":id"), &company) {
// 		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
// 		return
// 	}
// 	var companyFeatures CompanyFeatures
// 	db.Get("company-features", r.FormValue("id"), &companyFeatures)
// 	FormToStruct(&companyFeatures, r.Form, "")
// 	if companyFeatures.Id == "" {
// 		companyFeatures.Id = strconv.Itoa(int(time.Now().UnixNano()))
// 		companyFeatures.CompanyId = company.Id
// 	}
// 	db.Set("company-features", companyFeatures.Id, companyFeatures)
// 	active, _ := strconv.ParseBool(r.FormValue("login"))
// 	if active {
// 		if company.Email == "" {
// 			company.Active = false
// 			db.Set("company", company.Id, company)
// 			web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/feature", "Customer must have a valid email address before enabling login.<br>All other features were saved")
// 			return
// 		}
// 		if company.Password == "" {
// 			company.Password = company.Email
// 		}
// 	}
// 	company.Active = active
// 	db.Set("company", company.Id, company)
// 	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/feature", "Successfully saved features")
// 	return
// }}

var companyDel = web.Route{"POST", "/cns/company/:id/del", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}

	// delete all documents
	var documents []Document
	db.TestQuery("document", &documents, adb.Eq("companyId", `"`+company.Id+`"`))
	for _, document := range documents {
		db.Del("document", document.Id)
	}

	// // delete companyFeatures
	// var companyFeatures CompanyFeatures
	// if db.TestQueryOne("company-features", &companyFeatures, adb.Eq("companyId", `"`+company.Id+`"`)) {
	// 	db.Del("company-features", companyFeatures.Id)
	// }

	// delete all notes
	var notes []Note
	db.TestQuery("note", &notes, adb.Eq("companyId", `"`+company.Id+`"`))
	for _, note := range notes {
		db.Del("note", note.Id)
	}

	// delete all notifications
	var notifications []Notification
	db.TestQuery("notification", &notifications, adb.Eq("modelId", `"`+company.Id+`"`), adb.Eq("type", "COMPANY"))
	for _, notification := range notifications {
		db.Del("notification", notification.Id)
	}

	// delete all Drivers
	var drivers []Driver
	db.TestQuery("driver", &drivers, adb.Eq("companyId", `"`+company.Id+`"`))
	for _, driver := range drivers {
		DeleteDriver(driver.Id)
	}

	// delete all vehicles
	var vehicles []Vehicle
	db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`))
	for _, vehicle := range vehicles {
		DeleteVehicle(vehicle.Id)
	}

	// delete company files
	os.RemoveAll("upload/company/" + company.Id + "/")

	// delete company
	db.Del("company", company.Id)

	web.SetSuccessRedirect(w, r, "/cns/company", "Successfully deleted customer")
	return
}}

/* --- Company Vehicle Management --- */

var companyVehicle = web.Route{"GET", "/cns/company/:id/vehicle", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var vehicles []Vehicle
	var state string
	switch r.FormValue("state") {
	case "active":
		db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("active", "true"))
		state = "active"
	case "inactive":
		db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("active", "false"))
		state = "inactive"
	default:
		db.TestQuery("vehicle", &vehicles, adb.Eq("companyId", `"`+company.Id+`"`))
		state = "all"
	}

	tc.Render(w, r, "company-vehicle-all.tmpl", web.Model{
		"company":  company,
		"vehicles": vehicles,
		"state":    state,
	})
	return
}}

var companyVehicleView = web.Route{"GET", "/cns/company/:compId/vehicle/:vId", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":compId"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	vehicleId := r.FormValue(":vId")
	var vehicle Vehicle
	if vehicleId != "new" && !db.Get("vehicle", vehicleId, &vehicle) {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/vehicle", "Error finding vehicle")
		return
	}

	tc.Render(w, r, "company-vehicle.tmpl", web.Model{
		"company": company,
		"vehicle": vehicle,
	})
	return
}}

var companyVehicleSave = web.Route{"POST", "/cns/company/:compId/vehicle", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":compId"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var vehicle Vehicle
	db.Get("vehicle", r.FormValue("id"), &vehicle)
	FormToStruct(&vehicle, r.Form, "")
	if vehicle.Id == "" {
		vehicle.Id = strconv.Itoa(int(time.Now().UnixNano()))
		vehicle.CompanyId = company.Id
	}
	vehicle.PlateExpire = ""
	if vehicle.PlateExpireMonth != "" && vehicle.PlateExpireYear != "" {
		vehicle.PlateExpire = vehicle.PlateExpireMonth + "/" + vehicle.PlateExpireYear
	}
	db.Set("vehicle", vehicle.Id, vehicle)
	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/vehicle/"+vehicle.Id, "Successfully saved vehicle")
	return
}}

var companyVehicleFile = web.Route{"GET", "/cns/company/:compId/vehicle/:vId/file", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":compId"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	var vehicle Vehicle
	if !db.Get("vehicle", r.FormValue(":vId"), &vehicle) || vehicle.CompanyId != company.Id {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/vehicle", "Error finding vehicle")
		return
	}

	var files []map[string]interface{}
	if fileInfos, err := ioutil.ReadDir("upload/vehicle/" + vehicle.Id); err == nil {
		for _, fileInfo := range fileInfos {
			var info = make(map[string]interface{})
			info["name"] = fileInfo.Name()
			info["size"] = fileInfo.Size()
			files = append(files, info)
		}
	}
	tc.Render(w, r, "company-vehicle-file.tmpl", web.Model{
		"company": company,
		"vehicle": vehicle,
		"files":   files,
	})
	return
}}

/* --- Company Driver Management --- */

var companyDriverAll = web.Route{"GET", "/cns/company/:id/driver", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	var drivers []Driver
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	db.TestQuery("driver", &drivers, adb.Eq("companyId", `"`+company.Id+`"`))
	tc.Render(w, r, "company-driver.tmpl", web.Model{
		"company": company,
		"drivers": drivers,
	})
	return
}}

var companyDriverView = web.Route{"GET", "/cns/company/:compId/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":compId"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	driverId := r.FormValue(":id")
	var driver Driver
	if driverId != "new" {
		if !db.Get("driver", driverId, &driver) {
			web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error finding driver")
			return
		}
		if driver.CompanyId != company.Id {
			web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error finding driver")
			return
		}
	}
	var companies []Company
	db.All("company", &companies)

	tc.Render(w, r, "driver.tmpl", web.Model{
		"driver":    driver,
		"company":   company,
		"companies": companies,
	})
	return
}}

var companyDriverFormAll = web.Route{"GET", "/cns/company/:compId/driver/:id/form", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	var docs []Document
	if !db.Get("driver", r.FormValue(":id"), &driver) {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":compId")+"/driver", "Error finding driver")
		return
	}
	if driver.CompanyId != r.FormValue(":compId") {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":compId")+"/driver", "Error finding driver")
		return
	}
	db.TestQuery("document", &docs, adb.Eq("driverId", `"`+driver.Id+`"`), adb.Eq("companyId", `"`+driver.CompanyId+`"`))
	tc.Render(w, r, "driver-form.tmpl", web.Model{
		"driver": driver,
		"dqfs":   DQFS,
		"docs":   docs,
	})
	return
}}

var companyDriverFileAll = web.Route{"GET", "/cns/company/:compId/driver/:id/file", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	if !db.Get("driver", r.FormValue(":id"), &driver) {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":compId")+"/driver", "Error finding driver")
		return
	}
	if driver.CompanyId != r.FormValue(":compId") {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":compId")+"/driver", "Error finding driver")
		return
	}
	var files []map[string]interface{}
	if fileInfos, err := ioutil.ReadDir("upload/driver/" + driver.Id); err == nil {
		for _, fileInfo := range fileInfos {
			var info = make(map[string]interface{})
			info["name"] = fileInfo.Name()
			info["size"] = fileInfo.Size()
			files = append(files, info)
		}
	}
	tc.Render(w, r, "driver-file.tmpl", web.Model{
		"driver": driver,
		"files":  files,
	})
	return
}}

var companyDriverSave = web.Route{"POST", "/cns/driver", func(w http.ResponseWriter, r *http.Request) {
	var driver Driver
	db.Get("driver", r.FormValue("id"), &driver)
	FormToStruct(&driver, r.Form, "")
	if driver.Id == "" {
		driver.Id = strconv.Itoa(int(time.Now().UnixNano()))
		driver.Password = driver.Email
		driver.Role = "DRIVER"
	}
	var company Company
	db.Get("company", driver.CompanyId, &company)
	db.Set("driver", driver.Id, driver)
	web.SetSuccessRedirect(w, r, "/cns/company/"+driver.CompanyId+"/driver/"+driver.Id, "Successfully saved driver")
	return
}}

var companyDriverDel = web.Route{"POST", "/cns/driver/:id", func(w http.ResponseWriter, r *http.Request) {
	driverId := r.FormValue(":id")
	/*var documents []Document
	db.TestQuery("document", &documents, adb.Eq("DriverId", `"`+driverId+`"`))

	for _, doc := range documents {
		db.Del("document", doc.Id)
	}

	os.RemoveAll("upload/driver/" + driverId + "/")

	db.Del("driver", driverId)*/

	DeleteDriver(driverId)

	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue("companyId")+"/driver", "Successfully deleted driver and all of the associated forms and files")
	return

}}

var companyDriverTransfer = web.Route{"POST", "/cns/driver/:id/transfer", func(w http.ResponseWriter, r *http.Request) {

	newDriverId := strconv.Itoa(int(time.Now().UnixNano()))

	driverId := r.FormValue(":id")
	var driver Driver
	db.Get("driver", driverId, &driver)

	companyId := r.FormValue("companyId")

	var documents []Document
	db.TestQuery("document", &documents, adb.Eq("driverId", `"`+driverId+`"`))
	for _, doc := range documents {
		doc.CompanyId = companyId
		db.Set("document", doc.Id, doc)

		newDocId := strconv.Itoa(int(time.Now().UnixNano()))
		doc.Id = newDocId
		doc.DriverId = newDriverId
		doc.CompanyId = driver.CompanyId
		db.Set("document", doc.Id, doc)
	}

	oldCompanyId := driver.CompanyId

	driver.CompanyId = companyId
	db.Set("driver", driverId, driver)

	driver.Id = newDriverId
	driver.Status = "TRANSFERED"
	driver.OriginalId = driverId
	driver.CompanyId = oldCompanyId
	db.Set("driver", newDriverId, driver)

	if err := CopyDir("upload/driver/"+driverId+"/", "upload/driver/"+newDriverId+"/"); err != nil {
		log.Printf("cnsRoutes.go >> transferDriver >> CopyDir() >> %v\n\n", err)
	}

	web.SetSuccessRedirect(w, r, "/cns/company/"+companyId+"/driver/"+driverId, "Successfully transfered driver")
	return
}}

var companyDriverImportUpload = web.Route{"POST", "/cns/company/:id/driver/import", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	path := "upload/company/" + company.Id + "/driver-import/"
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("main.go -> uploadCompanyFile() -> os.MkdirAll() -> %v\n", err)
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error uploading csv file")
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("main.go -> uploadCompanyFile() -> r.FormFile() -> %v\n", err)
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error uploading csv file")
		return
	}
	defer file.Close()
	f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("main.go -> uploadCompanyFile() -> os.OpenFile() -> %v\n", err)
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error uploading csv file")
		return
	}
	defer f.Close()
	io.Copy(f, file)

	web.SetSuccessRedirect(w, r, "/cns/company/"+company.Id+"/driver/import/"+handler.Filename, "Successfully uploaded csv file")
	return
}}

var companyDriverImport = web.Route{"GET", "/cns/company/:id/driver/import/:file", func(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("File: ", r.FormValue(":file"))
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}
	path := "upload/company/" + company.Id + "/driver-import/" + r.FormValue(":file")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error finding file "+r.FormValue(":file"))
		return
	}

	file := NewCSVFile(path)

	tc.Render(w, r, "company-import.tmpl", web.Model{
		"company": company,
		"file":    r.FormValue(":file"),
		"head":    file.GetHeader(),
	})
}}

var companydriverConvert = web.Route{"POST", "/cns/company/:id/driver/convert", func(w http.ResponseWriter, r *http.Request) {
	var company Company
	if !db.Get("company", r.FormValue(":id"), &company) {
		web.SetErrorRedirect(w, r, "/cns/company", "Error finding company")
		return
	}

	path := "upload/company/" + company.Id + "/driver-import/" + r.FormValue("file")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		web.SetErrorRedirect(w, r, "/cns/company/"+company.Id+"/driver", "Error finding file "+r.FormValue(":file"))
		return
	}

	r.ParseForm()
	var drivers []Driver

	csvFile := NewCSVFile(path)

	if err := csvFile.ConvertFromForm(r.Form, &drivers); err != nil {
		web.SetErrorRedirect(w, r, "/cns/company/"+r.FormValue(":id")+"/driver", err.Error())
		return
	}

	for _, driver := range drivers {
		driver.Id = strconv.Itoa(int(time.Now().UnixNano()))
		driver.CompanyId = company.Id
		db.Set("driver", driver.Id, driver)
	}

	web.SetSuccessRedirect(w, r, "/cns/company/"+r.FormValue(":id")+"/driver", "Successfully imported drivers")
	return

}}
