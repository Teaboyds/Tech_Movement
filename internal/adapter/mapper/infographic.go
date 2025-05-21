package mapper

import "backend_tech_movement_hex/internal/core/domain"

func InfographicRequestToDomain(info *domain.InfographicRequestDTO) *domain.Infographic {
	return &domain.Infographic{
		Image:    info.Image,
		Title:    info.Title,
		Category: info.Category,
		Tags:     info.Tags,
		Status:   info.Status,
	}
}

func EnrichInfographics(infographic []*domain.Infographic, categoryMap map[string]domain.Category, imageMap map[string]domain.UploadFile) []*domain.InfographicRespose {
	var result []*domain.InfographicRespose

	for _, info := range infographic {
		category := categoryMap[info.Category] // media.Category เป็น categoryID
		image := imageMap[info.Image]

		resp := &domain.InfographicRespose{
			ID:    info.ID,
			Title: info.Title,
			Image: domain.UploadFile{
				ID:        image.ID,
				Path:      image.Path,
				Name:      image.Name,
				FileType:  image.FileType,
				Type:      image.Type,
				CreatedAt: image.CreatedAt,
				UpdatedAt: image.UpdatedAt,
			},
			Tags:   info.Tags,
			Status: info.Status,
			Category: domain.Category{
				ID:           category.ID,
				Name:         category.Name,
				CategoryType: category.CategoryType,
				CreatedAt:    category.CreatedAt,
				UpdatedAt:    category.UpdatedAt,
			},
			PageView:  info.PageView,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
		}

		result = append(result, resp)
	}

	return result
}

func EnrichInfographic(info *domain.Infographic, category *domain.Category, image *domain.UploadFile) *domain.InfographicRespose {
	return &domain.InfographicRespose{
		ID:    info.ID,
		Title: info.Title,
		Image: domain.UploadFile{
			ID:        image.ID,
			Path:      image.Path,
			Name:      image.Name,
			FileType:  image.FileType,
			Type:      image.Type,
			CreatedAt: image.CreatedAt,
			UpdatedAt: image.UpdatedAt,
		},
		Category: domain.Category{
			ID:           category.ID,
			Name:         category.Name,
			CategoryType: category.CategoryType,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		},
		Tags:      info.Tags,
		Status:    info.Status,
		PageView:  info.PageView,
		CreatedAt: info.CreatedAt,
	}
}

// helper //
func ExtractCategoryAndImageIDs(infographics []*domain.Infographic) []string {
	seen := make(map[string]bool)
	var ids []string

	for _, i := range infographics {
		if !seen[i.Category] {
			ids = append(ids, i.Category)
			seen[i.Category] = true
		}

		if !seen[i.Image] {
			ids = append(ids, i.Image)
			seen[i.Image] = true
		}
	}
	return ids
}
