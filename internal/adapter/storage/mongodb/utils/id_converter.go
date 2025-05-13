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

func ConvertStringToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid object ID format: %w", err)
	}
	return objectID, nil
}
