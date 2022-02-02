<script>
import { send } from '../../../../utils/icons.js'
import { createEventDispatcher } from 'svelte'

const dispatch = createEventDispatcher()

export let noInput;

function sendMessage() {
    dispatch("send", true)
}

</script>

<div class="inc">


    <div class="sep gr-center mr3">
    </div>


    <div class="sic gr-center mr3" 
        class:ac={!noInput}
        on:click={sendMessage}>
        {@html send}
    </div>

</div>


<style>

.inc {
    display: grid;
    grid-template-columns: auto auto;
}

.sep {
    height: 26px;
    border-left: 1px solid var(--background-10);
}

.sic {
    width: 20px;
    height: 20px;
    fill: var(--background-10);
    cursor: not-allowed;
    transition: 0.1s;
}

.ac {
    fill: var(--green);
    cursor: pointer;
}

.ac:hover {
    opacity: 0.5;
}

</style>
