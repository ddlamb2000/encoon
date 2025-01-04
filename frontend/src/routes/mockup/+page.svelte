<script>
  import '$lib/app.css'
  const sample = "Test"

  const textSample = [
    'For the time being, εncooη is an experimental web-based application oriented toward data management.',
    'The object of the reasearches εncooη is based on is to propose a true engine that allows data structuration, data relationship management, data presentation, and, mainly, to make obvious data navigation.',
    'This brand new application doesnt require database system learning or skills, because everything in the application is managed through a very simple web-based user interface, in real time.',
    'Primary dedicated to small business structures, the application lets manage a large amount of multi-purposes information  within the organization, and share this information across business stakeholders  in an easy and natural way, using one Internet browser only.',
    'εncooη is also the foundation for the development of a true business-oriented application development toolkit.'
  ]

</script>

<svelte:head>
  <title>Mockup</title>
</svelte:head>

<main class="global-container grid h-full [grid-template-rows:auto_1fr]">
  <nav class="global header bg-violet-200">
    global header
  </nav>
  <section class="main-container bg-sky-200 grid [grid-template-columns:1fr_4fr] overflow-auto">
    <aside class="side-bar bg-teal-200 grid overflow-auto">
      <div class="overflow-auto h-[500px]">
        sidebar
      </div>
    </aside>
    <section class="content bg-lime-200 grid [grid-template-rows:auto_auto_1fr_auto] overflow-auto">
      <div class="bg-yellow-200">
        tabs
      </div>
      <aside class="bg-orange-200">
        toolbar
      </aside>
      <div class="bg-red-200 grid overflow-auto">
        <article class="h-[500px]">
          content
        </article>
      </div>
      <footer class="bg-slate-200">
        footer
      </footer>
    </section>
  </section>
</main>

<!-- 
<header class="sticky top-0 h-14 bg-blue-500 text-white text-center py-4">
  <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"> εncooη </span>
  Header
</header>
<section class="flex">
  <nav class="fixed top-14 w-64 h-dvh overflow-y-auto bg-gray-300 text-white text-center py-4">
    {#each textSample as sample}<p>Nav {sample}</p>{/each}
  </nav>
  <article class="fixed left-64 h-dvh overflow-auto bg-gray-200 text-center py-4">
    {#each textSample as sample}<p>Article {sample}</p>{/each}
    {#each textSample as sample}<p>Article {sample}</p>{/each}
    {#each textSample as sample}<p>Article {sample}</p>{/each}
    {#each textSample as sample}<p>Article {sample}</p>{/each}
  </article>
  <aside class="fixed flex-none left-full top-14 w-60 h-screen overflow-y-auto bg-gray-100 text-center py-4">
    {#each textSample as sample}<p>Side {sample}</p>{/each}
  </aside>
</section>
<footer class="fixed inset-x-0 bottom-0 h-14 bg-blue-500 text-white text-center py-4">
  Footer {sample}
</footer>
 -->
