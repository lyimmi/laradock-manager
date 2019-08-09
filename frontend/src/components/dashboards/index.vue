<template>
    <v-container
            fluid
    >
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
                snackbarText: ""
            }
        },
        mounted() {
            this.getMessage();
            this.getContainers();
            this.checkDotEnv();
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
                window.backend.containers().then(result => {
                    return result
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
            }
        }
    }
</script>

<style scoped>

</style>