package mapper

import "backend_tech_movement_hex/internal/core/domain"

func EnrichMedias(medias []*domain.Media, categoryMap map[string]domain.Category) []*domain.MediaResponse {
	var result []*domain.MediaResponse

	for _, media := range medias {
		category := categoryMap[media.Category] // media.Category เป็น categoryID

		resp := &domain.MediaResponse{
			ID:       media.ID,
			Title:    media.Title,
			Content:  media.Content,
			VideoUrl: media.VideoURL,
			Thumnail: media.Thumnail,
			Category: domain.Category{
				ID:           category.ID,
				Name:         category.Name,
				CategoryType: category.CategoryType,
				CreatedAt:    category.CreatedAt,
				UpdatedAt:    category.UpdatedAt,
			},
			Tags:      media.Tags,
			PageView:  media.View,
			Action:    media.Action,
			CreatedAt: media.CreatedAt,
		}

		result = append(result, resp)
	}

	return result
}

func EnrichMedia(media *domain.Media, category *domain.Category) *domain.MediaResponse {

	resp := &domain.MediaResponse{
		ID:       media.ID,
		Title:    media.Title,
		Content:  media.Content,
		VideoUrl: media.VideoURL,
		Thumnail: media.Thumnail,
		Category: domain.Category{
			ID:           category.ID,
			Name:         category.Name,
			CategoryType: category.CategoryType,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		},
		Tags:      media.Tags,
		PageView:  media.View,
		Action:    media.Action,
		CreatedAt: media.CreatedAt,
	}

	return resp
}

func MediaRequestToDomain(req domain.MediaRequest) *domain.Media {
	return &domain.Media{
		Title:    req.Title,
		Content:  req.Content,
		VideoURL: req.VideoURL,
		Thumnail: req.Thumnail,
		Category: req.Category,
		Tags:     req.Tags,
		Action:   req.Action,
	}
}

// helper //
func ExtractCategoryIDs(medias []*domain.Media) []string {
	seen := make(map[string]bool)
	var ids []string

	for _, m := range medias {
		if !seen[m.Category] {
			ids = append(ids, m.Category)
			seen[m.Category] = true
		}
	}
	return ids
}
