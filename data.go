package main

/*
import (
	"fmt"
	"strconv"
	"strings"
)

func converVehicles() {
	var vehicle2s []Vehicle2
	db.All("vehicle", &vehicle2s)
	for _, vehicle2 := range vehicle2s {
		vehicle := Vehicle{
			Id:               vehicle2.Id,
			CompanyId:        vehicle2.CompanyId,
			VehicleType:      vehicle2.VehicleType,
			UnitNumber:       vehicle2.UnitNumber,
			Make:             vehicle2.Make,
			VIN:              vehicle2.VIN,
			Title:            vehicle2.Title,
			GVW:              vehicle2.GVW,
			GCR:              vehicle2.GCR,
			UnladenWeight:    vehicle2.UnladenWeight,
			PurchasePrice:    vehicle2.PurchasePrice,
			PurchaseDate:     vehicle2.PurchaseDate,
			CurrentValue:     vehicle2.CurrentValue,
			AxleAmount:       strconv.Itoa(vehicle2.AxleAmount),
			Active:           vehicle2.Active,
			Owner:            vehicle2.Owner,
			Year:             vehicle2.Year,
			PlateNum:         vehicle2.PlateNum,
			PlateExpire:      vehicle2.PlateExpire,
			PlateExpireMonth: vehicle2.PlateExpireMonth,
			PlateExpireYear:  vehicle2.PlateExpireYear,
			BodyType:         ConvertBodyType(vehicle2.BodyType),
			BodyTypeOther:    vehicle2.BodyTypeOther,
			FuelType:         convertFuelType(vehicle2.FuelType),
		}
		db.Set("vehicle", vehicle.Id, vehicle)
	}
	var vehicles []Vehicle
	db.All("vehicle", &vehicles)
	fmt.Printf("old vehicles: %d, converted vehicles %d\n", len(vehicle2s), len(vehicles))
}

func ConvertBodyType(bt BodyType) string {
	switch bt {
	case TT:
		return "TT"
	case TK:
		return "TK"
	case TRL:
		return "TRL"
	case BUS:
		return "BS"
	case SW:
		return "SW"
	case BODY_OTHER:
		return "O"
	}
	return ""
}

func convertFuelType(ft string) string {
	if ft == "Diesel" || ft == "d" || ft == "D" {
		return "D"
	}
	return ""
}

type BodyType int

const (
	TT BodyType = iota
	TK
	TRL
	BUS
	SW
	BODY_OTHER
)

type FuelType int

const (
	DIESEL FuelType = iota
	GAS
	HYBRID
	NATURAL_GAS
	PROPANE
	FUEL_OTHER
)

type Vehicle2 struct {
	Id               string   `json:"id"`
	CompanyId        string   `json:"companyId,omitempty"`
	VehicleType      string   `json:"vehicleType,omitempty"`
	UnitNumber       string   `json:"unitNumber,omitempty"`
	Make             string   `json:"make,omitempty"`
	VIN              string   `json:"vin,omitempty"`
	Title            string   `json:"title,omitempty"`
	GVW              int      `json:"gvw,omitempty"`
	GCR              int      `json:"gcr,omitempty"`
	UnladenWeight    int      `json:"unladenWeight,omitempty"`
	PurchasePrice    float32  `json:"purchasePrice,omitempty"`
	PurchaseDate     string   `json:"purchaseDate,omitempty"`
	CurrentValue     float32  `json:"currentValue,omitempty"`
	AxleAmount       int      `json:"axleAmount,omitempty"`
	FuelType         string   `json:"fuelType,omitempty"`
	Active           bool     `json:"active"`
	Owner            string   `json:"owner,omitempty"`
	Year             string   `json:"year,omitempty"`
	PlateNum         string   `json:"plateNum,omitempty"`
	PlateExpire      string   `json:"plateExpire,omitempty"`
	PlateExpireMonth string   `json:"plateExpireMonth,omitempty"`
	PlateExpireYear  string   `json:"plateExpireYear,omitempty"`
	BodyType         BodyType `json:"bodyType,omitempty"`
	BodyTypeOther    string   `json:"bodyTypeOther,omitempty"`
}

func convertCCExpire() {
	var companies []Company
	db.All("company", &companies)
	for _, company := range companies {
		if company.CreditCard.ExpirationDate == "" {
			continue
		}
		ss := strings.Split(company.CreditCard.ExpirationDate, "/")
		if len(ss) != 3 {
			continue
		}
		if ss[1] == "" {
			continue
		}
		m, err := strconv.Atoi(ss[1])
		if err != nil {
			continue
		}
		if ss[2] == "" {
			continue
		}
		y, err := strconv.Atoi(ss[2])
		if err != nil {
			continue
		}
		if m < 1 || m > 12 || y < 0 {
			continue
		}
		company.CreditCard.ExpirationMonth = m
		company.CreditCard.ExpirationYear = y
		db.Set("company", company.Id, company)
	}
	var companies2 []Company
	db.All("company", &companies2)
	fmt.Printf("Old companies: %d, modified companies: %d\n", len(companies), len(companies2))
}

func convertBusinessType(bt BusinessType) string {
	switch bt {
	case SOLE_PROPRIETOR:
		return "Sole Proprietor"
	case CORPORATION:
		return "Corporation"
	case PARTNERSHIP:
		return "Partnership"
	case LLC:
		return "LLC"
	case LLP:
		return "LLP"
	case BUSINESS_OTHER:
		return "OTHER"
	}
	return ""
}

type BusinessType int

const (
	SOLE_PROPRIETOR BusinessType = iota
	CORPORATION
	PARTNERSHIP
	LLC
	LLP
	BUSINESS_OTHER
)

func convertCompanies() {
	var companies2 []Company2
	db.All("company", &companies2)
	for _, c := range companies2 {
		company := Company{
			Id:                c.Id,
			Auth:              c.Auth,
			DOTNum:            c.DOTNum,
			Name:              c.Name,
			DBA:               c.DBA,
			ContactName:       c.ContactName,
			ContactTitle:      c.ContactTitle,
			ContactSSN:        c.ContactSSN,
			ContactPhone:      c.ContactPhone,
			ContactAddress:    c.ContactAddress,
			SecondName:        c.SecondName,
			SecondTitle:       c.SecondTitle,
			SecondPhone:       c.SecondPhone,
			SameAddress:       c.SameAddress,
			PhysicalAddress:   c.PhysicalAddress,
			MailingAddress:    c.MailingAddress,
			BusinessType:      convertBusinessType(c.BusinessType),
			BusinessTypeOther: c.BusinessTypeOther,
			MCNum:             c.MCNum,
			PUCNum:            c.PUCNum,
			Phone:             c.Phone,
			Fax:               c.Fax,
			Email:             c.Email,
			EINNum:            c.EINNum,
			ARPAccountNum:     c.ARPAccountNum,
			CarrierType:       c.CarrierType,
			CarrierTypeOther:  c.CarrierTypeOther,
			EntityNum:         c.EntityNum,
			CreditCard:        c.CreditCard,
			NYHutUsername:     c.NYHutUsername,
			NYHutPassword:     c.NYHutPassword,
			NYOscarUsername:   c.NYOscarUsername,
			NYOscarPassword:   c.NYOscarPassword,
			KYUseNum:          c.KYUseNum,
			NMHutUsername:     c.NMHutUsername,
			NMHutPassword:     c.NMHutPassword,
			DOTPin:            c.DOTPin,
			MCPin:             c.MCPin,
			FMCSAUsername:     c.FMCSAUsername,
			FMCSAPassword:     c.FMCSAPassword,
			IRPNum:            c.IRPNum,
			InsuranceCompany:  c.InsuranceCompany,
			InsuranceNaic:     c.InsuranceNaic,
			PolicyNum:         c.PolicyNum,
			EffectiveDate:     c.EffectiveDate,
			ExpirationDate:    c.ExpirationDate,
			Service:           c.Service,
			OregonNum:         c.OregonNum,
			GPSProvider:       c.GPSProvider,
			GPSUsername:       c.GPSUsername,
			GPSPassword:       c.GPSPassword,
			FuelCardProvider:  c.FuelCardProvider,
			FuelCardUsername:  c.FuelCardUsername,
			FuelCardPassword:  c.FuelCardPassword,
		}
		db.Set("company", company.Id, company)
	}
	var companies []Company
	db.All("company", &companies)
	fmt.Printf("old vehicles: %d, converted vehicles %d\n", len(companies2), len(companies))
}

type Company2 struct {
	Id string `json:"id"`
	Auth
	DOTNum            string         `json:"dotNum,omitempty"`
	Name              string         `json:"name,omitempty"`
	DBA               string         `json:"dba,omitempty"`
	ContactName       string         `json:"contactName,omitempty"`
	ContactTitle      string         `json:"contactTitle,omitempty"`
	ContactSSN        string         `jsni:"contactSSN,omitempty"`
	ContactPhone      string         `json:"contactPhone,omitempty"`
	ContactAddress    Address        `json:"contactAddress,omitempty"`
	SecondName        string         `json:"secondName,omitempty"`
	SecondTitle       string         `json:"secondTitle,omitempty"`
	SecondPhone       string         `json:"secondPhone,omitempty"`
	SameAddress       bool           `json:"sameAddress"`
	PhysicalAddress   Address        `json:"pysicalAddress,omitempty"`
	MailingAddress    Address        `json:"mailingAddress,omitempty"`
	BusinessType      BusinessType   `json:"businessType,omitempty"`
	BusinessTypeOther string         `json:"businessTypeOther,omitempty"`
	MCNum             string         `json:"mcNum,omitempty"`
	PUCNum            string         `json:"pucNum,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	Fax               string         `json:"fax,omitempty"`
	Email             string         `json:"email,omitempty"`
	EINNum            string         `json:"einNum,omitempty"`
	ARPAccountNum     string         `json:"arpAccountNum,omitempty"`
	CarrierType       CarrierType    `json:"carrierType,omitempty"`
	CarrierTypeOther  string         `json:"carrierTypeOther,omitempty"`
	EntityNum         string         `jaon:"entityNum,omitempty"`
	CreditCard        CreditCard     `json:"crediCard,omitempty"`
	NYHutUsername     string         `json:"nyHutUsername,omitempty"`
	NYHutPassword     string         `json:"nyHutPassword,omitempty"`
	NYOscarUsername   string         `json:"nyOrcarUsername,omitempty"`
	NYOscarPassword   string         `json:"nyOscarUsername,omitempty"`
	KYUseNum          string         `json:"kyUseNum,omitempty"`
	NMHutUsername     string         `json:"nmHutUsername,omitempty"`
	NMHutPassword     string         `json:"nmHutPassword,omitempty"`
	DOTPin            string         `json:"dotPin,omitempty"`
	MCPin             string         `json:"mcPin,omitempty"`
	FMCSAUsername     string         `json:"fmcsaUsername,omitempty"`
	FMCSAPassword     string         `json:"fmcsaPassword,omitempty"`
	IRPNum            string         `json:"irpNum,omitempty"`
	InsuranceCompany  string         `json:"insuranceCompany,omitempty"`
	InsuranceNaic     string         `json:"insuranceNaic,omitempty"`
	PolicyNum         string         `json:"policyNum,omitempty"`
	EffectiveDate     string         `json:"effectiveDate,omitempty"`
	ExpirationDate    string         `json:"expirationDate,omitempty"`
	Service           CompanyService `json:"service,omitempty"`
	OregonNum         string         `json:"oregonNum,omiyempty"`
	GPSProvider       string         `json:"gpsProvider,omiyempty"`
	GPSUsername       string         `json:"gpsUsername,omiyempty"`
	GPSPassword       string         `json:"gpsPassword,omiyempty"`
	FuelCardProvider  string         `json:"fuelCardProvider,omiyempty"`
	FuelCardUsername  string         `json:"fuelCardUsername,omiyempty"`
	FuelCardPassword  string         `json:"fuelCardPassword,omiyempty"`
}
*/
