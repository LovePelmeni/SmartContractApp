package rest

import (
	"errors"
	"net/http"

	"github.com/LovePelmeni/ContractApp/exceptions"
	"github.com/LovePelmeni/ContractApp/models"
	"github.com/LovePelmeni/ContractApp/transactions/ethereum"
	"github.com/gin-gonic/gin"
)

var Customer models.Customer

// Customer Rest API Endpoints

func LoginRestController(context *gin.Context) {

}

func LogoutRestController(context *gin.Context) {

}

func CreateCustomer(context *gin.Context) {

}

func DeleteCustomer(context *gin.Context) {

}

// Smart Contracts Rest API Endpoints

func PurchaseContractRestController(context *gin.Context) {
	// Rest Controller, Responsible for Purchasing Smart Contracts

	PurchaserId := context.Query("customerId")
	SmartContractId := context.Query("smartContractId")

	Customer := Customer.Get(CustomerId)
	SmartContractManager := ethereum.NewSmartContractManager(Customer.AccountBlockChainAddress)
	ContractTransaction, TransactError := SmartContractManager.TransactSmartContract(SmartContractId)

	// Checking for Errors

	if errors.Is(TransactError, exceptions.SmartContractDoesNotExist) {
		context.AbortWithStatusJSON(http.StatusBadgateway, gin.H{
			"Error": "Oops, Looks Like, it has been already Purchased by Someone Else :("})
	}

	if errors.Is(TransactError, exceptions.UnknownSmartContractError) {
		context.AbortWithStatusJSON(http.StatusBadGateway,
			gin.H{"Error": "Unknown Smart Contract Error, Please Write Feedback to the Support"})
	}

	if TransactError == nil {
		context.JSON(http.StatusCreated, gin.H{"Status": "Smart Contract has been Purchased :)"})
	}
}

func CreateContractRestController(context *gin.Context) {
	// Rest Controller, That is Responsible for Creating Smart Contracts

	CustomerId := context.Query("customerId")
	Customer := Customer.Get(CustomerId)
	SmartContractManager := ethereum.NewSmartContractManager(Customer.AccountBlockChainAddress)
	ContractAddress, Transaction, BoundContract, Error := SmartContractManager.CreateSmartContract()

	if errors.Is(Error, exceptions.UnknownSmartContractError) {
		context.AbortWithStatusJSON(http.StatusBadGateway,
			gin.H{"Error": "Unknown Smart Contract Error, Please Write Feedback to Support"})
	}

	if Error == nil {
		context.JSON(http.StatusCreated,
			gin.H{"Status": "Smart Contract has been Deployed :)"})
	}
}

func RollbackContractRestController(context *gin.Context) {
	// Rest Controller, Responsible for Rolling Back Smart Contract

}
