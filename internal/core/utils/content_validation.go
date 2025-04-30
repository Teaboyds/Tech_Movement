package utils

/////////////// case Content ใช้ใน handler news (Create , Update) ///////////////

var validContentStatus = map[string]bool{
	"draft":     true,
	"published": true,
	"archived":  true,
}

func IsValidContentStatus(status string) bool {
	return validContentStatus[status]
}

var ContentType = map[string]bool{
	"general":  true,
	"breaking": true,
	"video":    true,
}

func IsContentType(status string) bool {
	return ContentType[status]
}

var ImageType = map[string]bool{
	"news":        true,
	"banner":      true,
	"infographic": true,
}

func IsImageType(status string) bool {
	return ImageType[status]
}

///////////////////////////////// case Content //////////////////////////////
