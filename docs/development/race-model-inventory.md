# Race model inventory

The existing `/race-viewer` combines Maudigan's model previews with the bundled
Rain of Fear 2 inventory mirrored at https://races.clumsysworld.com/ (originally
Shendare's EQ Race Inventory). This is client model availability, not live NPC
spawn data. Archives may be local to a zone, imported by a zone, or loaded
globally (sometimes conditionally through client settings).

The viewer uses `/api/v1/static-map/race-inventory-map.json`. Browsing does not
contact the reference website or require a zone database query. Race names, IDs,
model codes, source archive names, and zone names are searchable. A zone filter
can include global models; each location retains its source archive association.
Missing previews, missing source archives, and archives with no known zones are
shown explicitly. The URL preserves the search, zone, global toggle, and race.
The model locations layout is the default; the Classic gallery toggle restores
the full appearance gallery and persists as `view=classic` in the URL. Switching
layouts retains filters and race selection. Each source model has a copy action
for `CODE,archive_name`, omitting the `.s3d` or `.eqg` extension (for example,
`CLM,steamfont_chr`), with inline success or failure feedback.

Refresh the generated snapshot from the repository root:

```sh
go run . make:race-model-maps
```

To reproduce from a saved copy of the reference:

```sh
go run . make:race-model-maps --source-file /path/to/reference.html
```

Do not hand-edit the JSON. The parser lives in
`internal/generators/race_model_maps_cmd.go`; its focused fixtures preserve
single-gender codes, high race IDs, hyphenated names, and archive-zone links.

```sh
go test ./internal/generators -run TestRaceInventoryReference
npx playwright test tests/race-viewer.spec.ts
```
