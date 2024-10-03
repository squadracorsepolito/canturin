<script lang="ts">
	type Props = {
		label: string;
	};

	let { label }: Props = $props();
</script>

<button type="submit" class="btn btn-primary">
	{label}
</button>
