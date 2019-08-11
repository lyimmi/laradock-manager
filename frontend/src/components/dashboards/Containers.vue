<template>
  <v-card class="mx-auto" :loading="containersLoading">
    <v-card-title>Containers</v-card-title>
    <v-card-text>
      <v-simple-table fixed-header v-if="containers.length > 0">
        <thead>
          <tr>
            <th class="text-left">Name</th>
            <th class="text-center">State</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in containers" :key="item.name">
            <td>
              <b>{{ item.code }}</b>
              <br />
              <small>{{ item.name }}</small>
            </td>
            <td class="text-center">
              <v-chip :color="item.state === 'Up' ? 'success' : 'warning'">{{item.state}}</v-chip>
            </td>
            <td>
              <v-btn
                small
                color="warning"
                v-if="item.state === 'Up'"
                :loading="containersLoading"
                @click="toggleContainer('stop', item.code)"
              >
                <v-icon>stop</v-icon>
              </v-btn>
              <v-btn
                small
                color="success"
                v-else
                :loading="containersLoading"
                @click="toggleContainer('start', item.code)"
              >
                <v-icon>play_arrow</v-icon>
              </v-btn>
            </td>
          </tr>
        </tbody>
      </v-simple-table>
      <span v-else>No containers found</span>
    </v-card-text>
  </v-card>
</template>

<script>
import dockerCompose from "../../shared/dockerCompose";

export default {
  name: "Containers",
  mixins: [dockerCompose],
  mounted() {
    this.getContainers();
    this.$root.$on("refreshData", () => {
      this.getContainers();
    });
  }
};
</script>

<style scoped>
</style>