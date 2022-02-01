<script>
import { onMount, onDestroy, tick } from 'svelte'
import { store } from '../../store/store.js'
import { twemoji } from '../../utils/twemoji'
import { makeid } from '../../utils/utils.js'
import { search, close, em0, em1, em2, em3, em4, em5, em6, em7, em8 } from '../../utils/icons.js'

import EMOJIBASE from 'emojibase-data/en/compact.json';
import SHORTCODES from 'emojibase-data/en/shortcodes/joypixels.json';

EMOJIBASE.forEach(item => {
    let code = ``
    let shortcodes = SHORTCODES[item.hexcode]
    if(Array.isArray(shortcodes)) {
        shortcodes.forEach(x => {
            code = code + `:${x}: `
        })
    } else {
        code = `:${shortcodes}:`
    }
    item.shortcode = code
    item.baseunicode = item.unicode
})

let CONSTRUCTED = EMOJIBASE
window.emoji = CONSTRUCTED

$: people = CONSTRUCTED.filter(x => x.group == 0 || x.group == 1)
$: nature = CONSTRUCTED.filter(x => x.group == 3)
$: food = CONSTRUCTED.filter(x => x.group == 4)
$: travel = CONSTRUCTED.filter(x => x.group == 5)
$: activities = CONSTRUCTED.filter(x => x.group == 6)
$: objects = CONSTRUCTED.filter(x => x.group == 7)
$: symbols = CONSTRUCTED.filter(x => x.group == 8)
$: flags = CONSTRUCTED.filter(x => x.group == 9)

$: emojis = {
    people: {
        emoji: people,
        icon: em1,
    },
    nature: {
        emoji: nature,
        icon: em2,
    },
    food: {
        emoji: food,
        icon: em3,
    },
    travel: {
        emoji: travel,
        icon: em4,
    },
    activities: {
        emoji: activities,
        icon: em5,
    },
    objects: {
        emoji: objects,
        icon: em6,
    },
    symbols: {
        emoji: symbols,
        icon: em7,
    },
    flags: {
        emoji: flags,
        icon: em8,
    },
}

export let container;
export let isActive;

function formatEmoji(e) {
    if(e) {
        //let code = String.fromCodePoint(parseInt (e, 16))
        return twemoji.parse(e, {
            base: '/static/img/emoji/',
            className: 'emo',
            ext: '.svg',
            size: '72x72'
        })
    }
    return ``
}

$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

function react(e) {
    let key = e.target.getAttribute('alt')
    let shortcode = e.target.getAttribute('title')

    //ignore if user already reacted with the same key
    if(userReactions?.length > 0) {
        let reacted = false;
        console.log("checking")
        userReactions.forEach(reaction => {
            if(reaction?.content?.['m.relates_to']?.key == key){
                reacted = true
            }
        })
        if(reacted) {
            kill()
            return
        }
    }


    if(isDispatch) {
        window.dispatchEmoji(key)
        kill()
        return
    }

    let tempid = makeid(16)

    let content = {
        'm.relates_to': {
            event_id: eventID,
            key: key,
            rel_type: 'm.annotation',
        }
    }


    let newEvent= {
        "age": 0,
        "type": "m.reaction",
        "room_id": roomID,
        "sender": $store.active_account,
        "content": content,
        "origin_server_ts": new Date(),
        "unsigned": {
            "age": 0
        },
        "transaction_id": tempid,
        "user_id": $store.active_account,
        "delivered": false,
    }

    //store.updateReactions(room.room_id, newEvent)
    store.addEventToRoom(roomID, newEvent)

    let eventType = `m.reaction`
    let room_id = roomID

    if(isTopic) {
        console.log(isTopic, slug)
        //eventType = `${slug}.reaction`
        room_id = topicRoomID
        //newEvent.content['m.relates_to']['topic_reply']
    }


    matrix.sendEvent(room_id, eventType, content, tempid, (err, res) => {
        console.log(res)
    });

    active = false

    let storedEmojis = localStorage.getItem('emojis')
    if(!storedEmojis) {
        let em = [{emoji: key, shortcode: shortcode, frequency: 1}]
        localStorage.setItem('emojis', JSON.stringify(em))
    } else {
        let parsed = JSON.parse(storedEmojis)

        let ind = parsed.findIndex(x => x.emoji == key)
        if(ind == -1) {
            parsed.push({emoji: key, shortcode: shortcode, frequency: 1})
        } else {
            parsed[ind].frequency = parsed[ind].frequency + 1
        }
        localStorage.setItem('emojis', JSON.stringify(parsed))
    }
    resetFrequent()
}

