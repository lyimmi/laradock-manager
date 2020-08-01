export default class Ticker {

    /**
     * Ticker handles multiple interval dependent callback to save performance
     * 
     * @param {Number} interval Number of milliseconds between the ticker calls the ticks
     */
    constructor(interval = 1000) {
        this.tickers = [];
        setInterval(() => {
            this.tick()
        }, interval)
    }

    /**
     * Execure ticks
     */
    tick() {
        this.tickers.forEach((t) => {
            t.run()
        });
    }

    /**
     * Run a tick on the next tick, then remove it
     * 
     * @param {Function} callback 
     */
    nextTick(callback) {
        try {
            const tick = new Tick(callback, 1, this, true)
            this.tickers.push(tick)
            let id = this.tickers.length - 1
            tick.setId(id)
        } catch (error) {
            console.error("An error occured on tick set!")
            console.error(error)
        }
    }

    /**
     * Set a new tick
     * 
     * @param {Number} interval 
     * @param {Function} callback 
     * @returns {Number} returns the ticker's id
     */
    setTicker(callback, interval) {
        try {
            const tick = new Tick(callback, interval, this)
            this.tickers.push(tick)
            let id = this.tickers.length - 1
            tick.setId(id)
            return tick.id
        } catch (error) {
            console.error("An error occured on tick set!")
            console.error(error)
        }
    }

    /**
     * Clear a tick by it's id
     * @param {*} id 
     */
    clearTicker(id) {
        try {
            this.tickers.splice(id, 1)
        } catch (error) {
            console.error(`Ticker with id: ${id} cannot be removed, no such ticker.`)
            console.error(error)
        }
    }
}

/**
 * Tick
 */
class Tick {
    /**
     * Construct
     * 
     * @param {Function} callback 
     * @param {Number} interval 
     * @param {Ticker} ticker 
     * @param {Boolean} once 
     */
    constructor(callback, interval = 10, ticker, once = false) {
        this.interval = interval
        this.callback = callback
        this.tick = 0
        this.once = once
        this.ticker = ticker
        this.id = null
    }

    /**
     * 
     * @param {Number} id 
     */
    setId(id) {
        this.id = id
    }

    increment() {
        this.tick++
    }

    isActual() {
        return this.interval === this.tick
    }

    resetTick() {
        this.tick = 0
    }

    run() {
        try {
            this.increment()
            if (this.isActual()) {
                this.resetTick()
                this.callback()
                if (this.once) {
                    this.ticker.clearTicker(this.id)
                }
            }
        } catch (error) {
            console.error(`An error occurred while calling tick's callback, id: ${this.id} removed from the list.`)
            console.error(error)
            this.ticker.clearTicker(this.id)
        }
    }
}