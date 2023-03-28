package kbhandler

import (
	"testing"

	"github.com/disaster37/go-kibana-rest/v8"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const baseURL = "http://localhost:5601"

type KibanaHandlerTestSuite struct {
	suite.Suite
	kbHandler KibanaHandler
}

func TestKibanahHandlerSuite(t *testing.T) {
	suite.Run(t, new(KibanaHandlerTestSuite))
}

func (t *KibanaHandlerTestSuite) SetupTest() {

	cfg := kibana.Config{
		Address: baseURL,
	}
	client, err := kibana.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	client.Client.SetTransport(httpmock.DefaultTransport)

	t.kbHandler = &KibanaHandlerImpl{
		client: client,
		log:    logrus.NewEntry(logrus.New()),
	}

	httpmock.Activate()

}

func (t *KibanaHandlerTestSuite) BeforeTest(suiteName, testName string) {
	httpmock.Reset()
}
