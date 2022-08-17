package ethereum

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/LovePelmeni/ContractApp/models"
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

var StoreABIObject = abi.ABI{
	Methods: map[string]abi.Method{},
	Events:  map[string]abi.Event{},
	Errors:  map[string]abi.Error{},
}

var StoreABI = os.Getenv("SMART_CONTRACT_STORE_ABI")
var StoreBin = os.Getenv("SMART_CONTRACT_STORE_BIN")

func init() {
	LogFile, Error := os.OpenFile("Transactions.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Error != nil {
		panic(Error)
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ltime|log.Ldate|log.Llongfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ltime|log.Ldate|log.Llongfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ltime|log.Ldate|log.Llongfile)
}

type SmartContractParametersInterface interface {
	// Interface, that represents Smart Contract Parameters
}
type SmartContractInterface interface {
	// Interface, that represents Smart Contract
}
type SmartContractTransactionInterface interface {
	// Interface, represents Interface of the Smart Contract Transaction
}
type SmartContractTransactionManagerInterface interface {
	// Class, represents Interface, for managing Smart Contract Transactions
	SaveContractTransaction(Transaction SmartContractTransactionInterface) bool
}





type SmartContractParameters struct {
	SmartContractParametersInterface
}

func NewSmartContractParameters() *SmartContractParameters {
	return &SmartContractParameters{}
}

// Smart Contract Class, that represents Deployed Smart Contract
type SmartContract struct {
	SmartContractInterface
	Address string `json:"Address"`
	OwnerId string `json:"OwnerId"`
}

func NewSmartContract() *SmartContract {
	return &SmartContract{}
}

// Smart Contract Class, that represents Purchased Smart Conract
type SmartContractTransaction struct {
	SmartContractTransactionInterface
	TransactionInputData []byte          `json:"TransactionInputData"`
	Cost                 int             `json:"Cost"`
	Purchaser            models.Customer `json:"Purchaser"`
	Owner                models.Customer `json:"Owner"`
}

func NewSmartContractTransaction(Data []byte, Cost int,
	Purchaser models.Customer, Owner models.Customer) *SmartContractTransaction {
	return &SmartContractTransaction{
		TransactionInputData: Data,
		Cost:                 Cost,
		Purchaser:            Purchaser,
		Owner:                Owner,
	}
}

// Smart Contract Store Caller

type SmartContractStoreCaller struct {
	bind.ContractCaller
	// Class, that will be call the Store of the blocks, to
}

func NewSmartContractStoreCaller() *SmartContractStoreCaller {
	return &SmartContractStoreCaller{}
}

func (this *SmartContractStoreCaller) CallSmartContractStore(
	TransactionOptions *bind.CallOpts,
	Contract *bind.BoundContract,
	results *[]interface{},
	ContractParameters interface{},
	Method string,

) error {
	// Performing Smart Contract Call...
	Error := Contract.Call(TransactionOptions, results, Method, ContractParameters)
	return Error
}

// Smart Contract Transfer Class
type SmartContractStoreFilterer struct {
	bind.ContractFilterer
	// Class, that will be Handling Transfer Operation of the Smart Contract to the Other's Account.
}

func NewSmartContractStoreFilterer() *SmartContractStoreFilterer {
	return &SmartContractStoreFilterer{}
}

func (this *SmartContractStoreFilterer) Transfer(Contract *bind.BoundContract, Options *bind.TransactOpts) (*types.Transaction, error) {
	Transaction, Error := Contract.Transfer(Options)
	return Transaction, Error
}

// Smart Contract Transaction Class

type SmartContractStoreTransactor struct {
	bind.ContractTransactor
	// Class, that will be implementing Transaction
	// on the Remote Block of the BlockChain
}

func NewSmartContractStoreTransactor() *SmartContractStoreTransactor {
	return &SmartContractStoreTransactor{}
}

