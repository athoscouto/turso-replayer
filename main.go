package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

var (
	TURSO_DOMAIN = ".turso.io"
)

var ipv4Dbs = make(map[string]bool)

func main() {
	runHttpProxy()
}

func runServer(name string, handler http.Handler, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v", addr, err)
		return
	}

	log.Printf("%v listening on %v", name, listener.Addr())
	server := http.Server{Handler: handler}
	err = server.Serve(listener)
	log.Fatalf("%v failed: %v", name, err)
}

func runHttpProxy() {
	handler := http.HandlerFunc(handleHttp)
	runServer("HTTP proxy", handler, ":8080")
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	route, ok := parseRoute(r.Host)
	if !ok {
		http.Error(w, "Invalid host", http.StatusBadRequest)
		return
	}

	w.Header().Set("fly-replay", flyReplayHeader(route))
	w.WriteHeader(200)
}

type Route struct {
	app     string
	machine string
}

func dbToApp(db string) string {
	return "turso-" + strings.ToLower(db)
}

func parseRoute(hostHeader string) (Route, bool) {
	hostname, err := parseHostname(hostHeader)
	if err != nil {
		return Route{}, false
	}

	if !strings.HasSuffix(hostname, TURSO_DOMAIN) {
		return Route{}, false
	}

	segment := strings.TrimSuffix(hostname, TURSO_DOMAIN)
	if len(segment) == 0 {
		return Route{}, false
	}

	machine, db, ok := strings.Cut(segment, "-")
	if ok && looksLikeMachineId(machine) {
		return Route{app: dbToApp(db), machine: machine}, true
	}
	return Route{app: dbToApp(segment)}, true
}

func parseHostname(hostHeader string) (string, error) {
	if strings.Contains(hostHeader, ":") {
		hostname, _, err := net.SplitHostPort(hostHeader)
		return hostname, err
	} else {
		return hostHeader, nil
	}
}

func internalAddr(route Route) string {
	if route.machine == "" {
		return fmt.Sprintf("%s.internal", route.app)
	}
	return fmt.Sprintf("%s.vm.%s.internal", route.machine, route.app)
}

func flyReplayHeader(route Route) string {
	if route.machine == "" {
		return fmt.Sprintf("app=%s", route.app)
	}
	return fmt.Sprintf("app=%s;instance=%s", route.app, route.machine)
}

func looksLikeMachineId(token string) bool {
	if len(token) != 14 {
		return false
	}
	for _, r := range token {
		if !(isDigit(r) || isABCDEF(r)) {
			return false
		}
	}
	return true
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isABCDEF(r rune) bool {
	return 'a' <= r && r <= 'f'
}
