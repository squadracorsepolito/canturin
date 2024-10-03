<script lang="ts">
	type Props = {
		label: String;
		onclick: () => void;
	};

	let { label, onclick }: Props = $props();
</script>

<button {onclick} class="btn">
	{label}
</button>
