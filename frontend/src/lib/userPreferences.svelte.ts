export class UserPreferences {
  #showEvents = $state(false)
  #expandSidebar = $state(false)

  readUserPreferences() {
    const showEvents = localStorage.getItem("encoon-showEvents")
    if(showEvents) this.#showEvents = showEvents === "true"
    const expandSidebar = localStorage.getItem("encoon-expandSidebar")
    if(expandSidebar) this.#expandSidebar = expandSidebar === "true"
  }

  get showEvents() { return this.#showEvents }

  get expandSidebar() { return this.#expandSidebar }

  set showEvents(value: boolean) {
    localStorage.setItem("encoon-showEvents", value ? "true" : "false")
    this.#showEvents = value 
  }

  set expandSidebar(value: boolean) {
    localStorage.setItem("encoon-expandSidebar", value ? "true" : "false")
    this.#expandSidebar = value
  }
}
