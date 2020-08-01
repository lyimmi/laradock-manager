<template>
  <v-card flat tile>
    <v-form v-model="valid" ref="appSettings">
      <v-container fluid>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              required
              label="Container prefix"
              v-model="containerPrefix"
              :rules="containerPrefixRules"
            ></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              readonly
              required
              label="Laradock path"
              v-model="laradockPath"
              :rules="laradockPathRules"
              @click.prevent="selectLaradockDirectory"
            ></v-text-field>
          </v-col>

          <v-col cols="12" md="4">
            <v-text-field
              v-model="terminalPath"
              :rules="terminalPathRules"
              label="Terminal path"
              required
            ></v-text-field>
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="4">
            <v-switch v-model="darkTheme" label="Dark theme" />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols="12" md="12">
            <v-btn color="error">Clear app settings</v-btn>
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </v-card>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import DockerMixin from "@/shared/dockerMixin";
export default {
  name: "app",
  mixins: [DockerMixin],
  data: () => {
    return {
      valid: true,
      laradockPathRules: [v => !!v || "Laradock path is required"],
      terminalPathRules: [v => !!v || "Terminal path is required"],
      containerPrefixRules: [v => !!v || "Container prefix is required"]
    };
  },
  mounted() {
    this.$refs.appSettings.validate();
  },
  computed: {
    ...mapGetters("Settings", ["laradockPath", "terminalPath"]),
    containerPrefix: {
      set(value) {
        this.$store.dispatch("Settings/setContainerPrefix", value);
      },
      get() {
        return this.$store.getters["Settings/containerPrefix"];
      }
    },
    darkTheme: {
      set(value) {
        this.$store.dispatch("Settings/setDarkTheme", value);
      },
      get() {
        return this.$store.getters["Settings/darkTheme"];
      }
    }
  },
  watch: {
    darkTheme(val) {
      this.$vuetify.theme.dark = val;
    }
  },
  methods: {
    ...mapActions("Settings", ["setLaradockPath", "setTerminalPath"]),
    /**
     * Select a directory return it's path
     */
    selectTerminalExecutable() {
      window.backend.App.SelectFile().then(res => {
        this.terminalPathTmp = res;
      });
    },

    /**
     * Select Laradock directory
     */
    selectLaradockDirectory() {
      window.backend.App.SelectDirectory().then(res => {
        this.setLaradockPath(res);
      });
    },

    /**
     * Store laradock path
     */
    storeLaradockPath(path) {
      this.setLaradockPath(path);
      this.applyLaradockPath(path);
      this.hasDotEnv = path !== "";
    }
  }
};
</script>