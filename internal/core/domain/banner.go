package domain

type Banner struct {
	ID           string
	DesktopImage ImageInfo
	MobileImage  ImageInfo
	Status       StatusType
	LinkUrl      string
	Action       string
	CreatedAt    string
	UpdatedAt    string
}

type BannerV2 struct {
	ID           string
	DesktopImage string
	MobileImage  string
	Status       StatusType
	LinkUrl      string
	Action       string
	CreatedAt    string
	UpdatedAt    string
}

type BannerRequestV2 struct {
	DesktopImage string     `json:"desktop_image" form:"desktop_image"`
	MobileImage  string     `json:"mobile_image" form:"mobile_image"`
	Status       StatusType `json:"status" form:"status"`
	LinkUrl      string     `json:"link_url" form:"link_url"`
	Action       string     `json:"action" form:"action"`
	CreatedAt    string     `json:"created_at" form:"created_at"`
	UpdatedAt    string     `json:"updated_at" form:"updated_at"`
	BannerType   string     //Home//
}

type StatusType struct {
	Home        bool
	Media       bool
	News        bool
	Infographic bool
}

type BannerResponseV2 struct {
	ID              string     `json:"id"`
	DesktopImageUrl MetaData   `json:"desktop_image_url"`
	MobileImageUrl  MetaData   `json:"mobile_image_url"`
	Status          StatusType `json:"status"`
	Action          string     `json:"action"`
	CreatedAt       string     `json:"created_at"`
	UpdatedAt       string     `json:"updated_at"`
}

type BannerRequest struct {
	DesktopImage ImageInfo  `json:"desktop_image" form:"desktop_image"`
	MobileImage  ImageInfo  `json:"mobile_image" form:"mobile_image"`
	Status       StatusType `json:"status" form:"status"`
	LinkUrl      string     `json:"link_url" form:"link_url"`
	Action       string     `json:"action" form:"action"`
	CreatedAt    string     `json:"created_at" form:"created_at"`
	UpdatedAt    string     `json:"updated_at" form:"updated_at"`
}

type BannerResponse struct {
	ID              string     `json:"id"`
	DesktopImageUrl string     `json:"desktop_image_url"`
	MobileImageUrl  string     `json:"mobile_image_url"`
	Status          StatusType `json:"status"`
	Action          string
}

type ImageInfo struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	FileType string `form:"file_type" validate:"required,oneof=banner infographic news" json:"file_type"`
	Type     string `json:"type"`
}

type MetaData struct {
	Alt      string `json:"alt"`
	Url      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	Type     string `json:"type"`
}
