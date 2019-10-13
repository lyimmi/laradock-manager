<template>
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
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details></v-text-field>
    </v-card-title>
    <v-card-text>
      <v-data-table :headers="headers" :items="availableContainers" :search="search">
        <template v-slot:item.favorite="{ item }">
          <v-btn text icon :color="item.favorite ? 'yellow' : 'grey'" @click="toggleFavorite(item)">
            <v-icon>star</v-icon>
          </v-btn>
        </template>
        <template v-slot:item.state="{ item }">
          <v-chip small :color="stateColor(item.state)">{{item.state}}</v-chip>
        </template>
        <template v-slot:item.action="{ item }">
          <actions :item="item"></actions>
        </template>
      </v-data-table>

      <exec></exec>
    </v-card-text>
  </v-card>
</template>

<script>
import Actions from "../containers/Actions";
import Exec from "../containers/Exec";
import dockerCompose from "../../shared/dockerCompose";
import { mapActions, mapGetters } from "vuex";

export default {
  name: "containers-list",
  components: { Exec, Actions },
  mixins: [dockerCompose],
  data() {
    return {
      loading: true,
      search: "",
      headers: [
        {
          text: "Favorite",
          align: "center",
          sortable: true,
          value: "favorite"
        },
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
  computed: {
    ...mapGetters("Containers", ["favoritContainers", "availableContainers"])
  },
  methods: {
    ...mapActions("Containers", [
      "addFavoriteContiner",
      "updateAvailableContainer",
      "removeFavoriteContiner"
    ]),
    toggleFavorite(item) {
      let isFavorite = true;
      if (this.favoritContainers.findIndex(x => x === item.name) === -1) {
        this.addFavoriteContiner(item.name);
        isFavorite = true;
      } else {
        this.removeFavoriteContiner(item.name);
        isFavorite = false;
      }
      this.updateAvailableContainer({ item: item, isFavorite: isFavorite });
    },
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