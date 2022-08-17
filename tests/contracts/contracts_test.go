package contracts_test

// NOTE: THIS TESTS IS RUNNING AT THE BUILDING STAGE 
// IT DOES NOT BELONG TO THE INTEGRATION TESTS!

import (
	"testing"
	"github.com/stretchr/testify/suite"
)

type SmartContractTestSuite struct {
	suite.Suite 
}
func TestSmartContractTestSuite(t *testing.T) {
	suite.Run(t, new(SmartContractTestSuite))
}

func (this *SmartContractTestSuite) TestSmartContract(t *testing.T) {
	// Tests Main Functionality of the Smart Contracts, that has been Implemented 
}

