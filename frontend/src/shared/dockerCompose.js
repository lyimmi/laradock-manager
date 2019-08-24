import { mapActions, mapGetters } from "vuex";

export default {
  data() {
    return {
      dockerComposeStatus: null,
      dockerError: "",
      dotEnv: null,
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
    applyLaradockPath(laradockPath) {
      window.backend.Compose.SetLaradockPath(laradockPath);
    },
    getAvailableContainers(callback) {
      let self = this;
      this.waitForSettings(() => {
        self.getContainers(() => {
          window.backend.Compose.GetAvailables().then(result => {
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
    checkDotEnv() {
      this.waitForSettings(() => {
        let self = this;
        window.backend.Compose.CheckDotEnv().then(result => {
          self.dotEnv = result;
          self.setAppStatus({ dotEnv: result });
        });
      });
    },
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
    execContainer(container, user, callback) {
      window.backend.Compose.Exec(container, user).then(() => {
        if (typeof callback === "function") {
          callback();
        }
      });
    },
    stopExecContiner(callback) {
      window.backend.Compose.StopExec().then(res => {
        if (res === "disconnected" && typeof callback === "function") {
          callback();
        }
      });
    }
  }
};
