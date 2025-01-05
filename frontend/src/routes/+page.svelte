<script>
  import '$lib/app.css'
  import { Dropdown, Search, Checkbox, Button } from 'flowbite-svelte';
  import { ChevronDownOutline, UserRemoveSolid } from 'flowbite-svelte-icons';
  const dbs = ['master', 'test', 'sandbox']
  let searchTerm = ''
  const people = [
    { name: 'Robert Gouth', checked: false },
    { name: 'Jese Leos', checked: false },
    { name: 'Bonnie Green', checked: true },
  ]
  $: filteredItems = people.filter((person) => person.name.toLowerCase().indexOf(searchTerm?.toLowerCase()) !== -1);
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
      <div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
        <div class="p-6 space-y-4 md:space-y-6 sm:p-8">
          {#each dbs as dbname}
            <p><a href={"/" + dbname} data-sveltekit-reload>{dbname}</a></p>
          {/each}
   
          <Button class="sizes" size="xs">Dropdown search<ChevronDownOutline class="w-6 h-6 ms-2 text-white dark:text-white" /></Button>
          <Dropdown triggeredBy=".sizes" class="overflow-y-auto px-3 pb-3 text-xs h-44">
            <div slot="header" class="p-3">
              <Search size="md" bind:value={searchTerm}/>
            </div>
            {#each filteredItems as person (person.name)}
              <li class="rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600">
                <Checkbox bind:checked={person.checked}>{person.name}</Checkbox>
              </li>
            {/each}
            <a slot="footer" href="/" class="flex items-center p-3 -mb-1 text-sm font-medium text-red-600 bg-gray-50 hover:bg-gray-100 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-red-500 hover:underline">
              <UserRemoveSolid class="w-4 h-4 me-2 text-primary-700 dark:text-primary-700" />Delete user
            </a>
          </Dropdown>    
    
        </div>
      </div>
    </div>
  </section>
</main>