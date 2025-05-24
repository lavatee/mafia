# **Online Mafia — многопользовательская игра в мафию**  

**Проект**: Веб-приложение для игры в "Мафию" с реальными игроками.  
**Стек**:  Go (Gin), PostgreSQL, MongoDB, Redis, Unit Testing

---

## **🚀 Основной функционал**  
✅ **Создание комнат**
✅ **Друзья игроков**: пользователи могут отправлять запросы в друзья и играть вместе со своими друзьями
✅ **Unit-тесты** для регистрации

---

## **🛠 Технологии**  
| Компонент       | Технологии                          |  
|-----------------|-------------------------------------|  
| **Backend**     | Go 1.22, Gin    |  
| **Базы данных** | PostgreSQL (аккаунты, игры, запросы в друзья), MongoDB (друзья игроков) |  
| **Кеширование** | Redis    |  
| **Тестирование** | `go test`, `gomock` (моки)        |  
| **Деплой**      | Docker          |  

## **Unit-тестирование**
Пример unit-теста для метода регистрации:
```Go
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
```

## **🎯 Дальнейшее развитие**
Логика самой игры: обмен данными через WebSockets, распределение ролей, ход игры
