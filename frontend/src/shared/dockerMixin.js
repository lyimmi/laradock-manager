import { mapActions, mapGetters } from "vuex";
// import _ from "lodash";
import ErrorHandler from "@/shared/errorHandlerMixin"
export default {
    mixins: [ErrorHandler],
    data: () => {
        return {
            containersLoading: false,
        }
    },
    computed: {
        ...mapGetters("Containers", ["favoritContainers", "availableContainers"]),
    },
    methods: {
        ...mapActions("Containers", ["setAvailableContainers", "updateAvailableContainers"]),

        /**
         * load containers with statuses
         */
        loadConstainers() {
            this.containersLoading = true
            window.backend.Compose.GetContainersWithStatuses()
                .then(res => {
                    try {
                        let response = this.handleResponse(res)
                        this.setAvailableContainers(JSON.parse(response))
                    } catch (error) {
                        this.setError(error)
                    }
                })
                .then(() => {
                    this.containersLoading = false
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })
        },

        /**
         * Stop containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        stopContainer(container) {
            this.containersLoading = true
            window.backend.Compose.Toggle("stop", container)
                .then(() => {
                    this.loadConstainers()
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })
        },

        /**
         * Start containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        startContainer(container) {
            this.containersLoading = true
            window.backend.Compose.Toggle("start", container)
                .then(() => {
                    this.loadConstainers()
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })
        },

        /**
         * Up containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        upContainer(container) {
            this.containersLoading = true
            window.backend.Compose.Up(container)
                .then(() => {
                    this.loadConstainers()
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })
        },

        /**
         * Build containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         * @param {Boolean} force true to use --no-cache false to use cache
         */
        buildContainer(container, force = false) {
            this.$root.$refs.confirm
                .open('Build container', 'Are you sure to build the container?', { color: 'default' })
                .then((confirm) => {
                    if (confirm) {
                        this.containersLoading = true
                        window.backend.Compose.Build(container, force)
                            .then(() => {
                                this.loadConstainers()
                            })
                            .catch(error => {
                                this.setError(error)
                                this.containersLoading = false
                            })
                    }
                })
        },

        /**
         * Check if any container is up
         */
        hasRunning() {
            this.containersLoading = true
            return new Promise(resolve => {
                window.backend.Compose.HasRunning()
                .then(res => {
                    resolve(res)
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                    resolve(false)
                })
            })
        },

        /**
         * Start docker status events (goroutine)
         */
        getStats() {
            window.backend.Compose.Stats()
        },

        /**
         * Stop docker status events (goroutine)
         */
        stopStats() {
            window.backend.Compose.StatsStop()
        },

        /**
         * Open a terminall to show container's logs
         * 
         * @param {String} container 
         */
        logContainer(container) {
            window.backend.Compose.Logs(container)
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })
        },

        /**
         * Execute a container
         */
        execContainer() {
            if (this.executableContainer === "" || this.executableUser === "") {
                return
            }
            window.backend.Compose.Exec(this.executableContainer, this.executableUser)
                .then(() => {
                    this.executableContainer = ""
                    this.executableUser = ""
                })
                .catch(error => {
                    this.setError(error)
                    this.containersLoading = false
                })

        },

        /**
         * Handle wails response
         * 
         * @param {String} res
         */
        handleResponse(res) {
            try {
                res = JSON.parse(res)
                if (res.success) {
                    return res.message
                } else {
                    console.log(res)
                    this.setError(res.message)
                    return null
                }
            } catch (error) {
                return {}
            }
        }
    }
}