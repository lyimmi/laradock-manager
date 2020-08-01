<template>
  <div v-if="hasDotEnv">
    <v-card>
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

    <v-alert
      v-if="hasContainersUp === null"
      icon="mdi-database-sync"
      prominent
      text
      type="info"
      transition="fade-transition"
    >Checking Containers please wait...</v-alert>
    <v-alert
      v-else-if="hasContainersUp === true && stats.length === 0"
      icon="mdi-database-sync"
      prominent
      text
      type="info"
      transition="fade-transition"
    >Waiting for container statistics, it will take a moment.</v-alert>
    <v-alert
      v-else-if="hasContainersUp === false"
      icon="mdi-database-remove"
      prominent
      text
      type="warning"
      transition="fade-transition"
    >No container is up, please start one first!</v-alert>
  </div>
</template>
<script>
import DockerMixin from "@/shared/dockerMixin";
export default {
  name: "stats",
  mixins: [DockerMixin],
  data: () => {
    return {
      hasContainersUp: null,
      stats: [],
    };
  },
  created() {
    // Attach an event handler to <div>
    window.wails.Events.On("stats", (stats) => {
      this.stats = JSON.parse(atob(stats));
    });
  },
  destroyed() {
    this.stopStats();
  },
  mounted() {
    this.checkDotEnv().then((resenv) => {
      if (!resenv) return;
      this.hasRunning().then((res) => {
        this.hasContainersUp = res;
        if (res) {
          this.getStats();
        }
      });
    });
  },
};
</script>