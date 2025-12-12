variable "render_api_key" {
  description = "Render API key for authentication"
  type        = string
  sensitive   = true
}

variable "render_owner_id" {
  description = "Render owner ID (username)"
  type        = string
}
