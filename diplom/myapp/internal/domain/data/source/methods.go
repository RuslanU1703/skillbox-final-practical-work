package source

import (
	"encoding/json"
	"io/ioutil"
	"myapp/internal/domain/data/myerrors"
	"myapp/internal/entity"
	"sort"
	"strconv"
	"strings"
)

const dataFile = "data.json"

func (s *source) GetData() (entity.ResultSetT, error) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	info, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return entity.ResultSetT{}, myerrors.ErrReadData
	}
	data := entity.ResultSetT{}
	err = json.Unmarshal(info, &data)
	if err != nil {
		return entity.ResultSetT{}, myerrors.ErrReadData
	}
	return data, nil
}
func (s *source) AddData(data entity.ResultSetT) error {
	rwLock.Lock()
	defer rwLock.Unlock()

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return myerrors.ErrWriteData
	}
	err = ioutil.WriteFile(dataFile, file, 0644)
	if err != nil {
		return myerrors.ErrWriteData
	}
	return nil
}

func (s *source) GetSms(path string) ([][]entity.SMSData, error) {
	verifiedData := make([]entity.SMSData, 0)
	resultData := make([][]entity.SMSData, 0)
	dataByte, err := readGetResponse(path)
	if err != nil {
		return [][]entity.SMSData{}, err
	}
	data := strings.Split(string(string(dataByte)), "\n")
	for i := range data {
		dataLine := strings.Split(data[i], ";")
		ok, err := dataValidation(dataLine, "sms")
		// проверим правильность запуска валидации ("sms")
		if i == 0 {
			if err != nil {
				return [][]entity.SMSData{}, err
			}
		}

		if ok {
			verifiedData = append(verifiedData,
				entity.SMSData{
					Country:      dataLine[0],
					Bandwidth:    dataLine[1],
					ResponseTime: dataLine[2],
					Provider:     dataLine[3],
				})
		}
	}
	// sort
	sort.Slice(verifiedData, func(i, j int) bool {
		return verifiedData[i].Provider < verifiedData[j].Provider
	})
	providerAscending := make([]entity.SMSData, len(verifiedData))
	copy(providerAscending, verifiedData)
	resultData = append(resultData, providerAscending)
	sort.Slice(verifiedData, func(i, j int) bool {
		return verifiedData[i].Country < verifiedData[j].Country
	})
	resultData = append(resultData, verifiedData)
	return resultData, nil
}
func (s *source) GetMms(path string) ([][]entity.MMSData, error) {
	verifiedData := make([]entity.MMSData, 0)
	resultData := make([][]entity.MMSData, 0)
	dataSlice := make([]entity.MMSData, 0)
	data, err := readGetResponse(path)
	if err != nil {
		return [][]entity.MMSData{}, err
	}

	err = json.Unmarshal(data, &dataSlice)
	if err != nil {
		return [][]entity.MMSData{}, err
	}

	for i := range dataSlice {
		if ok := MMSValidation(dataSlice[i]); ok {
			dataSlice[i].Country = entity.Countries[dataSlice[i].Country]
			verifiedData = append(verifiedData, dataSlice[i])
		}
	}

	sort.Slice(verifiedData, func(i, j int) bool {
		return verifiedData[i].Provider < verifiedData[j].Provider
	})
	providerAscending := make([]entity.MMSData, len(verifiedData))
	copy(providerAscending, verifiedData)
	resultData = append(resultData, providerAscending)
	sort.Slice(verifiedData, func(i, j int) bool {
		return verifiedData[i].Country < verifiedData[j].Country
	})
	resultData = append(resultData, verifiedData)
	return resultData, nil
}
func (s *source) GetVoiceCall(path string) ([]entity.VoiceCallData, error) {
	verifiedData := make([]entity.VoiceCallData, 0)
	dataByte, err := readGetResponse(path)
	if err != nil {
		return []entity.VoiceCallData{}, err
	}
	data := strings.Split(string(string(dataByte)), "\n")

	for i := range data {
		dataLine := strings.Split(data[i], ";")
		ok, err := dataValidation(dataLine, "voiceCall")
		if err != nil {
			return []entity.VoiceCallData{}, err
		}
		if ok {
			ConnectionStability, err := strconv.ParseFloat(dataLine[4], 32)
			if err != nil {
				return []entity.VoiceCallData{}, err
			}
			TTFB, err := strconv.Atoi(dataLine[5])
			if err != nil {
				return []entity.VoiceCallData{}, err
			}
			VoicePurity, err := strconv.Atoi(dataLine[6])
			if err != nil {
				return []entity.VoiceCallData{}, err
			}
			MedianOfCallsTime, err := strconv.Atoi(dataLine[7])
			if err != nil {
				return []entity.VoiceCallData{}, err
			}

			verifiedData = append(verifiedData,
				entity.VoiceCallData{
					Country:             dataLine[0],
					Bandwidth:           dataLine[1],
					ResponseTime:        dataLine[2],
					Provider:            dataLine[3],
					ConnectionStability: float32(ConnectionStability),
					TTFB:                TTFB,
					VoicePurity:         VoicePurity,
					MedianOfCallsTime:   MedianOfCallsTime,
				})
		}
	}

	return verifiedData, nil
}
func (s *source) GetEmail(path string) (map[string][][]entity.EmailData, error) {
	verifiedData := make([]entity.EmailData, 0)
	resultData := make(map[string][][]entity.EmailData)
	dataByte, err := readGetResponse(path)
	if err != nil {
		return nil, err
	}
	data := strings.Split(string(string(dataByte)), "\n")

	for i := range data {
		dataLine := strings.Split(data[i], ";")
		ok, err := dataValidation(dataLine, "email")
		if err != nil {
			return nil, err
		}
		if ok {
			DeliveryTime, err := strconv.Atoi(dataLine[2])
			if err != nil {
				return nil, err
			}

			verifiedData = append(verifiedData,
				entity.EmailData{
					Country:      dataLine[0],
					Provider:     dataLine[1],
					DeliveryTime: DeliveryTime,
				})
		}
	}

	sort.Slice(verifiedData, func(i, j int) bool {
		return verifiedData[i].Country < verifiedData[j].Country
	})
	mySortEmail := mySplitSlice(verifiedData) // [{RU}{ENG}{UK}{RU}] -> [[{RU}{RU}][{UK}][{ENG}]]

	for i := range mySortEmail {
		emailSlice := mySortEmail[i]
		countryMapKey := emailSlice[0].Country
		returnMapValue := make([][]entity.EmailData, 0)

		sort.Slice(emailSlice, func(k, j int) bool {
			return emailSlice[k].DeliveryTime > emailSlice[j].DeliveryTime
		})
		emailFast := make([]entity.EmailData, len(emailSlice))
		copy(emailFast, emailSlice)
		sort.Slice(emailSlice, func(k, j int) bool {
			return emailSlice[k].DeliveryTime < emailSlice[j].DeliveryTime
		})

		if len(emailSlice) >= 3 {
			returnMapValue = append(returnMapValue, emailFast[:3], emailSlice[:3])
		} else {
			returnMapValue = append(returnMapValue, emailFast, emailSlice)
		}
		resultData[countryMapKey] = returnMapValue
	}
	return resultData, nil
}
func (s *source) GetBilling(path string) (entity.BillingData, error) {
	var verifiedBilling entity.BillingData
	dataByte, err := readGetResponse(path)
	if err != nil {
		return entity.BillingData{}, err
	}
	data := strings.Split(string(string(dataByte)), "\n")

	for i := range data {
		sum, err := readBillingData(data[i])
		if err != nil {
			return entity.BillingData{}, err
		}
		bd, err := decodeBillingData(sum)
		if err != nil {
			return entity.BillingData{}, err
		}

		verifiedBilling = bd
	}

	return verifiedBilling, nil
}
func (s *source) GetSupport(path string) ([]int, error) {
	returnSlice := make([]int, 0)
	dataSlice := make([]entity.SupportData, 0)
	data, err := readGetResponse(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &dataSlice)
	if err != nil {
		return nil, err
	}

	ticketSum := 0
	for i := range dataSlice {
		ticketSum += dataSlice[i].ActiveTickets
	}
	switch {
	case ticketSum < 9:
		returnSlice = append(returnSlice, 1)
	case ticketSum > 16:
		returnSlice = append(returnSlice, 3)
	default:
		returnSlice = append(returnSlice, 2)
	}

	supportSpeed := 60.0 / 18.0
	waitTime := supportSpeed * float64(ticketSum)
	returnSlice = append(returnSlice, int(waitTime))

	return returnSlice, nil
}
func (s *source) GetIncident(path string) ([]entity.IncidentData, error) {
	dataSlice := make([]entity.IncidentData, 0)
	data, err := readGetResponse(path)
	if err != nil {
		return []entity.IncidentData{}, err
	}

	err = json.Unmarshal(data, &dataSlice)
	if err != nil {
		return []entity.IncidentData{}, err
	}
	sort.Slice(dataSlice, func(i, j int) bool {
		return dataSlice[i].Status < dataSlice[j].Status
	})
	return dataSlice, nil
}
