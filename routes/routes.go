package routes

import (
	"invitation-app/handlers"
	"invitation-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.GlobalRateLimit())

	api := r.Group("/api/v1")

	// ─── PUBLIC ──────────────────────────────────────────────

	// User auth
	userAuth := api.Group("/auth")
	userAuth.Use(middleware.AuthRateLimit())
	{
		// Magic Link
		userAuth.POST("/magic-link", handlers.RequestMagicLink)
		userAuth.GET("/verify", handlers.VerifyMagicLink)
		userAuth.GET("/verify-api", handlers.VerifyMagicLinkAPI)

		// Google OAuth
		userAuth.GET("/google", handlers.GoogleLogin)
		userAuth.GET("/google/callback", handlers.GoogleCallback)
	}

	// Admin auth — password only
	adminAuth := api.Group("/admin/auth")
	adminAuth.Use(middleware.AuthRateLimit())
	{
		adminAuth.POST("/login", handlers.AdminLogin)
	}

	// Public invitation
	public := api.Group("")
	public.Use(middleware.APIRateLimit())
	{
		public.GET("/templates/public", handlers.GetPublicTemplates)
		public.GET("/templates/public/:id", handlers.GetTemplate)
		public.GET("/invitations/:slug", handlers.GetInvitationBySlug)
		public.GET("/invitations/:slug/messages", handlers.GetMessages)
		public.POST("/invitations/:slug/messages",
			middleware.MessageRateLimit(),
			handlers.CreateMessage,
		)
	}

	// Webhook
	r.POST("/webhook/xendit", handlers.XenditWebhook)

	// ─── PROTECTED (user & admin) ─────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(), middleware.APIRateLimit())
	{
		protected.GET("/auth/me", handlers.Me)

		protected.GET("/users/:id", handlers.GetUser)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.GET("/invitations/me", handlers.GetMyInvitations)
		protected.GET("/invitations/me/:id", handlers.GetMyInvitation)
		protected.POST("/invitations", handlers.CreateInvitation)
		protected.PUT("/invitations/me/:id", handlers.UpdateInvitation)
		protected.DELETE("/invitations/me/:id", handlers.DeleteInvitation)

		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
		protected.GET("/orders/:id", handlers.GetOrder)
		protected.POST("/orders/:id/cancel", handlers.CancelOrder)
		protected.GET("/templates/me", handlers.GetMyTemplates)
	}

	// ─── ADMIN ONLY ───────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.RoleRequired("admin"))
	{
		// Auth admin
		admin.POST("/auth/create", handlers.CreateAdmin)
		admin.POST("/auth/change/password", handlers.ChangePassword)

		// Users
		admin.GET("/users", handlers.GetUsers)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		// Templates
		admin.GET("/templates", handlers.GetTemplates)
		admin.GET("/templates/:id", handlers.GetTemplate)
		admin.POST("/templates", handlers.CreateTemplate)
		admin.PUT("/templates/:id", handlers.UpdateTemplate)
		admin.DELETE("/templates/:id", handlers.DeleteTemplate)

		// Invitations
		admin.GET("/invitations", handlers.GetInvitations)

		// Messages
		admin.GET("/messages", handlers.GetAllMessages)
		admin.DELETE("/messages/:id", handlers.DeleteMessage)

		// Orders
		admin.GET("/orders", handlers.GetAllOrders)
	}
}
