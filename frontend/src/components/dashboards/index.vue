<template>
    <v-container fluid>
        <v-layout
                align-center
                justify-center
        >
            <v-flex sm6 md3>
                <div v-if="dockerComposeStatus !== null" class="px-2">
                    <v-alert :type="dockerComposeStatus ? 'success': 'error'">
                        docker-compose
                    </v-alert>
                </div>
            </v-flex>
            <v-flex sm6 md2>
                <div v-if="dotEnv !== null" class="px-23">
                    <v-alert :type="dotEnv ? 'success': 'error'">
                        .env
                        <span v-if="!dotEnv">
                            <v-btn @click="copyFromExample">Create from example</v-btn>
                        </span>
                    </v-alert>
                </div>
            </v-flex>
        </v-layout>

        <v-layout>
            <v-flex xs6>
                <v-card
                        class="mx-auto"
                        :loading="containersLoading"
                >
                    <v-card-title>
                        Containers
                    </v-card-title>
                    <v-card-text>
                        <v-simple-table fixed-header v-if="containers.length > 0">
                            <thead>
                            <tr>
                                <th class="text-left">Name</th>
                                <th class="text-center">State</th>
                                <th>Action</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr v-for="item in containers" :key="item.name">
                                <td>
                                    <b>{{ item.code }}</b><br>
                                    <small>{{ item.name }}</small>
                                </td>
                                <td class="text-center">
                                    <v-chip :color="item.state === 'Up' ? 'success' : 'warning'">
                                        {{item.state}}
                                    </v-chip>
                                </td>
                                <td>
                                    <v-btn
                                            small
                                            color="warning"
                                            v-if="item.state === 'Up'"
                                            @click="toggleContainer('stop', item.code)"
                                    >
                                        <v-icon>stop</v-icon>
                                    </v-btn>
                                    <v-btn
                                            small
                                            color="success"
                                            v-else
                                            @click="toggleContainer('start', item.code)"
                                    >
                                        <v-icon>play_arrow</v-icon>
                                    </v-btn>
                                </td>
                            </tr>
                            </tbody>
                        </v-simple-table>
                        <span v-else>
                            No containers found
                        </span>

                    </v-card-text>
                </v-card>
            </v-flex>

        </v-layout>
        <v-snackbar
                v-model="snackbar"
                color="error"
                :top="true"
        >
            {{ snackbarText }}
            <v-btn
                    text
                    @click="snackbar = false"
            >
                Close
            </v-btn>
        </v-snackbar>
    </v-container>
</template>

<script>
    export default {
        name: "index",
        data() {
            return {
                dockerComposeStatus: null,
                dockerError: "",
                dotEnv: null,
                snackbar: false,
                snackbarText: "",
                containers: [],
                containersLoading: true
            }
        },
        mounted() {
            this.getMessage();
            this.getContainers();
            this.checkDotEnv();
            this.$root.$on("refreshData", () => {
                this.getContainers();
                this.checkDotEnv();
            })
        },
        methods: {
            getMessage: function () {
                var self = this;
                window.backend.basic().then(result => {
                    self.dockerComposeStatus = result.startsWith("OK:");
                    if (!self.dockerComposeStatus) {
                        self.dockerError = result.includes('executable file not found') ?
                            "docker-compose is not installed. help: https://docs.docker.com/compose/install/" :
                            result;
                    }
                });
            },
            getContainers: function () {
                var self = this;
                self.containersLoading = true
                window.backend.containers().then(result => {
                    let lines = result.split("\n")
                    let containers = [];
                    // console.log(lines)
                    lines.forEach((line, k) => {
                        if (k > 1 && line !== "") {
                            let seg = line.split(/\s+/g)
                            containers.push({
                                name: seg[0],
                                code: seg[0].replace("laradock_", "").replace("_1", ""),
                                command: seg[1] + " " + seg[2] + " " + seg[3],
                                state: seg[4],
                                ports: seg[5]
                            })
                        }
                    })
                    self.containers = containers
                    self.containersLoading = false
                });
            },
            checkDotEnv: function () {
                var self = this;
                window.backend.checkDotEnv().then(result => {
                    console.log(result)
                    self.dotEnv = result
                });
            },
            copyFromExample: function () {
                var self = this;
                window.backend.copyEnv().then(result => {
                    console.log(result)
                    if (result) {
                        self.dotEnv = true
                    } else {
                        self.snackbarText = "Copy failed"
                        self.snackbar = true
                    }
                });
            },
            toggleContainer(state, container) {
                this.containersLoading = true
                window.backend.toggleContainers(state, container)
                    .then(result => {
                        console.log(result)
                        this.getContainers()
                    });
            }
        }
    }
</script>

<style scoped>

</style>