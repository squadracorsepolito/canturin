<script lang="ts">
	type Props = {
		checked: boolean;
	};

	let { checked }: Props = $props();
</script>

<input type="checkbox" class="toggle toggle-primary" {checked} readonly />
