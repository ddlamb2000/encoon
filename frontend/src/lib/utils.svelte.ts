export const newUuid = () => crypto.randomUUID()

export const numberToLetters = (num: number) => {
  let letters = ''
  while(num >= 0) {
    letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'[num % 26] + letters
    num = Math.floor(num / 26) - 1
  }
  return letters
}

export const debounce = <T extends unknown[], U>(callback: (...args: T) => PromiseLike<U> | U, wait: number) => {
  let timer: NodeJS.Timeout
  return (...args: T): Promise<U> => {
    clearTimeout(timer)
    return new Promise((resolve) => {
      timer = setTimeout(() => resolve(callback(...args)), wait)
    })
  }
}