<script>
import {onMount, createEventDispatcher} from 'svelte'
import Item from './item.svelte'

const dispatch = createEventDispatcher();

export let multiple;
export let items;

function select(item) { 
    let state = items[item.detail]?.selected
    if(!state) {
        items[item.detail].selected = !state
    }
    if(!multiple) {
        items.forEach((x, i) => {
            if(i != item.detail) {
                delete x.selected
            }
        })
    }
    dispatch("selected", items)
}

onMount(() => {
    let isSelected  = items?.filter(x => x.selected)?.length > 0
    // if no items are selected by default, select the first item
    if(!isSelected && !multiple) {
        items[0].selected = true
    }
})

</script>


{#if items}

    <div class="select-container fl-co">

        {#each items as item, i}

            <Item 
            item={item} 
            index={i}
            multiple={multiple}
            on:selected={select}/>

        {/each}

    </div>

{:else}
    <div class="warn">
    No items in Select component
    </div>
{/if}



