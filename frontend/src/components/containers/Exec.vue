
<template>
  <div>
    <v-layout justify-center wrap>
      <v-dialog v-model="selectUserDialog" max-width="500px">
        <v-card>
          <v-card-title>Executing container: {{executableContiner}}</v-card-title>
          <v-card-text>
            <v-select :items="execUsers" v-model="selectedUser" label="User"></v-select>
          </v-card-text>
          <v-card-actions>
            <v-btn color="primary" text @click="selectUserDialog = false">Cancel</v-btn>
            <v-spacer></v-spacer>
            <v-btn color="primary" text @click="conainerExec">Ok</v-btn>
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
      execUsers: [
        {
          text: "laradock",
          value: "laradock"
        },
        {
          text: "root",
          value: "root"
        }
      ],
      selectedUser: "laradock"
    };
  },
  methods: {
    conainerExec() {
      this.selectUserDialog = false;
      console.log(this.selectedUser);
      this.execContainer(this.executableContiner, this.selectedUser);
    }
  }
};
</script>