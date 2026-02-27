package api

// InstantlyError is returned when the API responds with an error.
type InstantlyError struct {
	StatusCode int
	Message    string
}

func (e *InstantlyError) Error() string {
	return e.Message
}

// --- Generic list response ---

// ListResponse wraps paginated list responses from the Instantly API.
type ListResponse[T any] struct {
	Items             []T    `json:"items"`
	NextStartingAfter string `json:"next_starting_after"`
}

// --- Campaign ---

type Campaign struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Status          int      `json:"status"` // 0=draft, 1=active, 2=paused, 3=completed, 4=error
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	DailyLimit      int      `json:"daily_limit"`
	EmailList       []string `json:"email_list"`
	Timezone        string   `json:"timezone"`
	StopOnReply     bool     `json:"stop_on_reply"`
	StopOnAutoReply bool     `json:"stop_on_auto_reply"`
	LinkTracking    bool     `json:"link_tracking"`
	OpenTracking    bool     `json:"open_tracking"`
	TextOnly        bool     `json:"text_only"`
}

type CampaignAnalytics struct {
	CampaignID   string  `json:"campaign_id"`
	CampaignName string  `json:"campaign_name"`
	Sent         int     `json:"total_sent"`
	Opened       int     `json:"total_opened"`
	Clicked      int     `json:"total_clicked"`
	Replied      int     `json:"total_replied"`
	Bounced      int     `json:"total_bounced"`
	Unsubscribed int     `json:"total_unsubscribed"`
	NewLeads     int     `json:"new_leads_contacted"`
	OpenRate     float64 `json:"open_rate"`
	ClickRate    float64 `json:"click_rate"`
	ReplyRate    float64 `json:"reply_rate"`
}

// --- Campaign Subsequence ---

type CampaignSubsequence struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CampaignID string `json:"campaign_id"`
	Status     int    `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// --- Account ---

type Account struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Status          int    `json:"status"` // 1=active, 2=paused, -1=error
	DailyLimit      int    `json:"daily_limit"`
	WarmupEnabled   bool   `json:"warmup_enabled"`
	WarmupLimit     int    `json:"warmup_limit"`
	WarmupReplyRate int    `json:"warmup_reply_rate"`
	SmtpHost        string `json:"smtp_host"`
	SmtpPort        int    `json:"smtp_port"`
	ImapHost        string `json:"imap_host"`
	ImapPort        int    `json:"imap_port"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Alias           string `json:"alias"`
	TrackingDomain  string `json:"tracking_domain_name"`
}

type WarmupAnalytics struct {
	Email    string `json:"email"`
	Sent     int    `json:"warmup_emails_sent_count"`
	Received int    `json:"warmup_emails_received_count"`
	Date     string `json:"date"`
}

// --- Lead ---

type Lead struct {
	ID              string                 `json:"id"`
	Email           string                 `json:"email"`
	FirstName       string                 `json:"first_name"`
	LastName        string                 `json:"last_name"`
	CompanyName     string                 `json:"company_name"`
	Phone           string                 `json:"phone"`
	Website         string                 `json:"website"`
	LinkedinURL     string                 `json:"linkedin_url"`
	CustomVariables map[string]interface{} `json:"custom_variables"`
	CampaignID      string                 `json:"campaign_id"`
	ListID          string                 `json:"list_id"`
	Status          string                 `json:"lt_interest_status"`
	AssignedTo      string                 `json:"assigned_to"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// --- Lead List ---

type LeadList struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// --- Lead Label ---

type LeadLabel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// --- Email ---

type Email struct {
	ID         string `json:"id"`
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	FromEmail  string `json:"from_email"`
	ToEmail    string `json:"to_email"`
	Type       string `json:"type"`
	CampaignID string `json:"campaign_id"`
	Timestamp  string `json:"timestamp"`
	IsRead     bool   `json:"is_read"`
	ThreadID   string `json:"eaccount_reply_gmail_message_thread_id"`
}

// --- Email Verification ---

type EmailVerification struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
	Valid  bool   `json:"valid"`
}

// --- Webhook ---

type Webhook struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	Active     bool     `json:"active"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type WebhookEventType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// --- Custom Tag ---

type CustomTag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	WorkspaceID string `json:"workspace_id"`
	CreatedAt   string `json:"created_at"`
}

// --- Blocklist Entry ---

type BlocklistEntry struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

// --- API Key ---

type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
}

// --- Workspace ---

type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// --- Workspace Member ---

type WorkspaceMember struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// --- Background Job ---

type BackgroundJob struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Progress  float64 `json:"progress"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
