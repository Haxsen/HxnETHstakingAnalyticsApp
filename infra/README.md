# Infrastructure

This directory contains the Terraform configuration for deploying the HxnETHstakingAnalyticsApp to Render.

## Overview

The infrastructure consists of:

- **PostgreSQL Database**: Free tier database for storing token data, events, and snapshots
- **Backend Web Service**: Node.js/TypeScript API server with indexer
- **Frontend Static Site**: React SPA served as static site

## Prerequisites

- Render account with API access
- Render API key (get from Render dashboard → Account → API)
- Environment variables configured (see .env.example)

## Files

- `main.tf`: Main Terraform configuration defining the Render resources
- `variables.tf`: Input variables (API key)
- `outputs.tf`: Output values (service URLs)
- `terraform.tfvars`: Variable values (add your actual Render API key here)
- `.env.example`: Example Render API key for reference

## Deployment

1. Add your Render API key to `terraform.tfvars`:
   ```hcl
   render_api_key = "your_actual_api_key_here"
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

- **Backend Service**: Add `DATABASE_URL` with the database connection string from the PostgreSQL service
- **Backend Service**: Add `NODE_ENV=production`
- **Backend Service**: Add any other required env vars (CoinGecko API key, RPC endpoints, etc.)

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
PostgreSQL + Backend API + Frontend Site
