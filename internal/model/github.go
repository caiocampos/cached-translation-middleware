package model

type ListUserReposRequest struct {
	Username string `json:"username" validate:"required"`

	Type      string `json:"type,omitempty"`      // all | owner | member — default: owner
	Sort      string `json:"sort,omitempty"`      // created | updated | pushed | full_name — default: full_name
	Direction string `json:"direction,omitempty"` // asc | desc — default: asc
	PerPage   int    `json:"per_page,omitempty"`  // máx 100 — default: 30
	Page      int    `json:"page,omitempty"`      // default: 1
}

type RepoOwner struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type RepoLicense struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	SPDXID  string `json:"spdx_id"`
	URL     string `json:"url"`
	NodeID  string `json:"node_id"`
	HTMLURL string `json:"html_url,omitempty"`
}

type RepoPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type SecurityFeatureStatus struct {
	Status string `json:"status"` // "enabled" | "disabled"
}

type RepoSecurityAndAnalysis struct {
	AdvancedSecurity                      *SecurityFeatureStatus `json:"advanced_security,omitempty"`
	SecretScanning                        *SecurityFeatureStatus `json:"secret_scanning,omitempty"`
	SecretScanningPushProtection          *SecurityFeatureStatus `json:"secret_scanning_push_protection,omitempty"`
	SecretScanningNonProviderPatterns     *SecurityFeatureStatus `json:"secret_scanning_non_provider_patterns,omitempty"`
	SecretScanningDelegatedAlertDismissal *SecurityFeatureStatus `json:"secret_scanning_delegated_alert_dismissal,omitempty"`
}

type RepoItem struct {
	ID          int64      `json:"id"`
	NodeID      string     `json:"node_id"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Owner       *RepoOwner `json:"owner"`
	Private     bool       `json:"private"`
	HTMLURL     string     `json:"html_url"`
	Description *string    `json:"description"`
	Fork        bool       `json:"fork"`
	URL         string     `json:"url"`

	ArchiveURL       string  `json:"archive_url"`
	AssigneesURL     string  `json:"assignees_url"`
	BlobsURL         string  `json:"blobs_url"`
	BranchesURL      string  `json:"branches_url"`
	CollaboratorsURL string  `json:"collaborators_url"`
	CommentsURL      string  `json:"comments_url"`
	CommitsURL       string  `json:"commits_url"`
	CompareURL       string  `json:"compare_url"`
	ContentsURL      string  `json:"contents_url"`
	ContributorsURL  string  `json:"contributors_url"`
	DeploymentsURL   string  `json:"deployments_url"`
	DownloadsURL     string  `json:"downloads_url"`
	EventsURL        string  `json:"events_url"`
	ForksURL         string  `json:"forks_url"`
	GitCommitsURL    string  `json:"git_commits_url"`
	GitRefsURL       string  `json:"git_refs_url"`
	GitTagsURL       string  `json:"git_tags_url"`
	GitURL           string  `json:"git_url"`
	IssueCommentURL  string  `json:"issue_comment_url"`
	IssueEventsURL   string  `json:"issue_events_url"`
	IssuesURL        string  `json:"issues_url"`
	KeysURL          string  `json:"keys_url"`
	LabelsURL        string  `json:"labels_url"`
	LanguagesURL     string  `json:"languages_url"`
	MergesURL        string  `json:"merges_url"`
	MilestonesURL    string  `json:"milestones_url"`
	NotificationsURL string  `json:"notifications_url"`
	PullsURL         string  `json:"pulls_url"`
	ReleasesURL      string  `json:"releases_url"`
	SSHURL           string  `json:"ssh_url"`
	StargazersURL    string  `json:"stargazers_url"`
	StatusesURL      string  `json:"statuses_url"`
	SubscribersURL   string  `json:"subscribers_url"`
	SubscriptionURL  string  `json:"subscription_url"`
	TagsURL          string  `json:"tags_url"`
	TeamsURL         string  `json:"teams_url"`
	TreesURL         string  `json:"trees_url"`
	CloneURL         string  `json:"clone_url"`
	MirrorURL        *string `json:"mirror_url"`
	HooksURL         string  `json:"hooks_url"`
	SvnURL           string  `json:"svn_url"`
	Homepage         *string `json:"homepage"`

	Language        *string  `json:"language"`
	ForksCount      int      `json:"forks_count"`
	StargazersCount int      `json:"stargazers_count"`
	WatchersCount   int      `json:"watchers_count"`
	Size            int      `json:"size"`
	DefaultBranch   string   `json:"default_branch"`
	OpenIssuesCount int      `json:"open_issues_count"`
	IsTemplate      bool     `json:"is_template"`
	Topics          []string `json:"topics"`

	HasIssues      bool `json:"has_issues"`
	HasProjects    bool `json:"has_projects"`
	HasWiki        bool `json:"has_wiki"`
	HasPages       bool `json:"has_pages"`
	HasDownloads   bool `json:"has_downloads"`
	HasDiscussions bool `json:"has_discussions"`
	Archived       bool `json:"archived"`
	Disabled       bool `json:"disabled"`

	Visibility string `json:"visibility"` // "public" | "private"
	PushedAt   string `json:"pushed_at"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`

	License             *RepoLicense             `json:"license,omitempty"`
	Permissions         *RepoPermissions         `json:"permissions,omitempty"`
	SecurityAndAnalysis *RepoSecurityAndAnalysis `json:"security_and_analysis,omitempty"`
}

type ListUserReposResponse []RepoItem
