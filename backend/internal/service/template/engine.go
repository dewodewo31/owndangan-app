package template

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/owndangan/backend/internal/model"
	"github.com/owndangan/backend/internal/pkg/errors"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service/entitlement"
)

type Engine struct {
	templateRepo repository.TemplateRepository
	pkgRepo      repository.PackageRepository
	storage      Storage
}

type Storage interface {
	Upload(ctx context.Context, key string, data []byte, opts UploadOpts) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) string
}

type UploadOpts struct {
	ContentType string
	Extension   string
	MaxSize     int64
}

type UploadResult struct {
	URL string
	Key string
}

func NewEngine(templateRepo repository.TemplateRepository, pkgRepo repository.PackageRepository, storage Storage) *Engine {
	return &Engine{
		templateRepo: templateRepo,
		pkgRepo:      pkgRepo,
		storage:      storage,
	}
}

func (e *Engine) GetTemplate(ctx context.Context, templateID uuid.UUID) (*model.Template, error) {
	t, err := e.templateRepo.GetByID(ctx, templateID)
	if err != nil || t == nil {
		return nil, errors.ErrNotFound
	}
	if !t.IsActive {
		return nil, fmt.Errorf("%w: template is not active", errors.ErrConflict)
	}
	return t, nil
}

func (e *Engine) ListTemplates(ctx context.Context, userID uuid.UUID) ([]model.Template, error) {
	sub, err := e.getActiveSubscription(ctx, userID)
	if err != nil {
		return e.templateRepo.ListByGroups(ctx, []string{"standard"})
	}

	pkg, err := e.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return e.templateRepo.ListByGroups(ctx, []string{"standard"})
	}

	resolver := entitlement.NewResolver(pkg)
	groups := []string{"standard"}

	if resolver.CanAccessPremiumTemplates() {
		groups = append(groups, "premium")
	}
	if resolver.CanAccessAllTemplates() {
		groups = append(groups, "all")
	}

	return e.templateRepo.ListByGroups(ctx, groups)
}

func (e *Engine) CanUseTemplate(ctx context.Context, userID uuid.UUID, templateID uuid.UUID) (bool, error) {
	t, err := e.templateRepo.GetByID(ctx, templateID)
	if err != nil || t == nil {
		return false, nil
	}

	sub, err := e.getActiveSubscription(ctx, userID)
	if err != nil {
		return false, nil
	}

	pkg, err := e.pkgRepo.GetByID(ctx, sub.PackageID)
	if err != nil || pkg == nil {
		return false, nil
	}

	resolver := entitlement.NewResolver(pkg)

	switch t.GroupName {
	case "standard":
		return true, nil
	case "premium":
		return resolver.CanAccessPremiumTemplates(), nil
	case "all":
		return resolver.CanAccessAllTemplates(), nil
	default:
		return false, nil
	}
}

func (e *Engine) RenderInvitation(ctx context.Context, event *model.Event, template *model.Template) (*RenderedInvitation, error) {
	if template == nil {
		return nil, fmt.Errorf("%w: template is required", errors.ErrInvalidInput)
	}

	sections := e.buildSections(event, template)

	return &RenderedInvitation{
		Event:    event,
		Template: template,
		Sections: sections,
	}, nil
}

func (e *Engine) buildSections(event *model.Event, template *model.Template) []Section {
	var sections []Section

	if event.Sections != nil {
		if event.Sections.HeroEnabled {
			sections = append(sections, Section{Type: "cover", Enabled: true})
		}
		if event.Sections.OpeningMessage != "" {
			sections = append(sections, Section{Type: "opening", Enabled: true})
		}
		if event.Sections.CoupleEnabled {
			sections = append(sections, Section{Type: "couple", Enabled: true})
		}
		if event.Sections.EventDetailsEnabled {
			sections = append(sections, Section{Type: "events", Enabled: true})
		}
		if event.Sections.GalleryEnabled {
			sections = append(sections, Section{Type: "gallery", Enabled: true})
		}
		if event.Sections.RSVPEnabled {
			sections = append(sections, Section{Type: "rsvp", Enabled: true})
		}
		if event.Sections.DigitalGiftsEnabled {
			sections = append(sections, Section{Type: "digital_gift", Enabled: true})
		}
		if event.Sections.DressCode != "" {
			sections = append(sections, Section{Type: "dress_code", Enabled: true})
		}
		if event.Sections.ClosingMessage != "" {
			sections = append(sections, Section{Type: "closing", Enabled: true})
		}
	}

	return sections
}

func (e *Engine) getActiveSubscription(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	return nil, fmt.Errorf("not found")
}

type RenderedInvitation struct {
	Event    *model.Event
	Template *model.Template
	Sections []Section
}

type Section struct {
	Type    string
	Enabled bool
}
