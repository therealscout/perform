package main

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/cagnosolutions/adb"
)

type Auth struct {
	Email    string `json:"email,omitempty" auth:"username"`
	Password string `json:"password,omitempty" auth:"password"`
	Active   bool   `json:"active,omitempty" auth:"active"`
	Role     string `json:"role,omitempty"`
}

type Address struct {
	Street string `json:"street,omitempty"`
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
	Zip    string `json:"zip,omitempty"`
	County string `json:"county,omitempty"`
}

func (a Address) AddrHTML() string {
	if a.Street == "" && a.City == "" && a.State == "" && a.Zip == "" && a.County == "" {
		return ""
	}
	address := a.Street + "<br>" + a.City + ", "
	if a.County != "" {
		address += a.County + ", "
	}
	address += a.State + " " + a.Zip
	return address
}

type CreditCard struct {
	Number          string `json:"number,omitempty"`
	ExpirationDate  string `json:"expirationDate,omitempty"`
	ExpirationMonth int    `json:"expirationMonth,omitempty"`
	ExpirationYear  int    `json:"expirationYear,omitempty"`
	SecurityCode    string `json:"securityCode,omitempty"`
}

type Employee struct {
	Auth
	Id        string `json:"id"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Home      string `json:"home,omitempty"`
	Address
}

type Company struct {
	Id                    string `json:"id"`
	Name                  string `json:"name,omitempty"`
	RegisteredDate        int64  `json:"registeredDate,omitempty"`
	RegistrationFee       int    `json:"registrationPaid,omitempty"`
	RegistrationPaid      bool   `json:"registrationPaid,omitempty"`
	RegistrationPaidDate  int64  `json:"registrationPaidDate,omitempty"`
	CustomerExperienceRep string `json:"customerExperienceRep,omitempty"`
	SalesRep              string `json:"salesRep,omitempty"`
	DUNSNumber            int    `json:"dunsNumber,omitempty"`
}

type Driver struct {
	Auth
	Address
	Id                     string `json:"id"`
	EmployeeId             string `json:"employeeId,omitempty"`
	OriginalId             string `json:"originalId,omitempty"`
	FirstName              string `json:"firstName,omitempty"`
	LastName               string `json:"lastName,omitempty"`
	Phone                  string `json:"phone,omitempty"`
	EmergencyContactName   string `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone  string `json:"emergencyContactPhone,omitempty"`
	LicenseNum             string `json:"licenseNum,omitempty"`
	DOB                    string `json:"dob,omitempty"`
	LicenseState           string `json:"licenseState,omitempty"`
	LicenseExpire          string `json:"licenseExpire,omitempty"`
	LicenseExpireEmailId   string `json:"licenseExpireEmailId, omitempty"`
	LicenseExpireNotify    bool   `json:"licenseExpireNotify"`
	MedCardExpiry          string `json:"medCardExpiry,omitempty"`
	MedCardExpireEmailId   string `json:"medCardExpireEmailId, omitempty"`
	MedCardExpireNotify    bool   `json:"medCardExpireNotify"`
	MVRExpiry              string `json:"mVRExpiry,omitempty"`
	MVRExpireEmailId       string `json:"mVRExpireEmailId, omitempty"`
	MVRExpireNotify        bool   `json:"mVRExpireNotify"`
	ReviewExpiry           string `json:"reviewExpiry,omitempty"`
	ReviewExpireEmailId    string `json:"reviewExpireEmailId, omitempty"`
	ReviewExpireNotify     bool   `json:"reviewExpireNotify"`
	OneEightyExpiry        string `json:"oneEightyExpiry,omitempty"`
	OneEightyExpireEmailId string `json:"oneEightyExpireEmailId, omitempty"`
	OneEightyExpireNotify  bool   `json:"oneEightyExpireNotify"`
	HireDate               string `json:"hireDate,omitempty"`
	TermDate               string `json:"termDate,omitempty"`
	CompanyId              string `json:"companyId,omitempty"`
	Status                 string `json:"driverStatus,omitempty"`
}

func DeleteDriver(driverId string) {
	var documents []Document
	db.TestQuery("document", &documents, adb.Eq("driverId", `"`+driverId+`"`))
	for _, doc := range documents {
		db.Del("document", doc.Id)
	}

	os.RemoveAll("upload/driver/" + driverId + "/")

	db.Del("driver", driverId)
}

