# 🚀 Database Fix Quick Start Guide

**For:** Developers implementing the database refactoring  
**Duration:** 6 weeks across 3 phases  
**Start Here:** This guide gets you started in 5 minutes

---

## 📚 What Was Created?

I've created a **complete implementation plan** with:

### **✅ Master Planning Documents:**
- [`MASTER-FIX-PLAN.md`](./MASTER-FIX-PLAN.md) - 6-week roadmap, resource allocation, KPIs
- [`DATABASE-ARCHITECTURE-REVIEW.md`](./DATABASE-ARCHITECTURE-REVIEW.md) - 14 issues identified
- [`ER-DIAGRAM.md`](./ER-DIAGRAM.md) - Current database structure

### **✅ Implementation Tickets:**
- [`fixes/phase1/T1.1-ticket-assignment.md`](./fixes/phase1/T1.1-ticket-assignment.md) - Detailed example ticket
- 11 more tickets (to be created for remaining issues)

### **✅ Migration Scripts:**
- [`migrations/phase1/001_ticket_assignment.up.sql`](./migrations/phase1/001_ticket_assignment.up.sql) - Forward migration
- [`migrations/phase1/001_ticket_assignment.down.sql`](./migrations/phase1/001_ticket_assignment.down.sql) - Rollback script

### **✅ Todo Tracking:**
- 12 tracked items across 3 phases (visible in your Factory todo list)

---

## 🎯 How to Use This Plan

### **Phase 1: Review (30 minutes)**

1. **Read Master Plan:**
   ```bash
   cat docs/database/MASTER-FIX-PLAN.md
   ```
   - Understand 3-phase approach
   - Review timeline (6 weeks)
   - Check resource requirements

2. **Review Example Ticket:**
   ```bash
   cat docs/database/fixes/phase1/T1.1-ticket-assignment.md
   ```
   - See detailed problem statement
   - Review SQL migrations
   - Check backend code changes
   - Look at frontend examples

3. **Check Database Review:**
   ```bash
   cat docs/database/DATABASE-ARCHITECTURE-REVIEW.md
   ```
   - Understand all 14 issues
   - Review impact analysis
   - See correct design solutions

---

### **Phase 2: Start Implementation (Week 1-2)**

#### **Step 1: Setup Environment**
```bash
# Ensure you're on latest code
cd /path/to/aby-med
git checkout -b database-refactor-phase1

# Backup database
pg_dump -U postgres -h localhost aby_med_db > backup_$(date +%Y%m%d).sql

# Test migrations on dev database
psql -U postgres -h localhost aby_med_db_test < docs/database/migrations/phase1/001_ticket_assignment.up.sql
```

#### **Step 2: Implement T1.1 (Ticket Assignment)**

**Day 1: Database Migration**
```bash
# Run migration
psql -U postgres aby_med_dev < docs/database/migrations/phase1/001_ticket_assignment.up.sql

# Verify results
psql -U postgres aby_med_dev -c "
  SELECT COUNT(*) FROM engineer_assignments;
  SELECT COUNT(DISTINCT ticket_id) FROM engineer_assignments;
"

# Test rollback
psql -U postgres aby_med_dev < docs/database/migrations/phase1/001_ticket_assignment.down.sql
psql -U postgres aby_med_dev < docs/database/migrations/phase1/001_ticket_assignment.up.sql
```

**Day 2: Backend Changes**
```bash
# Create new assignment domain file
touch internal/service-domain/service-ticket/domain/assignment.go
# Copy code from T1.1-ticket-assignment.md section "Backend Code Changes"

# Create repository implementation
touch internal/service-domain/service-ticket/infra/assignment_repository.go
# Copy code from ticket

# Update service layer
vim internal/service-domain/service-ticket/app/service.go
# Update AssignTicket, add EscalateTicket functions

# Write tests
go test ./internal/service-domain/service-ticket/...
```

**Day 3: Frontend Changes**
```bash
# Create API client
touch admin-ui/src/lib/api/assignments.ts
# Copy from ticket

# Create UI component
touch admin-ui/src/components/tickets/AssignmentHistory.tsx
# Copy from ticket

# Test in dev
cd admin-ui
npm run dev
```

**Day 4: Testing & Review**
```bash
# Run full test suite
go test ./... -v
cd admin-ui && npm test

# Manual testing checklist (from ticket):
# - Create ticket and assign engineer
# - Escalate ticket to next tier
# - View assignment history in UI
# - Verify backfilled data

# Code review
git add .
git commit -m "feat(tickets): Refactor assignment tracking to use engineer_assignments table"
git push origin database-refactor-phase1
# Create PR with ticket link
```

---

### **Phase 3: Track Progress**

#### **Use Factory Todo List:**
Your Factory AI assistant is tracking 12 items:

**Phase 1 (Week 1-2):**
- [ ] T1.1: Service ticket assignment refactor
- [ ] T1.2: Create customers table
- [ ] T1.3: Normalize RFQ/Quote items
- [ ] T1.4: Equipment relationships history

**Phase 2 (Week 3-4):**
- [ ] T2.1: Org relationship terms versioning
- [ ] T2.2: Standardize IDs to UUID
- [ ] T2.3: Normalize engineer coverage
- [ ] T2.4: Price rules temporal constraints

**Phase 3 (Week 5-6):**
- [ ] T3.1: Certification renewal tracking
- [ ] T3.2: Ticket status sync mechanism
- [ ] T3.3: Contact person history
- [ ] T3.4: Territory assignments many-to-many

