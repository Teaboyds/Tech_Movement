// entity //
package domain

import "time"

type News struct {
	ID          string
	ThumnailID  string
	Title       string
	Description string
	Content     string
	ImageIDs    []string
	CategoryID  string
	Tags        []string
	Status      string
	ContentType string
	View        string
	CreatedAt   string
	UpdatedAt   string
}

//  HtmlContent string `json:"html_content"`
// Read string `json:"read"`

////////////////////////////////// Response // Request // Models /////////////////////////////////////////////////

type NewsResponse struct {
	ID          string               `json:"id"`
	ThumnailID  UploadFileResponse   `json:"thumnail_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Content     string               `json:"content"`
	ImageIDs    []UploadFileResponse `json:"image_ids"`
	CategoryID  CategoryResponse     `json:"category_id"`
	Tags        []string             `json:"tag"`
	Status      string               `json:"status"`
	ContentType string               `json:"content_type"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type NewsResponseV2 struct {
	ID          string             `json:"id"`
	ThumnailID  UploadFileResponse `json:"thumnail_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Content     string             `json:"content"`
	CategoryID  CategoryResponse   `json:"category_id"`
	Tags        []string           `json:"tag"`
	Status      string             `json:"status"`
	ContentType string             `json:"content_type"`
	View        string             `json:"view"`
	CreatedAt   string             `json:"created_at"`
}

type NewsRequest struct {
	Title       string   `json:"title" validate:"required"`
	ThumnailID  string   `json:"thumnail_id" form:"thumnail_id"`
	Description string   `json:"description" validate:"required"`
	Content     string   `json:"content" validate:"required"`
	ImageIDs    []string `json:"image_ids"`
	CategoryID  string   `json:"category_id" form:"category_id"`
	Tags        []string `json:"tags" validate:"required,min=1,required"`
	Status      string   `json:"status" validate:"required"`
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
	Title       string               `json:"title"`
	Detail      string               `json:"detail"`
	Image       []UploadFileResponse `json:"image"`
	Category    CategoryResponse     `json:"category"`
	ContentType string               `json:"content_type"`
	CreatedAt   time.Time            `json:"created_at"`
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

type PaginationResp struct {
	TotalItems  string
	TotalPages  string
	CurrentPage string
	PageSize    string
}

type DeleteManyID struct {
	IDs []string `json:"ids"`
}

////////////////////////////////// Response Models /////////////////////////////////////////////////
