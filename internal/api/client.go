package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const apiBase = "https://api.instantly.ai/api/v2"

// Client is the Instantly API client.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new API client with the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, &InstantlyError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
		}
	}
	return body, nil
}

func (c *Client) buildURL(path string, params url.Values) string {
	u, _ := url.Parse(apiBase + path)
	if params != nil {
		u.RawQuery = params.Encode()
	}
	return u.String()
}

func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.buildURL(path, params), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *Client) Post(path string, params url.Values, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encoding request: %w", err)
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPost, c.buildURL(path, params), body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.doRequest(req)
}

func (c *Client) Patch(path string, params url.Values, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("encoding request: %w", err)
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPatch, c.buildURL(path, params), body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.doRequest(req)
}

func (c *Client) Delete(path string, params url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, c.buildURL(path, params), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

// ===== Campaign =====

func (c *Client) ListCampaigns(params url.Values) ([]Campaign, string, error) {
	body, err := c.Get("/campaigns", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[Campaign]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetCampaign(id string) (*Campaign, error) {
	body, err := c.Get("/campaigns/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item Campaign
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateCampaign(payload map[string]interface{}) (*Campaign, error) {
	body, err := c.Post("/campaigns", nil, payload)
	if err != nil {
		return nil, err
	}
	var item Campaign
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateCampaign(id string, payload map[string]interface{}) (*Campaign, error) {
	body, err := c.Patch("/campaigns/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item Campaign
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteCampaign(id string) error {
	_, err := c.Delete("/campaigns/"+id, nil)
	return err
}

func (c *Client) ActivateCampaign(id string) error {
	_, err := c.Post("/campaigns/"+id+"/activate", nil, map[string]interface{}{})
	return err
}

func (c *Client) PauseCampaign(id string) error {
	_, err := c.Post("/campaigns/"+id+"/pause", nil, map[string]interface{}{})
	return err
}

func (c *Client) GetCampaignAnalytics(id string, params url.Values) (*CampaignAnalytics, error) {
	body, err := c.Get("/campaigns/"+id+"/analytics", params)
	if err != nil {
		return nil, err
	}
	var item CampaignAnalytics
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) GetCampaignAnalyticsOverview(params url.Values) ([]CampaignAnalytics, error) {
	body, err := c.Get("/campaigns/analytics/overview", params)
	if err != nil {
		return nil, err
	}
	var resp struct {
		Data []CampaignAnalytics `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		var items []CampaignAnalytics
		return items, json.Unmarshal(body, &items)
	}
	return resp.Data, nil
}

func (c *Client) DuplicateCampaign(id string) (*Campaign, error) {
	body, err := c.Post("/campaigns/"+id+"/duplicate", nil, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var item Campaign
	return &item, json.Unmarshal(body, &item)
}

// ===== Campaign Subsequence =====

func (c *Client) ListCampaignSubsequences(params url.Values) ([]CampaignSubsequence, string, error) {
	body, err := c.Get("/campaignsubsequences", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[CampaignSubsequence]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetCampaignSubsequence(id string) (*CampaignSubsequence, error) {
	body, err := c.Get("/campaignsubsequences/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item CampaignSubsequence
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateCampaignSubsequence(payload map[string]interface{}) (*CampaignSubsequence, error) {
	body, err := c.Post("/campaignsubsequences", nil, payload)
	if err != nil {
		return nil, err
	}
	var item CampaignSubsequence
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateCampaignSubsequence(id string, payload map[string]interface{}) (*CampaignSubsequence, error) {
	body, err := c.Patch("/campaignsubsequences/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item CampaignSubsequence
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteCampaignSubsequence(id string) error {
	_, err := c.Delete("/campaignsubsequences/"+id, nil)
	return err
}

func (c *Client) PauseCampaignSubsequence(id string) error {
	_, err := c.Post("/campaignsubsequences/"+id+"/pause", nil, map[string]interface{}{})
	return err
}

func (c *Client) ResumeCampaignSubsequence(id string) error {
	_, err := c.Post("/campaignsubsequences/"+id+"/resume", nil, map[string]interface{}{})
	return err
}

func (c *Client) DuplicateCampaignSubsequence(id string) (*CampaignSubsequence, error) {
	body, err := c.Post("/campaignsubsequences/"+id+"/duplicate", nil, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	var item CampaignSubsequence
	return &item, json.Unmarshal(body, &item)
}

// ===== Account =====

func (c *Client) ListAccounts(params url.Values) ([]Account, string, error) {
	body, err := c.Get("/accounts", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[Account]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetAccount(id string) (*Account, error) {
	body, err := c.Get("/accounts/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item Account
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateAccount(payload map[string]interface{}) (*Account, error) {
	body, err := c.Post("/accounts", nil, payload)
	if err != nil {
		return nil, err
	}
	var item Account
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateAccount(id string, payload map[string]interface{}) (*Account, error) {
	body, err := c.Patch("/accounts/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item Account
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteAccount(id string) error {
	_, err := c.Delete("/accounts/"+id, nil)
	return err
}

func (c *Client) EnableWarmup(emails []string) error {
	_, err := c.Post("/accounts/warmup/enable", nil, map[string]interface{}{"emails": emails})
	return err
}

func (c *Client) DisableWarmup(emails []string) error {
	_, err := c.Post("/accounts/warmup/disable", nil, map[string]interface{}{"emails": emails})
	return err
}

func (c *Client) GetWarmupAnalytics(params url.Values) ([]WarmupAnalytics, error) {
	body, err := c.Get("/accounts/warmup/analytics", params)
	if err != nil {
		return nil, err
	}
	var resp ListResponse[WarmupAnalytics]
	if err := json.Unmarshal(body, &resp); err != nil {
		var items []WarmupAnalytics
		return items, json.Unmarshal(body, &items)
	}
	return resp.Items, nil
}

func (c *Client) PauseAccount(email string) error {
	_, err := c.Post("/accounts/pause", nil, map[string]interface{}{"email": email})
	return err
}

func (c *Client) ResumeAccount(email string) error {
	_, err := c.Post("/accounts/resume", nil, map[string]interface{}{"email": email})
	return err
}

// ===== Lead =====

func (c *Client) ListLeads(params url.Values) ([]Lead, string, error) {
	body, err := c.Get("/leads", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[Lead]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetLead(id string) (*Lead, error) {
	body, err := c.Get("/leads/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item Lead
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateLead(payload map[string]interface{}) (*Lead, error) {
	body, err := c.Post("/leads", nil, payload)
	if err != nil {
		return nil, err
	}
	var item Lead
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateLead(id string, payload map[string]interface{}) (*Lead, error) {
	body, err := c.Patch("/leads/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item Lead
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteLead(id string) error {
	_, err := c.Delete("/leads/"+id, nil)
	return err
}

func (c *Client) UpdateLeadInterest(id string, status string) (*Lead, error) {
	body, err := c.Patch("/leads/"+id+"/interest", nil, map[string]interface{}{"interest_status": status})
	if err != nil {
		return nil, err
	}
	var item Lead
	return &item, json.Unmarshal(body, &item)
}

// ===== Lead List =====

func (c *Client) ListLeadLists(params url.Values) ([]LeadList, string, error) {
	body, err := c.Get("/leadlists", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[LeadList]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetLeadList(id string) (*LeadList, error) {
	body, err := c.Get("/leadlists/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item LeadList
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateLeadList(payload map[string]interface{}) (*LeadList, error) {
	body, err := c.Post("/leadlists", nil, payload)
	if err != nil {
		return nil, err
	}
	var item LeadList
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateLeadList(id string, payload map[string]interface{}) (*LeadList, error) {
	body, err := c.Patch("/leadlists/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item LeadList
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteLeadList(id string) error {
	_, err := c.Delete("/leadlists/"+id, nil)
	return err
}

// ===== Lead Label =====

func (c *Client) ListLeadLabels(params url.Values) ([]LeadLabel, string, error) {
	body, err := c.Get("/leadlabels", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[LeadLabel]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetLeadLabel(id string) (*LeadLabel, error) {
	body, err := c.Get("/leadlabels/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item LeadLabel
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateLeadLabel(payload map[string]interface{}) (*LeadLabel, error) {
	body, err := c.Post("/leadlabels", nil, payload)
	if err != nil {
		return nil, err
	}
	var item LeadLabel
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateLeadLabel(id string, payload map[string]interface{}) (*LeadLabel, error) {
	body, err := c.Patch("/leadlabels/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item LeadLabel
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteLeadLabel(id string) error {
	_, err := c.Delete("/leadlabels/"+id, nil)
	return err
}

// ===== Email =====

func (c *Client) ListEmails(params url.Values) ([]Email, string, error) {
	body, err := c.Get("/emails", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[Email]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetEmail(id string) (*Email, error) {
	body, err := c.Get("/emails/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item Email
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) ReplyToEmail(payload map[string]interface{}) error {
	_, err := c.Post("/emails/reply", nil, payload)
	return err
}

func (c *Client) ForwardEmail(payload map[string]interface{}) error {
	_, err := c.Post("/emails/forward", nil, payload)
	return err
}

func (c *Client) MarkThreadAsRead(threadID string) error {
	_, err := c.Post("/emails/thread/read", nil, map[string]interface{}{"thread_id": threadID})
	return err
}

// ===== Email Verification =====

func (c *Client) CreateEmailVerification(email string) (*EmailVerification, error) {
	body, err := c.Post("/emailverifications", nil, map[string]interface{}{"email": email})
	if err != nil {
		return nil, err
	}
	var item EmailVerification
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) GetEmailVerification(id string) (*EmailVerification, error) {
	body, err := c.Get("/emailverifications/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item EmailVerification
	return &item, json.Unmarshal(body, &item)
}

// ===== Webhook =====

func (c *Client) ListWebhooks(params url.Values) ([]Webhook, string, error) {
	body, err := c.Get("/webhooks", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[Webhook]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetWebhook(id string) (*Webhook, error) {
	body, err := c.Get("/webhooks/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item Webhook
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateWebhook(payload map[string]interface{}) (*Webhook, error) {
	body, err := c.Post("/webhooks", nil, payload)
	if err != nil {
		return nil, err
	}
	var item Webhook
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateWebhook(id string, payload map[string]interface{}) (*Webhook, error) {
	body, err := c.Patch("/webhooks/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item Webhook
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteWebhook(id string) error {
	_, err := c.Delete("/webhooks/"+id, nil)
	return err
}

func (c *Client) TestWebhook(id string) error {
	_, err := c.Post("/webhooks/"+id+"/test", nil, map[string]interface{}{})
	return err
}

func (c *Client) ResumeWebhook(id string) error {
	_, err := c.Post("/webhooks/"+id+"/resume", nil, map[string]interface{}{})
	return err
}

func (c *Client) ListWebhookEventTypes() ([]WebhookEventType, error) {
	body, err := c.Get("/webhooks/event-types", nil)
	if err != nil {
		return nil, err
	}
	var items []WebhookEventType
	if err := json.Unmarshal(body, &items); err != nil {
		var resp struct {
			EventTypes []WebhookEventType `json:"event_types"`
		}
		if err2 := json.Unmarshal(body, &resp); err2 != nil {
			return nil, err
		}
		return resp.EventTypes, nil
	}
	return items, nil
}

// ===== Custom Tag =====

func (c *Client) ListCustomTags(params url.Values) ([]CustomTag, string, error) {
	body, err := c.Get("/customtags", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[CustomTag]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetCustomTag(id string) (*CustomTag, error) {
	body, err := c.Get("/customtags/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item CustomTag
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateCustomTag(payload map[string]interface{}) (*CustomTag, error) {
	body, err := c.Post("/customtags", nil, payload)
	if err != nil {
		return nil, err
	}
	var item CustomTag
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateCustomTag(id string, payload map[string]interface{}) (*CustomTag, error) {
	body, err := c.Patch("/customtags/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item CustomTag
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteCustomTag(id string) error {
	_, err := c.Delete("/customtags/"+id, nil)
	return err
}

func (c *Client) ToggleCustomTag(tagID, resourceID, resourceType string) error {
	_, err := c.Post("/customtags/toggle", nil, map[string]interface{}{
		"tag_id":        tagID,
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})
	return err
}

// ===== Blocklist Entry =====

func (c *Client) ListBlocklistEntries(params url.Values) ([]BlocklistEntry, string, error) {
	body, err := c.Get("/blocklistentries", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[BlocklistEntry]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetBlocklistEntry(id string) (*BlocklistEntry, error) {
	body, err := c.Get("/blocklistentries/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item BlocklistEntry
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateBlocklistEntry(payload map[string]interface{}) (*BlocklistEntry, error) {
	body, err := c.Post("/blocklistentries", nil, payload)
	if err != nil {
		return nil, err
	}
	var item BlocklistEntry
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateBlocklistEntry(id string, payload map[string]interface{}) (*BlocklistEntry, error) {
	body, err := c.Patch("/blocklistentries/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item BlocklistEntry
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteBlocklistEntry(id string) error {
	_, err := c.Delete("/blocklistentries/"+id, nil)
	return err
}

// ===== API Key =====

func (c *Client) ListAPIKeys(params url.Values) ([]APIKey, string, error) {
	body, err := c.Get("/apikeys", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[APIKey]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) CreateAPIKey(name string) (*APIKey, error) {
	body, err := c.Post("/apikeys", nil, map[string]interface{}{"name": name})
	if err != nil {
		return nil, err
	}
	var item APIKey
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteAPIKey(id string) error {
	_, err := c.Delete("/apikeys/"+id, nil)
	return err
}

// ===== Workspace =====

func (c *Client) GetWorkspace() (*Workspace, error) {
	body, err := c.Get("/workspaces", nil)
	if err != nil {
		return nil, err
	}
	var item Workspace
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateWorkspace(payload map[string]interface{}) (*Workspace, error) {
	body, err := c.Patch("/workspaces", nil, payload)
	if err != nil {
		return nil, err
	}
	var item Workspace
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) ListWorkspaceMembers(params url.Values) ([]WorkspaceMember, string, error) {
	body, err := c.Get("/workspacemembers", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[WorkspaceMember]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetWorkspaceMember(id string) (*WorkspaceMember, error) {
	body, err := c.Get("/workspacemembers/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item WorkspaceMember
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) CreateWorkspaceMember(payload map[string]interface{}) (*WorkspaceMember, error) {
	body, err := c.Post("/workspacemembers", nil, payload)
	if err != nil {
		return nil, err
	}
	var item WorkspaceMember
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) UpdateWorkspaceMember(id string, payload map[string]interface{}) (*WorkspaceMember, error) {
	body, err := c.Patch("/workspacemembers/"+id, nil, payload)
	if err != nil {
		return nil, err
	}
	var item WorkspaceMember
	return &item, json.Unmarshal(body, &item)
}

func (c *Client) DeleteWorkspaceMember(id string) error {
	_, err := c.Delete("/workspacemembers/"+id, nil)
	return err
}

// ===== Background Job =====

func (c *Client) ListBackgroundJobs(params url.Values) ([]BackgroundJob, string, error) {
	body, err := c.Get("/backgroundjobs", params)
	if err != nil {
		return nil, "", err
	}
	var resp ListResponse[BackgroundJob]
	return resp.Items, resp.NextStartingAfter, json.Unmarshal(body, &resp)
}

func (c *Client) GetBackgroundJob(id string) (*BackgroundJob, error) {
	body, err := c.Get("/backgroundjobs/"+id, nil)
	if err != nil {
		return nil, err
	}
	var item BackgroundJob
	return &item, json.Unmarshal(body, &item)
}
