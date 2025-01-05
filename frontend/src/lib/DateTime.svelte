<script lang="ts">
  let { dateTime, showDate = true } = $props()

  let localDate = $state(new Date(dateTime))

  const getTimeAgo = () => {
		const localNow = new Date
    const localDateUTC =  Date.UTC(localDate.getUTCFullYear(),
                                    localDate.getUTCMonth(),
                                    localDate.getUTCDate(),
                                    localDate.getUTCHours(),
                                    localDate.getUTCMinutes(),
                                    localDate.getUTCSeconds())
		const localNowUTC =  Date.UTC(localNow.getUTCFullYear(),
                                  localNow.getUTCMonth(),
                                  localNow.getUTCDate(),
                                  localNow.getUTCHours(),
                                  localNow.getUTCMinutes(),
                                  localNow.getUTCSeconds())
		const seconds = (localNowUTC - localDateUTC) / 1000
		const MINUTE = 60, HOUR = MINUTE * 60, DAY = HOUR * 24, WEEK = DAY * 7, MONTH = DAY * 30, YEAR = DAY * 365		
		if(seconds < MINUTE) return `${Math.round(seconds)}&nbsp;sec&nbsp;ago`
		if(seconds < HOUR) return `${Math.round(seconds / MINUTE)}&nbsp;min&nbsp;ago`
		if(seconds < DAY) return `${Math.round(seconds / HOUR)}&nbsp;hour&nbsp;ago`
		if(seconds < WEEK) return `${Math.round(seconds / DAY)}&nbsp;day&nbsp;ago`
		if(seconds < MONTH) return `${Math.round(seconds / WEEK)}&nbsp;week&nbsp;ago`
		if(seconds < YEAR) return `${Math.round(seconds / MONTH)}&nbsp;month&nbsp;ago`
		return `${Math.round(seconds / YEAR)}&nbsp;year&nbsp;ago`
	}

  let timeAgo = $state(getTimeAgo())

  $effect(() => {
		const interval = setInterval(() => {
      if(dateTime !== undefined) {
        localDate = new Date(dateTime)
        timeAgo = getTimeAgo()
      } else timeAgo = ""
    }, 1000)
		return () => { clearInterval(interval) }
	})  

</script>
 
{#if showDate}
  {localDate.toLocaleDateString()} {localDate.toLocaleTimeString()}
{/if}
<small><em>{@html timeAgo}</em></small>