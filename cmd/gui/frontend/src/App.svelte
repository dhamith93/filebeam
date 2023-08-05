<script>
  import {notifications} from './notifications.js'
  import Toast from './components/Toast.svelte'
  import TransferList from './components/TransferList.svelte';
  import SelectedFileList from './components/SelectedFileList.svelte';
  import FileTree from './components/FileTree.svelte';
  import DeviceTree from './components/DeviceTree.svelte';
  import {GetHomeDir} from '../wailsjs/go/main/App.js'
  import {AddToQueue} from '../wailsjs/go/main/App.js'
  import {GetKey} from '../wailsjs/go/main/App.js';
  import {GetIp} from '../wailsjs/go/main/App.js';
  import {AmIRunningOnMacos} from '../wailsjs/go/main/App.js';
  import {GetPendingDownloads} from '../wailsjs/go/main/App.js';

  let ip = '';
  let key = '';
  let selectedFilesModal;
  let transfersModal;
  let homeDir = '';
  let fileTree;
  let selected = [];
  let amIrunningOnMacos = false;
  let pendingDownloadCount = 0;

  GetKey().then(result => key = result);
  GetIp().then(result => ip = result);
  AmIRunningOnMacos().then(result => amIrunningOnMacos = result);
  GetPendingDownloads().then(result => pendingDownloadCount = result.length);
  setInterval(() => {
    GetPendingDownloads().then(result => pendingDownloadCount = result.length);
  }, 500);

  function displayModal(modal) {
    modal.classList.add('is-active');
  }

  GetHomeDir().then(result => homeDir = result);

  function handleFileUpload(host, key) {
    if (key.trim() === '') {
        notifications.danger('Key is empty!', 3000);
        return
    }
    AddToQueue(selected, host, key).then(() => {
        fileTree.clearSelection();
        selected = [];
        notifications.success('Selected file(s) added to transfer queue', 3000);
    }).catch(e => {
        notifications.danger(e, 3000);
    });
  }

  function handleFileSelect(item) {
    if (!selected.includes(item)) {
        notifications.success(`${item.Name} selected`, 3000)
        selected.push(item);
    }
    selected = selected;
  }
  function removeSelectedFile(file) {
    if (selected.includes(file)) {
        selected.splice(selected.indexOf(file), 1);
    }
    selected = selected;
  }
</script>

<div>
    <main style="--wails-draggable:drag">
        {#if amIrunningOnMacos}
            <div id="header">
                <p>FileBeam</p>
            </div>
        {/if}
        <div id="meta">
            <div>
                <h3>IP: {ip}</h3>  
            </div>
            <div>
                <h3>Key: {key}</h3>
            </div>
            <div>
                <button on:click={() => displayModal(transfersModal)} class="button is-info">
                    Transfers
                    <span class="badge">{pendingDownloadCount}</span>
                </button>
            </div>
            <div>
                <button on:click={() => displayModal(selectedFilesModal)} class="button is-primary">
                    Selected Files
                    <span class="badge">{selected.length}</span>
                </button>
            </div>
        </div>
        <div id="devices-main" class="columns">
            <FileTree path="{homeDir}" onFileSelect={handleFileSelect} bind:this={fileTree}/>
            <DeviceTree onFileUpload={handleFileUpload} />
        </div>
    </main>

    <div bind:this={transfersModal} class="modal">
        <div class="modal-background" on:click={() => transfersModal.classList.remove('is-active')}></div>
        <div class="modal-content">
            <TransferList />
        </div>
        <button on:click={() => transfersModal.classList.remove('is-active')} class="modal-close is-large" aria-label="close"></button>
    </div>

    <div bind:this={selectedFilesModal} class="modal">
        <div class="modal-background" on:click={() => selectedFilesModal.classList.remove('is-active')}></div>
        <div class="modal-content">
            <SelectedFileList removeSelectedFile={removeSelectedFile} selected={selected} />
        </div>
        <button on:click={() => selectedFilesModal.classList.remove('is-active')} class="modal-close is-large" aria-label="close"></button>
    </div>
</div>

<Toast />

{#if !amIrunningOnMacos}
<style>
    #meta {
        margin-top: 25px; /*WINDOWS*/
    }

    #devices-main {
        padding-left: 50px; /*WINDOWS*/
        padding-right: 50px; /*WINDOWS*/
    }
</style>
{/if}