func (this *SmartContractStoreTransactor) Transact(
	Options *bind.TransactOpts, Contract *bind.BoundContract, Method string,
	ContractParameters SmartContractParametersInterface) (*types.Transaction, error) {
	Transaction, Error := Contract.Transact(Options, Method, ContractParameters)
	return Transaction, Error
}

type SmartContractBackend struct {
	// Structure, represents Smart Contract Backend
	SmartContractStoreCaller
	SmartContractStoreTransactor
	SmartContractStoreFilterer
}

func NewSmartContractBackend() *SmartContractBackend {
	return &SmartContractBackend{}
}

// Smart Contract Manager Class

type SmartContractManager struct {
	// Class, that rules all of the Operations Through the Smart Contracts
	SmartContractStoreTransactor
	SmartContractStoreCaller

	Contract *bind.BoundContract
}

func NewSmartContractManager(

	AccountBlockChainAddress string, // NFT Token of Customer, In Order to Perform Any Operations with Smart Contracts

) *SmartContractManager {

	SmartContractCaller := NewSmartContractStoreCaller()
	SmartContractStoreTransactor := NewSmartContractStoreTransactor()
	SmartContractFilterer := NewSmartContractStoreFilterer()

	Contract := bind.NewBoundContract(
		common.HexToAddress(AccountBlockChainAddress),
		StoreABIObject,
		SmartContractCaller,          // Interface, that represents Smart Contract Caller (see above)
		SmartContractStoreTransactor, // Interface, that represents Contract Transactor (see above)
		SmartContractFilterer,
	)

	return &SmartContractManager{
		Contract: Contract,
	}
}

func (this *SmartContractManager) TransactSmartContract(
	// params of the Smart Contract Target
	AuthCredentials *bind.TransactOpts,
	SmartContractId string,
	ContractParameters SmartContractParametersInterface,
	Method string,
) (SmartContractTransactionInterface, error) {

	Transaction, TransactError := this.Transact(AuthCredentials,
		this.Contract, Method, ContractParameters)

	// Checking if the Contract performed Well, and there is no Errors

	if TransactError != nil || Transaction == nil {
		return nil, TransactError
	}

	var SmartContractPurchaser models.Customer // Contract Purchaser
	var SmartContractOwner models.Customer     // Contract Owner

	// Obtaining the Owner Of the Smart Contract, by

	models.Database.Model(&models.Customer{}).Where(
		"blockchain_account_address = ?", AuthCredentials.From.Hash()).Find(&SmartContractOwner)

	// Obtaining the Purchaser of the Smart Contract
	models.Database.Model(&models.Customer{}).Where(
		"blockchain_account_address = ?", Transaction.To().Hash()).Find(&SmartContractPurchaser)

	// Initializing Smart Contract Model
	newTransaction := NewSmartContractTransaction(
		Transaction.Data(), int(Transaction.Cost().Int64()),
		SmartContractPurchaser,
		SmartContractOwner,
	)
	return newTransaction, nil
}

func (this *SmartContractManager) CreateSmartContract(AuthCredentials bind.TransactOpts,
	ContractBackend bind.ContractBackend) (common.Address, *types.Transaction, *bind.BoundContract, error) {

	// Parsing Smart Contract Credentials, in order to perform Operation

	ParsedStoreABI, Error := abi.JSON(strings.NewReader(StoreABI))
	if Error != nil {
		ErrorLogger.Printf("Failed to Parse Store ABI")
		return common.Address{}, nil, nil, errors.New("Invalid Blockchain Store Credentials")
	}

	// Deploying Contrat
	ContractAddress, Transaction, Contract, ContractError := bind.DeployContract(
		&AuthCredentials, ParsedStoreABI, common.FromHex(StoreBin), ContractBackend)
	if ContractError != nil {
		ErrorLogger.Printf("Failed to Deploy Contract")

		// Returning Smart Contract Deployed Credentials
		return common.Address{}, &types.Transaction{}, nil,
			errors.New("Contract Deploy Failure")
	} else {
		return ContractAddress, Transaction, Contract, nil
	}
}