func (d *Driver) GenNotifications() {
	now := time.Now()
	beg := time.Date(now.Year(), (now.Month() + 1), now.Day(), 0, 0, 0, 0, now.Location())
	end := beg.AddDate(0, 1, 0)
	var check time.Time
	var err error
	var nId string
	var n Notification

	// license expire
	check, err = time.Parse("01/02/2006", d.LicenseExpire)
	if err != nil {
		check, err = time.Parse("1/02/2006", d.LicenseExpire)
	}
	if err == nil {
		//check = check.AddDate(beg.Year(), 0, 0).In(beg.Location())
		if check.After(beg) && check.Before(end) {
			if !d.LicenseExpireNotify {
				// create new notification
				nId = strconv.Itoa(int(time.Now().UnixNano()))
				n = Notification{
					Id:      nId,
					ModelId: d.CompanyId,
					Type:    "COMPANY",
					SubType: "DRIVER-FORM",
					Title:   "License Expiring",
					Body:    d.FirstName + " " + d.LastName + "'s license is expiring on " + d.LicenseExpire,
					Manual:  false,
				}
				db.Add("notification", nId, n)
				d.LicenseExpireNotify = true
			}
		} else {
			d.LicenseExpireNotify = false
		}
	}

	// med card expire
	check, err = time.Parse("01/02/2006", d.MedCardExpiry)
	if err != nil {
		check, err = time.Parse("1/02/2006", d.MedCardExpiry)
	}
	if err == nil {
		//check = check.AddDate(beg.Year(), 0, 0).In(beg.Location())
		if check.After(beg) && check.Before(end) {
			if !d.MedCardExpireNotify {
				// create new notification
				nId = strconv.Itoa(int(time.Now().UnixNano()))
				n = Notification{
					Id:      nId,
					ModelId: d.CompanyId,
					Type:    "COMPANY",
					SubType: "DRIVER-FORM",
					Title:   "Medical Card Expiring",
					Body:    d.FirstName + " " + d.LastName + "'s medical card is expiring on " + d.MedCardExpiry,
					Manual:  false,
				}
				db.Add("notification", nId, n)
				d.MedCardExpireNotify = true
			}
		} else {
			d.MedCardExpireNotify = false
		}
	}

	// mvr expire
	check, err = time.Parse("01/02/2006", d.MVRExpiry)
	if err != nil {
		check, err = time.Parse("1/02/2006", d.MVRExpiry)
	}
	if err == nil {
		//check = check.AddDate(beg.Year(), 0, 0).In(beg.Location())
		if check.After(beg) && check.Before(end) {
			if !d.MVRExpireNotify {
				// create new notification
				nId = strconv.Itoa(int(time.Now().UnixNano()))
				n = Notification{
					Id:      nId,
					ModelId: d.CompanyId,
					Type:    "COMPANY",
					SubType: "DRIVER-FORM",
					Title:   "MVR Expiring",
					Body:    d.FirstName + " " + d.LastName + "'s MVR is expiring on " + d.MVRExpiry,
					Manual:  false,
				}
				db.Add("notification", nId, n)
				d.MVRExpireNotify = true
			}
		} else {
			d.MVRExpireNotify = false
		}
	}

	// review expire
	check, err = time.Parse("01/02/2006", d.ReviewExpiry)
	if err != nil {
		check, err = time.Parse("1/02/2006", d.ReviewExpiry)
	}
	if err == nil {
		//check = check.AddDate(beg.Year(), 0, 0).In(beg.Location())
		if check.After(beg) && check.Before(end) {
			if !d.ReviewExpireNotify {
				// create new notification
				nId = strconv.Itoa(int(time.Now().UnixNano()))
				n = Notification{
					Id:      nId,
					ModelId: d.CompanyId,
					Type:    "COMPANY",
					SubType: "DRIVER-FORM",
					Title:   "Review due",
					Body:    d.FirstName + " " + d.LastName + "'s review is due on " + d.ReviewExpiry,
					Manual:  false,
				}
				db.Add("notification", nId, n)
				d.ReviewExpireNotify = true
			}
		} else {
			d.ReviewExpireNotify = false
		}
	}

	// 180 expire
	check, err = time.Parse("01/02/2006", d.OneEightyExpiry)
	if err != nil {
		check, err = time.Parse("1/02/2006", d.OneEightyExpiry)
	}
	if err == nil {
		//check = check.AddDate(beg.Year(), 0, 0).In(beg.Location())
		if check.After(beg) && check.Before(end) {
			if !d.OneEightyExpireNotify {
				// create new notification
				nId = strconv.Itoa(int(time.Now().UnixNano()))
				n = Notification{
					Id:      nId,
					ModelId: d.CompanyId,
					Type:    "COMPANY",
					SubType: "DRIVER-FORM",
					Title:   "180 Expiring",
					Body:    d.FirstName + " " + d.LastName + "'s 180 is expiring on " + d.OneEightyExpiry,
					Manual:  false,
				}
				db.Add("notification", nId, n)
				d.OneEightyExpireNotify = true
			}
		} else {
			d.OneEightyExpireNotify = false
		}
	}

}

