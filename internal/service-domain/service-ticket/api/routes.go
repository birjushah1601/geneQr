package api

import (
	"github.com/go-chi/chi/v5"
)

// RegisterAssignmentRoutes registers all assignment-related routes
func RegisterAssignmentRoutes(r chi.Router, handler *AssignmentHandler) {
	// Ticket-based assignment endpoints
	r.Route("/tickets/{ticketId}", func(r chi.Router) {
		r.Post("/assign", handler.AssignTicket)
		r.Post("/escalate", handler.EscalateTicket)
		r.Get("/current-assignment", handler.GetCurrentAssignment)
		r.Get("/assignments", handler.GetAssignmentHistory)
	})

	// Assignment-based endpoints
	r.Route("/assignments/{assignmentId}", func(r chi.Router) {
		r.Post("/accept", handler.AcceptAssignment)
		r.Post("/reject", handler.RejectAssignment)
		r.Post("/start", handler.StartAssignment)
		r.Post("/complete", handler.CompleteAssignment)
		r.Post("/feedback", handler.AddCustomerFeedback)
	})

	// Engineer-based endpoints
	r.Route("/engineers/{engineerId}", func(r chi.Router) {
		r.Get("/assignments", handler.GetEngineerAssignments)
		r.Get("/assignments/active", handler.GetActiveEngineerAssignments)
		r.Get("/workload", handler.GetEngineerWorkload)
	})
}
