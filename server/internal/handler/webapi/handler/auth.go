// handler パッケージ: Web フロントエンド向け認証 API
// app/api/login.ts / logout.ts / me.ts に相当
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"orlder-car-line/server/internal/domain/auth"
	"orlder-car-line/server/internal/domain/session"
)

// HandleLogin POST /api/login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "リクエストの形式が不正です"})
		return
	}

	users := auth.ParseUsers()
	if body.Username == "" || body.Password == "" || users[body.Username] != body.Password {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "ユーザー名またはパスワードが正しくありません"})
		return
	}

	token, err := session.Create(r.Context(), body.Username)
	if err != nil {
		log.Printf("[login] セッション作成エラー: %v", err)
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "内部エラー"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(session.TTL.Seconds()),
		SameSite: http.SameSiteLaxMode,
	})
	respondJSON(w, http.StatusOK, map[string]string{"username": body.Username})
}

// HandleLogout POST /api/logout
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		if delErr := session.Delete(r.Context(), cookie.Value); delErr != nil {
			log.Printf("[logout] セッション削除エラー: %v", delErr)
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	w.WriteHeader(http.StatusNoContent)
}

// HandleMe GET /api/me
func HandleMe(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		if cookie, err := r.Cookie("session"); err == nil {
			authHeader = "Bearer " + cookie.Value
		}
	}

	result := auth.ValidateHeader(r.Context(), authHeader)
	if !result.Valid {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": result.Err})
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"username": result.Username})
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("[respondJSON] encode error: %v", err)
	}
}