# Plan: Bulk PUT for Tournament Predictions

## Goal

Replace the two per-category `PUT .../players/{category}` and `PUT .../teams/{category}` endpoints with two bulk `PUT .../players` and `PUT .../teams` endpoints that accept the full set of picks at once. Sending `null` for a `player_id` / `team_id` clears (deletes) that slot. The lock check and wildcard cap apply once to the whole batch. The response mirrors the existing GET view so clients always receive current state after a write.

## Motivation

- The mobile UI saves all categories in one form submission — a bulk API matches that UX with a single atomic call.
- The wildcard cap (max 8 playoff slot_index=2 picks) is a global constraint; it is cleaner to enforce it across the whole batch in memory than one item at a time with per-request DB queries.
- Eliminates the need for a separate DELETE endpoint — `null` encodes "clear this slot" directly.
- The GET already returns all slots at once; the write mirrors read.

## Data model changes

No schema/migration changes. Two new SQL delete queries are added:

- `DeletePlayerPredictionGroup` — deletes a group-scoped player prediction by (user, tournament, category, group_letter)
- `DeletePlayerPredictionNoGroup` — deletes a non-group player prediction by (user, tournament, category)
- `DeleteTeamPredictionGroup` — deletes a group-scoped team prediction by (user, tournament, category, group_letter, slot_index)
- `DeleteTeamPredictionNoGroup` — deletes a no-group team prediction by (user, tournament, category, slot_index)

Run `make sqlc-gen` after adding these queries.

## API Contract

### Removed paths

Delete these entries from `docs/openapi.yaml`:

```
/tournaments/{tournamentId}/predictions/players/{category}  (entire path object)
/tournaments/{tournamentId}/predictions/teams/{category}    (entire path object)
```

Remove schemas `UpsertPlayerPredictionRequest`, `PlayerPredictionResponse`, `UpsertTeamPredictionRequest`, `TeamPredictionResponse` if no longer referenced.

### New schemas (add to `components/schemas`)

```yaml
BulkPlayerPredictionItem:
  type: object
  required: [category]
  properties:
    category:
      $ref: '#/components/schemas/PlayerHandicapCategory'
    group_letter:
      type: string
      nullable: true
    player_id:
      type: string
      format: uuid
      nullable: true
      description: "null clears the slot"

BulkUpsertPlayerPredictionsRequest:
  type: array
  items:
    $ref: '#/components/schemas/BulkPlayerPredictionItem'

BulkTeamPredictionItem:
  type: object
  required: [category, slot_index]
  properties:
    category:
      $ref: '#/components/schemas/TeamHandicapCategory'
    group_letter:
      type: string
      nullable: true
    slot_index:
      type: integer
    team_id:
      type: string
      format: uuid
      nullable: true
      description: "null clears the slot"

BulkUpsertTeamPredictionsRequest:
  type: array
  items:
    $ref: '#/components/schemas/BulkTeamPredictionItem'
```

### New `put` operation on `/tournaments/{tournamentId}/predictions/players`

```yaml
put:
  summary: Bulk upsert player tournament predictions
  operationId: bulkUpsertPlayerPredictions
  tags: [predictions]
  security:
    - bearerAuth: []
  parameters:
    - name: tournamentId
      in: path
      required: true
      schema:
        type: string
        format: uuid
  requestBody:
    required: true
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/BulkUpsertPlayerPredictionsRequest'
  responses:
    '200':
      description: Current player predictions after applying the batch
      content:
        application/json:
          schema:
            type: array
            items:
              $ref: '#/components/schemas/PlayerPredictionView'
    '400':
      description: Invalid input
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '401':
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '403':
      description: Predictions locked
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '404':
      description: Player not found in this tournament
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '500':
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
```

### New `put` operation on `/tournaments/{tournamentId}/predictions/teams`

```yaml
put:
  summary: Bulk upsert team tournament predictions
  operationId: bulkUpsertTeamPredictions
  tags: [predictions]
  security:
    - bearerAuth: []
  parameters:
    - name: tournamentId
      in: path
      required: true
      schema:
        type: string
        format: uuid
  requestBody:
    required: true
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/BulkUpsertTeamPredictionsRequest'
  responses:
    '200':
      description: Current team predictions after applying the batch
      content:
        application/json:
          schema:
            type: array
            items:
              $ref: '#/components/schemas/TeamPredictionView'
    '400':
      description: Invalid input
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '401':
      description: Unauthorized
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '403':
      description: Predictions locked or wildcard cap exceeded
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '404':
      description: Team not found in this tournament
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
    '500':
      description: Internal server error
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
```

## Edge cases

1. **Empty batch (`[]`)** — valid no-op; returns current state.
2. **Duplicate slots in same batch** — last writer wins (items processed in order).
3. **`null` pick for a slot that doesn't exist** — no-op (delete is idempotent); do not return 404.
4. **Wildcard cap check** — computed as `dbCount + batchAdds - batchClears` where adds/clears refer to playoff slot_index=2 items in the batch. Conservative: may slightly overcount when a batch add updates an already-existing wildcard slot, but this is a rare edge case.
5. **Lock check fires once** for the whole batch before any persistence.

## Implementation steps

1. Add delete SQL queries to `internal/repository/queries/tournament_predictions.sql`.
2. Run `make sqlc-gen`.
3. Add `BulkPlayerPredictionItem` and `BulkTeamPredictionItem` input types to `internal/domain/tournament_prediction.go`.
4. Add `DeletePlayer` / `DeleteTeam` methods to repo interfaces in `internal/service/tournament_prediction.go`; add `BulkUpsertPlayerPredictions` and `BulkUpsertTeamPredictions` service methods.
5. Implement `DeletePlayer` / `DeleteTeam` in `internal/repository/tournament_prediction.go`.
6. Replace `UpsertPlayerPredictions` / `UpsertTeamPredictions` handler methods with bulk variants; update `TournamentPredictionSvc` interface.
7. Update `internal/server/router.go` — swap per-category routes with bulk routes.
8. Update `docs/openapi.yaml`; run `make generate`.
9. Replace Postman request files.
10. Update tests.

## Test plan

- **Service**: bulk happy path (mix of upsert + clear), lock check fires on whole batch, wildcard cap enforced, empty batch returns current state, invalid player/team returns ErrNotFound.
- **Handler**: valid JSON decoded and forwarded, null pick accepted, malformed JSON → 400, service ErrForbidden → 403, service ErrNotFound → 404.
- **Repository (integration)**: `DeletePlayerPrediction` deletes matching row and is idempotent; `DeleteTeamPrediction` same.

## Acceptance criteria

- `PUT /tournaments/{id}/predictions/players` with mixed upsert+clear batch returns 200 with full view.
- `PUT /tournaments/{id}/predictions/teams` with 9 wildcard adds returns 403.
- Old `/{category}` PUT routes return 404 (not registered).
- `make fmt`, `make lint`, `make test` all pass.
- `make generate` produces in-sync `docs/openapi.json` and `internal/server/oapi/models.gen.go`.
- Postman bulk request files replace the old per-category ones.

