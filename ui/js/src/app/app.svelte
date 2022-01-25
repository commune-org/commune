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
    let theme = localStorage.getItem("theme");
    if(theme && theme == "light") {
        document.documentElement.classList.add('light')
    } 
    document.title = `Commune`
    store.activate()


})


$: isMobile = $store.isMobile
$: mobileViewToggled = $store.mobileViewToggled

$: alertsExist = $store.alerts?.length > 0

</script>


<Router>
{#if $store.active}
<div class="foundation" 
    class:alert={alertsExist}
    transition:fade="{{duration: 100}}">

    <Alert/>

    <div class="container"
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
.container {
    height: 100vh;
    width: 100vw;
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: [switcher] auto [sidebar] auto [view] 1fr;
}

.isMobile {
    grid-template-columns: [view] 1fr;
}

</style>


