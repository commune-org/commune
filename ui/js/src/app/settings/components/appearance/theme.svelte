<script>
import Select from '../../../components/ui/select/select.svelte'
import {onMount} from 'svelte'

let ready = false;

onMount(() => {
    ready = true
})

let lightMode = localStorage.getItem("light-mode");
let isLight = lightMode && lightMode == "true"

let modes = [
    {
        name: "dark",
        caption: "Dark",
    },
    {
        name: "light",
        caption: "Light",
        selected: true
    },
    {
        name: "sync",
        caption: "Sync with computer"
    },
]

function selected(e) {
console.log(e.detail)
}

</script>

<div class="">  
    {#if ready}
        <Select 
            items={modes} 
            on:selected={selected}
            multiple={true}
        />
    {/if}
</div>
