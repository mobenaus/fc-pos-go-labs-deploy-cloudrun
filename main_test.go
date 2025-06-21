package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func httpResponse(response string) func(url string) (*http.Response, error) {
	return func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(response)),
		}, nil
	}
}

func httpError() func(url string) (*http.Response, error) {
	return func(url string) (*http.Response, error) {
		return nil, errors.New("request error")
	}
}

type ApiClientMock struct {
	mock.Mock
}

func (c *ApiClientMock) getCityByCEP(cep string) (string, error) {
	args := c.Called(cep)
	return args.String(0), args.Error(1)
}

func (c *ApiClientMock) getTemperatureByCity(cep string) (float64, error) {
	args := c.Called(cep)
	return float64(args.Int(0)), args.Error(1)
}

func TestIsValidCEP(t *testing.T) {
	assert.True(t, isValidCEP("12345678"))
	assert.True(t, isValidCEP("87654321"))
	assert.False(t, isValidCEP("1234567"))
	assert.False(t, isValidCEP("abcdefgh"))
	assert.False(t, isValidCEP("1234-567"))
	assert.False(t, isValidCEP(""))
}

func TestGetCityByCEP(t *testing.T) {
	// mock de respota correta
	httpGet := httpResponse(`{"localidade":"São Paulo"}`)
	city, err := NewClient(httpGet, "").getCityByCEP("01001000")
	assert.NoError(t, err)
	assert.Equal(t, "São Paulo", city)

	// Mock de resposta com erro

	httpGet = httpResponse(`{"erro":true}`)
	city, err = NewClient(httpGet, "").getCityByCEP("00000000")
	assert.Error(t, err)
	assert.Empty(t, city)

	// Mock de resposta invalida
	httpGet = httpResponse(`not a valid JSON`)
	city, err = NewClient(httpGet, "").getCityByCEP("00000000")
	assert.Error(t, err)
	assert.Empty(t, city)

	// Mock de erro de requisição
	httpGet = httpError()
	city, err = NewClient(httpGet, "").getCityByCEP("00000000")
	assert.Error(t, err)
	assert.Empty(t, city)
}

func TestGetTemperatureByCity(t *testing.T) {
	// Caso de sucesso
	httpGet := httpResponse(`{"current":{"temp_c":22.5}}`)
	temp, err := NewClient(httpGet, "").getTemperatureByCity("São Paulo")
	assert.NoError(t, err)
	assert.Equal(t, 22.5, temp)

	// Caso de erro na resposta JSON
	httpGet = httpResponse(`not a valid JSON`)
	_, err = NewClient(httpGet, "").getTemperatureByCity("São Paulo")
	assert.Error(t, err)

	// Caso de erro de requisição
	httpGet = httpError()
	_, err = NewClient(httpGet, "").getTemperatureByCity("São Paulo")
	assert.Error(t, err)
}

func TestWeatherHandlerInvalidZipCode(t *testing.T) {
	client := &ApiClientMock{}
	wh := NewWeatherHandler(client)
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=123", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wh.weatherHandler)
	handler.ServeHTTP(rr, req)
	assert.Exactly(t, 422, rr.Code)
	assert.Equal(t, "invalid zipcode", strings.TrimSpace(rr.Body.String()))
}

func TestWeatherHandlerCannotFindZipCode(t *testing.T) {
	client := &ApiClientMock{}
	client.On("getCityByCEP", "12345678").Return("234", errors.New("not found"))

	wh := NewWeatherHandler(client)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=12345678", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wh.weatherHandler)
	handler.ServeHTTP(rr, req)
	assert.Exactly(t, 404, rr.Code)
	assert.Equal(t, "can not find zipcode", strings.TrimSpace(rr.Body.String()))
}
func TestWeatherHandlerCannotFindCityTemp(t *testing.T) {
	client := &ApiClientMock{}
	client.On("getCityByCEP", "12345678").Return("Cidade Fake", nil)
	client.On("getTemperatureByCity", "Cidade Fake").Return(0, errors.New("not found"))

	wh := NewWeatherHandler(client)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=12345678", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wh.weatherHandler)
	handler.ServeHTTP(rr, req)
	assert.Exactly(t, 404, rr.Code)
	assert.Equal(t, "can not find temperature", strings.TrimSpace(rr.Body.String()))
}

func TestWeatherHandlerWithSuccess(t *testing.T) {
	client := &ApiClientMock{}
	client.On("getCityByCEP", "12345678").Return("Cidade Fake", nil)
	client.On("getTemperatureByCity", "Cidade Fake").Return(30, nil)

	wh := NewWeatherHandler(client)

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=12345678", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(wh.weatherHandler)
	handler.ServeHTTP(rr, req)
	assert.Exactly(t, 200, rr.Code)
	assert.JSONEq(t, `{"temp_C":30, "temp_F":86, "temp_K":303}`, rr.Body.String())
}
