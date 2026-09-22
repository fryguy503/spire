<template>
  <div class="row" style="max-height: 85vh; overflow-y: scroll; top: -15px; position: inherit">
    <div class="col-lg-6">

      <div class="row">
        <div class="col-sm-6 col-lg-6">
          <dashboard-process-counts v-if="canRead('eqemuserver/server-stats')"/>

          <dashboard-networking-info v-if="canRead('admin/serverconfig')"/>
          <dashboard-system-info-v2 v-if="canRead('admin/system/host')"/>
        </div>
        <div class="col-sm-6 col-lg-6">
          <dashboard-cpu-info v-if="canRead('admin/system/cpu')"/>
        </div>
      </div>

    </div>

    <!-- Right side -->
    <div class="col-lg-6">
      <players-online-component v-if="canRead('eqemuserver/client-list')"/>
    </div>
  </div>
</template>

<script>
import {canReadAdminApi} from "@/app/user/server-admin-access";
import ServerProcessButtonComponent from "@/views/admin/components/ServerProcessButtonComponent";
import DashboardProcessCounts       from "@/views/admin/components/DashboardProcessCounts";
import DashboardCpuInfo             from "@/views/admin/components/DashboardCpuInfo";
import DashboardSystemInfo          from "@/views/admin/components/DashboardSystemInfo";
import DashboardCounter             from "@/views/admin/components/DashboardCounter.vue";
import PlayersOnlineComponent       from "@/views/admin/components/PlayersOnlineComponent.vue";
import DashboardSystemInfoV2        from "@/views/admin/components/DashboardSystemInfoV2.vue";
import DashboardNetworkingInfo      from "@/views/admin/components/DashboardNetworkingInfo.vue";

export default {
  methods: { canRead: canReadAdminApi },
  components: {
    DashboardNetworkingInfo,
    DashboardSystemInfoV2,
    PlayersOnlineComponent,
    DashboardCounter,
    DashboardSystemInfo,
    DashboardCpuInfo,
    DashboardProcessCounts,
    ServerProcessButtonComponent
  },
}
</script>
