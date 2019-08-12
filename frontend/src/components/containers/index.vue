<template>
  <v-container fluid grid-list-md id="containersContainer">
    <exec></exec>
    <v-layout>
      <v-flex xs12>
        <v-card class="mx-auto" :loading="containersLoading">
          <v-card-title>Containers</v-card-title>
          <v-card-text>
            <v-simple-table fixed-header dense v-if="containers.length > 0">
              <thead>
                <tr>
                  <th class="text-left">Name</th>
                  <th class="text-center">State</th>
                  <th class="text-center">Action</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in availableContainers" :key="item.name">
                  <td>{{ item.name }}</td>
                  <td class="text-center">
                    <v-chip small :color="stateColor(item.state)">{{item.state}}</v-chip>
                  </td>
                  <td class="text-center">
                    <v-btn
                      icon
                      small
                      v-if="item.state === 'Up'"
                      class="ma-2"
                      dark
                      v-on:click="$root.$emit('execContiner', item.name)"
                    >
                      <v-icon :size="14">fas fa-terminal</v-icon>
                    </v-btn>

                    <v-btn
                      :loading="loadingContainer == item.name"
                      small
                      icon
                      :size="14"
                      v-if="item.state === 'Up'"
                      @click="toggleContainer('stop', item.name, true)"
                    >
                      <v-icon>stop</v-icon>
                    </v-btn>

                    <v-btn
                      :loading="loadingContainer == item.name"
                      small
                      icon
                      :size="14"
                      v-else-if="item.state !== 'DOWN'"
                      @click="toggleContainer('start', item.name, true)"
                    >
                      <v-icon>play_arrow</v-icon>
                    </v-btn>

                    <v-menu offset-y>
                      <template v-slot:activator="{ on }">
                        <v-btn small icon v-on="on" :loading="loadingContainer == item.name">
                          <v-icon :size="14">build</v-icon>
                        </v-btn>
                      </template>
                      <v-list>
                        <v-list-item @click="buildContainer(item.name, false, true)">
                          <v-list-item-title>Build with cache</v-list-item-title>
                        </v-list-item>
                        <v-list-item @click="buildContainer(item.name, true, true)">
                          <v-list-item-title>Build without cache</v-list-item-title>
                        </v-list-item>
                      </v-list>
                    </v-menu>
                  </td>
                </tr>
              </tbody>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import Exec from "../containers/Exec";
import dockerCompose from "../../shared/dockerCompose";

export default {
  name: "index",
  components: { Exec },
  mixins: [dockerCompose],
  data() {
    return {
      loading: true
    };
  },
  mounted() {
    this.getAvailableContainers(() => {
      this.containersLoading = false;
    });
    this.$root.$on("refreshData", () => {
      this.containersLoading = true;
      this.getAvailableContainers(() => {
        this.containersLoading = false;
      });
    });
  },
  methods: {
    stateColor(state) {
      if (state === "Up") {
        return "success";
      } else if (state === "DOWN") {
        return "error";
      } else {
        return "warning";
      }
    }
  }
};
</script>

<style scoped>
</style>