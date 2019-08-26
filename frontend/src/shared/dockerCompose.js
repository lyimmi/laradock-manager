import { mapActions, mapGetters } from "vuex";

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
      containers: [],
      containersLoading: true,
      availableContainers: [],
      loadingContainer: "",
      dotEnvContents: {},
      dotEnvContentGroups: {}
    };
  },
  computed: {
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
  methods: {
    ...mapActions("Status", ["setAppStatus"]),

    /**
     * Wait for wails to load settings to store
     *
     * @param {*} callback
     * @param {*} i
     */
    waitForSettings(callback, i) {
      i = typeof i === "undefined" ? 0 : i;
      if (i > 3) {
        this.$root.$emit(
          "showError",
          "Laradock's path is not set. Please set it on the settings page."
        );
        return false;
      } else if (
        typeof this.laradockPath === "undefined" ||
        this.laradockPath === ""
      ) {
        setTimeout(() => {
          this.waitForSettings(callback, i + 1);
        }, 350);
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
      window.backend.Compose.SetLaradockPath(laradockPath);
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
                  return co.code === c;
                });
                containers.push({
                  name: c,
                  state: typeof co === "undefined" ? "DOWN" : co.state
                });
              }
            });
            // Sort by built containers
            containers.sort((a, b) => {
              return ("" + b.state).localeCompare(a.state);
            });
            self.availableContainers = containers;
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
            this.$root.$emit("showError", result);
            return;
          }
          result = JSON.parse(result);
          let containers = [];
          result.forEach((line, k) => {
            if (k > 1 && line[0] !== "") {
              containers.push({
                name: line[0],
                code: line[0].replace("laradock_", "").replace("_1", ""),
                command: line[1],
                state: line[2],
                ports: line[3]
              });
            }
          });
          self.containers = containers;
          self.$root.$emit("containersNotLoading");
          self.$root.$emit("resetRefreshConter");
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
        });
      });
    },

    /**
     * Get .env file's contents
     *
     * @param {*} callback
     */
    getDotEnv(callback) {
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
     */
    writeDotEnv(data) {
      this.waitForSettings(() => {
        // let self = this;
        let sData = "";
        Object.keys(data).forEach(e => {
          sData += e + "=" + data[e] + "\n";
        });
        window.backend.Compose.SaveDotEnvContent(sData).then(() => {});
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
          self.loadingContainer = "";
        });
      });
    },

    /**
     * Up a container
     *
     * @param {*} container
     */
    upContainer(container) {
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