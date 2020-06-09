const merge = require('deepmerge')

class VuexPersist {
  /**
   * Constructor.
   * 
   * @param {object} options
   */
  constructor(options) {
    this.options = Object.assign({
      reducer: null,
      mutations: [],
      driver: null
    }, options)
    this.driver = null
    this.driverCheckCount = 0
    this.driverCheckInterval = null
  }

  /**
   * Check if vuex is ready
   * 
   * @returns {Boolean}
   */
  checkState() {
    return this.driver !== null
  }

  /**
   * Check if driver is ready
   * 
   * @param {Function} callback 
   */
  checkDriver(callback) {
    if (this.driver === null && this.driverCheckInterval === null) {
      this.driverCheckInterval = setInterval(() => {
        this.driverCheckCount++
        if (this.driver === null) {

          //Limit driver checking to 30*250ms
          if (this.driverCheckCount >= 30) {
            clearInterval(this.driverCheckInterval)
            alert("Vuex driver check failed, fatal error!")
            return;
          }

          //Check if Go vuex is available
          if (typeof window.backend !== "undefined" && this.driver === null) {
            this.driver = window.backend.VuexState;
          } else if (this.driver !== null) {
            if (typeof callback === "function") {
              clearInterval(this.driverCheckInterval)
              callback()
            }
          }
        } else {
          clearInterval(this.driverCheckInterval)
        }
      }, 250)
    }
  }

  /**
   * Persist the state to file.
   * @param {object} state
   */
  saveState(state) {
    this.checkDriver(() => {
      let data = JSON.stringify(this.options.reducer ? this.options.reducer(state) : state)
      this.driver.Write(data)
    })
  }

  /**
   * Load the state from file.
   * @param {object} store
   */
  loadState(store) {
    try {
      this.checkDriver(() => {
        this.driver.Read()
          .then((data) => {
            store.replaceState(merge(store.state, JSON.parse(data)))
          })
      })
    } catch (err) {
      alert(err)
    }
  }

  /**
   * Subscribe to the Vuex store.
   * @returns {function}
   */
  subscribe() {
    this.checkDriver()
    return (store) => {
      this.loadState(store)

      store.subscribe((mutation, state) => {
        if (this._mutation(mutation.type)) {
          this.saveState(state)
        }
      })
    }
  }

  /**
   * Checks if a mutation is in the list of allowed
   * mutations.
   * @param {string} type
   * @returns {boolean}
   * @private
   */
  _mutation(type) {
    this.checkDriver()
    return !this.options.mutations.length ||
      this.options.mutations.includes(type)
  }
}

export default VuexPersist