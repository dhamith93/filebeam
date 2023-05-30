<script>
  import DevicesTab from "./components/DevicesTab.svelte";
  import {GetKey} from '../wailsjs/go/main/App.js'

  let key = '';
  GetKey().then(result => key += result)
</script>

<main>
    {key}
    <DevicesTab /> 
</main>

<style>

</style>
