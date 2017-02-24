package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cagnosolutions/web"
)

var transferModels = web.Route{"GET", "/transfer/models", func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-------------------------------------\n")
	transferCompanies()
	fmt.Println("-------------------------------------\n")
	transferDrivers()
	fmt.Fprint(w, "Done check terminal for results")
}}

func transferCompanies() {
	var companies2 []Company2
	db.All("company", &companies2)

	fmt.Printf(">> %d original companies <<\n", len(companies2))

	for _, company2 := range companies2 {
		company := Company{
			Auth:              company2.Auth,
			Id:                company2.Id,
			DOTNum:            company2.DOTNum,
			Name:              company2.Name,
			DBA:               company2.DBA,
			ContactName:       company2.ContactName,
			ContactTitle:      company2.ContactTitle,
			ContactSSN:        company2.ContactSSN,
			ContactPhone:      company2.ContactPhone,
			ContactAddress:    company2.ContactAddress,
			SecondName:        company2.SecondName,
			SecondTitle:       company2.SecondTitle,
			SecondPhone:       company2.SecondPhone,
			SameAddress:       company2.SameAddress,
			PhysicalAddress:   company2.PhysicalAddress,
			MailingAddress:    company2.MailingAddress,
			BusinessType:      company2.BusinessType,
			BusinessTypeOther: company2.BusinessTypeOther,
			MCNum:             company2.MCNum,
			PUCNum:            company2.PUCNum,
			Phone:             company2.Phone,
			Fax:               company2.Fax,
			//Email:                   company2.Email,
			EINNum:                  company2.EINNum,
			ARPAccountNum:           company2.ARPAccountNum,
			CarrierType:             company2.GetCarrierType(),
			CarrierTypeOther:        company2.CarrierTypeOther,
			EntityNum:               company2.EntityNum,
			CreditCard:              company2.CreditCard,
			NYHutUsername:           company2.NYHutUsername,
			NYHutPassword:           company2.NYHutPassword,
			NYOscarUsername:         company2.NYOscarUsername,
			NYOscarPassword:         company2.NYOscarPassword,
			KYUseNum:                company2.KYUseNum,
			NMHutUsername:           company2.NMHutUsername,
			NMHutPassword:           company2.NMHutPassword,
			DOTPin:                  company2.DOTPin,
			MCPin:                   company2.MCPin,
			FMCSAUsername:           company2.FMCSAUsername,
			FMCSAPassword:           company2.FMCSAPassword,
			IRPNum:                  company2.IRPNum,
			InsuranceCompany:        company2.InsuranceCompany,
			InsuranceNaic:           company2.InsuranceNaic,
			InsurancePolicyNum:      company2.PolicyNum,
			InsuranceEffectiveDate:  company2.EffectiveDate,
			InsuranceExpirationDate: company2.ExpirationDate,
			OregonNum:               company2.OregonNum,
			GPSProvider:             company2.GPSProvider,
			GPSUsername:             company2.GPSUsername,
			GPSPassword:             company2.GPSPassword,
			FuelCardProvider:        company2.FuelCardProvider,
			FuelCardUsername:        company2.FuelCardUsername,
			FuelCardPassword:        company2.FuelCardPassword,
		}
		db.Set("company", company.Id, company)
		serviceId := strconv.Itoa(int(time.Now().UnixNano()))
		companyService := CompanyService{
			Id:                           serviceId,
			CompanyId:                    company.Id,
			Apportion:                    company2.Service.Apportion,
			ApportionDateOne:             company2.Service.ApportionDateOne,
			ApportionOneComplete:         company2.Service.ApportionOneComplete,
			ApportionDateTwo:             company2.Service.ApportionDateTwo,
			ApportionTwoComplete:         company2.Service.ApportionTwoComplete,
			FuelTaxProgram:               company2.Service.FuelTaxProgram,
			FuelTaxProgramComplete:       company2.Service.FuelTaxProgramComplete,
			FuelTaxNY:                    company2.Service.FuelTaxNY,
			FuelTaxNYComplete:            company2.Service.FuelTaxNYComplete,
			FuelTaxKY:                    company2.Service.FuelTaxKY,
			FuelTaxKYComplete:            company2.Service.FuelTaxKYComplete,
			FuelTaxNM:                    company2.Service.FuelTaxNM,
			FuelTaxNMComplete:            company2.Service.FuelTaxNMComplete,
			DrugConsortium:               company2.Service.DrugConsortium,
			DrugConsortiumDate:           company2.Service.DrugConsortiumDate,
			DrugConsortiumComplete:       company2.Service.DrugConsortiumComplete,
			DriverFileManagement:         company2.Service.DriverFileManagement,
			DriverFileManagementDate:     company2.Service.DriverFileManagementDate,
			DriverFileManagementComplete: company2.Service.DriverFileManagementComplete,
			DOTUpdate:                    company2.Service.DOTUpdate,
			DOTUpdateDate:                company2.Service.DOTUpdateDate,
			DOTUpdateComplete:            company2.Service.DOTUpdateComplete,
			TwentyTwoNinety:              company2.Service.TwentyTwoNinety,
			TwentyTwoNinetyComplete:      company2.Service.TwentyTwoNinetyComplete,
			UCR:                 company2.Service.UCR,
			UCRComplete:         company2.Service.UCRComplete,
			LogAuditing:         company2.Service.LogAuditing,
			LogAuditingComplete: company2.Service.LogAuditingComplete,
			CSAService:          company2.Service.CSAService,
			CSAServiceDate:      company2.Service.CSAServiceDate,
			CSAServiceComplete:  company2.Service.CSAServiceComplete,
			NY:                  company2.Service.NY,
			NYDate:              company2.Service.NYDate,
			NYComplete:          company2.Service.NYComplete,
			GPS:                 company2.Service.GPS,
			GPSDate:             company2.Service.GPSDate,
			GPSComplete:         company2.Service.GPSComplete,
			Training:            company2.Service.Training,
			TrainingDate:        company2.Service.TrainingDate,
			TrainingComplete:    company2.Service.TrainingComplete,
			IFTARenewal:         company2.Service.IFTARenewal,
			IFTARenewalComplete: company2.Service.IFTARenewalComplete,
		}

		db.Set("company-service", serviceId, companyService)
	}

	var companies []Company
	db.All("company", &companies)
	fmt.Printf(">> %d companies after transfer <<\n", len(companies))

	var companyServices []CompanyService
	db.All("company-service", &companyServices)
	fmt.Printf(">> %d company services after transfer <<\n", len(companyServices))
}

