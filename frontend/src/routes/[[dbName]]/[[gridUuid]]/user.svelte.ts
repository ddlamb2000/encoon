// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

export class User {
	#token: string = $state("")
  #userUuid: string = $state("")
  #user: string = $state("")
  #userFirstName: string = $state("")
  #userLastName: string = $state("")
  #loggedIn: boolean = $state(false)

  reset() {
    this.#token = ""
    this.#loggedIn = false
    this.#userUuid = ""
    this.#user = ""
    this.#userFirstName = ""
    this.#userLastName = ""
  }

  getUserUuid(): string { return this.#userUuid }
  getUser(): string { return this.#user }
  getFirstName(): string { return this.#userFirstName }
  getLastName(): string { return this.#userLastName }
  getToken(): string { return this.#token }
  getIsLoggedIn(): boolean { return this.#loggedIn }

  checkToken(token: string | null): boolean {
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