package types

type TaskWithFormattedTime struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Model     string    `json:"model"`
	Dataset   string    `json:"dataset"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Progress  string    `json:"progress"`
	Accuracy  []float32 `json:"accuracy"`
	Loss      []float32 `json:"loss"`
}

type ServerWithFormattedTime struct {
	ServerID  string  `json:"server_id"`
	TaskID    int     `json:"task_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Model     string  `json:"model"`
	Dataset   string  `json:"dataset"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Progress  string  `json:"progress"`
	Accuracy  float32 `json:"accuracy"`
	Loss      float32 `json:"loss"`
	Cpu       string  `json:"cpu"`
	Memory    string  `json:"memory"`
	Disk      string  `json:"disk"`
}

type ClientWithFormattedTime struct {
	ClientID  string  `json:"client_id"`
	TaskID    int     `json:"task_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Model     string  `json:"model"`
	Dataset   string  `json:"dataset"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Progress  string  `json:"progress"`
	Accuracy  float32 `json:"accuracy"`
	Loss      float32 `json:"loss"`
	Cpu       string  `json:"cpu"`
	Memory    string  `json:"memory"`
	Disk      string  `json:"disk"`
}
