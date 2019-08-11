
<template>
  <div>
    <v-dialog v-model="loaderDialog" hide-overlay persistent width="300">
      <v-card>
        <v-card-text>
          Connecting to container, please stand by
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>

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
          <v-card-text style="min-height:calc(100% - 48px)">
            <v-container grid-list-sm style="height:100%">
              <v-layout style="height:90%">
                <v-flex xs12>
                  <v-sheet
                    style="height:100%;"
                    class="pa-5"
                    color="grey darken-4"
                  >{{containerContent}}</v-sheet>
                </v-flex>
              </v-layout>
              <v-layout>
                <v-flex xs12>
                  <v-text-field autofocus label="Command" outlined background-color="grey darken-4"></v-text-field>
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
            <v-select :items="select" v-model="selectedUser" label="User" item-value="text"></v-select>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" text @click="selectUserDialog = false">Cancel</v-btn>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              text
              @click="loaderDialog= true; conainerExecDialog = true;  selectUserDialog = false"
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
    wails.events.on("containerExecOutputChange", execData => {
      if (execData) {
        let data = (this.containerContent += execData).t;
        this.containerContent = data;
      }
    });
  },
  data() {
    return {
      executableContiner: "",
      conainerExecDialog: false,
      containerContent: "",
      loaderDialog: false,
      selectUserDialog: false,
      selectedUser: { text: "laradock" },
      select: [{ text: "root" }, { text: "laradock" }]
    };
  },
  methods: {},
  watch: {
    conainerExecDialog(current) {
      if (current === false) {
        console.log("container stopped");
        this.stopExecContiner(() => {
          this.executableContiner = "";
          this.containerContent = "";
          this.selectedUser = { text: "laradock" };
        });
      } else if (current === true && this.executableContiner !== "") {
        console.log(this.executableContiner, this.selectedUser.text);
        console.log("container started");
        this.execContainer(
          this.executableContiner,
          this.selectedUser.text,
          () => {
            this.loaderDialog = false;
          }
        );
      }
    }
  }
};
</script>