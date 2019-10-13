import Vue from "vue";
const namespaced = true;

const state = {
  favorites: [],
  availableContainers: []
};

const mutations = {
  SET_AVAILABLE_CONTAINERS(state, payload) {
    Vue.set(state, "availableContainers", payload);
  },
  UPDATE_AVAILABLE_CONTAINER(state, payload) {
    const index = state.availableContainers.findIndex(x => x === payload.item.name);
    payload.item.favorite = payload.isFavorite;
    state.availableContainers[index] = payload.item;
  },
  ADD_FAVORITE(state, payload) {
    state.favorites.push(payload);
  },
  REMOVE_FAVORITE(state, payload) {
    const index = state.favorites.findIndex(x => x === payload);
    state.favorites.splice(index, 1);
  }
};

const getters = {
  favoritContainers(state) {
    return state.favorites;
  },
  availableContainers(state) {
    return state.availableContainers;
  }
};

const actions = {
  setAvailableContainers(context, payload) {
    context.commit("SET_AVAILABLE_CONTAINERS", payload);
  },
  updateAvailableContainer(context, payload) {
    context.commit("UPDATE_AVAILABLE_CONTAINER", payload);
  },
  addFavoriteContiner(context, payload) {
    context.commit("ADD_FAVORITE", payload);
  },
  removeFavoriteContiner(context, payload) {
    context.commit("REMOVE_FAVORITE", payload);
  }
};

export default {
  namespaced,
  state,
  mutations,
  getters,
  actions
};
