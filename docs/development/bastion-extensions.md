# Combined master and Bastion extensions

This fork combines Valorith master at `1e41e986` with the achievement release branch at `8b601230`. The original histories remain in Git. Releases and updater downloads use `fryguy503/spire`.

## Kinbound item settings

Open an existing item and select **Kinbound**. Its **Save Kinbound** action saves independently from the ordinary item form and checks for concurrent changes.

| Value | Meaning |
| --- | --- |
| 0 | Unlisted; original item behavior |
| 1 | Explicitly Always Kinbound |
| 2 | Never Kinbound |

Always requires the saved item to be NO TRADE (`nodrop = 0`), have NO TRANSFER disabled, and have no known epic exclusion in the Kinbound catalog. Expansion metadata does not enable items. Save ordinary flag changes first, then refresh Kinbound before selecting Always.

The selected content database needs the matching Source migrations `2026_09_19_kinbound_catalog.sql` (9397) and `2026_09_19_kinbound_item_flag.sql` (9398). Spire reports missing support and does not apply migrations. Full gameplay requires the matching Source/client deployment, its additional main-database migrations, and `Items:KinboundEnabled`. Existing permanent attunement remains intact.

To apply selection changes, follow Source's `docs/kinbound-dev-testing.md`: stop world/zones, regenerate the shared item cache with the matching `shared_memory` binary, and restart world/zones with the same active ruleset. A live rule reload alone does not apply list changes.

Restricted connection users need **Item Kinbound** read/write permissions in addition to their item permissions. Kinbound routes are hand-authored and separate from generated item CRUD, so older item schemas and API regeneration remain compatible.

## Saved expedition history

Open **Server Admin > Player Operations**, select a character, then open **Connections > Saved expeditions**. The list provides an eligible-only filter, membership counts, expiration, and blocker descriptions.

The character database needs Source's `2026_09_20_dynamic_zone_resume.sql` and `2026_09_21_dynamic_zone_resume_latest.sql`, including the matching InnoDB tables. `Expedition:EnableDzResume` controls rejoining. Spire reads the configured world ruleset, season bucket, historical eligibility, active roster, and lockouts. A newer copy's durable record excludes the older copy even after the newer instance disappears.

The list is a database snapshot. The running server rechecks loaded rules, invitations, current location, and admission when a player uses `/dzmanager`. Spire never rejoins or modifies an expedition through this view. Missing migrations leave the rest of character administration usable. Access uses the existing **Player Operations** read permission.

## Build and release

Use Node 20 and Go 1.23.12. The GitHub **Manual Release** workflow installs locked frontend dependencies, builds the embedded EQ Sage assets and Vue application, packs the frontend, and publishes Linux/Windows Spire ZIPs plus server installers. The top changelog entry is marked Beta, so the combined release is a GitHub prerelease.

The repository is public. Publication to `fryguy503/spire` does not deploy a running Bastion server or apply its database migrations. Discord release messages are off unless the repository owner explicitly enables `SPIRE_DISCORD_RELEASE_ANNOUNCEMENTS`.
