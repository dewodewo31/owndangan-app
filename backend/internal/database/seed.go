package database

import (
	"github.com/owndangan/backend/internal/model"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedPackages(db *gorm.DB) error {
	var count int64
	db.Model(&model.Package{}).Where("code = ?", "free").Count(&count)
	if count > 0 {
		return nil
	}

	freePkg := &model.Package{
		Name:          "Free Trial",
		Code:          "free",
		Price:         0,
		GuestLimit:    intPtr(50),
		TemplateGroup: "standard",
		Features:      datatypes.JSON(`{"guest.max":50,"gallery.max":10,"music.preset":true,"music.upload":false,"digital_gift.qris":false,"custom_domain":false,"template_group":"standard"}`),
		IsActive:      true,
	}
	if err := db.Create(freePkg).Error; err != nil {
		return err
	}

	starterPkg := &model.Package{
		Name:          "Starter",
		Code:          "starter",
		Price:         99000,
		DurationDays:  intPtr(30),
		GuestLimit:    intPtr(100),
		TemplateGroup: "standard",
		Features:      datatypes.JSON(`{"guest.max":100,"gallery.max":100,"music.preset":true,"music.upload":true,"digital_gift.qris":true,"custom_domain":false,"rsvp.export":true,"template_group":"standard"}`),
		IsActive:      true,
	}
	if err := db.Create(starterPkg).Error; err != nil {
		return err
	}

	premiumPkg := &model.Package{
		Name:          "Premium",
		Code:          "premium",
		Price:         299000,
		DurationDays:  intPtr(60),
		GuestLimit:    nil,
		TemplateGroup: "all",
		Features:      datatypes.JSON(`{"guest.max":null,"gallery.max":null,"music.preset":true,"music.upload":true,"video.enabled":true,"digital_gift.qris":true,"custom_domain":true,"rsvp.export":true,"watermark.removed":true,"template_group":"all"}`),
		IsActive:      true,
	}
	if err := db.Create(premiumPkg).Error; err != nil {
		return err
	}

	allPkg := &model.Package{
		Name:          "All Access",
		Code:          "all",
		Price:         999000,
		GuestLimit:    nil,
		TemplateGroup: "all",
		Features:      datatypes.JSON(`{"guest.max":null,"gallery.max":null,"music.preset":true,"music.upload":true,"video.enabled":true,"digital_gift.qris":true,"custom_domain":true,"rsvp.export":true,"watermark.removed":true,"template.all":true,"template_group":"all"}`),
		IsActive:      true,
	}
	if err := db.Create(allPkg).Error; err != nil {
		return err
	}

	var tmplCount int64
	db.Model(&model.Template{}).Where("name = ?", "Modern Minimalist").Count(&tmplCount)
	if tmplCount == 0 {
		templates := []*model.Template{
			{
				Name:         "Modern Minimalist",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#5b6470","secondary_color":"#9aa3af","background_color":"#ffffff","font_family":"sans-serif","hero_image":"https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"split","gallery":"grid"}`),
				IsActive:     true,
			},
			{
				Name:         "Romantic Elegant",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#b8860b","secondary_color":"#8b6914","background_color":"#fbf8f1","font_family":"serif","hero_image":"https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"cover","gallery":"grid"}`),
				IsActive:     true,
			},
			{
				Name:         "Rustic Bohemian",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#a0522d","secondary_color":"#c08552","background_color":"#faf3ec","font_family":"handwritten","hero_image":"https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"full","gallery":"masonry"}`),
				IsActive:     true,
			},
			{
				Name:         "Contemporary Editorial",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#2e7d52","secondary_color":"#4caf7d","background_color":"#eef6f0","font_family":"sans-serif","hero_image":"https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"cover","gallery":"grid"}`),
				IsActive:     true,
			},
			{
				Name:         "Luxury Black & Gold",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#7b1e3b","secondary_color":"#b03a5b","background_color":"#fbf0f3","font_family":"serif","hero_image":"https://images.unsplash.com/photo-1465495976272-a759b7a8e978?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"full","gallery":"carousel"}`),
				IsActive:     true,
			},
			{
				Name:         "Islamic",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1511285560929-80b456fea0bc?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#0f766e","secondary_color":"#14b8a6","background_color":"#f0fdfa","font_family":"sans-serif","hero_image":"https://images.unsplash.com/photo-1511285560929-80b456fea0bc?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"cover","gallery":"grid"}`),
				IsActive:     true,
			},
			{
				Name:         "Japanese Zen",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1520854229188-3da5971b5e0c?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#374151","secondary_color":"#9ca3af","background_color":"#f9fafb","font_family":"sans-serif","hero_image":"https://images.unsplash.com/photo-1520854229188-3da5971b5e0c?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"cover","gallery":"grid"}`),
				IsActive:     true,
			},
			{
				Name:         "Javanese Traditional",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#92400e","secondary_color":"#b45309","background_color":"#fffbeb","font_family":"serif","hero_image":"https://images.unsplash.com/photo-1519741497674-9b181a268d28?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"full","gallery":"masonry"}`),
				IsActive:     true,
			},
			{
				Name:         "Sundanese Traditional",
				GroupName:    "standard",
				ThumbnailURL: "https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=800&q=60",
				CSSConfig:    datatypes.JSON(`{"primary_color":"#1e3a8a","secondary_color":"#3b82f6","background_color":"#eff6ff","font_family":"serif","hero_image":"https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=1600&q=60"}`),
				LayoutConfig: datatypes.JSON(`{"hero":"cover","gallery":"masonry"}`),
				IsActive:     true,
			},
		}
		for _, t := range templates {
			if err := db.Create(t).Error; err != nil {
				return err
			}
		}
	}

	var botanicalCount int64
	db.Model(&model.Template{}).Where("name = ?", "Modern Botanical").Count(&botanicalCount)
	if botanicalCount == 0 {
		if err := db.Create(&model.Template{
			Name:         "Modern Botanical",
			GroupName:    "standard",
			ThumbnailURL: "https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=800&q=60",
			CSSConfig:    datatypes.JSON(`{"primary_color":"#8A9A86","secondary_color":"#6B7B66","background_color":"#FDFBF7","font_family":"serif","hero_image":"https://images.unsplash.com/photo-1519225421980-715cb0215aed?auto=format&fit=crop&w=1600&q=60"}`),
			LayoutConfig: datatypes.JSON(`{"hero":"centered","gallery":"columns"}`),
			IsActive:     true,
		}).Error; err != nil {
			return err
		}
	}

	var corporateCount int64
	db.Model(&model.Template{}).Where("name = ?", "Modern Corporate").Count(&corporateCount)
	if corporateCount == 0 {
		if err := db.Create(&model.Template{
			Name:         "Modern Corporate",
			GroupName:    "standard",
			ThumbnailURL: "https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&w=800&q=60",
			CSSConfig:    datatypes.JSON(`{"primary_color":"#818cf8","secondary_color":"#1e293b","background_color":"#0f172a","font_family":"sans-serif","hero_image":"https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&w=1600&q=60"}`),
			LayoutConfig: datatypes.JSON(`{"hero":"centered","gallery":"grid"}`),
			IsActive:     true,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func intPtr(i int) *int {
	return &i
}