let roomID;
let eventID;

let isDispatch = false;
let isInline = false;

let isTopic = false;
let topicRoomID = null;
let slug = null;

let frequent = []

function resetFrequent() {
    let storedEmojis = localStorage.getItem('emojis')
    if(storedEmojis) {
        let parsed = JSON.parse(storedEmojis)
        parsed?.sort((a, b) => (a.frequency < b.frequency) ? 1 : -1)

        if(parsed?.length > 0) {
            frequent = parsed
        }
    } else {
        frequent = [
            {
                emoji: '👍️', 
                shortcode: ":thumbsup: :+1: :thumbup: ", 
                frequency: 1,
            },
            {
                emoji: '👀', 
                shortcode: ":eyes:",
                frequency: 1,
            },
            {
                emoji: '😆', 
                shortcode: ":laughing: :satisfied: ",
                frequency: 1,
            },
            {
                emoji: '💯', 
                shortcode: ":100:",
                frequency: 1,
            },
            {
                emoji: '🍉', 
                shortcode: ":watermelon:",
                frequency: 1,
            },
            {
                emoji: '🍴', 
                shortcode: ":fork_and_knife:",
                frequency: 1,
            },
            {
                emoji: '😋', 
                shortcode: ":yum:",
                frequency: 1,
            },
            {
                emoji: '😩', 
                shortcode: ":weary:",
                frequency: 1,
            },
            {
                emoji: '😫', 
                shortcode: ":tired_face:",
                frequency: 1,
            },
            {
                emoji: '💩', 
                shortcode: ":poop: :shit: :hankey: :poo: ",
                frequency: 1,
            },
        ]
        localStorage.setItem('emojis', JSON.stringify(frequent))
    }
}


let userReactions;

onMount(() => {

    resetFrequent()

    if(isActive) {
        active = true
        return
    }
    window.toggleEmojiPicker = (opts) => {
        isInline = false
        isDispatch = false
        if(opts?.dispatch) {
            isDispatch = true
        }
        if(opts?.inline) {
            isInline = true
        }
        if(opts.topics) {
            isTopic = true
            topicRoomID = opts.room_id
        }
        if(opts.slug) {
            slug = opts.slug
        }
        userReactions = opts?.userReactions
        roomID = opts?.room?.room_id
        eventID = opts?.event?.event_id
        container = opts?.container
        query = ''
        searchInput.value = null
        resetHovered()
        emojiContainer.scrollTop = 0
        if(!active) {
            active = true
        } else {
            repeat = true
        }
    }
    window.killEmojiPicker = () => {
        kill()
    }
})

async function focusSearchInput() {
    await tick()
    searchInput.focus()
}

$: if(active) {
    $store.emoji = {
        active: true,
        event: eventID,
    }
    setTimeout(() => {
        document.addEventListener('click', killMe)
        emojiContainer.addEventListener('scroll', highlightNav)
        document.addEventListener('keydown', escape)
    }, 10)
    focusSearchInput()
} else {
    $store.emoji = {
        active: false,
        event: null,
    }
    document.removeEventListener('click', killMe)
    if(emojiContainer) {
        emojiContainer.removeEventListener('scroll', highlightNav)
    }
    document.removeEventListener('keydown', escape)
}

onDestroy(() => {
    document.removeEventListener('click', killMe)
    if(emojiContainer) {
        emojiContainer.removeEventListener('scroll', highlightNav)
    }
    document.removeEventListener('keydown', escape)
})

