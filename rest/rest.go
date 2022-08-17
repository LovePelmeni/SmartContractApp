package rest

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"context"

	"github.com/LovePelmeni/ContractApp/blockchain/transactions/ethereum"
	"github.com/LovePelmeni/ContractApp/exceptions"
	"github.com/LovePelmeni/ContractApp/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

var Customer models.Customer
var SmartContract models.SmartContract

// Customer Rest API Endpoints

func LoginRestController(context *gin.Context) {

}

func LogoutRestController(context *gin.Context) {

}

func CreateCustomerRestController(context *gin.Context) {

}

func ChangePasswordRestController(context *gin.Context) {

}

func DeleteCustomerRestController(context *gin.Context) {

}

// Smart Contracts Rest API Endpoints

func PurchaseContractRestController(context *gin.Context) {
	// Rest Controller, Responsible for Purchasing Smart Contracts

	PurchaserId := context.Query("customerId")
	SmartContractId := context.Query("smartContractId")

	ContractTransactionMetadata, _ := json.Marshal(context.PostForm(""))
	Method := context.PostForm("Method")

	SmartContractManager := ethereum.NewSmartContractManager(Customer.AccountBlockChainAddress)
	ContractTransaction, TransactError := SmartContractManager.TransactSmartContract(
		SmartContractId,
		ContractTransactionMetadata,
		Method,
	)
	// Checking for Errors

	switch {

	// Case Smart Contract, that Purchaser is trying to Buy, Does Not Exist
	case errors.Is(TransactError, exceptions.SmartContractDoesNotExist()):
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"Error": "Oops, Looks Like, it has been already Purchased by Someone Else :("})

	// Case of Some Unknown Error, Occurred, while Performing Smart Contract Transaction
	case errors.Is(TransactError, exceptions.UnknownSmartContractError()):
		context.AbortWithStatusJSON(http.StatusBadGateway,
			gin.H{"Error": "Unknown Smart Contract Error, Please Write Feedback to the Support"})

	// Case Transaction has been Executed Successfully and Contract Transaction is not None
	case TransactError == nil && !reflect.ValueOf(
		ContractTransaction.(ethereum.SmartContractTransaction)).IsNil():

		// Creating New ORM Object of the Smart Contract Transaction and Storing into SQL Database
		Transaction := ContractTransaction.(ethereum.SmartContractTransaction)

		newSmartContractTransaction := models.NewSmartContractTransaction(
			PurchaserId,
			OwnerId,
		)
		_, Saved := newSmartContractTransaction.Save(*newSmartContractTransaction)

		if Saved == nil {
			context.JSON(http.StatusCreated,
				gin.H{"Status": "Smart Contract has been Purchased :)"})

		} else {
			context.JSON(http.StatusBadGateway, gin.H{"Error": "Don't Worry, You successfully Charged Smart Contract," +
				"But You won't Temporarily See Your Payment At Profile"})
		}
	}
}

func CreateContractRestController(RequestContext *gin.Context) {
	// Rest Controller, That is Responsible for Creating Smart Contracts

	// Initializing Smart Contract Parameters....
	CustomerId := RequestContext.Query("customerId")
	ContractMetadata, _ := json.Marshal(RequestContext.PostForm("ContractMetadata"))
	ContractValue, _ := strconv.Atoi(RequestContext.PostForm("ContractValue"))

	Customer := Customer.Get(CustomerId)
	SmartContractManager := ethereum.NewSmartContractManager(Customer.AccountBlockChainAddress)

	TimeoutContext, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer CancelFunc()

	AuthTransactOpts := bind.TransactOpts{
		From:     common.HexToAddress(Customer.AccountBlockChainAddress),
		Nonce:    nil,
		Value:    big.NewInt(int64(ContractValue)),
		GasPrice: nil,
		Context:  TimeoutContext,
	}

	// Creating New Smart Contract...
	ContractBackend := ethereum.NewSmartContractBackend()
	ContractAddress, Transaction, _, CreationError := SmartContractManager.CreateSmartContract(AuthTransactOpts, ContractBackend)

	// Comparing Received Response with different Cases
	switch {
	// In Case there is Some Unknown Errors, Responding with Fixing Advice.
	case CreationError != nil || errors.Is(CreationError, exceptions.UnknownSmartContractError()) && !errors.Is(CreationError, exceptions.SmartContractDoesNotExist()):
		InfoLogger.Printf("Failed to Create Smart Contract, Error: %s", CreationError)
		RequestContext.JSON(http.StatusBadGateway, gin.H{"Error": "Failed to Create Smart Contract, Due to Unknown Error, Send Feedback to Support Please"})

	case CreationError == nil: // In Case Smart Contract has been Deployed Successfully, We
		// Creating New ORM Object and Putting it into the Database

		newSmartContract := models.NewSmartContract(
			ContractAddress.Hex(),
			int(Transaction.Value().Int64()),
			ContractMetadata,
			Transaction.Protected(),
			CustomerId,
		)
		_, Saved := newSmartContract.SaveContract(*newSmartContract)

		if Saved {
			RequestContext.JSON(http.StatusCreated,
				gin.H{"Status": "Smart Contract has been Deployed :)"})
		} else {
			ErrorLogger.Printf("Failed to Create New Smart Contract ORM Object, Error: %s", Saved)
			RequestContext.JSON(http.StatusBadGateway,
				gin.H{"Error": "Don't worry, your Contract has been Deployed, but there is Some Issues," +
					"Locally, So you are temporarily not able to see new Contracts at your Profile"})
		}
	}
}

func RollbackContractRestController(context *gin.Context) {
	// Rest Controller, Responsible for Rolling Back Smart Contract

}
