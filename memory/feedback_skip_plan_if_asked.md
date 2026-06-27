---
name: skip-plan-if-asked
description: User wants to skip writing a plan and jump straight to implementation when they say so
metadata:
  type: feedback
---

If the user explicitly says to skip writing a plan and jump to implementation, do so — don't invoke the planner agent or write a plan file.

**Why:** User prefers to iterate on code directly rather than read an intermediate plan doc.

**How to apply:** When the user says "jump straight to implementation" or "no plan needed" or similar, go directly to code changes.
