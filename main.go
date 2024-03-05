package main

import (
  "fmt"
  "net/http"
  "sync"
  "crypto/rand"
  "encoding/base64"
)

var (
  // Map for storing session information
  sessions = make(map[string]bool)

  // Mutex to synchronize access to session informstion
  sessionMutex = &sync.Mutex{}
)

const sessionIDLength = 32

func generateSessionID() (string, error) {
  // Create a byte slice to hold random bytes
  randomBytes := make([]byte, sessionIDLength)
  // Generate random byte using crypto/rand
  _, err := rand.Read(randomBytes)
  if err != nil {
    return "", err
  }

  // Encode the randomBytes using base64 
  sessionID := base64.URLEncoding.EncodeToString(randomBytes)

  return sessionID, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request)  {
  // Normally, implement user authentication logic with ID and passowrd here

  // Generete a session using ID (here, simply using a username)
  sessionID, err := generateSessionID()
  if err !=  nil {
    fmt.Printf("Error generation sesson ID : %q", err)
  }

  // Save the session on the server side
  sessionMutex.Lock()
  sessions[sessionID] = true

  for id := range sessions {
    fmt.Println("loginHandler: Current session ID is ", id)
  }
  sessionMutex.Unlock()

  // Send the session ID to this client as cookie
  http.SetCookie(w, &http.Cookie{
    Name: "session_id",
    Value: sessionID,
    Path: "/",
    MaxAge: 60,  // set the expiration time to 60 seconds (1 minute)
  })

  w.Write([]byte("Login succesfully!"))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
  // Retrive session DI from Cookies
  cookie, err := r.Cookie("session_id")
  if err != nil {
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
  }

  sessionMutex.Lock()
  authenticated, ok := sessions[cookie.Value]
  sessionMutex.Unlock()

  // Check if the session is Valid
  if !ok || !authenticated {
      http.Error(w, "Forbidden", http.StatusForbidden)
      return
  }

  w.Write([]byte("Welcome to your dashboard"))}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
  // Retrive session ID from cookies and delete the session onthe server side 
  cookie, err := r.Cookie("session_id")
  if err == nil {
    sessionMutex.Lock()
    delete(sessions, cookie.Value)
    if len(sessions) == 0 {
      fmt.Println("logoutHandler: current sessions ID is nothing")
    }
    sessionMutex.Unlock()
  }

  http.SetCookie(w, &http.Cookie{
    Name: "session_id",
    Value: "",
    Path: "/",
    MaxAge: -1,
  })

  w.Write([]byte("Logout succsessful!"))
}

func main() {
  http.HandleFunc("/login", loginHandler)
  http.HandleFunc("/dashboard", dashboardHandler)
  http.HandleFunc("/logout", logoutHandler)

  http.ListenAndServe(":6969", nil)
}
