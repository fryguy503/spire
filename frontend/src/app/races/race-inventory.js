// The inventory is a client reference, not a query of live NPC spawn locations.
export function buildRaceInventory (entries, images, names) {
  const byId = new Map((entries || []).map(entry => [Number(entry.race_id), entry]))
  Object.keys(images).forEach(id => {
    if (!byId.has(Number(id))) byId.set(Number(id), { race_id: Number(id), sources: [] })
  })
  return Array.from(byId.values()).map(entry => {
    const sources = (entry.sources || []).map(source => ({
      ...source,
      models: source.models || [],
      zones: source.zones || []
    }))
    const codes = [
      { gender: 'Male', code: entry.male_gender_model_code },
      { gender: 'Female', code: entry.female_gender_model_code },
      { gender: 'Neutral', code: entry.neutral_gender_model_code }
    ].filter(model => model.code)
    const zones = Array.from(new Map(sources.flatMap(source => source.zones).map(zone => [zone.id, zone])).values())
      .sort((a, b) => a.long_name.localeCompare(b.long_name) || a.id - b.id)
    const name = entry.description && entry.description !== 'N/A'
      ? entry.description
      : names[entry.race_id] || 'Unknown race'
    return {
      ...entry,
      name,
      sources,
      codes,
      zones,
      images: images[entry.race_id] || [],
      global: sources.some(source => source.is_global),
      searchText: [entry.race_id, name, names[entry.race_id], ...codes.map(model => model.code),
        ...sources.flatMap(source => [source.source_file, ...source.models.map(model => model.code)]),
        ...zones.flatMap(zone => [zone.long_name, zone.short_name])].join(' ').toLowerCase()
    }
  }).sort((a, b) => a.race_id - b.race_id)
}

export function matchesZone (race, zoneId, includeGlobal) {
  return !zoneId || race.zones.some(zone => zone.id === zoneId) || (includeGlobal && race.global)
}

export function inventoryZones (races) {
  const zones = new Map()
  races.forEach(race => race.zones.forEach(zone => {
    if (!zones.has(zone.id)) zones.set(zone.id, { ...zone, count: 0 })
    zones.get(zone.id).count++
  }))
  return Array.from(zones.values()).sort((a, b) => a.long_name.localeCompare(b.long_name) || a.id - b.id)
}
