<script>
  import DevicesTab from "./components/DevicesTab.svelte";
  import TransferList from "./components/TransferList.svelte";
  import Modal from "./components/Modal.svelte";
  import {GetKey} from '../wailsjs/go/main/App.js'
  import {GetIp} from '../wailsjs/go/main/App.js'

  let ip = '';
  let key = '';
  let showTransfers = false;
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
        <div>
            <button on:click={() => (showTransfers = true)}>Transfers</button>
        </div>
    </div>
    <DevicesTab /> 
</main>

<Modal bind:showTransfers>
    <TransferList />
</Modal>

<style>
    #meta {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        gap: 50px;
        height: 40px;
    }

</style>