let escape = (e) => {
    if(e.code == 'Escape' || e.key == 'Escape') {
        kill()
    }
}

let repeat = false;

let highlighted = 'em-frequently'

let highlightNav = (e) => {
    let el0 = document.getElementById('em-frequently')
    let el1 = document.getElementById('em-people')
    let el2 = document.getElementById('em-nature')
    let el3 = document.getElementById('em-food')
    let el4 = document.getElementById('em-travel')
    let el5 = document.getElementById('em-activities')
    let el6 = document.getElementById('em-objects')
    let el7 = document.getElementById('em-symbols')
    let el8 = document.getElementById('em-flags')
    if(emojiContainer.scrollTop >= el0.offsetTop) {
        highlighted = `em-frequently`
    }
    if(emojiContainer.scrollTop >= el1.offsetTop) {
        highlighted = `em-people`
    }
    if(emojiContainer.scrollTop >= el2.offsetTop) {
        highlighted = `em-nature`
    }
    if(emojiContainer.scrollTop >= el3.offsetTop) {
        highlighted = `em-food`
    }
    if(emojiContainer.scrollTop >= el4.offsetTop) {
        highlighted = `em-travel`
    }
    if(emojiContainer.scrollTop >= el5.offsetTop) {
        highlighted = `em-activities`
    }
    if(emojiContainer.scrollTop >= el6.offsetTop) {
        highlighted = `em-objects`
    }
    if(emojiContainer.scrollTop >= el7.offsetTop) {
        highlighted = `em-symbols`
    }
    if(emojiContainer.scrollTop >= el8.offsetTop) {
        highlighted = `em-flags`
    }
}


let killMe = (e) => {
    if(e.target !== picker && 
        !picker?.contains(e.target) && 
        e.target !== searchInput && 
        !e.target.classList.contains('tone-option') && 
        !e.target.classList.contains('emoji-nav-item') && 
        !e.target.classList.contains('search-icon') && 
        e.target !== searchContainer){
        if(!repeat) {
            kill()
        } else {
            repeat = false
        }
    }
}


let active = false

function kill(e) {
    active = false
    query = ''
    searchInput.value = null
    emojiContainer.scrollTop = 0
    resetHovered()
    if(isDispatch) {
        window.killEmoji(true)
    }
    isDispatch = false
    isTopic = false;
    isInline = false;
    topicRoomID = null;
    slug = null;
}

let picker;

$: bounding = container?.getBoundingClientRect()

$: inside = (bounding?.top + picker?.getBoundingClientRect()?.height) < document.body.clientHeight 

$: at = document.body.clientHeight - picker?.getBoundingClientRect()?.height-80

$: topNormal = inside ? bounding?.top : at

$: calcTop = bounding?.top - 444

$: top = isDispatch ? calcTop : topNormal

$: rightOffset = isDispatch ? -14 : isInline ? -444 : 40
$: right = document.body.clientWidth - bounding?.right  + rightOffset

let emojiContainer;

let searchContainer;
let searchInput;
let query = '';

$: filtering = query.length > 0

$: if(filtering) {
    highlighted = ``
}

$: filtered = EMOJIBASE.filter(x => x.shortcode.includes(query))

function killFilter() {
    searchInput.value = null
    query = ''
    filterEmoji()
    focusSearchInput()
}

function filterEmoji() {
    if(searchInput.value.length == 0) {
        highlighted = 'em-frequently'
    }
    let el = document.getElementById(`em-frequently`)
    if(el) {
        el.scrollIntoView(true)
    }
}

let moved = false;

let placeholder = `Find the perfect emoji`
$: hoveredEmoji = frequent[0]?.emoji
$: hoveredShortcode = frequent[0]?.shortcode

function resetHovered() {
    placeholder = `Find the perfect emoji`
    hoveredEmoji = frequent[0]?.emoji
    hoveredShortcode = frequent[0]?.shortcode
}

