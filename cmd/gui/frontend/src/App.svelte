<script>
  import DevicesTab from "./components/DevicesTab.svelte";
  import {GetKey} from '../wailsjs/go/main/App.js'
  import {GetIp} from '../wailsjs/go/main/App.js'

  let ip = '';
  let key = '';
  GetKey().then(result => key = result)
  GetIp().then(result => ip = result);
</script>

<main>
    <div id="meta">
        <div>
            <h3>IP: {ip}</h3>  
        </div>
        <div>
            <h3>Key: {key}</h3>
        </div>
    </div>
    <DevicesTab /> 
</main>

<style>
    #meta {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        gap: 50px;
        height: 30px;
    }

</style>
