# Infrastructure

This directory contains the Terraform configuration for deploying the HxnETHstakingAnalyticsApp to Render.

## Overview

The infrastructure consists of:

- **PostgreSQL Database**: Free tier database for storing token data
- **Redis Cache**: Free tier key-value store for API caching
- **Backend Web Service**: Go API server with REST endpoints (Hobby plan)
- **Frontend Static Site**: Next.js app served as static site (Free plan)

## Prerequisites

- Render account with API access
- Render API key (get from Render dashboard → Account → API)
- Render username/owner ID
- Environment variables configured (see .env.example)

## Files

- `main.tf`: Main Terraform configuration defining the Render resources
- `variables.tf`: Input variables (API key and owner ID)
- `outputs.tf`: Output values (service URLs and names)
- `terraform.tfvars`: Variable values (add your actual credentials here)
- `.env.example`: Example configuration reference

## Deployment

1. Add your Render credentials to `terraform.tfvars`:
   ```hcl
   render_api_key  = "your_actual_api_key_here"
   render_owner_id = "your_render_username"
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Plan the deployment:
   ```bash
   terraform plan
   ```

4. Apply the infrastructure:
   ```bash
   terraform apply
   ```

## Post-Deployment Setup

After Terraform creates the resources, manually set environment variables in the Render dashboard:

### Backend Service Environment Variables:
- `DATABASE_URL`: PostgreSQL connection string (from Render dashboard)
- `REDIS_URL`: Redis connection string (from Render dashboard)
- `COINGECKO_API_KEY`: Your CoinGecko API key
- `ETHEREUM_RPC_URL`: Ethereum RPC endpoint (default: https://ethereum-rpc.publicnode.com)
- `PORT`: 10000 (Render default)

### Database Setup:
After deployment, run the database schema:
```bash
# Connect to your Render PostgreSQL database and run:
psql -h [your-db-host] -U [your-db-user] -d [your-db-name] < backend/schema.sql
```

## Notes

- The database connection string is not exported by the Terraform provider, so it must be copied from the Render dashboard
- Free tier resources have limitations (1GB storage, may expire after 30 days)
- Services auto-deploy on git pushes to main branch

## Architecture

```
GitHub Repo
    ↓
Terraform → Render Services
    ↓
PostgreSQL + Redis + Backend API + Frontend Site
