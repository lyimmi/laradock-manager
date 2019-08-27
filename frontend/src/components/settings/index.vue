<template>
  <v-container fluid>
    <v-layout align-center justify-center wrap>
      <v-flex xs12>
        <v-card :loading="containersLoading">

          <v-tabs v-model="tab" vertical>
            <v-tab href="#tab-basic" class="ma-0">Basics</v-tab>
            <v-tab href="#tab-env">Env</v-tab>
            <v-fab-transition>
              <v-btn
                  fab
                  small
                  color="pink"
                  @click="saveDotEnv"
                  class="ml-6 mt-6"
                  v-if="tab === 'tab-env'"
              >
                <v-icon>fa-save</v-icon>
              </v-btn>
            </v-fab-transition>
            <v-tabs-items v-model="tab">
              <v-tab-item :value="'tab-basic'">
                <v-card flat>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" sm="6">
                        <v-text-field
                            v-model="dockerComposeYmlPath"
                            :value="laradockPath"
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
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-switch
                            label="Dark theme"
                            color="primary"
                            v-model="darkThemeSwitch"
                        ></v-switch>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-tab-item>

              <v-tab-item :value="'tab-env'">
                <v-row>
                  <v-col cols="12" xs="12">
                    <v-text-field
                        class="mr-10"
                        prepend-icon="search"
                        label="Search in .env file"
                        v-model="envFilter"
                        hide-details
                        clearable
                    ></v-text-field>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="12" xs="12">
                    <v-expansion-panels>
                      <v-expansion-panel v-for="(item, key) in dotEnvContentGroupsFiltered" :key="key">
                        <v-expansion-panel-header>{{key}}</v-expansion-panel-header>
                        <v-expansion-panel-content>
                          <v-row>
                            <v-col cols="12" md="3" v-for="input in item" :key="input.field">
                              <v-text-field
                                  placeholder=" "
                                  :label="input.fieldName"
                                  v-model="dotEnvContents[input.field]"
                                  type="text"
                              ></v-text-field>
                            </v-col>
                          </v-row>
                        </v-expansion-panel-content>
                      </v-expansion-panel>
                    </v-expansion-panels>
                  </v-col>
                </v-row>
              </v-tab-item>
            </v-tabs-items>
          </v-tabs>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import {mapActions, mapGetters} from 'vuex'
  import dockerCompose from '../../shared/dockerCompose'

  export default {
    name: 'index',
    mixins: [dockerCompose],
    data() {
      return {
        dockerComposeYmlPath: '',
        tab: null,
        form: {},
        darkThemeSwitch: true
      }
    },
    mounted() {
      this.getDotEnv(() => {
        Object.keys(this.dotEnvContents).forEach(element => {
          this.form[element] = this.dotEnvContents[element]
        })
      })
    },
    computed: {
      ...mapGetters('Settings', [
        'laradockPath',
        'darkTheme'
      ])
    },
    watch: {
      darkThemeSwitch(val) {
        console.log(val)
        this.$vuetify.theme.dark = val
        this.setDarkTheme(val)
      }
    },
    methods: {
      ...mapActions('Settings', [
        'setLaradockPath',
        'setDarkTheme'
      ]),
      storeLaradockPath() {
        this.setLaradockPath(this.dockerComposeYmlPath)
        this.applyLaradockPath(this.dockerComposeYmlPath)
      },
      saveDotEnv() {
        this.$root.$emit('containersLoading')
        this.$nextTick(() => {
          this.writeDotEnv(this.dotEnvContents, () => {
            setTimeout(() => {
              this.$root.$emit('containersNotLoading')
            }, 500)
          })
        })
      }
    }
  }
</script>

<style scoped>
</style>