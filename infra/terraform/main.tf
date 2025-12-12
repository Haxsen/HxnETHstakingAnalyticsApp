terraform {
  required_providers {
    render = {
      source  = "render-oss/render"
      version = "~> 1.0"
    }
  }
}

provider "render" {
  api_key  = var.render_api_key
  owner_id = var.render_owner_id
}

# PostgreSQL Database
resource "render_postgres" "database" {
  name     = "eth-staking-analytics-db"
  plan     = "free"
  region   = "frankfurt"
  version  = "16"
}

# Redis Cache
resource "render_redis" "cache" {
  name               = "eth-staking-analytics-redis"
  plan               = "free"
  region             = "frankfurt"
  max_memory_policy  = "allkeys_lru"
}

# Backend Web Service
resource "render_web_service" "backend" {
  name     = "eth-staking-analytics-backend"
  region   = "frankfurt"
  plan     = "free"

  runtime_source = {
    native_runtime = {
      runtime = "go"
      repo_url = "https://github.com/Haxsen/HxnETHstakingAnalyticsApp"
      branch = "main"
      build_command = "cd backend && go build -o app main.go"
    }
  }
  start_command = "./app"
}

# Frontend Static Site
resource "render_static_site" "frontend" {
  name         = "eth-staking-analytics-frontend"
  repo_url     = "https://github.com/Haxsen/HxnETHstakingAnalyticsApp"
  branch       = "main"
  build_command = "cd frontend && pnpm install && pnpm run build"
  publish_path = "frontend/out"
}
