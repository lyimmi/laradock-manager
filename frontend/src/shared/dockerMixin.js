import { mapActions, mapGetters } from "vuex";
// import _ from "lodash";
import ErrorHandler from "@/shared/errorHandlerMixin"
export default {
    mixins: [ErrorHandler],
    data: () => {
        return { containersLoading: false }
    },
    computed: {
        ...mapGetters("Containers", ["favoritContainers", "availableContainers"]),
    },
    methods: {
        ...mapActions("Containers", ["setAvailableContainers", "updateAvailableContainers"]),
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
        getStats() {
            window.backend.Compose.Stats()
        },
        stopStats() {
            window.backend.Compose.StatsStop()
        },
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