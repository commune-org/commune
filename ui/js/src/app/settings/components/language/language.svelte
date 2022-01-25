<script>
import { store } from '../../../store/store.js'
import Select from '../../../components/ui/select/select.svelte'

$: language = $store.settings.language

$: localeList = $store.locale?.settings?.language?.locale

$: languages = [
    {
        name: "en_US", 
        caption: "English, US",
        selected: language === "en_US",
        image: "🇺🇸",
        imageLabel: localeList["en_US"],
    },
    {
        name: "fr", 
        caption: "French",
        selected: language === "fr",
        image: "🇫🇷",
        imageLabel: localeList["fr"],
    },
]


function selected(e) {
    languages = e.detail
    let lang = languages.filter(x => x.selected)[0]?.name
    //localStorage.setItem("language", lang)
    //store.updateLanguage(lang)

    import(`./locale/${lang}.js`).then(d => {
        if(d) {
            localStorage.setItem("language", lang)
            store.updateLanguage(lang, d.default)
        }
    })
}

</script>


<div class="language-settings-container fl-co">

    <div class="title">
        Language
    </div>

    <div class="fl-co mt4">  
        <div class="n-t">
            select a language
        </div>


        <div class="">
            <Select items={languages} 
            on:selected={selected}
            />
        </div>
    </div>


</div>

<style>
.title {
    font-size: 1.2rem;
    font-weight: bold;
    color: var(--white);
}
.n-t {
    text-transform: uppercase;
    font-weight: bold;
    font-size: 0.7rem;
    letter-spacing: 1px;
    color: var(--text);
    margin-bottom: 0.5rem;
}
</style>

