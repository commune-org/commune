<script>
import { onMount } from 'svelte';
import { store } from './store/store.js'
import { fade } from 'svelte/transition'
import Switcher from './switcher/switcher.svelte'
import Sidebar from './sidebar/sidebar.svelte'
import View from './view/view.svelte'
import { Router } from 'svelte-navigator'
import Settings from './settings/settings.svelte'
import ImageViewer from './components/image-viewer/image-viewer.svelte'
import EmojiPicker from './components/emoji/emoji.svelte'
import Tooltip from './components/tooltip/tooltip.svelte'
import Alert from './components/alert/alert.svelte'

onMount(() => {
    document.title = `Commune`
    store.activate()
})


$: isMobile = $store.isMobile
$: mobileViewToggled = $store.mobileViewToggled

$: alertsExist = $store.alerts?.length > 0

$: limitDisplayWidth = $store.settings.accessibility.limitDisplayWidth

</script>


<Router>
{#if $store.active}
<div class="foundation" 
    class:alert={alertsExist}
    transition:fade="{{duration: 100}}">

    <Alert/>

    <div class="container"
    class:mw-1600={limitDisplayWidth}
    class:isMobile={isMobile && !mobileViewToggled}>

        <Switcher />

        <Sidebar />

        <View />

        <Settings />

        <ImageViewer />

        <EmojiPicker />


    </div>
</div>
{:else}
    <div class="gr-center fl-co">
        <div class="spinner gr-center"></div>
        <div class="loading-message pa3 center">
            {$store.loadingMessage}...
        </div>
    </div>
{/if}
</Router>

        <Tooltip/>


<style>

.foundation {
    height: 100vh;
    width: 100vw;
    display: grid;
    grid-template-columns: 100%;
    grid-template-rows: 100%;
}

.container {
    height: 100vh;
    width: 100vw;
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: [switcher] auto [sidebar] auto [view] 1fr;
    justify-self: center;
    align-self: center;
    border-right: 1px solid var(--background-3);
    border-left: 1px solid var(--background-3);
}

.mw-1600 {
    max-width: 1600px;
}

.isMobile {
    grid-template-columns: [view] 1fr;
}

</style>


