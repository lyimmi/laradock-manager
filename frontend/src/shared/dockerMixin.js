import { mapActions, mapGetters } from "vuex";
// import _ from "lodash";

export default {
    data: () => { return {}},
    computed: {
        ...mapGetters("Containers", ["favoritContainers", "availableContainers"]),
        ...mapGetters("Status", ["appStatus"]),
        ...mapGetters("Settings", ["laradockPath", "containerPrefix"])
    },
    methods: {
        ...mapActions("Status", ["setAppStatus"]),
        ...mapActions("Containers", ["setAvailableContainers", "updateAvailableContainers"]),
    }
}