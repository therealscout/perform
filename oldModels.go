package main

type CarrierType int

const (
	PRIVATE CarrierType = iota
	COMMON
	CONTRACT
	CARRIER_OTHER
)

type Company2 struct {
	Auth
	Id                string          `json:"id"`
	DOTNum            string          `json:"dotNum,omitempty"`
	Name              string          `json:"name,omitempty"`
	DBA               string          `json:"dba,omitempty"`
	ContactName       string          `json:"contactName,omitempty"`
	ContactTitle      string          `json:"contactTitle,omitempty"`
	ContactSSN        string          `jsni:"contactSSN,omitempty"`
	ContactPhone      string          `json:"contactPhone,omitempty"`
	ContactAddress    Address         `json:"contactAddress,omitempty"`
	SecondName        string          `json:"secondName,omitempty"`
	SecondTitle       string          `json:"secondTitle,omitempty"`
	SecondPhone       string          `json:"secondPhone,omitempty"`
	SameAddress       bool            `json:"sameAddress"`
	PhysicalAddress   Address         `json:"pysicalAddress,omitempty"`
	MailingAddress    Address         `json:"mailingAddress,omitempty"`
	BusinessType      string          `json:"businessType,omitempty"`
	BusinessTypeOther string          `json:"businessTypeOther,omitempty"`
	MCNum             string          `json:"mcNum,omitempty"`
	PUCNum            string          `json:"pucNum,omitempty"`
	Phone             string          `json:"phone,omitempty"`
	Fax               string          `json:"fax,omitempty"`
	Email             string          `json:"email,omitempty"`
	EINNum            string          `json:"einNum,omitempty"`
	ARPAccountNum     string          `json:"arpAccountNum,omitempty"`
	CarrierType       CarrierType     `json:"carrierType,omitempty"`
	CarrierTypeOther  string          `json:"carrierTypeOther,omitempty"`
	EntityNum         string          `jaon:"entityNum,omitempty"`
	CreditCard        CreditCard      `json:"crediCard,omitempty"`
	NYHutUsername     string          `json:"nyHutUsername,omitempty"`
	NYHutPassword     string          `json:"nyHutPassword,omitempty"`
	NYOscarUsername   string          `json:"nyOrcarUsername,omitempty"`
	NYOscarPassword   string          `json:"nyOscarUsername,omitempty"`
	KYUseNum          string          `json:"kyUseNum,omitempty"`
	NMHutUsername     string          `json:"nmHutUsername,omitempty"`
	NMHutPassword     string          `json:"nmHutPassword,omitempty"`
	DOTPin            string          `json:"dotPin,omitempty"`
	MCPin             string          `json:"mcPin,omitempty"`
	FMCSAUsername     string          `json:"fmcsaUsername,omitempty"`
	FMCSAPassword     string          `json:"fmcsaPassword,omitempty"`
	IRPNum            string          `json:"irpNum,omitempty"`
	InsuranceCompany  string          `json:"insuranceCompany,omitempty"`
	InsuranceNaic     string          `json:"insuranceNaic,omitempty"`
	PolicyNum         string          `json:"policyNum,omitempty"`
	EffectiveDate     string          `json:"effectiveDate,omitempty"`
	ExpirationDate    string          `json:"expirationDate,omitempty"`
	Service           CompanyService2 `json:"service,omitempty"`
	OregonNum         string          `json:"oregonNum,omiyempty"`
	GPSProvider       string          `json:"gpsProvider,omiyempty"`
	GPSUsername       string          `json:"gpsUsername,omiyempty"`
	GPSPassword       string          `json:"gpsPassword,omiyempty"`
	FuelCardProvider  string          `json:"fuelCardProvider,omiyempty"`
	FuelCardUsername  string          `json:"fuelCardUsername,omiyempty"`
	FuelCardPassword  string          `json:"fuelCardPassword,omiyempty"`
}

func (c Company2) GetCarrierType() string {
	switch c.CarrierType {
	case PRIVATE:
		return "Private"
	case COMMON:
		return "Common"
	case CONTRACT:
		return "Contract"
	case CARRIER_OTHER:
		return "OTHER"
	}
	return ""
}

var COMPANY_CONSTS map[string]interface{} = map[string]interface{}{
	"PRIVATE":       PRIVATE,
	"COMMON":        COMMON,
	"CONTRACT":      CONTRACT,
	"CARRIER_OTHER": CARRIER_OTHER,
}

