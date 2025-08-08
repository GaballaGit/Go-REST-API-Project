package middlewares

import (
	"net/http"
)

func SecurityHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off") // Reduces DNS leaks/traffic when being preresolved in an <h ref>

		w.Header().Set("X-Frame-Options", "DENY")                                                 // Prevents site from being displayed in <iframe>. This blocks clickjacking attacks
		w.Header().Set("X-XSS-Protection", "1:mode=block")                                        // Tells legacy browsers to include XSS filtering
		w.Header().Set("X-Content-Type-Options", "nosniff")                                       //Stops content type confusion (e.g a browser treats a js file as html)
		w.Header().Set("Strict Transport Security", "max-age=63072000;includeSubDomains;preload") // FORCES HTTPS. prevents SSL strip attacks
		w.Header().Set("Content-Security-Policy", "default-src `self`")                           //Restricts where content can be loaded. Here it can only be loaded in my own domain
		w.Header().Set("Referrer-Policy", "no-referrer")                                          // Prevents leaking path/query data to other external sites
		next.ServeHTTP(w, r)
	})

}
