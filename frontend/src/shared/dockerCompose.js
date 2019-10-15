import { mapActions, mapGetters } from "vuex";
import _ from "lodash";

export default {
  data() {
    return {
      dockerComposeStatus: null,
      dockerError: "",
      dotEnv: null,
      dockerVersion: "",
      dockerComposeVersion: "",
      snackbar: false,
      snackbarText: "",
      containersLoading: true,
      containers: [],
      envFilter: "",
      loadingContainer: "",
      dotEnvContents: {},
      dotEnvContentGroups: {},
      dotEnvContentGroupsFiltered: {},
      settingsHasError: false
    };
  },
  computed: {
    ...mapGetters("Containers", ["favoritContainers", "availableContainers"]),
    ...mapGetters("Status", ["appStatus"]),
    ...mapGetters("Settings", ["laradockPath"])
  },
  mounted() {
    this.$root.$on("containersLoading", () => {
      this.containersLoading = true;
    });
    this.$root.$on("containersNotLoading", () => {
      this.containersLoading = false;
    });
  },
  watch: {
    envFilter() {
      this.filterEnvGroup();
    }
  },
  methods: {
    ...mapActions("Status", ["setAppStatus"]),
    ...mapActions("Containers", ["setAvailableContainers", "updateAvailableContainers"]),

    filterEnvGroup: _.debounce(function() {
      if (this.envFilter === null || this.envFilter === "") {
        this.dotEnvContentGroupsFiltered = this.dotEnvContentGroups;
        return;
      }
      let result = {},
        conditions = this.envFilter.split(" ");
      for (let key in this.dotEnvContentGroups) {
        if (this.dotEnvContentGroups.hasOwnProperty(key)) {
          for (let k in this.dotEnvContentGroups[key]) {
            if (
              this.dotEnvContentGroups[key].hasOwnProperty(k) &&
              conditions.every(el =>
                this.dotEnvContentGroups[key][k].field.includes(
                  el.toUpperCase()
                )
              )
            ) {
              if (typeof result[key] === "undefined") {
                result[key] = [];
              }
              result[key].push(this.dotEnvContentGroups[key][k]);
            }
          }
        }
      }
      this.dotEnvContentGroupsFiltered = result;
    }, 350),

    /**
     * Wait for wails to load settings to store
     *
     * @param {*} callback
     * @param {*} i
     */
    waitForSettings(callback, i = 0) {
      if (this.$root.settingsHasError) {
        this.$root.$emit("containersNotLoading");
        return;
      }
      i = typeof i === "undefined" ? 0 : i;
      if (i > 3) {
        this.$root.$emit(
          "showError",
          "Laradock's path is not set. Please set it on the settings page."
        );
        if (!this.$root.settingsHasError) {
          this.$root.$emit("containersNotLoading");
          this.$root.settingsHasError = true;
        }
        return false;
      } else if (
        typeof this.laradockPath === "undefined" ||
        this.laradockPath === ""
      ) {
        setTimeout(() => {
          this.waitForSettings(callback, i + 1);
        }, 100);
      } else if (typeof callback === "function") {
        callback();
      }
    },

    /**
     * Apply laradock path to backend
     *
     * @param {*} laradockPath
     */
    applyLaradockPath(laradockPath) {
      let self = this;
      window.backend.Compose.SetLaradockPath(laradockPath).then(() => {
        self.$root.settingsHasError = false;
      });
    },

    /**
     * Apply terminal path to backend
     *
     * @param {*} terminalPath
     */
    applyTerminalPath(terminalPath) {
      let self = this;
      window.backend.Compose.SetTerminalPath(terminalPath).then(() => {
        self.$root.settingsHasError = false;
      });
    },

    /**
     * Get all available containers
     *
     * @param {*} callback
     */
    getAvailableContainers(callback) {
      let self = this;
      this.waitForSettings(() => {
        self.getContainers(() => {
          window.backend.Compose.GetAvailables().then(result => {
            // Parse response into an object for the table
            let data = JSON.parse(result);
            let containers = [];
            data.forEach(c => {
              if (c !== "") {
                let co = self.containers.find(co => {
                  return co.name === c;
                });
                let isFavorite = false;
                if (this.favoritContainers.findIndex(x => x === c) !== -1) {
                  isFavorite = true;
                }
                containers.push({
                  favorite: isFavorite,
                  name: c,
                  state: typeof co === "undefined" ? "DOWN" : co.state
                });
              }
            });
            // Sort by built containers
            containers.sort((a, b) => {
              if (a.favorite !== b.favorite) {
                return  a.favorite ? -1 : 1;
              }
              return ("" + b.state).localeCompare(a.state);
            });
            self.setAvailableContainers(containers);
            self.$root.$emit("resetRefreshConter");
            if (typeof callback === "function") {
              callback();
            }
          });
        });
      });
    },

    /**
     * Get built containers
     */
    getContainers(callback) {
      this.waitForSettings(() => {
        let self = this;
        self.$root.$emit("containersLoading");
        window.backend.Compose.Get().then(result => {
          if (result.startsWith("Error:")) {
            self.$root.$emit("containersNotLoading");
            self.$root.$emit("showError", result);
            return;
          }
          result = JSON.parse(result);
          let containers = [];
          result.forEach((line, k) => {
            if (k > 1 && line[0] !== "") {
              containers.push({
                code: line[0],
                name: line[0].replace("laradock_", "").replace("_1", ""),
                command: line[1],
                state: line[2],
                ports: line[3]
              });
            }
          });
          self.containers = containers;
          self.$root.$emit("containersNotLoading");
          if (typeof callback === "function") {
            self.$root.$emit("containersLoading");
            callback();
          }
        });
      });
    },

    /**
     * Check if the .env file exists
     */
    checkDotEnv() {
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.CheckDotEnv().then(result => {
          self.dotEnv = result;
          self.setAppStatus({ dotEnv: result });
        });
      });
    },

    /**
     * Check docker executable's version
     */
    checkDockerVersion() {
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.CheckDockerVersion().then(result => {
          let sPartial = "Docker version ";
          let v = result.startsWith(sPartial)
            ? result.replace(sPartial, "")
            : "";
          self.dockerVersion = v;
          let versionNumber = Number(self.dockerVersion.substring(0, 2));
          if (versionNumber < this.$config.docker.dockerMinVersion) {
            self.$root.$emit(
              "showError",
              "Docker executable is too old please update it!"
            );
          }
        });
      });
    },

    /**
     * Check docker-compose executable's version
     */
    checkDockerComposeVersion() {
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.CheckDockerComposeVersion().then(result => {
          let sPartial = "docker-compose version ";
          let v = result.startsWith(sPartial)
            ? result.replace(sPartial, "")
            : "";
          self.dockerComposeVersion = v;
          let versionNumber = Number(self.dockerComposeVersion.substring(0, 3));
          if (versionNumber < this.$config.docker.composeMinVersion) {
            self.$root.$emit(
              "showError",
              "Docker Compose executable is too old please update it!"
            );
          }
        });
      });
    },

    /**
     * Get .env file's contents
     *
     * @param {*} callback
     */
    getDotEnv(callback) {
      this.$root.$emit("containersLoading");
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.DotEnvContent().then(result => {
          let groups = {};
          Object.keys(result).forEach(k => {
            let group = k.split("_", 1)[0];
            if (typeof groups[group] === "undefined") {
              groups[group] = [];
            }
            groups[group].push({
              field: k,
              value: result[k],
              fieldName: k.split(group + "_").pop()
            });
          });
          self.dotEnvContents = result;
          self.dotEnvContentGroups = groups;
          self.dotEnvContentGroupsFiltered = groups;
          this.$root.$emit("containersNotLoading");
          if (typeof callback !== "undefined") {
            callback();
          }
        });
      });
    },

    /**
     * Write .env file's content
     *
     * @param {*} data
     * @param callback
     */
    writeDotEnv(data, callback) {
      this.waitForSettings(() => {
        // let self = this;
        let sData = "";
        Object.keys(data).forEach(e => {
          sData += e + "=" + data[e] + "\n";
        });
        window.backend.Compose.SaveDotEnvContent(sData).then(() => {
          if (typeof callback === "function") {
            callback();
          }
        });
      });
    },

    /**
     * Copy .env file from the example file
     */
    copyFromExample() {
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.CopyEnv().then(result => {
          if (result) {
            self.dotEnv = true;
          } else {
            self.$root.$emit("showError", "Copy failed");
          }
        });
      });
    },

    /**
     * Toggle a continer on/off
     *
     * @param {*} state
     * @param {*} container
     */
    toggleContainer(state, container) {
      this.loadingContainer = container;
      this.waitForSettings(() => {
        let self = this;
        self.$root.$emit("containersLoading");
        window.backend.Compose.Toggle(state, container).then(() => {
          self.$root.$emit("refreshData");
        });
      });
    },

    /**
     * Up a container
     *
     * @param {*} container
     */
    upContainer(container) {
      let self = this;
      this.waitForSettings(() => {
        this.$root.$emit("containersLoading");
        window.backend.Compose.Up(container).then(res => {
          if (res) {
            self.$root.$emit("refreshData");
          } else {
            self.$root.$emit(
              "showError",
              "Container cannot be uppped. Try building it first!"
            );
          }
          self.loadingContainer = "";
        });
      });
    },

    /**
     * Down all containers
     */
    downContainers() {
      let self = this;
      this.waitForSettings(() => {
        this.$root.$emit("containersLoading");
        window.backend.Compose.Down().then(res => {
          if (res) {
            self.$root.$emit("refreshData");
          } else {
            self.$root.$emit("showError", "Container cannot be downed.");
          }
        });
      });
    },

    /**
     * Build a container
     */
    buildContainer(container, force) {
      force = typeof force === "undefined" ? false : true;
      let self = this;
      this.waitForSettings(() => {
        this.$root.$emit("containersLoading");
        window.backend.Compose.Build(container, force).then(res => {
          if (res) {
            self.$root.$emit("refreshData");
          } else {
            self.$root.$emit(
              "showError",
              "Container cannot be uppped. Try building it first!"
            );
          }
          self.loadingContainer = "";
        });
      });
    },

    /**
     * Call docker-compose exec on a container
     *
     * @param {*} container
     * @param {*} user
     * @param {*} callback
     */
    execContainer(container, user, callback) {
      window.backend.Compose.Exec(container, user).then(() => {
        if (typeof callback === "function") {
          callback();
        }
      });
    },

    logsContainer(container, callback) {
      window.backend.Compose.Logs(container).then(() => {
        if (typeof callback === "function") {
          callback();
        }
      });
    },

    /**
     * Stop container exec
     *
     * @param {*} callback
     */
    stopExecContiner(callback) {
      window.backend.Compose.StopExec().then(res => {
        if (res === "disconnected" && typeof callback === "function") {
          callback();
        }
      });
    }
  }
};
