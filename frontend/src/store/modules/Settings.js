import Vue from "vue";

const namespaced = true;

const state = {
  laradockPath: "",
  terminalPath: "/usr/bin/gnome-terminal",
  containerPrefix: "laradock",
  darkTheme: true
};

const mutations = {
  SET_CONTAINER_PREFIX(state, payload) {
    Vue.set(state, "containerPrefix", payload);
  },
  SET_LARADOCK_PATH(state, payload) {
    Vue.set(state, "laradockPath", payload);
  },
  SET_TERMINAL_PATH(state, payload) {
    Vue.set(state, "terminalPath", payload);
  },
  SET_DARK_THEME(state, payload) {
    Vue.set(state, "darkTheme", payload);
  }
};

const getters = {
  containerPrefix(state) {
    return state.containerPrefix;
  },
  laradockPath(state) {
    return state.laradockPath;
  },
  terminalPath(state) {
    return state.terminalPath;
  },
  darkTheme(state) {
    return state.darkTheme;
  }
};

const actions = {
  setContainerPrefix(context, payload) {
    context.commit("SET_CONTAINER_PREFIX", payload);
  },
  setLaradockPath(context, payload) {
    context.commit("SET_LARADOCK_PATH", payload);
  },
  setTerminalPath(context, payload) {
    context.commit("SET_TERMINAL_PATH", payload);
  },
  setDarkTheme(context, payload) {
    context.commit("SET_DARK_THEME", payload);
  }
};

export default {
  namespaced,
  state,
  mutations,
  getters,
  actions
};
