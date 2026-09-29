/*
 * Copyright 2026 Jonas Kaninda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package manifest

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// RouteSpec exposes an application on the public hostname of a URL the installer asked for. The route
// is best-effort: a host that is empty, not a public name, or under no domain of the workspace is
// reported and skipped, never failing the install.
type RouteSpec struct {
	// App names the application the route serves; it must declare a port.
	App string `yaml:"app" json:"app"`
	// Host renders to a URL or a bare hostname; only the hostname is kept. It must come from an input.
	Host string `yaml:"host" json:"host"`
	// Port is the container port to route to; empty means the application's first port.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`
}

func (m *Manifest) validateRoutes() error {
	apps := map[string]AppSpec{}
	for _, a := range m.Applications {
		apps[a.Name] = a
	}
	seen := map[string]bool{}
	for i, r := range m.Routes {
		a, ok := apps[r.App]
		if !ok {
			return fmt.Errorf("routes[%d]: unknown application %q", i, r.App)
		}
		if len(a.Ports) == 0 {
			return fmt.Errorf("routes[%d]: application %q declares no port to route to", i, r.App)
		}
		if r.Port != 0 && !hasPort(a, r.Port) {
			return fmt.Errorf("routes[%d]: application %q does not declare port %d", i, r.App, r.Port)
		}
		// A literal hostname would let a template claim a name the installer never chose.
		if !strings.Contains(r.Host, ".inputs.") {
			return fmt.Errorf("routes[%d]: host must be taken from an input, e.g. \"{{ .inputs.public_url }}\"", i)
		}
		key := r.App + "|" + r.Host
		if seen[key] {
			return fmt.Errorf("routes[%d]: duplicate route for %q", i, r.App)
		}
		seen[key] = true
	}
	return nil
}

func hasPort(a AppSpec, port int) bool {
	for _, p := range a.Ports {
		if p.Container == port {
			return true
		}
	}
	return false
}

// RouteHost reduces a rendered URL or host to the bare, lowercase hostname a route matches. It returns
// "" for a value that names no public host: empty, localhost, an IP address, or a single label.
func RouteHost(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.Contains(v, "://") {
		v = "http://" + v
	}
	u, err := url.Parse(v)
	if err != nil {
		return ""
	}
	h := strings.ToLower(strings.TrimSuffix(u.Hostname(), "."))
	if h == "" || h == "localhost" || net.ParseIP(h) != nil || !strings.Contains(h, ".") {
		return ""
	}
	return h
}
