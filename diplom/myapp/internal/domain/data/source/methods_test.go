package source

import (
	"encoding/json"
	"log"
	"myapp/internal/entity"
	"net"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

const testServerURL = "http://127.0.0.1:61972/"

var (
	testSource = New()
)

func TestMain(m *testing.M) {
	// Setup
	var (
		mmsObj = entity.MMSData{
			Country:      "RU",
			Provider:     "Rond",
			Bandwidth:    "50",
			ResponseTime: "1935",
		}
		mmsObj2 = entity.MMSData{
			Country:      "BG",
			Provider:     "Topolo",
			Bandwidth:    "58",
			ResponseTime: "1755",
		}
		supportObj = entity.SupportData{
			Topic:         "SMS",
			ActiveTickets: 3,
		}
		accendentObj = entity.IncidentData{
			Topic:  "SMS delivery in EU",
			Status: "closed",
		}
		accendentObjSecond = entity.IncidentData{
			Topic:  "MMS connection stability",
			Status: "active",
		}
	)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/static/files/sms.data", func(c *gin.Context) {
		c.Writer.Write([]byte("RU;45;411;Topolo\nUS;50;1716;Rond"))
	})
	router.GET("/mms", func(c *gin.Context) {
		dataSlice := make([]entity.MMSData, 0)
		dataSlice = append(dataSlice, mmsObj, mmsObj2)
		b, _ := json.Marshal(dataSlice)
		c.Writer.Write(b)
	})
	router.GET("/static/files/voice.data", func(c *gin.Context) {
		c.Writer.Write([]byte("RU;79;1317;TransparentCalls;0.88;154;9;33\nUS;23;830;E-Voice;0.94;540;85;19"))
	})
	router.GET("/static/files/email.data", func(c *gin.Context) {
		c.Writer.Write([]byte("RU;Gmail;373\nRU;Yahoo;103\nGB;Gmail;173\nGB;Yandex;511\nGB;Mail.ru;120"))
	})
	router.GET("/static/files/billing.data", func(c *gin.Context) {
		c.Writer.Write([]byte("101111"))
	})
	router.GET("/support", func(c *gin.Context) {
		dataSlice := make([]entity.SupportData, 0)
		dataSlice = append(dataSlice, supportObj)
		b, _ := json.Marshal(dataSlice)
		c.Writer.Write(b)
	})
	router.GET("/accendent", func(c *gin.Context) {
		dataSlice := make([]entity.IncidentData, 0)
		dataSlice = append(dataSlice, accendentObj, accendentObjSecond)
		b, _ := json.Marshal(dataSlice)
		c.Writer.Write(b)
	})

	l, err := net.Listen("tcp", "127.0.0.1:61972")
	if err != nil {
		log.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(router)
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	log.Println("Start test server on PORT: ", ts.URL)

	// Tests running
	exitVal := m.Run()

	// Finish
	ts.Close()
	os.Exit(exitVal)
}
func TestGetSms(t *testing.T) {
	expectedData := `[[{"country":"United States of America","bandwidth":"50","response_time":"1716","provider":"Rond"},{"country":"Russian Federation","bandwidth":"45","response_time":"411","provider":"Topolo"}],[{"country":"Russian Federation","bandwidth":"45","response_time":"411","provider":"Topolo"},{"country":"United States of America","bandwidth":"50","response_time":"1716","provider":"Rond"}]]`
	res, err := testSource.GetSms(testServerURL + "static/files/sms.data")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetMms(t *testing.T) {
	expectedData := `[[{"country":"Russian Federation","provider":"Rond","bandwidth":"50","response_time":"1935"},{"country":"Bulgaria","provider":"Topolo","bandwidth":"58","response_time":"1755"}],[{"country":"Bulgaria","provider":"Topolo","bandwidth":"58","response_time":"1755"},{"country":"Russian Federation","provider":"Rond","bandwidth":"50","response_time":"1935"}]]`
	res, err := testSource.GetMms(testServerURL + "mms")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetVoiceCall(t *testing.T) {
	expectedData := `[{"country":"RU","bandwidth":"79","response_time":"1317","provider":"TransparentCalls","connection_stability":0.88,"ttfb":154,"voice_purity":9,"median_of_call_time":33},{"country":"US","bandwidth":"23","response_time":"830","provider":"E-Voice","connection_stability":0.94,"ttfb":540,"voice_purity":85,"median_of_call_time":19}]`
	res, err := testSource.GetVoiceCall(testServerURL + "static/files/voice.data")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetEmail(t *testing.T) {
	expectedData := `{"GB":[[{"country":"GB","provider":"Yandex","delivery_time":511},{"country":"GB","provider":"Gmail","delivery_time":173},{"country":"GB","provider":"Mail.ru","delivery_time":120}],[{"country":"GB","provider":"Mail.ru","delivery_time":120},{"country":"GB","provider":"Gmail","delivery_time":173},{"country":"GB","provider":"Yandex","delivery_time":511}]],"RU":[[{"country":"RU","provider":"Gmail","delivery_time":373},{"country":"RU","provider":"Yahoo","delivery_time":103}],[{"country":"RU","provider":"Yahoo","delivery_time":103},{"country":"RU","provider":"Gmail","delivery_time":373}]]}`
	res, err := testSource.GetEmail(testServerURL + "static/files/email.data")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetBilling(t *testing.T) {
	expectedData := `{"create_customer":true,"purchase":false,"payout":true,"recurring":true,"fraud_control":true,"checkout_page":true}`
	res, err := testSource.GetBilling(testServerURL + "static/files/billing.data")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetSupport(t *testing.T) {
	expectedData := `[1,10]`
	res, err := testSource.GetSupport(testServerURL + "support")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
func TestGetIncident(t *testing.T) {
	expectedData := `[{"topic":"MMS connection stability","status":"active"},{"topic":"SMS delivery in EU","status":"closed"}]`
	res, err := testSource.GetIncident(testServerURL + "accendent")
	if err != nil {
		t.Fail()
	}
	resultData, err := json.Marshal(res)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if expectedData != string(resultData) {
		t.Errorf("Want: %s Got: %s", expectedData, string(resultData))
		t.Fail()
	}
}
