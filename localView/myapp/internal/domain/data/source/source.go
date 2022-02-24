package source

import (
	"myapp/internal/entity"
	"sync"
)

var (
	rwLock sync.RWMutex
)

type Source interface {
	AddData(entity.ResultSetT) error
	GetData() (entity.ResultSetT, error)

	GetSms(string) ([][]entity.SMSData, error)
	GetMms(string) ([][]entity.MMSData, error)
	GetVoiceCall(string) ([]entity.VoiceCallData, error)
	GetEmail(string) (map[string][][]entity.EmailData, error)
	GetBilling(string) (entity.BillingData, error)
	GetSupport(string) ([]int, error)
	GetIncident(string) ([]entity.IncidentData, error)
}

type source struct {
}

func New() Source {
	return &source{}
}
