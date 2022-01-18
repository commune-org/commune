<script>
import { store } from '../store/store.js'
import SidebarHeader from './header/header.svelte'
import SidebarProfile from './profile/profile.svelte'
import SidebarView from './view/view.svelte'

let sidebar;


$: isMobile = $store.isMobile
$: mobileViewToggled = $store.mobileViewToggled

$: hideSidebar = isMobile && !mobileViewToggled

</script>

<div class="sidebar no-select" 
    class:no-dis={hideSidebar}
    bind:this={sidebar}>

    <div class="sidebar-container">
        <SidebarHeader />


        <SidebarView />


        <SidebarProfile />

    </div>


</div>

<style>
.sidebar {
    background-color: var(--background-2);
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: [content] 1fr [expander] auto;
    width: 240px;
    position: relative;
}


.sidebar-container {
    background-color: var(--background-2);
    display: grid;
    grid-template-columns: 100%;
    grid-template-rows: [sidebar-header] 48px [sidebar-content] 1fr [sidebar-profile] auto;
    width: 240px;
}

.expand-sidebar {
    position: absolute;
    right: -3px;
    height: 100%;
    width: 8px;
    cursor: ew-resize;
    transition: 0.2s;
}

.expand-sidebar-indicator {
    width: 6px;
    height: 30px;
    border-radius: 4px;
    background-color: var(--background-5);
    position: absolute;
    top: calc(50vh - 15px);
    right: 0;
    cursor: ew-resize;
    transition: 0.2s;
}

.expand-sidebar:hover .expand-sidebar-indicator {
    background-color: var(--primary);
}

.no-dis {
    display: none;
}
</style>
