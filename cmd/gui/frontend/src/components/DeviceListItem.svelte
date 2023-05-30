<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import FiShare from "svelte-icons-pack/fi/FiShare";
	export let host;
</script>

<style>
    div {
        height: 35px;
        width: 360px;
        padding: 5px;
        background-color: rgb(247, 247, 247);
		color: #3f3f3f;
        border: 2px solid rgb(108, 108, 108);
        display: flex;
        align-items: center;
        margin: 5px;
	}
    div input, div p {
        margin-left: 10px;
        margin-right: 10px;
    }

    div input {
        width: 150px;
        height: 30px;
        font-size: 1em;
    }
</style>

<div>
    <p>{host}</p>
    <input type="text" placeholder="key">
    <button>
        <Icon src="{FiShare}" />
    </button>
</div>