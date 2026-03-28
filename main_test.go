package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T){
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
        count int    // передаваемое значение count
        want  int  	 // ожидаемое количество кафе в ответе
		city  string // город поиска
    }{
        {0, min(0, len(cafeList["moscow"])), "moscow"},
		{1, min(1, len(cafeList["tula"])), "tula"},
		{2, min(2, len(cafeList["moscow"])), "moscow"},
		{100, min(100, len(cafeList["tula"])), "tula"},
    } 

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", v.count, v.city), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
						
		if v.count == 0 && response.Body.String() == "" {
			continue
		}

		foundCofe := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		
		assert.Equal(t, v.want, len(foundCofe))
	}
}

func TestCafeSearch(t *testing.T){
	handler := http.HandlerFunc(mainHandle)
	 
	requests := []struct {
        search    string   // передаваемое значение search 
        wantCount int      // ожидаемое количество кафе в ответе
    }{
        {"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
		{"lojka", 0},
		{"Студент", 1},
		{"404", 0},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?search=%s&city=moscow", v.search) , nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		if v.wantCount == 0 && response.Body.String() == "" {
			continue
		}

		foundCofe := strings.Split(strings.TrimSpace(response.Body.String()), ",")

		for _, cofe := range foundCofe {
			assert.Equal(t, true, strings.Contains(strings.ToLower(cofe), strings.ToLower(v.search)))
			assert.Equal(t, v.wantCount, len(foundCofe))
		}
	}
	
}