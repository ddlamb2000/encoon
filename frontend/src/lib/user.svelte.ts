// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

export class User {
	#token: string = $state("")
  #userUuid: string = $state("")
  #user: string = $state("")
  #userFirstName: string = $state("")
  #userLastName: string = $state("")
  #loggedIn: boolean = $state(false)
  #tokenName = ""

  constructor(dbName: string) {
    this.#tokenName = `access_token_${dbName}`
  }

  reset() {
    this.#token = ""
    this.#loggedIn = false
    this.#userUuid = ""
    this.#user = ""
    this.#userFirstName = ""
    this.#userLastName = ""
  }

  getUserUuid = (): string => this.#userUuid
  getUser = (): string => this.#user
  getFirstName = (): string => this.#userFirstName
  getLastName = (): string => this.#userLastName
  getToken = (): string => this.#token
  getIsLoggedIn = (): boolean => this.#loggedIn
  removeToken = () => localStorage.removeItem(this.#tokenName)
  setToken = (jwt: string) => localStorage.setItem(this.#tokenName, jwt)

  checkToken = (): boolean => {
    const token = localStorage.getItem(this.#tokenName)
    if(token) {
      try {
        const arrayToken = token.split('.')
        const tokenPayload = JSON.parse(atob(arrayToken[1]))
        const now = (new Date).toISOString()
        const nowDate = Date.parse(now)
        const tokenExpirationDate = Date.parse(tokenPayload.expires)
        if(nowDate < tokenExpirationDate) {
          this.#token = token
          this.#loggedIn = true
          this.#userUuid = tokenPayload.userUuid
          this.#user = tokenPayload.user
          this.#userFirstName = tokenPayload.userFirstName
          this.#userLastName = tokenPayload.userLastName
          return true
        }
      } catch (error) {
        console.error(`Error checking token:`, error)
      }
    }
    this.reset()
    return false
  }
}