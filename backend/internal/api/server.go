package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/owndangan/backend/internal/api/handler"
	adminHandler "github.com/owndangan/backend/internal/api/handler/admin"
	"github.com/owndangan/backend/internal/api/middleware"
	"github.com/owndangan/backend/internal/config"
	"github.com/owndangan/backend/internal/pkg/jwt"
	"github.com/owndangan/backend/internal/pkg/storage"
	"github.com/owndangan/backend/internal/repository"
	"github.com/owndangan/backend/internal/service"
	"github.com/owndangan/backend/internal/service/admin"
	"github.com/owndangan/backend/internal/service/email"
	"github.com/owndangan/backend/internal/service/guest"
	"github.com/owndangan/backend/internal/service/guestbook"
	"github.com/owndangan/backend/internal/service/rsvp"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	router *chi.Mux
}

type Dependencies struct {
	UserRepo               repository.UserRepository
	RefreshTokenRepo       repository.RefreshTokenRepository
	PackageRepo            repository.PackageRepository
	TransactionRepo        repository.TransactionRepository
	SubscriptionRepo       repository.SubscriptionRepository
	EventRepo              repository.EventRepository
	EventSectionRepo       repository.EventSectionRepository
	TemplateRepo           repository.TemplateRepository
	MusicRepo              repository.MusicRepository
	GuestRepo              repository.GuestRepository
	RSVPRepo               repository.RSVPRepository
	GuestbookRepo          repository.GuestbookRepository
	LoveStoryRepo          repository.LoveStoryRepository
	DigitalGiftRepo        repository.DigitalGiftRepository
	GalleryPhotoRepo       repository.GalleryPhotoRepository
	AnalyticsRepo          repository.AnalyticsEventRepository
	AuditLogRepo           repository.AuditLogRepository
	WebhookIdempotencyRepo repository.WebhookIdempotencyRepository
	JWTService             *jwt.Service
	Storage                storage.Storage
	EmailSender            service.EmailSender
}

