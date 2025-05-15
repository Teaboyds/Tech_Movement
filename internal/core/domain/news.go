// entity //
package domain

import "time"

type News struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Image       []string `json:"image"`
	CategoryID  string   `json:"category"`
	Tag         []string `json:"tag"`
	Status      bool     `json:"status"`
	ContentType string   `json:"content_type"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

//  HtmlContent string `json:"html_content"`
// Read string `json:"read"`

////////////////////////////////// Response // Request // Models /////////////////////////////////////////////////

type NewsResponse struct {
	ID          string                       `json:"id"`
	Title       string                       `json:"title"`
	Description string                       `json:"description"`
	Content     string                       `json:"content"`
	Image       []UploadFileResponseHomePage `json:"image"`
	CategoryID  CategoryResponse             `json:"category"`
	Tag         []string                     `json:"tag"`
	Status      bool                         `json:"status"`
	ContentType string                       `json:"content_type"`
	CreatedAt   string                       `json:"created_at"`
	UpdatedAt   string                       `json:"updated_at"`
}

type NewsRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Content     string   `json:"content" validate:"required"`
	Image       []string `json:"news_image" validate:"required"`
	Category    string   `json:"category"`
	Tag         []string `json:"tag" validate:"required,min=1,required"`
	Status      bool     `json:"status" validate:"required"`
	ContentType string   `form:"content_type" validate:"required,oneof=Global_Tech Local_Tech"`
}

type UpdateNewsRequest struct {
	Title       *string   `json:"title" validate:"required"`
	Description *string   `json:"description" validate:"required"`
	Content     *string   `json:"content" validate:"required"`
	Image       *[]string `json:"news_image" validate:"required"`
	Category    *string   `json:"category"`
	Tag         *[]string `json:"tag" validate:"required,min=1,required"`
	Status      *bool     `json:"status" validate:"required"`
	ContentType *string   `json:"content_type" form:"content_type" validate:"required,oneof=Global_Tech Local_Tech"`
}

// News Home::LastedNews//
type HomePageLastedNewResponse struct {
	Title       string                       `json:"title"`
	Detail      string                       `json:"detail"`
	Image       []UploadFileResponseHomePage `json:"image"`
	Category    CategoryResponse             `json:"category"`
	ContentType string                       `json:"content_type"`
	CreatedAt   time.Time                    `json:"created_at"`
}

// landing page //
type Home struct {
	Message        string      `json:"message"`
	Video          interface{} `json:"video"`
	LastedNews     interface{} `json:"lasted_news"`
	TechnologyNews interface{} `json:"technology_news"`
	Short          interface{} `json:"short_video"`
	Infographic    interface{} `json:"infographic"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

////////////////////////////////// Response Models /////////////////////////////////////////////////
