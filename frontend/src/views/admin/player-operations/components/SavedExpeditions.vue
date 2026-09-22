<template>
  <section class="saved-expeditions" data-testid="saved-expeditions" aria-label="Saved expeditions">
    <div class="spire-editor-section-heading">
      <div>
        <span class="spire-editor-section-kicker">DZManager history</span>
        <h3>Saved expeditions</h3>
      </div>
      <b-button size="sm" variant="outline-secondary" :disabled="loading" @click="load">
        <i :class="loading ? 'fa fa-spinner fa-spin' : 'fa fa-refresh'" class="mr-1"></i>Refresh expeditions
      </b-button>
    </div>
    <p v-if="loading" role="status">Loading saved expeditions…</p>
    <div v-else-if="error" class="alert alert-warning" role="alert">{{ error }}</div>
    <template v-else-if="snapshot">
      <p class="saved-expeditions-note">{{ snapshot.message }}</p>
      <template v-if="snapshot.available">
        <p v-if="!snapshot.configured_enabled" class="alert alert-warning">
          Saved expedition rejoining is disabled in the configured world ruleset.
        </p>
        <div class="saved-expeditions-summary">
          <span>{{ eligibleCount }} eligible pending server check · {{ entries.length }} saved</span>
          <label><input v-model="eligibleOnly" type="checkbox"> Show eligible only</label>
        </div>
        <div v-if="visibleEntries.length" class="saved-expeditions-table">
          <table class="table table-sm">
            <thead><tr><th>Expedition</th><th>Members</th><th>Time remaining</th><th>Rejoin status</th></tr></thead>
            <tbody>
              <tr v-for="entry in visibleEntries" :key="entry.uuid" :data-dz-id="entry.dz_id">
                <td>
                  <strong>{{ entry.name }}</strong>
                  <small>DZ #{{ entry.dz_id }} · instance {{ entry.instance_id }}</small>
                  <small>Zone #{{ entry.zone_id }} · version {{ entry.zone_version }} · season {{ entry.season_id }}</small>
                  <small>Leader: {{ entry.leader || 'Unassigned' }}</small>
                </td>
                <td>{{ entry.members }} / {{ entry.max_members }}</td>
                <td :title="expiresAt(entry.expires)">{{ remaining(entry.expires) }}<small>{{ expiresAt(entry.expires) }}</small></td>
                <td><span class="saved-expeditions-status" :class="{ eligible: entry.database_eligible }">{{ statusLabel(entry.status) }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
        <p v-else class="saved-expeditions-empty">
          {{ eligibleOnly && entries.length ? 'No saved expeditions currently pass the database checks.' : 'No unexpired saved expeditions remain for this character’s latest copies.' }}
        </p>
        <small class="saved-expeditions-note">Players rejoin through /dzmanager in game, then use the normal expedition entrance. Refresh to update this snapshot.</small>
      </template>
    </template>
  </section>
</template>

<script>
import { SpireApi } from '@/app/api/spire-api'

const statuses = {
  server_check: 'Eligible pending server check',
  invalid: 'Invalid saved expedition identity',
  disabled: 'Rejoining disabled by configured rules',
  expired: 'Expired',
  already_member: 'Already a member',
  in_expedition: 'Leave current expedition first',
  in_instance: 'Leave current instance first (saved location)',
  locked: 'Expedition locked',
  full: 'Expedition full',
  season: 'Season does not match configured season',
  lockout: 'Conflicting replay or event lockout'
}

export default {
  name: 'SavedExpeditions',
  props: { characterId: { type: Number, required: true } },
  data () {
    return { loading: false, error: '', snapshot: null, eligibleOnly: false, requestID: 0 }
  },
  computed: {
    entries () { return this.snapshot && Array.isArray(this.snapshot.entries) ? this.snapshot.entries : [] },
    eligibleCount () { return this.entries.filter(entry => entry.database_eligible).length },
    visibleEntries () { return this.eligibleOnly ? this.entries.filter(entry => entry.database_eligible) : this.entries }
  },
  watch: {
    characterId: { immediate: true, handler () { this.eligibleOnly = false; this.load() } }
  },
  beforeDestroy () { this.requestID++ },
  methods: {
    async load () {
      const requestID = ++this.requestID
      this.loading = true
      this.error = ''
      this.snapshot = null
      try {
        const response = await SpireApi.v1().get(`/player-operations/character/${this.characterId}/saved-expeditions`)
        if (requestID === this.requestID) this.snapshot = response.data
      } catch (error) {
        if (requestID === this.requestID) this.error = 'Unable to load saved expeditions. Refresh to try again.'
      } finally {
        if (requestID === this.requestID) this.loading = false
      }
    },
    statusLabel (status) { return statuses[status] || 'Server check required' },
    expiresAt (expires) { return new Date(Number(expires) * 1000).toLocaleString() },
    remaining (expires) {
      const seconds = Math.max(0, Number(expires) - Number(this.snapshot.server_time))
      if (seconds < 60) return `${seconds}s`
      const minutes = Math.floor(seconds / 60)
      return minutes < 60 ? `${minutes}m` : `${Math.floor(minutes / 60)}h ${minutes % 60}m`
    }
  }
}
</script>

<style scoped>
.saved-expeditions { margin-top: 24px; }
.saved-expeditions-note, .saved-expeditions-empty { color: #b3b8c5; }
.saved-expeditions-summary { display: flex; justify-content: space-between; gap: 12px; flex-wrap: wrap; margin: 12px 0; }
.saved-expeditions-summary label { margin: 0; }
.saved-expeditions-summary input { margin-right: 6px; }
.saved-expeditions-table { overflow-x: auto; }
.saved-expeditions-table table { margin-bottom: 12px; }
.saved-expeditions-table th, .saved-expeditions-table td { padding: 10px; border-color: rgba(185, 160, 105, .2); }
.saved-expeditions-table small { display: block; color: #b3b8c5; }
.saved-expeditions-status { color: #efc77b; }
.saved-expeditions-status.eligible { color: #94d4aa; }
</style>
