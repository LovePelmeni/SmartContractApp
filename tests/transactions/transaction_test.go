package transaction_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SmartContractTransactionTestSuite struct {
	suite.Suite
}

func TestSmartContractTransactionSuite(t *testing.T) {
	suite.Run(t, new(SmartContractTransactionTestSuite))
}

func (this *SmartContractTransactionTestSuite) SetupTest() {

}

func (this *SmartContractTransactionTestSuite) Teardown() {

}
