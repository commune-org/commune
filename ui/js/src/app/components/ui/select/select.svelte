<script>
import {onMount, createEventDispatcher} from 'svelte'
import Item from './item.svelte'

const dispatch = createEventDispatcher();

export let multiple;
export let items;

function select(item) { 
    let state = _items[item.detail]?.selected
    _items[item.detail].selected = !state
    if(!multiple) {
        _items.forEach((x, i) => {
            if(i != item.detail) {
                delete x.selected
            }
        })
    }
    dispatch("selected", _items)
}

let _items = items;
onMount(() => {
    let isSelected  = _items?.filter(x => x.selected)?.length > 0
    // if no items are selected by default, select the first item
    if(!isSelected && !multiple) {
        _items[0].selected = true
    }
})

</script>


{#if _items}

    <div class="select-container fl-co">

        {#each _items as item, i}

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