func NewServer(cfg *config.Config, deps *Dependencies, db *gorm.DB, log zerolog.Logger) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	r.Use(middleware.RateLimitPerIP(100, 60*time.Second))

	jwtSvc := deps.JWTService
	authRequired := middleware.Authenticate(jwtSvc)
	adminRequired := middleware.RequireRole("admin")

	emailSender := deps.EmailSender
	if emailSender == nil {
		emailSender = email.NewService(cfg.SMTP, log)
	}

	healthHandler := handler.NewHealthHandler()

	authSvc := service.NewAuthService(
		deps.UserRepo, deps.RefreshTokenRepo, deps.PackageRepo, deps.SubscriptionRepo, jwtSvc, deps.AuditLogRepo, emailSender, cfg.FrontendURL,
	)
	authHandler := handler.NewAuthHandler(authSvc)

	pkgSvc := service.NewPackageService(deps.PackageRepo, deps.AuditLogRepo)
	pkgHandler := handler.NewPackageHandler(pkgSvc)

	userSvc := service.NewUserService(deps.UserRepo, deps.SubscriptionRepo, deps.PackageRepo, deps.AuditLogRepo)
	userHandler := handler.NewUserHandler(userSvc)

	subSvc := service.NewSubscriptionService(deps.SubscriptionRepo, deps.PackageRepo, deps.TransactionRepo, deps.UserRepo, deps.AuditLogRepo)
	subHandler := handler.NewSubscriptionHandler(subSvc)

	paySvc := service.NewPaymentService(
		deps.TransactionRepo, deps.PackageRepo, deps.UserRepo, deps.AuditLogRepo,
		deps.WebhookIdempotencyRepo, cfg.Midtrans, subSvc, emailSender,
	)
	payHandler := handler.NewPaymentHandler(paySvc)

	eventSvc := service.NewEventService(
		db, deps.EventRepo, deps.EventSectionRepo, deps.DigitalGiftRepo,
		deps.SubscriptionRepo, deps.PackageRepo, deps.GuestRepo, deps.RSVPRepo,
		deps.GuestbookRepo, deps.LoveStoryRepo, deps.TemplateRepo, deps.MusicRepo, deps.GalleryPhotoRepo,
		deps.AuditLogRepo, deps.AnalyticsRepo, deps.Storage,
	)
	EventHandler := handler.NewEventHandler(eventSvc)

	analyticsSvc := service.NewAnalyticsService(
		deps.EventRepo, deps.AnalyticsRepo, deps.RSVPRepo, deps.SubscriptionRepo, deps.PackageRepo,
	)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)

	guestSvc := guest.NewService(
		deps.GuestRepo, deps.EventRepo, deps.SubscriptionRepo, deps.PackageRepo, deps.AuditLogRepo,
	)
	guestHandler := handler.NewGuestHandler(guestSvc)

	rsvpSvc := rsvp.NewService(
		deps.RSVPRepo, deps.GuestRepo, deps.EventRepo, deps.UserRepo, deps.AuditLogRepo, emailSender,
	)
	rsvpHandler := handler.NewRSVPHandler(rsvpSvc)

	guestbookSvc := guestbook.NewService(
		deps.GuestbookRepo, deps.EventRepo, deps.UserRepo, deps.AuditLogRepo, emailSender,
	)
	guestbookHandler := handler.NewGuestbookHandler(guestbookSvc)

	adminSvc := admin.NewService(
		deps.UserRepo, deps.PackageRepo, deps.TransactionRepo, deps.TemplateRepo, deps.SubscriptionRepo, deps.EventRepo, deps.AuditLogRepo,
	)
	adminHandler := adminHandler.NewHandler(adminSvc)

	r.Get("/health", healthHandler.HealthCheck)

	if cfg.Storage.Provider == "local" && cfg.Storage.LocalPath != "" {
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.Storage.LocalPath))))
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r, authRequired)
		})

		r.Route("/packages", func(r chi.Router) {
			r.Get("/", pkgHandler.ListActive)
			r.Group(func(r chi.Router) {
				r.Use(authRequired, adminRequired)
				r.Get("/all", pkgHandler.ListAll)
				r.Post("/", pkgHandler.Create)
				r.Get("/{id}", pkgHandler.GetByID)
				r.Put("/{id}", pkgHandler.Update)
				r.Delete("/{id}", pkgHandler.Delete)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authRequired)
				r.Get("/me", userHandler.Profile)
				r.Put("/me", userHandler.UpdateProfile)
				r.Get("/me/subscription", userHandler.GetUserSubscription)
			})
		})

		r.Route("/subscriptions", func(r chi.Router) {
			subHandler.RegisterRoutes(r, authRequired)
		})

		r.Route("/payments", func(r chi.Router) {
			payHandler.RegisterRoutes(r, authRequired)
		})

		r.Route("/events", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authRequired)
				r.Post("/", EventHandler.Create)
				r.Get("/", EventHandler.List)
				r.Get("/{id}", EventHandler.GetByID)
				r.Put("/{id}", EventHandler.Update)
				r.Delete("/{id}", EventHandler.Delete)
				r.Post("/{id}/publish", EventHandler.Publish)
				r.Post("/{id}/unpublish", EventHandler.Unpublish)
				r.Get("/{id}/analytics", analyticsHandler.GetEventAnalytics)

				r.Put("/{id}/template", EventHandler.AssignTemplate)
				r.Get("/{id}/sections", EventHandler.GetSections)
				r.Put("/{id}/sections", EventHandler.UpdateSections)
				r.Get("/{id}/digital-gifts", EventHandler.GetDigitalGift)
				r.Put("/{id}/digital-gifts", EventHandler.UpdateDigitalGift)
				r.Route("/{id}/gallery", func(r chi.Router) {
					r.Get("/", EventHandler.ListGallery)
					r.Post("/upload", EventHandler.UploadGallery)
					r.Delete("/{photoID}", EventHandler.DeleteGallery)
					r.Put("/reorder", EventHandler.ReorderGallery)
				})
				r.Route("/{id}/music", func(r chi.Router) {
					r.Get("/", EventHandler.GetMusic)
					r.Post("/upload", EventHandler.UploadMusic)
					r.Get("/presets", EventHandler.ListMusicPresets)
					r.Post("/presets", EventHandler.AssignPresetMusic)
					r.Delete("/", EventHandler.RemoveMusic)
				})
				r.Route("/{id}/love-stories", func(r chi.Router) {
					r.Get("/", EventHandler.ListLoveStories)
					r.Post("/", EventHandler.CreateLoveStory)
					r.Put("/{storyID}", EventHandler.UpdateLoveStory)
					r.Delete("/{storyID}", EventHandler.DeleteLoveStory)
				})
			})
		})

		r.Route("/templates", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authRequired)
				r.Get("/", EventHandler.ListTemplates)
			})
		})

		r.Get("/e/{slug}", EventHandler.PublicView)

		r.Get("/invitations/public", EventHandler.PublicList)

		r.Post("/analytics/events", analyticsHandler.TrackEvent)

		r.Route("/admin", func(r chi.Router) {
			adminHandler.RegisterRoutes(r, authRequired, adminRequired)
		})

		r.Route("/events/{eventID}/guests", func(r chi.Router) {
			guestHandler.RegisterRoutes(r, authRequired)
		})

		r.Route("/rsvp", func(r chi.Router) {
			rsvpHandler.RegisterRoutes(r, authRequired, func(h http.Handler) http.Handler { return h })
		})

		r.Route("/guestbook", func(r chi.Router) {
			guestbookHandler.RegisterRoutes(r, authRequired, func(h http.Handler) http.Handler { return h })
		})
	})

	s := &Server{router: r}
	return s
}

func (s *Server) Handler() http.Handler {
	return s.router
}
