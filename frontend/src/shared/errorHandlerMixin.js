export default {
    data: () => {
        return {
            errors: [],
            isMaster: false,
        }
    },
    methods: {
        /**
         * Set a component as a master error handler
         */
        setUpMasterErrorHandler() {
            this.isMaster = true
            this.$root.$on("errorSet", (error) => {
                this.addError(error.text, error.timeout);
            })
            this.$ticker.setTicker(() => {
                this.errors.forEach((err, i) => {
                    if (err.state === false) {
                        this.errors.splice(i, 1)
                    }
                });
            }, 10)
        },

        /**
         * Add an error to the main error list
         * 
         * @param {String} text 
         * @param {Number} timout 
         */
        addError(text, timeout) {
            if (this.isMaster) {
                this.errors.push({
                    text: text,
                    timeout: timeout,
                    state: true
                })
            } else {
                console.error("Only master component can add an error, user setError!")
            }
        },

        /**
         * Set an error message
         * 
         * @param {String} text 
         * @param {Number} timout 
         */
        setError(text, timeout = 10000) {
            this.$root.$emit("errorSet", {
                text: text,
                timeout: timeout,
                state: true
            })
        },

        clearError(id) {
            if (typeof this.errors[id] !== "undefined") {
                this.errors.splice(id, 1)
            }
        },

        /**
         * Get all errors
         * 
         * @returns {Array} return an array of strings
         */
        getErrors() {
            return this.errors
        },

        /**
         * Check if there are any errors
         * 
         * @returns {Boolean}
         */
        hasError() {
            return this.errors.length > 0
        }
    }
}