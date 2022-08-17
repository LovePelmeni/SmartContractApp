package models

import (
	"errors"
	"fmt"
	"regexp"

	"log"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DATABASE_USER     = os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")
	DATABASE_NAME     = os.Getenv("DATABASE_NAME")
	DATABASE_PORT     = os.Getenv("DATABASE_PORT")
	DATABASE_HOST     = os.Getenv("DATABASE_HOST")
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

var (
	Database *gorm.DB
)

var smartContract SmartContract
var customer Customer

func InitializeDatabase() {
	Db, Error := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", DATABASE_HOST, DATABASE_PORT,
			DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME),
	}))
	if Error != nil {
		fmt.Sprint(
			"Please Check, that you've correctly set up Database Environment Variables, Otherwise it would break up")
		panic(fmt.Sprintf("Database Error: %s", Error))
	}
	Database = Db
}

func init() {
	LogFile, LogError := os.OpenFile("Database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Llongfile|log.Lshortfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	if LogError != nil {
		panic(LogError)
	}
}

type BaseModelValidatorInterface interface {
	// Interface, represents Class, That Validates Model Input Data, before saving
	// it to the Database
	GetRegexValidationPatterns() map[string]string
	Validate(InputData interface{}) (bool, []error)
}

// Validator, for Validating Customer SQL Model

type CustomerModelValidator struct {
	BaseModelValidatorInterface
	Patterns map[string]string
}

func NewBaseModelValidator() *CustomerModelValidator {
	return &CustomerModelValidator{}
}

func (this *CustomerModelValidator) Validate(InputData interface{}) (bool, []error) {
	var ValidationErrors []error
	Patterns := this.GetRegexValidationPatterns()
	for Field, Value := range InputData.(map[string]string) {
		if Matches, Error := regexp.MatchString(Patterns[Field], Value); Matches != true || Error != nil {
			ValidationErrors = append(ValidationErrors,
				errors.New(fmt.Sprintf("Invalid Value for Field `%s`", Field)))
		}
	}
	if len(ValidationErrors) != 0 {
		return false, ValidationErrors
	} else {
		return true, nil
	}
}

func (this *CustomerModelValidator) GetRegexValidationPatterns() map[string]string {
	return map[string]string{}
}

// Validator, for Validating Smart Contract Model

type SmartContractModelValidator struct{}

func NewSmartContractModelValidator() *SmartContractModelValidator {
	return &SmartContractModelValidator{}
}

func (this *SmartContractModelValidator) GetRegexValidationPatterns() map[string]string {
	return map[string]string{
		"Cost":        "^[0-9]{1,}$",
		"PurchaserId": "^[0-9]{1,}$",
		"OwnerId":     "^[0-9]{1,}$",
	}
}

func (this *SmartContractModelValidator) Validate(InputData interface{}) (bool, []error) {
	var ValidationErrors []error
	Patterns := this.GetRegexValidationPatterns()
	for Field, Value := range InputData.(map[string]string) {
		if Matches, Error := regexp.MatchString(Patterns[Field], Value); Matches != true || Error != nil {
			ValidationErrors = append(ValidationErrors,
				errors.New(fmt.Sprintf("Invalid Value for Field `%s`", Field)))
		}
	}
	if len(ValidationErrors) != 0 {
		return false, ValidationErrors
	} else {
		return true, nil
	}
}

type Customer struct {
	gorm.Model
	AccountBlockChainAddress string `json:"AccountBlockChainAddress" gorm:"type:varchar(1000); not null;unique;"`
	Username                 string `json"Username" gorm:"type:varchar(100);not null;unique;"`
	Email                    string `json:"Email" gorm:"type:varchar(100); not null;unique;`
	Password                 string `json:"Password" gorm:"type:varchar(100); not null;`
}

func NewCustomer(Username string, Email string) *Customer {
	return &Customer{
		Username: Username,
		Email:    Email,
	}
}

func (this *Customer) AfterSave(Obj *gorm.DB) error {
	// performing email notification about event
}

func (this *Customer) AfterDelete() {
	// performing email notification about event
}

func (this *Customer) Create() error {
	// Method, that creates Customer Profile
}

func (this *Customer) ChangePassword(Password ...string) {
	// Method, that updates Customer Profile
}

