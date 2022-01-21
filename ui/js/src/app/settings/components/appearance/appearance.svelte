<script>
import Preview from './preview.svelte'
import Theme from './theme.svelte'
import Display from './display.svelte'

let displayMode;

function updateMessageDisplay(e) {
    displayMode = e.detail
}

</script>

<div class="appearance-settings-container fl-co">

    <div class="title">
        Appearance
    </div>
    
    <div class="mt3">
        <Preview mode={displayMode}/>
    </div>


    <div class="mt4">
        <Theme/>
    </div>

    <div class="mt3">
        <Display on:message-display={updateMessageDisplay}/>
    </div>


    <div class="">
    </div>

</div>

<style>
.title {
    font-size: 1.2rem;
    font-weight: bold;
    color: var(--white);
}
</style>
