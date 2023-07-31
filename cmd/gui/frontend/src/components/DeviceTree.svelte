<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscSync from "svelte-icons-pack/vsc/VscSync";
    import {GetDevices} from '../../wailsjs/go/main/App.js'
    import DeviceListItem from './DeviceListItem.svelte';

    export let onFileUpload;
    let devices = [];
    function refreshDeviceList() {
        GetDevices().then(result => {
                let newDevices = []
                result.forEach(d => {
                    newDevices.push({
                        host: d,
                        onFileUpload: onFileUpload,
                        component: DeviceListItem
                    });
            });
            devices = newDevices;
        })
    }
    refreshDeviceList()
</script>

<div id="main" class="column is-one-quarter">
    <h3>Devices on network <span><button class="refresh-btn" on:click={refreshDeviceList}><Icon src={VscSync} /></button></span></h3>
    <div id="device-tree">
        {#each devices as device}
            <svelte:component this={device.component} {...device} />
        {/each}
    </div>
</div>

<style>
    #main {
        display: flex;
        flex-direction: column;
        padding: 5px;
        width: 356px;
    }

    .refresh-btn {
        margin-left: 20px;
        width: 30px;
        height: 30px;
        border-radius: 50%;
        cursor: pointer;
    }

    #device-tree {
        height: calc(100vh - 129px);
        width: 340px;
        overflow-y: scroll;
    }
</style>