<template>
  <v-card>
    <v-skeleton-loader class="mx-auto" type="table" :loading="availableContainers.length === 0">
      <v-card-title>
        <v-text-field
          v-model="search"
          append-icon="mdi-search"
          label="Search"
          single-line
          hide-details
        ></v-text-field>
      </v-card-title>
      <v-data-table
        :loading="containersLoading"
        sortBy="state"
        :sortDesc="true"
        :search="search"
        v-model="selected"
        :headers="headers"
        :items="availableContainers"
        item-key="name"
        show-select
        class="containers"
      >
        <template v-slot:item.favorite="{ item }">
          <v-btn icon @click="toggleFavorite(item)">
            <v-icon v-if="item.favorite" color="yellow">mdi-star</v-icon>
            <v-icon v-else>mdi-star-outline</v-icon>
          </v-btn>
        </template>
        <template v-slot:item.state="{ item }">
          <v-chip :color="getColor(item.state)" small>{{ item.state }}</v-chip>
        </template>
        <template v-slot:item.actions="{item}">
          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                small
                v-show="item.state === 'Down'"
                @click="upContainer(item.name)"
                v-bind="attrs"
                v-on="on"
              >
                <v-icon>mdi-arrow-up</v-icon>
              </v-btn>
            </template>
            <span>Up container</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                small
                v-show="item.state === 'Stopped'"
                @click="startContainer(item.name)"
                v-bind="attrs"
                v-on="on"
              >
                <v-icon>mdi-play</v-icon>
              </v-btn>
            </template>
            <span>Start container</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                small
                v-show="item.state === 'Up'"
                @click="stopContainer(item.name)"
                v-bind="attrs"
                v-on="on"
              >
                <v-icon>mdi-stop</v-icon>
              </v-btn>
            </template>
            <span>Stop container</span>
          </v-tooltip>

          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn small class="ml-5 mr-1" @click="buildContainer(item.name)" v-bind="attrs" v-on="on">
                <v-icon>mdi-docker</v-icon>
              </v-btn>
            </template>
            <span>Build container (cached)</span>
          </v-tooltip>

          <v-tooltip top>
            <template v-slot:activator="{ on, attrs }">
              <v-btn small @click="buildContainer(item.name, true)" v-bind="attrs" v-on="on">
                <v-icon color="red">mdi-docker</v-icon>
              </v-btn>
            </template>
            <span>Build container (without cache)</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-skeleton-loader>
  </v-card>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import DockerMixin from "@/shared/dockerMixin";
export default {
  name: "container-list",
  mixins: [DockerMixin],
  data: () => {
    return {
      confirmDialog: false,
      search: "",
      selected: [],
      headers: [
        {
          text: "Favorite",
          value: "favorite",
          sortable: true,
          align: "center"
        },
        {
          text: "Name",
          sortable: true,
          value: "name",
          align: "center"
        },
        { text: "State", value: "state", sortable: true, align: "center" },
        { text: "Actions", value: "actions", sortable: false, align: "center" }
      ]
    };
  },
  mounted() {
    if (this.laradockPath !== "") {
      this.loadConstainers();
    }
  },
  computed: {
    ...mapGetters("Settings", ["laradockPath", "containerPrefix"])
  },
  methods: {
    ...mapActions("Containers", [
      "addFavoriteContiner",
      "removeFavoriteContiner"
    ]),
    getColor(state) {
      if (state === "Up") {
        return "green";
      } else if (state === "down") {
        return "red";
      } else if (state === "Stopped") {
        return "orange";
      } else {
        return "default";
      }
    },
    toggleFavorite(container) {
      const promise = new Promise(resolve => {
        if (container.favorite) {
          this.removeFavoriteContiner(container.name).then(() => {
            resolve(true);
          });
        } else {
          this.addFavoriteContiner(container.name).then(() => {
            resolve(true);
          });
        }
      });

      promise.then(() => {
        this.loadConstainers();
      });
    }
  }
};
</script>