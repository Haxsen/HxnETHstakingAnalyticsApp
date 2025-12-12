output "backend_service_url" {
  description = "Backend service URL"
  value       = render_web_service.backend.url
}

output "frontend_site_url" {
  description = "Frontend static site URL"
  value       = render_static_site.frontend.url
}

output "database_name" {
  description = "PostgreSQL database name"
  value       = render_postgres.database.name
}

output "redis_name" {
  description = "Redis cache name"
  value       = render_redis.cache.name
}
