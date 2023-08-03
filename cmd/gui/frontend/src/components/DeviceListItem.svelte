<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscCloudUpload from "svelte-icons-pack/vsc/VscCloudUpload";
    export let host;
    export let onFileUpload;
    let key;

    function addDashes(event) {
        if (event.key !== 'Backspace') {
            if (key.length === 4 || key.length === 9) {
                key += '-';
            }
        }
    } 
</script>

<style>
    div {
        height: 80px;
        width: 320px;
        padding: 5px;
        display: block;
        align-items: center;
        margin: 5px;
        overflow: hidden;
	}
    div input, div p {
        margin-left: 10px;
        margin-right: 10px;
    }

    div input {
        width: 200px;
        height: 30px;
        font-size: 1em;
        font-family: monospace;
        border: 1px #3e3e3e solid;
        padding: 0 5px;
        text-align: center;
    }

    .upload-btn {
        background-color: #3c6db2;
        color: #fff;
        font-weight: 800;
        height: 35px;
        width: 35px;
        border-radius: 10%;
    }
</style>

<div class="box">
    <p>{host}</p>
    <input type="text" placeholder="key" bind:value={key} on:keyup={addDashes} maxlength="16">
    <button class="upload-btn" on:click={onFileUpload(host, key)}>
        <Icon src="{VscCloudUpload}" />
    </button>
</div>