func (d Driver) GetAge() int {
	dobT, err := time.Parse("01/02/2006", d.DOB)
	if err != nil {
		dobT, err = time.Parse("1/02/2006", d.DOB)
		if err != nil {
			return 0
		}
	}
	dob := dobT.UnixNano()
	diff := time.Now().UnixNano() - dob
	return int(math.Floor((float64(diff) / float64(1000) / float64(1000) / float64(1000) / float64(60) / float64(60) / float64(24) / float64(365.25))))
}

type Vehicle struct {
	Id               string  `json:"id"`
	CompanyId        string  `json:"companyId,omitempty"`
	VehicleType      string  `json:"vehicleType,omitempty"`
	UnitNumber       string  `json:"unitNumber,omitempty"`
	Make             string  `json:"make,omitempty"`
	VIN              string  `json:"vin,omitempty"`
	Title            string  `json:"title,omitempty"`
	GVW              int     `json:"gvw,omitempty"`
	GCR              int     `json:"gcr,omitempty"`
	UnladenWeight    int     `json:"unladenWeight,omitempty"`
	PurchasePrice    float32 `json:"purchasePrice,omitempty"`
	PurchaseDate     string  `json:"purchaseDate,omitempty"`
	CurrentValue     float32 `json:"currentValue,omitempty"`
	AxleAmount       string  `json:"axleAmount,omitempty"`
	FuelType         string  `json:"fuelType,omitempty"`
	FuelTypeOther    string  `json:"fuelTypeOther,omitempty"`
	Active           bool    `json:"active"`
	Owner            string  `json:"owner,omitempty"`
	Year             string  `json:"year,omitempty"`
	PlateNum         string  `json:"plateNum,omitempty"`
	PlateExpire      string  `json:"plateExpire,omitempty"`
	PlateExpireMonth string  `json:"plateExpireMonth,omitempty"`
	PlateExpireYear  string  `json:"plateExpireYear,omitempty"`
	BodyType         string  `json:"bodyType,omitempty"`
	BodyTypeOther    string  `json:"bodyTypeOther,omitempty"`
}

func DeleteVehicle(vehicleId string) {
	os.RemoveAll("upload/vehicle/" + vehicleId + "/")
	db.Del("vehicle", vehicleId)
}

func (v Vehicle) HigherWeight() int {
	if v.GVW > v.GCR {
		return v.GVW
	}
	return v.GCR
}

func (v Vehicle) GetBodyType() string {
	if v.BodyType == "O" {
		return v.BodyTypeOther
	}
	return v.BodyType
}

func (v Vehicle) GetFuelType() string {
	switch v.FuelType {
	case "D":
		return "Diesel"
	case "G":
		return "Gas"
	case "H":
		return "Hybrid"
	case "N":
		return "Natural Gas"
	case "P":
		return "Propane"
	case "O":
		return v.FuelTypeOther
	}
	return "NONE"
}

type Note struct {
	Id              string `json:"id,omitempty"`
	CompanyId       string `json:"companyId,omitempty"`
	EmployeeId      string `json:"employeeId,omitempty"`
	Communication   string `json:"communication,omitempty"`
	Purpose         string `json:"purpose,omitempty"`
	StartTime       int64  `json:"startTime,omitempty"`
	StartTimePretty string `json:"startTimePretty,omitempty"`
	EndTime         int64  `json:"endTime,omitempty"`
	EndTimePretty   string `json:"endTimePretty,omitempty"`
	Representative  string `json:"representative,omitempty"`
	CallBack        string `json:"callBack,omitempty"`
	EmailEmployee   bool   `json:"emailEmployee,omitempty"`
	Billable        bool   `json:"billable,omitempty"`
	Body            string `json:"body,omitempty"`
}

type NoteSort []Note

func (n NoteSort) Len() int {
	return len(n)
}

func (n NoteSort) Less(i, j int) bool {
	return n[i].StartTime < n[j].StartTime
}

