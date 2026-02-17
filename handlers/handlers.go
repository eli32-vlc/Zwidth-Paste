package handlers

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/eli32-vlc/Zwidth-Paste/database"
	"github.com/eli32-vlc/Zwidth-Paste/hashcash"
	"github.com/eli32-vlc/Zwidth-Paste/markdown"
	"github.com/eli32-vlc/Zwidth-Paste/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	DB            *database.DB
	SessionStore  *sessions.CookieStore
	MDRenderer    *markdown.Renderer
	AdminUsername string
	AdminPassword string
}

func NewHandler(db *database.DB, sessionSecret, adminUser, adminPass string) *Handler {
	return &Handler{
		DB:            db,
		SessionStore:  sessions.NewCookieStore([]byte(sessionSecret)),
		MDRenderer:    markdown.NewRenderer(),
		AdminUsername: adminUser,
		AdminPassword: adminPass,
	}
}

// HomePage displays the main editor
func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "home.html"),
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Zwidth Paste - Create New Entry",
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// CreateEntry handles entry creation
func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify hashcash
	hashcashValue := r.FormValue("hashcash")
	if err := hashcash.Verify(hashcashValue, "create", hashcash.DefaultBits); err != nil {
		http.Error(w, "Invalid hashcash: "+err.Error(), http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	customURL := strings.ToLower(strings.TrimSpace(r.FormValue("url")))
	editCode := strings.TrimSpace(r.FormValue("edit_code"))
	modifyCode := strings.TrimSpace(r.FormValue("modify_code"))

	// Validate content
	if len(content) == 0 {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}
	if len(content) > 200000 {
		http.Error(w, "Content exceeds 200,000 character limit", http.StatusBadRequest)
		return
	}

	// Generate or validate URL
	var url string
	if customURL != "" {
		if !utils.ValidateURL(customURL) {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		exists, err := h.DB.URLExists(customURL)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "URL already taken", http.StatusConflict)
			return
		}
		url = customURL
	} else {
		// Generate random URL
		for i := 0; i < 10; i++ {
			url = utils.GenerateRandomURL(8)
			exists, _ := h.DB.URLExists(url)
			if !exists {
				break
			}
		}
	}

	// Generate or validate edit code
	if editCode == "" {
		editCode = utils.GenerateRandomCode(16)
	} else if !utils.ValidateCode(editCode) {
		http.Error(w, "Invalid edit code", http.StatusBadRequest)
		return
	}

	// Validate modify code if provided
	if modifyCode != "" && !utils.ValidateCode(modifyCode) {
		http.Error(w, "Invalid modify code", http.StatusBadRequest)
		return
	}

	// Get client IP
	clientIP := utils.GetClientIP(r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-IP"))

	// Create entry
	entry := &database.Entry{
		URL:         url,
		Content:     content,
		EditCode:    editCode,
		ModifyCode:  modifyCode,
		CreatedByIP: clientIP,
	}

	if err := h.DB.CreateEntry(entry); err != nil {
		http.Error(w, "Failed to create entry: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success with edit code
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"url":       url,
		"edit_code": editCode,
	})
}

// ViewEntry displays an entry
func (h *Handler) ViewEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := strings.ToLower(vars["url"])

	entry, err := h.DB.GetEntryByURL(url)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Increment view count
	go func() {
		if err := h.DB.IncrementViewCount(entry.ID); err != nil {
			log.Printf("Failed to increment view count for entry %d: %v", entry.ID, err)
		}
	}()

	// Render markdown
	htmlContent, err := h.MDRenderer.Render(entry.Content)
	if err != nil {
		htmlContent = entry.Content // Fallback to raw content
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "view.html"),
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":       "View Entry - " + url,
		"URL":         url,
		"Content":     template.HTML(htmlContent),
		"RawContent":  entry.Content,
		"CreatedAt":   entry.CreatedAt.Format("2006-01-02 15:04:05"),
		"UpdatedAt":   entry.UpdatedAt.Format("2006-01-02 15:04:05"),
		"ViewCount":   entry.ViewCount,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// EditEntryPage displays the edit form
func (h *Handler) EditEntryPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := strings.ToLower(vars["url"])

	entry, err := h.DB.GetEntryByURL(url)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "edit.html"),
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Edit Entry - " + url,
		"URL":     url,
		"Content": entry.Content,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// UpdateEntry handles entry updates
