<template>

    <v-container fluid>
        <v-layout align-center justify-center wrap>
            <v-flex xs12>
                <v-card>
                    <v-tabs v-model="tab">
                        <v-tabs-slider></v-tabs-slider>

                        <v-tab href="#tab-1">
                            Basics
                        </v-tab>

                        <v-tab href="#tab-2">
                            Env
                        </v-tab>

                        <v-tab href="#tab-3">
                            Containers
                        </v-tab>
                    </v-tabs>

                    <v-tabs-items v-model="tab">
                        <v-tab-item
                                v-for="i in 3"
                                :key="i"
                                :value="'tab-' + i"
                        >
                            <v-card flat>
                                <v-card-text>{{ text }}</v-card-text>
                            </v-card>
                        </v-tab-item>
                    </v-tabs-items>
                </v-card>
            </v-flex>
        </v-layout>
    </v-container>
</template>

<script>
    export default {
        name: "index",
        data () {
            return {
                tab: null,
                text: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
            }
        },
    }
</script>

<style scoped>

</style>