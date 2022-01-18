<script>
import {onMount, createEventDispatcher} from 'svelte'
import Item from './item.svelte'

const dispatch = createEventDispatcher();

export let multiple;
export let items;

function select(item) { 
    dispatch("selected", item.detail)
}

let _items = items;
onMount(() => {
    let isSelected  = _items?.filter(x => x.selected)?.length > 0
    // if no items are selected by default, select the first item
    if(!isSelected) {
        _items[0].selected = true
    }
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
    console.log(isSelected)
})

</script>


{#if _items}

    <div class="select-container fl-co">

        {#each _items as item}

            <Item 
            item={item} 
            multiple={multiple}
            on:selected={select}/>

        {/each}

    </div>

{:else}
    <div class="warn">
    No items in Select component
    </div>
{/if}