type CompanyService2 struct {
	Apportion                    bool   `json:"apportion"`
	ApportionDateOne             string `json:"apportionDateOne,omitempty"`
	ApportionOneComplete         bool   `json:"apportionOneComplete"`
	ApportionDateTwo             string `json:"apportionDateTwo,omitempty"`
	ApportionTwoComplete         bool   `json:"apportionTwoComplete"`
	FuelTaxProgram               bool   `json:"fuelTaxProgram"`
	FuelTaxProgramComplete       bool   `json:"fuelTaxProgramComplete"`
	FuelTaxNY                    bool   `json:"fuelTaxNY"`
	FuelTaxNYComplete            bool   `json:"fuelTaxNYComplete"`
	FuelTaxKY                    bool   `json:"fuelTaxKY"`
	FuelTaxKYComplete            bool   `json:"fuelTaxKYComplete"`
	FuelTaxNM                    bool   `json:"fuelTaxNM"`
	FuelTaxNMComplete            bool   `json:"FuelTaxNMComplete"`
	DrugConsortium               bool   `json:"drugConsortium"`
	DrugConsortiumDate           string `json:"drugConsortiumDate,omitempty"`
	DrugConsortiumComplete       bool   `json:"drugConsortiumComplete"`
	DriverFileManagement         bool   `json:"driverFileManagement"`
	DriverFileManagementDate     string `json:"driverFileManagementDate,omitempty"`
	DriverFileManagementComplete bool   `json:"driverFileManagementComplete"`
	DOTUpdate                    bool   `json:"dotUpdate"`
	DOTUpdateDate                string `json:"dotUpdateDate,omitempty"`
	DOTUpdateComplete            bool   `json:"dotUpdateComplete"`
	TwentyTwoNinety              bool   `json:"twentyTwoNinety"`
	TwentyTwoNinetyComplete      bool   `json:"twentyTwoNinetyComplete"`
	UCR                          bool   `json:"ucr"`
	UCRComplete                  bool   `json:"ucrComplete"`
	LogAuditing                  bool   `json:"logAuditing"`
	LogAuditingComplete          bool   `json:"logAuditingComplete"`
	CSAService                   bool   `json:"csaService"`
	CSAServiceDate               string `json:"csaServiceDate,omitempty"`
	CSAServiceComplete           bool   `json:"csaServiceComplete"`
	NY                           bool   `json:"ny"`
	NYDate                       string `json:"nyDate"`
	NYComplete                   bool   `json:"nyComplete"`
	GPS                          bool   `json:"gps"`
	GPSDate                      string `json:"gpsDate,omitempty"`
	GPSComplete                  bool   `json:"gpsComplete"`
	Training                     bool   `json:"training"`
	TrainingDate                 string `json:"trainingDate,omitempty"`
	TrainingComplete             bool   `json:"trainingComplete"`
	IFTARenewal                  bool   `json:"iftaRenewal"`
	IFTARenewalComplete          bool   `json:"iftaRenewalComplete"`
}

type DriverStatus int

const (
	WORKING DriverStatus = iota
	FIRED
	TRANSFERED
	LAID_OFF
)

type Driver2 struct {
	Auth
	Address
	Id                     string       `json:"id"`
	EmployeeId             string       `json:"employeeId,omitempty"`
	OriginalId             string       `json:"originalId,omitempty"`
	FirstName              string       `json:"firstName,omitempty"`
	LastName               string       `json:"lastName,omitempty"`
	Phone                  string       `json:"phone,omitempty"`
	EmergencyContactName   string       `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone  string       `json:"emergencyContactPhone,omitempty"`
	LicenseNum             string       `json:"licenseNum,omitempty"`
	LicenseState           string       `json:"licenseState,omitempty"`
	LicenseExpire          string       `json:"licenseExpire,omitempty"`
	DOB                    string       `json:"dob,omitempty"`
	MedCardExpiry          string       `json:"medCardExpiry,omitempty"`
	MVRExpiry              string       `json:"mVRExpiry,omitempty"`
	ReviewExpiry           string       `json:"reviewExpiry,omitempty"`
	OneEightyExpiry        string       `json:"oneEightyExpiry,omitempty"`
	HireDate               string       `json:"hireDate,omitempty"`
	TermDate               string       `json:"termDate,omitempty"`
	CompanyId              string       `json:"companyId,omitempty"`
	LicenseExpireEmailId   string       `json:"licenseExpireEmailId, omitempty"`
	MedCardExpireEmailId   string       `json:"medCardExpireEmailId, omitempty"`
	MVRExpireEmailId       string       `json:"mVRExpireEmailId, omitempty"`
	ReviewExpireEmailId    string       `json:"reviewExpireEmailId, omitempty"`
	OneEightyExpireEmailId string       `json:"oneEightyExpireEmailId, omitempty"`
	Status                 DriverStatus `json:"driverStatus,omitempty"`
}

func (d Driver2) GetStatus() string {
	switch d.Status {
	case WORKING:
		return "Working"
	case FIRED:
		return "Fired"
	case TRANSFERED:
		return "Transfered"
	case LAID_OFF:
		return "Laid Off"
	}
	return ""
}

var DRIVER_CONSTS = map[string]DriverStatus{
	"WORKING":    WORKING,
	"FIRED":      FIRED,
	"TRANSFERED": TRANSFERED,
	"LAID_OFF":   LAID_OFF,
}