func (this *Customer) Delete() error {
	// Method, that deletes Customer Profile
}

func (this *Customer) Get(CustomerId string) Customer {
	// Method, that returns Customer Profile Object
}

type SmartContract struct {
	gorm.Model
	Address      string `json:"Address" gorm:"type:varchar(100); not null;"`
	Cost         int    `json:"Cost" gorm:"type:integer;not null;"`
	InputData    []byte `json:"InputData; omitempty;" gorm:"type:varchar(1000);default:null;"`
	Protected    bool   `json:"Protected" gorm:"type:boolean;default:false;"`
	OwnerId      string `json:"OwnerId" gorm:"type:varchar(100); not null;"`
	Rollbackable bool   `json:"Rollbackable" gorm"type:boolean;default:true;"` // this Field will be changed to false, after Smart Contract has been Purchased.
}

func NewSmartContract(Address string, Cost int, InputData []byte, Protected bool, OwnerId string) *SmartContract {
	return &SmartContract{
		Address:   Address,
		InputData: InputData,
		Cost:      Cost,
		Protected: Protected,
		OwnerId:   OwnerId,
	}
}

func (this *SmartContract) Save(SmartContractObj SmartContract) (*SmartContract, bool) {
	var SmartContract SmartContract
	// Initializing Creation Transaction for Smart Contract...
	Saved := Database.Model(smartContract).Create(&SmartContractObj)
	if Saved.Error != nil {
		Saved.Rollback()
		DebugLogger.Printf(
			"Failed to Create New Smart Contract ORM Object, GORM ERROR: %s", Saved.Error)
		return nil, false
	}
	// Committing Transaction
	Committed := Saved.Commit().Find(&SmartContract)

	// If Any Errors Occurred, Rolling Back and returns Error
	if Committed.Error != nil {
		DebugLogger.Printf(
			"Failed to Commit Smart Contract ORM Object: Error: %s", Committed.Error)
		return nil, false
	} else {
		return &SmartContract, true
	}
}

func (this *SmartContract) Get(SmartContractId string) (SmartContract, error) {
	// Returning object of the Smart Contract...
}

func (this *SmartContract) Rollback(smartContractObject *gorm.DB) (bool, error) {
	// Rolling Back Creation of the Smart Contract, if it's possible and `Rollbackable` is set to True

	// Getting Info About Smart Contract
	var smartContract SmartContract
	smartContractObject.Find(&smartContract)

	if smartContract.Rollbackable != true {
		return false, errors.New("Rollback is Not Allowed, Contract is Already Purchased")
	} else {
		RolledBack := smartContractObject.Rollback()
		if RolledBack.Error != nil {
			DebugLogger.Printf(
				"Failed to Rollback Smart Contract Transaction, Error: %s", RolledBack.Error)
			return false, RolledBack.Error
		}
		return true, nil
	}
}

// Annotation Methods

func (this *SmartContract) GetOwnedSmartContracts(CustomerId string) []SmartContract {
	var smartContracts []SmartContract
	ObtainError := Database.Model(&SmartContract{}).Where(
		"owner_id = ?", CustomerId).Find(&smartContracts)
	if ObtainError != nil {
		DebugLogger.Printf(
			"Failed to Obtain Owner Smart Contracts, Error: %s", ObtainError.Error)
	}
	return smartContracts
}

func (this *SmartContract) GetPurchasedSmartContracts(CustomerId string) []SmartContract {
	var smartContracts []SmartContract
	ObtainError := Database.Model(&SmartContract{}).Where(
		"purchaser_id = ?", CustomerId).Find(&smartContracts)
	if ObtainError.Error != nil {
		DebugLogger.Printf(
			"Failed to Obtain Purchased Smart Contracts, Error: %s", ObtainError.Error)
	}
	return smartContracts
}

type SmartContractTransaction struct {
	gorm.Model
}

func NewSmartContractTransaction() *SmartContractTransaction {
	return &SmartContractTransaction{}
}

func (this *SmartContractTransaction) GetCustomerTransactions(CustomerId string) []SmartContractTransaction

func (this *SmartContractTransaction) Get(transactionId string) SmartContractTransaction

func (this *SmartContractTransaction) Create() bool

func (this *SmartContractTransaction) Rollback() bool
