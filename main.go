package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eli32-vlc/Zwidth-Paste/database"
	"github.com/eli32-vlc/Zwidth-Paste/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables (or use defaults)
	dbPath := getEnv("DB_PATH", "./zwidth.db")
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")
	sessionSecret := getEnv("SESSION_SECRET", "change-this-secret-key")
	adminUsername := getEnv("ADMIN_USERNAME", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "changeme")

	// Initialize database
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database initialized successfully")

	// Initialize handlers
	h := handlers.NewHandler(db, sessionSecret, adminUsername, adminPassword)

	// Set up router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Public routes
	r.HandleFunc("/", h.HomePage).Methods("GET")
	r.HandleFunc("/create", h.CreateEntry).Methods("POST")
	r.HandleFunc("/hashcash", h.GetHashcashChallenge).Methods("GET")
	r.HandleFunc("/register", h.Register).Methods("GET", "POST")
	r.HandleFunc("/login", h.Login).Methods("GET", "POST")
	r.HandleFunc("/logout", h.Logout).Methods("GET")
	r.HandleFunc("/about", h.AboutPage).Methods("GET")

	// Admin routes
	r.HandleFunc("/admin/login", h.AdminLogin).Methods("GET", "POST")
	r.HandleFunc("/admin", h.AdminMiddleware(h.AdminDashboard)).Methods("GET")
	r.HandleFunc("/admin/search", h.AdminMiddleware(h.AdminSearch)).Methods("GET")
	r.HandleFunc("/admin/delete", h.AdminMiddleware(h.AdminDeleteEntry)).Methods("POST")

	// Entry routes
	r.HandleFunc("/raw/{url}", h.RawEntry).Methods("GET")
	r.HandleFunc("/edit/{url}", h.EditEntryPage).Methods("GET")
	r.HandleFunc("/edit/{url}", h.UpdateEntry).Methods("POST")
	r.HandleFunc("/delete/{url}", h.DeleteEntry).Methods("POST")
	r.HandleFunc("/{url}", h.ViewEntry).Methods("GET")

	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("Starting server on http://%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
