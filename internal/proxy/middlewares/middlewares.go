package middlewares

import "fmt"

func RedirectToHTTPSLabels(namespace string) map[string]string {
	labels := make(map[string]string, 0)

	labels[fmt.Sprintf("traefik.http.routers.%s-insecure.middlewares", namespace)] = "redirect-to-https@docker"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.scheme"] = "https"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.port"] = "443"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.permanent"] = "true"

	return labels
}