<script>
    import Icon from 'svelte-icons-pack/Icon.svelte';
    import VscClose from "svelte-icons-pack/vsc/VscClose";

	export let showTransfers; // boolean

	let dialog; // HTMLDialogElement

	$: if (dialog && showTransfers) dialog.showModal();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog
	bind:this={dialog}
	on:close={() => (showTransfers = false)}
	on:click|self={() => dialog.close()}
>
    <div id="header">
        <button on:click={() => dialog.close()}><Icon src="{VscClose}" /></button>
    </div>
    <div on:click|stopPropagation>
        <slot />
    </div>
</dialog>

<style>
	dialog {
        background-color: rgba(27, 38, 54, 1);
		width: 600px;
        height: 570px;
		border-radius: 0.2em;
		border: none;
		padding: 0;
        color: #fff;
	}
	dialog::backdrop {
		background: rgba(0, 0, 0, 0.5);
	}
	dialog > div {
		padding: 1em;
	}
	dialog[open] {
		animation: zoom 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
	}
	@keyframes zoom {
		from {
			transform: scale(0.95);
		}
		to {
			transform: scale(1);
		}
	}
	dialog[open]::backdrop {
		animation: fade 0.2s ease-out;
	}
	@keyframes fade {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

    #header {
        position: fixed;
        height: 29px;
        background-color: rgba(137, 172, 243, 0.2);
        padding: 5px;
    }

    button {
        height: 25px;
        color: #fff;
        font-weight: bold;
        cursor: pointer;
        transition: all 0.3s ease;
        background: #d90429;
        float: right;
        margin-top: 3px;
    }
    button:hover {
        background: #e5e5e5;
        color: #d90429
    }
</style>
