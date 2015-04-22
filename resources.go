package github

type Repository struct {
	Id                int
	Owner             Owner
	Name              string
	Full_name         string
	Description       string
	Private           bool
	Fork              bool
	Url               string
	Html_url          string
	Clone_url         string
	Git_url           string
	Ssh_url           string
	Svn_url           string
	Mirror_url        string
	Homepage          string
	Language          string
	Forks_count       int
	Stargazers_count  int
	Watchers_count    int
	Size              int
	Default_branch    string
	Open_issues_count int
	Has_issues        bool
	Has_wiki          bool
	Has_pages         bool
	Has_downloads     bool
	Pushed_at         string
	Created_at        string
	Updated_at        string
	Permissions       Permissions
}

type Permissions struct {
	Admin bool
	Push  bool
	Pull  bool
}

type Owner struct {
	Login               string
	Id                  int
	Avatar_url          string
	Gravatar_id         string
	Url                 string
	Html_url            string
	Followers_url       string
	Following_url       string
	Gists_url           string
	Starred_url         string
	Subscriptions_url   string
	Organizations_url   string
	Repos_url           string
	Events_url          string
	Received_events_url string
	Type                string
	Site_admin          bool
}

type Commit struct {
	Url          string
	Sha          string
	Html_url     string
	Comments_url string
	Commit       SubCommit
	Author       User
	Committer    User
	Parents      []Parent
}

type Parent struct {
	Url string
	Sha string
}

type SubCommit struct {
	Url           string
	Author        Author
	Committer     Author
	Message       string
	Tree          Tree
	Comment_count int
}

type Author struct {
	Name  string
	Email string
	Date  string
}

type Tree struct {
	Url string
	Sha string
}

type User struct {
	Login               string
	Id                  int
	Avatar_url          string
	Gravatar_id         string
	Url                 string
	Html_url            string
	Followers_url       string
	Following_url       string
	Gists_url           string
	Starred_url         string
	Subscriptions_url   string
	Organizations_url   string
	Repos_url           string
	Events_url          string
	Received_events_url string
	Type                string
	Site_admin          bool
}
