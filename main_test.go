package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testerror struct {
	Error string `json:"error"`
}

func TestAmazonScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/amazon", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var amazonError testerror
		json.NewDecoder(resp.Body).Decode(&amazonError)
		t.Errorf("/amazon failed. Error message: %s", amazonError.Error)
	}
}

func TestAtlassianScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/atlassian", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var atlassianError testerror
		json.NewDecoder(resp.Body).Decode(&atlassianError)
		t.Errorf("/atlassian failed. Error message: %s", atlassianError.Error)
	}
}

func TestCourseraScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/coursera", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var courseraError testerror
		json.NewDecoder(resp.Body).Decode(&courseraError)
		t.Errorf("/coursera failed. Error message: %s", courseraError.Error)
	}
}

func TestFreshWorksScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/freshworks", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var freshworksError testerror
		json.NewDecoder(resp.Body).Decode(&freshworksError)
		t.Errorf("/freshworks failed. Error message: %s", freshworksError.Error)
	}
}

func TestGojekScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/gojek", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var gojekError testerror
		json.NewDecoder(resp.Body).Decode(&gojekError)
		t.Errorf("/gojek failed. Error message: %s", gojekError.Error)
	}
}

func TestMplScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/mpl", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var mplError testerror
		json.NewDecoder(resp.Body).Decode(&mplError)
		t.Errorf("/mpl failed. Error message: %s", mplError.Error)
	}
}

func TestGoogleScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/google", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var googleError testerror
		json.NewDecoder(resp.Body).Decode(&googleError)
		t.Errorf("/google failed. Error message: %s", googleError.Error)
	}
}

func TestFiScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/fi", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var fiError testerror
		json.NewDecoder(resp.Body).Decode(&fiError)
		t.Errorf("/fi failed. Error message: %s", fiError.Error)
	}
}

func TestFrontrowScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/frontrow", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var frontrowError testerror
		json.NewDecoder(resp.Body).Decode(&frontrowError)
		t.Errorf("/frontrow failed. Error message: %s", frontrowError.Error)
	}
}

func TestSardineScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/sardine", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var sardineError testerror
		json.NewDecoder(resp.Body).Decode(&sardineError)
		t.Errorf("/sardine failed. Error message: %s", sardineError.Error)
	}
}

func TestZohoScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/zoho", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var zohoError testerror
		json.NewDecoder(resp.Body).Decode(&zohoError)
		t.Errorf("/zoho failed. Error message: %s", zohoError.Error)
	}
}

func TestJarScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/jar", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var jarError testerror
		json.NewDecoder(resp.Body).Decode(&jarError)
		t.Errorf("/jar failed. Error message: %s", jarError.Error)
	}
}

func TestPaytmScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/paytm", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var paytmError testerror
		json.NewDecoder(resp.Body).Decode(&paytmError)
		t.Errorf("/paytm failed. Error message: %s", paytmError.Error)
	}
}

func TestFincentScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/fincent", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var fincentError testerror
		json.NewDecoder(resp.Body).Decode(&fincentError)
		t.Errorf("/fincent failed. Error message: %s", fincentError.Error)
	}
}

func TestPaypalScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/paypal", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var paypalError testerror
		json.NewDecoder(resp.Body).Decode(&paypalError)
		t.Errorf("/paypal failed. Error message: %s", paypalError.Error)
	}
}

func TestNiyoScrapersHandler(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/niyo", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, http.StatusOK, resp.Code); check == false {
		var niyoError testerror
		json.NewDecoder(resp.Body).Decode(&niyoError)
		t.Errorf("/niyo failed. Error message: %s", niyoError.Error)
	}
}

func TestCredScraper(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/cred", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if check := assert.Equal(t, resp.Code, http.StatusOK); check == false {
		var credError testerror
		json.NewDecoder(resp.Body).Decode(&credError)
		t.Errorf("/cred failed. Error message: %s", credError.Error)
	}
}
