<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscCloudUpload from "svelte-icons-pack/vsc/VscCloudUpload";
	export let host;
    export let onFileUpload;
    let key;
</script>

<style>
    div {
        height: 35px;
        width: 300px;
        padding: 5px;
        background-color: rgb(247, 247, 247);
		color: #3f3f3f;
        border: 2px solid rgb(108, 108, 108);
        display: flex;
        align-items: center;
        margin: 5px;
        overflow: hidden;
	}
    div input, div p {
        margin-left: 10px;
        margin-right: 10px;
    }

    div input {
        width: 80px;
        height: 30px;
        font-size: 1em;
        font-family: monospace;
        border: 1px #3e3e3e solid;
        padding: 0 5px;
    }

    .upload-btn {
        background-color: #3c6db2;
        color: #fff;
        font-size: 2em;
        font-weight: 800;
        height: 35px;
        width: 35px;
        border-radius: 10%;
    }
</style>

<div>
    <p>{host}</p>
    <input type="text" placeholder="key" bind:value={key}>
    <button class="upload-btn" on:click={onFileUpload(host, key)}>
        <Icon src="{VscCloudUpload}" />
    </button>
</div>