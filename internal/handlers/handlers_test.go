package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophermart/internal/middleware"
	"gophermart/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	defaultSignatureKey = "super_secret"
	defaultPassKey      = "myverystrongpasswordo32bitlength"
	cookieName          = "authToken"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body interface{}, userID string) (*http.Response, string) {
	var reqBody io.Reader = nil

	switch v := body.(type) {
	case string:
		reqBody = strings.NewReader(v)
	default:
		jsonData, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	require.NoError(t, err)

	if userID != "" {
		authToken, err := middleware.MakeAuthToken(defaultSignatureKey, userID)
		require.NoError(t, err)

		cookie := http.Cookie{Name: cookieName, Value: authToken}
		req.AddCookie(&cookie)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := NewMockgmService(ctrl)

	handler, err := New(mockService, defaultSignatureKey, []byte(defaultPassKey))
	require.NoError(t, err)

	ts := httptest.NewServer(handler.InitRouter())
	defer ts.Close()

	oneOrder := []model.Order{{
		Number:     "1234",
		Status:     "NEW",
		UploadedAt: time.Now(),
	}}

	oneOrderByte, _ := json.Marshal(oneOrder)

	balance := model.Balance{
		Current:   500.5,
		Withdrawn: 42,
	}

	balanceByte, _ := json.Marshal(balance)

	type want struct {
		statusCode  int
		contentType string
		respBody    string
	}

	tests := []struct {
		name        string
		method      string
		path        string
		body        interface{}
		userForAuth string
		expectCall  func()
		want        want
	}{
		{
			name:   "simple register",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: model.LogoPass{
				Login:    "login1",
				Password: "pass1",
			},
			expectCall: func() {
				mockService.EXPECT().AddAuthInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(int64(1), nil)
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name:   "register with login exist",
			method: http.MethodPost,
			path:   "/api/user/register",
			body: model.LogoPass{
				Login:    "login2",
				Password: "pass2",
			},
			expectCall: func() {
				mockService.EXPECT().AddAuthInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(int64(0), model.ErrLoginAlreadyExist)
			},
			want: want{
				statusCode:  http.StatusConflict,
				contentType: "text/plain; charset=utf-8",
				respBody:    "login already exist\n",
			},
		},
		{
			name:   "simple login",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: model.LogoPass{
				Login:    "login1",
				Password: "pass1",
			},
			expectCall: func() {
				mockService.EXPECT().GetAuthInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(int64(1), nil)
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name:   "login wrong login",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: model.LogoPass{
				Login:    "login3",
				Password: "pass3",
			},
			expectCall: func() {
				mockService.EXPECT().GetAuthInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(int64(0), model.ErrWrongLogin)
			},
			want: want{
				statusCode:  http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
				respBody:    "login does not exist\n",
			},
		},
		{
			name:   "login wrong password",
			method: http.MethodPost,
			path:   "/api/user/login",
			body: model.LogoPass{
				Login:    "login1",
				Password: "pass2",
			},
			expectCall: func() {
				mockService.EXPECT().GetAuthInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(int64(0), model.ErrWrongPas)
			},
			want: want{
				statusCode:  http.StatusUnauthorized,
				contentType: "text/plain; charset=utf-8",
				respBody:    "wrong password\n",
			},
		},
		{
			name:        "simple add order",
			method:      http.MethodPost,
			path:        "/api/user/orders",
			body:        "79927398713",
			userForAuth: "4",
			expectCall: func() {
				mockService.EXPECT().AddOrder(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			want: want{
				statusCode:  http.StatusAccepted,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "add order incorrect order",
			method:      http.MethodPost,
			path:        "/api/user/orders",
			body:        "79927398712",
			userForAuth: "4",
			expectCall: func() {
			},
			want: want{
				statusCode:  http.StatusUnprocessableEntity,
				contentType: "text/plain; charset=utf-8",
				respBody:    "incorrect order id\n",
			},
		},
		{
			name:        "get orders simple",
			method:      http.MethodGet,
			path:        "/api/user/orders",
			body:        nil,
			userForAuth: "4",
			expectCall: func() {
				mockService.EXPECT().GetOrders(gomock.Any(), gomock.Any()).Times(1).Return(oneOrder, nil)
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				respBody:    string(oneOrderByte),
			},
		},
		{
			name:        "get balance simple",
			method:      http.MethodGet,
			path:        "/api/user/balance",
			body:        nil,
			userForAuth: "4",
			expectCall: func() {
				mockService.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Times(1).Return(balance, nil)
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				respBody:    string(balanceByte),
			},
		},
	}

	for _, tt := range tests {
		fmt.Println("Test: ", tt.name)
		tt.expectCall()

		resp, get := testRequest(t, ts, tt.method, tt.path, tt.body, tt.userForAuth)

		assert.Equal(t, tt.want.respBody, get)

		resp.Body.Close()

		assert.Equal(t, tt.want.statusCode, resp.StatusCode)
		assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
	}
}