func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	url := strings.ToLower(vars["url"])

	// Get existing entry
	entry, err := h.DB.GetEntryByURL(url)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Verify edit code
	providedCode := r.FormValue("edit_code")
	isEditCode := providedCode == entry.EditCode
	isModifyCode := entry.ModifyCode != "" && providedCode == entry.ModifyCode

	if !isEditCode && !isModifyCode {
		http.Error(w, "Invalid edit code", http.StatusForbidden)
		return
	}

	// Update content
	content := r.FormValue("content")
	if len(content) == 0 {
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}
	if len(content) > 200000 {
		http.Error(w, "Content exceeds 200,000 character limit", http.StatusBadRequest)
		return
	}

	entry.Content = content
	entry.LastEditedIP = utils.GetClientIP(r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-IP"))

	// Only allow full edit code holders to change URL and codes
	if isEditCode {
		newURL := strings.ToLower(strings.TrimSpace(r.FormValue("new_url")))
		newEditCode := strings.TrimSpace(r.FormValue("new_edit_code"))
		newModifyCode := strings.TrimSpace(r.FormValue("new_modify_code"))

		if newURL != "" && newURL != url {
			if !utils.ValidateURL(newURL) {
				http.Error(w, "Invalid URL format", http.StatusBadRequest)
				return
			}
			exists, _ := h.DB.URLExists(newURL)
			if exists {
				http.Error(w, "URL already taken", http.StatusConflict)
				return
			}
			entry.URL = newURL
		}

		if newEditCode != "" {
			if !utils.ValidateCode(newEditCode) {
				http.Error(w, "Invalid edit code", http.StatusBadRequest)
				return
			}
			entry.EditCode = newEditCode
		}

		if newModifyCode != "" {
			if !utils.ValidateCode(newModifyCode) {
				http.Error(w, "Invalid modify code", http.StatusBadRequest)
				return
			}
			entry.ModifyCode = newModifyCode
		}
	}

	if err := h.DB.UpdateEntry(entry); err != nil {
		http.Error(w, "Failed to update entry", http.StatusInternalServerError)
		return
	}

	// Redirect to the (possibly new) URL
	http.Redirect(w, r, "/"+entry.URL, http.StatusSeeOther)
}

// DeleteEntry handles entry deletion
func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	url := strings.ToLower(vars["url"])

	entry, err := h.DB.GetEntryByURL(url)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Verify edit code (only full edit code can delete)
	providedCode := r.FormValue("edit_code")
	if providedCode != entry.EditCode {
		http.Error(w, "Invalid edit code", http.StatusForbidden)
		return
	}

	if err := h.DB.DeleteEntry(entry.ID); err != nil {
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GetHashcashChallenge returns a hashcash challenge
func (h *Handler) GetHashcashChallenge(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		resource = "create"
	}

	challenge := hashcash.GenerateChallenge(resource)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"challenge": challenge,
	})
}

// RawEntry returns raw markdown content
func (h *Handler) RawEntry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := strings.ToLower(vars["url"])

	entry, err := h.DB.GetEntryByURL(url)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(entry.Content))
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("templates", "base.html"),
			filepath.Join("templates", "register.html"),
		)
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title": "Register - Zwidth Paste",
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// POST - Handle registration
	// Verify hashcash
	hashcashValue := r.FormValue("hashcash")
	if err := hashcash.Verify(hashcashValue, "register", hashcash.DefaultBits); err != nil {
		http.Error(w, "Invalid hashcash: "+err.Error(), http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if len(username) < 3 || len(username) > 50 {
		http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
		return
	}

	if len(password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if err := h.DB.CreateUser(username, hashedPassword); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("templates", "base.html"),
			filepath.Join("templates", "login.html"),
		)
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title": "Login - Zwidth Paste",
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// POST - Handle login
	// Verify hashcash
	hashcashValue := r.FormValue("hashcash")
	if err := hashcash.Verify(hashcashValue, "login", hashcash.DefaultBits); err != nil {
		http.Error(w, "Invalid hashcash: "+err.Error(), http.StatusBadRequest)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	user, err := h.DB.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, _ := h.SessionStore.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["username"] = user.Username
	session.Values["is_admin"] = user.IsAdmin
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, "session")
	session.Values["user_id"] = nil
	session.Values["username"] = nil
	session.Values["is_admin"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AdminLogin handles admin authentication
func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles(
			filepath.Join("templates", "base.html"),
			filepath.Join("templates", "admin_login.html"),
		)
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Title": "Admin Login - Zwidth Paste",
		}

		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// POST - Handle login
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	// Use constant-time comparison to prevent timing attacks
	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(h.AdminUsername)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(h.AdminPassword)) == 1

	if !usernameMatch || !passwordMatch {
		http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
		return
	}

	// Create admin session
	session, _ := h.SessionStore.Get(r, "admin_session")
	session.Values["admin_authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// AdminMiddleware checks if user is admin
func (h *Handler) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.SessionStore.Get(r, "admin_session")
		if auth, ok := session.Values["admin_authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// AdminDashboard displays the admin panel
func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.DB.GetStatistics()
	if err != nil {
		http.Error(w, "Error fetching statistics", http.StatusInternalServerError)
		return
	}

	// Get recent entries
	entries, err := h.DB.GetAllEntries(20, 0)
	if err != nil {
		http.Error(w, "Error fetching entries", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "admin.html"),
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":   "Admin Dashboard - Zwidth Paste",
		"Stats":   stats,
		"Entries": entries,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// AdminSearch handles entry search
func (h *Handler) AdminSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	entries, err := h.DB.SearchEntries(query, 50)
	if err != nil {
		http.Error(w, "Error searching entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

// AdminDeleteEntry handles entry deletion from admin panel
func (h *Handler) AdminDeleteEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	if err := h.DB.HardDeleteEntry(id); err != nil {
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// AboutPage displays the about page
func (h *Handler) AboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", "about.html"),
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "About - Zwidth Paste",
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
