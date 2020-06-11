<template>
  <div>
    <v-alert
      v-if="stats.length === 0"
      icon="mdi-database-sync"
      prominent
      text
      type="info"
    >No container data available, plase wait.</v-alert>
    <v-card v-else>
      <v-simple-table>
        <template v-slot:default>
          <thead>
            <tr>
              <th class="text-left">Container</th>
              <th class="text-left">CPU usage</th>
              <th class="text-left">Memory usage</th>
              <th class="text-left">Memory percent</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="stat in stats" :key="stat.name">
              <td>{{ stat.name }}</td>
              <td>
                <v-progress-linear v-model="stat.cpu_perc" height="25" reactive rounded>
                  <strong>{{ stat.cpu_perc_string }}</strong>
                </v-progress-linear>
              </td>
              <td>{{ stat.memory_usage }}</td>
              <td>
                <v-progress-linear v-model="stat.memory_percent" height="25" reactive rounded>
                  <strong>{{ stat.memory_percent_string }}</strong>
                </v-progress-linear>
              </td>
            </tr>
          </tbody>
        </template>
      </v-simple-table>
    </v-card>
  </div>
</template>
<script>
import DockerMixin from "@/shared/dockerMixin";
export default {
  name: "stats",
  mixins: [DockerMixin],
  data: () => {
    return {
      stats: []
    };
  },
  created() {
    // Attach an event handler to <div>
    window.wails.Events.On("stats", stats => {
      this.stats = JSON.parse(atob(stats));
    });
  },
  destroyed() {
    this.stopStats();
  },
  mounted() {
    this.getStats();
  }
};
</script>