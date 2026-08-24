package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	domainfipe "cvmc/internal/domain/fipe"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = "https://fipe.parallelum.com.br/api/v2"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		baseURL: baseURL,
		token:   strings.TrimSpace(token),
		httpClient: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.token != "" {
		req.Header.Set("X-Subscription-Token", c.token)
	}
	return req, nil
}

func (c *Client) FetchBrands(ctx context.Context, vehicleType domainfipe.VehicleType) ([]domainfipe.Brand, error) {
	endpoint := fmt.Sprintf("/%s/brands", vehicleType)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching brands from fipe api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fipe api returned status %d: %s", resp.StatusCode, string(body))
	}

	var brands []domainfipe.Brand
	if err := json.NewDecoder(resp.Body).Decode(&brands); err != nil {
		return nil, fmt.Errorf("failed to decode brands response: %w", err)
	}

	for i := range brands {
		brands[i].VehicleType = string(vehicleType)
	}

	return brands, nil
}

func (c *Client) FetchModels(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode string) ([]domainfipe.Model, error) {
	endpoint := fmt.Sprintf("/%s/brands/%s/models", vehicleType, brandCode)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching models from fipe api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fipe api returned status %d: %s", resp.StatusCode, string(body))
	}

	var models []domainfipe.Model
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}
	return models, nil
}

func (c *Client) FetchYears(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode string) ([]domainfipe.Year, error) {
	endpoint := fmt.Sprintf("/%s/brands/%s/models/%s/years", vehicleType, brandCode, modelCode)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching years from fipe api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fipe api returned status %d: %s", resp.StatusCode, string(body))
	}

	var years []domainfipe.Year
	if err := json.NewDecoder(resp.Body).Decode(&years); err != nil {
		return nil, fmt.Errorf("failed to decode years response: %w", err)
	}
	return years, nil
}

func (c *Client) FetchVehicleDetail(ctx context.Context, vehicleType domainfipe.VehicleType, brandCode, modelCode, yearCode string) (domainfipe.VehicleDetail, error) {
	endpoint := fmt.Sprintf("/%s/brands/%s/models/%s/years/%s", vehicleType, brandCode, modelCode, yearCode)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return domainfipe.VehicleDetail{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domainfipe.VehicleDetail{}, fmt.Errorf("error fetching vehicle detail from fipe api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return domainfipe.VehicleDetail{}, fmt.Errorf("fipe api returned status %d: %s", resp.StatusCode, string(body))
	}

	var detail domainfipe.VehicleDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return domainfipe.VehicleDetail{}, fmt.Errorf("failed to decode vehicle detail response: %w", err)
	}

	detail.PriceValue = parsePriceValue(detail.Price)
	return detail, nil
}

func parsePriceValue(raw string) float64 {
	// Ex: "R$ 55.400,00" -> 55400.00
	cleaned := strings.ReplaceAll(raw, "R$", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	val, err := strconv.ParseFloat(strings.TrimSpace(cleaned), 64)
	if err != nil {
		return 0
	}
	return val
}
