package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/dto"
	"log"
	"net/http"
)

type CatalogClient struct {
	BaseURL string
}

func NewCatalogClient(baseURL string) *CatalogClient {
	return &CatalogClient{
		BaseURL: baseURL,
	}
}

func (c *CatalogClient) GetProductByID(productID uint) (*dto.ProductResponse, error) {
	url := fmt.Sprintf("%s/products/%d", c.BaseURL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get product from catalog service")
	}

	var response struct {
		Data *dto.ProductResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	log.Println(response.Data)
	return response.Data, nil
}
