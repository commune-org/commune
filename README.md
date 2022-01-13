# Commune

Commune is a communications suite built on top of [matrix](https://matrix.org). Commune aims to bring together chat, discussions, email and other interactive apps into a single matrix client. Commune is under active development and not yet ready for production.

Dev instance: [dev.commune.chat](https://dev.commune.chat)

## Roadmap

- [ ] E2E encryption
- [ ] Unified matrix + email identities
- [ ] Easy self-hosting setup
- [ ] Federation + interoperability with [Element](https://app.element.io)
- [ ] OAuth2 layer
- [ ] Websocket layer
- [ ] Directory API
- [ ] Replace Synapse with [Conduit](https://gitlab.com/famedly/conduit)

## Documentation
See docs [here](https://docs.commune.chat).

### Screenshots

![Alt text](/docs/screenshots/screenshot3.png?raw=true "")
![Alt text](/docs/screenshots/screenshot2.png?raw=true "")
![Alt text](/docs/screenshots/screenshot1.png?raw=true "")

### LICENSE
Commune is released under the [PolyForm Noncommercial License](https://polyformproject.org/wp-content/uploads/2020/05/PolyForm-Noncommercial-1.0.0.txt).

## Q&A

> What is Commune actually from the end user perspective?

Commune is a space for two or more people to communicate together, safely and expressively.

It’s basically WhatsApp extended with ‘build your own comms space’ modules. What starts as ephemeral chatter can progressively turn into lasting artifacts of knowledge.

>So when I sign up on commune.chat I actually get a @username:commune.chat MxID?

Commune accounts are regular Matrix accounts, so a Commune @username is actually just @username:commune.chat. We intentionally use only the `localpart` of the MxID to make it easier for new users who are not familiar with the Matrix ID format.

> How much is Matrix behind the scenes?

Matrix is Commune's events and federation layer, so in a sense all the important parts of Commune are Matrix. Commune does have a backend layer (in Go) that sits between the frontend Svelte client and Synapse. This backend is necessary for extending Matrix without modifying Synapse itself.

> Are you going to distribute Commune as a frontend, or as a whole package?

We want to offer Commune as a single package that includes Commune, Synapse, and anything else that might make self-hosting dead simple. The eventual goal is have a one-click installer that shouldn't require advanced tech/server knowledge. Currently, Commune is best suited to those looking to start self-hosting a new matrix homeserver from scratch. Making Commune work on top of existing homeservers is something we'd like to support by the time Commune reaches beta. Commune will still work if pointed at an existing homeserver, but making existing rooms place nicely with Commune's server -> channel -> stream design is going to take additional work.

> How does it compare to leading community platforms like Discord and Discourse?

Since the Commune product manager @erlend-sh also works as a product manager at @Discourse, a disclaimer is necessary: These two applications are not in direct competition. They solve similar problems, but in very different ways.

#### Commune is…

- Personal; 1:1 as default starting point.
- Primarily Decentralised.
- Optimizing for self-hosting (even p2p longterm).
- Ecosystem-driven (Matrix spec).
- Software-license seller ([PolyForm NonCommercial](https://polyformproject.org/licenses/noncommercial/1.0.0/)).

#### Discord/Discourse is…

- Impersonal; many-to-many as default starting point.
- Primarily Centralised.
- Optimizing for cloud-hosting (SaaS).
- Customer-driven
- Software-hosting seller.

> How long has it been in development?

We’ve been quietly working on the Commune app since mid-2021 and discussed product plans ever since Hummingbard (the prototypical predecessor of Commune) [was publicly announced](https://news.ycombinator.com/item?id=26277602).

> Are you going to distribute Commune as a frontend, or as a whole package?

Whole package. While we expect most users will want to use Commune merely as a window into their favorite people and communities, it’s important that deploying your own server is an equally accessible option for users of this software. Come for the client convenience, stay for the server sovereignty.

> Who are you?

We are a team of two. [ahq](https://github.com/ChurchOfTheSubgenius) is developing it full-time. [Erlend](https://github.com/erlend-sh) is directing the product development of the app part-time, whilst also paying ahq a living wage, effectively as a seed investor.

> What are the implications of the PolyForm NonCommercial license?

The two key restrictions are:

- for-profit, commercial users of Commune must pay a license fee to use our software
- for-profit, commercial forks/extensions of Commune must pay a (larger) license fee to build a service with our software.

Seeing as the software code is entirely transparent, it’s pretty hard for us to artificially raise our license fee past what is reasonable. 

Any general-purpose libraries we come up with will for sure be shared as MIT-style open source.

More about PolyForm in this podcast episode:
https://player.captivate.fm/episode/75724001-14e8-48ef-8fc4-c2b1a6b04001

> Why did you choose the PolyForm NonCommercial license?

We want the software project to be self-sustainable. The best way to do that is to [charge for it somehow](https://meta.discourse.org/t/which-is-better-discourse-or-flarum/71726/7?u=erlend_sh).

We are trying out a new kind of open license because traditional open source licenses do not create a strong incentive for an exceptionally good self-hosted experience.

Self-hosting a WordPress site today isn’t any easier than it was 15 years ago. It should be like installing just another app on your computer, but it isn’t.

Standard open source licenses are locked into favoring the SaaS model, and just the SaaS model. So long as self-hosters are *not* a customer segment, they won't get the same level of service and product priority as thr SaaS users.

ahq doesn’t live in a country with welfare and safety nets, so his continued development of Commune relies on him being able to make a living with it. Therefore we plan to be charging for it as early as possible.

Commune will always be free for anyone who can't afford to pay for it, but the existence of that free service is entirely dependent on also having paying users, from the very beginning. It's either that or VC money.
