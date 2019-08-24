<template>
  <v-container fluid>
    <v-layout align-center justify-center wrap>
      <v-flex xs12>
        <v-card>
          <v-tabs v-model="tab">
            <v-tabs-slider></v-tabs-slider>

            <v-tab href="#tab-basic">Basics</v-tab>

            <v-tab href="#tab-env">Env</v-tab>

            <v-tab href="#tab-3">Containers</v-tab>
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
                      </template>
                    </v-flex>
                  </v-layout>
                </v-card-text>
              </v-card>
            </v-tab-item>

            <v-tab-item :value="'tab-env'">
              <v-fab-transition>
                <v-btn
                  fixed
                  fab
                  small
                  bottom
                  right
                  color="pink"
                  style="z-index:2;"
                  @click="saveDotEnv"
                >
                  <v-icon>fa-save</v-icon>
                </v-btn>
              </v-fab-transition>
              <v-expansion-panels>
                <v-expansion-panel v-for="(item, key) in dotEnvContentGroups" :key="key">
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
            </v-tab-item>
          </v-tabs-items>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { mapGetters } from "vuex";
import { mapActions } from "vuex";
import dockerCompose from "../../shared/dockerCompose";

export default {
  name: "index",
  data() {
    return {
      dockerComposeYmlPath: "",
      tab: null,
      form: {}
    };
  },
  created() {
    this.getDotEnv(() => {
      Object.keys(this.dotEnvContents).forEach(element => {
        this.form[element] = this.dotEnvContents[element];
      });
    });
  },
  mounted() {},
  mixins: [dockerCompose],
  computed: {
    ...mapGetters("Settings", ["laradockPath"])
  },
  methods: {
    ...mapActions("Settings", ["setLaradockPath"]),
    storeLaradockPath() {
      this.setLaradockPath(this.dockerComposeYmlPath);
      this.applyLaradockPath(this.dockerComposeYmlPath);
    },
    saveDotEnv() {
      this.$nextTick(() => {
        this.writeDotEnv(this.dotEnvContents);
      });
    }
  },
  watch: {}
};
</script>

<style scoped>
</style>