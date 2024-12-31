<script lang="ts">
  let { set, row, column, value = $bindable(), isFocused, changeFocus, changeCell } = $props()
</script>
<td contenteditable
    class="{isFocused(set, column, row) ? 'focus' : 'cell'}"
    onfocus={() => changeFocus(set, row, column)}
    oninput={() => changeCell(set, row)}
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