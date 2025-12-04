terraform {
  required_providers {
    render = {
      source  = "render-oss/render"
      version = "~> 1.0"
    }
  }
}

provider "render" {
  api_key = var.render_api_key
}

# PostgreSQL Database
resource "render_postgres" "database" {
  name     = "eth-staking-analytics-db"
  plan     = "free"
  region   = "oregon"
  version  = "16"
}

# Backend Web Service
resource "render_web_service" "backend" {
  name     = "eth-staking-analytics-backend"
  region   = "oregon"
  plan     = "free"

  runtime_source = {
    native_runtime = {
      runtime = "node"
      repo_url = "https://github.com/Haxsen/HxnETHstakingAnalyticsApp"
      branch = "main"
      build_command = "cd backend && pnpm install"
    }
  }
  start_command = "cd backend && pnpm start"
}

# Frontend Static Site
resource "render_static_site" "frontend" {
  name         = "eth-staking-analytics-frontend"
  repo_url     = "https://github.com/Haxsen/HxnETHstakingAnalyticsApp"
  branch       = "main"
  build_command = "cd frontend && pnpm install && pnpm run build"
  publish_path = "frontend/dist"
}
