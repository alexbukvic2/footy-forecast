# Backend

## Core
- [ ] Refresh token handling - caching user, how do we manage cache memory? What if there are many users?
- [ ] Admin role for entity creation
- [ ] Select statements outside SQL files, e.g. in go files, is it okay? Should we move all SQL to files?

## Features
- [ ] Historical leaderboard snapshots (point-in-time rank history).
- [ ] Pagination of leaderboard results.
- [ ] Fixtures in group have group, but what about quarter finals, etc?
- [x] *Locked fixtures - lock by worker n minutes in advance
- [x] *Same points breaking criteria
- [ ] *Default values for missing handicaps
- [x] Worker should fetch fixtures every day

### Profile page
- [x] Change display name

---

# Mobile

## Core
- [ ] Share button: share link to entity, or deep link to app with entity preloaded
- [x] Auto-join on valid code entry — skip the manual "Join" button confirm
- [ ] *If player tries to submit result or outright but it is now locked (it locked between opening screen and submitting), show label "locked" for outrights and hide match from matches
- [ ] At least one KO must be selected from each group, implement both on mobile and BE
- [x] *Highlight third teams of group that advance to playoffs

## Onboarding
- [x] *Sign-up flow: prompt for display name, pre-fill from SSO profile first name

## Notifications
- [ ] Player should land on profile page after signup
- [ ] Click on push notification leads to appropriate screen

## Important
- [ ] Create migration for players
- [ ] Rewards for teams that go to playoffs as third ones, how does API show those
- [ ] Playoff matches - different handicap, if player sets draw, needs to choose winner
- [ ] Players that do not have handicap, set default
- [ ] Change app name and icon


"be protagonist, not just a spectator"
