package controller

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/challenge/pkg/models"
)

const (
	// Query Params
	Recipient = "recipient"
	Start     = "start"
	Limit     = "limit"

	// Error messages
	MandatoryField = "%v is a mandatory field"
	InvalidNumber  = "invalid %v number: %v"
	MissingRecord  = "record not found"

	UnknownError = "Error retrieving user information"
)

//validateUserData validates that the given input data is not empty
func validateUserData(user models.User) error {
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Password == "" || len(user.Password) < 8 {
		return errors.New("password should be at least 8 characters")
	}
	return nil
}

func validateMessageData(message models.Message) error {
	if message.UserID == 0 {
		return errors.New("sender cannot be empty")
	}
	if message.Recipient == 0 {
		return errors.New("recipient cannot be empty")
	}

	return nil
}

// checkInt64QueryParams reads a query param and return the corresponding int64
func checkInt64QueryParams(paramName string, paramValue string) (int64, error) {
	parsedInt64, err := strconv.ParseInt(paramValue, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(InvalidNumber, paramName, err)

	} else if parsedInt64 < 0 {
		return 0, fmt.Errorf(InvalidNumber, paramName, parsedInt64)
	}
	return parsedInt64, nil
}

// checkIntQueryParams reads a query param and return the corresponding int64
func checkIntQueryParams(paramName string, paramValue string) (int, error) {
	parsedInt, err := strconv.Atoi(paramValue)
	if err != nil {
		return 0, fmt.Errorf(InvalidNumber, paramName, err)

	} else if parsedInt < 0 {
		return 0, fmt.Errorf(InvalidNumber, paramName, parsedInt)
	}
	return parsedInt, nil
}

// getQueryParams given the query params of  GET /messages, checks if the given values are valid or not, and assigns it
func getQueryParams(queryParams url.Values, filter *models.MessageFilter) error {
	recipientValue := queryParams.Get(Recipient)
	startValue := queryParams.Get(Start)
	limitValue := queryParams.Get(Limit)

	if recipientValue == "" {
		return fmt.Errorf(MandatoryField, Recipient)
	}

	if startValue == "" {
		return fmt.Errorf(MandatoryField, Start)
	}

	var errGettingRecipient error
	filter.Recipient, errGettingRecipient = checkInt64QueryParams(Recipient, recipientValue)
	if errGettingRecipient != nil {
		return errGettingRecipient
	}

	var errGettingStartID error
	filter.Start, errGettingStartID = checkInt64QueryParams(Start, startValue)
	if errGettingStartID != nil {
		return errGettingStartID
	}

	if limitValue != "" {
		var errGettingLimit error
		filter.Limit, errGettingLimit = checkIntQueryParams(Limit, limitValue)
		if errGettingLimit != nil {
			return errGettingLimit
		}
	}

	return nil
}
