package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
	"github.com/cagnosolutions/web"
)

/* --- Employee Management --- */

var employeeAll = web.Route{"GET", "/admin/employee", func(w http.ResponseWriter, r *http.Request) {
	var employees []Employee
	// get all "employees" except the default logins
	db.TestQuery("employee", &employees, adb.Gt("id", `"1"`))
	tc.Render(w, r, "admin-employee-all.tmpl", web.Model{
		"employees": employees,
	})
	return
}}

var employeeView = web.Route{"GET", "/admin/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	employeeId := r.FormValue(":id")
	if employeeId != "new" && !db.Get("employee", employeeId, &employee) {
		web.SetErrorRedirect(w, r, "/admin/employee", "Error finding employee")
		return
	}
	tc.Render(w, r, "admin-employee.tmpl", web.Model{
		"employee": employee,
	})
	return
}}

var employeeSave = web.Route{"POST", "/admin/employee", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue("id")
	var employee Employee
	db.Get("employee", empId, &employee)

	// if employee.Email != r.FormValue("email") {
	// 	var employees []Employee
	// 	db.TestQuery("employee", &employees, adb.Eq("email", employee.Email))
	// 	if len(employees) > 0 {
	// 		web.SetErrorRedirect(w, r, "/cns/employee/"+employee.Id, "Error saving employee. Email is already registered")
	// 		return
	// 	}
	// }

	FormToStruct(&employee, r.Form, "")
	if employee.Id == "" && empId == "" {
		employee.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}

	var employees []Employee
	db.TestQuery("employee", &employees, adb.Eq("email", employee.Email), adb.Ne("id", `"`+employee.Id+`"`))
	if len(employees) > 0 {
		web.SetErrorRedirect(w, r, "/admin/employee/"+employee.Id, "Error saving employee. Email is already registered")
		return
	}
	db.Set("employee", employee.Id, employee)
	web.SetSuccessRedirect(w, r, "/admin/employee/"+employee.Id, "Successfully saved employee")
	return
}}

var adminEmployeeTask = web.Route{"GET", "/admin/employee/:id/task", func(w http.ResponseWriter, r *http.Request) {
	employeeId := r.FormValue(":id")
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.SetErrorRedirect(w, r, "/admin/employee", "Error finding employee")
		return
	}
	beg, end := Today()
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Gt("assignedTime", strconv.Itoa(int(beg))), adb.Lt("assignedTime", strconv.Itoa(int(end))))
	GetTaskEmployeeView(tasks)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-employee-task.tmpl", web.Model{
		"employee":  employee,
		"companies": companies,
		"tasks":     tasks,
		"page":      "today",
	})
}}

var adminEmployeeTaskIncomplete = web.Route{"GET", "/admin/employee/:id/task/incomplete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := r.FormValue(":id")
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.SetErrorRedirect(w, r, "/admin/employee", "Error finding employee")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Eq("complete", "false"))
	GetTaskEmployeeView(tasks)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-employee-task.tmpl", web.Model{
		"employee":  employee,
		"companies": companies,
		"tasks":     tasks,
		"page":      "incomplete",
	})
}}

var adminEmployeeTaskComplete = web.Route{"GET", "/admin/employee/:id/task/complete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := r.FormValue(":id")
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.SetErrorRedirect(w, r, "/admin/employee", "Error finding employee")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("employeeId", `"`+employee.Id+`"`), adb.Eq("complete", "true"))
	GetTaskEmployeeView(tasks)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-employee-task.tmpl", web.Model{
		"employee":  employee,
		"companies": companies,
		"tasks":     tasks,
		"page":      "complete",
	})
}}

var employeeDel = web.Route{"POST", "/admin/employee/:id", func(w http.ResponseWriter, r *http.Request) {
	empId := r.FormValue(":id")
	db.Del("employee", empId)
	web.SetSuccessRedirect(w, r, "/admin/employee", "Successfully deleted employee")
	return
}}

var adminTask = web.Route{"GET", "/admin/task", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	var tasks []Task
	db.All("task", &tasks)
	GetTaskAdminView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-task.tmpl", web.Model{
		"companies": companies,
		"employee":  employee,
		"employees": employees,
		"tasks":     tasks,
		"page":      "all",
	})
}}

var adminTasksave = web.Route{"POST", "/admin/task", func(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	var task Task
	if taskId != "" {
		db.Get("task", taskId, &task)
	} else {
		task.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("01/02/2006", r.FormValue("assignedDate"), loc)
	if err != nil {
		t, err = time.ParseInLocation("1/02/2006", r.FormValue("assignedDate"), loc)
		if err != nil {
			web.SetErrorRedirect(w, r, "/admin/task", "Error saving task.")
			return
		}
	}
	task.Description = r.FormValue("description")
	task.EmployeeId = r.FormValue("employeeId")
	task.CompanyId = r.FormValue("companyId")
	task.AssignedTime = t.Unix()
	task.CreatedTime = time.Now().Unix()
	db.Set("task", task.Id, task)
	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/admin/task"
	}
	web.SetSuccessRedirect(w, r, redirect, "Successfully saved task")
	return

}}

var adminTaskToday = web.Route{"GET", "/admin/task/today", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	beg, end := Today()
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Gt("assignedTime", strconv.Itoa(int(beg))), adb.Lt("assignedTime", strconv.Itoa(int(end))))
	GetTaskAdminView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-task.tmpl", web.Model{
		"companies": companies,
		"employee":  employee,
		"employees": employees,
		"tasks":     tasks,
		"page":      "today",
	})
}}

