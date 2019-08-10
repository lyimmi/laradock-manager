<template>

  <v-container fluid>
    <v-layout align-center justify-center wrap>
      <v-flex xs12>
        <v-card>
          <v-tabs v-model="tab">
            <v-tabs-slider></v-tabs-slider>

            <v-tab href="#tab-basic">
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
            <v-tab-item :value="'tab-basic'">
              <v-card flat>
                <v-card-text>
                  <v-layout>
                    <v-flex sm6>
                      <template>
                        <v-text-field
                            v-model="dockerComposeYmlPath"
                            label="Absolute path to laradock folder"
                            :placeholder="laradockPath"
                            :disabled="laradockPath !== '' && dockerComposeYmlPath === laradockPath"
                            ref="laradockPathInput"
                        >
                          <template v-slot:append>
                            <v-btn
                                depressed
                                tile
                                color="secondary"
                                class="ma-0"
                                @click="storeLaradockPath"
                                :disabled="dockerComposeYmlPath === ''"
                            >
                              <v-icon>done</v-icon>
                            </v-btn>
                          </template>
                        </v-text-field>
                      </template>
                    </v-flex>
                  </v-layout>
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import {mapGetters} from 'vuex'
  import {mapActions} from 'vuex'
  import dockerCompose from '../../shared/dockerCompose'

  export default {
    name: 'index',
    data() {
      return {
        dockerComposeYmlPath: "",
        tab: null,
      }
    },
    mixins: [dockerCompose],
    computed: {
      ...mapGetters('Settings', [
        'laradockPath'
      ])
    },
    methods: {
      ...mapActions('Settings', [
        'setLaradockPath'
      ]),
      storeLaradockPath() {
        this.setLaradockPath(this.dockerComposeYmlPath)
        this.applyLaradockPath(this.dockerComposeYmlPath)
      }
    },
    watch: {}
  }
</script>

<style scoped>

</style>