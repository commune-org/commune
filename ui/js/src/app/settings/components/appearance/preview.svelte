<script>
import { store } from '../../../store/store'

import { logo as user } from '../../../switcher/home/logo'

$: compact = $store.settings.displayMode === "compact"

let events = [
    {
        content: "TAKE, O take those lips away",
    },
    {
        content: "That so sweetly were forsworn,",
        skip:true,
    },
    {
        content: "And those eyes, the break of day,",
    },
    {
        content: "Lights that do mislead the morn:",
        skip:true,
    },
    {
        content: "But my kisses bring again,",
    },
    {
        content: "Bring again—",
        skip:true,
    },
    {
        content: "Seals of love, but seal’d in vain,",
    },
    {
        content: "Seal’d in vain!",
        skip:true,
    },
]

$: avatarExists = identity?.avatar_url?.length > 0

$: avatar = `${homeServer}/_matrix/media/r0/download/${identity?.avatar_url?.substring(6)}`

$: username =  strip(identity?.user_id)

function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}
$: displayNameExists = identity?.display_name?.length > 0
$: name = displayNameExists ? identity.display_name : username

function when() {
    let now = new Date()
    let ampm= now.getHours() >= 12 ? `PM` : `AM`
    return `${now.getHours()}:${now.getMinutes()} ${ampm}`
}

</script>

<div class="p-c"> 
    <div class="p-i fl-co pa3"> 
        {#each events as event}
            <div class="event-item mb1 flex" 
            class:mb3={event.skip && !compact}> 

            {#if !compact}
                

                <div class="pi-co mr3">

                    {#if !event.skip && !compact}
                        

                        {#if avatarExists}
                            <div class="profile-avatar ncol bg-img"
                                style="background-image: url({avatar});">
                            </div>
                        {:else}
                            <div class="profile-avatar gr-default">
                                <div class="log gr-default gr-center">
                                    {@html user}
                                </div>
                            </div>
                        {/if}

                    {/if}

                </div>

            {:else}
                        <div class="mb1 flex mr2"> 
                            <div class="when"> 
                                {when()}
                            </div>
                            <div class="name ml2"> 
                                <strong>{name}</strong>
                            </div>
                        </div>

            {/if}


                <div class="fl-co flex-one"> 
                    {#if !event.skip && !compact}
                        <div class="mb1 flex"> 
                            <div class="name"> 
                                <strong>{name}</strong>
                            </div>
                            <div class="ml2 when"> 
                            Today at {when()}
                            </div>
                        </div>
                    {/if}
                    <div class=""> 
                        {event.content}
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>


<style> 
.p-c {
    height: 180px;
    width: 100%;
    border: 1px solid var(--background-1);
    background-color: var(--background-2);
    border-radius: 5px;
    overflow-y: hidden;
}
.p-i {
    word-wrap: break-word;
    pointer-events: none !important;
    margin-top: -2.5rem;
}

.pi-co {
    min-width: 40px;
}

.profile-avatar {
    width: 40px;
    height: 40px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
    cursor: pointer;
}
.ncol {
    background-color: transparent;
}
.when {
    margin-top: 4px;
    font-size: 0.8rem;
    color: var(--text-muted);
}
.name {
    color: var(--white);
}
</style>
