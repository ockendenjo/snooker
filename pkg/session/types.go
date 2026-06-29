package session

type Data struct {
	ID          string `json:"ID"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Points      int    `json:"points"`
}
