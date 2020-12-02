import { mapActions, mapGetters } from "vuex";
import ErrorHandler from "@/shared/errorHandlerMixin"
export default {
    mixins: [ErrorHandler],
    data: () => {
        return {
            containersLoading: false,
            hasDotEnv: false,
            startup: true
        }
    },
    computed: {
        ...mapGetters("Containers", ["favoritContainers", "availableContainers"]),
        ...mapGetters("Settings", ["laradockPath", "terminalPath"]),
    },
    methods: {
        ...mapActions("Containers", ["setAvailableContainers", "updateAvailableContainers"]),
        checkDotEnv() {
            return new Promise(res => {

                if (typeof window.startup === "undefined") {
                    window.startup = true;
                }

                let timeout = 10;

                if (window.startup) {
                    timeout = 1000;
                    window.startup = false;
                }
                setTimeout(() => {
                    if (this.laradockPath === "") {
                        this.setError("Laradock path is not set, please set it first!")
                        this.hasDotEnv = false
                        res(false);
                    }
                    this.hasDotEnv = true
                    res(true);
                }, timeout);
            })
        },
        /**
         * load containers with statuses
         */
        loadConstainers() {
            this.checkDotEnv().then(res => {
                if (!res) return
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
            });
        },

        /**
         * Stop containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        stopContainer(container) {
            this.checkDotEnv().then(res => {
                if (!res) return
                this.containersLoading = true
                window.backend.Compose.Toggle("stop", container)
                    .then(() => {
                        this.loadConstainers()
                    })
                    .catch(error => {
                        this.setError(error)
                        this.containersLoading = false
                    })
            });
        },

        /**
         * Start containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        startContainer(container) {
            this.checkDotEnv().then(res => {
                if (!res) return
                this.containersLoading = true
                window.backend.Compose.Toggle("start", container)
                    .then(() => {
                        this.loadConstainers()
                    })
                    .catch(error => {
                        this.setError(error)
                        this.containersLoading = false
                    })
            });
        },

        /**
         * Up containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         */
        upContainer(container) {
            this.checkDotEnv().then(res => {
                if (!res) return
                this.containersLoading = true
                window.backend.Compose.Up(container)
                    .then(() => {
                        this.loadConstainers()
                    })
                    .catch(error => {
                        this.setError(error)
                        this.containersLoading = false
                    })
            });
        },

        /**
         * Build containers
         * 
         * @param {String} container list of containers separated with {|} pipes
         * @param {Boolean} force true to use --no-cache false to use cache
         */
        buildContainer(container, force = false) {
            this.checkDotEnv().then(res => {
                if (!res) return
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
            });
        },

        /**
         * Check if any container is up
         */
        hasRunning() {
            return new Promise(resolve => {
                this.checkDotEnv().then(res => {
                    if (!res) resolve(false)
                    this.containersLoading = true
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
            });
        },

        /**
         * Start docker status events (goroutine)
         */
        getStats() {
            this.checkDotEnv().then(res => {
                if (!res) return
                window.backend.Compose.Stats()
            });
        },

        /**
         * Stop docker status events (goroutine)
         */
        stopStats() {
            this.checkDotEnv().then(res => {
                if (!res) return
                window.backend.Compose.StatsStop()
            });
        },

        /**
         * Open a terminall to show container's logs
         * 
         * @param {String} container 
         */
        logContainer(container) {
            this.checkDotEnv().then(res => {
                if (!res) return
                window.backend.Compose.Logs(container)
                    .catch(error => {
                        this.setError(error)
                        this.containersLoading = false
                    })
            });
        },

        /**
         * Execute a container
         */
        execContainer() {
            this.checkDotEnv().then(res => {
                if (!res) return
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
            });
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
                    this.setError(res.message)
                    return null
                }
            } catch (error) {
                return {}
            }
        }
    }
}