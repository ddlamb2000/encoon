export class UserPreferences {
  #showEvents = $state(false)
  #showPrompt = $state(false)
  #expandSidebar = $state(false)

  readUserPreferences() {
    const showEvents = localStorage.getItem("encoon-showEvents")
    if(showEvents) this.#showEvents = showEvents === "true"
    const showPrompt = localStorage.getItem("encoon-showPrompt")
    if(showPrompt) this.#showPrompt = showPrompt === "true"
    const expandSidebar = localStorage.getItem("encoon-expandSidebar")
    if(expandSidebar) this.#expandSidebar = expandSidebar === "true"
  }

  get showEvents() { return this.#showEvents }  
  set showEvents(value: boolean) {
    localStorage.setItem("encoon-showEvents", value ? "true" : "false")
    this.#showEvents = value 
    if(value) this.showPrompt = false
  }
  
  get showPrompt() { return this.#showPrompt }
  set showPrompt(value: boolean) {
    localStorage.setItem("encoon-showPrompt", value ? "true" : "false")
    this.#showPrompt = value 
    if(value) this.showEvents = false
  }
  
  get expandSidebar() { return this.#expandSidebar }
  set expandSidebar(value: boolean) {
    localStorage.setItem("encoon-expandSidebar", value ? "true" : "false")
    this.#expandSidebar = value
  }

  toggleSidebar = () => {
    this.expandSidebar = !this.expandSidebar
  }

  toggleShowEvents = () => {
    this.showEvents = !this.showEvents
  }

  toggleShowPrompt = () => {
    this.showPrompt = !this.showPrompt
  }
}
