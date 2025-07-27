package api

type WebHookRequest struct {
	Text interface{} `json:"text"`
}

type WebHookAlertPayloadRequest struct {
	RequestID   string `json:"request_id"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

type SendEmailRequest struct {
	From       string `json:"From"`
	Receiver   string `json:"Receiver"`
	Subject    string `json:"Subject"`
	Message    string `json:"Message"`
	Preference string `json:"Preference"`
}

type SendNotificationRequest struct {
	Title     string `schema:"title" validate:"required,omitempty"`
	Codebuild string `schema:"codebuild" validate:"required,omitempty"`
	Status    string `schema:"status" validate:"required,omitempty"`
	Type      string `schema:"type" validate:"required,omitempty"`
	HotelId   string `schema:"hotel_id"`
	BrandId   string `schema:"brand_id"`
	Website   string `schema:"website"`
	Author    string `schema:"author"`
	Notes     string `schema:"notes"`
}

type NotificationCard struct {
	Header struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		ImageUrl string `json:"imageUrl"`
	} `json:"header"`
	Sections []struct {
		Widgets []struct {
			TextParagraph struct {
				Text string `json:"text"`
			} `json:"textParagraph"`
		} `json:"widgets"`
	} `json:"sections"`
}

type NotificationCardRequest struct {
	NotificationCard []NotificationCard `json:"cards"`
}

type PaginationRequest struct {
	Page      int64  `schema:"page" validate:"number"`
	PageSize  int64  `schema:"page_size" validate:"number"`
	SortBy    string `schema:"sort_by"`
	SortOrder string `schema:"sort_order"`
}
