package source

import (
	"io/ioutil"
	"math"
	"myapp/internal/domain/data/myerrors"
	"myapp/internal/entity"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// GET DATA HELPER FUNC
func readGetResponse(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, myerrors.ErrSendRqst
	}

	if resp.StatusCode != 200 {
		return nil, myerrors.ErrWrongResponse
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, myerrors.ErrReadResponse
	}
	defer resp.Body.Close()

	return data, nil
}

// VALIDATION DATA HELPER FUNC
func dataValidation(data []string, dataName string) (bool, error) {
	switch dataName {
	case "sms":
		if len(data) != 4 {
			return false, nil
		}
		if _, ok := entity.ProvidersMS[data[3]]; !ok {
			return false, nil
		}
	case "email":
		if len(data) != 3 {
			return false, nil
		}
		if _, ok := entity.ProvidersEmail[data[1]]; !ok {
			return false, nil
		}
	case "voiceCall":
		if len(data) != 8 {
			return false, nil
		}
		if _, ok := entity.ProvidersCall[data[3]]; !ok {
			return false, nil
		}
	default:
		return false, myerrors.ErrValidation
	}

	if _, ok := entity.Countries[data[0]]; !ok {
		return false, nil
	}
	// change country name for SMS
	if dataName == "sms" {
		data[0] = entity.Countries[data[0]]
	}

	return true, nil
}
func MMSValidation(data entity.MMSData) bool {
	if fields := reflect.ValueOf(data).NumField(); fields != 4 {
		return false
	}
	if _, ok := entity.Countries[data.Country]; !ok {
		return false
	}

	if _, ok := entity.ProvidersMS[data.Provider]; !ok {
		return false
	}

	return true
}

// BILLING HELPER FUNC
func readBillingData(data string) (uint8, error) {
	// reverse data
	runes := []rune(data)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	dataIntSlice, err := stringToSliceInt(string(runes))
	if err != nil {
		return 0, err
	}
	var sum uint8
	// interpreted as num
	for i := range dataIntSlice {
		if float64(dataIntSlice[i]) == 0 {
			continue
		}
		sum += uint8(math.Pow(2, float64(i)))
	}
	return sum, nil
}
func decodeBillingData(sum uint8) (entity.BillingData, error) {
	data := entity.BillingData{}

	s := strconv.FormatInt(int64(sum), 2)
	dataIntSlice, err := stringToSliceInt(s)
	if err != nil {
		return entity.BillingData{}, err
	}

	for i := range dataIntSlice {
		switch i {
		case 0:
			if dataIntSlice[i] == 1 {
				data.CreateCustomer = true
			}
		case 1:
			if dataIntSlice[i] == 1 {
				data.Purchase = true
			}
		case 2:
			if dataIntSlice[i] == 1 {
				data.Payout = true
			}
		case 3:
			if dataIntSlice[i] == 1 {
				data.Recurring = true
			}
		case 4:
			if dataIntSlice[i] == 1 {
				data.FraudControl = true
			}
		case 5:
			if dataIntSlice[i] == 1 {
				data.CheckoutPage = true
			}
		default:
			return entity.BillingData{}, myerrors.ErrBilling
		}
	}

	return data, nil
}

// CONVERT HELPER FUNC
func stringToSliceInt(s string) ([]int, error) {
	returnSlice := strings.Split(string(s), "")
	dataIntSlice := make([]int, 0)

	for i := range returnSlice {
		intElem, err := strconv.Atoi(returnSlice[i])
		if err != nil {
			return nil, myerrors.ErrConvert
		}
		dataIntSlice = append(dataIntSlice, intElem)
	}

	return dataIntSlice, nil
}

// Email helper func
func mySplitSlice(inputSlice []entity.EmailData) [][]entity.EmailData {
	fixedCountry := ""
	var returnSlice []entity.EmailData
	var bigReturnSlice [][]entity.EmailData
	for i := range inputSlice {
		// Массив для страны уже начал собираться?
		switch inputSlice[i].Country {
		case fixedCountry:
			returnSlice = append(returnSlice, inputSlice[i])
			if i == len(inputSlice)-1 {
				bigReturnSlice = append(bigReturnSlice, returnSlice)
			}
			continue
		default:
			// Начать собирать
			switch {
			case len(returnSlice) > 0:
				// Есть данные
				bigReturnSlice = append(bigReturnSlice, returnSlice)
				returnSlice = []entity.EmailData{}
				fallthrough
			default:
				// Нет данных
				fixedCountry = inputSlice[i].Country
				returnSlice = append(returnSlice, inputSlice[i])
			}
		}
	}
	return bigReturnSlice
}
