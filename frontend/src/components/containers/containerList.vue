<template>
  <v-card v-if="hasDotEnv">
    <v-skeleton-loader
      class="mx-auto"
      type="table"
      :loading="availableContainers === null || availableContainers.length === 0"
    >
      <v-card-title>
        <v-row>
          <v-col class="py-0" cols="12" sm="4">
            <v-menu offset-y>
              <template v-slot:activator="{ on, attrs }">
                <v-btn v-bind="attrs" v-on="on" small>Actions</v-btn>
              </template>
              <v-list>
                <v-list-item
                  v-for="(item, index) in massActions"
                  :key="index"
                  @click="runMassAction(item)"
                >
                  <v-list-item-title>{{ item.name }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
          <v-col class="py-0" cols="12" sm="8">
            <v-text-field
              v-model="search"
              class="pa-0 ma-0"
              append-icon="mdi-search"
              label="Search"
              single-line
              hide-details
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-title>
      <v-data-table
        calculate-widths
        :loading="containersLoading"
        sortBy="favorite"
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
          <v-row>
            <v-col cols="12" sm="6">
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
                  <v-btn
                    class="ml-1"
                    small
                    v-show="item.state === 'Up'"
                    @click="logContainer(item.name)"
                    v-bind="attrs"
                    v-on="on"
                  >
                    <v-icon>mdi-view-list</v-icon>
                  </v-btn>
                </template>
                <span>Open container logs</span>
              </v-tooltip>

              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    class="ml-1"
                    small
                    v-show="item.state === 'Up'"
                    @click="execDialog = true; executableContainer=item.name"
                    v-bind="attrs"
                    v-on="on"
                  >
                    <v-icon>mdi-console</v-icon>
                  </v-btn>
                </template>
                <span>Exec container</span>
              </v-tooltip>
            </v-col>
            <v-col cols="12" sm="6">
              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    small
                    class="mr-1"
                    @click="buildContainer(item.name)"
                    v-bind="attrs"
                    v-on="on"
                  >
                    <v-icon>mdi-account-hard-hat</v-icon>
                  </v-btn>
                </template>
                <span>Build container (cached)</span>
              </v-tooltip>

              <v-tooltip top>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn small @click="buildContainer(item.name, true)" v-bind="attrs" v-on="on">
                    <v-icon color="red">mdi-account-hard-hat</v-icon>
                  </v-btn>
                </template>
                <span>Build container (without cache)</span>
              </v-tooltip>
            </v-col>
          </v-row>
        </template>
      </v-data-table>
    </v-skeleton-loader>
    <v-fab-transition>
      <v-btn
        @click="loadConstainers"
        fab
        dark
        bottom
        right
        fixed
        small
        color="primary"
        class="v-btn--example"
      >
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
    </v-fab-transition>
    <v-dialog v-model="execDialog" persistent max-width="350">
      <v-card>
        <v-card-title>
          <span class="headline">Select user</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12">
                <v-select
                  v-model="executableUser"
                  :items="['root', 'laradock']"
                  label="User"
                  required
                ></v-select>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="execDialog = false">Cancel</v-btn>
          <v-btn
            color="blue darken-1"
            text
            @click="execContainer(); execDialog = false"
            :disabled="executableUser === ''"
          >Exec</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
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
      execDialog: false,
      executableContainer: "",
      executableUser: "",
      massActions: [
        { value: "up", name: "Up Containers", action: "upContainer" },
        { value: "start", name: "Start Containers", action: "startContainer" },
        { value: "stop", name: "Stop Containers", action: "stopContainer" },
      ],
      headers: [
        {
          text: "Favorite",
          value: "favorite",
          sortable: true,
          align: "center",
        },
        {
          text: "Name",
          sortable: true,
          value: "name",
          align: "center",
        },
        { text: "State", value: "state", sortable: true, align: "center" },
        {
          text: "Actions",
          value: "actions",
          sortable: false,
          align: "center",
          width: 400,
        },
      ],
    };
  },
  mounted() {
    this.loadConstainers();
  },
  computed: {
    ...mapGetters("Settings", ["laradockPath", "containerPrefix"]),
  },
  methods: {
    ...mapActions("Containers", [
      "addFavoriteContiner",
      "removeFavoriteContiner",
    ]),
    runMassAction(item) {
      if (typeof this[item.action] === "function") {
        this[item.action](
          this.selected
            .map(function (elem) {
              return elem.name;
            })
            .join("|")
        );
      }
    },
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
      const promise = new Promise((resolve) => {
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
    },
  },
};
</script>