var adminTaskIncomplete = web.Route{"GET", "/admin/task/incomplete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("complete", "false"))
	GetTaskAdminView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-task.tmpl", web.Model{
		"companies": companies,
		"employee":  employee,
		"employees": employees,
		"tasks":     tasks,
		"page":      "incomplete",
	})
}}

var adminTaskComplete = web.Route{"GET", "/admin/task/complete", func(w http.ResponseWriter, r *http.Request) {
	employeeId := web.GetId(r)
	var employee Employee
	if !db.Get("employee", employeeId, &employee) {
		web.Logout(w)
		web.SetErrorRedirect(w, r, "/login", "Error finding your account")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("complete", "true"))
	GetTaskAdminView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	var companies []Company
	db.All("company", &companies)
	tc.Render(w, r, "admin-task.tmpl", web.Model{
		"companies": companies,
		"employee":  employee,
		"employees": employees,
		"tasks":     tasks,
		"page":      "complete",
	})
}}

var adminCompanyTask = web.Route{"GET", "/admin/company/:id/task", func(w http.ResponseWriter, r *http.Request) {
	companyId := r.FormValue(":id")
	var company Company
	if !db.Get("company", companyId, &company) {
		web.SetErrorRedirect(w, r, "/cns/copany", "Error finding company")
		return
	}
	beg, end := Today()
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("companyId", `"`+company.Id+`"`), adb.Gt("assignedTime", strconv.Itoa(int(beg))), adb.Lt("assignedTime", strconv.Itoa(int(end))))
	GetTaskCompanyView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	tc.Render(w, r, "admin-company-task.tmpl", web.Model{
		"company":   company,
		"employees": employees,
		"tasks":     tasks,
		"page":      "today",
	})
}}

var adminCompanyTaskIncomplete = web.Route{"GET", "/admin/company/:id/task/incomplete", func(w http.ResponseWriter, r *http.Request) {
	companyId := r.FormValue(":id")
	var company Company
	if !db.Get("company", companyId, &company) {
		web.SetErrorRedirect(w, r, "/cns/copany", "Error finding company")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("complete", "false"))
	GetTaskCompanyView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	tc.Render(w, r, "admin-company-task.tmpl", web.Model{
		"company":   company,
		"employees": employees,
		"tasks":     tasks,
		"page":      "incomplete",
	})
}}

var adminCompanyTaskComplete = web.Route{"GET", "/admin/company/:id/task/complete", func(w http.ResponseWriter, r *http.Request) {
	companyId := r.FormValue(":id")
	var company Company
	if !db.Get("company", companyId, &company) {
		web.SetErrorRedirect(w, r, "/cns/copany", "Error finding company")
		return
	}
	var tasks []Task
	db.TestQuery("task", &tasks, adb.Eq("companyId", `"`+company.Id+`"`), adb.Eq("complete", "true"))
	GetTaskCompanyView(tasks)
	var employees []Employee
	db.All("employee", &employees)
	tc.Render(w, r, "admin-company-task.tmpl", web.Model{
		"company":   company,
		"employees": employees,
		"tasks":     tasks,
		"page":      "complete",
	})
}}

/* --- Email Template Management --- */

var emailTemplateAll = web.Route{"GET", "/admin/template", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	var emailTemplates []EmailTemplate
	db.All("emailTemplate", &emailTemplates)
	tc.Render(w, r, "email-templates.tmpl", web.Model{
		"emailTemplate":  emailTemplate,
		"emailTemplates": emailTemplates,
	})
	return
}}

var emailTemplateView = web.Route{"GET", "/admin/template/:id", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	if !db.Get("emailTemplate", r.FormValue(":id"), &emailTemplate) {
		web.SetErrorRedirect(w, r, "/admin/template", "Error finding template")
		return
	}
	var emailTemplates []EmailTemplate
	db.All("emailTemplate", &emailTemplates)
	tc.Render(w, r, "email-templates.tmpl", web.Model{
		"emailTemplate":  emailTemplate,
		"emailTemplates": emailTemplates,
	})
	return
}}

var emailTemplateSave = web.Route{"POST", "/admin/template", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplate EmailTemplate
	db.Get("emailTemplate", r.FormValue("id"), &emailTemplate)
	FormToStruct(&emailTemplate, r.Form, "")
	if emailTemplate.Id == "" {
		emailTemplate.Id = strconv.Itoa(int(time.Now().UnixNano()))
	}
	var emailTemplates []EmailTemplate
	db.TestQuery("emailTemplate", &emailTemplates, adb.Eq("name", emailTemplate.Name), adb.Ne("id", `"`+emailTemplate.Id+`"`))
	if len(emailTemplates) > 0 {
		web.SetErrorRedirect(w, r, "/admin/template/"+r.FormValue("id"), "Error saving email template. Name is already in use")
		return
	}
	db.Set("emailTemplate", emailTemplate.Id, emailTemplate)
	web.SetSuccessRedirect(w, r, "/admin/template/"+emailTemplate.Id, "Successfully saved email template")
	return

}}

var emailTest = web.Route{"GET", "/admin/email/test", func(w http.ResponseWriter, r *http.Request) {
	var emailTemplates []EmailTemplate
	db.All("emailTemplate", &emailTemplates)
	tc.Render(w, r, "email-test.tmpl", web.Model{
		"emailTemplates": emailTemplates,
	})
	return
}}
