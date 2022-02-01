import { store } from '../../../store/store.js'
import {keymap} from "prosemirror-keymap"
import {history} from "prosemirror-history"
import {baseKeymap} from "prosemirror-commands"
import {Plugin} from "prosemirror-state"
import {dropCursor} from "prosemirror-dropcursor"
import {gapCursor} from "prosemirror-gapcursor"
import {menuBar} from "../plugins/menu/index.js"

import {buildMenuItems} from "./menu"

import {buildKeymap} from "./keymap"
import {buildBoardKeymap} from "./keymap-board.js"
import {buildInputRules} from "./inputrules"

import Placeholder from '../plugins/placeholder.js'

import {getMentionsPlugin} from '../plugins/mentions/index'

export {buildMenuItems, buildKeymap, buildInputRules, buildBoardKeymap}

import { hash, topics, user } from '../../../utils/icons.js'

// !! This module exports helper functions for deriving a set of basic
// menu items, input rules, or key bindings from a schema. These
// values need to know about the schema for two reasons—they need
// access to specific instances of node and mark types, and they need
// to know which of the node and mark types that they know about are
// actually present in the schema.
//
// The `exampleSetup` plugin ties these together into a plugin that
// will automatically enable this basic functionality in an editor.

// :: (Object) → [Plugin]
// A convenience plugin that bundles together a simple menu with basic
// key bindings, input rules, and styling for the example schema.
// Probably only useful for quickly setting up a passable
// editor—you'll need more control over your settings in most
// real-world situations.
//
//   options::- The following options are recognized:
//
//     schema:: Schema
//     The schema to generate key bindings and menu items for.
//
//     mapKeys:: ?Object
//     Can be used to [adjust](#example-setup.buildKeymap) the key bindings created.
//
//     menuBar:: ?bool
//     Set to false to disable the menu bar.
//
//     history:: ?bool
//     Set to false to disable the history plugin.
//
//     floatingMenu:: ?bool
//     Set to false to make the menu bar non-floating.
//
//     menuContent:: [[MenuItem]]
//     Can be used to override the menu content.
  //
  //
  //
  //


/*
var getTagSuggestionsHTML = items => '<div class="suggestion-item-list scrl">'+
  items.map(i => '<div class="suggestion-item pv2 ph3">'+i.unicode+i.shortcode+'</div>').join('')+
'</div>';

var getEmojiSuggestionsHTML = items => '<div class="suggestion-item-list scrl">'+
  items.map(i => '<div class="suggestion-item pv2 ph3">'+i.unicode+i.shortcode+'</div>').join('')+
'</div>';
*/

let getTagSuggestionsHTML = (items, text) => {
  let el = ``
  items.forEach(item => {
    let icon = item?.room_type == `topics` ? topics : 
      item?.room_type == 'chat' ? hash : ``
    let e = `<div class="suggestion-item s-ri flex pv1 ph3">${icon}<div class="gr-center"></div><div class="gr-center ml3">${item.name}</div></div>`
    el = el + e
  })

  let sub = `<div class="smaller flex ph3 pv2"><div class="upcs">rooms matching</div><div class="wh ml2">#${text}</div></div>`

  return `<div class="suggestion-item-container flex flex-column">${sub}<div class="suggestion-item-list">${el}</div></div>`
}

function strip(url) {
  return `${homeServer}/_matrix/media/r0/download/${url.substring(6)}`
}


let getMentionSuggestionsHTML = (items, text) => {
  let el = ``
  items.forEach(item => {
    let name = item.display_name?.length > 0 ?
      item.display_name : item?.username


    let pi = ``
    
    if(item.avatar_url?.length > 0) {
      let avatar = strip(item.avatar_url)
      pi = `<div class="p-av bg-img" style="background-image: url(${avatar});"></div>`
    } else {
      pi = `<div class="p-av gr-default" style="background-color: var(--avatar);"><div class="log gr-center">${user}</div></div>`
    }


    let e = `<div class="suggestion-item flex pv1 ph3"><div class="gr-center">${pi}</div><div class="fw-500 gr-center ml3">${name}</div></div>`
    el = el + e
  })

  let sub = `<div class="smaller flex ph3 pv2"><div class="slb upcs">members</div></div>`

  return `<div class="suggestion-item-container flex flex-column">${sub}<div class="suggestion-item-list">${el}</div></div>`
}


