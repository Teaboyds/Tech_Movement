package utils

/////////////// case ContentStatus ใช้ใน handler news (Create , Update) ///////////////

var validContentStatus = map[string]bool{
	"draft":     true,
	"published": true,
	"archived":  true,
}

func IsValidContentStatus(status string) bool {
	return validContentStatus[status]
}

///////////////////////////////// case ContentStatus //////////////////////////////
