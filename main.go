package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	sync.RWMutex
	users  map[string]User
	posts  map[int]Post
	nextID int
}

var store = Store{
	users:  make(map[string]User),
	posts:  make(map[int]Post),
	nextID: 1,
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	store.Lock()
	store.users[creds.Username] = User{Username: creds.Username, Password: string(hashedPassword)}
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	store.RLock()
	user, exists := store.users[creds.Username]
	store.RUnlock()

	if !exists || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	post.Author = username
	post.ID = store.nextID
	post.CreatedAt = time.Now()
	post.Comments = []Comment{}

	store.Lock()
	store.posts[store.nextID] = post
	store.nextID++
	store.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	store.RLock()
	posts := make([]Post, 0, len(store.posts))
	for _, post := range store.posts {
		posts = append(posts, post)
	}
	store.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	username := getUsernameFromRequest(w, r)
	if username == "" {
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	comment.Author = username
	comment.CreatedAt = time.Now()

	store.Lock()
	defer store.Unlock()
	post, exists := store.posts[postID]
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	comment.ID = len(post.Comments) + 1
	post.Comments = append(post.Comments, comment)
	store.posts[postID] = post

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func getUsernameFromRequest(w http.ResponseWriter, r *http.Request) string {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}

	tokenStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return ""
	}
	if !tkn.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return ""
	}
	return claims.Username
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.HandleFunc("/posts", CreatePost).Methods("POST")
	r.HandleFunc("/posts", GetPosts).Methods("GET")
	r.HandleFunc("/posts/{id}/comments", AddComment).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
