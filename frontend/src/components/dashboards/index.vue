<template>
    <v-container
            fluid
    >
        <v-layout
                align-center
                justify-center
        >
            <v-flex xs3>
                <div v-if="dockerComposeStatus !== null">
                    <v-alert type="success" v-if="dockerComposeStatus">
                        docker-compose is accessable
                    </v-alert>
                    <v-alert type="error" v-else>
                        docker-compose not accessible or not installed
                    </v-alert>
                </div>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
    export default {
        name: "index",
        data() {
            return {
                dockerComposeStatus: null,
            }
        },
        mounted() {
            this.getMessage();
        },
        methods: {
            getMessage: function () {
                var self = this;
                window.backend.basic().then(result => {
                    self.dockerComposeStatus = result.startsWith("OK:");
                });
            }
        }
    }
</script>

<style scoped>

</style>