<script>
	import { Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell } from 'flowbite-svelte';
	import { Button, Dropdown, DropdownItem, Search, Checkbox, ToolbarButton, DropdownDivider } from 'flowbite-svelte';
  import { DotsHorizontalOutline, DotsVerticalOutline } from 'flowbite-svelte-icons';
  import { ChevronDownOutline, UserRemoveSolid } from 'flowbite-svelte-icons';
  let searchTerm = ''
  const people = [
    { name: 'Robert Gouth', checked: false },
    { name: 'Jese Leos', checked: false },
    { name: 'Bonnie Green', checked: true },
  ]
  $: filteredItems = people.filter((person) => person.name.toLowerCase().indexOf(searchTerm?.toLowerCase()) !== -1);
  </script>


<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<h1>εncooη</h1>
<p>For the time being, εncooη is an experimental web-based application 
oriented toward data management.</p>
<p>The object of the reasearches εncooη is based on is to propose 
a true engine that allows data structuration, data relationship 
management, data presentation, and, mainly, to make obvious data navigation.</p>	
<p>This brand new application doesn't require database system 
learning or skills, because everything in the application 
is managed through a very simple web-based user interface, 
in real time. </p>
<p>Primary dedicated to small business structures, 
the application lets manage a large amount of multi-purposes information 
within the organization, and share this information across business stakeholders 
in an easy and natural way, using one Internet browser only.</p>
<p>εncooη is also the foundation for the development of a true business-oriented 
application development toolkit.</p>
<p>Copyright David Lambert 2024</p>

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


  <Table>
	<TableHead>
	  <TableHeadCell>Product name</TableHeadCell>
	  <TableHeadCell>Color</TableHeadCell>
	  <TableHeadCell>Category</TableHeadCell>
	  <TableHeadCell>Price</TableHeadCell>
	</TableHead>
	<TableBody tableBodyClass="divide-y">
	  <TableBodyRow>
		<TableBodyCell>Apple MacBook Pro 17"</TableBodyCell>
		<TableBodyCell>Sliver</TableBodyCell>
		<TableBodyCell>Laptop</TableBodyCell>
		<TableBodyCell>$2999</TableBodyCell>
	  </TableBodyRow>
	  <TableBodyRow>
		<TableBodyCell>Microsoft Surface Pro</TableBodyCell>
		<TableBodyCell>White</TableBodyCell>
		<TableBodyCell>Laptop PC</TableBodyCell>
		<TableBodyCell>$1999</TableBodyCell>
	  </TableBodyRow>
	  <TableBodyRow>
		<TableBodyCell>Magic Mouse 2</TableBodyCell>
		<TableBodyCell>Black</TableBodyCell>
		<TableBodyCell>Accessories</TableBodyCell>
		<TableBodyCell>$99</TableBodyCell>
	  </TableBodyRow>
	</TableBody>
  </Table>