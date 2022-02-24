package myerrors

import "errors"

var (
	// GETTING DATA ERRORS
	ErrReadFile = errors.New("could not read the file at the given path")

	ErrSendRqst      = errors.New("failed to send request")
	ErrWrongResponse = errors.New("the server did not respond")
	ErrReadResponse  = errors.New("could not read the response body at the given route")

	// VALIDATION ERRORS
	ErrValidation = errors.New("incorrect input data name")

	// BILLING DATA ERRORS
	ErrBilling = errors.New("incorrect input sum; more data than fields")

	// CONVERT ERRORS
	ErrConvert = errors.New("failed to convert one of the string to int")

	// ERROR DURING FURTHER PROCESSING OF DATA
	ErrWriteData = errors.New("failed to write data in file")
	ErrReadData  = errors.New("failed to read data from file")
)
