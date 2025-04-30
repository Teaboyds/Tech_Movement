package mongoUtils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(hexStr, label string) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(hexStr)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid %s ID: %w", label, err)
	}
	return id, nil
}
