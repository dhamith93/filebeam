<script>
    import {getNotificationsContext} from 'svelte-notifications';
    import {GetHomeDir} from '../../wailsjs/go/main/App.js'
    import {AddToQueue} from '../../wailsjs/go/main/App.js'
    import FileTree from './FileTree.svelte';
    import DeviceTree from './DeviceTree.svelte';
    const {addNotification} = getNotificationsContext();
    let homeDir = '';
    let selected = [];
    let fileTree;
    GetHomeDir().then(result => homeDir = result);

    function handleFileUpload(host, key) {
        if (key.trim() === '') {
            addNotification({
                text: 'Key is empty!',
                position: 'bottom-center',
                type: 'error',
                removeAfter: 5000
            });
            return
        }
        AddToQueue(selected, host, key).then(() => {
            fileTree.clearSelection();
            selected = [];
            addNotification({
                text: 'Selected file(s) added to transfer queue',
                position: 'bottom-center',
                type: 'success',
                removeAfter: 5000
            });
        }).catch(e => {
            addNotification({
                text: e,
                position: 'bottom-center',
                type: 'error',
                removeAfter: 5000
            });
        });
    }

    function handleFileSelect(item) {
        if (selected.includes(item)) {
            selected.splice(selected.indexOf(item), 1);
        } else {
            selected.push(item);
        }
    }
    function handlePathRefresh() {
        selected = [];
    }
</script>
<div id="devices-main" class="columns">
    <FileTree path="{homeDir}" onFileSelect={handleFileSelect} onRefreshPath={handlePathRefresh} bind:this={fileTree}/>
    <DeviceTree onFileUpload={handleFileUpload} />
</div>

<style>
    #devices-main {
        display: flex;
        justify-content: center;
        padding-left: 25px;
        padding-right: 25px;
    }
</style>