package service

import (
	"log"
	"myapp/internal/entity"
	"sync"
)

const myRoute = "https://salty-inlet-33171.herokuapp.com/"
const smsPath = myRoute + "static/files/sms.data"
const mmsPath = myRoute + "mms"
const voicePath = myRoute + "static/files/voice.data"
const emailPath = myRoute + "static/files/email.data"
const billingPath = myRoute + "static/files/billing.data"
const supportPath = myRoute + "support"
const accendentPath = myRoute + "accendent"

func (s *service) Collect() (entity.ResultSetT, error) {
	data := entity.ResultSetT{}
	var collectError error = nil
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		defer wg.Done()
		sms, err := s.source.GetSms(smsPath)
		if err != nil {
			collectError = err
			return
		}
		data.SMS = sms
	}()
	go func() {
		defer wg.Done()
		mms, err := s.source.GetMms(mmsPath)
		if err != nil {
			collectError = err
			return
		}
		data.MMS = mms
	}()
	go func() {
		defer wg.Done()
		vcall, err := s.source.GetVoiceCall(voicePath)
		if err != nil {
			collectError = err
			return
		}
		data.VoiceCall = vcall
	}()
	go func() {
		defer wg.Done()
		email, err := s.source.GetEmail(emailPath)
		if err != nil {
			collectError = err
			return
		}
		data.Email = email
	}()
	go func() {
		defer wg.Done()
		billing, err := s.source.GetBilling(billingPath)
		if err != nil {
			collectError = err
			return
		}
		data.Billing = billing
	}()
	go func() {
		defer wg.Done()
		support, err := s.source.GetSupport(supportPath)
		if err != nil {
			collectError = err
			return
		}
		data.Support = support
	}()
	go func() {
		defer wg.Done()
		incident, err := s.source.GetIncident(accendentPath)
		if err != nil {
			collectError = err
			return
		}
		data.Incidents = incident
	}()
	wg.Wait()
	if collectError != nil {
		return entity.ResultSetT{}, collectError
	}
	err := s.source.AddData(data)
	if err != nil {
		// non-critical error
		log.Println("Полученные данные не записаны в файл, err: ", err.Error())
	}
	return data, collectError
}
func (s *service) ShowData() (entity.ResultSetT, error) {
	data, err := s.source.GetData()
	if err != nil {
		log.Println("файл пуст или поврежден")
		return entity.ResultSetT{}, err
	}
	return data, nil
}
