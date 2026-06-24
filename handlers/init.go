package handlers

import (
	"invitation-app/services"
	"sync"
)

var (
	_xenditOnce          sync.Once
	_emailOnce           sync.Once
	_templateAccessOnce  sync.Once

	_xenditService         *services.XenditService
	_emailService          *services.EmailService
	_templateAccessService *services.TemplateAccessService
)

func xenditSvc() *services.XenditService {
	_xenditOnce.Do(func() {
		_xenditService = services.NewXenditService()
	})
	return _xenditService
}

func emailSvc() *services.EmailService {
	_emailOnce.Do(func() {
		_emailService = services.NewEmailService()
	})
	return _emailService
}

func templateAccessSvc() *services.TemplateAccessService {
	_templateAccessOnce.Do(func() {
		_templateAccessService = services.NewTemplateAccessService()
	})
	return _templateAccessService
}