func transferDrivers() {
	var drivers2 []Driver2
	db.All("driver", &drivers2)
	fmt.Printf(">> %d original drivers <<\n", len(drivers2))
	for _, driver2 := range drivers2 {
		driver := Driver{
			Auth:                   driver2.Auth,
			Address:                driver2.Address,
			Id:                     driver2.Id,
			EmployeeId:             driver2.EmployeeId,
			OriginalId:             driver2.OriginalId,
			FirstName:              driver2.FirstName,
			LastName:               driver2.LastName,
			Phone:                  driver2.Phone,
			EmergencyContactName:   driver2.EmergencyContactName,
			EmergencyContactPhone:  driver2.EmergencyContactPhone,
			LicenseNum:             driver2.LicenseNum,
			LicenseState:           driver2.LicenseState,
			LicenseExpire:          driver2.LicenseExpire,
			DOB:                    driver2.DOB,
			MedCardExpiry:          driver2.MedCardExpiry,
			MVRExpiry:              driver2.MVRExpiry,
			ReviewExpiry:           driver2.ReviewExpiry,
			OneEightyExpiry:        driver2.OneEightyExpiry,
			HireDate:               driver2.HireDate,
			TermDate:               driver2.TermDate,
			CompanyId:              driver2.CompanyId,
			LicenseExpireEmailId:   driver2.LicenseExpireEmailId,
			MedCardExpireEmailId:   driver2.MedCardExpireEmailId,
			MVRExpireEmailId:       driver2.MVRExpireEmailId,
			ReviewExpireEmailId:    driver2.ReviewExpireEmailId,
			OneEightyExpireEmailId: driver2.OneEightyExpireEmailId,
			Status:                 driver2.GetStatus(),
		}
		db.Set("driver", driver.Id, driver)
	}

	var drivers []Driver
	db.All("driver", &drivers)
	fmt.Printf(">> %d drivers after transfer <<\n", len(drivers))
}
