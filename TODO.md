# Backend

## Core
- [ ] Refresh token handling - caching user, how do we manage cache memory? What if there are many users?
- [ ] Admin role for entity creation

## Features
- [ ] Historical leaderboard snapshots (point-in-time rank history).
- [ ] Pagination of leaderboard results.
- [ ] Fixtures in group have group, but what about quarter finals, etc?
- [ ] *Locked fixtures - lock by worker n minutes in advance
- [ ] *Same points breaking criteria
- [ ] *Default values for missing handicaps
- [ ] Worker should fetch fixtures every day

### Profile page
- [ ] Change display name

---

# Mobile

## Core
- [ ] Share button: share link to entity, or deep link to app with entity preloaded
- [x] Auto-join on valid code entry — skip the manual "Join" button confirm
- [ ] *If player tries to submit result or outright but it is now locked (it locked between opening screen and submitting), show label "locked" for outrights and hide match from matches
- [ ] At least one KO must be selected from each group, implement both on mobile and BE
- [ ] *Highlight third teams of group that advance to playoffs

## Onboarding
- [ ] *Sign-up flow: prompt for display name, pre-fill from SSO profile first name


"be protagonist, not just a spectator"
