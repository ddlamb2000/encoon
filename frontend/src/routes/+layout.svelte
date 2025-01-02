<script lang="ts">
  import Header from './Header.svelte'
  import '$lib/app.css'
  let { children } = $props()
</script>

<Header></Header>
{@render children()}

<style>
</style>