package ncp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ListSubAccounts retrieves all sub accounts for the authenticated root account.
func (c *Client) ListSubAccounts() ([]SubAccount, error) {
	var all []SubAccount
	page := 0
	pageSize := 100

	for {
		path := fmt.Sprintf("/api/v1/sub-accounts?pageSize=%d&page=%d", pageSize, page)
		body, statusCode, err := c.doRequest("GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("listing sub accounts: %w", err)
		}
		if statusCode != 200 {
			return nil, fmt.Errorf("listing sub accounts: HTTP %d - %s", statusCode, string(body))
		}

		var resp SubAccountListResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}

		all = append(all, resp.Content...)

		if len(all) >= resp.TotalRows || len(resp.Content) == 0 {
			break
		}
		page++
	}

	return all, nil
}

// UpdateSubAccount updates a sub account by ID.
func (c *Client) UpdateSubAccount(subAccountId int64, req *SubAccountUpdateRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/sub-accounts/%d", subAccountId)
	body, statusCode, err := c.doRequest("PUT", path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("updating sub account: %w", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("updating sub account: HTTP %d - %s", statusCode, string(body))
	}

	return nil
}

// ActivateSubAccount activates a sub account and resets its password.
func (c *Client) ActivateSubAccount(subAccountId int64, password string) error {
	active := true
	needReset := true
	req := &SubAccountUpdateRequest{
		Active:            &active,
		Password:          &password,
		NeedPasswordReset: &needReset,
	}
	return c.UpdateSubAccount(subAccountId, req)
}

// DeactivateSubAccount deactivates (suspends) a sub account.
func (c *Client) DeactivateSubAccount(subAccountId int64) error {
	active := false
	req := &SubAccountUpdateRequest{
		Active: &active,
	}
	return c.UpdateSubAccount(subAccountId, req)
}
