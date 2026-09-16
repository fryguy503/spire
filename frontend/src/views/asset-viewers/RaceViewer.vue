<template>
  <div class="race-viewer" :class="{ 'locations-layout': viewMode === 'locations' }">
    <eq-window title="Race Viewer" class="race-toolbar">
      <div class="race-heading">
        <div>
          <h1>Race models <span>&amp; locations</span></h1>
          <p>Find a model, its client archive, and the zones that load it.</p>
        </div>
        <a :href="referenceUrl" target="_blank" rel="noopener noreferrer" class="reference-link">
          RoF2 reference <i class="fa fa-external-link" aria-hidden="true"></i>
        </a>
      </div>
      <div class="view-toggle" role="group" aria-label="Race viewer layout">
        <button class="btn btn-sm" :class="{ active: viewMode === 'locations' }" :aria-pressed="viewMode === 'locations'" @click="setView('locations')">
          <i class="fa fa-columns" aria-hidden="true"></i> Model locations
        </button>
        <button class="btn btn-sm" :class="{ active: viewMode === 'classic' }" :aria-pressed="viewMode === 'classic'" @click="setView('classic')">
          <i class="fa fa-th" aria-hidden="true"></i> Classic gallery
        </button>
      </div>
      <div class="race-filters">
        <div class="race-search">
          <label for="race-search">Search models</label>
          <input id="race-search" v-model="raceSearch" class="form-control" type="search"
                 placeholder="Race, ID, model code, archive or zone…" @input="queueQuery()">
        </div>
        <div class="zone-search">
          <label for="race-zone">Available in zone</label>
          <select id="race-zone" v-model.number="zoneSearch" class="form-control" @change="updateQuery()">
            <option :value="0">All zones</option>
            <option v-if="zoneSearch && !zoneList.some(zone => zone.id === zoneSearch)" :value="zoneSearch">
              Zone #{{ zoneSearch }} — not in reference
            </option>
            <option v-for="zone in zoneList" :key="zone.id" :value="zone.id">
              {{ zone.long_name }} · {{ zone.short_name }} (#{{ zone.id }})
            </option>
          </select>
        </div>
        <button class="btn btn-outline-light reset-button" @click="reset">Reset</button>
      </div>
      <div class="race-filter-footer">
        <span aria-live="polite">{{ filteredRaces.length }} {{ filteredRaces.length === 1 ? 'race' : 'races' }}</span>
        <label v-if="zoneSearch" class="global-toggle">
          <input v-model="includeGlobal" type="checkbox" @change="updateQuery()"> Include global models
        </label>
        <span class="reference-note">Client model availability · not NPC spawn locations</span>
      </div>
    </eq-window>

    <div v-if="loading" class="race-message" role="status">Loading the model inventory…</div>
    <div v-else-if="inventoryError" class="race-message" role="alert">
      <p>Model locations could not be loaded. {{ races.length ? 'Image previews are still available.' : 'Please try again.' }}</p>
      <button class="btn btn-outline-light" @click="loadInventory">Retry model inventory</button>
    </div>
    <p v-if="previewError" class="preview-notice" role="status">Image previews are unavailable. Model codes and locations are still searchable.</p>

    <eq-window v-if="!loading && races.length && viewMode === 'classic'" class="classic-window" data-testid="classic-gallery">
      <div v-if="!filteredRaces.length" class="race-message">
        <h2>No matching races</h2>
        <p>Try another model code or clear the zone filter.</p>
        <button class="btn btn-outline-light" @click="reset">Clear filters</button>
      </div>
      <div class="classic-gallery">
        <div v-for="race in visibleRaces" :key="race.race_id" class="classic-race">
          <div class="classic-previews">
            <span v-for="img in race.images" :key="img" :class="'race-models-ctn-' + img"
                  :title="imageDescription(img)" role="img" :aria-label="imageDescription(img)"></span>
            <span v-if="!race.images.length" class="muted">No preview available</span>
          </div>
          <h2 class="eq-header">{{ race.name }} ({{ race.race_id }})</h2>
        </div>
      </div>
      <button v-if="visibleRaces.length < filteredRaces.length" class="btn btn-outline-light load-more" @click="visibleLimit += 60">Show more races</button>
      <p class="image-credit">Model previews by Maudigan</p>
    </eq-window>

    <div v-if="!loading && races.length && viewMode === 'locations'" class="race-workspace">
      <eq-window class="race-gallery-window">
        <div class="gallery-heading">
          <span>MODEL GALLERY</span><span>{{ visibleRaces.length }} of {{ filteredRaces.length }}</span>
        </div>
        <div v-if="!filteredRaces.length" class="race-message">
          <h2>No matching races</h2>
          <p>Try another model code or clear the zone filter.</p>
          <button class="btn btn-outline-light" @click="reset">Clear filters</button>
        </div>
        <div class="race-gallery" id="race-viewer-viewport">
          <button v-for="race in visibleRaces" :key="race.race_id" class="race-tile"
                  :class="{ selected: selectedRace && selectedRace.race_id === race.race_id }"
                  :aria-pressed="!!selectedRace && selectedRace.race_id === race.race_id"
                  :aria-label="`${race.name}, race ${race.race_id}. View model locations`"
                  @click="selectRace(race.race_id)">
            <span class="tile-id">#{{ race.race_id }}</span>
            <span class="tile-preview" aria-hidden="true">
              <span v-for="img in previewImages(race).slice(0, 2)" :key="img" :class="'race-models-ctn-' + img"></span>
              <span v-if="!race.images.length" class="no-preview"><i class="ra ra-aware"></i><span>No preview</span></span>
            </span>
            <span class="tile-name">{{ race.name }}</span>
            <span class="tile-codes">{{ codeSummary(race) || 'Code not listed' }}</span>
            <span class="tile-location"><i :class="race.global ? 'fa fa-globe' : 'fa fa-map-marker'" aria-hidden="true"></i>
              {{ locationSummary(race) }}
            </span>
          </button>
        </div>
        <button v-if="visibleRaces.length < filteredRaces.length" class="btn btn-outline-light load-more" @click="visibleLimit += 60">
          Show more races
        </button>
        <p class="image-credit">Model previews by Maudigan</p>
      </eq-window>

      <eq-window v-if="selectedRace" class="race-inspector" ref="raceInspector" data-testid="race-inspector">
        <div class="inspector-heading">
          <div><span class="eyebrow">RACE {{ selectedRace.race_id }}</span><h2>{{ selectedRace.name }}</h2></div>
          <span v-if="selectedRace.is_playable" class="availability-tag">Playable</span>
        </div>
        <div class="inspector-scroll" :key="selectedRace.race_id">
          <div v-if="selectedRace.images.length" class="detail-preview">
            <span v-for="img in (showAppearances ? selectedRace.images : previewImages(selectedRace))" :key="img" :title="imageDescription(img)"
                  :class="'race-models-ctn-' + img" role="img" :aria-label="imageDescription(img)"></span>
          </div>
          <button v-if="selectedRace.images.length > previewImages(selectedRace).length" class="appearance-toggle btn btn-sm btn-outline-light" @click="showAppearances = !showAppearances">
            {{ showAppearances ? 'Show fewer appearances' : `View all ${selectedRace.images.length} appearances` }}
          </button>
          <div class="model-codes">
            <div v-for="model in selectedRace.codes" :key="model.gender" class="model-code">
              <span>{{ model.gender }}</span>
              <button :aria-label="`Copy ${model.gender.toLowerCase()} model code ${model.code}`" @click="copyCode(model.code)">
                {{ model.code }} <i class="fa fa-copy" aria-hidden="true"></i>
              </button>
            </div>
            <p v-if="!selectedRace.codes.length" class="muted">No model code listed in the reference.</p>
          </div>
          <p class="copy-status" role="status" aria-live="polite">{{ copyStatus }}</p>
          <div class="location-heading">
            <h3>Model locations</h3>
            <span>{{ selectedSources.length }} {{ selectedSources.length === 1 ? 'archive' : 'archives' }}</span>
          </div>
          <p v-if="zoneSearch" class="scope-note">Showing sources available in {{ selectedZoneName }}.</p>
          <p v-if="!selectedSources.length" class="empty-sources">No source archives or zone locations are listed for this race.</p>
          <section v-for="source in selectedSources" :key="source.source_file" class="model-source">
            <div class="source-heading">
              <code>{{ source.source_file }}</code>
              <span v-if="source.is_global" class="availability-tag"><i class="fa fa-globe" aria-hidden="true"></i> Global</span>
            </div>
            <div v-if="source.models.length" class="source-models">
              <span v-for="(model, index) in source.models" :key="index">
                <button v-if="model.code" class="source-copy" :aria-label="`Copy ${genderName(model.gender).toLowerCase()} model location ${modelLocation(model, source)}`"
                        :title="modelLocation(model, source)" @click="copyCode(modelLocation(model, source))">
                  <code>{{ model.code }}</code> {{ genderName(model.gender) }}
                  <i :class="copiedValue === modelLocation(model, source) ? 'fa fa-check' : 'fa fa-copy'" aria-hidden="true"></i>
                  <span v-if="copiedValue === modelLocation(model, source)">Copied</span>
                </button>
                <span v-else>Unspecified {{ genderName(model.gender) }}</span>
                <span v-if="copyErrorValue === modelLocation(model, source)" class="copy-error" role="alert">Could not copy. Select: {{ copyErrorValue }}</span>
              </span>
            </div>
            <p v-if="source.is_global" class="global-description">
              All zones<span v-if="source.loaded_via_description"> · {{ source.loaded_via_description }}</span>
            </p>
            <p v-else-if="source.loaded_via_description" class="scope-note">{{ source.loaded_via_description }}</p>
            <ul v-if="source.zones.length" class="source-zones">
              <li v-for="zone in source.zones" :key="zone.id + '-' + zone.import_type" :class="{ 'matched-zone': zone.id === zoneSearch }">
                <button class="zone-location" :aria-label="`Find models in ${zone.long_name} (${zone.short_name})`" @click="filterZone(zone.id)">
                  <span>{{ zone.long_name }}</span><small>{{ zone.short_name }} <span>· #{{ zone.id }}</span></small>
                </button>
                <span class="zone-kind" :class="{ imported: zone.import_type === 'Imported' }">{{ zone.import_type || 'Listed' }}</span>
              </li>
            </ul>
            <p v-else-if="!source.is_global" class="scope-note">No zones listed for this archive.</p>
          </section>
          <div class="inventory-attribution">
            <p><strong>Local</strong> is a zone’s own model archive. <strong>Imported</strong> is loaded from another archive. Global loading can depend on client settings.</p>
            <a :href="referenceUrl" target="_blank" rel="noopener noreferrer">Clumsy’s World · Shendare’s RoF2 race inventory <i class="fa fa-external-link" aria-hidden="true"></i></a>
            <p>Bundled reference; your client files may differ.</p>
          </div>
        </div>
      </eq-window>
    </div>
  </div>
</template>

<script>
  import { RACES } from '@/app/constants/eq-race-constants'
  import EqWindow from '@/components/eq-ui/EQWindow'
  import EqAssets from '@/app/eq-assets/eq-assets'
  import { SpireApi } from '@/app/api/spire-api'
  import { buildRaceInventory, inventoryZones, matchesZone } from '@/app/races/race-inventory'

  export default {
    components: { EqWindow },
    data () {
      return {
        races: [],
        zoneList: [],
        raceSearch: '',
        zoneSearch: 0,
        includeGlobal: true,
        selectedId: 0,
        visibleLimit: 60,
        loading: true,
        inventoryError: false,
        previewError: false,
        copyStatus: '',
        copiedValue: '',
        copyErrorValue: '',
        viewMode: 'locations',
        showAppearances: false,
        referenceUrl: 'https://races.clumsysworld.com/'
      }
    },
    computed: {
      filteredRaces () {
        const terms = this.raceSearch.trim().toLowerCase().split(/\s+/).filter(Boolean)
        return this.races.filter(race => matchesZone(race, this.zoneSearch, this.includeGlobal) &&
          terms.every(term => race.searchText.includes(term)))
      },
      visibleRaces () { return this.filteredRaces.slice(0, this.visibleLimit) },
      selectedRace () {
        return this.filteredRaces.find(race => race.race_id === this.selectedId) || this.filteredRaces[0] || null
      },
      selectedSources () {
        if (!this.selectedRace) return []
        return this.selectedRace.sources.filter(source => !this.zoneSearch ||
          source.zones.some(zone => zone.id === this.zoneSearch) || (this.includeGlobal && source.is_global))
      },
      selectedZoneName () {
        const zone = this.zoneList.find(zone => zone.id === this.zoneSearch)
        return zone ? zone.long_name : `zone #${this.zoneSearch}`
      }
    },
    watch: {
      '$route.query': 'loadQuery',
      raceSearch () { this.visibleLimit = 60 },
      zoneSearch () { this.visibleLimit = 60 },
      includeGlobal () { this.visibleLimit = 60 },
      selectedRace () { this.copyStatus = ''; this.copiedValue = ''; this.copyErrorValue = ''; this.showAppearances = false }
    },
    methods: {
      updateViewportOffset () {
        // Spire applies CSS zoom. Convert the remaining screen space to layout pixels.
        const bounds = this.$el.getBoundingClientRect()
        const scale = bounds.width / parseFloat(window.getComputedStyle(this.$el).width) || 1
        const top = bounds.top + window.scrollY
        this.$el.style.setProperty('--race-page-height', `${Math.max(0, (window.innerHeight - top - 24) / scale)}px`)
      },
      previewImages (race) {
        const genders = new Set()
        return race.images.filter(image => {
          const gender = image.split('-')[1]
          if (genders.has(gender)) return false
          genders.add(gender)
          return true
        })
      },
      codeSummary (race) { return [...new Set(race.codes.map(model => model.code))].join(' / ') },
      locationSummary (race) {
        if (race.global) return 'Global model'
        if (!race.zones.length) return 'No zones listed'
        if (race.zones.length === 1) return race.zones[0].long_name
        return `${race.zones.length} zones · ${race.zones[0].short_name}`
      },
      genderName (gender) { return ['Male', 'Female', 'Neutral'][gender] || 'Unspecified' },
      imageDescription (img) {
        const [race, gender, texture, helm] = img.split('-')
        return `Race ${race} · ${this.genderName(Number(gender))} · Texture ${texture} · Helm ${helm}`
      },
      loadQuery () {
        const q = this.$route.query
        this.viewMode = q.view === 'classic' ? 'classic' : 'locations'
        this.raceSearch = typeof q.raceSearch === 'string' ? q.raceSearch : ''
        this.zoneSearch = /^\d+$/.test(q.zoneSearch) ? Number(q.zoneSearch) : 0
        this.selectedId = /^\d+$/.test(q.race) ? Number(q.race) : 0
        this.includeGlobal = q.global !== '0'
        const index = this.filteredRaces.findIndex(race => race.race_id === this.selectedId)
        this.visibleLimit = Math.max(60, index + 1)
      },
      queueQuery () {
        clearTimeout(this.queryTimer)
        this.queryTimer = setTimeout(() => this.updateQuery(true), 250)
      },
      updateQuery (replace = false) {
        clearTimeout(this.queryTimer)
        const query = {}
        if (this.viewMode === 'classic') query.view = 'classic'
        if (this.raceSearch) query.raceSearch = this.raceSearch
        if (this.zoneSearch) query.zoneSearch = String(this.zoneSearch)
        if (!this.includeGlobal) query.global = '0'
        if (this.selectedRace) query.race = String(this.selectedRace.race_id)
        this.$router[replace ? 'replace' : 'push']({ path: this.$route.path, query }).catch(() => {})
      },
      selectRace (id) {
        this.selectedId = id
        this.updateQuery()
        // Keep the selected model in view when the panes stack on small screens.
        if (window.matchMedia('(max-width: 760px)').matches) {
          this.$nextTick(() => {
            const inspector = this.$refs.raceInspector
            if (inspector) inspector.$el.scrollIntoView({ block: 'start', behavior: 'auto' })
          })
        }
      },
      filterZone (id) { this.raceSearch = ''; this.zoneSearch = id; this.updateQuery() },
      reset () {
        this.raceSearch = ''; this.zoneSearch = 0; this.selectedId = 0; this.includeGlobal = true
        this.updateQuery()
      },
      setView (view) { this.viewMode = view; this.updateQuery() },
      modelLocation (model, source) {
        return `${model.code},${source.source_file.replace(/\.(s3d|eqg)$/i, '')}`
      },
      async copyCode (code) {
        this.copiedValue = ''
        this.copyErrorValue = ''
        try {
          await navigator.clipboard.writeText(code)
          this.copyStatus = `Copied ${code}`
          this.copiedValue = code
        } catch (_) {
          this.copyStatus = `Could not copy. Select: ${code}`
          this.copyErrorValue = code
        }
      },
      async loadInventory () {
        this.loading = true
        this.inventoryError = false
        this.previewError = false
        const images = {}
        let entries = []
        await Promise.all([
          SpireApi.v1().get('/static-map/race-inventory-map.json').then(result => {
            if (!Array.isArray(result.data.races) || !result.data.races.length) throw new Error('Invalid inventory')
            entries = result.data.races
          }).catch(() => { this.inventoryError = true }),
          EqAssets.getNpcModels().then(models => {
            models.forEach(model => {
              const match = model.fileName.match(/^CTN_(\d+)_(\d+)_(\d+)_(\d+)\.png$/)
              if (!match) return
              const id = Number(match[1])
              if (!images[id]) images[id] = []
              images[id].push(match.slice(1).join('-'))
            })
          }).catch(() => { this.previewError = true })
        ])
        if (this.disposed) return
        this.races = Object.freeze(buildRaceInventory(entries, images, RACES))
        this.zoneList = Object.freeze(inventoryZones(this.races))
        this.loading = false
        this.loadQuery()
      }
    },
    mounted () {
      this.loadQuery()
      this.loadInventory()
      this.updateViewportOffset()
      this.layoutObserver = new ResizeObserver(this.updateViewportOffset)
      this.layoutObserver.observe(this.$el)
      window.addEventListener('resize', this.updateViewportOffset)
    },
    beforeDestroy () {
      this.disposed = true
      clearTimeout(this.queryTimer)
      this.layoutObserver.disconnect()
      window.removeEventListener('resize', this.updateViewportOffset)
    }
  }
</script>

<style scoped>
.race-viewer { --race-accent: #d5bd84; --race-muted: #a9b1bc; color: #e5e8ee; }
.race-toolbar { padding: 18px 22px !important; }
.race-heading, .race-filters, .race-filter-footer, .gallery-heading, .inspector-heading, .location-heading, .source-heading { display: flex; align-items: center; justify-content: space-between; gap: 16px; }
.race-heading { margin: 10px 0 22px; }
.race-heading h1 { font-size: 24px; margin: 0 0 6px; font-weight: 600; color: #f0e4c9; }
.race-heading h1 span { font-weight: 400; color: #c4cbd4; }
.race-heading p, .scope-note, .global-description { margin: 0; font-size: 13px; color: var(--race-muted); }
.reference-link { font-size: 12px; white-space: nowrap; color: var(--race-accent); }
.reference-link i { margin-left: 6px; font-size: 10px; }
.view-toggle { display: inline-flex; gap: 4px; margin: 0 0 18px; padding: 3px; border: 1px solid #ffffff25; border-radius: 5px; }
.race-viewer .view-toggle button.btn { background: transparent; border: 1px solid transparent; color: var(--race-muted); padding: 7px 11px; }
.race-viewer .view-toggle button.btn.active { background: #d5bd8418; border-color: #d5bd8460; color: #f0e4c9; }
.view-toggle i { margin-right: 5px; }
.classic-window { margin-top: 24px; }
.classic-gallery { display: flex; flex-wrap: wrap; align-items: center; justify-content: center; gap: 24px; max-height: calc(100vh - 380px); min-height: 220px; overflow-y: auto; padding: 8px; }
.classic-race { border: 2px solid #dadada1a; border-radius: 5px; padding: 24px 18px 18px; text-align: center; max-width: 100%; }
.classic-previews { display: flex; flex-wrap: wrap; align-items: center; justify-content: center; }
.classic-previews > span[role="img"] { flex-shrink: 0; filter: drop-shadow(10px 5px 7px #000); }
.classic-race h2 { font-size: 14px; margin: 32px 0 0; }
.race-viewer .source-models button.source-copy { display: inline-flex; align-items: center; gap: 6px; padding: 5px 7px; border: 1px solid #d5bd8440; background: transparent; color: var(--race-muted); font-size: 11px; border-radius: 4px; }
.race-viewer .source-models button.source-copy:hover { border-color: var(--race-accent); box-shadow: none; background: #d5bd8418; }
.source-copy i, .source-copy > span { color: var(--race-accent); }
.copy-error { display: block; margin-top: 5px; user-select: text; overflow-wrap: anywhere; }
.race-filters { align-items: flex-end; }
.race-search { flex: 1.2; min-width: 0; }
.zone-search { flex: 1; min-width: 0; }
.race-filters label { display: block; margin-bottom: 7px; font-size: 12px; font-weight: 600; }
.race-filters .form-control { width: 100%; color: #e5e8ee; background: #121923; border: 1px solid #47505d; height: 42px; }
.race-filters .form-control::placeholder { color: #a9b1bc; }
.reset-button { height: 42px; }
.race-filter-footer { justify-content: flex-start; margin-top: 15px; font-size: 12px; color: var(--race-muted); flex-wrap: wrap; }
.race-filter-footer > span:first-child { color: var(--race-accent); font-variant-numeric: tabular-nums; }
.reference-note { margin-left: auto; }
.global-toggle { margin: 0; display: flex; gap: 7px; align-items: center; cursor: pointer; }
.global-toggle input { accent-color: var(--race-accent); }
.race-workspace { display: grid; grid-template-columns: minmax(0, 1fr) minmax(340px, 40%); gap: 18px; margin-top: 18px; align-items: start; }
.race-gallery-window { padding: 18px !important; min-width: 0; }
.gallery-heading { color: var(--race-muted); font-size: 11px; letter-spacing: .08em; margin: 0 2px 18px; }
.gallery-heading span:last-child { letter-spacing: 0; }
.race-gallery { display: grid; align-content: start; grid-template-columns: repeat(auto-fill, minmax(170px, 1fr)); gap: 12px; }
.race-viewer .race-gallery button.race-tile { position: relative; display: flex; flex-direction: column; align-items: center; padding: 15px 10px; min-width: 0; border: 1px solid #ffffff16; border-radius: 5px; background: #10172270; color: inherit; cursor: pointer; text-align: center; transition: border-color .18s, background .18s, transform .18s; }
.race-viewer .race-gallery button.race-tile:hover { border-color: #d5bd8480; background: #202a36; transform: translateY(-2px); }
.race-viewer .race-gallery button.race-tile.selected { border-color: var(--race-accent); background: #d5bd8410; box-shadow: inset 0 -2px 0 var(--race-accent); }
.race-viewer button:focus-visible, .race-viewer a:focus-visible { outline: 2px solid var(--race-accent); outline-offset: 3px; }
.tile-id { align-self: flex-start; color: var(--race-muted); font-size: 11px; }
.tile-preview { display: flex; align-items: center; justify-content: center; height: 136px; width: 100%; overflow: hidden; }
.tile-preview > span:not(.no-preview) { flex-shrink: 0; filter: drop-shadow(4px 5px 4px #000); }
.tile-name { font-size: 16px; color: #f0e4c9; line-height: 1.3; margin-top: 6px; }
.tile-codes { font-family: monospace; font-size: 12px; color: #d0d8e2; letter-spacing: .05em; margin: 7px 0 9px; }
.tile-location { font-size: 11px; color: var(--race-muted); line-height: 1.5; }
.tile-location i { margin-right: 4px; color: var(--race-accent); }
.no-preview { display: flex; flex-direction: column; gap: 9px; color: #8b96a6; font-size: 11px; }
.no-preview i { font-size: 32px; opacity: .55; }
.race-inspector { position: sticky; top: 18px; padding: 0 !important; min-width: 0; }
.inspector-heading { padding: 22px 22px 16px; border-bottom: 1px solid #ffffff15; }
.eyebrow { color: var(--race-accent); font-size: 10px; letter-spacing: .13em; }
.inspector-heading h2 { font-size: 24px; margin: 5px 0 0; color: #f0e4c9; }
.inspector-scroll { padding: 0 22px 22px; max-height: calc(100vh - 150px); overflow-y: auto; animation: inspector-in .2s ease-out; }
.detail-preview { display: flex; flex-wrap: wrap; align-items: center; justify-content: center; gap: 8px; padding: 20px 0 12px; max-height: 260px; overflow: auto; }
.detail-preview > span { flex-shrink: 0; filter: drop-shadow(5px 8px 5px #000); }
.model-codes { display: flex; gap: 24px; padding-top: 18px; flex-wrap: wrap; }
.model-code > span { display: block; color: var(--race-muted); font-size: 11px; margin-bottom: 3px; }
.race-viewer .model-codes .model-code button { color: var(--race-accent); background: transparent; border: 0; padding: 2px 0; font: 600 21px monospace; cursor: pointer; }
.model-code i { font-size: 11px; opacity: .6; margin-left: 6px; }
.race-viewer .source-zones button.zone-location:hover, .race-viewer .model-codes .model-code button:hover { background: transparent; border: 0; box-shadow: none; }
.appearance-toggle { display: block; margin: 5px auto 0; font-size: 11px; }
.copy-status { min-height: 17px; margin: 5px 0 14px; font-size: 11px; color: var(--race-accent); }
.location-heading { margin-bottom: 12px; }
.location-heading h3 { font-size: 15px; margin: 0; color: #e5e8ee; }
.location-heading > span { color: var(--race-muted); font-size: 11px; }
.model-source { border-top: 1px solid #ffffff1c; padding: 16px 0 4px; margin-top: 10px; }
.source-heading { gap: 8px; align-items: flex-start; }
.source-heading > code { font-size: 13px; color: #f0e4c9; overflow-wrap: anywhere; }
.availability-tag { display: inline-block; padding: 3px 6px; border: 1px solid #d5bd8440; border-radius: 3px; color: var(--race-accent); font-size: 10px; white-space: nowrap; }
.source-models { display: flex; flex-wrap: wrap; gap: 7px 14px; margin: 7px 0 10px; color: var(--race-muted); font-size: 11px; }
.source-models code { color: #d5dae1; font-size: 11px; }
.global-description { padding-bottom: 10px; }
.source-zones { list-style: none; margin: 7px 0 0; padding: 0; }
.source-zones li { display: flex; align-items: center; justify-content: space-between; gap: 8px; border-top: 1px solid #ffffff0a; padding: 8px 0; }
.source-zones li.matched-zone { border-left: 2px solid var(--race-accent); padding-left: 10px; }
.race-viewer .source-zones button.zone-location { text-align: left; background: transparent; padding: 0; border: 0; color: #e5e8ee; cursor: pointer; min-width: 0; }
.zone-location > span { display: block; font-size: 12px; line-height: 1.5; }
.zone-location small { display: block; font: 11px monospace; color: var(--race-muted); margin-top: 2px; overflow-wrap: anywhere; }
.race-viewer .source-zones button.zone-location:hover > span { color: var(--race-accent); text-decoration: underline; }
.zone-kind { font-size: 10px; color: var(--race-muted); white-space: nowrap; }
.zone-kind.imported { color: #d5bd84; }
.inventory-attribution { border-top: 1px solid #ffffff20; margin-top: 20px; padding-top: 15px; font-size: 11px; line-height: 1.7; color: var(--race-muted); }
.inventory-attribution a { color: var(--race-accent); }
.inventory-attribution p:last-child { margin: 6px 0 0; }
.race-message { padding: 35px 20px; text-align: center; color: var(--race-muted); }
.race-message h2 { font-size: 20px; color: #e5e8ee; }
.empty-sources, .muted, .preview-notice { color: var(--race-muted); font-size: 13px; }
.preview-notice { margin: 14px 0; }
.image-credit { text-align: center; color: var(--race-muted); font-size: 10px; margin: 20px 0 0; }
.load-more { display: block; margin: 22px auto 0; }
@keyframes inspector-in { from { opacity: .5; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
@media (min-width: 761px) {
  .race-viewer.locations-layout { display: flex; flex-direction: column; height: var(--race-page-height, calc(100dvh - 44px)); }
  .locations-layout > :not(.race-workspace) { flex-shrink: 0; }
  .race-workspace { flex: 1; min-height: 0; grid-template-rows: minmax(0, 1fr); align-items: stretch; }
  .race-gallery-window, .race-inspector { min-height: 0; height: 100%; margin-top: 0; }
  .race-inspector { position: relative; top: 0; }
  .race-gallery-window ::v-deep > div, .race-inspector ::v-deep > div { display: flex; flex-direction: column; height: 100%; min-height: 0; }
  .gallery-heading, .inspector-heading, .image-credit, .load-more { flex-shrink: 0; }
  .race-gallery, .inspector-scroll { flex: 1; min-height: 0; max-height: none; overflow-y: auto; overscroll-behavior-y: contain; }
  .race-gallery { padding: 3px; grid-auto-rows: max-content; }
}
@media (max-width: 1099px) { .race-workspace { grid-template-columns: minmax(0, 1fr) minmax(300px, 45%); } .reference-note { margin-left: 0; } }
@media (max-width: 760px) { .classic-gallery { max-height: none; } .race-viewer { margin: 0 12px; } .race-workspace { display: flex; flex-direction: column; } .race-inspector { order: -1; position: relative; top: 0; width: 100%; } .race-gallery-window { width: 100%; } .inspector-scroll { max-height: 550px; } .race-heading { align-items: flex-start; flex-direction: column; gap: 10px; } .race-filters { flex-wrap: wrap; } .race-search, .zone-search { flex-basis: 100%; } .race-toolbar { padding: 16px !important; } .race-gallery { grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); } .race-heading h1 { font-size: 21px; } }
@media (prefers-reduced-motion: reduce) { .race-viewer .race-gallery button.race-tile { transition: none; } .race-viewer .race-gallery button.race-tile:hover { transform: none; } .inspector-scroll { animation: none; } }
</style>
