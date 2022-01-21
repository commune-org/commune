<script>
import { logo as user } from '../../../switcher/home/logo'

export let user_id;
export let info;

$: username =  strip(user_id)

function strip(id) {
    let x= id?.split(":")[0]
    return x?.substring(1)
}

$: displayNameExists = info?.display_name?.length > 0
$: avatarExists = info?.avatar_url?.length > 0

$: name = displayNameExists ? info.display_name : username

$: avatar = `${homeServer}/_matrix/media/r0/download/${info?.avatar_url?.substring(6)}`

</script>

<div class="member-item ph2 flex">
    <div class="pr3">
        {#if avatarExists}
            <div class="profile-avatar ncol bg-img"
                style="background-image: url({avatar});">
            </div>
        {:else}
            <div class="profile-avatar gr-default">
                <div class="log gr-center">
                    {@html user}
                </div>
            </div>
        {/if}
    </div>
    <div class="flex gr-center">
        {name} 
    </div>
</div>

<style>
.member-item {
    padding-top: 0.5rem;
    padding-bottom: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
}
.member-item:hover {
    background-color: var(--background-6);
}

.profile-avatar {
    width: 32px;
    height: 32px;
    background-color: var(--avatar);
    border-radius: 50%;
    transition: 0.1s;
}
.profile-avatar:hover {
}
.ncol {
    background-color: transparent;
}
.username {
    color: var(--white);
    cursor: pointer;
}
.username:hover {
    text-decoration: underline;
}
.log {
    width: 16px;
    fill: var(--text);
}
</style>
