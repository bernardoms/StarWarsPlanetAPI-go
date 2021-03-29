package mock

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type LoggerMock struct {
	mock.Mock
}

func (m *LoggerMock) LogWithFields(req *http.Request, level string, fields map[string]interface{}, message string) {
	m.Called(req, level, fields, message)
}
