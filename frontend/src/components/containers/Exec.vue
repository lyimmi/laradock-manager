
<template>
  <div>
    <v-layout justify-center wrap>
      <v-dialog
        v-model="conainerExecDialog"
        fullscreen
        hide-overlay
        transition="dialog-bottom-transition"
        scrollable
      >
        <v-card tile>
          <v-toolbar flat dense color="primary">
            <v-btn icon @click="conainerExecDialog = false">
              <v-icon>close</v-icon>
            </v-btn>
            <v-toolbar-title>Executing container: {{executableContiner}}</v-toolbar-title>
          </v-toolbar>
          <v-card-text style="height:100%">
            <v-container grid-list-sm style="height:100%">
              <v-layout style="height:100%">
                <v-flex xs12>
                  <iframe
                    v-if="conainerExecDialog"
                    src="http://127.0.0.1:5000/"
                    border="0"
                    height="100%"
                    width="100%"
                    style="border: 1px solid #2b2b2b"
                    id="containerExecConsole"
                  ></iframe>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>

          <div style="flex: 1 1 auto;"></div>
        </v-card>
      </v-dialog>

      <v-dialog v-model="selectUserDialog" max-width="500px">
        <v-card>
          <v-card-title>Executing container: {{executableContiner}}</v-card-title>
          <v-card-text>
            <v-text-field
                v-model="selectedUser"
                label="User" @keyup.enter="conainerExecDialog = true; selectUserDialog = false"
                placeholder="root"
                autofocus
            ></v-text-field>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" text @click="selectUserDialog = false">Cancel</v-btn>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              text
              @click="conainerExecDialog = true; selectUserDialog = false"
            >Ok</v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-layout>
  </div>
</template>
<script>
import dockerCompose from "../../shared/dockerCompose";
export default {
  name: "exec",
  mixins: [dockerCompose],
  mounted() {
    this.$root.$on("execContiner", container => {
      this.selectUserDialog = true;
      this.executableContiner = container;
    });
  },
  data() {
    return {
      executableContiner: "",
      conainerExecDialog: false,
      containerContent: "",
      selectUserDialog: false,
      selectedUser: ""
    };
  },
  methods: {},
  watch: {
    conainerExecDialog(current) {
      if (current === false) {
        this.stopExecContiner(() => {
          this.executableContiner = "";
          this.containerContent = "";
          this.selectedUser = "";
          document.getElementById("inspire").setAttribute("style", "");
        });
      } else if (current === true && this.executableContiner !== "") {
        document
          .getElementById("inspire")
          .setAttribute("style", "height:100px; overflow: hidden;");
        this.execContainer(
          this.executableContiner,
          this.selectedUser,
          () => {}
        );
      }
    }
  }
};
</script>