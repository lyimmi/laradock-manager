<template>
  <v-container fluid grid-list-md id="containersContainer">
    <exec></exec>
    <v-layout>
      <v-flex xs12>
        <v-card class="mx-auto" :loading="containersLoading">
          <v-card-title>
            Containers &nbsp;
            <v-tooltip bottom>
              <template v-slot:activator="{ on }">
                <v-btn icon small @click="downContainers()" v-on="on">
                  <v-icon>fas fa-arrow-down</v-icon>
                </v-btn>
              </template>
              <span>Down all containers: docker-compose down</span>
            </v-tooltip>
            <v-spacer></v-spacer>
            <v-text-field
              v-model="search"
              append-icon="search"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
          </v-card-title>
          <v-card-text>
            <v-data-table :headers="headers" :items="availableContainers" :search="search">
              <template v-slot:item.state="{ item }">
                <v-chip small :color="stateColor(item.state)">{{item.state}}</v-chip>
              </template>
              <template v-slot:item.action="{ item }">
                <actions :item="item"></actions>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import Actions from "../containers/Actions";
import Exec from "../containers/Exec";
import dockerCompose from "../../shared/dockerCompose";

export default {
  name: "index",
  components: { Exec, Actions },
  mixins: [dockerCompose],
  data() {
    return {
      loading: true,
      search: "",
      headers: [
        {
          text: "Name",
          align: "left",
          sortable: true,
          value: "name"
        },
        { text: "State", align: "center", sortable: true, value: "state" },
        { text: "Actions", align: "center", value: "action", sortable: false }
      ]
    };
  },
  mounted() {
    this.$root.$emit("containersLoading");
    this.getAvailableContainers(() => {
      this.$root.$emit("containersNotLoading");
    });

    this.$root.$on("refreshData", () => {
      this.$root.$emit("containersLoading");
      this.getAvailableContainers(() => {
        this.$root.$emit("containersNotLoading");
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