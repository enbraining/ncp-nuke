package ncp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	SubAccountBaseURL   = "https://subaccount.apigw.ntruss.com"
	VServerBaseURL      = "https://ncloud.apigw.ntruss.com/vserver/v2"
	VNASBaseURL         = "https://ncloud.apigw.ntruss.com/vnas/v2"
	VLBBaseURL          = "https://ncloud.apigw.ntruss.com/vloadbalancer/v2"
	VCloudDBBaseURL     = "https://ncloud.apigw.ntruss.com/clouddb/v2"
	VVPCBaseURL         = "https://ncloud.apigw.ntruss.com/vpc/v2"
	VNKSBaseURL         = "https://nks.apigw.ntruss.com/vnks/v2"
	VAutoScalingBaseURL = "https://ncloud.apigw.ntruss.com/autoscaling/v2"
	VMongoDBBaseURL     = "https://ncloud.apigw.ntruss.com/vmongodb/v2"
	VPostgreSQLBaseURL  = "https://ncloud.apigw.ntruss.com/vpostgresql/v2"
	VMariaDBBaseURL     = "https://ncloud.apigw.ntruss.com/vmariadb/v2"
	VMySQLBaseURL       = "https://ncloud.apigw.ntruss.com/vmysql/v2"
	VRedisBaseURL       = "https://ncloud.apigw.ntruss.com/vredis/v2"
)

// Client is the NCP API client with HMAC-SHA256 authentication.
type Client struct {
	accessKey  string
	secretKey  string
	httpClient *http.Client
}

// NewClient creates a new NCP API client.
func NewClient(accessKey, secretKey string) *Client {
	return &Client{
		accessKey: accessKey,
		secretKey: secretKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeSignature generates the HMAC-SHA256 signature for NCP API authentication.
// Format: {method} {url}\n{timestamp}\n{accessKey}
func (c *Client) makeSignature(method, url, timestamp string) string {
	message := fmt.Sprintf("%s %s\n%s\n%s", method, url, timestamp, c.accessKey)
	mac := hmac.New(sha256.New, []byte(c.secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// doRequest executes an HTTP request with NCP authentication headers.
// Uses SubAccountBaseURL by default.
func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, int, error) {
	return c.doRequestWithBase(SubAccountBaseURL, method, path, body)
}

// doRequestWithBase executes an HTTP request against a specific base URL.
func (c *Client) doRequestWithBase(baseURL, method, path string, body io.Reader) ([]byte, int, error) {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// For signature, we need the full path including the base path portion
	fullURL := baseURL + path
	// Extract the path portion from the full URL for signature
	signPath := extractPath(fullURL)
	signature := c.makeSignature(method, signPath, timestamp)

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("x-ncp-apigw-timestamp", timestamp)
	req.Header.Set("x-ncp-iam-access-key", c.accessKey)
	req.Header.Set("x-ncp-apigw-signature-v2", signature)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// extractPath extracts the path (including query string) from a full URL.
func extractPath(fullURL string) string {
	// Find the third slash (after https://)
	count := 0
	for i, ch := range fullURL {
		if ch == '/' {
			count++
			if count == 3 {
				return fullURL[i:]
			}
		}
	}
	return "/"
}
