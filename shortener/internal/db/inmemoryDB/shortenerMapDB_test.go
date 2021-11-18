package inmemoryDB

//func TestShortnerMapDB_CreateShort(t *testing.T) {
//	cases := map[string]struct {
//		shortener *shortenerBL.Shortener
//		expected *shortenerBL.Shortener
//		errorMessage string
//	}{
//		"Test 1": {
//			// Создаём тестовый shortner.
//			shortener: shortenerBL.Shortener{
//				FullLink: "http://test.local/test",
//				ID: uuid.Parse("9ba87644-294d-4f3e-b5df-51c9d8c330b3"),
//			},
//			//// Создаем запрос с указанием нашего хендлера. Так как мы тестируем GET-эндпоинт
//			//// то нам не нужно передавать тело, поэтому третьим аргументом передаем nil
//			//req: httptest.NewRequest("GET", "/home/?name=John", nil),
//			//// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
//			//// и используем его для получения ответа
//			//rr: httptest.NewRecorder(),
//			//// Указываем какой хендлер будем тестировать
//			//handler: &handlers.Handler{},
//			//// Прогнозируемы код ответа
//			//status: http.StatusOK,
//			//// тело ответа
//			//expected: `Parsed query-param with key "name": John`,
//		},
//		"Test 2": {
//			// Создаем запрос с указанием нашего хендлера. Так как мы тестируем GET-эндпоинт
//			// то нам не нужно передавать тело, поэтому третьим аргументом передаем nil
//			req: httptest.NewRequest("GET", "/home/?name=Вано", nil),
//			// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
//			// и используем его для получения ответа
//			rr: httptest.NewRecorder(),
//			// Указываем какой хендлер будем тестировать
//			handler: &handlers.Handler{},
//			// Прогнозируемы код ответа
//			status: http.StatusOK,
//			// тело ответа
//			expected: `Parsed query-param with key "name": Вано`,
//		},
//	}
//}
