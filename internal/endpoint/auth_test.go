package endpoint

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lavatee/mafia"
	"github.com/lavatee/mafia/internal/service"
	mock_service "github.com/lavatee/mafia/internal/service/mocks"
	"github.com/magiconair/properties/assert"
)

func TestEndpoint_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, user mafia.User)
	testTable := []struct {
		Name               string
		InputBody          string
		InputUser          mafia.User
		MockBehavior       mockBehavior
		ExpectedStatusCode int
		ExpectedResponse   string
	}{
		{
			Name:      "OK",
			InputBody: `{"name": "test", "email": "test", "password": "test"}`,
			InputUser: mafia.User{Name: "test", Email: "test", Password: "test"},
			MockBehavior: func(s *mock_service.MockAuth, user mafia.User) {
				s.EXPECT().SignUp(user.Email, user.Name, user.Password).Return(1, nil)
			},
			ExpectedStatusCode: 200,
			ExpectedResponse:   `{"id":1}`,
		},
	}
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			auth := mock_service.NewMockAuth(c)
			test.MockBehavior(auth, test.InputUser)
			svc := &service.Service{Auth: auth}
			endp := &Endpoint{service: svc}

			r := gin.New()
			r.POST("/signup", endp.SignUp)
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(test.InputBody))
			r.ServeHTTP(recorder, req)
			assert.Equal(t, recorder.Code, test.ExpectedStatusCode)
			assert.Equal(t, recorder.Body.String(), test.ExpectedResponse)
		})
	}
}
