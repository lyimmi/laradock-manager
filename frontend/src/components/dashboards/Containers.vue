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
              <b>{{ item.name }}</b>
              <br />
              <small>{{ item.code }}</small>
            </td>
            <td class="text-center">
              <v-chip :color="item.state === 'Up' ? 'success' : 'warning'">{{item.state}}</v-chip>
            </td>
            <td class="text-center">
              <actions :item="item"></actions>
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
import Actions from '../containers/Actions'

export default {
  name: "Containers",
  components: {Actions},
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