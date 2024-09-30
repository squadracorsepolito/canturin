<script lang="ts">
	type Props = {
		label: string;
		value: string | number;
	};

	let { label, value }: Props = $props();
</script>

<div class="flex flex-col items-center border">
	<div class="opacity-75">{label}</div>
	<div class="text-3xl font-extrabold">{value}</div>
</div>
