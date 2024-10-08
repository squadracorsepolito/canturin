<script lang="ts">
	type Props = {
		name: string;
		checked: boolean;
	};

	let { name, checked = $bindable() }: Props = $props();
</script>

<label>
	<input type="checkbox" {name} class="toggle toggle-primary" bind:checked />
</label>
