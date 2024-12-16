<script module>

	export function newUuid() { return crypto.randomUUID() }

	export function numberToLetters(num) {
		let letters = ''
		while (num >= 0) {
			letters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'[num % 26] + letters
			num = Math.floor(num / 26) - 1
		}
		return letters
	}

</script>