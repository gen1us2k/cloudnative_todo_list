package models

type (
	Todo struct {
		ID      int64  `json:"id,omitempty"`
		Title   string `json:"title"`
		Status  string `json:"status"`
		OwnerID string `json:"owner_id"`
	}

	User struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)
