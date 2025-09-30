# Medical Equipment Procurement Platform - Architecture

## Overview
A complete medical equipment procurement platform with RFQ management, supplier profiles, quote comparison, and contract award capabilities.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        React Frontend                            │
│  (RFQ Management, Quote Comparison, Supplier Portal, Dashboard) │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                  ┌─────────▼──────────┐
                  │   API Gateway      │
                  │   (Port 8081)      │
                  └─────────┬──────────┘
                            │
          ┌─────────────────┼─────────────────┐
          │                 │                 │
┌─────────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
│  RFQ Service   │  │  Supplier   │  │   Quote     │
│                │  │  Service    │  │   Service   │
└────────┬───────┘  └──────┬──────┘  └──────┬──────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
┌─────────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
│  Comparison    │  │  Contract   │  │   Catalog   │
│  Service       │  │  Service    │  │   Service   │
└────────┬───────┘  └──────┬──────┘  └──────┬──────┘
         │                 │                 │
         └─────────────────┼─────────────────┘
                           │
                  ┌────────▼─────────┐
                  │  Event Bus       │
                  │  (Kafka)         │
                  └──────────────────┘
```

## Service Responsibilities

### 1. **RFQ Service** ✅ (Completed)
- Create and manage RFQs
- Define requirements and specifications
- Set deadlines and terms
- Publish RFQs to suppliers
- Manage RFQ lifecycle (draft → published → closed → awarded)

### 2. **Supplier Service** (Next)
- Supplier registration and profiles
- Certifications and compliance documents
- Performance ratings and history
- Category specializations
- Contact information and verification

### 3. **Quote Service** (After Supplier)
- Submit quotes in response to RFQs
- Line-item pricing
- Quote revisions
- Terms and conditions
- Attachments and supporting documents
- Quote lifecycle (draft → submitted → accepted/rejected)

### 4. **Comparison Service** (After Quote)
- Compare multiple quotes
- Score quotes based on criteria
- Generate comparison matrices
- Recommendation engine

### 5. **Contract Service** (After Comparison)
- Generate contracts from awarded quotes
- Contract templates
- Contract signing workflow
- Milestone tracking
