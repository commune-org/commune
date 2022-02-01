<script>
import { onMount } from 'svelte'
import { navigate } from 'svelte-navigator'

let popupVisible = false
let show = false;

let fetched = false;

let top;
let left;

onMount(() => {
    document.addEventListener('click', e => {
        if(e.target.classList.contains("mention")) {
            e.preventDefault()
            e.stopPropagation()

            let username = e.target.dataset.username

            let path = `/@${username}`
            navigate(path)
        }
        if(e.target.classList.contains("tag")) {
            e.preventDefault()
            e.stopPropagation()
            navigate(e.target.dataset.path, {replace: true})
        }
    })
    document.addEventListener('mouseover', e => {
        if(e.target.classList.contains("mention") ||
        e.target.classList.contains("username")) {
            e.preventDefault()
            e.stopPropagation()
            fetched = false
            popupVisible = false
            show = false

            let username = e.target.dataset.username

            let rect = e.target.getBoundingClientRect()
            top = rect.top
            left = rect.left + rect.width + 10
            show = true
            setTimeout(() => {
                activate(e)
            }, 500)
        }
        if(!e.target.classList.contains("mention") &&
        !e.target.classList.contains("username")) {
            show = false
        }
    })
})

function activate(e) {
    if(show) {
        popupVisible = true
        setTimeout(() => {
            fetched = true
            let rect = e.target.getBoundingClientRect()
            top = rect.top
            if((top + 254) > document.body.clientHeight) {
                top = top - ((top + 254) - (document.body.clientHeight)) - 20
            }
            left = rect.left + rect.width + 10
        }, 100)
    }
}

let container;

$: if(popupVisible) {
    document.addEventListener('click', killMe)
} else {
    document.removeEventListener('click', killMe)
}

let killMe = (e) => {
    if(e.target !== container && 
        !container?.contains(e.target)) {
        show = false
        popupVisible = false
        fetched = false
    }
}


</script>

{#if popupVisible}
<div class="t-layer">
</div>
    <div class="user-info pa2 gr-default"
        class:fetched={fetched}
        bind:this={container}
        style={`--top:${top}px;--left:${left}px;`}>
        {#if fetched} 
            Test
        {:else}
            <div class="spinner gr-center">
            </div>
        {/if}
    </div>
{/if}

<style>
.t-layer {
    position: fixed;
    pointer-events: none!important;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: none!important;
    z-index: 2000;
}

.user-info {
    position: fixed;
    top: var(--top);
    left: var(--left);
    background-color: var(--background-1);
    border-radius: 7px;
    box-shadow: 0 10px 40px rgba(0,0,0,.14);
    width: 60px;
    height: 60px;
    z-index: 2001;
}
.fetched {
    height: 254px;
    width: 300px;
}
</style>
