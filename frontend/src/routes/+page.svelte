<script>
  import '$lib/app.css'
  const dbs = ['master', 'test', 'sandbox']
</script>

<svelte:head><title>εncooη</title></svelte:head>
<main class="global-container grid h-full [grid-template-rows:auto_1fr]">
  <nav class="p-2 global header bg-gray-900 text-gray-100">
    <div class="relative flex items-center">
      <span class="ms-2 text-xl font-extrabold">
        <a href="/">εncooη</a>
      </span>
      <span class="lg:flex ml-auto">
      </span>
    </div>
  </nav>
  <section>
    <div class="flex flex-col items-center justify-center px-6 py-8 mx-auto mt-20 md:h-fit lg:py-0">
      <a href="#top" class="flex items-center mb-6 text-2xl font-extrabold text-gray-900 dark:text-white">
        εncooη
      </a>
      <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
        <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
          <h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white">
              Select a database
          </h1>
          {#each dbs as dbname}
            <p><a href={"/" + dbname} data-sveltekit-reload>{dbname}</a></p>
          {/each}
        </div>
      </div>
    </div>
  </section>
</main>