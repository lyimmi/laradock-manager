const namespaced = true;

const state = {
  laradockPath: "",
  terminalPath: "/usr/bin/gnome-terminal",
  darkTheme: true
};

const mutations = {
  SET_LARADOCK_PATH(state, payload) {
    state.laradockPath = payload;
  },
  SET_TERMINAL_PATH(state, payload) {
    state.terminalPath = payload;
  },
  SET_DARK_THEME(state, payload) {
    state.darkTheme = payload;
  }
};

const getters = {
  laradockPath(state) {
    return state.laradockPath;
  },
  terminalPath(state){
      return state.terminalPath;
  },
  darkTheme(state) {
    return state.darkTheme;
  }
};

const actions = {
  setLaradockPath(context, payload) {
    context.commit("SET_LARADOCK_PATH", payload);
  },
  setTerminalPath(context, payload){
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
