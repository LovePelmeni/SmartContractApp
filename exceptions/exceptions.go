package exceptions

import (
	"errors"
)

func SmartContractDoesNotExist() error {
	return errors.New("Smart Contract Does Not Exist.")
}

func UnknownSmartContractError() error {
	return errors.New("Unknown Contract Error")
}
