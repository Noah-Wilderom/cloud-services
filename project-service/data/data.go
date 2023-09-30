package data

type Project struct {
	Id           string  `json:"id"`
	Domain       Domain  `json:"domain"`
	Status       string  `json:"status"`
	Name         string  `json:"name"`
	Subdomain    string  `json:"subdomain"`
	Docker       bool    `json:"docker"`
	SFTP         *SFTP   `json:"sftp"`
	PreviewImage *string `json:"preview_image"`
	User         *string `json:"user"`
	SshKeyPath   *string `json:"ssh_key_path"`
	FilesPath    *string `json:"files_path"`
	Git          *Git    `json:"git"`
	Monitoring   bool    `json:"monitoring"`
}

type SFTP struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Path     string `json:"path"`
}

type Git struct {
	Repository string `json:"repository"`
	WebhookUrl string `json:"webhook_url"`
	SshKeyPath string `json:"ssh_key_path"`
}

type Domain struct {
	Id     string `json:"id"`
	Domain string `json:"domain"`
}
