package middlewares

import (
	"net/http"

	"github.com/LovePelmeni/ContractApp/models"
	"github.com/gin-gonic/gin"
)

func JwtAuthenticationMiddleware() gin.HandlerFunc {
	// Middleware for Validating Jwt Auth Tokens
	return func(context *gin.Context) {
	}
}

func CheckSmartContractIsRollbackableMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {

		var smartContract models.SmartContract
		SmartContractId := context.Query("smartContractId")
		CustomerId := context.Query("customerId")

		SmartContract := models.Database.Model(
			&models.SmartContract{}).Where("id = ? AND owner_id = ?", SmartContractId, CustomerId).Find(&smartContract)

		switch {

		case SmartContract.Error != nil:
			context.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"Error": "You are not Owner of this Smart Contract"})
			return

		case smartContract.Rollbackable != true:
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "This Contract has been Already Purchased, You are not able to Roll It Back"})
			return

		default:
			context.Next()
		}
	}
}

func CheckPurchaserIsNotSmartContractOwnerMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Checking If Purchaser not trying to Purchase His Own Smart Contract
		PurchaserId := context.Query("PurchaserId")
		SmartContract := models.Database.Model(
			&models.SmartContract{}).Where("owner_id = ?", PurchaserId)
		switch {

		case SmartContract != nil:
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "It is not Allowed, to Purchase Your Own Smart Contract"})
			return

		default:
			context.Next()
		}
	}
}
