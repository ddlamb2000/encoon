<script lang="ts">
  let { context, set, row, column, value = $bindable() } = $props()
</script>

<td contenteditable
    class="{context.isFocused(set, column, row) ? 'focus' : 'cell'}"
    onfocus={() => context.changeFocus(set, row, column)}
    oninput={() => context.changeCell(set, row)}
    bind:innerHTML={value}>
  {row[column.name]}
</td>

<style>
  .cell { border: 0.5px dotted gray; }  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>