let getEmojiSuggestionsHTML = (items, text) => {
  let el = ``
  items.forEach(item => {
    let e = `<div class="suggestion-item flex pv1 ph3"><div class="gr-center">${item.unicode}</div><div class="fw-500 gr-center ml3">${item.shortcode}</div></div>`
    el = el + e
  })

  let sub = `<div class="smaller flex ph3 pv2"><div class="slb upcs">emoji matching</div><div class="fw-500 wh ml2">:${text}</div></div>`

  return `<div class="suggestion-item-container flex flex-column">${sub}<div class="suggestion-item-list">${el}</div></div>`
}


var mentionOpts = {
    getSuggestions: (type, text, done) => {
      setTimeout(() => {
        if (type == 'mention') {
          let x = []
          store.subscribe(v => {
            let users = v.members[room?.server_id]
            for (const [user_id, info] of Object.entries(users)) {
              if(info.username.includes(text)) {
                x.push({
                  user_id: user_id,
                  username: info.username,
                  display_name: info.display_name,
                  avatar_url: info.avatar_url,
                })
              }
            }
          })
          x?.sort((a, b) => (a.username > b.username) ? 1 : -1)
          if(x.length > 0) {
            done(x)
          }
        } else if(type === "emoji"){
          if(text?.length > 1) {
            let m = window.emoji.filter(x => x.shortcode.includes(text))
            if(m.length > 0) {
              done(m.slice(0, 9))
            }
          }
        } else if(type === 'tag'){
          console.log("yesssss")
          let x = []
          store.subscribe(v => {
            let account = v.accounts.filter(x => x.user_id == v.active_account)[0]
            let server = account.servers.filter(x => x.room_id == room?.server_id)[0]
            console.log(server)
            for (const room of server.rooms) {
              if(room.name.includes(text)) {
                x.push({
                  name: room.name,
                  room_type: room.room_type,
                  path: `/${room.server_alias}/${room.alias}`
                })
              }
            }
          })
          if(x.length > 0) {
            done(x)
          }
        }
      }, 0);
    },
    getSuggestionsHTML: (items, type, text) =>  {
      if (type === 'mention') {
        return getMentionSuggestionsHTML(items, text)
      } else if (type === 'tag') {
        return getTagSuggestionsHTML(items, text)
      } else if (type === 'emoji') {
        return getEmojiSuggestionsHTML(items, text)
      }
    }
}

let room;

export function setup(options) {

  room = options.room

  let keymaps = keymap(buildKeymap(options.schema, options.mapKeys))

  if(options.topics || options.show) {
    keymaps = keymap(buildBoardKeymap(options.schema, options.mapKeys))
  }

  let plugins = [
    buildInputRules(options.schema),
    keymaps,
    keymap(baseKeymap),
    dropCursor(),
    gapCursor(),
    Placeholder(options.placeholder)
  ]

  if (options.menuBar !== false && options.topics && !options.show)
    plugins.push(menuBar({floating: options.floatingMenu !== false,
                          content: options.menuContent || buildMenuItems(options.schema).fullMenu}))

  if (options.history !== false)
    plugins.push(history())


  mentionOpts.id = options.id
  plugins.unshift(getMentionsPlugin(mentionOpts))

  let props = {
      attributes: {class: "chat-input"}
  }

  if(options.inline) {
    options.containerClass = `inline-editor scrl ${options.containerClass}`
  }

  if(options.topics && !options.show) {
    options.containerClass = `topics-editor scrl ${options.containerClass}`
  }

  if(options.containerClass && options.containerClass.length > 0) {
    props.attributes = {class: options.containerClass}
  }

  return plugins.concat(new Plugin({
    props: props,
  }))
}
