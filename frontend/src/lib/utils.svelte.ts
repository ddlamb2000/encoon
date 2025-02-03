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
  let timer: number
  return (...args: T): Promise<U> => {
    clearTimeout(timer)
    return new Promise((resolve) => {
      timer = setTimeout(() => resolve(callback(...args)), wait)
    })
  }
}

export const convertStreamTextToHtml = (input: string) => {
  const expBold = /\*\*([^\*]*)\*\*/g
  const replaceBold = (match: string, p1: string) => `<span class="font-bold">${p1}</span>`
  const expCode = /`([^`]*)`/g
  const replaceCode = (match: string, p1: string) => `<span class="font-mono text-xs">${p1}</span>`
  const expReference = /\{URI_REFERENCE:(\S+)\/(\S+)\/(\S+)\}/g
  const replaceReference = (match: string, p1: string, p2: string, p3: string) => {
    return `<a data-sveltekit-reload href="/${p1}/${p2}/${p3}" `
            + `class="cursor-pointer underline text-blue-800 hover:bg-gray-100 dark:hover:bg-gray-600">`
            + `Link</a>`
  }
    
  return input
          .replaceAll('\n', "<br/>")
          .replaceAll(expBold, replaceBold)
          .replaceAll(expReference, replaceReference)
          // .replaceAll(expCode, replaceCode)
}

export const convertMsToText = (input: number) => {
  if(input > 1000) return `${input/1000} s`
  else return `${input} ms`
}