func (n NoteSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type Comment struct {
	Id     string `json:"id"`
	Body   string `json:"body"`
	Url    string `json:"url"`
	Page   string `json:"page"`
	Closed bool   `json:"closed"`
}

type QuickNote struct {
	Name string
	Body string
}

type Document struct {
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	DocumentId string   `json:"documentId,omitempty"`
	Complete   bool     `json:"complete"`
	Data       string   `json:"data,omitempty"`
	CompanyId  string   `json:"companyId,omitempty"`
	DriverId   string   `json:"driverId,omitempty"`
	VehicleIds []string `json:"vehicleIds,omitempty"`
	StateForm  bool     `json:"stateForm,omitempty"`
}

type Notification struct {
	Id      string `json:"id"`
	ModelId string `json:"modelId,omitempty"`
	Type    string `json:"type,omitempty"`
	SubType string `json:"subType,omitempty"`
	Title   string `json:"title,omitempty"`
	Body    string `json:"body,omitempty"`
	Manual  bool   `json:"manual"`
}

var DQFS = [][]string{
	[]string{"100", "Driver's Application"},
	[]string{"180", "Certification of Violations"},
	[]string{"200", "Annual Inquery & Review"},
	[]string{"250", "Road Test Certication"},
	[]string{"300", "Previous Driver Inquires"},
	[]string{"400", "Drug & Alcohol Records Request"},
	[]string{"425", "Drug & Alcohol Pre-Employment Statement"},
	[]string{"475", "Alcohol and/or Drug Test Notification"},
	[]string{"450", "Drug & Alcohol Certified Receipt"},
	[]string{"500", "Certification Compliance"},
	[]string{"600", "Confictions for a Driver Violation"},
	[]string{"700", "New Hire Stmt On Duty Hours"},
	[]string{"750", "Other Ompensated Work"},
	[]string{"775", "Fair Credit Reporting Act"},
}

var CompanyForms = [][]string{
	[]string{"MV-550", "2", ""},
	[]string{"MV-550A", "1", ""},
	//[]string{"MV-551", "2"},
	[]string{"MV-552A", "ALL", ""},
	[]string{"MV-558", "4", "none"},
	[]string{"MV-41", "1", ""},
	[]string{"TMT-39", "", ""},
	[]string{"PUC App", "", ""},
	//[]string{"MCS-150", "", ""},
}

type EmailTemplate struct {
	Id   string
	Name string
	Body string
}

type ViolationCache struct {
	LastUpdate           string `json:"lastUpdate,omitempty"`
	Cache                bool   `json:"cache,omitempty"`
	UnsafeDriving        string `json:"unsafeDriving,omitempty"`
	HOSCompliance        string `json:"hosCompliance,omitempty"`
	VehicleMaintenance   string `json:"vehicleMaintenance,omitempty"`
	ControlledSubstances string `json:"controlledSubstances,omitempty"`
	DriverFitness        string `json:"driverFitness,omitempty"`
}

type Task struct {
	Id            string `json:"id"`
	EmployeeId    string `json:"employeeId,omitempty"`
	CompanyId     string `json:"companyId,omitempty"`
	CreatedTime   int64  `json:"createdTime,omitempty"`   // time.Time.Unix()
	AssignedTime  int64  `json:"assignedTime,omitempty"`  // time.Time.Unix()
	StartedTime   int64  `json:"startedTime,omitempty"`   // time.Time.Unix()
	CompletedTime int64  `json:"completedTime,omitempty"` // time.Time.Unix()
	Complete      bool   `json:"complete"`
	Description   string `json:"description,omitempty"`
	Notes         string `json:"notes,omitempty"`
	EmployeeName  string `json:"employeeName, omitempty"`
	CompanyName   string `json:"companyName, omitempty"`
}

func GetTaskEmployeeView(tasks []Task) {
	for i, task := range tasks {
		var company Company
		db.Get("company", task.CompanyId, &company)
		task.CompanyName = company.Name
		tasks[i] = task
	}
}

func GetTaskCompanyView(tasks []Task) {
	for i, task := range tasks {
		var employee Employee
		db.Get("employee", task.EmployeeId, &employee)
		task.EmployeeName = employee.FirstName + " " + employee.LastName
		tasks[i] = task
	}
}

func GetTaskAdminView(tasks []Task) {
	for i, task := range tasks {
		var company Company
		db.Get("company", task.CompanyId, &company)
		task.CompanyName = company.Name
		var employee Employee
		db.Get("employee", task.EmployeeId, &employee)
		task.EmployeeName = employee.FirstName + " " + employee.LastName
		tasks[i] = task
	}
}
