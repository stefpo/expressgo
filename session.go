package expressgo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var serverSessions map[string]*HTTPSession
var sessionIndex int
var sessionCleanerKeepActive = true
var sessionStoreLock sync.Mutex

// HTTPSession structure contains the session data and control information
type HTTPSession struct {
	ID      string
	LastUse time.Time
	Timeout int
	Values  map[string]string
}

// Set Sets a session string variable
func (s *HTTPSession) Set(key string, value string) {
	sessionStoreLock.Lock()
	s.Values[key] = value
	sessionStoreLock.Unlock()
}

// Get Gets a session variable by name. It returns an empty string if the value doesn ot exist
func (s *HTTPSession) Get(key string) string {
	sessionStoreLock.Lock()
	var rv string
	var found bool
	if rv, found = s.Values[key]; !found {
		rv = ""
	}
	sessionStoreLock.Unlock()
	return rv
}

// Delete Deletes a session variable by name.
func (s *HTTPSession) Delete(key string) {
	sessionStoreLock.Lock()
	delete(s.Values, key)
	sessionStoreLock.Unlock()
}

// Clear removes all variables from current session
func (s *HTTPSession) Clear() {
	sessionStoreLock.Lock()
	s.Values = make(map[string]string)
	sessionStoreLock.Unlock()
}

func (s *HTTPSession) init(sessionID string) {
	s.Timeout = config.Timeout
	s.ID = sessionID
	s.Values = make(map[string]string)
}

func sessionStoreCleaner() {
	for sessionCleanerKeepActive {
		time.Sleep(time.Duration(config.CleanupInterval) * time.Second)
		sessionStoreLock.Lock()
		LogDebug("Expired session cleanup")
		for k, v := range serverSessions {
			exp := v.LastUse.Add(time.Duration(v.Timeout) * time.Second)
			if time.Now().After(exp) {
				LogDebug("Session " + k + "Last used on" + v.LastUse.String() + " expired on " + exp.String())
				delete(serverSessions, k)
			}
		}
		sessionStoreLock.Unlock()
	}
}

func getHTTPSession(writer http.ResponseWriter, req *http.Request) *HTTPSession {
	var session *HTTPSession
	var found bool
	var sessionID string

	sessionStoreLock.Lock() // Prevent unikely parallel initialization

	if serverSessions == nil {
		LogDebug("Init Session Manager")
		serverSessions = make(map[string]*HTTPSession)
		sessionIndex = 12345
		go sessionStoreCleaner()
	}

	if c, e := req.Cookie("XprGo-Session-Id"); e == nil {
		sessionID = c.Value
	} else {
		sessionID = ""
	}

	LogDebug(fmt.Sprintf("Read session(id:%s)\n", sessionID))

	if session, found = serverSessions[sessionID]; !found || sessionID == "" {
		sessionIndex = sessionIndex + 1

		sessionID = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%d", sessionIndex, time.Now().UnixNano())))

		LogDebug(fmt.Sprintf("Start new session(id:%s)\n", sessionID))
		session = new(HTTPSession)
		session.init(sessionID)
		serverSessions[sessionID] = session
	} else {
		if time.Now().After(session.LastUse.Add(time.Duration(session.Timeout) * time.Second)) {
			LogDebug(fmt.Sprintf("Session expired. Creating new one(id:%s)\n", sessionID))
			session = new(HTTPSession)
			session.init(sessionID)
			serverSessions[sessionID] = session
		}
	}

	session.LastUse = time.Now()

	sessionStoreLock.Unlock()

	c := http.Cookie{Name: "XprGo-Session-Id", Value: sessionID, Path: "/"}
	writer.Header().Add("Set-cookie", c.String())

	return session
}

// SessionConfig defines session manager parameters
type SessionConfig struct {
	Timeout         int
	CleanupInterval int
}

var config SessionConfig

// Session is the session middleware generator.
func Session(conf SessionConfig) func(*Request, *Response, func(...Error)) {
	config = conf
	return func(req *Request, resp *Response, next func(...Error)) {
		req.Vars["Session"] = getHTTPSession(resp.writer, req.Request)
		next()
	}
}
