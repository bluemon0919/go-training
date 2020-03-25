package main

// UserSession manages linebot users session
type UserSession struct {
	UserID string
	Count  int
}

// Start starts session
func (us *UserSession) Start(userID string) int {
	if us.UserID == "" {
		us.UserID = userID
	}
	if us.UserID != userID {
		return -1
	}
	return us.Count
}

// Close closes session
func (us *UserSession) Close(userID string) {
	if us.UserID != userID {
		return
	}
	us.Count++
	if us.Count > 3 {
		us.Count = 0
	}
}
