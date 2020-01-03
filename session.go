/* *************************************************************
http extensions for Go web server

This package contains extensions to facilitate the use of
GO http server:

- A wrapper for for the ResponseWriter.
- A session manager.
************************************************************** */

// The httpext package contains utilities to facilitate
// the development of web application.
// - A session manager
// - An object to encapsulate all the HTTP Context
// - An easier prototype for handler functions

package expressgo

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Default session timeout (May be overwritten individually for each session)
//var SessionTimeout int = 600

// Frequency at which the session store removes expired sessions
//var CleanupFrequency int = 60
var serverSessions map[string]*HttpSession
var sessionIndex int
var sessionCleanerKeepActive = true
var sessionStoreLock sync.Mutex

// The session object
type HttpSession struct {
	Id      string
	LastUse time.Time
	Timeout int
	Values  map[string]string
}

// Sets a session variable
func (s *HttpSession) Set(key string, value string) {
	sessionStoreLock.Lock()
	s.Values[key] = value
	sessionStoreLock.Unlock()
}

// Gets a session variable
func (s *HttpSession) Get(key string) string {
	sessionStoreLock.Lock()
	var rv string
	var found bool
	if rv, found = s.Values[key]; !found {
		rv = ""
	}
	sessionStoreLock.Unlock()
	return rv
}

// Deletes a session variable
func (s *HttpSession) Delete(key string) {
	sessionStoreLock.Lock()
	delete(s.Values, key)
	sessionStoreLock.Unlock()
}

// Clears all variables from a session
func (s *HttpSession) Clear() {
	sessionStoreLock.Lock()
	s.Values = make(map[string]string)
	sessionStoreLock.Unlock()
}

func (s *HttpSession) init(sessionId string) {
	s.Timeout = config.Timeout
	s.Id = sessionId
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

func getHttpSession(writer http.ResponseWriter, req *http.Request) *HttpSession {
	var session *HttpSession
	var found bool
	var sessionId string

	sessionStoreLock.Lock() // Prevent unikely parallel initialization

	if serverSessions == nil {
		LogDebug("Init Session Manager")
		serverSessions = make(map[string]*HttpSession)
		sessionIndex = 12345
		go sessionStoreCleaner()
	}

	if c, e := req.Cookie("GOSVR_SessionId"); e == nil {
		sessionId = c.Value
	} else {
		sessionId = ""
	}

	LogDebug(fmt.Sprintf("Read session(id:%s)\n", sessionId))

	if session, found = serverSessions[sessionId]; !found || sessionId == "" {
		sessionIndex = sessionIndex + 1

		sessionId = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d%d", sessionIndex, time.Now().UnixNano())))

		LogDebug(fmt.Sprintf("Start new session(id:%s)\n", sessionId))
		session = new(HttpSession)
		session.init(sessionId)
		serverSessions[sessionId] = session
	} else {
		if time.Now().After(session.LastUse.Add(time.Duration(session.Timeout) * time.Second)) {
			LogDebug(fmt.Sprintf("Session expired. Creating new one(id:%s)\n", sessionId))
			session = new(HttpSession)
			session.init(sessionId)
			serverSessions[sessionId] = session
		}
	}

	session.LastUse = time.Now()

	sessionStoreLock.Unlock()

	c := http.Cookie{Name: "GOSVR_SessionId", Value: sessionId, Path: "/"}
	writer.Header().Add("Set-cookie", c.String())

	return session
}

type SessionConfig struct {
	Timeout         int
	CleanupInterval int
}

var config SessionConfig

func Session(conf SessionConfig) func(req *Request, resp *Response) Status {
	config = conf
	return func(req *Request, resp *Response) Status {
		req.Vars["x_session"] = getHttpSession(resp.ResponseWriter, req.Request)
		return resp.OK()
	}
}
