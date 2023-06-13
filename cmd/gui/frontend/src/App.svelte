<script>
  import Notifications from 'svelte-notifications';
  import DevicesTab from './components/DevicesTab.svelte';
  import TransferList from './components/TransferList.svelte';
  import Modal from './components/Modal.svelte';
  import {GetKey} from '../wailsjs/go/main/App.js';
  import {GetIp} from '../wailsjs/go/main/App.js';
  import {AmIRunningOnMacos} from '../wailsjs/go/main/App.js';
  import {GetPendingDownloads} from '../wailsjs/go/main/App.js';

  let ip = '';
  let key = '';
  let showTransfers = false;
  let amIrunningOnMacos = false;
  let pendingDownloadCount = 0;
  GetKey().then(result => key = result);
  GetIp().then(result => ip = result);
  AmIRunningOnMacos().then(result => amIrunningOnMacos = result);
  GetPendingDownloads().then(result => pendingDownloadCount = result.length);
  setInterval(() => {
    GetPendingDownloads().then(result => pendingDownloadCount = result.length);
  }, 500);
</script>

<Notifications>
    <main style="--wails-draggable:drag">
        {#if amIrunningOnMacos}
            <div id="header">
                <p>FileBeam</p>
            </div>
        {/if}
        <div id="meta">
            <div>
                <h3>IP: <span class="mono">{ip}</span></h3>  
            </div>
            <div>
                <h3>Key: <span class="mono">{key}</span></h3>
            </div>
            <div>
                <button on:click={() => (showTransfers = true)}>
                    Transfers 
                    {#if pendingDownloadCount > 0}
                        <span class="badge">{pendingDownloadCount}</span>
                    {/if}
                </button>
            </div>
        </div>
        <DevicesTab /> 
    </main>
    
    <Modal bind:showTransfers>
        <TransferList bind:showTransfers />
    </Modal>
</Notifications>


<style>
    #meta {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        gap: 50px;
        height: 40px;
        margin-bottom: 5px;
        user-select: auto;
        -webkit-user-select: auto; /* Safari */
        -ms-user-select: auto; /* IE 10 and IE 11 */
    }

    #header {
        height: 1.5em;
        margin: -10px 0 10px;
        font-weight: bold;
        cursor: pointer;
    }
    h3 {
        font-family: 'JetBrains Mono';
    }
</style>
