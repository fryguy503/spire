<template>
  <div class="mt-3" data-testid="item-kinbound-editor">
    <h4 class="eq-header">Kinbound</h4>
    <p>Choose which NO TRADE items may be shared between characters linked to the same forum account before permanent attunement.</p>
    <b-alert v-if="error" show variant="danger">{{ error }}</b-alert>
    <b-alert v-if="notification" show variant="success">{{ notification }}</b-alert>
    <div v-if="loading">Loading Kinbound setting...</div>
    <b-alert v-else-if="state && !state.supported" show variant="info">{{ state.message }}</b-alert>
    <template v-else-if="state && state.supported">
      <label for="kinbound">Kinbound selection</label>
      <select
        id="kinbound"
        v-model.number="selection"
        class="form-control"
        :disabled="saving"
        :title="fieldDescription"
        @change="markModified"
      >
        <option :value="0">0) Unlisted (original item behavior)</option>
        <option :value="1" :disabled="!state.can_always">1) Always Kinbound</option>
        <option :value="2">2) Never Kinbound</option>
      </select>
      <b-alert v-if="state.blocked_reason" show variant="warning" class="mt-3">
        {{ state.blocked_reason }}
      </b-alert>
      <p class="mt-3 mb-2">This selection uses saved item flags. After changing No Drop or No Transfer, save the item and refresh this panel.</p>
      <p class="mb-3">Expansion unlocks never enable Kinbound. The server must have Kinbound enabled and reload its item data before this change takes effect. Existing permanent attunement is preserved.</p>
      <div class="d-flex align-items-center">
        <button
          id="save-item-kinbound"
          :disabled="saving || !modified || (selection === 1 && !state.can_always)"
          @click="save"
        >{{ saving ? 'Saving...' : 'Save Kinbound' }}</button>
        <span v-if="modified" class="ml-3">Unsaved Kinbound selection</span>
      </div>
    </template>
    <button class="mt-3" :disabled="loading || saving" @click="load">Refresh Kinbound</button>
  </div>
</template>

<script>
import {SpireApi} from '../../../app/api/spire-api';
import {Items} from '../../../app/items';
import {EditFormFieldUtil} from '../../../app/forms/edit-form-field-util';

export default {
  name: 'ItemKinboundEditor',
  props: { itemId: { type: [Number, String], required: true } },
  data() {
    return { state: null, selection: 0, loading: false, saving: false, error: '', notification: '', requestID: 0 };
  },
  computed: {
    modified() { return this.state && this.state.supported && this.selection !== this.state.kinbound; },
    fieldDescription() { return Items.getFieldDescriptions().kinbound; }
  },
  watch: {
    itemId: { immediate: true, handler() { this.load(); } }
  },
  beforeDestroy() { this.requestID++; },
  methods: {
    async load() {
      const requestID = ++this.requestID;
      this.error = '';
      this.notification = '';
      this.state = null;
      this.loading = true;
      this.saving = false;
      try {
        const {data} = await SpireApi.v1().get(`item-kinbound/${this.itemId}`);
        if (requestID !== this.requestID) return;
        this.state = data;
        this.selection = Number(data.kinbound || 0);
        this.$nextTick(this.markModified);
      } catch (error) {
        if (requestID === this.requestID) this.error = this.errorMessage(error);
      } finally {
        if (requestID === this.requestID) this.loading = false;
      }
    },
    markModified() {
      if (this.modified) EditFormFieldUtil.setFieldModifiedById('kinbound');
      else EditFormFieldUtil.clearFieldModifiedById('kinbound');
    },
    async save() {
      if (!this.modified || this.saving) return;
      const requestID = this.requestID;
      this.error = '';
      this.notification = '';
      this.saving = true;
      try {
        const {data} = await SpireApi.v1().patch(`item-kinbound/${this.itemId}`, {
          kinbound: this.selection,
          expected_kinbound: this.state.kinbound
        });
        if (requestID !== this.requestID) return;
        this.state = data;
        this.selection = Number(data.kinbound);
        this.notification = 'Kinbound setting saved. Reload the server item data to apply it.';
        this.markModified();
      } catch (error) {
        if (requestID === this.requestID) this.error = this.errorMessage(error);
      } finally {
        if (requestID === this.requestID) this.saving = false;
      }
    },
    errorMessage(error) {
      return error.response && error.response.data && error.response.data.error || 'Unable to load or save Kinbound.';
    }
  }
};
</script>
