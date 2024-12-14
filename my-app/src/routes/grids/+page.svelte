<script lang="ts">
    import { seedData } from '$lib/data.js'
    import { newUuid, numberToLetters } from "$lib/utils.svelte"

	const grids = $state(seedData)
	let focus = $state({grid: null, i: -1, j: -1})

	function pushTransaction(payload) {
		console.log(payload)
	}

	function initGrid(grid) {
		grid.search = ''
		grid.columnSeq = grid.cols.length
		applyFilters(grid)
	}

	function applyFilters(grid) {
		if (grid.search === '') grid.rows.forEach((row) => row.filtered = true)
		else {
			const regex = new RegExp(grid.search, 'i')
			grid.rows.forEach((row) => row.filtered = regex.test(row.data[0]))
		}
	}

	grids.forEach((grid) => initGrid(grid))

	async function newGrid() {
		const grid = {uuid: newUuid(), title: 'Untitled', 
								  cols: [{uuid: newUuid(), title: 'A', type: 'coltypes-row-1'}],
								  rows: [{uuid: newUuid(), data: ['***']}]
								 }
		initGrid(grid)
		grids.push(grid)
		pushTransaction({action: 'newgrid', griduuid: grid.uuid})
	}

	async function addRow(grid) {
		const uuid = newUuid()
		const data = []
		grid.cols.forEach((col) => data.push(''))
		grid.rows.push({uuid: uuid, data: data, filtered: true})
		pushTransaction({action: 'addrow', griduuid: grid.uuid, rowuuid: uuid})
	}

	async function removeRow(grid, rowuuid) {
		grid.rows = grid.rows.filter((t) => t.uuid !== rowuuid)
		pushTransaction({action: 'delrow', griduuid: grid.uuid, rowuuid: rowuuid})
	}

	async function addColumn(grid) {
		const col = {uuid: newUuid(), title: numberToLetters(grid.columnSeq), type: 'coltypes-row-1'}
		grid.cols.push(col)
		grid.columnSeq += 1
		grid.rows.forEach((row) => row.data.push(''))
		pushTransaction({action: 'addcol', griduuid: grid.uuid, col: col})
	}

	async function removeColumn(grid, coluuid) {
		const colindex = grid.cols.findIndex((col) => col.uuid === coluuid)
		grid.cols.splice(colindex, 1)
		grid.rows.forEach((row) => row.data.splice(colindex, 1))
		pushTransaction({action: 'delcol', griduuid: grid.uuid, coluuid: coluuid})
	}

	async function changeCell(grid, rowuuid, coluuid, value) {
		pushTransaction({action: 'chgcell',
										 griduuid: grid.uuid,
										 rowuuid: rowuuid,
										 coluuid: coluuid,
										 value: value})
	}

	function changeFocus(grid, i, j) { focus = {grid: grid, i: i, j: j} }

	function findGrid(uuid) { return grids.find((grid) => grid.uuid === uuid) }
	
	const coltypesGrid = findGrid('coltypes')</script>


<div class="layout">
	<main>
		<ul>
			{#each grids as grid}
				{#key grid.uuid}
					<li>
						<h1>{grid.title}</h1>
						Filter: 
						<span
							bind:innerHTML={grid.search}
							oninput={() => applyFilters(grid)}
							contenteditable>				
						</span>
						<table>
							<thead>
								<tr>
									<th></th>
									{#each grid.cols as col, j}
										<th class='header'>
											<span bind:innerHTML={grid.cols[j].title} contenteditable>{col}</span>
											<select bind:value={col.type} onchange={() => console.log(col.type)}>
												{#each coltypesGrid.rows as row}
													<option value={row.uuid}>{row.data[0]}</option>
												{/each}
											</select>
											<button onclick={() => removeColumn(grid, col.uuid)}>-</button>
										</th>
									{/each}
									<th><button onclick={() => addColumn(grid)}>+</button></th>
								</tr>
							</thead>
							<tbody>
							{#each grid.rows as row, i}
								{#if row.filtered}
									{#key row.uuid}
										<tr>
											<td>
												<button onclick={() => addRow(grid)}>+</button>
												<button onclick={() => removeRow(grid, row.uuid)}>-</button>
											</td>
											{#each grid.cols as col, j}
												<td
													class={
														(focus.grid !== null && focus.grid.uuid === grid.uuid
														 && focus.i === i && focus.j === j) 
														? 'focus' : 'cell'}
													>
													<div>
														<span
															bind:innerHTML={grid.rows[i].data[j]}
															oninput={() => changeCell(grid, row.uuid, grid.cols[j].uuid, grid.rows[i].data[j])}
															onfocus={() => changeFocus(grid, i, j)}
															contenteditable>
														</span>
													</div>
												</td>
											{/each}
										</tr>
									{/key}
								{/if}
							{/each}
							</tbody>
						</table>
						{grid.rows.length} rows in total
					</li>
				{/key}
			{/each}
			<button onclick={() => newGrid()}>New Grid</button>
		</ul>	
	</main>
	<aside>
		<h2>Infos</h2>
		{#if focus.grid !== null}
			<ul>
				<li>Grid: {focus.grid.title}</li>
				<li>Columns
					<ul>
						{#each focus.grid.cols as col}
							<li>{col.title}</li>
						{/each}
					</ul>
				</li>
				<li>Content: {focus.grid.rows[focus.i].data[focus.j]}</li>
			</ul>
		{/if}
	</aside>
</div>

<style>
	@media (min-width: 640px) {
		.layout {
			display: grid;
			gap: 2em;
			grid-template-columns: 1fr 16em;
		}
	}

	table, th, td {
		border-collapse: collapse;
		padding: 2px;
	}
	
	li {
	  list-style: none;
		margin-bottom: 10px;
	}
	
	div {
		position: relative;
		display: inline-block;
	}

	.header {
		border: 1px solid;
	}
	
	.cell {
		border: 0.5px dotted gray;
	}

	.focus {
		border: 0.5px solid;
	}
</style>