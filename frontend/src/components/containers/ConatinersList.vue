<template>
  <v-card class="mx-auto" :loading="containersLoading">
    <v-card-title>
      Containers &nbsp;
      <v-menu offset-y>
        <template v-slot:activator="{ on }">
          <v-btn color="default" x-small dark v-on="on">Mass actions</v-btn>
        </template>
        <v-list dense>
          <v-list-item @click="massBuild">
            <v-list-item-icon>
              <v-icon>mdi-progress-wrench</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Build</v-list-item-content>
          </v-list-item>
          <v-list-item @click="massBuild(true)">
            <v-list-item-icon>
              <v-icon>mdi-progress-wrench</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Build without cache</v-list-item-content>
          </v-list-item>
          <v-list-item @click="massUp">
            <v-list-item-icon>
              <v-icon>mdi-arrow-up-bold</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Up</v-list-item-content>
          </v-list-item>
          <v-list-item @click="massToggle('start')">
            <v-list-item-icon>
              <v-icon>mdi-play</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Start</v-list-item-content>
          </v-list-item>
          <v-list-item @click="massToggle('stop')">
            <v-list-item-icon>
              <v-icon>mdi-stop</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Stop</v-list-item-content>
          </v-list-item>
          <v-list-item @click="downContainers()">
            <v-list-item-icon>
              <v-icon>mdi-arrow-down-bold</v-icon>
            </v-list-item-icon>
            <v-list-item-content>Down all</v-list-item-content>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-spacer></v-spacer>
      <v-text-field
        class="mt-0 pt-0"
        v-model="search"
        append-icon="mdi-database-search"
        label="Search"
        single-line
        hide-details
      ></v-text-field>
    </v-card-title>
    <v-card-text>
      <v-data-table
        :headers="headers"
        :single-select="false"
        item-key="name"
        :items="availableContainers"
        :search="search"
        v-model="selectedContainers"
        show-select
      >
        <template v-slot:item.favorite="{ item }">
          <v-btn text icon :color="item.favorite ? 'yellow' : 'grey'" @click="toggleFavorite(item)">
            <v-icon>mdi-star</v-icon>
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
        { text: "Actions", align: "left", value: "action", sortable: false }
      ],
      selectedContainers: []
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
    pluckCotnainers(containers) {
      return Object.keys(containers)
        .map(f => containers[f].name)
        .join("|");
    },
    massBuild(cache = false) {
      this.buildContainer(this.pluckCotnainers(this.selectedContainers), cache);
    },
    massToggle(state) {
      this.toggleContainer(
        state,
        this.pluckCotnainers(this.selectedContainers)
      );
    },
    massUp() {
      this.upContainer(this.pluckCotnainers(this.selectedContainers));
    },
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