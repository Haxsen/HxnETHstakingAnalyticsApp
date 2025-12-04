output "backend_service_url" {
  description = "Backend service URL"
  value       = render_web_service.backend.url
}

output "frontend_site_url" {
  description = "Frontend static site URL"
  value       = render_static_site.frontend.url
}
