package accrual

import (
	"encoding/json"
	"fmt"
	"gophermart/internal/model"
	"io/ioutil"
	"net/http"
	"time"
)

const getAccrualPath = "/api/orders/%s"

type client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *client {
	return &client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *client) GetAccrual(orderID string) (int, model.Accrual, error) {
	path := fmt.Sprintf(getAccrualPath, orderID)
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return 0, model.Accrual{}, fmt.Errorf("GetAccrual Get-err: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, model.Accrual{}, fmt.Errorf("GetAccrual ReadBody-err: %w", err)
	}

	var accrual model.Accrual
	err = json.Unmarshal(body, &accrual)
	if err != nil {
		return 0, model.Accrual{}, fmt.Errorf("GetAccrual UnmarshalBody-err: %w", err)
	}

	return resp.StatusCode, accrual, nil
}