#### **Update Status:**
Tell your Factory assistant:
- "Mark T1.1 as in_progress" → When you start
- "Mark T1.1 as completed" → When done
- "Show me todo list" → See current status

---

## 🔥 The Fastest Way to Start

### **Option A: Start with T1.1 Today (Recommended)**

This is the most critical fix with the clearest implementation path:

```bash
# 1. Run migration (5 minutes)
psql -U postgres aby_med_dev < docs/database/migrations/phase1/001_ticket_assignment.up.sql

# 2. Copy backend code (30 minutes)
# Follow T1.1-ticket-assignment.md section "Backend Code Changes"

# 3. Copy frontend code (30 minutes)
# Follow T1.1-ticket-assignment.md section "Frontend Code Changes"

# 4. Test (1 hour)
# Follow testing checklist in T1.1

# Total: ~2.5 hours for first fix!
```

### **Option B: Request Remaining Tickets**

Tell your Factory assistant:
> "Create detailed tickets for T1.2, T1.3, and T1.4 with SQL scripts and code examples"

Then implement all Phase 1 tickets in parallel.

### **Option C: Adjust Priorities**

If certain issues are more urgent for your business:
> "Move T1.4 (Equipment Relationships) to top priority because we need ownership tracking ASAP"

---

## 📊 What Each Ticket Contains

Every ticket has the same structure (T1.1 is the example):

1. **📋 Problem Statement** - What's broken and why
2. **🎯 Objective** - What success looks like
3. **📊 Current vs Target** - Schema changes visualized
4. **🗃️ SQL Migration** - Complete up/down scripts
5. **💻 Backend Code** - Go code with before/after
6. **🌐 Frontend Code** - React/TypeScript components
7. **✅ Testing Checklist** - Unit, integration, manual tests
8. **📊 Acceptance Criteria** - Definition of done
9. **🔄 Rollback Procedure** - How to revert safely
10. **📝 Documentation** - What needs updating

---

## 🚨 Important Principles

### **1. Zero Downtime Migrations**
Every migration supports:
- ✅ Dual-write pattern (old + new columns)
- ✅ Backward compatibility
- ✅ Safe rollback

### **2. Test Before Production**
Always:
```bash
# Test on dev database first
psql -U postgres aby_med_dev < migration.up.sql

# Test rollback
psql -U postgres aby_med_dev < migration.down.sql

# Test again
psql -U postgres aby_med_dev < migration.up.sql
```

### **3. Incremental Changes**
- Each ticket is independently deployable
- Don't wait for all Phase 1 to be complete
- Ship T1.1, then T1.2, etc.

---

## 🆘 Need Help?

### **Ask Your Factory Assistant:**
- "Show me the SQL for T1.2 (customers table)"
- "Generate backend code for equipment relationships"
- "Create frontend component for org terms history"
- "What's the status of Phase 1?"

### **Review Documents:**
- **Issues unclear?** → Read `DATABASE-ARCHITECTURE-REVIEW.md`
- **Need timeline?** → Check `MASTER-FIX-PLAN.md`
- **Want example?** → Study `T1.1-ticket-assignment.md`

### **Common Questions:**

**Q: Do I need to implement all 12 tickets?**  
A: Start with Phase 1 (4 tickets). Re-assess after 2 weeks.

**Q: Can I change the priorities?**  
A: Yes! Tell your assistant: "I want to prioritize X over Y because..."

**Q: What if migration fails?**  
A: Run the `.down.sql` rollback script immediately. All migrations are reversible.

**Q: Can I deploy T1.1 without waiting for T1.2?**  
A: Absolutely! Each ticket is independent.

---

## ✅ Quick Checklist: Am I Ready?

Before starting implementation:

- [ ] I've read the Master Plan
- [ ] I understand the 3-phase approach
- [ ] I've reviewed T1.1 example ticket
- [ ] I have a backed-up dev database
- [ ] I know how to run PostgreSQL migrations
- [ ] I can write Go code for backend
- [ ] I can update React/TypeScript frontend
- [ ] I've marked T1.1 as "in_progress" in my todo list

---

## 🎉 Success Metrics

**After Phase 1 (Week 2):**
- ✅ Ticket escalation tracking works (L1 → L2 → L3)
- ✅ Customer data normalized (no duplicates)
- ✅ RFQ/Quote items fully queryable
- ✅ Equipment ownership changes tracked

**After Phase 2 (Week 4):**
- ✅ Historical commission calculations accurate
- ✅ All IDs standardized to UUID
- ✅ Engineer coverage queries < 50ms
- ✅ Price history works for billing

**After Phase 3 (Week 6):**
- ✅ All 14 issues resolved
- ✅ System audit-ready for compliance
- ✅ Query performance excellent
- ✅ Zero data integrity issues

---

## 🚀 Let's Go!

**Recommended First Step:**  
Tell your Factory assistant:
> "Let's start implementing T1.1 (Service Ticket Assignment). Show me the migration script."

Or:
> "Create detailed tickets for all of Phase 1 (T1.2, T1.3, T1.4)"

---

**Questions?** Ask your Factory AI assistant anytime!  
**Blocked?** Review the relevant ticket documentation.  
**Unsure?** Check the Database Architecture Review for the "why" behind each fix.

**Good luck! You've got this! 💪**
