import {mapActions, mapGetters} from "vuex"

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
      availableContainers: []
    }
  },
  computed: {
    ...mapGetters("Status", [
      "appStatus"
    ]),
    ...mapGetters("Settings", [
      "laradockPath"
    ])
  },
  methods: {
    ...mapActions("Status", [
      "setAppStatus"
    ]),
    waitForSettings(callback, i) {
      i = typeof i === "undefined" ? 0 : i
      if (i > 3) {
        this.$root.$emit("showError", "Laradock's path is not set. Please set it on the settings page.")
        return false
      } else if (typeof this.laradockPath === "undefined" || this.laradockPath === "") {
        setTimeout(() => {
          this.waitForSettings(callback, i + 1)
        }, 350)
      } else if (typeof callback === "function") {
        callback()
      }
    },
    applyLaradockPath(laradockPath) {
      window.backend.DockerCompose.SetLaradockPath(laradockPath)
    },
    getAvailableContainers(callback) {
      let self = this
      this.waitForSettings(() => {
        self.getContainers()
        window.backend.DockerCompose.GetAvailableContainers().then(result => {
          let data = JSON.parse(result)
          let containers = []
          data.forEach((c) => {
            let co = self.containers.find((co) => {
              return co.code === c
            })
            containers.push({
              name: c,
              state: typeof co === "undefined" ? "DOWN" : co.state
            })
          })
          containers.sort((a, b) => {
            return ("" + b.state).localeCompare(a.state)
          })
          self.availableContainers = containers

          if (typeof callback === "function") {
            callback()
          }
        })
      })
    },
    getContainers() {
      this.waitForSettings(() => {
        let self = this
        self.containersLoading = true
        window.backend.DockerCompose.GetContainers().then(result => {
          if (result.startsWith("Error:")) {
            self.containersLoading = false
            this.$root.$emit("showError", result)
            return
          }
          result = JSON.parse(result)
          let containers = []
          result.forEach((line, k) => {
            if (k > 1 && line[0] !== "") {
              containers.push({
                name: line[0],
                code: line[0].replace("laradock_", "").replace("_1", ""),
                command: line[1],
                state: line[2],
                ports: line[3]
              })
            }
          })
          self.containers = containers
          self.containersLoading = false
        })
      })
    },
    checkDotEnv() {
      this.waitForSettings(() => {
        let self = this
        window.backend.DockerCompose.CheckDotEnv().then(result => {
          self.dotEnv = result
          self.setAppStatus({dotEnv: result})
        })
      })
    },
    copyFromExample() {
      this.waitForSettings(() => {
        let self = this
        window.backend.DockerCompose.CopyEnv().then(result => {
          if (result) {
            self.dotEnv = true
          } else {
            self.$root.$emit("showError", "Copy failed")
          }
        })
      })
    },
    toggleContainer(state, container) {
      this.waitForSettings(() => {
        this.containersLoading = true
        window.backend.DockerCompose.ToggleContainer(state, container)
          .then(() => {
            this.getContainers()
          })
      })
    }
  }
}