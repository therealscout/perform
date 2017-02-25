package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/web"
)

const (
	CperE = 1
	COMP  = 10
	DperC = 10
	VperC = 5
)

func testDrivers() {
	var drivers []Driver
	fmt.Println("Getting all drivers...")
	db.All("driver", &drivers)
	fmt.Println("Compiling list of driver ids...")
	var ids []string
	for _, driver := range drivers {
		ids = append(ids, driver.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all drivers individually by id...")
	i := 0
	for _, id := range ids {
		var driver Driver
		if !db.Get("driver", id, &driver) {
			fmt.Printf("Failed to get driver with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d drivers\n\n", i)
}

func testEmployees() {
	var employees []Employee
	fmt.Println("Getting all employees...")
	db.All("employee", &employees)
	fmt.Println("Compiling list of employee ids...")
	var ids []string
	for _, employee := range employees {
		ids = append(ids, employee.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all employees individually by id...")
	i := 0
	for _, id := range ids {
		var employee Employee
		if !db.Get("employee", id, &employee) {
			fmt.Printf("Failed to get employee with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d employees\n\n", i)
}

func testCompanies() {
	var companies []Company
	fmt.Println("Getting all companies...")
	db.All("company", &companies)
	fmt.Println("Compiling list of company ids...")
	var ids []string
	for _, company := range companies {
		ids = append(ids, company.Id)
	}
	fmt.Println("Waiting 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("Getting all companies individually by id...")
	i := 0
	for _, id := range ids {
		var company Company
		if !db.Get("company", id, &company) {
			fmt.Printf("Failed to get company with id %s\n", id)
			i++
		}
	}
	fmt.Printf("\nFailed to get %d companies\n\n", i)
}

func defaultUsers() {

	developer := Employee{
		FirstName: "developer",
		LastName:  "developer",
	}

	developer.Id = "0"
	developer.Email = "developer@perform.com"
	developer.Password = "developer"
	developer.Active = true
	developer.Role = "DEVELOPER"

	admin := Employee{
		Id:        "1",
		FirstName: "admin",
		LastName:  "admin",
		Auth: Auth{
			Email:    "admin@perform.com",
			Password: "admin",
			Active:   true,
			Role:     "ADMIN",
		},
	}

	company := Company{}

	company.Id = "0"
	company.Name = "Test Company"
	company.RegisteredDate = time.Now().UnixNano()

	db.Set("employee", "0", developer)
	db.Set("employee", "1", admin)
	db.Add("company", "0", company)

}

var makeUsers = web.Route{"GET", "/makeUsers", func(w http.ResponseWriter, r *http.Request) {
	MakeEmployees()
	compIds := MakeCompanies()
	MakeDrivers(compIds)
	MakeVehicles(compIds)
	web.SetSuccessRedirect(w, r, "/", "Successfully made users")
	return
}}

func MakeEmployees() {
	for i := 0; i < (COMP / CperE); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))

		employee := Employee{
			FirstName: "John",
			LastName:  fmt.Sprintf("Smith the %dth", (i + 4)),
			Phone:     fmt.Sprintf("717-777-777%d", i),
		}

		employee.Id = id
		employee.Email = fmt.Sprintf("%d@cns.com", i)
		employee.Password = fmt.Sprintf("Password-%d", i)
		employee.Active = (i%2 == 0)
		employee.Role = "EMPLOYEE"

		employee.Street = fmt.Sprintf("12%d Main Street", 1)
		employee.City = fmt.Sprintf("%dville", i)
		employee.State = fmt.Sprintf("%d state", i)
		employee.Zip = fmt.Sprintf("1234%d", i)

		db.Add("employee", id, employee)
	}
}

func MakeCompanies() [COMP]string {
	compIds := [COMP]string{}
	for i := 0; i < COMP; i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIds[i] = id

		company := Company{}
		company.Id = id
		company.Name = fmt.Sprintf("Company %d", i)

		db.Add("company", id, company)
	}
	return compIds
}

func MakeDrivers(compIds [COMP]string) {
	for i := 0; i < (COMP * DperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / DperC
		d := i % 10
		driver := Driver{
			FirstName:             "Daniel",
			LastName:              fmt.Sprintf("Jones the %dth", (i + 4)),
			Phone:                 fmt.Sprintf("717-777-777%d", d),
			EmergencyContactName:  "Samuel Johnson",
			EmergencyContactPhone: fmt.Sprintf("222-222-222%d", d),
			LicenseNum:            fmt.Sprintf("1234567%d", i),
			LicenseState:          fmt.Sprintf("%d state", i),
			LicenseExpire:         fmt.Sprintf("03/1%d/202%d", d, d),
			DOB:                   fmt.Sprintf("01/1%d/198%d", d, d),
			MedCardExpiry:         fmt.Sprintf("02/1%d/202%d", d, d),
			MVRExpiry:             fmt.Sprintf("03/1%d/202%d", d, d),
			ReviewExpiry:          fmt.Sprintf("04/1%d/202%d", d, d),
			OneEightyExpiry:       fmt.Sprintf("05/1%d/202%d", d, d),
			HireDate:              fmt.Sprintf("06/1%d/199%d", d, d),
			TermDate:              fmt.Sprintf("07/1%d/202%d", d, d),
			CompanyId:             compIds[compIdx],
		}

		driver.Id = id
		driver.Email = fmt.Sprintf("%d@%d.com", i, i)
		driver.Password = fmt.Sprintf("Password-%d", i)
		driver.Active = (i%2 == 0)
		driver.Role = "DRIVER"

		driver.Street = fmt.Sprintf("12%d Main Street", 1)
		driver.City = fmt.Sprintf("%dville", i)
		driver.State = fmt.Sprintf("%d state", i)
		driver.Zip = fmt.Sprintf("1234%d", d)

		db.Add("driver", id, driver)
	}
}

func MakeVehicles(compIds [COMP]string) {
	for i := 0; i < (COMP * VperC); i++ {
		id := strconv.Itoa(int(time.Now().UnixNano()))
		compIdx := i / VperC
		vType := "TRUCK"
		if i%3 == 0 {
			vType = "TRACTOR"
		} else if i%2 == 0 {
			vType = "TRAILER"
		}
		vehicle := Vehicle{
			Id:            id,
			CompanyId:     compIds[compIdx],
			VehicleType:   vType,
			UnitNumber:    fmt.Sprintf("%d", i),
			Make:          fmt.Sprintf("make-%d", i),
			VIN:           fmt.Sprintf("%d%d%d", i, i, i),
			Title:         fmt.Sprintf("title-%d", i),
			GVW:           1000 * i,
			GCR:           1155 * i,
			UnladenWeight: 1357 * i,
			PurchasePrice: float32(i),
			PurchaseDate:  fmt.Sprintf("03/1%d/199%d", i, i),
			CurrentValue:  float32(i),
			Active:        i%2 == 0,
			Owner:         fmt.Sprintf("Vinny P number %d", i),
			Year:          fmt.Sprintf("%d", 1980+compIdx),
			PlateNum:      fmt.Sprintf("%d", 658231+compIdx),
			PlateExpire:   fmt.Sprintf("03/1%d/%d", 1992+compIdx, i),
		}
		if i%3 == 0 {
			vehicle.AxleAmount = "2"
		} else if i%2 == 0 {
			vehicle.AxleAmount = "3"
		} else {
			vehicle.AxleAmount = "4"
		}
		db.Add("vehicle", id, vehicle)
	}
}

var httpError = web.Route{"GET", "/http/error", func(w http.ResponseWriter, r *http.Request) {
	tc.Render(w, r, "error.tmpl", nil)
	return
}}