function changePlaceholder(e) {
    if(e.target?.className?.includes?.('emoji-key')) {
        moved = true
        placeholder = e.target.title
        hoveredEmoji = e.target.getAttribute('alt')
        hoveredShortcode = e.target.title
    }
}

function toggleSection(e) {
    let el = document.querySelector(`#${e}`)
    if(el) {
        if(el?.classList.contains('no-dis')) {
            el.classList.remove('no-dis')
        } else {
            el.classList.add('no-dis')
        }
    }
}

function cat1() {
    let el = document.getElementById(`em-frequently`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat2() {
    let el = document.getElementById(`em-people`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat3() {
    let el = document.getElementById(`em-nature`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat4() {
    let el = document.getElementById(`em-food`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat5() {
    let el = document.getElementById(`em-travel`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat6() {
    let el = document.getElementById(`em-activities`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat7() {
    let el = document.getElementById(`em-objects`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat8() {
    let el = document.getElementById(`em-symbols`)
    if(el) {
        el.scrollIntoView(true)
    }
}
function cat9() {
    let el = document.getElementById(`em-flags`)
    if(el) {
        el.scrollIntoView(true)
    }
}

let baseTones = [
    "👏",
    "👏🏻",
    "👏🏼",
    "👏🏽",
    "👏🏾",
    "👏🏿"
]


let tones = [
    "👏",
    "👏🏻",
    "👏🏼",
    "👏🏽",
    "👏🏾",
    "👏🏿"
]


let skin = localStorage.getItem('emoji-tone')
if(skin) {
    let t = tones.filter(t => t != tones[skin])
    t.unshift(tones[skin])
    tones = t
    setTones()
}

function setTones() {
    let skin = localStorage.getItem('emoji-tone')

    if(skin == 0) {
        console.log("we are resetting to base tone")
        CONSTRUCTED.forEach(item => {
            item.unicode = item.baseunicode
        })
    } else {
        CONSTRUCTED.forEach(item => {
            if(skin && item?.skins?.length > 0) {
                let ind = parseInt(skin) - 1
                if(Array.isArray(item.skins)) {
                    item.unicode = item?.skins?.[ind]?.unicode
                }
            }
        })
    }

    CONSTRUCTED = CONSTRUCTED
    window.emoji = CONSTRUCTED
}

/*
$: skinTone = findSkin()

function findSkin() {
    let clap = people.filter(x => x.hexcode == "1F44F")[0]
    let skin = localStorage.getItem('emoji-tone')
    if(skin) {
        clap.unicode = clap?.skins[skin].unicode
        return clap
    } else {
        return clap
    }
}
*/

let tonesToggled = false;

function toggleTones() {
    tonesToggled = !tonesToggled
}

$: if(tonesToggled) {
    setTimeout(() => {
        document.addEventListener('click', killTones)
    }, 10)
} else {
    document.removeEventListener('click', killTones)
}

let killTones= (e) => {
    if(!e.target.classList.contains('tone-option')){
        tonesToggled = false
    }
}

function selectTone(tone) {

    tonesToggled = false
    let ind = baseTones.findIndex(t => t == tone)
    localStorage.setItem('emoji-tone', ind)

    tones = tones.filter(t => t != tone)
    tones.unshift(tone)
    tones = tones
    setTones()
}

</script>

<div class="layer" class:inactive={!active}>

    <div class="emoji-container" 
    class:nd={!isDispatch && !isInline}
    class:ndr={isInline}
    class:fly={active && !isDispatch}
    style="--top:{top};--right:{right}"
    bind:this={picker}>

        <div class="emoji-search pa2 flex" 
            bind:this={searchContainer}>

            <div class="gr-center flex-one relative">
                <input 
                class="pa2"
                bind:this={searchInput}
                bind:value={query}
                on:input={filterEmoji}
                placeholder={placeholder} />
                {#if filtering}
                    <div class="search-icon pointer sia" 
                    on:click|stopPropagation={killFilter}>
                        {@html close}
                    </div>
                {:else}
                    <div class="search-icon">
                        {@html search}
                    </div>
                {/if}
            </div>

            <div class="gr-center pa2 skintone relative" >

                {#if !tonesToggled}
                    <div class="tone-options"
                        on:click|self={toggleTones}>
                        <div class="tone-option gr-center"
                            on:click|self={toggleTones}>
                            {tones[0]}
                        </div>
                    </div>
                {:else}
                    <div class="tone-options toggled">
                        {#each tones as tone, i (i)}
                            <div class="tone-option gr-center" 
                                on:click={selectTone(tone)}
                                class:pt2={i!=0}>
                                {tone}
                            </div>
                        {/each}
                    </div>
                {/if}
            </div>

        </div>

        <div class="inner-container" >

            <div class="emoji-nav flex flex-column">
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-frequently'}
                    on:click={cat1}>
                    {@html em0}
                </div>
                <div class="nav-sep mv2"></div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-people'}
                    on:click={cat2}>
                    {@html em1}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-nature'}
                    on:click={cat3}>
                    {@html em2}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-food'}
                    on:click={cat4}>
                    {@html em3}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-travel'}
                    on:click={cat5}>
                    {@html em4}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-activities'}
                    on:click={cat6}>
                    {@html em5}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-objects'}
                    on:click={cat7}>
                    {@html em6}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-symbols'}
                    on:click={cat8}>
                    {@html em7}
                </div>
                <div class="emoji-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-flags'}
                    on:click={cat9}>
                    {@html em8}
                </div>
            </div>

            <div class="emojis-container">

                {#if filtering}
                    <div class="emojis scrl" >
                        <div class="flex flex-column">
                            <div class="con">
                                {#each filtered as item (item.order)}
                                    <li class="emoji-item gr-default" >
                                        <div class="emoji-key"
                                        alt={item.unicode}
                                        title={item.shortcode}
                                        on:click={react}>
                                            {@html item.unicode}
                                        </div>
                                    </li>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/if}

                <div class="emojis scrl" 
                    class:no-dis={filtering}
                    on:mouseover={changePlaceholder}
                    bind:this={emojiContainer}>

                    {#if frequent.length > 0}

                        <div id="em-frequently" class="flex flex-column mb3">
                            <div class="emoji-title pl1 flex">
                                <div class="tin-c gr-center">
                                    {@html em0}
                                </div>
                                <div class="tt ml1 gr-center">
                                    frequently used
                                </div>
                            </div>
                            <div class="con">
                                {#each frequent as item, i (item.frequency)}
                                    <li class="emoji-item gr-default"
                                        class:high={i == 0 && !moved}>
                                        <div class="emoji-key"
                                        alt={item.emoji}
                                        title={item.shortcode}
                                        on:click={react}>
                                            {@html item.emoji}
                                        </div>
                                    </li>
                                {/each}
                            </div>
                        </div>

                    {/if}

                    {#each Object.entries(emojis) as [title, set] ,i}
                        <div id="em-{title}" class="flex flex-column mb3">
                            <div class="emoji-title pl1 flex">
                                <div class="tin-c gr-center">
                                    {@html set.icon}
                                </div>
                                <div class="tt ml1 gr-center">
                                    {title}
                                </div>
                            </div>
                            <div class="con">
                                {#each set.emoji as item (item.order)}
                                    <li class="emoji-item gr-default">
                                        <div class="emoji-key"
                                        alt={item.unicode}
                                        title={item.shortcode}
                                        on:click={react}>
                                            {@html item.unicode}
                                        </div>
                                    </li>
                                {/each}
                            </div>
                        </div>
                    {/each}
                </div>

                <div class="emoji-preview flex">
                    <div class="gr-center ml2 hov-e">
                        {hoveredEmoji} 
                    </div>
                    <div class="hov-s ml2 gr-center flex-one">
                        {hoveredShortcode}
                    </div>
                </div>

            </div>


        </div>

    </div>

</div>

<style>
.inactive {
    visibility: hidden;
    pointer-events: none!important;
}
.layer {
    position: fixed;
    pointer-events: none!important;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: none!important;
    z-index: 1002;
}

.emoji-search {
    border-bottom: 1px solid var(--background-1);
}

.emoji-container {
    position: absolute;
    top: var(--top);
    right: var(--right);
    pointer-events: auto;
    overflow: hidden;
    height: 420px;
    width: 424px;
    border-radius: 8px;
    background-color: var(--background-2);
    color: var(--text);
    font-size: 14px;
    line-height: 1.4;
    outline: 0;
    transition-property: transform, visibility, opacity;
    box-shadow: 0 10px 20px rgba(0,0,0,.18);
    z-index: 10000;
    overflow: hidden;
    display: grid;
    grid-template-rows: auto 1fr;
}

.nd {
    margin-right: 1rem;
    transition: margin-right 0.05s;
}

.ndr {
    margin-left: 1rem;
    transition: margin-right 0.05s;
}


.fly {
    margin-right: 0.25rem
}

.emojis {
    padding: 0 0.5rem;
    overflow: hidden auto;
    display: grid;
    grid-template-rows: repeat();
    width: 100%;
    height: 100%;
    scroll-behavior: auto;
    background-color: var(--background-2);
    position: relative;
}

.inner-container {
    display: grid;
    grid-template-columns: auto 1fr;
    overflow: hidden;
}

.emojis-container {
    display: grid;
    grid-template-rows: 1fr auto;
    overflow: hidden;
    background-color: var(--background-5);
}

.emoji-preview {
    height: 48px;
}

.emoji-nav {
    background-color: var(--background-1);
    width: 48px;
    padding: 0.5rem;
}

.emoji-nav-item {
    width: 32px;
    height: 32px;
    cursor: pointer;
    padding: 0.25rem;
}

.emoji-nav-item:hover {
    background-color: var(--background-2);
    border-radius: 3px;
}

.nav-sep {
    height: 1px;
    background-color: var(--background-3);
    width: 100%;
}

.con {
    display: grid;
    grid-template-columns: repeat(9, 40px);
}

.emoji-title {
    text-transform: uppercase;
    font-size: 12px;
    font-weight: 600;
    color: var(--text-light);
    line-height: 1;
    position: sticky;
    position: -webkit-sticky;
    top: 0;
    background-color: var(--background-2);
    padding: 0.5rem 0;
    z-index: 25;
}

.tt {
    vertical-align: bottom;
}

.hov-e {
    font-size: 28px;
}

.hov-s {
    font-weight: bold;
    font-size: 1.1rem;
}

.emoji-item {
    padding: 0.25rem;
    cursor: pointer;
    width: 40px;
    height: 40px;
    border-radius: 3px;
}

.emoji-key {
    font-size: 32px;
    line-height: 32px;
    font-family: Twemoji;
    justify-self: center;
    align-self: center;
}

.high {
    background-color: var(--background-4);
}

.emoji-item:hover {
    background-color: var(--background-4);
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--background-1) transparent;
    scroll-behavior: auto;
}

.scrl::-webkit-scrollbar {
  width: 4px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: var(--background-1);
}

input {
    font-size: 0.9rem;
    padding-left: 0.5rem;
    width: 100%;
    height: 100%;
}

.no-dis {
    display: none;
}
.dis {
    display: block;
}
.skintone {
    font-size: 24px;
    line-height: 24px;
    cursor: pointer;
    width: 36px;
}
.tone-options {
    position: absolute;
    top: -8px;
    left: 4px;
    z-index: 35;
    padding: 0.25rem;
    border: 1px solid transparent;
}
.toggled {
    border: 1px solid var(--background-1);
    background-color: var(--background-5);
    border-radius: 3px;
}
.highlight {
    background-color: var(--background-3);
    border-radius: 3px;
}
.search-icon {
    position:absolute;
    top: 6px;
    right: 8px;
    fill: var(--text-muted);
    height: 20;
    width: 20;
}
.sia {
    fill: var(--text);
}
.tin-c {
    width: 16px;
    height: 16px;
}
</style>

