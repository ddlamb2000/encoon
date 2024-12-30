<script lang="ts">
  let { set, row, column, innerHTML = $bindable(), isFocused, changeFocus, changeCell } = $props()
</script>
<td class={isFocused(set, column, row) ? 'focus' : 'cell'}  
    bind:innerHTML={innerHTML}
    onfocus={() => changeFocus(set, row, column)}
    oninput={() => changeCell(set, row)}
    contenteditable>
    {row[column.name]}
</td>
<style>
  .cell { border: 0.5px dotted gray; }  
  .focus {
    border: 0.5px solid; 
    background-color: lightyellow;
  }
</style>